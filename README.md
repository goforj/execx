<p align="center">
  <img src="./docs/images/logo.png?v=2" width="400" alt="str logo">
</p>

<p align="center">
    execx is an ergonomic, fluent wrapper around Go’s `os/exec` package.
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/goforj/execx"><img src="https://pkg.go.dev/badge/github.com/goforj/execx.svg" alt="Go Reference"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License: MIT"></a>
    <a href="https://github.com/goforj/execx/actions"><img src="https://github.com/goforj/execx/actions/workflows/test.yml/badge.svg" alt="Go Test"></a>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/go-1.23+-blue?logo=go" alt="Go version"></a>
    <img src="https://img.shields.io/github/v/tag/goforj/execx?label=version&sort=semver" alt="Latest tag"> 
    <a href="https://codecov.io/gh/goforj/execx" ><img src="https://codecov.io/github/goforj/execx/graph/badge.svg?token=RBB8T6WQ0U"/></a>
<!-- test-count:embed:start -->
    <img src="https://img.shields.io/badge/tests-99-brightgreen" alt="Tests">
<!-- test-count:embed:end -->
    <a href="https://goreportcard.com/report/github.com/goforj/execx"><img src="https://goreportcard.com/badge/github.com/goforj/execx" alt="Go Report Card"></a>
</p>

## What execx is

execx is a small, explicit wrapper around `os/exec`. It keeps the `exec.Cmd` model but adds fluent construction and consistent result handling.

There is no shell interpolation. Arguments, environment, and I/O are set directly, and nothing runs until you call `Run`, `Output`, or `Start`.

## Installation

```bash
go get github.com/goforj/execx
```

## Quick Start

```go
out, _ := execx.Command("echo", "hello").OutputTrimmed()
fmt.Println(out)
// #string hello
```

On Windows, use `cmd /c echo hello` or `powershell -Command "echo hello"` for shell built-ins.

## Basic usage

Build a command and run it:

```go
cmd := execx.Command("echo").Arg("hello")
res, _ := cmd.Run()
fmt.Print(res.Stdout)
// hello
```

Arguments are appended deterministically and never shell-expanded.

## Output handling

Use `Output` variants when you only need stdout:

```go
out, _ := execx.Command("echo", "hello").OutputTrimmed()
fmt.Println(out)
// #string hello
```

`Output`, `OutputBytes`, `OutputTrimmed`, and `CombinedOutput` differ only in how they return data.

## Pipelining

Pipelines run on all platforms; command availability is OS-specific.

```go
out, _ := execx.Command("printf", "go").
	Pipe("tr", "a-z", "A-Z").
	OutputTrimmed()
fmt.Println(out)
// #string GO
```

On Windows, use `cmd /c` or `powershell -Command` for shell built-ins.

`PipeStrict` (default) stops at the first failing stage and returns that error.  
`PipeBestEffort` runs all stages, returns the last stage output, and surfaces the first error if any stage failed.

## Context & cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
res, _ := execx.Command("go", "env", "GOOS").WithContext(ctx).Run()
fmt.Println(res.ExitCode == 0)
// #bool true
```

## Environment & I/O control

Environment is explicit and deterministic:

```go
cmd := execx.Command("echo", "hello").Env("MODE=prod")
fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "MODE=prod"))
// #bool true
```

Standard input is opt-in:

```go
out, _ := execx.Command("cat").
	StdinString("hi").
	OutputTrimmed()
fmt.Println(out)
// #string hi
```

## Advanced features

For process control, use `Start` with the `Process` helpers:

```go
proc := execx.Command("go", "env", "GOOS").Start()
res, _ := proc.Wait()
fmt.Println(res.ExitCode == 0)
// #bool true
```

Signals, timeouts, and OS controls are documented in the API section below.

## Kitchen Sink Chaining Example

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

res := execx.
    Command("printf", "hello\nworld\n").
    Pipe("tr", "a-z", "A-Z").
    Env("MODE=demo").
    WithContext(ctx).
    OnStdout(func(line string) {
        fmt.Println("OUT:", line)
    }).
    OnStderr(func(line string) {
        fmt.Println("ERR:", line)
    }).
    Run()

if !res.OK() {
    log.Fatalf("command failed: %v", res.Err)
}
```

## Non-goals and design principles

Design principles:

* Explicit over implicit
* No shell interpolation
* Composable, deterministic behavior

Non-goals:

* Shell scripting replacement
* Command parsing or glob expansion
* Task runners or build systems
* Automatic retries or heuristics

All public APIs are covered by runnable examples under `./examples`, and the test suite executes them to keep docs and behavior in sync.

<!-- api:embed:start -->

## API Index

| Group | Functions |
|------:|:-----------|
| **Arguments** | [Arg](#arg) |
| **Construction** | [Command](#command) |
| **Context** | [WithContext](#withcontext) [WithDeadline](#withdeadline) [WithTimeout](#withtimeout) |
| **Debugging** | [Args](#args) [ShellEscaped](#shellescaped) [String](#string) |
| **Environment** | [Env](#env) [EnvAppend](#envappend) [EnvInherit](#envinherit) [EnvList](#envlist) [EnvOnly](#envonly) |
| **Errors** | [Error](#error) [Unwrap](#unwrap) |
| **Execution** | [CombinedOutput](#combinedoutput) [Output](#output) [OutputBytes](#outputbytes) [OutputTrimmed](#outputtrimmed) [Run](#run) [Start](#start) |
| **Input** | [StdinBytes](#stdinbytes) [StdinFile](#stdinfile) [StdinReader](#stdinreader) [StdinString](#stdinstring) |
| **OS Controls** | [CreationFlags](#creationflags) [HideWindow](#hidewindow) [Pdeathsig](#pdeathsig) [Setpgid](#setpgid) [Setsid](#setsid) |
| **Pipelining** | [Pipe](#pipe) [PipeBestEffort](#pipebesteffort) [PipeStrict](#pipestrict) [PipelineResults](#pipelineresults) |
| **Process** | [GracefulShutdown](#gracefulshutdown) [Interrupt](#interrupt) [KillAfter](#killafter) [Send](#send) [Terminate](#terminate) [Wait](#wait) |
| **Results** | [IsExitCode](#isexitcode) [IsSignal](#issignal) [OK](#ok) |
| **Streaming** | [OnStderr](#onstderr) [OnStdout](#onstdout) [StderrWriter](#stderrwriter) [StdoutWriter](#stdoutwriter) |
| **WorkingDir** | [Dir](#dir) |


## Arguments

### <a id="arg"></a>Arg

Arg appends arguments to the command.

```go
cmd := execx.Command("go", "env").Arg("GOOS")
out, _ := cmd.Output()
fmt.Println(out != "")
// #bool true
```

## Construction

### <a id="command"></a>Command

Command constructs a new command without executing it.

```go
cmd := execx.Command("go", "env", "GOOS")
out, _ := cmd.Output()
fmt.Println(out != "")
// #bool true
```

## Context

### <a id="withcontext"></a>WithContext

WithContext binds the command to a context.

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
res, _ := execx.Command("go", "env", "GOOS").WithContext(ctx).Run()
fmt.Println(res.ExitCode == 0)
// #bool true
```

### <a id="withdeadline"></a>WithDeadline

WithDeadline binds the command to a deadline.

```go
res, _ := execx.Command("go", "env", "GOOS").WithDeadline(time.Now().Add(2 * time.Second)).Run()
fmt.Println(res.ExitCode == 0)
// #bool true
```

### <a id="withtimeout"></a>WithTimeout

WithTimeout binds the command to a timeout.

```go
res, _ := execx.Command("go", "env", "GOOS").WithTimeout(2 * time.Second).Run()
fmt.Println(res.ExitCode == 0)
// #bool true
```

## Debugging

### <a id="args"></a>Args

Args returns the argv slice used for execution.

```go
cmd := execx.Command("go", "env", "GOOS")
fmt.Println(strings.Join(cmd.Args(), " "))
// #string go env GOOS
```

### <a id="shellescaped"></a>ShellEscaped

ShellEscaped returns a shell-escaped string for logging only.

```go
cmd := execx.Command("echo", "hello world", "it's")
fmt.Println(cmd.ShellEscaped())
// #string echo 'hello world' 'it'\\''s'
```

### <a id="string"></a>String

String returns a human-readable representation of the command.

```go
cmd := execx.Command("echo", "hello world", "it's")
fmt.Println(cmd.String())
// #string echo "hello world" it's
```

## Environment

### <a id="env"></a>Env

Env adds environment variables to the command.

```go
cmd := execx.Command("go", "env", "GOOS").Env("MODE=prod")
fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "MODE=prod"))
// #bool true
```

### <a id="envappend"></a>EnvAppend

EnvAppend merges variables into the inherited environment.

```go
cmd := execx.Command("go", "env", "GOOS").EnvAppend(map[string]string{"A": "1"})
fmt.Println(strings.Contains(strings.Join(cmd.EnvList(), ","), "A=1"))
// #bool true
```

### <a id="envinherit"></a>EnvInherit

EnvInherit restores default environment inheritance.

```go
cmd := execx.Command("go", "env", "GOOS").EnvInherit()
fmt.Println(len(cmd.EnvList()) > 0)
// #bool true
```

### <a id="envlist"></a>EnvList

EnvList returns the environment list for execution.

```go
cmd := execx.Command("go", "env", "GOOS").EnvOnly(map[string]string{"A": "1"})
fmt.Println(strings.Join(cmd.EnvList(), ","))
// #string A=1
```

### <a id="envonly"></a>EnvOnly

EnvOnly ignores the parent environment.

```go
cmd := execx.Command("go", "env", "GOOS").EnvOnly(map[string]string{"A": "1"})
fmt.Println(strings.Join(cmd.EnvList(), ","))
// #string A=1
```

## Errors

### <a id="error"></a>Error

Error returns the wrapped error message when available.

```go
err := execx.ErrExec{Err: fmt.Errorf("boom")}
fmt.Println(err.Error())
// #string boom
```

### <a id="unwrap"></a>Unwrap

Unwrap exposes the underlying error.

```go
err := execx.ErrExec{Err: fmt.Errorf("boom")}
fmt.Println(err.Unwrap() != nil)
// #bool true
```

## Execution

### <a id="combinedoutput"></a>CombinedOutput

CombinedOutput executes the command and returns stdout+stderr and any error.

```go
out, _ := execx.Command("go", "env", "GOOS").CombinedOutput()
fmt.Println(out != "")
// #bool true
```

### <a id="output"></a>Output

Output executes the command and returns stdout and any error.

```go
out, _ := execx.Command("go", "env", "GOOS").Output()
fmt.Println(out != "")
// #bool true
```

### <a id="outputbytes"></a>OutputBytes

OutputBytes executes the command and returns stdout bytes and any error.

```go
out, _ := execx.Command("go", "env", "GOOS").OutputBytes()
fmt.Println(len(out) > 0)
// #bool true
```

### <a id="outputtrimmed"></a>OutputTrimmed

OutputTrimmed executes the command and returns trimmed stdout and any error.

```go
out, _ := execx.Command("go", "env", "GOOS").OutputTrimmed()
fmt.Println(out != "")
// #bool true
```

### <a id="run"></a>Run

Run executes the command and returns the result and any error.

```go
res, _ := execx.Command("go", "env", "GOOS").Run()
fmt.Println(res.ExitCode == 0)
// #bool true
```

### <a id="start"></a>Start

Start executes the command asynchronously.

```go
proc := execx.Command("go", "env", "GOOS").Start()
res, _ := proc.Wait()
fmt.Println(res.ExitCode == 0)
// #bool true
```

## Input

### <a id="stdinbytes"></a>StdinBytes

StdinBytes sets stdin from bytes.

```go
out, _ := execx.Command("cat").
	StdinBytes([]byte("hi")).
	Output()
fmt.Println(out)
// #string hi
```

### <a id="stdinfile"></a>StdinFile

StdinFile sets stdin from a file.

```go
file, _ := os.CreateTemp("", "execx-stdin")
_, _ = file.WriteString("hi")
_, _ = file.Seek(0, 0)
out, _ := execx.Command("cat").
	StdinFile(file).
	Output()
fmt.Println(out)
// #string hi
```

### <a id="stdinreader"></a>StdinReader

StdinReader sets stdin from an io.Reader.

```go
out, _ := execx.Command("cat").
	StdinReader(strings.NewReader("hi")).
	Output()
fmt.Println(out)
// #string hi
```

### <a id="stdinstring"></a>StdinString

StdinString sets stdin from a string.

```go
out, _ := execx.Command("cat").
	StdinString("hi").
	Output()
fmt.Println(out)
// #string hi
```

## OS Controls

### <a id="creationflags"></a>CreationFlags

CreationFlags is a no-op on non-Windows platforms.

_Example: creation flags_

```go
fmt.Println(execx.Command("go", "env", "GOOS").CreationFlags(0) != nil)
// #bool true
```

_Example: creation flags_

```go
fmt.Println(execx.Command("go", "env", "GOOS").CreationFlags(0) != nil)
// #bool true
```

### <a id="hidewindow"></a>HideWindow

HideWindow is a no-op on non-Windows platforms.

_Example: hide window_

```go
fmt.Println(execx.Command("go", "env", "GOOS").HideWindow(true) != nil)
// #bool true
```

_Example: hide window_

```go
fmt.Println(execx.Command("go", "env", "GOOS").HideWindow(true) != nil)
// #bool true
```

### <a id="pdeathsig"></a>Pdeathsig

Pdeathsig sets a parent-death signal on Linux.

_Example: pdeathsig_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Pdeathsig(0) != nil)
// #bool true
```

_Example: pdeathsig_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Pdeathsig(0) != nil)
// #bool true
```

_Example: pdeathsig_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Pdeathsig(0) != nil)
// #bool true
```

### <a id="setpgid"></a>Setpgid

Setpgid sets the process group ID behavior.

_Example: setpgid_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Setpgid(true) != nil)
// #bool true
```

_Example: setpgid_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Setpgid(true) != nil)
// #bool true
```

_Example: setpgid_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Setpgid(true) != nil)
// #bool true
```

### <a id="setsid"></a>Setsid

Setsid sets the session ID behavior.

_Example: setsid_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Setsid(true) != nil)
// #bool true
```

_Example: setsid_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Setsid(true) != nil)
// #bool true
```

_Example: setsid_

```go
fmt.Println(execx.Command("go", "env", "GOOS").Setsid(true) != nil)
// #bool true
```

## Pipelining

### <a id="pipe"></a>Pipe

Pipe appends a new command to the pipeline. Pipelines run on all platforms.

```go
out, _ := execx.Command("printf", "go").
	Pipe("tr", "a-z", "A-Z").
	OutputTrimmed()
fmt.Println(out)
// #string GO
```

### <a id="pipebesteffort"></a>PipeBestEffort

PipeBestEffort sets best-effort pipeline semantics (run all stages, surface the first error).

```go
res, err := execx.Command("false").
	Pipe("printf", "ok").
	PipeBestEffort().
	Run()
fmt.Println(err == nil && res.Stdout == "ok")
// #bool true
```

### <a id="pipestrict"></a>PipeStrict

PipeStrict sets strict pipeline semantics (stop on first failure).

```go
res, _ := execx.Command("false").
	Pipe("printf", "ok").
	PipeStrict().
	Run()
fmt.Println(res.ExitCode != 0)
// #bool true
```

### <a id="pipelineresults"></a>PipelineResults

PipelineResults executes the command and returns per-stage results and any error.

```go
results, err := execx.Command("printf", "go").
	Pipe("tr", "a-z", "A-Z").
	PipelineResults()
fmt.Println(err == nil && len(results) == 2)
// #bool true
```

## Process

### <a id="gracefulshutdown"></a>GracefulShutdown

GracefulShutdown sends a signal and escalates to kill after the timeout.

```go
proc := execx.Command("sleep", "2").Start()
_ = proc.GracefulShutdown(os.Interrupt, 100*time.Millisecond)
res, err := proc.Wait()
fmt.Println(err != nil || res.ExitCode != 0)
// #bool true
```

### <a id="interrupt"></a>Interrupt

Interrupt sends an interrupt signal to the process.

```go
proc := execx.Command("sleep", "2").Start()
_ = proc.Interrupt()
res, err := proc.Wait()
fmt.Println(err != nil || res.ExitCode != 0)
// #bool true
```

### <a id="killafter"></a>KillAfter

KillAfter terminates the process after the given duration.

```go
proc := execx.Command("sleep", "2").Start()
proc.KillAfter(100 * time.Millisecond)
res, err := proc.Wait()
fmt.Println(err != nil || res.ExitCode != 0)
// #bool true
```

### <a id="send"></a>Send

Send sends a signal to the process.

```go
proc := execx.Command("sleep", "2").Start()
_ = proc.Send(os.Interrupt)
res, err := proc.Wait()
fmt.Println(err != nil || res.ExitCode != 0)
// #bool true
```

### <a id="terminate"></a>Terminate

Terminate kills the process immediately.

```go
proc := execx.Command("sleep", "2").Start()
_ = proc.Terminate()
res, err := proc.Wait()
fmt.Println(err != nil || res.ExitCode != 0)
// #bool true
```

### <a id="wait"></a>Wait

Wait waits for the command to complete and returns the result and any error.

```go
proc := execx.Command("go", "env", "GOOS").Start()
res, _ := proc.Wait()
fmt.Println(res.ExitCode == 0)
// #bool true
```

## Results

### <a id="isexitcode"></a>IsExitCode

IsExitCode reports whether the exit code matches.

```go
res, _ := execx.Command("go", "env", "GOOS").Run()
fmt.Println(res.IsExitCode(0))
// #bool true
```

### <a id="issignal"></a>IsSignal

IsSignal reports whether the command terminated due to a signal.

```go
res, _ := execx.Command("go", "env", "GOOS").Run()
fmt.Println(res.IsSignal(os.Interrupt))
// false
```

### <a id="ok"></a>OK

OK reports whether the command exited cleanly without errors.

```go
res, _ := execx.Command("go", "env", "GOOS").Run()
fmt.Println(res.OK())
// #bool true
```

## Streaming

### <a id="onstderr"></a>OnStderr

OnStderr registers a line callback for stderr.

```go
_, err := execx.Command("go", "env", "-badflag").
	OnStderr(func(line string) {
		fmt.Println(line)
	}).
	Run()
fmt.Println(err == nil)
// flag provided but not defined: -badflag
// usage: go env [-json] [-changed] [-u] [-w] [var ...]
// Run 'go help env' for details.
// false
```

### <a id="onstdout"></a>OnStdout

OnStdout registers a line callback for stdout.

```go
_, _ = execx.Command("printf", "hi\n").
	OnStdout(func(line string) { fmt.Println(line) }).
	Run()
// hi
```

### <a id="stderrwriter"></a>StderrWriter

StderrWriter sets a raw writer for stderr.

```go
var out strings.Builder
_, err := execx.Command("go", "env", "-badflag").
	StderrWriter(&out).
	Run()
fmt.Print(out.String())
fmt.Println(err == nil)
// flag provided but not defined: -badflag
// usage: go env [-json] [-changed] [-u] [-w] [var ...]
// Run 'go help env' for details.
// true
```

### <a id="stdoutwriter"></a>StdoutWriter

StdoutWriter sets a raw writer for stdout.

```go
var out strings.Builder
_, err := execx.Command("go", "env", "GOOS").
	StdoutWriter(&out).
	Run()
fmt.Println(err == nil && out.Len() > 0)
// #bool true
```

## WorkingDir

### <a id="dir"></a>Dir

Dir sets the working directory.

```go
dir := os.TempDir()
out, _ := execx.Command("pwd").
	Dir(dir).
	OutputTrimmed()
fmt.Println(out == dir)
// #bool true
```
<!-- api:embed:end -->
