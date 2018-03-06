package upstream

import (
	"bytes"
	"testing"
)

var s = Status{
	Package:  "spectre-meltdown-checker",
	Message:  "Package spectre-meltdown-checker 0.35-1 matches upstream version 0.35",
	Version:  "0.35-1",
	Upstream: "0.35",
	Status:   "UP-TO-DATE",
}

func TestStatusOutput(t *testing.T) {
	out := bytes.NewBuffer(nil)
	statusWriter = out
	s.Print()
	actual := string(out.Bytes())
	expected := "\x1b[32m          [UP-TO-DATE] [spectre-meltdown-checker] Package spectre-meltdown-checker 0.35-1 matches upstream version 0.35 \x1b[0m\n"
	if actual != expected {
		t.Errorf("Expecting '%s', but got '%s'", expected, actual)
	}
}

func TestStatusJSONOutput(t *testing.T) {
	out := bytes.NewBuffer(nil)
	statusWriter = out
	s.PrintJSONTextSequence()
	actual := string(out.Bytes())
	expected := "\u001e" + `{"name":"spectre-meltdown-checker","message":"Package spectre-meltdown-checker 0.35-1 matches upstream version 0.35","version":"0.35-1","upstream":"0.35","status":"UP-TO-DATE"}` + "\u000a"
	if actual != expected {
		t.Errorf("Expecting '%s', but got '%s'", expected, actual)
	}
}
