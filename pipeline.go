package execx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"
)

type stage struct {
	cmd         *exec.Cmd
	def         *Cmd
	stdoutBuf   bytes.Buffer
	stderrBuf   bytes.Buffer
	combinedBuf bytes.Buffer
	startErr    error
	setupErr    error
	waitErr     error
	startTime   time.Time
	pipeWriter  *io.PipeWriter
	ptyMaster   *os.File
	ptySlave    *os.File
	ptyWriter   io.Writer
	ptyDone     chan error
}

type pipeline struct {
	stages       []*stage
	withCombined bool
}

func (c *Cmd) newPipeline(withCombined bool, shadow *shadowContext) *pipeline {
	stages := c.pipelineStages()
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
			stage.ptyWriter = stage.def.ptyWriter(&stage.stdoutBuf, withCombined, &stage.combinedBuf, shadow)
			stage.cmd.Stdout = slave
			stage.cmd.Stderr = slave
		} else {
			stdoutWriter := stage.def.stdoutWriter(&stage.stdoutBuf, withCombined, &stage.combinedBuf, shadow)
			stderrWriter := stage.def.stderrWriter(&stage.stderrBuf, withCombined, &stage.combinedBuf, shadow)
			stage.cmd.Stdout = stdoutWriter
			stage.cmd.Stderr = stderrWriter
		}
	}

	for i := range stages {
		if i == 0 {
			stages[i].cmd.Stdin = stages[i].def.stdin
			continue
		}
		reader, writer := io.Pipe()
		stages[i-1].pipeWriter = writer
		stages[i].cmd.Stdin = reader
		stages[i-1].cmd.Stdout = io.MultiWriter(stages[i-1].cmd.Stdout, writer)
	}

	return &pipeline{stages: stages, withCombined: withCombined}
}

func (p *pipeline) start() {
	for i, stg := range p.stages {
		if stg.setupErr != nil {
			stg.startErr = stg.setupErr
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
			for j := i + 1; j < len(p.stages); j++ {
				p.stages[j].startErr = stg.startErr
			}
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

func (p *pipeline) wait() {
	for i := range p.stages {
		if p.stages[i].startErr != nil {
			if p.stages[i].pipeWriter != nil {
				_ = p.stages[i].pipeWriter.Close()
			}
			continue
		}
		p.stages[i].waitErr = p.stages[i].cmd.Wait()
		if p.stages[i].pipeWriter != nil {
			_ = p.stages[i].pipeWriter.Close()
		}
		if p.stages[i].ptyDone != nil {
			if err := <-p.stages[i].ptyDone; err != nil && p.stages[i].waitErr == nil {
				p.stages[i].waitErr = err
			}
		}
	}
}

func (p *pipeline) results() []Result {
	results := make([]Result, 0, len(p.stages))
	for _, stage := range p.stages {
		results = append(results, stage.result())
	}
	return results
}

func (p *pipeline) primaryResult(mode pipeMode) (Result, string) {
	results := p.results()
	primaryIndex := len(results) - 1
	if mode == pipeStrict {
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
	return res
}

func (c *Cmd) pipelineStages() []*stage {
	root := c.rootCmd()
	stages := []*stage{}
	for current := root; current != nil; current = current.next {
		stages = append(stages, &stage{def: current})
	}
	return stages
}
