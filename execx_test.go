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
		if runtime.GOOS == "windows" {
			os.Exit(3)
		}
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
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

func helperCommand(args ...string) *Cmd {
	full := append([]string{"-test.run=TestHelperProcess", "--"}, args...)
	cmd := Command(os.Args[0], full...)
	cmd.Env("EXECX_TEST_HELPER=1")
	return cmd
}

func helperPipe(cmd *Cmd, args ...string) *Cmd {
	full := append([]string{"-test.run=TestHelperProcess", "--"}, args...)
	stage := cmd.Pipe(os.Args[0], full...)
	stage.Env("EXECX_TEST_HELPER=1")
	return stage
}

type envStringer struct{}

func (envStringer) String() string {
	return "EXECX_ENV_VALUE=stringer"
}

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

func TestArgVariants(t *testing.T) {
	out, err := helperCommand("echo").Arg([]string{"a", "b"}, 123).Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out != "a b 123" {
		t.Fatalf("unexpected output: %q", out)
	}
}

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

func TestEnvList(t *testing.T) {
	cmd := helperCommand("env", "NONE").EnvOnly(map[string]string{"B": "2", "A": "1", "EXECX_TEST_HELPER": "1"})
	list := cmd.EnvList()
	if strings.Join(list, ",") != "A=1,B=2,EXECX_TEST_HELPER=1" {
		t.Fatalf("unexpected env list: %v", list)
	}
}

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

func TestWithDeadline(t *testing.T) {
	_, err := helperCommand("sleep", "100").WithDeadline(time.Now().Add(10 * time.Millisecond)).Run()
	if err == nil {
		t.Fatalf("expected deadline error")
	}

	_, err = helperCommand("echo", "ok").WithDeadline(time.Now().Add(200 * time.Millisecond)).WithDeadline(time.Now().Add(300 * time.Millisecond)).Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

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

func TestPipeBestEffortSetsError(t *testing.T) {
	res, err := helperPipe(helperCommand("sleep", "50").WithTimeout(10*time.Millisecond).PipeBestEffort(), "echo", "ok").Run()
	if err == nil || !errorsIsContext(err) {
		t.Fatalf("expected context error, got %v", err)
	}
	if res.Stdout != "ok" {
		t.Fatalf("expected stdout from last stage, got %q", res.Stdout)
	}
}

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

func TestStartAndWait(t *testing.T) {
	proc := helperCommand("sleep", "50").Start()
	res, err := proc.Wait()
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("expected clean exit, got code=%d err=%v", res.ExitCode, err)
	}
}

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

func TestLineWriterNil(t *testing.T) {
	writer := &lineWriter{}
	n, err := writer.Write([]byte("data"))
	if err != nil || n != 4 {
		t.Fatalf("unexpected write result n=%d err=%v", n, err)
	}
}

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

func TestIsTerminalWriterNonFile(t *testing.T) {
	var buf bytes.Buffer
	if isTerminalWriter(&buf) {
		t.Fatalf("expected non-file writer to be non-terminal")
	}
}

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

func TestSignalFromStateNil(t *testing.T) {
	if signalFromState(nil) != nil {
		t.Fatalf("expected nil signal")
	}
}

func TestRootCmd(t *testing.T) {
	cmd := &Cmd{}
	if cmd.rootCmd() != cmd {
		t.Fatalf("expected rootCmd to return self")
	}
}

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

func TestProcessSendAfterExit(t *testing.T) {
	proc := helperCommand("echo", "ok").Start()
	_, _ = proc.Wait()
	if err := proc.Send(os.Interrupt); err == nil {
		t.Fatalf("expected send error after exit")
	}
}

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

func (w *orderedWriter) Write(p []byte) (int, error) {
	if w.order != nil && len(*w.order) == 0 {
		*w.order = append(*w.order, w.tag)
	}
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func TestPipeStrictExplicit(t *testing.T) {
	res, err := helperPipe(helperCommand("exit", "2").PipeStrict(), "echo", "ok").Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ExitCode != 2 {
		t.Fatalf("expected strict pipeline to return first failure, got %d", res.ExitCode)
	}
}

func errorsIsContext(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}
