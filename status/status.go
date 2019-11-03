package status

import (
	"fmt"
	"io"
	"os"

	pkgbuild "github.com/mikkeloscar/gopkgbuild"
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
	Ignored          bool             `json:"ignored,omitempty"`
	Version          string           `json:"version,omitempty"`
	Upstream         upstream.Version `json:"upstream,omitempty"`
	Status           StatusType       `json:"status"`
}

// Compare to upstream version and set message and status accordingly
func (s *Status) Compare(upstreamVersion upstream.Version) {
	pkgVersion, err := pkgbuild.NewCompleteVersion(s.Version)
	if err != nil {
		s.Status = Unknown
		s.Message = fmt.Sprintf("Failed to parse pkg version: %v", err)
		return
	}

	upstreamCompleteVersion, err := pkgbuild.NewCompleteVersion(upstreamVersion.String())
	if err != nil {
		s.Status = Unknown
		s.Message = fmt.Sprintf("Failed to parse upstream version: %v", err)
		return
	}
	s.Upstream = upstreamVersion

	newer := upstreamCompleteVersion.Newer(pkgVersion)
	if s.FlaggedOutOfDate {
		s.Status = FlaggedOutOfDate
		s.Message = fmt.Sprintf("has been flagged out-of-date and should be updated to %v", upstreamVersion)
	} else if newer && s.Ignored {
		s.Status = Unknown
		s.Message = fmt.Sprintf("ignoring package upgrade to %v", upstreamVersion)
	} else if newer {
		s.Status = OutOfDate
		s.Message = fmt.Sprintf("should be updated to %v", upstreamVersion)
	} else if upstreamCompleteVersion.Equal(pkgVersion) {
		s.Status = UpToDate
		s.Message = fmt.Sprintf("matches upstream version %v", upstreamVersion)
	} else {
		s.Status = Unknown
		s.Message = fmt.Sprintf("upstream version is %v", upstreamVersion)
	}
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
