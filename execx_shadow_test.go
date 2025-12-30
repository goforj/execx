package execx

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
)

var stderrMu sync.Mutex

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

func stripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

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

func TestShadowPrintPrefix(t *testing.T) {
	out := captureStderr(t, func() {
		_, _ = Command("printf", "hi").ShadowPrintPrefix("run").Run()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "run > printf hi") {
		t.Fatalf("expected prefix, got %q", plain)
	}
}

func TestShadowPrintOff(t *testing.T) {
	out := captureStderr(t, func() {
		_, _ = Command("printf", "hi").ShadowPrint().ShadowPrintOff().Run()
	})
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected no output, got %q", out)
	}
}

func TestShadowPrintMask(t *testing.T) {
	out := captureStderr(t, func() {
		mask := func(cmd string) string {
			return strings.ReplaceAll(cmd, "secret", "***")
		}
		_, _ = Command("printf", "secret").ShadowPrintMask(mask).Run()
	})
	plain := stripANSI(out)
	if !strings.Contains(plain, "printf ***") {
		t.Fatalf("expected masked output, got %q", plain)
	}
}

func TestShadowPrintFormatter(t *testing.T) {
	out := captureStderr(t, func() {
		formatter := func(ev ShadowEvent) string {
			return "shadow:" + string(ev.Phase) + ":" + ev.RawCommand
		}
		_, _ = Command("printf", "hi").ShadowPrintFormatter(formatter).Run()
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

func TestShadowPrintFormatterEmpty(t *testing.T) {
	out := captureStderr(t, func() {
		formatter := func(ev ShadowEvent) string {
			return ""
		}
		_, _ = Command("printf", "hi").ShadowPrintFormatter(formatter).Run()
	})
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected no output, got %q", out)
	}
}

func TestShadowCommandPipeline(t *testing.T) {
	cmd := Command("printf", "go").Pipe("tr", "a-z", "A-Z")
	if got := cmd.shadowCommand(); got != "printf go | tr a-z A-Z" {
		t.Fatalf("unexpected shadow command: %q", got)
	}
}

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

func TestShadowPrintLineNil(t *testing.T) {
	shadowPrintLine(nil, ShadowBefore, 0, false)
}
