package execx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestHelperProcess provides deterministic subprocess behavior for stream, signal, and exit tests.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("EXECX_TEST_HELPER") != "1" {
		return
	}
	args := os.Args
	idx := 0
	for idx < len(args) && args[idx] != "--" {
		idx++
	}
	if idx >= len(args)-1 {
		os.Exit(1)
	}
	cmd := args[idx+1]
	cmdArgs := args[idx+2:]
	switch cmd {
	case "echo":
		_, _ = io.WriteString(os.Stdout, strings.Join(cmdArgs, " "))
	case "stderr":
		_, _ = io.WriteString(os.Stderr, strings.Join(cmdArgs, " "))
	case "cat":
		_, _ = io.Copy(os.Stdout, os.Stdin)
	case "exit":
		code, _ := strconv.Atoi(cmdArgs[0])
		os.Exit(code)
	case "mix":
		_, _ = io.WriteString(os.Stdout, "a")
		time.Sleep(10 * time.Millisecond)
		_, _ = io.WriteString(os.Stderr, "b")
		time.Sleep(10 * time.Millisecond)
		_, _ = io.WriteString(os.Stdout, "c")
	case "lines":
		_, _ = io.WriteString(os.Stdout, "a\nb\n")
		_, _ = io.WriteString(os.Stderr, "c\n")
	case "env":
		_, _ = io.WriteString(os.Stdout, os.Getenv(cmdArgs[0]))
	case "sleep":
		ms, _ := strconv.Atoi(cmdArgs[0])
		time.Sleep(time.Duration(ms) * time.Millisecond)
	case "pwd":
		wd, _ := os.Getwd()
		_, _ = io.WriteString(os.Stdout, wd)
	case "signal":
		terminateHelperProcess()
		time.Sleep(50 * time.Millisecond)
	case "ignore-term":
		if runtime.GOOS == "windows" {
			os.Exit(3)
		}
		signal.Ignore(syscall.SIGTERM, os.Interrupt)
		time.Sleep(200 * time.Millisecond)
	default:
		os.Exit(1)
	}
	os.Exit(0)
}

// helperCommand routes deterministic subprocess behavior through the current test binary.
func helperCommand(args ...string) *Cmd {
	full := append([]string{"-test.run=TestHelperProcess", "--"}, args...)
	cmd := Command(os.Args[0], full...)
	cmd.Env("EXECX_TEST_HELPER=1")
	return cmd
}

// helperPipe appends a deterministic helper subprocess as one pipeline stage.
func helperPipe(cmd *Cmd, args ...string) *Cmd {
	full := append([]string{"-test.run=TestHelperProcess", "--"}, args...)
	stage := cmd.Pipe(os.Args[0], full...)
	stage.Env("EXECX_TEST_HELPER=1")
	return stage
}

type envStringer struct{}

// String provides an environment entry through Env's fmt.Stringer fallback.
func (envStringer) String() string {
	return "EXECX_ENV_VALUE=stringer"
}

// TestArgOrderAndArgs ensures heterogeneous arguments retain deterministic command-line order.
func TestArgOrderAndArgs(t *testing.T) {
	cmd := helperCommand("echo").Arg("alpha").Arg(map[string]string{"--b": "2", "--a": "1"})
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "alpha --a 1 --b 2" {
		t.Fatalf("unexpected output: %q", out)
	}
	args := cmd.Args()
	if len(args) < 1 || args[0] != os.Args[0] {
		t.Fatalf("expected argv to include executable, got %v", args)
	}
}

// TestArgVariants ensures slices and scalar values share the fluent argument path.
func TestArgVariants(t *testing.T) {
	out, err := helperCommand("echo").Arg([]string{"a", "b"}, 123).Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "a b 123" {
		t.Fatalf("unexpected output: %q", out)
	}
}

// TestEnvModes ensures inheritance, replacement, and appended overrides remain distinct policies.
func TestEnvModes(t *testing.T) {
	key := "EXECX_ENV_VALUE"
	t.Setenv(key, "base")

	out, err := helperCommand("env", key).Output()
	if err != nil || out != "base" {
		t.Fatalf("expected inherited env, got %q err=%v", out, err)
	}

	out, err = helperCommand("env", key).EnvOnly(map[string]string{key: "only", "EXECX_TEST_HELPER": "1"}).Output()
	if err != nil || out != "only" {
		t.Fatalf("expected env only, got %q err=%v", out, err)
	}

	out, err = helperCommand("env", key).EnvAppend(map[string]string{key: "append"}).Output()
	if err != nil || out != "append" {
		t.Fatalf("expected env append override, got %q err=%v", out, err)
	}

	out, err = helperCommand("env", key).EnvOnly(map[string]string{key: "only", "EXECX_TEST_HELPER": "1"}).EnvInherit().Output()
	if err != nil || out != "only" {
		t.Fatalf("expected env inherit to keep overrides, got %q err=%v", out, err)
	}
}

// TestEnvVariants ensures supported environment input forms resolve with deterministic precedence.
func TestEnvVariants(t *testing.T) {
	cmd := Command(os.Args[0], "-test.run=TestHelperProcess", "--", "env", "EXECX_ENV_VALUE").
		Env(envStringer{}).
		Env([]string{"EXECX_TEST_HELPER=1"}).
		Env(map[string]string{"EXECX_ENV_VALUE": "map"})
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "map" {
		t.Fatalf("unexpected env output: %q", out)
	}
}

// TestEnvAppendEmpty ensures append mode starts from the process environment when no explicit map exists.
func TestEnvAppendEmpty(t *testing.T) {
	cmd := Command(os.Args[0], "-test.run=TestHelperProcess", "--", "env", "EXECX_ENV_VALUE").
		EnvAppend(map[string]string{"EXECX_ENV_VALUE": "append", "EXECX_TEST_HELPER": "1"})
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "append" {
		t.Fatalf("unexpected env output: %q", out)
	}
}

// TestEnvList ensures environment inspection is sorted and therefore reproducible.
func TestEnvList(t *testing.T) {
	cmd := helperCommand("env", "NONE").EnvOnly(map[string]string{"B": "2", "A": "1", "EXECX_TEST_HELPER": "1"})
	list := cmd.EnvList()
	if strings.Join(list, ",") != "A=1,B=2,EXECX_TEST_HELPER=1" {
		t.Fatalf("unexpected env list: %v", list)
	}
}

// TestStdinHelpers ensures every supported input source reaches the child process unchanged.
func TestStdinHelpers(t *testing.T) {
	cases := []struct {
		name string
		cmd  func() *Cmd
	}{
		{
			name: "string",
			cmd: func() *Cmd {
				return helperCommand("cat").StdinString("hello")
			},
		},
		{
			name: "bytes",
			cmd: func() *Cmd {
				return helperCommand("cat").StdinBytes([]byte("hello"))
			},
		},
		{
			name: "reader",
			cmd: func() *Cmd {
				return helperCommand("cat").StdinReader(strings.NewReader("hello"))
			},
		},
		{
			name: "file",
			cmd: func() *Cmd {
				file, err := os.CreateTemp(t.TempDir(), "stdin")
				if err != nil {
					t.Fatalf("temp file: %v", err)
				}
				if _, err := file.WriteString("hello"); err != nil {
					t.Fatalf("write temp: %v", err)
				}
				if _, err := file.Seek(0, io.SeekStart); err != nil {
					t.Fatalf("seek temp: %v", err)
				}
				return helperCommand("cat").StdinFile(file)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.cmd().Output()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if out != "hello" {
				t.Fatalf("unexpected output: %q", out)
			}
		})
	}
}

// TestStdinBytesCopiesInput ensures later caller mutation cannot alter a configured command.
func TestStdinBytesCopiesInput(t *testing.T) {
	input := []byte("before")
	cmd := helperCommand("cat").StdinBytes(input)
	copy(input, "mutate")

	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "before" {
		t.Fatalf("expected configured bytes to be stable, got %q", out)
	}
}

// TestOutputVariants ensures string, byte, trimmed, and combined capture preserve their distinct contracts.
func TestOutputVariants(t *testing.T) {
	out, err := helperCommand("echo", "  spaced  ").OutputTrimmed()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "spaced" {
		t.Fatalf("unexpected trimmed output: %q", out)
	}

	bytesOut, err := helperCommand("echo", "hi").OutputBytes()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(bytesOut) != "hi" {
		t.Fatalf("unexpected bytes output: %q", string(bytesOut))
	}

	combined, err := helperCommand("mix").CombinedOutput()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if combined != "abc" {
		t.Fatalf("unexpected combined output: %q", combined)
	}
}

// TestExitHelpers ensures ordinary nonzero exits remain results rather than execution failures.
func TestExitHelpers(t *testing.T) {
	res, err := helperCommand("exit", "2").Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.OK() {
		t.Fatalf("expected not OK")
	}
	if !res.IsExitCode(2) {
		t.Fatalf("expected exit code 2, got %d", res.ExitCode)
	}
	if res.IsSignal(syscall.SIGTERM) {
		t.Fatalf("expected no signal for exit")
	}
}

// TestIsSignal ensures signal termination remains distinguishable from numeric exit status.
func TestIsSignal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("signals not supported on windows")
	}
	res, err := helperCommand("signal").Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsSignal(syscall.SIGTERM) {
		t.Fatalf("expected SIGTERM, got %v", res.signal)
	}
}

// TestWithTimeout ensures the shortest configured timeout bounds process execution.
func TestWithTimeout(t *testing.T) {
	_, err := helperCommand("sleep", "200").WithTimeout(50 * time.Millisecond).Run()
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if !errorsIsContext(err) {
		t.Fatalf("expected context error, got %v", err)
	}

	_, err = helperCommand("sleep", "50").WithTimeout(10 * time.Millisecond).WithTimeout(5 * time.Millisecond).Run()
	if err == nil {
		t.Fatalf("expected timeout error on repeated call")
	}
}

// TestWithDeadline ensures expired deadlines stop work while later replacements permit valid commands.
func TestWithDeadline(t *testing.T) {
	_, err := helperCommand("sleep", "100").WithDeadline(time.Now().Add(10 * time.Millisecond)).Run()
	if err == nil {
		t.Fatalf("expected deadline error")
	}

	_, err = helperCommand("echo", "ok").WithDeadline(time.Now().Add(2 * time.Second)).WithDeadline(time.Now().Add(3 * time.Second)).Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestTimeoutPreservesCanceledParent ensures derived timing options cannot hide parent cancellation.
func TestTimeoutPreservesCanceledParent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := helperCommand("sleep", "200").
		WithContext(ctx).
		WithTimeout(time.Second).
		WithDeadline(time.Now().Add(2 * time.Second))
	cancel()

	_, err := cmd.Run()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled parent to remain authoritative, got %v", err)
	}
}

// TestWithContext ensures replacing command context consistently controls subsequent execution.
func TestWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := helperCommand("sleep", "50").WithContext(ctx).Run()
	if err == nil {
		t.Fatalf("expected canceled error")
	}

	_, err = helperCommand("echo", "ok").WithTimeout(500 * time.Millisecond).WithContext(context.Background()).Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestDir ensures the child process observes the configured working directory.
func TestDir(t *testing.T) {
	temp := t.TempDir()
	out, err := helperCommand("pwd").Dir(temp).Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	resolvedTemp, err := filepath.EvalSymlinks(temp)
	if err != nil {
		t.Fatalf("resolve temp: %v", err)
	}
	resolvedOut, err := filepath.EvalSymlinks(out)
	if err != nil {
		t.Fatalf("resolve out: %v", err)
	}
	if resolvedOut != resolvedTemp {
		t.Fatalf("expected dir %q, got %q", resolvedTemp, resolvedOut)
	}
}

// TestPipeModes ensures strict pipelines report any stage failure while best-effort mode follows the last stage.
func TestPipeModes(t *testing.T) {
	strictRes, err := helperPipe(helperCommand("exit", "2"), "echo", "ok").Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if strictRes.ExitCode != 2 {
		t.Fatalf("expected strict pipeline to return first failure, got %d", strictRes.ExitCode)
	}

	bestEffortRes, err := helperPipe(helperCommand("exit", "2").PipeBestEffort(), "echo", "ok").Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if bestEffortRes.ExitCode != 0 {
		t.Fatalf("expected best effort to return last stage, got %d", bestEffortRes.ExitCode)
	}
	if bestEffortRes.Stdout != "ok" {
		t.Fatalf("expected stdout from last stage, got %q", bestEffortRes.Stdout)
	}
}

// TestPipeChain ensures appended stages preserve pipeline order and final output ownership.
func TestPipeChain(t *testing.T) {
	root := helperCommand("echo", "a")
	stage := helperPipe(root, "echo", "b")
	final := helperPipe(stage, "echo", "c")
	res, err := final.Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Stdout != "c" {
		t.Fatalf("expected last stage output, got %q", res.Stdout)
	}
}

// TestPipeBestEffortSetsError ensures context failure remains visible even when the final stage succeeds.
func TestPipeBestEffortSetsError(t *testing.T) {
	res, err := helperPipe(helperCommand("sleep", "50").WithTimeout(10*time.Millisecond).PipeBestEffort(), "echo", "ok").Run()
	if err == nil || !errorsIsContext(err) {
		t.Fatalf("expected context error, got %v", err)
	}
	if res.Stdout != "ok" {
		t.Fatalf("expected stdout from last stage, got %q", res.Stdout)
	}
}

// TestPipeStartError ensures an unstartable stage produces an execution error and sentinel exit code.
func TestPipeStartError(t *testing.T) {
	bad := Command("execx-does-not-exist")
	stage := helperPipe(bad, "echo", "ok")
	res, err := stage.Run()
	if err == nil {
		t.Fatalf("expected start error")
	}
	var errExec ErrExec
	if !errors.As(err, &errExec) {
		t.Fatalf("expected ErrExec, got %T", err)
	}
	if res.ExitCode != -1 {
		t.Fatalf("expected exit code -1, got %d", res.ExitCode)
	}
}

// TestPipelineDownstreamStartErrorAbortsUpstream ensures partial startup cannot leave an earlier process running.
func TestPipelineDownstreamStartErrorAbortsUpstream(t *testing.T) {
	cmd := helperCommand("sleep", "5000").Pipe("execx-does-not-exist")
	type outcome struct {
		result Result
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		result, err := cmd.Run()
		done <- outcome{result: result, err: err}
	}()

	select {
	case got := <-done:
		if got.err == nil {
			t.Fatal("expected downstream start error")
		}
		var execErr ErrExec
		if !errors.As(got.err, &execErr) {
			t.Fatalf("expected ErrExec, got %T", got.err)
		}
		if got.result.ExitCode != -1 {
			t.Fatalf("expected failed stage result, got exit code %d", got.result.ExitCode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("pipeline did not abort after downstream start failure")
	}
}

// TestStringAndShellEscaped ensures display quoting and shell-safe quoting remain intentionally different.
func TestStringAndShellEscaped(t *testing.T) {
	cmd := Command("echo", "hello world", "it's")
	if cmd.String() != "echo \"hello world\" it's" {
		t.Fatalf("unexpected String(): %q", cmd.String())
	}
	if cmd.ShellEscaped() != "echo 'hello world' \"it's\"" {
		t.Fatalf("unexpected ShellEscaped(): %q", cmd.ShellEscaped())
	}

	empty := Command("").ShellEscaped()
	if empty != "''" {
		t.Fatalf("unexpected ShellEscaped empty: %q", empty)
	}
	noQuote := Command("echo", "plain").ShellEscaped()
	if noQuote != "echo plain" {
		t.Fatalf("unexpected ShellEscaped plain: %q", noQuote)
	}
}

// TestLineCallbacks ensures stdout and stderr lines reach only their respective callbacks.
func TestLineCallbacks(t *testing.T) {
	var stdoutLines []string
	var stderrLines []string
	_, err := helperCommand("lines").OnStdout(func(line string) {
		stdoutLines = append(stdoutLines, line)
	}).OnStderr(func(line string) {
		stderrLines = append(stderrLines, line)
	}).Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if strings.Join(stdoutLines, ",") != "a,b" {
		t.Fatalf("unexpected stdout lines: %v", stdoutLines)
	}
	if strings.Join(stderrLines, ",") != "c" {
		t.Fatalf("unexpected stderr lines: %v", stderrLines)
	}
}

// TestLineCallbacksFlushFinalLine ensures unterminated trailing output is not lost at process exit.
func TestLineCallbacksFlushFinalLine(t *testing.T) {
	var stdoutLines []string
	var stderrLines []string
	_, err := helperCommand("echo", "tail").
		OnStdout(func(line string) { stdoutLines = append(stdoutLines, line) }).
		Run()
	if err != nil {
		t.Fatalf("expected no stdout error, got %v", err)
	}
	_, err = helperCommand("stderr", "tail").
		OnStderr(func(line string) { stderrLines = append(stderrLines, line) }).
		Run()
	if err != nil {
		t.Fatalf("expected no stderr error, got %v", err)
	}
	if strings.Join(stdoutLines, ",") != "tail" || strings.Join(stderrLines, ",") != "tail" {
		t.Fatalf("expected final lines, got stdout=%v stderr=%v", stdoutLines, stderrLines)
	}
}

// TestOutputCallbacksAreSerialized ensures shared callback state is safe across concurrent stdout and stderr reads.
func TestOutputCallbacksAreSerialized(t *testing.T) {
	var lines []string
	callback := func(line string) {
		lines = append(lines, line)
	}
	_, err := helperCommand("lines").OnStdout(callback).OnStderr(callback).Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected every output line, got %v", lines)
	}
}

// TestWritersBeforeLineCallbacks ensures byte writers observe output before derived line notifications.
func TestWritersBeforeLineCallbacks(t *testing.T) {
	var order []string
	var stdoutLines []string
	writer := &orderedWriter{order: &order, tag: "writer"}
	_, err := helperCommand("lines").StdoutWriter(writer).OnStdout(func(line string) {
		if len(stdoutLines) == 0 {
			order = append(order, "line")
		}
		stdoutLines = append(stdoutLines, line)
	}).Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(order) == 0 || order[0] != "writer" {
		t.Fatalf("expected writer before line callback, got %v", order)
	}
	if len(writer.buf) == 0 {
		t.Fatalf("expected writer to receive output")
	}
}

// TestStderrWriter ensures raw stderr forwarding and line callbacks can coexist without data loss.
func TestStderrWriter(t *testing.T) {
	var stderrLines []string
	writer := &orderedWriter{tag: "stderr"}
	_, err := helperCommand("lines").StderrWriter(writer).OnStderr(func(line string) {
		stderrLines = append(stderrLines, line)
	}).Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(writer.buf) == 0 {
		t.Fatalf("expected stderr writer to receive output")
	}
	if strings.Join(stderrLines, ",") != "c" {
		t.Fatalf("unexpected stderr lines: %v", stderrLines)
	}
}

// TestStartAndWait ensures asynchronous execution reports the same clean result as Run.
func TestStartAndWait(t *testing.T) {
	proc := helperCommand("sleep", "50").Start()
	res, err := proc.Wait()
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("expected clean exit, got code=%d err=%v", res.ExitCode, err)
	}
}

// TestStartError ensures startup failures retain ErrExec identity and an unavailable exit code.
func TestStartError(t *testing.T) {
	res, err := Command("execx-does-not-exist").Run()
	if err == nil {
		t.Fatalf("expected start error")
	}
	var errExec ErrExec
	if !errors.As(err, &errExec) {
		t.Fatalf("expected ErrExec, got %T", err)
	}
	if res.ExitCode != -1 {
		t.Fatalf("expected exit code -1 for start error, got %d", res.ExitCode)
	}
}

// TestLineWriterNil ensures a line splitter without a callback still satisfies io.Writer safely.
func TestLineWriterNil(t *testing.T) {
	writer := &lineWriter{}
	n, err := writer.Write([]byte("data"))
	if err != nil || n != 4 {
		t.Fatalf("unexpected write result n=%d err=%v", n, err)
	}
}

// TestWithPTYPipelineUnsupported ensures PTY mode rejects pipelines whose stream topology cannot be represented safely.
func TestWithPTYPipelineUnsupported(t *testing.T) {
	prevCheck := ptyCheckFunc
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		ptyCheckFunc = prevCheck
	})
	_, err := Command("printf", "hi").
		WithPTY().
		Pipe("tr", "a-z", "A-Z").
		Run()
	if err == nil || !strings.Contains(err.Error(), "WithPTY is not supported") {
		t.Fatalf("expected WithPTY pipeline error, got %v", err)
	}
}

// TestWithPTYOpenError ensures terminal allocation failures abort execution with their cause intact.
func TestWithPTYOpenError(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		return nil, nil, errors.New("pty open failed")
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})
	_, err := Command("printf", "hi").WithPTY().Run()
	if err == nil || !strings.Contains(err.Error(), "pty open failed") {
		t.Fatalf("expected openpty error, got %v", err)
	}
}

// TestWithPTYCombinedStream ensures PTY output is captured once while notifying both logical stream callbacks.
func TestWithPTYCombinedStream(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		r, w, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		return r, w, nil
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})
	stdoutLines := []string{}
	stderrLines := []string{}
	res, err := Command("printf", "a\nb\n").
		WithPTY().
		OnStdout(func(line string) { stdoutLines = append(stdoutLines, line) }).
		OnStderr(func(line string) { stderrLines = append(stderrLines, line) }).
		Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Stdout != "a\nb\n" {
		t.Fatalf("expected stdout to contain output, got %q", res.Stdout)
	}
	if res.Stderr != "" {
		t.Fatalf("expected stderr to be empty, got %q", res.Stderr)
	}
	if strings.Join(stdoutLines, ",") != "a,b" {
		t.Fatalf("unexpected stdout lines: %v", stdoutLines)
	}
	if strings.Join(stderrLines, ",") != "a,b" {
		t.Fatalf("unexpected stderr lines: %v", stderrLines)
	}
}

// TestWithPTYCombinedOutput ensures the combined-output convenience API works with a single PTY stream.
func TestWithPTYCombinedOutput(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		r, w, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		return r, w, nil
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})
	out, err := Command("printf", "hi").WithPTY().CombinedOutput()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "hi" {
		t.Fatalf("expected combined output, got %q", out)
	}
}

// TestWithPTYPipelineResults ensures a one-stage PTY command retains the pipeline-results API contract.
func TestWithPTYPipelineResults(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		r, w, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		return r, w, nil
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})
	results, err := Command("printf", "ok").WithPTY().PipelineResults()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 1 || results[0].Stdout != "ok" {
		t.Fatalf("unexpected pipeline results: %+v", results)
	}
}

// TestWithPTYStart ensures asynchronous PTY capture is complete before Wait returns.
func TestWithPTYStart(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		r, w, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		return r, w, nil
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})
	proc := Command("printf", "hi").WithPTY().Start()
	res, err := proc.Wait()
	if err != nil || res.Stdout != "hi" {
		t.Fatalf("expected stdout from Start, got %q err=%v", res.Stdout, err)
	}
}

// TestWithPTYCheckError ensures unsupported terminals fail before process startup.
func TestWithPTYCheckError(t *testing.T) {
	prevCheck := ptyCheckFunc
	ptyCheckFunc = func() error { return errors.New("pty unsupported") }
	t.Cleanup(func() {
		ptyCheckFunc = prevCheck
	})
	_, err := Command("printf", "hi").WithPTY().CombinedOutput()
	if err == nil || !strings.Contains(err.Error(), "pty unsupported") {
		t.Fatalf("expected pty check error, got %v", err)
	}
}

// TestWithPTYStartCheckError ensures asynchronous PTY validation failures surface through Wait.
func TestWithPTYStartCheckError(t *testing.T) {
	prevCheck := ptyCheckFunc
	ptyCheckFunc = func() error { return errors.New("pty unsupported") }
	t.Cleanup(func() {
		ptyCheckFunc = prevCheck
	})
	proc := Command("printf", "hi").WithPTY().Start()
	res, err := proc.Wait()
	if err == nil || !strings.Contains(err.Error(), "pty unsupported") {
		t.Fatalf("expected pty check error, got %v", err)
	}
	if res.ExitCode != -1 {
		t.Fatalf("expected exit code -1, got %d", res.ExitCode)
	}
}

// TestWithPTYPipelineResultsCheckError ensures PTY validation also protects the pipeline-results entry point.
func TestWithPTYPipelineResultsCheckError(t *testing.T) {
	prevCheck := ptyCheckFunc
	ptyCheckFunc = func() error { return errors.New("pty unsupported") }
	t.Cleanup(func() {
		ptyCheckFunc = prevCheck
	})
	_, err := Command("printf", "hi").WithPTY().PipelineResults()
	if err == nil || !strings.Contains(err.Error(), "pty unsupported") {
		t.Fatalf("expected pty check error, got %v", err)
	}
}

type errWriter struct {
	called bool
}

// Write records the attempt before returning a deterministic output failure.
func (w *errWriter) Write(p []byte) (int, error) {
	w.called = true
	return 0, errors.New("write failed")
}

type errReader struct{}

// Read returns a deterministic stdin-copy failure.
func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

type sliceWriter []byte

// Write accepts bytes without mutating its non-comparable value receiver.
func (sliceWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// TestWithPTYWriterError ensures forwarding failures become ErrExec without discarding captured output.
func TestWithPTYWriterError(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		r, w, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		return r, w, nil
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})
	writer := &errWriter{}
	res, err := Command("printf", "hi").WithPTY().StdoutWriter(writer).Run()
	if err == nil {
		t.Fatal("expected writer error")
	}
	var execErr ErrExec
	if !errors.As(err, &execErr) {
		t.Fatalf("expected ErrExec, got %T", err)
	}
	if !writer.called {
		t.Fatalf("expected writer to be called")
	}
	if res.Stdout != "hi" {
		t.Fatalf("expected captured stdout before writer failure, got %q", res.Stdout)
	}
}

// TestWriterErrorIsExecutionError ensures ordinary stream-writer failures use the same ErrExec boundary as process failures.
func TestWriterErrorIsExecutionError(t *testing.T) {
	writer := &errWriter{}
	res, err := helperCommand("echo", "hi").StdoutWriter(writer).Run()
	if err == nil {
		t.Fatal("expected writer error")
	}
	var execErr ErrExec
	if !errors.As(err, &execErr) {
		t.Fatalf("expected ErrExec, got %T", err)
	}
	if res.Stdout != "hi" {
		t.Fatalf("expected result capture before writer failure, got %q", res.Stdout)
	}
}

// TestReaderErrorIsExecutionError ensures stdin-copy failures cannot be mistaken for successful execution.
func TestReaderErrorIsExecutionError(t *testing.T) {
	_, err := helperCommand("cat").StdinReader(errReader{}).Run()
	if err == nil {
		t.Fatal("expected reader error")
	}
	var execErr ErrExec
	if !errors.As(err, &execErr) {
		t.Fatalf("expected ErrExec, got %T", err)
	}
}

// TestWithPTYStartError ensures process startup failures remain visible after successful terminal allocation.
func TestWithPTYStartError(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		r, w, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		return r, w, nil
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})
	_, err := Command("execx-does-not-exist").WithPTY().Run()
	if err == nil {
		t.Fatalf("expected start error")
	}
}

// TestWithPTYWritersNoCallbacks ensures both configured writers receive the merged PTY stream without callbacks.
func TestWithPTYWritersNoCallbacks(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		r, w, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		return r, w, nil
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	res, err := Command("printf", "hi").
		WithPTY().
		StdoutWriter(&stdoutBuf).
		StderrWriter(&stderrBuf).
		Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Stdout != "hi" {
		t.Fatalf("expected stdout buffer to capture output, got %q", res.Stdout)
	}
	if stdoutBuf.String() != "hi" || stderrBuf.String() != "hi" {
		t.Fatalf("unexpected writers: stdout=%q stderr=%q", stdoutBuf.String(), stderrBuf.String())
	}
}

// TestWithPTYSharedWriterAndCallback ensures a shared writer receives merged bytes once while callbacks remain serialized.
func TestWithPTYSharedWriterAndCallback(t *testing.T) {
	prevOpen := openPTYFunc
	prevCheck := ptyCheckFunc
	openPTYFunc = func() (*os.File, *os.File, error) {
		r, w, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		return r, w, nil
	}
	ptyCheckFunc = func() error { return nil }
	t.Cleanup(func() {
		openPTYFunc = prevOpen
		ptyCheckFunc = prevCheck
	})

	var output bytes.Buffer
	var lines []string
	callback := func(line string) { lines = append(lines, line) }
	_, err := Command("printf", "a\nb").
		WithPTY().
		StdoutWriter(&output).
		StderrWriter(&output).
		OnStdout(callback).
		OnStderr(callback).
		Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.String() != "a\nb" {
		t.Fatalf("expected shared writer to receive one merged stream, got %q", output.String())
	}
	if strings.Join(lines, ",") != "a,a,b,b" {
		t.Fatalf("expected serialized PTY callbacks, got %v", lines)
	}
}

// TestPTYLineWriterNil ensures an unconfigured PTY splitter remains a safe io.Writer.
func TestPTYLineWriterNil(t *testing.T) {
	writer := &ptyLineWriter{}
	n, err := writer.Write([]byte("data"))
	if err != nil || n != 4 {
		t.Fatalf("unexpected write result n=%d err=%v", n, err)
	}
}

// TestPTYLineWriterFlush ensures a trailing carriage-return line is emitted exactly once.
func TestPTYLineWriterFlush(t *testing.T) {
	var stdoutLines []string
	var stderrLines []string
	writer := &ptyLineWriter{
		onStdout: func(line string) { stdoutLines = append(stdoutLines, line) },
		onStderr: func(line string) { stderrLines = append(stderrLines, line) },
	}
	_, _ = writer.Write([]byte("tail\r"))
	writer.Flush()
	writer.Flush()
	if strings.Join(stdoutLines, ",") != "tail" || strings.Join(stderrLines, ",") != "tail" {
		t.Fatalf("expected flushed PTY lines, got stdout=%v stderr=%v", stdoutLines, stderrLines)
	}
}

// TestSameWriter ensures writer deduplication avoids panics for non-comparable implementations.
func TestSameWriter(t *testing.T) {
	var first bytes.Buffer
	var second bytes.Buffer
	if !sameWriter(nil, nil) {
		t.Fatal("expected two nil writers to match")
	}
	if sameWriter(nil, &first) || sameWriter(&first, &second) {
		t.Fatal("expected distinct writers not to match")
	}
	if !sameWriter(&first, &first) {
		t.Fatal("expected identical writer pointers to match")
	}
	if sameWriter(sliceWriter{}, sliceWriter{}) {
		t.Fatal("expected non-comparable writer values not to be compared")
	}
}

// TestOnExecCmdApplied ensures the final os/exec command remains customizable before startup.
func TestOnExecCmdApplied(t *testing.T) {
	called := false
	cmd := Command("printf", "hi").OnExecCmd(func(ec *exec.Cmd) {
		called = true
		ec.Env = append(ec.Env, "EXECX_TEST=1")
	})
	execCmd := cmd.execCmd()
	if !called {
		t.Fatalf("expected OnExecCmd callback to run")
	}
	found := false
	for _, entry := range execCmd.Env {
		if entry == "EXECX_TEST=1" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected OnExecCmd to mutate env")
	}
}

// TestIsTerminalWriterNonFile ensures terminal probing never assumes arbitrary writers are files.
func TestIsTerminalWriterNonFile(t *testing.T) {
	var buf bytes.Buffer
	if isTerminalWriter(&buf) {
		t.Fatalf("expected non-file writer to be non-terminal")
	}
}

// TestStdoutWriterTTYPassthrough ensures interactive stdout is not hidden behind capture buffering.
func TestStdoutWriterTTYPassthrough(t *testing.T) {
	prev := isTerminalFunc
	isTerminalFunc = func(int) bool { return true }
	t.Cleanup(func() {
		isTerminalFunc = prev
	})
	cmd := Command("printf", "hi").StdoutWriter(os.Stdout)
	out := cmd.stdoutWriter(&bytes.Buffer{}, false, &bytes.Buffer{}, nil)
	if out != os.Stdout {
		t.Fatalf("expected stdout writer to passthrough tty")
	}
}

// TestStderrWriterTTYPassthrough ensures interactive stderr is not hidden behind capture buffering.
func TestStderrWriterTTYPassthrough(t *testing.T) {
	prev := isTerminalFunc
	isTerminalFunc = func(int) bool { return true }
	t.Cleanup(func() {
		isTerminalFunc = prev
	})
	cmd := Command("printf", "hi").StderrWriter(os.Stderr)
	out := cmd.stderrWriter(&bytes.Buffer{}, false, &bytes.Buffer{}, nil)
	if out != os.Stderr {
		t.Fatalf("expected stderr writer to passthrough tty")
	}
}

// TestSignalFromStateNil ensures absent process state cannot fabricate a termination signal.
func TestSignalFromStateNil(t *testing.T) {
	if signalFromState(nil) != nil {
		t.Fatalf("expected nil signal")
	}
}

// TestRootCmd ensures an unpiped command is its own pipeline root.
func TestRootCmd(t *testing.T) {
	cmd := &Cmd{}
	if cmd.rootCmd() != cmd {
		t.Fatalf("expected rootCmd to return self")
	}
}

// TestStageResultContextError ensures cancellation takes precedence when no process state exists.
func TestStageResultContextError(t *testing.T) {
	st := &stage{
		waitErr: context.Canceled,
		def:     &Cmd{},
		cmd:     &exec.Cmd{},
	}
	res := st.result()
	if !errors.Is(res.Err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", res.Err)
	}
	if res.ExitCode != -1 {
		t.Fatalf("expected exit code -1, got %d", res.ExitCode)
	}
}

// TestPipelineResults ensures callers receive one ordered result per pipeline stage.
func TestPipelineResults(t *testing.T) {
	root := helperCommand("echo", "a")
	stage := helperPipe(root, "echo", "b")
	final := helperPipe(stage, "echo", "c")
	results, err := final.PipelineResults()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[2].Stdout != "c" {
		t.Fatalf("expected last stage stdout, got %q", results[2].Stdout)
	}
}

// TestPipelineResultsError ensures a failed single stage still has an inspectable result.
func TestPipelineResultsError(t *testing.T) {
	results, err := Command("execx-does-not-exist").PipelineResults()
	if err == nil {
		t.Fatalf("expected error")
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Fatalf("expected result error")
	}
	var errExec ErrExec
	if !errors.As(err, &errExec) {
		t.Fatalf("expected ErrExec, got %T", err)
	}
}

// TestPipelineStartErrorPropagation ensures startup failure is recorded on the stage that could not run.
func TestPipelineStartErrorPropagation(t *testing.T) {
	results, err := Command("execx-does-not-exist").
		Pipe("printf", "ok").
		PipelineResults()
	if err == nil {
		t.Fatalf("expected error")
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[1].Err == nil {
		t.Fatalf("expected downstream start error")
	}
}

// TestProcessSignals ensures callers can deliver an explicit operating-system signal to a running process.
func TestProcessSignals(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("signals not supported on windows")
	}
	proc := helperCommand("sleep", "200").Start()
	if err := proc.Send(syscall.SIGTERM); err != nil {
		t.Fatalf("send signal: %v", err)
	}
	res, err := proc.Wait()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsSignal(syscall.SIGTERM) {
		t.Fatalf("expected SIGTERM, got %v", res.signal)
	}
}

// TestProcessInterrupt ensures the convenience API records interrupt termination accurately.
func TestProcessInterrupt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("signals not supported on windows")
	}
	proc := helperCommand("sleep", "200").Start()
	if err := proc.Interrupt(); err != nil {
		t.Fatalf("interrupt: %v", err)
	}
	res, err := proc.Wait()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsSignal(os.Interrupt) {
		t.Fatalf("expected interrupt, got %v", res.signal)
	}
}

// TestProcessTerminate ensures forced termination cannot be reported as a clean exit.
func TestProcessTerminate(t *testing.T) {
	proc := helperCommand("sleep", "200").Start()
	if err := proc.Terminate(); err != nil {
		t.Fatalf("terminate: %v", err)
	}
	res, err := proc.Wait()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ExitCode == 0 {
		t.Fatalf("expected non-zero exit")
	}
}

// TestGracefulShutdownKills ensures an uncooperative process is killed after its grace period.
func TestGracefulShutdownKills(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("signals not supported on windows")
	}
	proc := helperCommand("ignore-term").Start()
	time.Sleep(50 * time.Millisecond)
	if err := proc.GracefulShutdown(syscall.SIGTERM, 20*time.Millisecond); err != nil {
		t.Fatalf("graceful shutdown: %v", err)
	}
	res, err := proc.Wait()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsSignal(syscall.SIGKILL) {
		t.Fatalf("expected SIGKILL, got %v", res.signal)
	}
}

// TestKillAfter ensures repeated scheduling still terminates a long-running process safely.
func TestKillAfter(t *testing.T) {
	proc := helperCommand("sleep", "200").Start()
	proc.KillAfter(10 * time.Millisecond)
	proc.KillAfter(20 * time.Millisecond)
	res, err := proc.Wait()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ExitCode == 0 {
		t.Fatalf("expected killed process")
	}
}

// TestGracefulShutdownCompletes ensures cooperative signal exit avoids escalation to a kill.
func TestGracefulShutdownCompletes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("signals not supported on windows")
	}
	proc := helperCommand("sleep", "200").Start()
	if err := proc.GracefulShutdown(syscall.SIGTERM, 200*time.Millisecond); err != nil {
		t.Fatalf("graceful shutdown: %v", err)
	}
	res, err := proc.Wait()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.IsSignal(syscall.SIGTERM) {
		t.Fatalf("expected SIGTERM, got %v", res.signal)
	}
}

// TestGracefulShutdownImmediate ensures a zero grace period escalates without waiting.
func TestGracefulShutdownImmediate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("signals not supported on windows")
	}
	proc := helperCommand("sleep", "200").Start()
	if err := proc.GracefulShutdown(syscall.SIGTERM, 0); err != nil {
		t.Fatalf("graceful shutdown immediate: %v", err)
	}
	res, err := proc.Wait()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ExitCode == 0 {
		t.Fatalf("expected killed process")
	}
}

// TestProcessSendErrors ensures missing process and pipeline state return errors instead of panicking.
func TestProcessSendErrors(t *testing.T) {
	var proc *Process
	if err := proc.Send(os.Interrupt); err == nil {
		t.Fatalf("expected send error for nil process")
	}
	proc = &Process{pipeline: &pipeline{}}
	if err := proc.Send(os.Interrupt); err == nil {
		t.Fatalf("expected send error for empty pipeline")
	}
	if err := proc.GracefulShutdown(os.Interrupt, 10*time.Millisecond); err == nil {
		t.Fatalf("expected graceful shutdown error for empty pipeline")
	}
}

// TestProcessSendSkipsStages ensures incomplete pipeline stages are ignored while reporting that nothing was signaled.
func TestProcessSendSkipsStages(t *testing.T) {
	proc := &Process{
		pipeline: &pipeline{
			stages: []*stage{
				nil,
				{},
				{cmd: &exec.Cmd{}},
			},
		},
		done: make(chan struct{}),
	}
	if err := proc.Send(os.Interrupt); err == nil {
		t.Fatalf("expected send error for empty stages")
	}
}

// TestProcessSendAfterExit ensures signaling a reaped process returns an actionable error.
func TestProcessSendAfterExit(t *testing.T) {
	proc := helperCommand("echo", "ok").Start()
	_, _ = proc.Wait()
	if err := proc.Send(os.Interrupt); err == nil {
		t.Fatalf("expected send error after exit")
	}
}

// TestErrExecMethods ensures execution errors preserve wrapping semantics and a useful empty fallback.
func TestErrExecMethods(t *testing.T) {
	baseErr := errors.New("boom")
	execErr := ErrExec{Err: baseErr}
	if execErr.Error() != "boom" {
		t.Fatalf("unexpected error string: %q", execErr.Error())
	}
	if !errors.Is(execErr, baseErr) {
		t.Fatalf("expected unwrap to match")
	}
	empty := ErrExec{}
	if empty.Error() == "" {
		t.Fatalf("expected default error string")
	}
	if empty.Unwrap() != nil {
		t.Fatalf("expected nil unwrap")
	}
}

// TestSysProcAttrNoops ensures platform-specific process options remain harmless where unsupported.
func TestSysProcAttrNoops(t *testing.T) {
	cmd := Command("echo")
	cmd.CreationFlags(123).HideWindow(true).Pdeathsig(syscall.SIGTERM)
	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		if cmd.sysProcAttr == nil {
			t.Fatalf("expected sys proc attr on supported platform")
		}
		return
	}
	if cmd.sysProcAttr != nil {
		t.Fatalf("expected no sys proc attr on unsupported platform")
	}
}

type orderedWriter struct {
	order *[]string
	tag   string
	buf   []byte
}

// Write records first-observer ordering while retaining all received bytes.
func (w *orderedWriter) Write(p []byte) (int, error) {
	if w.order != nil && len(*w.order) == 0 {
		*w.order = append(*w.order, w.tag)
	}
	w.buf = append(w.buf, p...)
	return len(p), nil
}

// TestPipeStrictExplicit ensures explicitly selected strict mode reports an earlier stage failure.
func TestPipeStrictExplicit(t *testing.T) {
	res, err := helperPipe(helperCommand("exit", "2").PipeStrict(), "echo", "ok").Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ExitCode != 2 {
		t.Fatalf("expected strict pipeline to return first failure, got %d", res.ExitCode)
	}
}

// errorsIsContext accepts either cancellation outcome used by timing-sensitive process tests.
func errorsIsContext(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}
