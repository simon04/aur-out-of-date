package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/mikkeloscar/aur"
	"github.com/simon04/aur-out-of-date/pkg"
	"github.com/simon04/aur-out-of-date/upstream"
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

func handlePackage(pkg pkg.Pkg) {

	pkgVersion := pkg.Version()

	upstreamVersion, err := upstream.VersionForPkg(pkg)
	if err != nil {
		fmt.Printf("\x1b[37m[UNKNOWN]     [%s] %v \x1b[0m\n", pkg.Name(), err)
		statistics.Unknown++
		return
	}

	if pkg.OutOfDate() {
		fmt.Printf("\x1b[31m[OUT-OF-DATE] [%s] Package %s has been flagged out-of-date and should be updated from %v-%v to %v \x1b[0m\n", pkg.Name(), pkg.Name(), pkgVersion.Version, pkgVersion.Pkgrel, upstreamVersion)
		statistics.FlaggedOutOfDate++
	} else if pkgVersion.Older(string(upstreamVersion)) {
		fmt.Printf("\x1b[31m[OUT-OF-DATE] [%s] Package %s should be updated from %v-%v to %v \x1b[0m\n", pkg.Name(), pkg.Name(), pkgVersion.Version, pkgVersion.Pkgrel, upstreamVersion)
		statistics.OutOfDate++
		flagOnAur(pkg, upstreamVersion)
	} else {
		fmt.Printf("\x1b[32m[UP-TO-DATE]  [%s] Package %s %v-%v matches upstream version %v \x1b[0m\n", pkg.Name(), pkg.Name(), pkgVersion.Version, pkgVersion.Pkgrel, upstreamVersion)
		statistics.UpToDate++
	}
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
