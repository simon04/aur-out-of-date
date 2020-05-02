package upstream

import (
	"fmt"
	"time"
)

type gitHubAPIReleases struct {
	gitHub
}

func (g gitHubAPIReleases) releasesURL() string {
	// API documentation: https://developer.github.com/v3/repos/releases/
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", g.owner, g.repository)
}

func (g gitHubAPIReleases) errorWrap(err error) error {
	return fmt.Errorf("Failed to obtain GitHub release for %s from %s: %w", g.String(), g.releasesURL(), err)
}

func (g gitHubAPIReleases) errorNotFound() error {
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

func (g gitHubAPIReleases) latestVersion() (Version, error) {
	var release gitHubRelease
	err := g.request(g.releasesURL(), &release)
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
