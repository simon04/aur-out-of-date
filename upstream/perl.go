package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-errors/errors"
)

type cpanInfo struct {
	Version string `json:"version"`
}

func perlVersion(url string, re *regexp.Regexp) (Version, error) {
	match := re.FindSubmatch([]byte(url))
	if match == nil {
		return "", errors.Errorf("No CPAN release found for %s", url)
	}
	// API documentation: https://github.com/metacpan/metacpan-api/blob/master/docs/API-docs.md
	resp, err := http.Get(fmt.Sprintf("https://fastapi.metacpan.org/v1/release/%s", match[1]))
	if err != nil {
		return "", errors.WrapPrefix(err, "No CPAN release found for "+url, 0)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var cpan cpanInfo
	err = dec.Decode(&cpan)
	if err != nil || cpan.Version == "" {
		return "", errors.WrapPrefix(err, "No CPAN release found for "+url, 0)
	}
	return Version(cpan.Version), nil
}
