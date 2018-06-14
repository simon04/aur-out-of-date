package config

import (
	"testing"

	"github.com/simon04/aur-out-of-date/upstream"
)

func TestIsIgnored(t *testing.T) {
	conf := Config{
		Ignore: map[string][]upstream.Version{
			"foo": {"*"},
			"bar": {"0.1", "0.2-alpha.1", "0.2-beta.1"},
			"baz": {},
		},
	}
	if !conf.IsIgnored("foo", "1.0") {
		t.Errorf("foo-1.0 should be ignored")
	}
	if !conf.IsIgnored("foo", "2.0") {
		t.Errorf("foo-2.0 should be ignored")
	}
	if !conf.IsIgnored("foo", "*") {
		t.Errorf("foo-* should be ignored")
	}
	if !conf.IsIgnored("bar", "0.2-beta.1") {
		t.Errorf("bar-0.2-beta.1 should be ignored")
	}
	if conf.IsIgnored("bar", "*") {
		t.Errorf("bar-* should not be ignored")
	}
	if conf.IsIgnored("baz", "1.0") {
		t.Errorf("baz-1.0 should not be ignored")
	}
	if conf.IsIgnored("baz-bin", "1.0") {
		t.Errorf("baz-bin-1.0 should not be ignored")
	}
}
