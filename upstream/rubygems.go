package upstream

import (
	"fmt"
	"time"
)

type rubygemsVersions []struct {
	Authors         string        `json:"authors"`
	BuiltAt         time.Time     `json:"built_at"`
	CreatedAt       time.Time     `json:"created_at"`
	Description     string        `json:"description"`
	DownloadsCount  int           `json:"downloads_count"`
	Number          string        `json:"number"`
	Summary         string        `json:"summary"`
	Platform        string        `json:"platform"`
	RubygemsVersion string        `json:"rubygems_version"`
	RubyVersion     string        `json:"ruby_version"`
	Prerelease      bool          `json:"prerelease"`
	Licenses        []string      `json:"licenses"`
	Requirements    []interface{} `json:"requirements"`
	Sha             string        `json:"sha"`
}

type rubygem string

func (g rubygem) releasesURL() string {
	return fmt.Sprintf("https://rubygems.org/api/v1/versions/%s.json", g)
}

func (g rubygem) latestVersion() (Version, error) {
	var versions rubygemsVersions
	if err := fetchJSON(g, &versions); err != nil || len(versions) == 0 {
		return "", fmt.Errorf("No RubyGems release found for %v: %w", g, err)
	}
	return Version(versions[0].Number), nil
}
