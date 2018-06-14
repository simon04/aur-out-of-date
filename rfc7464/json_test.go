package rfc7464

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	out := bytes.NewBuffer(nil)
	NewEncoder(out).Encode(map[string]int{
		"foo": 1, "bar": 2, "baz": 3})
	actual := string(out.Bytes())
	expected := "\u001e" + "{\"bar\":2,\"baz\":3,\"foo\":1}" + "\u000a"
	if actual != expected {
		t.Errorf("Expecting '%s', but got '%s'", expected, actual)
	}
}
