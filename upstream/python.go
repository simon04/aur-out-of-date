package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	pkgbuild "github.com/mikkeloscar/gopkgbuild"
)

type pypiResponse struct {
	Info pypiInfo `json:"info"`
}

type pypiInfo struct {
	Version string `json:"version"`
}

func pythonVersion(url string, re *regexp.Regexp) (*pkgbuild.CompleteVersion, error) {
	match := re.FindSubmatch([]byte(url))
	if match == nil {
		return nil, fmt.Errorf("No PyPI release found for %s", url)
	}
	resp, err := http.Get(fmt.Sprintf("https://pypi.python.org/pypi/%s/json", match[1]))
	if err != nil {
		return nil, fmt.Errorf("No PyPI release found for %s: %v", url, err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var pypi pypiResponse
	err = dec.Decode(&pypi)
	if err != nil || pypi.Info.Version == "" {
		return nil, fmt.Errorf("No PyPI release found for %s: %v", url, err)
	}
	return pkgbuild.NewCompleteVersion(pypi.Info.Version)
}
