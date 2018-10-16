package upstream

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/go-errors/errors"
)

type debianVersion struct {
	Version string `json:"version"`
}

type debianResponse struct {
	Versions []debianVersion `json:"versions"`
}

type debian string

func (d debian) releasesURL() string {
	// API documentation: https://sources.debian.org/doc/api/
	return fmt.Sprintf("https://sources.debian.org/api/src/%s/", url.PathEscape(string(d)))
}

func (d debian) latestVersion() (Version, error) {
	var res debianResponse
	if err := fetchJSON(d, &res); err != nil {
		return "", errors.WrapPrefix(err, "No debian release found for "+string(d), 0)
	}
	match := regexp.MustCompile("([^-]+)[-|~]").FindStringSubmatch(res.Versions[0].Version)
	if len(match) > 0 {
		return Version(match[1]), nil
	}
	return Version(res.Versions[0].Version), nil
}
