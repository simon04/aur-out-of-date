package pkg

import (
	pkgbuild "github.com/mikkeloscar/gopkgbuild"
)

// Pkg is an interface representing an Arch Linux package.
type Pkg interface {
	Name() string
	Version() *pkgbuild.CompleteVersion
	LocalPKGBUILD() string
	URL() string
	Sources() ([]string, error)
	OutOfDate() bool
}

// New creates a Pkg from the given parameters. Mainly used for testing.
func New(name, version, url string, sources ...string) Pkg {
	pkg := pkgbuild.PKGBUILD{
		Pkgbase:  name,
		Pkgnames: []string{name},
		Pkgver:   pkgbuild.Version(version),
		Pkgrel:   "1",
		URL:      url,
		Source:   sources,
	}
	return &localPkg{pkg: &pkg, path: ""}
}
