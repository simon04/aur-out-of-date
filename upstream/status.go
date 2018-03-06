package upstream

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var statusWriter io.Writer = os.Stdout

// UpToDate means that the packaged version matches the upstream version
const UpToDate = "UP-TO-DATE"

// FlaggedOutOfDate means that package is outdated and flagged
const FlaggedOutOfDate = "FLAGGED-OUT-OF-DATE"

// OutOfDate means that package is outdated
const OutOfDate = "OUT-OF-DATE"

// Unknown represents an unknown upstream version
const Unknown = "UNKNOWN"

// Status holds the packaged and upstream version for a package
type Status struct {
	Package          string  `json:"name"`
	Message          string  `json:"message"`
	FlaggedOutOfDate bool    `json:"flagged,omitempty"`
	Version          string  `json:"version,omitempty"`
	Upstream         Version `json:"upstream,omitempty"`
	Status           string  `json:"status"`
}

// Print displays the status on the console
func (s *Status) Print() {
	ansiColor := ""
	switch s.Status {
	case UpToDate:
		ansiColor = "\x1b[32m"
	case FlaggedOutOfDate:
		ansiColor = "\x1b[31m"
	case OutOfDate:
		ansiColor = "\x1b[31m"
	default:
		ansiColor = "\x1b[37m"
	}
	fmt.Fprintf(statusWriter, "%s%22s [%s] %s \x1b[0m\n", ansiColor, "["+s.Status+"]", s.Package, s.Message)
}

// PrintJSONTextSequence outputs the status as JSON Text Sequences (RFC 7464)
func (s *Status) PrintJSONTextSequence() {
	// https://tools.ietf.org/html/rfc7464 JavaScript Object Notation (JSON) Text Sequences
	statusWriter.Write([]byte("\u001e")) // record separator
	bytes, _ := json.Marshal(s)
	statusWriter.Write(bytes)
	statusWriter.Write([]byte("\u000a")) // line feed
}
