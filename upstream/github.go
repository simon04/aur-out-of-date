package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type gitHub struct {
	owner      string
	repository string
}

func (g gitHub) String() string {
	return g.owner + "/" + g.repository
}

func parseGitHub(url string) *gitHub {
	match := regexp.MustCompile("github.com/([^/#.]+)/([^/#]+)").FindStringSubmatch(url)
	if len(match) > 0 {
		return &gitHub{match[1], match[2]}
	}
	match = regexp.MustCompile("([^/#.]+).github.io/([^/#]+)").FindStringSubmatch(url)
	if len(match) > 0 {
		return &gitHub{match[1], match[2]}
	}
	return nil
}

func (g gitHub) request(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	// Obtain GitHub token for higher request limits, see https://developer.github.com/v3/#rate-limiting
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode == http.StatusForbidden {
		var message gitHubMessage
		err = dec.Decode(&message)
		if err == nil && message.Message != "" {
			err = fmt.Errorf("%s", message.Message)
		}
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("No GitHub project found for %s on %s", g, url)
	}
	return dec.Decode(target)
}
