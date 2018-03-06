package rfc7464

import (
	"encoding/json"
	"io"
)

// An Encoder writes JSON Text Sequences to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// Encode writes v as JSON Text Sequence (RFC 7464) to the stream.
func (enc *Encoder) Encode(v interface{}) error {
	// https://tools.ietf.org/html/rfc7464 JavaScript Object Notation (JSON) Text Sequences
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	// record separator
	if _, err = enc.w.Write([]byte("\u001e")); err != nil {
		return err
	}
	if _, err = enc.w.Write(bytes); err != nil {
		return err
	}
	// line feed
	if _, err = enc.w.Write([]byte("\u000a")); err != nil {
		return err
	}
	return nil
}
