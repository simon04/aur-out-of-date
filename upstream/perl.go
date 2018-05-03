package upstream

import (
	"fmt"

	"github.com/go-errors/errors"
)

type cpanRelease struct {
	Version string `json:"version"`
}

type cpan string

func (p cpan) releasesURL() string {
	// API documentation: https://github.com/metacpan/metacpan-api/blob/master/docs/API-docs.md
	return fmt.Sprintf("https://fastapi.metacpan.org/v1/release/%s", p)
}

func (p cpan) latestVersion() (Version, error) {
	var info cpanRelease
	if err := fetchJSON(p, &info); err != nil || info.Version == "" {
		return "", errors.WrapPrefix(err, "No CPAN release found for "+string(p), 0)
	}
	return Version(info.Version), nil
}
