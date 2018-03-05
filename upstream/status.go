package upstream

import (
	"fmt"
)

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
	Package          string
	Message          string
	FlaggedOutOfDate bool
	Version          string
	Upstream         Version
	Status           string
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
	fmt.Printf("%s%22s [%s] %s \x1b[0m\n", ansiColor, "["+s.Status+"]", s.Package, s.Message)
}
