package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/mikkeloscar/aur"
	"github.com/simon04/aur-out-of-date/pkg"
	"github.com/simon04/aur-out-of-date/upstream"
	xdgbasedir "github.com/zchee/go-xdgbasedir"
)

var statistics struct {
	UpToDate         int
	FlaggedOutOfDate int
	OutOfDate        int
	Unknown          int
}

var commandline struct {
	user            string
	remote          bool
	local           bool
	includeVcsPkgs  bool
	printStatistics bool
	flagOnAur       bool
}

func handlePackage(pkg pkg.Pkg) upstream.Status {

	pkgVersion := pkg.Version()
	status := upstream.Status{
		Package:          pkg.Name(),
		FlaggedOutOfDate: pkg.OutOfDate(),
		Version:          pkgVersion.String(),
	}

	upstreamVersion, err := upstream.VersionForPkg(pkg)
	if err != nil {
		status.Status = upstream.Unknown
		status.Message = err.Error()
		statistics.Unknown++
		return status
	}
	status.Upstream = upstreamVersion

	if pkg.OutOfDate() {
		status.Status = upstream.FlaggedOutOfDate
		status.Message = fmt.Sprintf("Package %s has been flagged out-of-date and should be updated from %v-%v to %v", pkg.Name(), pkgVersion.Version, pkgVersion.Pkgrel, upstreamVersion)
		statistics.FlaggedOutOfDate++
	} else if pkgVersion.Older(string(upstreamVersion)) {
		status.Status = upstream.OutOfDate
		status.Message = fmt.Sprintf("Package %s should be updated from %v-%v to %v", pkg.Name(), pkgVersion.Version, pkgVersion.Pkgrel, upstreamVersion)
		statistics.OutOfDate++
		flagOnAur(pkg, upstreamVersion)
	} else {
		status.Status = upstream.UpToDate
		status.Message = fmt.Sprintf("Package %s %v-%v matches upstream version %v", pkg.Name(), pkgVersion.Version, pkgVersion.Pkgrel, upstreamVersion)
		statistics.UpToDate++
	}
	return status
}

func flagOnAur(pkg pkg.Pkg, upstreamVersion upstream.Version) {
	if !commandline.flagOnAur {
		return
	}
	fmt.Printf("Should the package %s be flagged out-of-date? [y/N] ", pkg.Name())
	var response string
	chars, err := fmt.Scanln(&response)
	if err != nil || chars == 0 {
		return
	}
	if response != "y" && response != "Y" {
		return
	}
	fmt.Printf("Flagging package %s out-of-date ...\n", pkg.Name())
	comment := fmt.Sprintf("Version %s is out. #simon04/aur-out-of-date", upstreamVersion)
	cmd := exec.Command("ssh", "aur@aur.archlinux.org", "flag", pkg.Name(), "\""+comment+"\"")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to flag out-of-date (running \"%v\"): %v\n%s\n", strings.Join(cmd.Args, "\" \""), err, output)
	} else {
		fmt.Printf("%s", output)
	}
}

func handlePackages(vcsPackages bool, packages []pkg.Pkg, err error) {
	if err != nil {
		panic(err)
	}
	sort.Slice(packages, func(i, j int) bool { return strings.Compare(packages[i].Name(), packages[j].Name()) == -1 })
	for _, pkg := range packages {
		isVcsPackage := strings.HasSuffix(pkg.Name(), "-git") || strings.HasSuffix(pkg.Name(), "-hg") || strings.HasSuffix(pkg.Name(), "-svn")
		if vcsPackages == isVcsPackage {
			status := handlePackage(pkg)
			status.Print()
		}
	}
}

func printStatistics() {
	fmt.Println()
	fmt.Println("[STATISTICS]")
	fmt.Printf("\x1b[32m%12s: %4d \x1b[0m\n", "UP-TO-DATE", statistics.UpToDate)
	fmt.Printf("\x1b[31m%12s: %4d \x1b[0m\n", "FLAGGED", statistics.FlaggedOutOfDate)
	fmt.Printf("\x1b[31m%12s: %4d \x1b[0m\n", "OUT-OF-DATE", statistics.OutOfDate)
	fmt.Printf("\x1b[37m%12s: %4d \x1b[0m\n", "UNKNOWN", statistics.Unknown)
}

func main() {
	flag.StringVar(&commandline.user, "user", "", "AUR username")
	flag.BoolVar(&commandline.remote, "pkg", false, "AUR package name(s)")
	flag.BoolVar(&commandline.local, "local", false, "Local .SRCINFO files")
	flag.BoolVar(&commandline.includeVcsPkgs, "devel", false, "Check -git/-svn/-hg packages")
	flag.BoolVar(&commandline.printStatistics, "statistics", false, "Print summary statistics")
	flag.BoolVar(&commandline.flagOnAur, "flag", false, "Flag out-of-date on AUR")
	flag.Parse()

	// cache HTTP requests (RFC 7234)
	cacheDir := path.Join(xdgbasedir.CacheHome(), "aur-out-of-date")
	http.DefaultClient = httpcache.NewTransport(diskcache.New(cacheDir)).Client()

	if commandline.user != "" {
		packages, err := aur.SearchByMaintainer(commandline.user)
		handlePackages(commandline.includeVcsPkgs, pkg.NewRemotePkgs(packages), err)
	} else if commandline.remote {
		pkgs := flag.Args()
		for len(pkgs) > 0 {
			limit := 100
			if len(pkgs) < limit {
				limit = len(pkgs)
			}
			packages, err := aur.Info(pkgs[:limit])
			handlePackages(false, pkg.NewRemotePkgs(packages), err)
			handlePackages(true, pkg.NewRemotePkgs(packages), err)
			pkgs = pkgs[limit:]
		}
	} else if commandline.local {
		packages, err := pkg.NewLocalPkgs(flag.Args())
		handlePackages(false, packages, err)
		handlePackages(true, packages, err)
	} else {
		fmt.Fprintln(os.Stderr, "Either -user or -pkg or -local is required!")
		flag.Usage()
		os.Exit(1)
	}
	if commandline.printStatistics {
		printStatistics()
	}
}
