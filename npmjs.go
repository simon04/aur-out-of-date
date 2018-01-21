package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	pkgbuild "github.com/mikkeloscar/gopkgbuild"
)

type npmDistTags struct {
	Latest string `json:"latest"`
}

func npmVersion(url string, re *regexp.Regexp) (*pkgbuild.CompleteVersion, error) {
	match := re.FindSubmatch([]byte(url))
	if match == nil {
		return nil, fmt.Errorf("No npm release found for %s", url)
	}
	resp, err := http.Get(fmt.Sprintf("https://registry.npmjs.org/-/package/%s/dist-tags", match[1]))
	if err != nil {
		return nil, fmt.Errorf("No npm release found for %s: %v", url, err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var distTags npmDistTags
	err = dec.Decode(&distTags)
	if err != nil || distTags.Latest == "" {
		return nil, fmt.Errorf("No npm release found for %s: %v", url, err)
	}
	return pkgbuild.NewCompleteVersion(distTags.Latest)
}
