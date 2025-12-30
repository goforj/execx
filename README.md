<p align="center">
  <img src="./docs/images/logo.png?v=2" width="400" alt="str logo">
</p>

<p align="center">
    **execx** is an ergonomic, fluent wrapper around Go’s `os/exec` package.
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/goforj/execx"><img src="https://pkg.go.dev/badge/github.com/goforj/execx.svg" alt="Go Reference"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License: MIT"></a>
    <a href="https://github.com/goforj/execx/actions"><img src="https://github.com/goforj/execx/actions/workflows/test.yml/badge.svg" alt="Go Test"></a>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/go-1.18+-blue?logo=go" alt="Go version"></a>
    <img src="https://img.shields.io/github/v/tag/goforj/execx?label=version&sort=semver" alt="Latest tag">
    <a href="https://codecov.io/gh/goforj/execx" ><img src="https://codecov.io/github/goforj/execx/graph/badge.svg?token=9KT46ZORP3"/></a>
<!-- test-count:embed:start -->
    <img src="https://img.shields.io/badge/tests-196-brightgreen" alt="Tests">
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

### Output

Return stdout directly:

```go
out, err := cmd.Output()
```

### Start (async)

```go
proc := cmd.Start()
proc.Wait()
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
* `Err` only indicates execution failure (spawn, context, signal)

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

## Context & Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

execx.Command("sleep", "5").
    WithContext(ctx).
    Run()
```

## Environment Control

```go
cmd.Env("DEBUG=true")
cmd.Env(map[string]string{"MODE": "prod"})
cmd.EnvClear()
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

## Exit Handling

```go
if res.IsExitCode(1) {
    log.Println("Command failed")
}
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
* Safe for concurrent use
