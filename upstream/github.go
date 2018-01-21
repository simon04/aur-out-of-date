package upstream

import (
	"fmt"
	"regexp"

	pkgbuild "github.com/mikkeloscar/gopkgbuild"
	"github.com/mmcdole/gofeed"
)

func githubVersion(url string, re *regexp.Regexp) (*pkgbuild.CompleteVersion, error) {
	match := re.FindSubmatch([]byte(url))
	if match == nil {
		return nil, fmt.Errorf("No GitHub release found for %s", url)
	}

	owner, repo := string(match[1]), string(match[2])
	feedURL := fmt.Sprintf("https://github.com/%s/%s/releases.atom", owner, repo)
	feed, err := gofeed.NewParser().ParseURL(feedURL)
	if err != nil || len(feed.Items) == 0 {
		return nil, fmt.Errorf("No GitHub release found for %s: %v", url, err)
	}

	link := regexp.MustCompile("/releases/tag/v?(.*)").FindSubmatch([]byte(feed.Items[0].Link))
	if link == nil {
		return nil, fmt.Errorf("No GitHub release found for %s: %v", url, err)
	}
	version := string(link[1])
	return pkgbuild.NewCompleteVersion(version)
}
