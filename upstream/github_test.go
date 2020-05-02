package upstream

import (
	"testing"
)

func TestParseGitHub(t *testing.T) {
	g := parseGitHub("https://github.com/foo/bar")
	if g.owner != "foo" || g.repository != "bar" {
		t.Errorf("Expecting foo/bar, but got %v", g.String())
	}
	g = parseGitHub("https://foo.github.io/bar/...")
	if g.owner != "foo" || g.repository != "bar" {
		t.Errorf("Expecting foo/bar, but got %v", g.String())
	}
}
