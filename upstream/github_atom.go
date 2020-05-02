package upstream

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

type gitHubAPIAtom struct {
	gitHub
}

type atomFeed struct {
	Items []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Link atomLink `xml:"link"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
}

func (g gitHubAPIAtom) atomURL() string {
	return fmt.Sprintf("https://github.com/%s/%s/releases.atom", g.owner, g.repository)
}

func (g gitHubAPIAtom) errorWrap(err error) error {
	return fmt.Errorf("Failed to obtain GitHub release for %s from %s: %w", g.String(), g.atomURL(), err)
}

func (g gitHubAPIAtom) errorNotFound() error {
	return fmt.Errorf("No GitHub release found for %s on %s", g, g.atomURL())
}

func (g gitHubAPIAtom) latestVersion() (Version, error) {
	resp, err := http.Get(g.atomURL())
	if err != nil {
		return "", g.errorWrap(err)
	}
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)
	var feed atomFeed
	err = dec.Decode(&feed)
	if err != nil {
		return "", g.errorWrap(err)
	} else if len(feed.Items) == 0 {
		return "", g.errorNotFound()
	}

	href, err := url.PathUnescape(feed.Items[0].Link.Href)
	if err != nil {
		return "", g.errorWrap(err)
	}
	link := regexp.MustCompile("/releases/tag/v?(.*)").FindSubmatch([]byte(href))
	if link == nil {
		return "", g.errorNotFound()
	}
	version := string(link[1])
	return Version(version), nil
}
