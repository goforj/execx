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
    <img src="https://img.shields.io/badge/tests-98-brightgreen" alt="Tests">
<!-- test-count:embed:end -->
    <a href="https://goreportcard.com/report/github.com/goforj/execx"><img src="https://goreportcard.com/badge/github.com/goforj/execx" alt="Go Report Card"></a>
</p>

It provides a clean, composable API for running system commands without sacrificing control, correctness, or transparency.  
No magic. No hidden behavior. Just a better way to work with processes.

## Why execx?

The standard library’s `os/exec` package is powerful, but verbose and easy to misuse.  
`execx` keeps the same underlying model while making the common cases obvious and safe.

**execx is for you if you want:**

- Clear, chainable command construction
- Predictable execution semantics
- Explicit control over arguments, environment, and I/O
- Zero shell interpolation or magic
- A small, auditable API surface

## Installation

```bash
go get github.com/goforj/execx
````

## Quick Start

```go
out, err := execx.
    Command("git", "status").
    Output()

fmt.Println(out)
```

Or with structured execution:

```go
res := execx.Command("ls", "-la").Run()

if res.Err != nil {
    log.Fatal(res.Err)
}

fmt.Println(res.Stdout)
```

## Fluent Command Construction

Commands are built fluently and executed explicitly.

```go
cmd := execx.
    Command("docker", "run").
    Arg("--rm").
    Arg("-p", "8080:80").
    Arg("nginx")
```

Nothing is executed until you call `Run`, `Output`, or `Start`.

## Argument Handling

Arguments are appended deterministically and never shell-expanded.

```go
cmd.Arg("--env", "PROD")
cmd.Arg(map[string]string{"--name": "api"})
```

This guarantees predictable behavior across platforms.

## Execution Modes

### Run

Execute and return a structured result:

```go
res := cmd.Run()
```

### Output Variants

```go
out, err := cmd.Output()
out, err := cmd.OutputBytes()
out, err := cmd.OutputTrimmed()
out, err := cmd.CombinedOutput()
```

### Output

Return stdout directly:

```go
out, err := cmd.Output()
```

### Start (async)

```go
proc := cmd.Start()
proc.Wait()
proc.KillAfter(5 * time.Second)
```

## Result Object

Every execution returns a `Result`:

```go
type Result struct {
    Stdout   string
    Stderr   string
    ExitCode int
    Err      error
    Duration time.Duration
}
```

* Non-zero exit codes do **not** imply failure
* `Err` indicates execution failure (spawn, context, signal)

## Pipelining

Chain commands safely:

```go
out, err := execx.
    Command("ps", "aux").
    Pipe("grep", "nginx").
    Pipe("awk", "{print $2}").
    Output()
```

Pipelines are explicit and deterministic.

```go
cmd := execx.Command("ps", "aux").Pipe("grep", "nginx")
cmd.PipeStrict()     // default
cmd.PipeBestEffort() // returns last stage, surfaces first error
```

## Context & Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

execx.Command("sleep", "5").
    WithContext(ctx).
    Run()

execx.Command("sleep", "5").
    WithTimeout(2 * time.Second).
    Run()

execx.Command("sleep", "5").
    WithDeadline(time.Now().Add(2 * time.Second)).
    Run()
```

## Environment Control

```go
cmd.Env("DEBUG=true")
cmd.Env(map[string]string{"MODE": "prod"})
cmd.EnvOnly(map[string]string{"MODE": "prod"})
cmd.EnvInherit()
cmd.EnvAppend(map[string]string{"DEBUG": "1"})
```

## Streaming Output

```go
cmd.
    OnStdout(func(line string) {
        fmt.Println("OUT:", line)
    }).
    OnStderr(func(line string) {
        fmt.Println("ERR:", line)
    }).
    Run()
```

## Raw Writers

```go
cmd.StdoutWriter(os.Stdout)
cmd.StderrWriter(os.Stderr)
```

## Exit Handling

```go
if res.IsExitCode(1) {
    log.Println("Command failed")
}
```

## Debugging Helpers

```go
cmd.Args()
cmd.EnvList()
cmd.ShellEscaped()
```

## Design Principles

* **Explicit over implicit**
* **No hidden behavior**
* **No shell magic**
* **Composable over clever**
* **Predictable over flexible**

`execx` is intentionally boring — in the best possible way.

## Non-Goals

* Shell scripting replacement
* Command parsing or glob expansion
* Task runners or build systems
* Automatic retries or heuristics

## Testing & Reliability

* 100% public API coverage
* Deterministic behavior
* No global state
* Safe for concurrent read-only use; mutation during execution is undefined

## Runnable examples

Every function has a corresponding runnable example under [`./examples`](./examples).

These examples are **generated directly from the documentation blocks** of each function, ensuring the docs and code never drift. These are the same examples you see here in the README and GoDoc.

An automated test executes **every example** to verify it builds and runs successfully.

This guarantees all examples are valid, up-to-date, and remain functional as the API evolves.

<!-- api:embed:start -->

## API Index

| Group | Functions |
|------:|-----------|
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
res := execx.Command("go", "env", "GOOS").WithContext(ctx).Run()
fmt.Println(res.Err == nil)
// #bool true
```

### <a id="withdeadline"></a>WithDeadline

WithDeadline binds the command to a deadline.

```go
res := execx.Command("go", "env", "GOOS").WithDeadline(time.Now().Add(2 * time.Second)).Run()
fmt.Println(res.Err == nil)
// #bool true
```

### <a id="withtimeout"></a>WithTimeout

WithTimeout binds the command to a timeout.

```go
res := execx.Command("go", "env", "GOOS").WithTimeout(2 * time.Second).Run()
fmt.Println(res.Err == nil)
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

CombinedOutput executes the command and returns stdout+stderr.

```go
out, _ := execx.Command("go", "env", "GOOS").CombinedOutput()
fmt.Println(out != "")
// #bool true
```

### <a id="output"></a>Output

Output executes the command and returns stdout.

```go
out, _ := execx.Command("go", "env", "GOOS").Output()
fmt.Println(out != "")
// #bool true
```

### <a id="outputbytes"></a>OutputBytes

OutputBytes executes the command and returns stdout bytes.

```go
out, _ := execx.Command("go", "env", "GOOS").OutputBytes()
fmt.Println(len(out) > 0)
// #bool true
```

### <a id="outputtrimmed"></a>OutputTrimmed

OutputTrimmed executes the command and returns trimmed stdout.

```go
out, _ := execx.Command("go", "env", "GOOS").OutputTrimmed()
fmt.Println(out != "")
// #bool true
```

### <a id="run"></a>Run

Run executes the command and returns the result.

```go
res := execx.Command("go", "env", "GOOS").Run()
fmt.Println(res.Stdout)
// darwin (or linux, windows, etc.)
```

### <a id="start"></a>Start

Start executes the command asynchronously.

```go
proc := execx.Command("go", "env", "GOOS").Start()
res := proc.Wait()
fmt.Println(res.ExitCode == 0)
// #bool true
```

## Input

### <a id="stdinbytes"></a>StdinBytes

StdinBytes sets stdin from bytes.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	buf := make([]byte, 8)
	n, _ := os.Stdin.Read(buf)
	_, _ = os.Stdout.Write(buf[:n])
	return
}
out, _ := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	StdinBytes([]byte("hi")).
	Output()
fmt.Println(out == "hi")
// #bool true
```

### <a id="stdinfile"></a>StdinFile

StdinFile sets stdin from a file.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	buf := make([]byte, 8)
	n, _ := os.Stdin.Read(buf)
	_, _ = os.Stdout.Write(buf[:n])
	return
}
file, _ := os.CreateTemp("", "execx-stdin")
_, _ = file.WriteString("hi")
_, _ = file.Seek(0, 0)
out, _ := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	StdinFile(file).
	Output()
fmt.Println(out == "hi")
// #bool true
```

### <a id="stdinreader"></a>StdinReader

StdinReader sets stdin from an io.Reader.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	buf := make([]byte, 8)
	n, _ := os.Stdin.Read(buf)
	_, _ = os.Stdout.Write(buf[:n])
	return
}
out, _ := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	StdinReader(strings.NewReader("hi")).
	Output()
fmt.Println(out == "hi")
// #bool true
```

### <a id="stdinstring"></a>StdinString

StdinString sets stdin from a string.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	buf := make([]byte, 8)
	n, _ := os.Stdin.Read(buf)
	_, _ = os.Stdout.Write(buf[:n])
	return
}
out, _ := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	StdinString("hi").
	Output()
fmt.Println(out == "hi")
// #bool true
```

## OS Controls

### <a id="creationflags"></a>CreationFlags

CreationFlags sets Windows creation flags.

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

HideWindow controls window visibility.

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

Pdeathsig is a no-op on Windows.

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

Setpgid is a no-op on Windows.

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

Setsid is a no-op on Windows.

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

Pipe appends a new command to the pipeline.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	mode := os.Getenv("EXECX_EXAMPLE_MODE")
	if mode == "emit" {
		fmt.Print("ok")
		return
	}
	if mode == "echo" {
		buf := make([]byte, 8)
		n, _ := os.Stdin.Read(buf)
		_, _ = os.Stdout.Write(buf[:n])
		return
	}
}
out, _ := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1", "EXECX_EXAMPLE_MODE=emit").
	Pipe(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1", "EXECX_EXAMPLE_MODE=echo").
	Output()
fmt.Println(out == "ok")
// #bool true
```

### <a id="pipebesteffort"></a>PipeBestEffort

PipeBestEffort sets best-effort pipeline semantics.

```go
cmd := execx.Command("go", "env", "GOOS").PipeBestEffort()
fmt.Println(cmd != nil)
// #bool true
```

### <a id="pipestrict"></a>PipeStrict

PipeStrict sets strict pipeline semantics.

```go
cmd := execx.Command("go", "env", "GOOS").PipeStrict()
fmt.Println(cmd != nil)
// #bool true
```

### <a id="pipelineresults"></a>PipelineResults

PipelineResults executes the command and returns per-stage results.

```go
results := execx.Command("go", "env", "GOOS").PipelineResults()
fmt.Println(len(results) == 1)
// #bool true
```

## Process

### <a id="gracefulshutdown"></a>GracefulShutdown

GracefulShutdown sends a signal and escalates to kill after the timeout.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	time.Sleep(2 * time.Second)
	return
}
proc := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	Start()
_ = proc.GracefulShutdown(os.Interrupt, 100*time.Millisecond)
res := proc.Wait()
fmt.Println(res.ExitCode != 0)
// #bool true
```

### <a id="interrupt"></a>Interrupt

Interrupt sends an interrupt signal to the process.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	time.Sleep(2 * time.Second)
	return
}
proc := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	Start()
_ = proc.Interrupt()
res := proc.Wait()
fmt.Println(res.ExitCode != 0)
// #bool true
```

### <a id="killafter"></a>KillAfter

KillAfter terminates the process after the given duration.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	time.Sleep(2 * time.Second)
	return
}
proc := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	Start()
proc.KillAfter(100 * time.Millisecond)
res := proc.Wait()
fmt.Println(res.ExitCode != 0)
// #bool true
```

### <a id="send"></a>Send

Send sends a signal to the process.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	time.Sleep(2 * time.Second)
	return
}
proc := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	Start()
_ = proc.Send(os.Interrupt)
res := proc.Wait()
fmt.Println(res.ExitCode != 0)
// #bool true
```

### <a id="terminate"></a>Terminate

Terminate kills the process immediately.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	time.Sleep(2 * time.Second)
	return
}
proc := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	Start()
_ = proc.Terminate()
res := proc.Wait()
fmt.Println(res.ExitCode != 0)
// #bool true
```

### <a id="wait"></a>Wait

Wait waits for the command to complete and returns the result.

```go
proc := execx.Command("go", "env", "GOOS").Start()
res := proc.Wait()
fmt.Println(res.ExitCode == 0)
// #bool true
```

## Results

### <a id="isexitcode"></a>IsExitCode

IsExitCode reports whether the exit code matches.

```go
res := execx.Command("go", "env", "GOOS").Run()
fmt.Println(res.IsExitCode(0))
// #bool true
```

### <a id="issignal"></a>IsSignal

IsSignal reports whether the command terminated due to a signal.

```go
res := execx.Command("go", "env", "GOOS").Run()
fmt.Println(res.IsSignal(os.Interrupt))
// #bool false
```

### <a id="ok"></a>OK

OK reports whether the command exited cleanly without errors.

```go
res := execx.Command("go", "env", "GOOS").Run()
fmt.Println(res.OK())
// #bool true
```

## Streaming

### <a id="onstderr"></a>OnStderr

OnStderr registers a line callback for stderr.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	_, _ = os.Stderr.WriteString("err\n")
	return
}
var lines []string
execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	OnStderr(func(line string) { lines = append(lines, line) }).
	Run()
fmt.Println(len(lines) == 1)
// #bool true
```

### <a id="onstdout"></a>OnStdout

OnStdout registers a line callback for stdout.

```go
var lines []string
execx.Command("go", "env", "GOOS").
	OnStdout(func(line string) { lines = append(lines, line) }).
	Run()
fmt.Println(len(lines) > 0)
// #bool true
```

### <a id="stderrwriter"></a>StderrWriter

StderrWriter sets a raw writer for stderr.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	_, _ = os.Stderr.WriteString("err\n")
	return
}
var out strings.Builder
execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	StderrWriter(&out).
	Run()
fmt.Println(out.Len() > 0)
// #bool true
```

### <a id="stdoutwriter"></a>StdoutWriter

StdoutWriter sets a raw writer for stdout.

```go
var out strings.Builder
execx.Command("go", "env", "GOOS").
	StdoutWriter(&out).
	Run()
fmt.Println(out.Len() > 0)
// #bool true
```

## WorkingDir

### <a id="dir"></a>Dir

Dir sets the working directory.

```go
if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	return
}
dir := os.TempDir()
out, _ := execx.Command(os.Args[0]).
	Env("EXECX_EXAMPLE_CHILD=1").
	Dir(dir).
	OutputTrimmed()
fmt.Println(out == dir)
// #bool true
```
<!-- api:embed:end -->
