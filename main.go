package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/mikkeloscar/aur"
	"github.com/simon04/aur-out-of-date/action"
	"github.com/simon04/aur-out-of-date/config"
	"github.com/simon04/aur-out-of-date/pkg"
	"github.com/simon04/aur-out-of-date/status"
	"github.com/simon04/aur-out-of-date/upstream"
)

var conf *config.Config
var statistics status.Statistics

var commandline struct {
	user            string
	config          string
	remote          bool
	local           bool
	includeVcsPkgs  bool
	printJSON       bool
	printStatistics bool
	flagOnAur       bool
	updatePKGBUILD  bool
}

func handlePackage(pkg pkg.Pkg) status.Status {

	pkgVersion := pkg.Version()
	s := status.Status{
		Package:          pkg.Name(),
		FlaggedOutOfDate: pkg.OutOfDate(),
		Version:          pkgVersion.String(),
	}

	upstreamVersion, err := upstream.VersionForPkg(pkg)
	if err != nil {
		s.Status = status.Unknown
		s.Message = err.Error()
		statistics.Unknown++
		return s
	}

	s.Ignored = conf.IsIgnored(pkg.Name(), upstreamVersion)
	s.Compare(upstreamVersion)
	statistics.Update(s.Status)
	return s
}

func handlePackages(vcsPackages bool, packages []pkg.Pkg, err error) {
	if err != nil {
		panic(err)
	}
	sort.Slice(packages, func(i, j int) bool { return strings.Compare(packages[i].Name(), packages[j].Name()) == -1 })
	for _, pkg := range packages {
		if vcsPackages == pkg.IsVcs() {
			s := handlePackage(pkg)
			if commandline.printJSON {
				s.PrintJSONTextSequence()
			} else {
				s.Print()
			}
			if s.Status == status.OutOfDate && commandline.flagOnAur {
				action.FlagOnAur(pkg, s.Upstream)
			}
			if s.Status == status.OutOfDate && commandline.updatePKGBUILD {
				action.UpdatePKGBUILD(pkg, s.Upstream)
			}
		}
	}
}

func main() {
	configDir, _ := os.UserConfigDir()
	defaultConfigFile := path.Join(configDir, "aur-out-of-date", "config.json")
	flag.StringVar(&commandline.user, "user", "", "AUR username")
	flag.StringVar(&commandline.config, "config", defaultConfigFile, "Config file")
	flag.BoolVar(&commandline.remote, "pkg", false, "AUR package name(s)")
	flag.BoolVar(&commandline.local, "local", false, "Local .SRCINFO files")
	flag.BoolVar(&commandline.includeVcsPkgs, "devel", false, "Check -git/-svn/-hg packages")
	flag.BoolVar(&commandline.printStatistics, "statistics", false, "Print summary statistics")
	flag.BoolVar(&commandline.flagOnAur, "flag", false, "Flag out-of-date on AUR")
	flag.BoolVar(&commandline.updatePKGBUILD, "update", false, "Update pkgver/pkgrel in local PKGBUILD files")
	flag.BoolVar(&commandline.printJSON, "json", false, "Generate JSON Text Sequences (RFC 7464)")
	flag.Parse()

	// cache HTTP requests (RFC 7234)
	cacheDir, _ := os.UserCacheDir()
	cacheDir = path.Join(cacheDir, "aur-out-of-date")
	http.DefaultClient = httpcache.NewTransport(diskcache.New(cacheDir)).Client()

	if c, err := config.FromFile(commandline.config); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read config:", err)
		os.Exit(1)
	} else {
		conf = c
	}

	if commandline.user != "" {
		packages, err := aur.SearchBy(commandline.user, aur.Maintainer)
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
		packages, err := pkg.NewLocalPkgs(flag.Args(), commandline.includeVcsPkgs)
		handlePackages(false, packages, err)
		handlePackages(true, packages, err)
	} else {
		fmt.Fprintln(os.Stderr, "Either -user or -pkg or -local is required!")
		flag.Usage()
		os.Exit(1)
	}
	if commandline.printStatistics && commandline.printJSON {
		statistics.PrintJSONTextSequence()
	} else if commandline.printStatistics {
		statistics.Print()
	}
	if statistics.OutOfDate > 0 {
		os.Exit(4)
	}
}
