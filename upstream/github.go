package upstream

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"

	pkgbuild "github.com/mikkeloscar/gopkgbuild"
)

type atomFeed struct {
	Items []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Link atomLink `xml:"link"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
}

func githubVersion(url string, re *regexp.Regexp) (*pkgbuild.CompleteVersion, error) {
	match := re.FindSubmatch([]byte(url))
	if match == nil {
		return nil, fmt.Errorf("No GitHub release found for %s", url)
	}

	owner, repo := string(match[1]), string(match[2])
	resp, err := http.Get(fmt.Sprintf("https://github.com/%s/%s/releases.atom", owner, repo))
	if err != nil {
		return nil, fmt.Errorf("No GitHub release found for %s: %v", url, err)
	}
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)
	var feed atomFeed
	err = dec.Decode(&feed)
	if err != nil || len(feed.Items) == 0 {
		return nil, fmt.Errorf("No GitHub release found for %s: %v", url, err)
	}

	link := regexp.MustCompile("/releases/tag/v?(.*)").FindSubmatch([]byte(feed.Items[0].Link.Href))
	if link == nil {
		return nil, fmt.Errorf("No GitHub release found for %s: %v", url, err)
	}
	version := string(link[1])
	return pkgbuild.NewCompleteVersion(version)
}
