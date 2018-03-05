package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

type gitHubMessage struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func githubVersion(u string, re *regexp.Regexp) (Version, error) {
	match := re.FindSubmatch([]byte(u))
	if match == nil {
		return "", errors.Errorf("No GitHub release found for %s", u)
	}

	owner, repo := string(match[1]), string(match[2])
	// API documentation: https://developer.github.com/v3/repos/releases/
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, err := http.NewRequest("GET", api, nil)

	// Obtain GitHub token for higher request limits, see https://developer.github.com/v3/#rate-limiting
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	if err != nil {
		return "", errors.WrapPrefix(err, "No GitHub release found for "+u, 0)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.WrapPrefix(err, "No GitHub release found for "+u, 0)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode == http.StatusForbidden {
		var message gitHubMessage
		err = dec.Decode(&message)
		if err == nil && message.Message != "" {
			err = errors.Wrap(message.Message, 0)
		}
		return "", errors.WrapPrefix(err, "Failed to obtain GitHub release for "+u, 0)
	}

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
