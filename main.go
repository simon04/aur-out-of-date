package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mikkeloscar/aur"
	"github.com/simon04/aur-out-of-date/pkg"
	"github.com/simon04/aur-out-of-date/upstream"
)

func handlePackage(pkg pkg.Pkg) {

	pkgVersion := pkg.Version()

	upstreamVersion, err := upstream.VersionForPkg(pkg)
	if err != nil {
		fmt.Printf("\x1b[37m[UNKNOWN]     [%s] %v \x1b[0m\n", pkg.Name(), err)
		return
	}
	upstreamVersion.Epoch = 0
	upstreamVersion.Pkgrel = ""

	if pkg.OutOfDate() {
		fmt.Printf("\x1b[31m[OUT-OF-DATE] [%s] Package %s has been flagged out-of-date and should be updated from %v-%v to %v \x1b[0m\n", pkg.Name(), pkg.Name(), pkgVersion, pkgVersion.Pkgrel, upstreamVersion.Version)
	} else if pkgVersion.Older(upstreamVersion.String()) {
		fmt.Printf("\x1b[31m[OUT-OF-DATE] [%s] Package %s should be updated from %v-%v to %v \x1b[0m\n", pkg.Name(), pkg.Name(), pkgVersion, pkgVersion.Pkgrel, upstreamVersion.Version)
	} else {
		fmt.Printf("\x1b[32m[UP-TO-DATE]  [%s] Package %s %v-%v matches upstream version %v \x1b[0m\n", pkg.Name(), pkg.Name(), pkgVersion, pkgVersion.Pkgrel, upstreamVersion.Version)
	}
}

// byName is used for sorting packages by their name
type byName []pkg.Pkg

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return strings.Compare(a[i].Name(), a[j].Name()) == -1 }

func handlePackages(vcsPackages bool, packages []pkg.Pkg, err error) {
	if err != nil {
		panic(err)
	}
	sort.Sort(byName(packages))
	for _, pkg := range packages {
		isVcsPackage := strings.HasSuffix(pkg.Name(), "-git") || strings.HasSuffix(pkg.Name(), "-hg") || strings.HasSuffix(pkg.Name(), "-svn")
		if vcsPackages == isVcsPackage {
			handlePackage(pkg)
		}
	}
}

// stringSlice is used for parsing multi-value string flags
type stringSlice []string

func (i *stringSlice) String() string {
	return strings.Join(*i, " ")
}
func (i *stringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var packages stringSlice
	flag.Var(&packages, "pkg", "AUR package name(s)")
	user := flag.String("user", "", "AUR username")
	vcsPackages := flag.Bool("devel", false, "Check -git/-svn/-hg packages")
	flag.Parse()
	if *user != "" {
		packages, err := aur.SearchByMaintainer(*user)
		handlePackages(*vcsPackages, pkg.NewRemotePkgs(packages), err)
	} else if len(packages) > 0 {
		packages, err := aur.Info(packages)
		handlePackages(false, pkg.NewRemotePkgs(packages), err)
		handlePackages(true, pkg.NewRemotePkgs(packages), err)
	} else {
		fmt.Fprintln(os.Stderr, "Either -user or -pkg is required!")
		flag.Usage()
		os.Exit(1)
	}
}
