package execx

import (
	"encoding/json"
	"errors"
	"runtime"
	"testing"
)

type testPayload struct {
	Name string `json:"name"`
}

func TestDecodeYAMLInto(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("decode yaml test uses printf")
	}
	type payload struct {
		Name string `yaml:"name"`
	}
	var out payload
	if err := Command("printf", "name: gopher").
		DecodeYAML().
		Into(&out); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Name != "gopher" {
		t.Fatalf("unexpected name: %q", out.Name)
	}
}

func TestDecodeFromStdout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("decode stdout test uses printf")
	}
	var out testPayload
	if err := Command("printf", `{"name":"gopher"}`).
		DecodeJSON().
		FromStdout().
		Into(&out); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Name != "gopher" {
		t.Fatalf("unexpected name: %q", out.Name)
	}
}

func TestDecodeNilCmd(t *testing.T) {
	var out testPayload
	err := (*Cmd)(nil).
		DecodeJSON().
		Into(&out)
	if err == nil {
		t.Fatalf("expected error for nil command")
	}
}

func TestDecodeWithJSON(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("output into test uses printf")
	}
	var out testPayload
	err := Command("printf", `{"name":"gopher"}`).
		DecodeWith(&out, DecoderFunc(json.Unmarshal))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Name != "gopher" {
		t.Fatalf("unexpected name: %q", out.Name)
	}
}

func TestDecodeTrim(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("decode trim test uses printf")
	}
	var out testPayload
	err := Command("printf", `  {"name":"gopher"}  `).
		DecodeJSON().
		Trim().
		Into(&out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Name != "gopher" {
		t.Fatalf("unexpected name: %q", out.Name)
	}
}

func TestDecodeCombined(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("combined output test uses sh")
	}
	var out testPayload
	err := Command("sh", "-c", `printf '{"name":"gopher"}'`).
		DecodeJSON().
		FromCombined().
		Into(&out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Name != "gopher" {
		t.Fatalf("unexpected name: %q", out.Name)
	}
}

func TestDecodeCombinedStartError(t *testing.T) {
	var out testPayload
	err := Command("execx-does-not-exist").
		DecodeJSON().
		FromCombined().
		Into(&out)
	if err == nil {
		t.Fatalf("expected error for combined start failure")
	}
}

func TestDecodeCombinedDecodeError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("combined output test uses sh")
	}
	var out testPayload
	err := Command("sh", "-c", "printf not-json").
		DecodeJSON().
		FromCombined().
		Into(&out)
	if err == nil || !strings.Contains(err.Error(), "combined output") {
		t.Fatalf("expected combined output decode error, got %v", err)
	}
}

func TestDecodeStderr(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stderr output test uses sh")
	}
	var out testPayload
	err := Command("sh", "-c", `printf '{"name":"gopher"}' 1>&2`).
		DecodeJSON().
		FromStderr().
		Into(&out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Name != "gopher" {
		t.Fatalf("unexpected name: %q", out.Name)
	}
}

func TestDecodeStderrStartError(t *testing.T) {
	var out testPayload
	err := Command("execx-does-not-exist").
		DecodeJSON().
		FromStderr().
		Into(&out)
	if err == nil {
		t.Fatalf("expected error for stderr start failure")
	}
}

func TestDecodeStderrDecodeError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stderr output test uses sh")
	}
	var out testPayload
	err := Command("sh", "-c", "printf not-json 1>&2").
		DecodeJSON().
		FromStderr().
		Into(&out)
	if err == nil || !strings.Contains(err.Error(), "stderr") {
		t.Fatalf("expected stderr decode error, got %v", err)
	}
}

func TestDecodeWithErrors(t *testing.T) {
	err := Command("printf", `{"name":"gopher"}`).
		DecodeWith(nil, DecoderFunc(json.Unmarshal))
	if err == nil {
		t.Fatalf("expected error for nil destination")
	}

	var out testPayload
	err = Command("printf", `{"name":"gopher"}`).
		DecodeWith(out, DecoderFunc(json.Unmarshal))
	if err == nil {
		t.Fatalf("expected error for non-pointer destination")
	}

	err = Command("printf", `{"name":"gopher"}`).
		DecodeWith(&out, nil)
	if err == nil {
		t.Fatalf("expected error for nil decoder")
	}
}

func TestDecodeErrorWrap(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("decode error test uses printf")
	}
	var out testPayload
	err := Command("printf", `not-json`).
		DecodeJSON().
		Into(&out)
	if err == nil {
		t.Fatalf("expected decode error")
	}
	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) && err.Error() == "" {
		t.Fatalf("expected decode error detail")
	}
}
