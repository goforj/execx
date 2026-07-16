package execx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

type stage struct {
	cmd         *exec.Cmd
	def         *Cmd
	stdoutBuf   bytes.Buffer
	stderrBuf   bytes.Buffer
	combinedBuf synchronizedBuffer
	startErr    error
	setupErr    error
	waitErr     error
	outputErr   error
	startTime   time.Time
	pipeReader  *io.PipeReader
	pipeWriter  *io.PipeWriter
	ptyMaster   *os.File
	ptySlave    *os.File
	ptyWriter   io.Writer
	ptyDone     chan error
	flushOutput []func()
}

type pipeline struct {
	stages       []*stage
	withCombined bool
	outputMu     sync.Mutex
	startErr     error
}

// newPipeline materializes fresh exec.Cmd values so a configured command can be inspected before execution.
func (c *Cmd) newPipeline(withCombined bool, shadow *shadowContext) *pipeline {
	stages := c.pipelineStages()
	pipe := &pipeline{stages: stages, withCombined: withCombined}
	for _, stage := range stages {
		stage.startTime = time.Now()
		stage.cmd = stage.def.execCmd()
		if stage.def.rootCmd().usePTY {
			master, slave, err := openPTYFunc()
			if err != nil {
				stage.setupErr = err
				continue
			}
			stage.ptyMaster = master
			stage.ptySlave = slave
			var flush func()
			stage.ptyWriter, flush = stage.def.ptyWriterWithFlush(&stage.stdoutBuf, withCombined, &stage.combinedBuf, shadow)
			stage.addFlusher(flush)
			stage.cmd.Stdout = slave
			stage.cmd.Stderr = slave
		} else {
			stdoutWriter, stdoutFlush := stage.def.stdoutWriterWithFlush(&stage.stdoutBuf, withCombined, &stage.combinedBuf, shadow)
			stderrWriter, stderrFlush := stage.def.stderrWriterWithFlush(&stage.stderrBuf, withCombined, &stage.combinedBuf, shadow)
			stage.addFlusher(stdoutFlush)
			stage.addFlusher(stderrFlush)
			stage.cmd.Stdout = &synchronizedWriter{mu: &pipe.outputMu, writer: stdoutWriter, stage: stage}
			stage.cmd.Stderr = &synchronizedWriter{mu: &pipe.outputMu, writer: stderrWriter, stage: stage}
		}
	}

	for i := range stages {
		if i == 0 {
			stages[i].cmd.Stdin = stages[i].def.stdin
			continue
		}
		reader, writer := io.Pipe()
		stages[i-1].pipeWriter = writer
		stages[i].pipeReader = reader
		stages[i].cmd.Stdin = reader
		stages[i-1].cmd.Stdout = io.MultiWriter(stages[i-1].cmd.Stdout, writer)
	}

	return pipe
}

// start launches every stage and aborts already-started work if the pipeline cannot be fully constructed.
func (p *pipeline) start() {
	for i, stg := range p.stages {
		if stg.setupErr != nil {
			stg.startErr = stg.setupErr
			p.abortStart(i, stg.setupErr)
			break
		}
		stg.startErr = stg.cmd.Start()
		if stg.startErr != nil {
			if stg.ptyMaster != nil {
				_ = stg.ptyMaster.Close()
			}
			if stg.ptySlave != nil {
				_ = stg.ptySlave.Close()
			}
			p.abortStart(i, stg.startErr)
			break
		}
		if stg.ptyMaster != nil {
			stg.ptyDone = make(chan error, 1)
			go func(st *stage) {
				_, err := io.Copy(st.ptyWriter, st.ptyMaster)
				if err != nil {
					st.ptyDone <- err
				} else {
					st.ptyDone <- nil
				}
				_ = st.ptyMaster.Close()
			}(stg)
			_ = stg.ptySlave.Close()
		}
	}
}

// abortStart releases pipeline pipes and processes because a partial pipeline cannot make progress safely.
func (p *pipeline) abortStart(failed int, cause error) {
	p.startErr = cause
	for i := failed + 1; i < len(p.stages); i++ {
		p.stages[i].startErr = cause
	}
	for _, stg := range p.stages {
		if stg.pipeReader != nil {
			_ = stg.pipeReader.CloseWithError(cause)
		}
		if stg.pipeWriter != nil {
			_ = stg.pipeWriter.CloseWithError(cause)
		}
	}
	for i := 0; i < failed; i++ {
		if proc := p.stages[i].cmd.Process; proc != nil {
			_ = proc.Kill()
		}
	}
}

// wait reaps each process and closes every in-memory pipe once its producer or consumer is done.
func (p *pipeline) wait() {
	for i := range p.stages {
		if p.stages[i].startErr != nil {
			if p.stages[i].pipeWriter != nil {
				_ = p.stages[i].pipeWriter.Close()
			}
			p.stages[i].flush()
			continue
		}
		p.stages[i].waitErr = p.stages[i].cmd.Wait()
		if p.stages[i].pipeWriter != nil {
			_ = p.stages[i].pipeWriter.Close()
		}
		if p.stages[i].ptyDone != nil {
			if err := <-p.stages[i].ptyDone; err != nil {
				p.stages[i].outputErr = err
			}
		}
		if p.stages[i].pipeReader != nil {
			_ = p.stages[i].pipeReader.Close()
		}
		p.stages[i].flush()
	}
}

// results snapshots stage outcomes after all output and process state has settled.
func (p *pipeline) results() []Result {
	results := make([]Result, 0, len(p.stages))
	for _, stage := range p.stages {
		results = append(results, stage.result())
	}
	return results
}

// primaryResult applies the selected pipeline policy without changing per-stage results.
func (p *pipeline) primaryResult(mode pipeMode) (Result, string) {
	results := p.results()
	primaryIndex := len(results) - 1
	if p.startErr != nil {
		for i, res := range results {
			if res.Err != nil {
				primaryIndex = i
				break
			}
		}
	} else if mode == pipeStrict {
		for i, res := range results {
			if res.ExitCode != 0 || res.Err != nil {
				primaryIndex = i
				break
			}
		}
	}

	primary := results[primaryIndex]
	if mode == pipeBestEffort && primary.Err == nil {
		for _, res := range results {
			if res.Err != nil {
				primary.Err = res.Err
				break
			}
		}
	}

	combined := ""
	if p.withCombined {
		combined = p.stages[primaryIndex].combinedBuf.String()
	}
	return primary, combined
}

// result translates os/exec state while preserving execx's non-zero-exit-is-data contract.
func (s *stage) result() Result {
	res := Result{
		Stdout:   s.stdoutBuf.String(),
		Stderr:   s.stderrBuf.String(),
		ExitCode: -1,
		Duration: time.Since(s.startTime),
	}
	if s.startErr != nil {
		res.Err = ErrExec{
			Err:      s.startErr,
			ExitCode: -1,
			Stderr:   res.Stderr,
		}
		return res
	}
	if s.waitErr != nil {
		if errors.Is(s.waitErr, context.Canceled) || errors.Is(s.waitErr, context.DeadlineExceeded) {
			res.Err = s.waitErr
		}
		if res.Err == nil && s.def.ctx != nil && s.def.ctx.Err() != nil {
			res.Err = s.def.ctx.Err()
		}
	}
	if s.cmd.ProcessState != nil {
		res.ExitCode = s.cmd.ProcessState.ExitCode()
		res.signal = signalFromState(s.cmd.ProcessState)
	}
	if res.Err == nil && s.outputErr != nil {
		res.Err = ErrExec{Err: s.outputErr, ExitCode: res.ExitCode, Signal: res.signal, Stderr: res.Stderr}
	}
	if res.Err == nil && s.waitErr != nil {
		var exitErr *exec.ExitError
		if !errors.As(s.waitErr, &exitErr) {
			res.Err = ErrExec{Err: s.waitErr, ExitCode: res.ExitCode, Signal: res.signal, Stderr: res.Stderr}
		}
	}
	return res
}

// pipelineStages walks from the root because fluent calls may execute from any stage in a chain.
func (c *Cmd) pipelineStages() []*stage {
	root := c.rootCmd()
	stages := []*stage{}
	for current := root; current != nil; current = current.next {
		stages = append(stages, &stage{def: current})
	}
	return stages
}

// addFlusher retains only callbacks that have buffered-line state to drain.
func (s *stage) addFlusher(flush func()) {
	if flush != nil {
		s.flushOutput = append(s.flushOutput, flush)
	}
}

// flush delivers unterminated callback lines only after their output stream is closed.
func (s *stage) flush() {
	for _, flush := range s.flushOutput {
		flush()
	}
	s.flushOutput = nil
}

// synchronizedWriter serializes callbacks and caller-provided writers across an entire pipeline.
type synchronizedWriter struct {
	mu     *sync.Mutex
	writer io.Writer
	stage  *stage
}

// Write preserves stream chunks while preventing concurrent callback or writer invocation.
func (w *synchronizedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, err := w.writer.Write(p)
	if err != nil && w.stage.outputErr == nil {
		w.stage.outputErr = err
	}
	return n, err
}

// synchronizedBuffer retains the relative order of concurrently arriving stdout and stderr chunks.
type synchronizedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write appends one complete stream chunk under the combined-output lock.
func (b *synchronizedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// String returns a detached combined-output string after execution has completed.
func (b *synchronizedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
