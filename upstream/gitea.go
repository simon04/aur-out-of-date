package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type gitea struct {
	Host       string
	Owner      string
	Repostiory string
}

type giteaRelease struct {
	Name       string `json:"name"`
	TagName    string `json:"tag_name"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

func parseGitea(host string, url string) *gitea {
	regex := fmt.Sprintf("%s/([^/#]+)/([^/#]+)", host)
	match := regexp.MustCompile(regex).FindStringSubmatch(url)
	if len(match) > 0 {
		return &gitea{host, match[1], match[2]}
	}
	return nil
}

func (g *gitea) latestVersion() (Version, error) {
	var releases []giteaRelease
	releaseURL := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/releases", g.Host, g.Owner, g.Repostiory)
	resp, err := http.Get(releaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&releases)
	if err != nil {
		return "", err
	}
	if releases[0].TagName != "" {
		return Version(releases[0].TagName), nil
	}
	return "", fmt.Errorf("No Gitea tag found for %s", g)
}
