package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/go-errors/errors"
)

type npmDistTags struct {
	Latest string `json:"latest"`
}

func npmVersion(u string, re *regexp.Regexp) (Version, error) {
	match := re.FindSubmatch([]byte(u))
	if match == nil {
		return "", errors.Errorf("No npm release found for %s", u)
	}
	pkg := url.PathEscape(string(match[1]))
	resp, err := http.Get(fmt.Sprintf("https://registry.npmjs.org/-/package/%s/dist-tags", pkg))
	if err != nil {
		return "", errors.WrapPrefix(err, "No npm release found for "+u, 0)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var distTags npmDistTags
	err = dec.Decode(&distTags)
	if err != nil || distTags.Latest == "" {
		return "", errors.WrapPrefix(err, "No npm release found for "+u, 0)
	}
	return Version(distTags.Latest), nil
}
