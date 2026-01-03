package execx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

// Decoder decodes serialized data into a destination value.
type Decoder interface {
	Decode(data []byte, dst any) error
}

// DecoderFunc adapts a function to a Decoder.
type DecoderFunc func(data []byte, dst any) error

type decodeSource int

const (
	decodeStdout decodeSource = iota
	decodeStderr
	decodeCombined
)

type decodeConfig struct {
	source decodeSource
	trim   bool
}

func defaultDecodeConfig() decodeConfig {
	return decodeConfig{source: decodeStdout}
}

// DecodeChain configures typed decoding for a command.
type DecodeChain struct {
	cmd     *Cmd
	decoder Decoder
	cfg     decodeConfig
}

// Decode configures a custom decoder for this command.
// Decoding reads from stdout by default; use FromStdout, FromStderr, or FromCombined to select a source.
// @group Decoding
//
// Example: decode custom
//
//	type payload struct {
//		Name string
//	}
//	decoder := execx.DecoderFunc(func(data []byte, dst any) error {
//		out, ok := dst.(*payload)
//		if !ok {
//			return fmt.Errorf("expected *payload")
//		}
//		_, val, ok := strings.Cut(string(data), "=")
//		if !ok {
//			return fmt.Errorf("invalid payload")
//		}
//		out.Name = val
//		return nil
//	})
//	var out payload
//	_ = execx.Command("printf", "name=gopher").
//		Decode(decoder).
//		Into(&out)
//	fmt.Println(out.Name)
//	// #string gopher
func (c *Cmd) Decode(decoder Decoder) *DecodeChain {
	return &DecodeChain{
		cmd:     c,
		decoder: decoder,
		cfg:     defaultDecodeConfig(),
	}
}

// Decode implements Decoder.
func (f DecoderFunc) Decode(data []byte, dst any) error {
	return f(data, dst)
}

// DecodeJSON configures JSON decoding for this command.
// Decoding reads from stdout by default; use FromStdout, FromStderr, or FromCombined to select a source.
// @group Decoding
//
// Example: decode json
//
//	type payload struct {
//		Name string `json:"name"`
//	}
//	var out payload
//	_ = execx.Command("printf", `{"name":"gopher"}`).
//		DecodeJSON().
//		Into(&out)
//	fmt.Println(out.Name)
//	// #string gopher
func (c *Cmd) DecodeJSON() *DecodeChain {
	return c.Decode(DecoderFunc(json.Unmarshal))
}

// DecodeYAML configures YAML decoding for this command.
// Decoding reads from stdout by default; use FromStdout, FromStderr, or FromCombined to select a source.
// @group Decoding
//
// Example: decode yaml
//
//	type payload struct {
//		Name string `yaml:"name"`
//	}
//	var out payload
//	_ = execx.Command("printf", "name: gopher").
//		DecodeYAML().
//		Into(&out)
//	fmt.Println(out.Name)
//	// #string gopher
func (c *Cmd) DecodeYAML() *DecodeChain {
	return c.Decode(DecoderFunc(yaml.Unmarshal))
}

// FromStdout decodes from stdout (default).
// @group Decoding
//
// Example: decode from stdout
//
//	type payload struct {
//		Name string `json:"name"`
//	}
//	var out payload
//	_ = execx.Command("printf", `{"name":"gopher"}`).
//		DecodeJSON().
//		FromStdout().
//		Into(&out)
//	fmt.Println(out.Name)
//	// #string gopher
func (d *DecodeChain) FromStdout() *DecodeChain {
	d.cfg.source = decodeStdout
	return d
}

// FromStderr decodes from stderr.
// @group Decoding
//
// Example: decode from stderr
//
//	type payload struct {
//		Name string `json:"name"`
//	}
//	var out payload
//	_ = execx.Command("sh", "-c", `printf '{"name":"gopher"}' 1>&2`).
//		DecodeJSON().
//		FromStderr().
//		Into(&out)
//	fmt.Println(out.Name)
//	// #string gopher
func (d *DecodeChain) FromStderr() *DecodeChain {
	d.cfg.source = decodeStderr
	return d
}

// FromCombined decodes from combined stdout+stderr.
// @group Decoding
//
// Example: decode combined
//
//	type payload struct {
//		Name string `json:"name"`
//	}
//	var out payload
//	_ = execx.Command("sh", "-c", `printf '{"name":"gopher"}'`).
//		DecodeJSON().
//		FromCombined().
//		Into(&out)
//	fmt.Println(out.Name)
//	// #string gopher
func (d *DecodeChain) FromCombined() *DecodeChain {
	d.cfg.source = decodeCombined
	return d
}

// Trim trims whitespace before decoding.
// @group Decoding
//
// Example: decode trim
//
//	type payload struct {
//		Name string `json:"name"`
//	}
//	var out payload
//	_ = execx.Command("printf", "  {\"name\":\"gopher\"}  ").
//		DecodeJSON().
//		Trim().
//		Into(&out)
//	fmt.Println(out.Name)
//	// #string gopher
func (d *DecodeChain) Trim() *DecodeChain {
	d.cfg.trim = true
	return d
}

// Into executes the command and decodes into dst.
// @group Decoding
//
// Example: decode into
//
//	type payload struct {
//		Name string `json:"name"`
//	}
//	var out payload
//	_ = execx.Command("printf", `{"name":"gopher"}`).
//		DecodeJSON().
//		Into(&out)
//	fmt.Println(out.Name)
//	// #string gopher
func (d *DecodeChain) Into(dst any) error {
	return decodeInto(d.cmd, dst, d.decoder, d.cfg)
}

// DecodeWith executes the command and decodes stdout into dst.
// @group Decoding
//
// Example: decode with
//
//	type payload struct {
//		Name string `json:"name"`
//	}
//	var out payload
//	_ = execx.Command("printf", `{"name":"gopher"}`).
//		DecodeWith(&out, execx.DecoderFunc(json.Unmarshal))
//	fmt.Println(out.Name)
//	// #string gopher
func (c *Cmd) DecodeWith(dst any, decoder Decoder) error {
	return decodeInto(c, dst, decoder, defaultDecodeConfig())
}

func decodeInto(c *Cmd, dst any, decoder Decoder, cfg decodeConfig) error {
	if c == nil {
		return errors.New("command is nil")
	}
	if decoder == nil {
		return errors.New("decoder is nil")
	}
	if dst == nil {
		return errors.New("destination is nil")
	}
	val := reflect.ValueOf(dst)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("destination must be a non-nil pointer")
	}
	data, err := c.decodeSource(cfg.source)
	if err != nil {
		return err
	}
	if cfg.trim {
		data = bytes.TrimSpace(data)
	}
	if err := decoder.Decode(data, dst); err != nil {
		return fmt.Errorf("decode %s: %w", decodeSourceName(cfg.source), err)
	}
	return nil
}

func (c *Cmd) decodeSource(source decodeSource) ([]byte, error) {
	switch source {
	case decodeCombined:
		out, err := c.CombinedOutput()
		if err != nil {
			return nil, err
		}
		return []byte(out), nil
	case decodeStderr:
		res, err := c.Run()
		if err != nil {
			return nil, err
		}
		return []byte(res.Stderr), nil
	case decodeStdout:
		fallthrough
	default:
		return c.OutputBytes()
	}
}

func decodeSourceName(source decodeSource) string {
	switch source {
	case decodeCombined:
		return "combined output"
	case decodeStderr:
		return "stderr"
	default:
		return "stdout"
	}
}
