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
// @group Construction
//
// Example: command
//
//	cmd := execx.Command("printf", "hello")
//	out, _ := cmd.Output()
//	fmt.Print(out)
//	// hello
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

	shadowPrint     bool
	shadowPrefix    string
	shadowFormatter func(ShadowEvent) string
	shadowMask      func(string) string
}

// Arg appends arguments to the command.
// @group Arguments
//
// Example: add args
//
//	cmd := execx.Command("printf").Arg("hello")
//	out, _ := cmd.Output()
//	fmt.Print(out)
//	// hello
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
// @group Environment
//
// Example: set env
//
//	cmd := execx.Command("go", "env", "GOOS").Env("MODE=prod")
//	fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "MODE=prod"))
//	// #bool true
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
// @group Environment
//
// Example: inherit env
//
//	cmd := execx.Command("go", "env", "GOOS").EnvInherit()
//	fmt.Println(len(cmd.EnvList()) > 0)
//	// #bool true
func (c *Cmd) EnvInherit() *Cmd {
	c.envMode = envInherit
	return c
}

// EnvOnly ignores the parent environment.
// @group Environment
//
// Example: replace env
//
//	cmd := execx.Command("go", "env", "GOOS").EnvOnly(map[string]string{"A": "1"})
//	fmt.Println(strings.Join(cmd.EnvList(), ","))
//	// #string A=1
func (c *Cmd) EnvOnly(values map[string]string) *Cmd {
	c.envMode = envOnly
	c.env = map[string]string{}
	for key, val := range values {
		c.env[key] = val
	}
	return c
}

// EnvAppend merges variables into the inherited environment.
// @group Environment
//
// Example: append env
//
//	cmd := execx.Command("go", "env", "GOOS").EnvAppend(map[string]string{"A": "1"})
//	fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "A=1"))
//	// #bool true
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
// @group WorkingDir
//
// Example: change dir
//
//	dir := os.TempDir()
//	out, _ := execx.Command("pwd").
//		Dir(dir).
//		OutputTrimmed()
//	fmt.Println(out == dir)
//	// #bool true
func (c *Cmd) Dir(path string) *Cmd {
	c.dir = path
	return c
}

// WithContext binds the command to a context.
// @group Context
//
// Example: with context
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	res, _ := execx.Command("go", "env", "GOOS").WithContext(ctx).Run()
//	fmt.Println(res.ExitCode == 0)
//	// #bool true
func (c *Cmd) WithContext(ctx context.Context) *Cmd {
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
	c.ctx = ctx
	return c
}

// WithTimeout binds the command to a timeout.
// @group Context
//
// Example: with timeout
//
//	res, _ := execx.Command("go", "env", "GOOS").WithTimeout(2 * time.Second).Run()
//	fmt.Println(res.ExitCode == 0)
//	// #bool true
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
// @group Context
//
// Example: with deadline
//
//	res, _ := execx.Command("go", "env", "GOOS").WithDeadline(time.Now().Add(2 * time.Second)).Run()
//	fmt.Println(res.ExitCode == 0)
//	// #bool true
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
// @group Input
//
// Example: stdin string
//
//	out, _ := execx.Command("cat").
//		StdinString("hi").
//		Output()
//	fmt.Println(out)
//	// #string hi
func (c *Cmd) StdinString(input string) *Cmd {
	c.stdin = strings.NewReader(input)
	return c
}

// StdinBytes sets stdin from bytes.
// @group Input
//
// Example: stdin bytes
//
//	out, _ := execx.Command("cat").
//		StdinBytes([]byte("hi")).
//		Output()
//	fmt.Println(out)
//	// #string hi
func (c *Cmd) StdinBytes(input []byte) *Cmd {
	c.stdin = bytes.NewReader(input)
	return c
}

// StdinReader sets stdin from an io.Reader.
// @group Input
//
// Example: stdin reader
//
//	out, _ := execx.Command("cat").
//		StdinReader(strings.NewReader("hi")).
//		Output()
//	fmt.Println(out)
//	// #string hi
func (c *Cmd) StdinReader(reader io.Reader) *Cmd {
	c.stdin = reader
	return c
}

// StdinFile sets stdin from a file.
// @group Input
//
// Example: stdin file
//
//	file, _ := os.CreateTemp("", "execx-stdin")
//	_, _ = file.WriteString("hi")
//	_, _ = file.Seek(0, 0)
//	out, _ := execx.Command("cat").
//		StdinFile(file).
//		Output()
//	fmt.Println(out)
//	// #string hi
func (c *Cmd) StdinFile(file *os.File) *Cmd {
	c.stdin = file
	return c
}

// OnStdout registers a line callback for stdout.
// @group Streaming
//
// Example: stdout lines
//
//	_, _ = execx.Command("printf", "hi\n").
//		OnStdout(func(line string) { fmt.Println(line) }).
//		Run()
//	// hi
func (c *Cmd) OnStdout(fn func(string)) *Cmd {
	c.onStdout = fn
	return c
}

// OnStderr registers a line callback for stderr.
// @group Streaming
//
// Example: stderr lines
//
//	_, err := execx.Command("go", "env", "-badflag").
//		OnStderr(func(line string) {
//			fmt.Println(line)
//		}).
//		Run()
//	fmt.Println(err == nil)
//	// flag provided but not defined: -badflag
//	// usage: go env [-json] [-changed] [-u] [-w] [var ...]
//	// Run 'go help env' for details.
//	// false
func (c *Cmd) OnStderr(fn func(string)) *Cmd {
	c.onStderr = fn
	return c
}

// StdoutWriter sets a raw writer for stdout.
// @group Streaming
//
// Example: stdout writer
//
//	var out strings.Builder
//	_, _ = execx.Command("printf", "hello").
//		StdoutWriter(&out).
//		Run()
//	fmt.Print(out.String())
//	// hello
func (c *Cmd) StdoutWriter(w io.Writer) *Cmd {
	c.stdoutW = w
	return c
}

// StderrWriter sets a raw writer for stderr.
// @group Streaming
//
// Example: stderr writer
//
//	var out strings.Builder
//	_, err := execx.Command("go", "env", "-badflag").
//		StderrWriter(&out).
//		Run()
//	fmt.Print(out.String())
//	fmt.Println(err == nil)
//	// flag provided but not defined: -badflag
//	// usage: go env [-json] [-changed] [-u] [-w] [var ...]
//	// Run 'go help env' for details.
//	// false
func (c *Cmd) StderrWriter(w io.Writer) *Cmd {
	c.stderrW = w
	return c
}

// Pipe appends a new command to the pipeline. Pipelines run on all platforms.
// @group Pipelining
//
// Example: pipe
//
//	out, _ := execx.Command("printf", "go").
//		Pipe("tr", "a-z", "A-Z").
//		OutputTrimmed()
//	fmt.Println(out)
//	// #string GO
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

// PipeStrict sets strict pipeline semantics (stop on first failure).
// @group Pipelining
//
// Example: strict
//
//	res, _ := execx.Command("false").
//		Pipe("printf", "ok").
//		PipeStrict().
//		Run()
//	fmt.Println(res.ExitCode != 0)
//	// #bool true
func (c *Cmd) PipeStrict() *Cmd {
	c.rootCmd().pipeMode = pipeStrict
	return c
}

// PipeBestEffort sets best-effort pipeline semantics (run all stages, surface the first error).
// @group Pipelining
//
// Example: best effort
//
//	res, _ := execx.Command("false").
//		Pipe("printf", "ok").
//		PipeBestEffort().
//		Run()
//	fmt.Print(res.Stdout)
//	// ok
func (c *Cmd) PipeBestEffort() *Cmd {
	c.rootCmd().pipeMode = pipeBestEffort
	return c
}

// Args returns the argv slice used for execution.
// @group Debugging
//
// Example: args
//
//	cmd := execx.Command("go", "env", "GOOS")
//	fmt.Println(strings.Join(cmd.Args(), " "))
//	// #string go env GOOS
func (c *Cmd) Args() []string {
	args := make([]string, 0, len(c.args)+1)
	args = append(args, c.name)
	args = append(args, c.args...)
	return args
}

// EnvList returns the environment list for execution.
// @group Environment
//
// Example: env list
//
//	cmd := execx.Command("go", "env", "GOOS").EnvOnly(map[string]string{"A": "1"})
//	fmt.Println(strings.Join(cmd.EnvList(), ","))
//	// #string A=1
func (c *Cmd) EnvList() []string {
	return buildEnv(c.envMode, c.env)
}

// String returns a human-readable representation of the command.
// @group Debugging
//
// Example: string
//
//	cmd := execx.Command("echo", "hello world", "it's")
//	fmt.Println(cmd.String())
//	// #string echo "hello world" it's
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
// @group Debugging
//
// Example: shell escaped
//
//	cmd := execx.Command("echo", "hello world", "it's")
//	fmt.Println(cmd.ShellEscaped())
//	// #string echo 'hello world' 'it'\\''s'
func (c *Cmd) ShellEscaped() string {
	parts := make([]string, 0, len(c.args)+1)
	parts = append(parts, shellEscape(c.name))
	for _, arg := range c.args {
		parts = append(parts, shellEscape(arg))
	}
	return strings.Join(parts, " ")
}

// ShadowPrint writes the shell-escaped command to stderr before and after execution.
// @group User Feedback
//
// Example: shadow print
//
//	_, _ = execx.Command("printf", "hi").ShadowPrint().Run()
//	// execx > printf hi
//	// execx > printf hi (1ms)
func (c *Cmd) ShadowPrint() *Cmd {
	c.rootCmd().shadowPrint = true
	return c
}

// ShadowPrintOff disables shadow printing for this command chain.
// @group User Feedback
//
// Example: shadow print off
//
//	_, _ = execx.Command("printf", "hi").ShadowPrint().ShadowPrintOff().Run()
func (c *Cmd) ShadowPrintOff() *Cmd {
	root := c.rootCmd()
	root.shadowPrint = false
	return c
}

// ShadowPrintPrefix sets the prefix used by ShadowPrint.
// @group User Feedback
//
// Example: shadow print prefix
//
//	_, _ = execx.Command("printf", "hi").ShadowPrintPrefix("run").Run()
//	// run > printf hi
//	// run > printf hi (1ms)
func (c *Cmd) ShadowPrintPrefix(prefix string) *Cmd {
	root := c.rootCmd()
	root.shadowPrint = true
	root.shadowPrefix = prefix
	return c
}

// ShadowPrintMask sets a command masker for ShadowPrint output.
// @group User Feedback
//
// Example: shadow print mask
//
//	mask := func(cmd string) string {
//		return strings.ReplaceAll(cmd, "secret", "***")
//	}
//	_, _ = execx.Command("printf", "secret").ShadowPrintMask(mask).Run()
//	// execx > printf ***
//	// execx > printf *** (1ms)
func (c *Cmd) ShadowPrintMask(fn func(string) string) *Cmd {
	root := c.rootCmd()
	root.shadowPrint = true
	root.shadowMask = fn
	return c
}

// ShadowPrintFormatter sets a formatter for ShadowPrint output.
// @group User Feedback
//
// Example: shadow print formatter
//
//	formatter := func(ev execx.ShadowEvent) string {
//		return fmt.Sprintf("shadow: %s %s", ev.Phase, ev.Command)
//	}
//	_, _ = execx.Command("printf", "hi").ShadowPrintFormatter(formatter).Run()
//	// shadow: before printf hi
//	// shadow: after printf hi
func (c *Cmd) ShadowPrintFormatter(fn func(ShadowEvent) string) *Cmd {
	root := c.rootCmd()
	root.shadowPrint = true
	root.shadowFormatter = fn
	return c
}

// Run executes the command and returns the result and any error.
// @group Execution
//
// Example: run
//
//	res, _ := execx.Command("go", "env", "GOOS").Run()
//	fmt.Println(res.ExitCode == 0)
//	// #bool true
func (c *Cmd) Run() (Result, error) {
	shadow := c.shadowPrintStart(false)
	pipe := c.newPipeline(false)
	pipe.start()
	pipe.wait()
	result, _ := pipe.primaryResult(c.rootCmd().pipeMode)
	shadow.finish()
	return result, result.Err
}

// Output executes the command and returns stdout and any error.
// @group Execution
//
// Example: output
//
//	out, _ := execx.Command("printf", "hello").Output()
//	fmt.Print(out)
//	// hello
func (c *Cmd) Output() (string, error) {
	result, err := c.Run()
	return result.Stdout, err
}

// OutputBytes executes the command and returns stdout bytes and any error.
// @group Execution
//
// Example: output bytes
//
//	out, _ := execx.Command("printf", "hello").OutputBytes()
//	fmt.Println(string(out))
//	// #string hello
func (c *Cmd) OutputBytes() ([]byte, error) {
	result, err := c.Run()
	return []byte(result.Stdout), err
}

// OutputTrimmed executes the command and returns trimmed stdout and any error.
// @group Execution
//
// Example: output trimmed
//
//	out, _ := execx.Command("printf", "hello\n").OutputTrimmed()
//	fmt.Println(out)
//	// #string hello
func (c *Cmd) OutputTrimmed() (string, error) {
	result, err := c.Run()
	return strings.TrimSpace(result.Stdout), err
}

// CombinedOutput executes the command and returns stdout+stderr and any error.
// @group Execution
//
// Example: combined output
//
//	out, err := execx.Command("go", "env", "-badflag").CombinedOutput()
//	fmt.Print(out)
//	fmt.Println(err == nil)
//	// flag provided but not defined: -badflag
//	// usage: go env [-json] [-changed] [-u] [-w] [var ...]
//	// Run 'go help env' for details.
//	// false
func (c *Cmd) CombinedOutput() (string, error) {
	shadow := c.shadowPrintStart(false)
	pipe := c.newPipeline(true)
	pipe.start()
	pipe.wait()
	result, combined := pipe.primaryResult(c.rootCmd().pipeMode)
	shadow.finish()
	return combined, result.Err
}

// PipelineResults executes the command and returns per-stage results and any error.
// @group Pipelining
//
// Example: pipeline results
//
//	results, _ := execx.Command("printf", "go").
//		Pipe("tr", "a-z", "A-Z").
//		PipelineResults()
//	fmt.Printf("%+v", results)
//	// [
//	//	{Stdout:go Stderr: ExitCode:0 Err:<nil> Duration:6.367208ms signal:<nil>}
//	//	{Stdout:GO Stderr: ExitCode:0 Err:<nil> Duration:4.976291ms signal:<nil>}
//	// ]
func (c *Cmd) PipelineResults() ([]Result, error) {
	shadow := c.shadowPrintStart(false)
	pipe := c.newPipeline(false)
	pipe.start()
	pipe.wait()
	results := pipe.results()
	shadow.finish()
	return results, firstResultErr(results)
}

// Start executes the command asynchronously.
// @group Execution
//
// Example: start
//
//	proc := execx.Command("go", "env", "GOOS").Start()
//	res, _ := proc.Wait()
//	fmt.Println(res.ExitCode == 0)
//	// #bool true
func (c *Cmd) Start() *Process {
	shadow := c.shadowPrintStart(true)
	pipe := c.newPipeline(false)
	pipe.start()

	proc := &Process{
		pipeline: pipe,
		mode:     c.rootCmd().pipeMode,
		done:     make(chan struct{}),
		shadow:   shadow,
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

func firstResultErr(results []Result) error {
	for _, res := range results {
		if res.Err != nil {
			return res.Err
		}
	}
	return nil
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

type shadowContext struct {
	cmd   *Cmd
	start time.Time
	async bool
}

// ShadowPhase describes whether the shadow print is before or after execution.
type ShadowPhase string

const (
	ShadowBefore ShadowPhase = "before"
	ShadowAfter  ShadowPhase = "after"
)

// ShadowEvent captures details for ShadowPrint formatting.
type ShadowEvent struct {
	Command    string
	RawCommand string
	Phase      ShadowPhase
	Duration   time.Duration
	Async      bool
}

func (c *Cmd) shadowPrintStart(async bool) *shadowContext {
	root := c.rootCmd()
	if root == nil || !root.shadowPrint {
		return nil
	}
	ctx := &shadowContext{
		cmd:   root,
		start: time.Now(),
		async: async,
	}
	shadowPrintLine(root, ShadowBefore, 0, async)
	return ctx
}

func (c *Cmd) shadowCommand() string {
	root := c.rootCmd()
	if root == nil {
		return ""
	}
	parts := []string{}
	for stage := root; stage != nil; stage = stage.next {
		parts = append(parts, stage.ShellEscaped())
	}
	return strings.Join(parts, " | ")
}

func (s *shadowContext) finish() {
	if s == nil || s.cmd == nil {
		return
	}
	duration := time.Since(s.start).Round(time.Millisecond)
	shadowPrintLine(s.cmd, ShadowAfter, duration, s.async)
}

func shadowPrintLine(cmd *Cmd, phase ShadowPhase, duration time.Duration, async bool) {
	if cmd == nil {
		return
	}
	rawCommand := cmd.shadowCommand()
	commandLine := rawCommand
	if cmd.shadowMask != nil {
		commandLine = cmd.shadowMask(commandLine)
	}
	event := ShadowEvent{
		Command:    commandLine,
		RawCommand: rawCommand,
		Phase:      phase,
		Duration:   duration,
		Async:      async,
	}
	if cmd.shadowFormatter != nil {
		line := cmd.shadowFormatter(event)
		if line != "" {
			fmt.Fprintln(os.Stderr, line)
		}
		return
	}
	prefix := cmd.shadowPrefix
	if prefix == "" {
		prefix = "execx"
	}
	timing := ""
	if phase == ShadowAfter {
		timing = " (" + duration.String()
		if async {
			timing += ", async"
		}
		timing += ")"
	} else if async {
		timing = " (async)"
	}
	fmt.Fprintf(
		os.Stderr,
		"%s%s > %s%s%s%s%s\n",
		shadowBoldHighIntensityBlack,
		prefix,
		shadowReset,
		shadowHighIntensityBlack,
		commandLine,
		shadowFadedGray,
		timing+shadowReset,
	)
}

// Process represents an asynchronously running command.
type Process struct {
	pipeline *pipeline
	mode     pipeMode
	done     chan struct{}
	result   Result

	shadow *shadowContext

	resultOnce sync.Once
	mu         sync.Mutex
	killTimer  *time.Timer
}

// Wait waits for the command to complete and returns the result and any error.
// @group Process
//
// Example: wait
//
//	proc := execx.Command("go", "env", "GOOS").Start()
//	res, _ := proc.Wait()
//	fmt.Printf("%+v", res)
//	// {Stdout:darwin
//	// Stderr: ExitCode:0 Err:<nil> Duration:1.234ms signal:<nil>}
func (p *Process) Wait() (Result, error) {
	<-p.done
	return p.result, p.result.Err
}

// KillAfter terminates the process after the given duration.
// @group Process
//
// Example: kill after
//
//	proc := execx.Command("sleep", "2").Start()
//	proc.KillAfter(100 * time.Millisecond)
//	res, _ := proc.Wait()
//	fmt.Printf("%+v", res)
//	// {Stdout: Stderr: ExitCode:-1 Err:<nil> Duration:100.456ms signal:killed}
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
// @group Process
//
// Example: send signal
//
//	proc := execx.Command("sleep", "2").Start()
//	_ = proc.Send(os.Interrupt)
//	res, _ := proc.Wait()
//	fmt.Printf("%+v", res)
//	// {Stdout: Stderr: ExitCode:-1 Err:<nil> Duration:80.123ms signal:interrupt}
func (p *Process) Send(sig os.Signal) error {
	return p.signalAll(func(proc *os.Process) error {
		return proc.Signal(sig)
	})
}

// Interrupt sends an interrupt signal to the process.
// @group Process
//
// Example: interrupt
//
//	proc := execx.Command("sleep", "2").Start()
//	_ = proc.Interrupt()
//	res, _ := proc.Wait()
//	fmt.Printf("%+v", res)
//	// {Stdout: Stderr: ExitCode:-1 Err:<nil> Duration:75.987ms signal:interrupt}
func (p *Process) Interrupt() error {
	return p.Send(os.Interrupt)
}

// Terminate kills the process immediately.
// @group Process
//
// Example: terminate
//
//	proc := execx.Command("sleep", "2").Start()
//	_ = proc.Terminate()
//	res, _ := proc.Wait()
//	fmt.Printf("%+v", res)
//	// {Stdout: Stderr: ExitCode:-1 Err:<nil> Duration:70.654ms signal:killed}
func (p *Process) Terminate() error {
	return p.signalAll(func(proc *os.Process) error {
		return proc.Kill()
	})
}

// GracefulShutdown sends a signal and escalates to kill after the timeout.
// @group Process
//
// Example: graceful shutdown
//
//	proc := execx.Command("sleep", "2").Start()
//	_ = proc.GracefulShutdown(os.Interrupt, 100*time.Millisecond)
//	res, _ := proc.Wait()
//	fmt.Println(res.IsSignal(os.Interrupt))
//	// #bool true
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
		p.shadow.finish()
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
