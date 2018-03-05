package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-errors/errors"
)

type gitHubRelease struct {
	URL         string    `json:"url"`
	Name        string    `json:"name"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
}

func githubVersion(u string, re *regexp.Regexp) (Version, error) {
	match := re.FindSubmatch([]byte(u))
	if match == nil {
		return "", errors.Errorf("No GitHub release found for %s", u)
	}

	owner, repo := string(match[1]), string(match[2])
	// API documentation: https://developer.github.com/v3/repos/releases/
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo))
	if err != nil {
		return "", errors.WrapPrefix(err, "No GitHub release found for "+u, 0)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var release gitHubRelease
	err = dec.Decode(&release)
	if err != nil {
		return "", errors.WrapPrefix(err, "No GitHub release found for "+u, 0)
	} else if release.Name == "" {
		return "", errors.Errorf("No GitHub release found for %s", u)
	} else if release.Prerelease {
		return "", errors.Errorf("Ignoring GitHub pre-release %s for %s", release.Name, u)
	} else if release.Draft {
		return "", errors.Errorf("Ignoring GitHub release draft %s for %s", release.Name, u)
	}

	return Version(release.Name), nil
}
