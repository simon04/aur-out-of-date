package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/mikkeloscar/aur"
	"github.com/mikkeloscar/gopkgbuild"
	"github.com/mmcdole/gofeed"
)

func githubVersion(url string) (*pkgbuild.CompleteVersion, error) {
	match := regexp.MustCompile("github.com/([^/#.]+)/([^/#.]+)").FindSubmatch([]byte(url))
	if match == nil {
		match = regexp.MustCompile("([^/#.]+).github.io/([^/#.]+)").FindSubmatch([]byte(url))
	}
	if match == nil {
		return nil, fmt.Errorf("No GitHub release found for %s", url)
	}

	owner, repo := string(match[1]), string(match[2])
	feedURL := fmt.Sprintf("https://github.com/%s/%s/releases.atom", owner, repo)
	feed, err := gofeed.NewParser().ParseURL(feedURL)
	if err != nil || len(feed.Items) == 0 {
		return nil, fmt.Errorf("No GitHub release found for %s: %v", url, err)
	}

	link := regexp.MustCompile("/releases/tag/v?(.*)").FindSubmatch([]byte(feed.Items[0].Link))
	if link == nil {
		return nil, fmt.Errorf("No GitHub release found for %s: %v", url, err)
	}
	version := string(link[1])
	return pkgbuild.NewCompleteVersion(version)
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
	gitVersion, err := githubVersion(pkg.URL)
	if err == nil {
		return gitVersion, nil
	}
	pkgbuild, err := fetchPkgbuild(pkg)
	if err != nil {
		return nil, err
	}
	if len(pkgbuild.Source) > 0 {
		return githubVersion(pkgbuild.Source[0])
	}
	return nil, fmt.Errorf("No release found for %s: %v", pkg.Name, err)
}

func handlePackage(pkg *aur.Pkg) error {

	pkgVersion, err := pkgbuild.NewCompleteVersion(pkg.Version)
	if err != nil {
		return err
	}

	gitVersion, err := obtainVersion(pkg)
	if err != nil {
		return err
	}
	// workaround for https://github.com/mikkeloscar/gopkgbuild/pull/8
	version := fmt.Sprintf("%d:%s-%s", pkgVersion.Epoch, gitVersion.Version, pkgVersion.Pkgrel)

	fmt.Println(pkg.Name, pkg.URL, pkgVersion, gitVersion, pkgVersion.Older(version))
	return nil
}

func handlePackageForMaintainer(maintainer string) {
	packages, err := aur.SearchByMaintainer(maintainer)
	if err != nil {
		panic(err)
	}
	for _, pkg := range packages {
		err := handlePackage(&pkg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
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
