package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type gitHub struct {
	owner      string
	repository string
}

func (g gitHub) String() string {
	return g.owner + "/" + g.repository
}

func (g gitHub) releasesURL() string {
	// API documentation: https://developer.github.com/v3/repos/releases/
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", g.owner, g.repository)
}

func (g gitHub) errorWrap(err error) error {
	return fmt.Errorf("Failed to obtain GitHub release for %s from %s: %w", g.String(), g.releasesURL(), err)
}

func (g gitHub) errorNotFound() error {
	return fmt.Errorf("No GitHub release found for %s on %s", g, g.releasesURL())
}

type gitHubRelease struct {
	URL         string    `json:"url"`
	Name        string    `json:"name"`
	TagName     string    `json:"tag_name"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
}

type gitHubMessage struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func (g gitHub) latestVersion() (Version, error) {
	req, err := http.NewRequest("GET", g.releasesURL(), nil)

	// Obtain GitHub token for higher request limits, see https://developer.github.com/v3/#rate-limiting
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	if err != nil {
		return "", g.errorWrap(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", g.errorWrap(err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode == http.StatusForbidden {
		var message gitHubMessage
		err = dec.Decode(&message)
		if err == nil && message.Message != "" {
			err = fmt.Errorf("%s", message.Message)
		}
		return "", g.errorWrap(err)
	} else if resp.StatusCode == http.StatusNotFound {
		return "", g.errorNotFound()
	}

	var release gitHubRelease
	err = dec.Decode(&release)
	if err != nil {
		return "", g.errorWrap(err)
	} else if release.Prerelease {
		return "", fmt.Errorf("Ignoring GitHub pre-release %s for %s", release.Name, g.String())
	} else if release.Draft {
		return "", fmt.Errorf("Ignoring GitHub release draft %s for %s", release.Name, g.String())
	} else if release.TagName != "" {
		return Version(release.TagName), nil
	} else if release.Name != "" {
		return Version(release.Name), nil
	}
	return "", g.errorNotFound()
}
