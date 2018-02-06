package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-errors/errors"
)

type pypiResponse struct {
	Info pypiInfo `json:"info"`
}

type pypiInfo struct {
	Version string `json:"version"`
}

func pythonVersion(url string, re *regexp.Regexp) (Version, error) {
	match := re.FindSubmatch([]byte(url))
	if match == nil {
		return "", errors.Errorf("No PyPI release found for %s", url)
	}
	resp, err := http.Get(fmt.Sprintf("https://pypi.python.org/pypi/%s/json", match[1]))
	if err != nil {
		return "", errors.WrapPrefix(err, "No PyPI release found for "+url, 0)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var pypi pypiResponse
	err = dec.Decode(&pypi)
	if err != nil || pypi.Info.Version == "" {
		return "", errors.WrapPrefix(err, "No PyPI release found for "+url, 0)
	}
	return Version(pypi.Info.Version), nil
}
