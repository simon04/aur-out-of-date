package pkg

import (
	"fmt"
	"strings"

	pkgbuild "github.com/mikkeloscar/gopkgbuild"
)

// NewLocalPkgs creates a Pkg slice from paths to .SRCINFO files.
func NewLocalPkgs(paths []string) ([]Pkg, error) {
	var r []Pkg
	for _, path := range paths {
		pkg, err := pkgbuild.ParseSRCINFO(path)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse %s: %w", path, err)
		}
		r = append(r, &localPkg{pkg, path})
	}
	return r, nil
}

type localPkg struct {
	pkg  *pkgbuild.PKGBUILD
	path string
}

func (p *localPkg) Name() string {
	return p.pkg.Pkgnames[0]
}

func (p *localPkg) Version() *pkgbuild.CompleteVersion {
	return &pkgbuild.CompleteVersion{
		Epoch:   uint8(p.pkg.Epoch),
		Version: p.pkg.Pkgver,
		Pkgrel:  p.pkg.Pkgrel,
	}
}

func (p *localPkg) LocalPKGBUILD() string {
	return strings.Replace(p.path, ".SRCINFO", "PKGBUILD", 1)
}

func (p *localPkg) URL() string {
	return p.pkg.URL
}

func (p *localPkg) Sources() ([]string, error) {
	return p.pkg.Source, nil
}

func (p *localPkg) OutOfDate() bool {
	return false
}
