package execx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type envMode int

type pipeMode int

const (
	envInherit envMode = iota
	envOnly
	envAppend
)

const (
	pipeStrict pipeMode = iota
	pipeBestEffort
)

// Command constructs a new command without executing it.
func Command(name string, args ...string) *Cmd {
	cmd := &Cmd{
		name:     name,
		args:     append([]string{}, args...),
		envMode:  envInherit,
		pipeMode: pipeStrict,
	}
	cmd.root = cmd
	return cmd
}

// Cmd represents a single command invocation or a pipeline stage.
type Cmd struct {
	name string
	args []string

	env     map[string]string
	envMode envMode
	ctx     context.Context
	cancel  context.CancelFunc
	dir     string

	stdin io.Reader

	onStdout func(string)
	onStderr func(string)
	stdoutW  io.Writer
	stderrW  io.Writer

	sysProcAttr *syscall.SysProcAttr

	next     *Cmd
	root     *Cmd
	pipeMode pipeMode
}

// Arg appends arguments to the command.
func (c *Cmd) Arg(values ...any) *Cmd {
	for _, value := range values {
		switch v := value.(type) {
		case string:
			c.args = append(c.args, v)
		case []string:
			c.args = append(c.args, v...)
		case map[string]string:
			keys := make([]string, 0, len(v))
			for key := range v {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				c.args = append(c.args, key, v[key])
			}
		default:
			c.args = append(c.args, fmt.Sprint(v))
		}
	}
	return c
}

// Env adds environment variables to the command.
func (c *Cmd) Env(values ...any) *Cmd {
	if c.env == nil {
		c.env = map[string]string{}
	}
	for _, value := range values {
		switch v := value.(type) {
		case string:
			key, val, _ := strings.Cut(v, "=")
			c.env[key] = val
		case []string:
			for _, entry := range v {
				key, val, _ := strings.Cut(entry, "=")
				c.env[key] = val
			}
		case map[string]string:
			for key, val := range v {
				c.env[key] = val
			}
		default:
			key, val, _ := strings.Cut(fmt.Sprint(v), "=")
			c.env[key] = val
		}
	}
	return c
}

// EnvInherit restores default environment inheritance.
func (c *Cmd) EnvInherit() *Cmd {
	c.envMode = envInherit
	return c
}

// EnvOnly ignores the parent environment.
func (c *Cmd) EnvOnly(values map[string]string) *Cmd {
	c.envMode = envOnly
	c.env = map[string]string{}
	for key, val := range values {
		c.env[key] = val
	}
	return c
}

// EnvAppend merges variables into the inherited environment.
func (c *Cmd) EnvAppend(values map[string]string) *Cmd {
	c.envMode = envAppend
	if c.env == nil {
		c.env = map[string]string{}
	}
	for key, val := range values {
		c.env[key] = val
	}
	return c
}

// Dir sets the working directory.
func (c *Cmd) Dir(path string) *Cmd {
	c.dir = path
	return c
}

// WithContext binds the command to a context.
func (c *Cmd) WithContext(ctx context.Context) *Cmd {
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
	c.ctx = ctx
	return c
}

// WithTimeout binds the command to a timeout.
func (c *Cmd) WithTimeout(d time.Duration) *Cmd {
	parent := c.ctx
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
		parent = nil
	}
	if parent == nil || parent.Err() != nil {
		parent = context.Background()
	}
	c.ctx, c.cancel = context.WithTimeout(parent, d)
	return c
}

// WithDeadline binds the command to a deadline.
func (c *Cmd) WithDeadline(t time.Time) *Cmd {
	parent := c.ctx
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
		parent = nil
	}
	if parent == nil || parent.Err() != nil {
		parent = context.Background()
	}
	c.ctx, c.cancel = context.WithDeadline(parent, t)
	return c
}

// StdinString sets stdin from a string.
func (c *Cmd) StdinString(input string) *Cmd {
	c.stdin = strings.NewReader(input)
	return c
}

// StdinBytes sets stdin from bytes.
func (c *Cmd) StdinBytes(input []byte) *Cmd {
	c.stdin = bytes.NewReader(input)
	return c
}

// StdinReader sets stdin from an io.Reader.
func (c *Cmd) StdinReader(reader io.Reader) *Cmd {
	c.stdin = reader
	return c
}

// StdinFile sets stdin from a file.
func (c *Cmd) StdinFile(file *os.File) *Cmd {
	c.stdin = file
	return c
}

// OnStdout registers a line callback for stdout.
func (c *Cmd) OnStdout(fn func(string)) *Cmd {
	c.onStdout = fn
	return c
}

// OnStderr registers a line callback for stderr.
func (c *Cmd) OnStderr(fn func(string)) *Cmd {
	c.onStderr = fn
	return c
}

// StdoutWriter sets a raw writer for stdout.
func (c *Cmd) StdoutWriter(w io.Writer) *Cmd {
	c.stdoutW = w
	return c
}

// StderrWriter sets a raw writer for stderr.
func (c *Cmd) StderrWriter(w io.Writer) *Cmd {
	c.stderrW = w
	return c
}

// Pipe appends a new command to the pipeline.
func (c *Cmd) Pipe(name string, args ...string) *Cmd {
	root := c.rootCmd()
	next := &Cmd{
		name:     name,
		args:     append([]string{}, args...),
		envMode:  envInherit,
		pipeMode: root.pipeMode,
		root:     root,
	}
	last := root
	for last.next != nil {
		last = last.next
	}
	last.next = next
	return next
}

// PipeStrict sets strict pipeline semantics.
func (c *Cmd) PipeStrict() *Cmd {
	c.rootCmd().pipeMode = pipeStrict
	return c
}

// PipeBestEffort sets best-effort pipeline semantics.
func (c *Cmd) PipeBestEffort() *Cmd {
	c.rootCmd().pipeMode = pipeBestEffort
	return c
}

// Args returns the argv slice used for execution.
func (c *Cmd) Args() []string {
	args := make([]string, 0, len(c.args)+1)
	args = append(args, c.name)
	args = append(args, c.args...)
	return args
}

// EnvList returns the environment list for execution.
func (c *Cmd) EnvList() []string {
	return buildEnv(c.envMode, c.env)
}

// String returns a human-readable representation of the command.
func (c *Cmd) String() string {
	parts := make([]string, 0, len(c.args)+1)
	parts = append(parts, c.name)
	for _, arg := range c.args {
		if strings.ContainsAny(arg, " \t\n\r") {
			parts = append(parts, strconv.Quote(arg))
			continue
		}
		parts = append(parts, arg)
	}
	return strings.Join(parts, " ")
}

// ShellEscaped returns a shell-escaped string for logging only.
func (c *Cmd) ShellEscaped() string {
	parts := make([]string, 0, len(c.args)+1)
	parts = append(parts, shellEscape(c.name))
	for _, arg := range c.args {
		parts = append(parts, shellEscape(arg))
	}
	return strings.Join(parts, " ")
}

// Run executes the command and returns the result.
func (c *Cmd) Run() Result {
	pipe := c.newPipeline(false)
	pipe.start()
	pipe.wait()
	result, _ := pipe.primaryResult(c.rootCmd().pipeMode)
	return result
}

// Output executes the command and returns stdout.
func (c *Cmd) Output() (string, error) {
	result := c.Run()
	return result.Stdout, result.Err
}

// OutputBytes executes the command and returns stdout bytes.
func (c *Cmd) OutputBytes() ([]byte, error) {
	result := c.Run()
	return []byte(result.Stdout), result.Err
}

// OutputTrimmed executes the command and returns trimmed stdout.
func (c *Cmd) OutputTrimmed() (string, error) {
	result := c.Run()
	return strings.TrimSpace(result.Stdout), result.Err
}

// CombinedOutput executes the command and returns stdout+stderr.
func (c *Cmd) CombinedOutput() (string, error) {
	pipe := c.newPipeline(true)
	pipe.start()
	pipe.wait()
	result, combined := pipe.primaryResult(c.rootCmd().pipeMode)
	return combined, result.Err
}

// PipelineResults executes the command and returns per-stage results.
func (c *Cmd) PipelineResults() []Result {
	pipe := c.newPipeline(false)
	pipe.start()
	pipe.wait()
	return pipe.results()
}

// Start executes the command asynchronously.
func (c *Cmd) Start() *Process {
	pipe := c.newPipeline(false)
	pipe.start()

	proc := &Process{
		pipeline: pipe,
		mode:     c.rootCmd().pipeMode,
		done:     make(chan struct{}),
	}
	go func() {
		pipe.wait()
		result, _ := pipe.primaryResult(proc.mode)
		proc.finish(result)
	}()
	return proc
}

func (c *Cmd) ctxOrBackground() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

func (c *Cmd) rootCmd() *Cmd {
	if c.root != nil {
		return c.root
	}
	return c
}

func (c *Cmd) execCmd() *exec.Cmd {
	cmd := exec.CommandContext(c.ctxOrBackground(), c.name, c.args...)
	if c.dir != "" {
		cmd.Dir = c.dir
	}
	cmd.Env = buildEnv(c.envMode, c.env)
	if c.sysProcAttr != nil {
		cmd.SysProcAttr = c.sysProcAttr
	}
	return cmd
}

func (c *Cmd) stdoutWriter(buf *bytes.Buffer, withCombined bool, combined *bytes.Buffer) io.Writer {
	writers := []io.Writer{}
	if c.stdoutW != nil {
		writers = append(writers, c.stdoutW)
	}
	writers = append(writers, buf)
	if withCombined {
		writers = append(writers, combined)
	}
	if c.onStdout != nil {
		writers = append(writers, &lineWriter{onLine: c.onStdout})
	}
	if len(writers) == 1 {
		return buf
	}
	return io.MultiWriter(writers...)
}

func (c *Cmd) stderrWriter(buf *bytes.Buffer, withCombined bool, combined *bytes.Buffer) io.Writer {
	writers := []io.Writer{}
	if c.stderrW != nil {
		writers = append(writers, c.stderrW)
	}
	writers = append(writers, buf)
	if withCombined {
		writers = append(writers, combined)
	}
	if c.onStderr != nil {
		writers = append(writers, &lineWriter{onLine: c.onStderr})
	}
	if len(writers) == 1 {
		return buf
	}
	return io.MultiWriter(writers...)
}

type lineWriter struct {
	onLine func(string)
	buf    bytes.Buffer
}

func (l *lineWriter) Write(p []byte) (int, error) {
	if l.onLine == nil {
		return len(p), nil
	}
	for _, b := range p {
		if b == '\n' {
			line := l.buf.String()
			l.buf.Reset()
			line = strings.TrimSuffix(line, "\r")
			l.onLine(line)
			continue
		}
		_ = l.buf.WriteByte(b)
	}
	return len(p), nil
}

func buildEnv(mode envMode, env map[string]string) []string {
	merged := map[string]string{}
	if mode != envOnly {
		for _, entry := range os.Environ() {
			key, val, _ := strings.Cut(entry, "=")
			merged[key] = val
		}
	}
	for key, val := range env {
		merged[key] = val
	}
	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	list := make([]string, 0, len(keys))
	for _, key := range keys {
		list = append(list, key+"="+merged[key])
	}
	return list
}

func shellEscape(arg string) string {
	if arg == "" {
		return "''"
	}
	needsQuote := strings.ContainsAny(arg, " \t\n\r'\"\\$`!")
	if !needsQuote {
		return arg
	}
	return "'" + strings.ReplaceAll(arg, "'", "'\\''") + "'"
}

// Process represents an asynchronously running command.
type Process struct {
	pipeline *pipeline
	mode     pipeMode
	done     chan struct{}
	result   Result

	resultOnce sync.Once
	mu         sync.Mutex
	killTimer  *time.Timer
}

// Wait waits for the command to complete and returns the result.
func (p *Process) Wait() Result {
	<-p.done
	return p.result
}

// KillAfter terminates the process after the given duration.
func (p *Process) KillAfter(d time.Duration) {
	p.mu.Lock()
	if p.killTimer != nil {
		p.killTimer.Stop()
	}
	p.killTimer = time.AfterFunc(d, func() {
		_ = p.Terminate()
	})
	p.mu.Unlock()
}

// Send sends a signal to the process.
func (p *Process) Send(sig os.Signal) error {
	return p.signalAll(func(proc *os.Process) error {
		return proc.Signal(sig)
	})
}

// Interrupt sends an interrupt signal to the process.
func (p *Process) Interrupt() error {
	return p.Send(os.Interrupt)
}

// Terminate kills the process immediately.
func (p *Process) Terminate() error {
	return p.signalAll(func(proc *os.Process) error {
		return proc.Kill()
	})
}

// GracefulShutdown sends a signal and escalates to kill after the timeout.
func (p *Process) GracefulShutdown(sig os.Signal, timeout time.Duration) error {
	if timeout <= 0 {
		return p.Terminate()
	}
	if err := p.Send(sig); err != nil {
		return err
	}
	select {
	case <-p.done:
		return nil
	case <-time.After(timeout):
	}
	_ = p.Terminate()
	<-p.done
	return nil
}

func (p *Process) finish(result Result) {
	p.resultOnce.Do(func() {
		p.result = result
		close(p.done)
	})
}

func (p *Process) signalAll(send func(*os.Process) error) error {
	if p == nil || p.pipeline == nil {
		return errors.New("process not started")
	}
	var firstErr error
	count := 0
	for _, stage := range p.pipeline.stages {
		if stage == nil || stage.cmd == nil || stage.cmd.Process == nil {
			continue
		}
		count++
		if err := send(stage.cmd.Process); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if count == 0 && firstErr == nil {
		return errors.New("process not started")
	}
	return firstErr
}
