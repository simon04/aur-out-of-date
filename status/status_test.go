package status

import (
	"bytes"
	"testing"

	"github.com/simon04/aur-out-of-date/upstream"
)

var s = Status{
	Package: "spectre-meltdown-checker",
	Version: "0.35-1",
}

func TestCompare(t *testing.T) {
	s.Compare(upstream.Version("0.35"))
	if s.Status != UpToDate {
		t.Errorf("Expecting status to be %s", UpToDate)
	}
	s.Compare(upstream.Version("0.37"))
	if s.Status != OutOfDate {
		t.Errorf("Expecting status to be %s", OutOfDate)
	}
	s.Compare(upstream.Version("0.30"))
	if s.Status != Unknown {
		t.Errorf("Expecting status to be %s", Unknown)
	}
	s.Compare(upstream.Version("foo"))
	if s.Status != Unknown {
		t.Errorf("Expecting status to be %s", Unknown)
	}
}

func TestStatusOutput(t *testing.T) {
	s.Compare(upstream.Version("0.35"))
	out := bytes.NewBuffer(nil)
	statusWriter = out
	s.Print()
	actual := string(out.Bytes())
	expected := "\x1b[32m          [UP-TO-DATE] [spectre-meltdown-checker][0.35-1] matches upstream version 0.35 \x1b[0m\n"
	if actual != expected {
		t.Errorf("Expecting '%s', but got '%s'", expected, actual)
	}
}

func TestStatusJSONOutput(t *testing.T) {
	s.Compare(upstream.Version("0.35"))
	out := bytes.NewBuffer(nil)
	statusWriter = out
	s.PrintJSONTextSequence()
	actual := string(out.Bytes())
	expected := "\u001e" + `{"type":"package","name":"spectre-meltdown-checker","message":"matches upstream version 0.35","version":"0.35-1","upstream":"0.35","status":"UP-TO-DATE"}` + "\u000a"
	if actual != expected {
		t.Errorf("Expecting '%s', but got '%s'", expected, actual)
	}
}
