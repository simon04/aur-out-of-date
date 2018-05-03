package upstream

import (
	"fmt"

	"github.com/go-errors/errors"
)

type pypiResponse struct {
	Info struct {
		Version string `json:"version"`
	} `json:"info"`
}

type pypi string

func (p pypi) releasesURL() string {
	return fmt.Sprintf("https://pypi.python.org/pypi/%s/json", p)
}

func (p pypi) latestVersion() (Version, error) {
	var response pypiResponse
	if err := fetchJSON(p, &response); err != nil || response.Info.Version == "" {
		return "", errors.WrapPrefix(err, "No PyPI release found for "+string(p), 0)
	}
	return Version(response.Info.Version), nil
}
