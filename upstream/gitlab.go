package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/go-errors/errors"
)

// Self-hosted GitLab instances use different domain names
type gitLab struct {
	domain     string
	owner      string
	repository string
}

func (g gitLab) String() string {
	return g.domain + "/" + g.owner + "/" + g.repository
}

func (g gitLab) encoded() string {
	// GitLab requires the owner + repository to be url-encoded together
	return url.PathEscape(g.owner) + "%2F" + url.PathEscape(g.repository)
}

func (g gitLab) releasesURL() string {
	// API documentation: https://docs.gitlab.com/ee/api/tags.html#list-project-repository-tags
	// Note that the second %s must be url-encoded (see gitLab.encoded())
	return fmt.Sprintf("https://%s/api/v4/projects/%s/repository/tags", g.domain, g.encoded())
}

func (g gitLab) errorWrap(err error) error {
	return errors.WrapPrefix(err, "Failed to obtain GitLab tag for "+g.String()+" from "+g.releasesURL(), 0)
}

func (g gitLab) errorNotFound() error {
	return errors.Errorf("No GitLab release found for %s on %s", g, g.releasesURL())
}

// Describes the individual tags in the returned taglist from the json call
type gitLabTag struct {
	Name string `json:"name"`
}

type gitLabMessage struct {
	Message string `json:"message"`
}

func (g gitLab) latestVersion() (Version, error) {
	req, err := http.NewRequest("GET", g.releasesURL(), nil)

	// Obtain GitLab token for higher request limits, see https://docs.gitlab.com/ee/api/#oauth2-tokens
	token := os.Getenv("GITLAB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
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
		var message gitLabMessage
		err = dec.Decode(&message)
		if err == nil && message.Message != "" {
			err = errors.Wrap(message.Message, 0)
		}
		return "", g.errorWrap(err)
	} else if resp.StatusCode == http.StatusNotFound {
		return "", g.errorNotFound()
	}

	// Can't get single tag, has to be an array
	// NOTE: If GitLab ever adds a "get newest tag" API call then change this
	var taglist []gitLabTag
	err = dec.Decode(&taglist)
	if err != nil {
		return "", g.errorWrap(err)
	} else if len(taglist) > 0 {
		// [0] will always be the newest, as its sorted by default
		if taglist[0].Name != "" {
			return Version(taglist[0].Name), nil
		}
	}
	return "", g.errorNotFound()
}
