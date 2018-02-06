package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-errors/errors"
)

type npmDistTags struct {
	Latest string `json:"latest"`
}

func npmVersion(url string, re *regexp.Regexp) (Version, error) {
	match := re.FindSubmatch([]byte(url))
	if match == nil {
		return "", errors.Errorf("No npm release found for %s", url)
	}
	resp, err := http.Get(fmt.Sprintf("https://registry.npmjs.org/-/package/%s/dist-tags", match[1]))
	if err != nil {
		return "", errors.WrapPrefix(err, "No npm release found for "+url, 0)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var distTags npmDistTags
	err = dec.Decode(&distTags)
	if err != nil || distTags.Latest == "" {
		return "", errors.WrapPrefix(err, "No npm release found for "+url, 0)
	}
	return Version(distTags.Latest), nil
}
