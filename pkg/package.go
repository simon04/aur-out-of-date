package pkg

import (
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
