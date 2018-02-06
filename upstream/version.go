package upstream

import (
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	pkgbuild "github.com/mikkeloscar/gopkgbuild"
	"github.com/simon04/aur-out-of-date/pkg"
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
		return nil, errors.Errorf("No release found for %s", url)
	}
}

// VersionForPkg determines the upstream version for the given package
func VersionForPkg(pkg pkg.Pkg) (*pkgbuild.CompleteVersion, error) {
	version, err := forURL(pkg.URL())
	if err == nil {
		return version, nil
	}
	sources, err := pkg.Sources()
	if err != nil {
		return nil, errors.WrapPrefix(err, "Failed to obtain sources for "+pkg.Name(), 0)
	}
	if len(sources) > 0 {
		return forURL(sources[0])
	}
	return nil, errors.WrapPrefix(err, "No release found for "+pkg.Name(), 0)
}
