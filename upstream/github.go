package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-errors/errors"
)

type gitHub struct {
	owner      string
	repository string
}

func (g gitHub) String() string {
	return g.owner + "/" + g.repository
}

func (g gitHub) tagsURL() string {
	// API documentation: https://developer.github.com/v3/repos/#list-tags
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", g.owner, g.repository)
}

func (g gitHub) errorWrap(err error) error {
	return errors.WrapPrefix(err, "Failed to obtain GitHub tag for "+g.String()+" from "+g.tagsURL(), 0)
}

func (g gitHub) errorNotFound() error {
	return errors.Errorf("No GitHub tag found for %s on %s", g, g.tagsURL())
}

type gitHubTag struct {
	Name        string    `json:"name"`
}

type gitHubMessage struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func (g gitHub) latestVersion() (Version, error) {
	req, err := http.NewRequest("GET", g.tagsURL(), nil)

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
			err = errors.Wrap(message.Message, 0)
		}
		return "", g.errorWrap(err)
	} else if resp.StatusCode == http.StatusNotFound {
		return "", g.errorNotFound()
	}

	var taglist []gitHubTag
	err = dec.Decode(&taglist)
	if err != nil {
		return "", g.errorWrap(err)
	} else if len(taglist) > 0 {
		if taglist[0].Name != "" {
			return Version(taglist[0].Name), nil
		}
	}
	return "", g.errorNotFound()
}
