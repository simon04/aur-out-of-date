package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-errors/errors"
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

func (g rubygem) errorWrap(err error) error {
	return errors.WrapPrefix(err, "No RubyGems release found for "+string(g)+" on "+g.releasesURL(), 0)
}

func rubygemsVersion(u string, re *regexp.Regexp) (Version, error) {
	match := re.FindSubmatch([]byte(u))
	if match == nil {
		return "", errors.Errorf("No RubyGems release found for %s", u)
	}
	gem := rubygem(string(match[1]))
	resp, err := http.Get(gem.releasesURL())
	if err != nil {
		return "", gem.errorWrap(err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var versions rubygemsVersions
	err = dec.Decode(&versions)
	if err != nil || len(versions) == 0 {
		return "", gem.errorWrap(err)
	}
	return Version(versions[0].Number), nil
}
