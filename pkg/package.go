package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mikkeloscar/aur"
	pkgbuild "github.com/mikkeloscar/gopkgbuild"
)

// Pkg is an interface representing an Arch Linux package.
type Pkg interface {
	Name() string
	Version() *pkgbuild.CompleteVersion
	URL() string
	Sources() ([]string, error)
	OutOfDate() bool
}

// NewRemotePkgs creates a Pkg slice from information returned from AUR RPC.
func NewRemotePkgs(pkg []aur.Pkg) []Pkg {
	var r []Pkg
	for i := range pkg {
		r = append(r, &remotePkg{&pkg[i]})
	}
	return r
}

type remotePkg struct {
	pkg *aur.Pkg
}

func (p *remotePkg) Name() string {
	return p.pkg.Name
}

func (p *remotePkg) Version() *pkgbuild.CompleteVersion {
	version, _ := pkgbuild.NewCompleteVersion(p.pkg.Version)
	return version
}

func (p *remotePkg) URL() string {
	return p.pkg.URL
}

func (p *remotePkg) Sources() ([]string, error) {
	resp, err := http.Get("https://aur.archlinux.org/cgit/aur.git/plain/.SRCINFO?h=" + p.pkg.Name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	pkgbuildBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	pkg, err := pkgbuild.ParseSRCINFOContent(pkgbuildBytes)
	if err != nil {
		return nil, err
	}
	return pkg.Source, nil
}

func (p *remotePkg) OutOfDate() bool {
	return p.pkg.OutOfDate > 0
}

// NewLocalPkgs creates a Pkg slice from paths to .SRCINFO files.
func NewLocalPkgs(paths []string) ([]Pkg, error) {
	var r []Pkg
	for _, path := range paths {
		pkg, err := pkgbuild.ParseSRCINFO(path)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse %s: %v", path, err)
		}
		r = append(r, &localPkg{pkg})
	}
	return r, nil
}

type localPkg struct {
	pkg *pkgbuild.PKGBUILD
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

func (p *localPkg) URL() string {
	return p.pkg.URL
}

func (p *localPkg) Sources() ([]string, error) {
	return p.pkg.Source, nil
}

func (p *localPkg) OutOfDate() bool {
	return false
}
