package execx

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var stderrMu sync.Mutex

// captureStderr serializes global descriptor replacement so shadow tests cannot interfere.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	stderrMu.Lock()
	defer stderrMu.Unlock()

	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stderr = w

	fn()

	_ = w.Close()
	os.Stderr = orig

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	return buf.String()
}

// stripANSI removes presentation codes before assertions inspect shadow text.
func stripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

// TestShadowPrintDefault ensures the default trace includes the command and elapsed duration.
func TestShadowPrintDefault(t *testing.T) {
	out := captureStderr(t, func() {
		_, _ = Command("printf", "hi").ShadowPrint().Run()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "execx > printf hi") {
		t.Fatalf("expected shadow print, got %q", plain)
	}
	if !strings.Contains(plain, "execx > printf hi (") {
		t.Fatalf("expected duration line, got %q", plain)
	}
}

// TestShadowPrintDefaultSpacing ensures command output remains visually separated from trace lines.
func TestShadowPrintDefaultSpacing(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("spacing test uses sh")
	}
	out := captureStderr(t, func() {
		_, _ = Command("sh", "-c", "printf 'hi\\n' 1>&2").
			ShadowPrint().
			StderrWriter(os.Stderr).
			Run()
	})
	plain := stripANSI(out)
	lines := strings.Split(plain, "\n")
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) < 5 {
		t.Fatalf("expected spacing lines, got %q", plain)
	}
	if !strings.Contains(lines[0], "execx > ") {
		t.Fatalf("expected before line, got %q", lines[0])
	}
	if lines[1] != "" {
		t.Fatalf("expected blank line before output, got %q", lines[1])
	}
	if lines[2] != "hi" {
		t.Fatalf("expected output line, got %q", lines[2])
	}
	if lines[3] != "" {
		t.Fatalf("expected blank line after output, got %q", lines[3])
	}
	if !strings.Contains(lines[4], "execx > ") {
		t.Fatalf("expected after line, got %q", lines[4])
	}
}

// TestShadowPrintPrefix ensures callers can identify traces with a custom prefix.
func TestShadowPrintPrefix(t *testing.T) {
	out := captureStderr(t, func() {
		_, _ = Command("printf", "hi").ShadowPrint(WithPrefix("run")).Run()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "run > printf hi") {
		t.Fatalf("expected prefix, got %q", plain)
	}
}

// TestShadowPrintOff ensures disabled tracing produces no diagnostic output.
func TestShadowPrintOff(t *testing.T) {
	out := captureStderr(t, func() {
		_, _ = Command("printf", "hi").ShadowPrint().ShadowOff().Run()
	})
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected no output, got %q", out)
	}
}

// TestShadowPrintMask ensures secrets can be redacted before commands reach diagnostic output.
func TestShadowPrintMask(t *testing.T) {
	out := captureStderr(t, func() {
		mask := func(cmd string) string {
			return strings.ReplaceAll(cmd, "secret", "***")
		}
		_, _ = Command("printf", "secret").ShadowPrint(WithMask(mask)).Run()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "printf ***") {
		t.Fatalf("expected masked output, got %q", plain)
	}
}

// TestShadowPrintFormatter ensures custom formatting receives both lifecycle phases.
func TestShadowPrintFormatter(t *testing.T) {
	out := captureStderr(t, func() {
		formatter := func(ev ShadowEvent) string {
			return "shadow:" + string(ev.Phase) + ":" + ev.RawCommand
		}
		_, _ = Command("printf", "hi").ShadowPrint(WithFormatter(formatter)).Run()
	})
	lines := strings.FieldsFunc(strings.TrimSpace(out), func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), out)
	}
	if !strings.HasPrefix(lines[0], "shadow:before:printf hi") {
		t.Fatalf("unexpected before line: %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "shadow:after:printf hi") {
		t.Fatalf("unexpected after line: %q", lines[1])
	}
}

// TestShadowPrintFormatterEmpty ensures an empty formatted line suppresses that trace event.
func TestShadowPrintFormatterEmpty(t *testing.T) {
	out := captureStderr(t, func() {
		formatter := func(ev ShadowEvent) string {
			return ""
		}
		_, _ = Command("printf", "hi").ShadowPrint(WithFormatter(formatter)).Run()
	})
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected no output, got %q", out)
	}
}

// TestShadowCommandPipeline ensures traces describe the complete pipeline in execution order.
func TestShadowCommandPipeline(t *testing.T) {
	cmd := Command("printf", "go").Pipe("tr", "a-z", "A-Z")
	if got := cmd.shadowCommand(); got != "printf go | tr a-z A-Z" {
		t.Fatalf("unexpected shadow command: %q", got)
	}
}

// TestShadowPrintAsync ensures Start traces are distinguishable from synchronous execution.
func TestShadowPrintAsync(t *testing.T) {
	out := captureStderr(t, func() {
		proc := Command("sleep", "0.01").ShadowPrint().Start()
		_, _ = proc.Wait()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "(async)") {
		t.Fatalf("expected async marker, got %q", plain)
	}
}

// TestShadowOffOnPreservesConfig ensures temporarily disabling tracing does not discard its options.
func TestShadowOffOnPreservesConfig(t *testing.T) {
	out := captureStderr(t, func() {
		cmd := Command("printf", "hi").ShadowPrint(WithPrefix("run"))
		cmd.ShadowOff()
		_, _ = cmd.ShadowOn().Run()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "run > printf hi") {
		t.Fatalf("expected preserved prefix, got %q", plain)
	}
}

// TestShadowOnDefaultConfig ensures enabling an unconfigured command installs safe defaults.
func TestShadowOnDefaultConfig(t *testing.T) {
	out := captureStderr(t, func() {
		cmd := Command("printf", "hi")
		cmd.ShadowOff()
		_, _ = cmd.ShadowOn().Run()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "execx > printf hi") {
		t.Fatalf("expected default prefix, got %q", plain)
	}
}

// TestShadowPrintMaskWithFormatter ensures formatters receive redacted and raw forms for deliberate handling.
func TestShadowPrintMaskWithFormatter(t *testing.T) {
	out := captureStderr(t, func() {
		mask := func(cmd string) string {
			return strings.ReplaceAll(cmd, "secret", "***")
		}
		formatter := func(ev ShadowEvent) string {
			return ev.Command + "|" + ev.RawCommand
		}
		_, _ = Command("printf", "secret").ShadowPrint(WithMask(mask), WithFormatter(formatter)).Run()
	})
	plain := stripANSI(strings.TrimSpace(out))
	if !strings.HasPrefix(plain, "printf ***|printf secret") {
		t.Fatalf("expected masked and raw values, got %q", plain)
	}
}

// TestShadowPrintEmptyPrefix ensures an empty custom prefix falls back to the recognizable default.
func TestShadowPrintEmptyPrefix(t *testing.T) {
	out := captureStderr(t, func() {
		_, _ = Command("printf", "hi").ShadowPrint(WithPrefix("")).Run()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "execx > printf hi") {
		t.Fatalf("expected default prefix, got %q", plain)
	}
}

// TestShadowPrintLineNil ensures absent shadow configuration remains a harmless no-op.
func TestShadowPrintLineNil(t *testing.T) {
	shadowPrintLine(nil, ShadowBefore, 0, false)
}
