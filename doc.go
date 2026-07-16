// Package execx provides explicit, fluent subprocess execution without shell
// parsing or interpolation.
//
// A Cmd is a mutable configuration value and must not be changed or executed
// concurrently. Each execution creates fresh os/exec commands, captures stdout
// and stderr, and reports process status in Result. A non-zero child exit is
// represented by Result.ExitCode rather than an error; errors are reserved for
// failures to start, cancellation, and I/O failures.
//
// Pipelines start every stage concurrently. Stage configuration methods apply
// to the stage on which they are called, while PipeStrict, PipeBestEffort,
// WithPTY, and shadow-print settings apply to the full chain.
package execx
