package upstream

import (
	"fmt"
	"io"
	"os"

	rfc7464 "github.com/simon04/aur-out-of-date/rfc7464go"
)

var statisticsWriter io.Writer = os.Stdout

// Statistics collects package status
type Statistics struct {
	Type             string `json:"type"`
	UpToDate         int    `json:"up_to_date"`
	FlaggedOutOfDate int    `json:"flagged_out_of_date"`
	OutOfDate        int    `json:"out_of_date"`
	Unknown          int    `json:"unknown"`
}

// Print displays the statistics on the console
func (s *Statistics) Print() {
	fmt.Fprintln(statisticsWriter)
	fmt.Fprintln(statisticsWriter, "STATISTICS")
	fmt.Fprintf(statisticsWriter, "%s%22s %d \x1b[0m\n", UpToDate.color(), "["+UpToDate+"]", s.UpToDate)
	fmt.Fprintf(statisticsWriter, "%s%22s %d \x1b[0m\n", FlaggedOutOfDate.color(), "["+FlaggedOutOfDate+"]", s.FlaggedOutOfDate)
	fmt.Fprintf(statisticsWriter, "%s%22s %d \x1b[0m\n", OutOfDate.color(), "["+OutOfDate+"]", s.OutOfDate)
	fmt.Fprintf(statisticsWriter, "%s%22s %d \x1b[0m\n", Unknown.color(), "["+Unknown+"]", s.Unknown)
	fmt.Fprintf(statisticsWriter, "%s%22s %d \x1b[0m\n", Unknown.color(), "[TOTAL]", s.UpToDate+s.FlaggedOutOfDate+s.OutOfDate+s.Unknown)
}

// PrintJSONTextSequence outputs the statistics as JSON Text Sequences (RFC 7464)
func (s *Statistics) PrintJSONTextSequence() {
	s.Type = "statistics"
	rfc7464.NewEncoder(statisticsWriter).Encode(s)
}
