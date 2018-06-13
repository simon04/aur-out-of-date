package status

import (
	"bytes"
	"strings"
	"testing"
)

var stat = Statistics{
	UpToDate:         2,
	FlaggedOutOfDate: 3,
	OutOfDate:        5,
	Unknown:          7,
}

func TestStatisticsOutput(t *testing.T) {
	out := bytes.NewBuffer(nil)
	statisticsWriter = out
	stat.Print()
	actual := string(out.Bytes())

	if !strings.Contains(actual, "[UP-TO-DATE] 2") ||
		!strings.Contains(actual, "[FLAGGED-OUT-OF-DATE] 3") ||
		!strings.Contains(actual, "[OUT-OF-DATE] 5") ||
		!strings.Contains(actual, "[UNKNOWN] 7") ||
		!strings.Contains(actual, "[TOTAL] 17") {
		t.Errorf("Unexpected '%s'", actual)
	}
}

func TestStatisticsJSONOutput(t *testing.T) {
	out := bytes.NewBuffer(nil)
	statisticsWriter = out
	stat.PrintJSONTextSequence()
	actual := string(out.Bytes())
	expected := "\u001e" + `{"type":"statistics","up_to_date":2,"flagged_out_of_date":3,"out_of_date":5,"unknown":7}` + "\u000a"
	if actual != expected {
		t.Errorf("Expecting '%s', but got '%s'", expected, actual)
	}
}
