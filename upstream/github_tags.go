package upstream

import (
	"fmt"
)

type gitHubAPITags struct {
	gitHub
}

func (g gitHubAPITags) tagsURL() string {
	// API documentation: https://developer.github.com/v3/repos/#list-tags
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", g.owner, g.repository)
}

func (g gitHubAPITags) errorWrap(err error) error {
	return fmt.Errorf("Failed to obtain GitHub tag for %s from %s: %w", g.String(), g.tagsURL(), err)
}

func (g gitHubAPITags) errorNotFound() error {
	return fmt.Errorf("No GitHub tag found for %s on %s", g, g.tagsURL())
}

type gitHubTag struct {
	Name string `json:"name"`
}

func (g gitHubAPITags) latestVersion() (Version, error) {
	var taglist []gitHubTag
	err := g.request(g.tagsURL(), &taglist)
	if err != nil {
		return "", g.errorWrap(err)
	} else if len(taglist) > 0 {
		if taglist[0].Name != "" {
			return Version(taglist[0].Name), nil
		}
	}
	return "", g.errorNotFound()
}
