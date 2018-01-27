package upstream

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/mikkeloscar/aur"
	pkgbuild "github.com/mikkeloscar/gopkgbuild"
)

func forURL(url string) (*pkgbuild.CompleteVersion, error) {
	switch {
	case strings.Contains(url, "github.com"):
		return githubVersion(url, regexp.MustCompile("github.com/([^/#.]+)/([^/#.]+)"))
	case strings.Contains(url, "github.io"):
		return githubVersion(url, regexp.MustCompile("([^/#.]+).github.io/([^/#.]+)"))
	case strings.Contains(url, "registry.npmjs.org"):
		return npmVersion(url, regexp.MustCompile("registry.npmjs.org/([^/#.]+)/"))
	case strings.Contains(url, "pypi.python.org"):
		return pythonVersion(url, regexp.MustCompile("pypi.python.org/packages/source/[^/#.]+/([^/#.]+)/"))
	case strings.Contains(url, "files.pythonhosted.org"):
		return pythonVersion(url, regexp.MustCompile("files.pythonhosted.org/packages/source/[^/#.]+/([^/#.]+)/"))
	case strings.Contains(url, "search.cpan.org"):
		fallthrough
	case strings.Contains(url, "search.mcpan.org"):
		return perlVersion(url, regexp.MustCompile("/([^/#.]+?)-v?([0-9.-]+)\\.(tgz|tar.gz)$"))
	default:
		return nil, fmt.Errorf("No release found for %s", url)
	}
}

func fetchPkgbuild(pkg *aur.Pkg) (*pkgbuild.PKGBUILD, error) {
	resp, err := http.Get("https://aur.archlinux.org/cgit/aur.git/plain/.SRCINFO?h=" + pkg.Name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	pkgbuildBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return pkgbuild.ParseSRCINFOContent(pkgbuildBytes)
}

// VersionForPkg determines the upstream version for the given package
func VersionForPkg(pkg *aur.Pkg) (*pkgbuild.CompleteVersion, error) {
	version, err := forURL(pkg.URL)
	if err == nil {
		return version, nil
	}
	pkgbuild, err := fetchPkgbuild(pkg)
	if err != nil {
		return nil, err
	}
	if len(pkgbuild.Source) > 0 {
		return forURL(pkgbuild.Source[0])
	}
	return nil, fmt.Errorf("No release found for %s: %v", pkg.Name, err)
}
