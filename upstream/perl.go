package upstream

import (
	"fmt"
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
		return "", fmt.Errorf("No CPAN release found for %v: %w", p, err)
	}
	return Version(info.Version), nil
}
