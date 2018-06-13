package status

import (
	"fmt"
	"io"
	"os"

	"github.com/simon04/aur-out-of-date/rfc7464"
	"github.com/simon04/aur-out-of-date/upstream"
)

var statusWriter io.Writer = os.Stdout

// StatusType represens the package up-to-date state
type StatusType string

// UpToDate means that the packaged version matches the upstream version
const UpToDate = StatusType("UP-TO-DATE")

// FlaggedOutOfDate means that package is outdated and flagged
const FlaggedOutOfDate = StatusType("FLAGGED-OUT-OF-DATE")

// OutOfDate means that package is outdated
const OutOfDate = StatusType("OUT-OF-DATE")

// Unknown represents an unknown upstream version
const Unknown = StatusType("UNKNOWN")

// Status holds the packaged and upstream version for a package
type Status struct {
	Type             string           `json:"type"`
	Package          string           `json:"name"`
	Message          string           `json:"message"`
	FlaggedOutOfDate bool             `json:"flagged,omitempty"`
	Version          string           `json:"version,omitempty"`
	Upstream         upstream.Version `json:"upstream,omitempty"`
	Status           StatusType       `json:"status"`
}

func (status StatusType) color() string {
	switch status {
	case UpToDate:
		return "\x1b[32m"
	case FlaggedOutOfDate:
		return "\x1b[31m"
	case OutOfDate:
		return "\x1b[31m"
	default:
		return "\x1b[37m"
	}
}

// Print displays the status on the console
func (s *Status) Print() {
	ansiColor := s.Status.color()
	fmt.Fprintf(statusWriter, "%s%22s [%s][%s] %s \x1b[0m\n", ansiColor, "["+s.Status+"]", s.Package, s.Version, s.Message)
}

// PrintJSONTextSequence outputs the status as JSON Text Sequences (RFC 7464)
func (s *Status) PrintJSONTextSequence() {
	s.Type = "package"
	rfc7464.NewEncoder(statusWriter).Encode(s)
}
