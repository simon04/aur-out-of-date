package upstream

import (
	"fmt"
	"net/url"

	"github.com/go-errors/errors"
)

type npmDistTags struct {
	Latest string `json:"latest"`
}

type npm string

func (n npm) releasesURL() string {
	// API documentation: https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md
	return fmt.Sprintf("https://registry.npmjs.org/-/package/%s/dist-tags", url.PathEscape(string(n)))
}

func (n npm) latestVersion() (Version, error) {
	var distTags npmDistTags
	if err := fetchJSON(n, &distTags); err != nil || distTags.Latest == "" {
		return "", errors.WrapPrefix(err, "No npm release found for "+string(n), 0)
	}
	return Version(distTags.Latest), nil
}
