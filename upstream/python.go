package upstream

import (
	"fmt"
)

type pypiResponse struct {
	Info struct {
		Version string `json:"version"`
	} `json:"info"`
}

type pypi string

func (p pypi) releasesURL() string {
	return fmt.Sprintf("https://pypi.org/pypi/%s/json", p)
}

func (p pypi) latestVersion() (Version, error) {
	var response pypiResponse
	if err := fetchJSON(p, &response); err != nil || response.Info.Version == "" {
		return "", fmt.Errorf("No PyPI release found for %v: %w", p, err)
	}
	return Version(response.Info.Version), nil
}
