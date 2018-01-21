package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/mikkeloscar/aur"
	"github.com/mikkeloscar/gopkgbuild"
)

func obtainVersionForURL(url string) (*pkgbuild.CompleteVersion, error) {
	switch {
	case strings.Contains(url, "github.com"):
		return githubVersion(url, regexp.MustCompile("github.com/([^/#.]+)/([^/#.]+)"))
	case strings.Contains(url, ("github.io")):
		return githubVersion(url, regexp.MustCompile("([^/#.]+).github.io/([^/#.]+)"))
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

func obtainVersion(pkg *aur.Pkg) (*pkgbuild.CompleteVersion, error) {
	version, err := obtainVersionForURL(pkg.URL)
	if err == nil {
		return version, nil
	}
	pkgbuild, err := fetchPkgbuild(pkg)
	if err != nil {
		return nil, err
	}
	if len(pkgbuild.Source) > 0 {
		return obtainVersionForURL(pkgbuild.Source[0])
	}
	return nil, fmt.Errorf("No release found for %s: %v", pkg.Name, err)
}

func handlePackage(pkg *aur.Pkg) {

	pkgVersion, err := pkgbuild.NewCompleteVersion(pkg.Version)
	if err != nil {
		fmt.Printf("\x1b[37m[UNKNOWN]     [%s] %v \x1b[0m\n", pkg.Name, err)
		return
	}

	gitVersion, err := obtainVersion(pkg)
	if err != nil {
		fmt.Printf("\x1b[37m[UNKNOWN]     [%s] %v \x1b[0m\n", pkg.Name, err)
		return
	}
	// workaround for https://github.com/mikkeloscar/gopkgbuild/pull/8
	version := fmt.Sprintf("%d:%s-%s", pkgVersion.Epoch, gitVersion.Version, pkgVersion.Pkgrel)

	if pkgVersion.Older(version) {
		fmt.Printf("\x1b[31m[OUT-OF-DATE] [%s] Package %s should be updated from %v-%v to %v \x1b[0m\n", pkg.Name, pkg.Name, pkgVersion.Version, pkgVersion.Pkgrel, gitVersion.Version)
	} else {
		fmt.Printf("\x1b[32m[UP-TO-DATE]  [%s] Package %s %v-%v matches upstream version %v \x1b[0m\n", pkg.Name, pkg.Name, pkgVersion.Version, pkgVersion.Pkgrel, gitVersion.Version)
	}
}

type byName []aur.Pkg

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) == -1 }

func handlePackageForMaintainer(maintainer string) {
	packages, err := aur.SearchByMaintainer(maintainer)
	if err != nil {
		panic(err)
	}
	sort.Sort(byName(packages))
	for _, pkg := range packages {
		handlePackage(&pkg)
	}
}

func main() {
	user := flag.String("user", "", "AUR username")
	flag.Parse()
	if *user == "" {
		fmt.Fprintln(os.Stderr, "-user is required")
		flag.Usage()
		os.Exit(1)
	}
	handlePackageForMaintainer(*user)
}
