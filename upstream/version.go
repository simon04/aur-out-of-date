package upstream

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/simon04/aur-out-of-date/pkg"
)

// Version represents the upstream version of a software project
type Version string

// String returns a sanitized version string
func (v Version) String() string {
	s := string(v)
	s = strings.TrimLeft(s, "v")
	return s
}

func forURL(url string) (Version, error) {
	switch {
	case strings.Contains(url, "github.com"):
		match := regexp.MustCompile("github.com/([^/#.]+)/([^/#]+)").FindStringSubmatch(url)
		if len(match) > 0 {
			return gitHub{match[1], match[2]}.latestVersion()
		}
	case strings.Contains(url, "github.io"):
		match := regexp.MustCompile("([^/#.]+).github.io/([^/#]+)").FindStringSubmatch(url)
		if len(match) > 0 {
			return gitHub{match[1], match[2]}.latestVersion()
		}
	case strings.Contains(url, "registry.npmjs.org"):
		match := regexp.MustCompile("registry.npmjs.org/((@[^/#.]+/)?[^/#.]+)").FindStringSubmatch(url)
		if len(match) > 0 {
			return npm(match[1]).latestVersion()
		}
	case strings.Contains(url, "npmjs.com/package"):
		fallthrough
	case strings.Contains(url, "npmjs.org/package"):
		match := regexp.MustCompile("/package/((@[^/#.]+/)?[^/#.]+)").FindStringSubmatch(url)
		if len(match) > 0 {
			return npm(match[1]).latestVersion()
		}
	case strings.Contains(url, "pypi.python.org"):
		fallthrough
	case strings.Contains(url, "files.pythonhosted.org"):
		fallthrough
	case strings.Contains(url, "pypi.io"):
		fallthrough
	case strings.Contains(url, "pypi.org"):
		match := regexp.MustCompile("/packages/source/[^/#.]+/([^/#.]+)/").FindStringSubmatch(url)
		if len(match) > 0 {
			return pypi(match[1]).latestVersion()
		}
		match = regexp.MustCompile("/([^/#.]+)-[0-9.]+(post.)?.tar.gz$").FindStringSubmatch(url)
		if len(match) > 0 {
			return pypi(match[1]).latestVersion()
		}
	case strings.Contains(url, "search.cpan.org"):
		fallthrough
	case strings.Contains(url, "search.mcpan.org"):
		fallthrough
	case strings.Contains(url, "cpan.metacpan.org"):
		match := regexp.MustCompile("/([^/#.]+?)-v?([0-9.-]+)\\.(tgz|tar.gz)$").FindStringSubmatch(url)
		if len(match) > 0 {
			return cpan(match[1]).latestVersion()
		}
	case strings.Contains(url, "rubygems.org"):
		fallthrough
	case strings.Contains(url, "gems.rubyforge.org"):
		match := regexp.MustCompile("/([^/#]+?)-[^-]+\\.gem$").FindStringSubmatch(url)
		if len(match) > 0 {
			return rubygem(match[1]).latestVersion()
		}
	case strings.Contains(url, "gitlab"):
		// Example: https://gitlab.com/gitlab-org/gitlab-ce/-/archive/v11.0.0-rc7/gitlab-ce-v11.0.0-rc7.tar.gz
		match := regexp.MustCompile("https?://([^/]+)/([^/]+)/([^/]+)(\\.git|/.*)?$").FindStringSubmatch(url)
		if len(match) > 0 {
			return gitLab{match[1], match[2], match[3]}.latestVersion()
		}
	}
	return "", errors.Errorf("No release found for %s", url)
}

// VersionForPkg determines the upstream version for the given package
func VersionForPkg(pkg pkg.Pkg) (Version, error) {
	version, err := forURL(pkg.URL())
	if err == nil {
		return version, nil
	}
	sources, err := pkg.Sources()
	if err != nil {
		return "", errors.WrapPrefix(err, "Failed to obtain sources for "+pkg.Name(), 0)
	}
	if len(sources) > 0 {
		return forURL(sources[0])
	}
	return "", errors.WrapPrefix(err, "No release found for "+pkg.Name(), 0)
}

type releasesAPI interface {
	releasesURL() string
	latestVersion() (Version, error)
}

func fetchJSON(a releasesAPI, target interface{}) error {
	url := a.releasesURL()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(target)
}
