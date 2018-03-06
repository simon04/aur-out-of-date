package upstream

import (
	"fmt"
	"io"
	"os"
)

var statisticsWriter io.Writer = os.Stdout

// Statistics collects package status
type Statistics struct {
	UpToDate         int
	FlaggedOutOfDate int
	OutOfDate        int
	Unknown          int
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
