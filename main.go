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
	pkgbuild "github.com/mikkeloscar/gopkgbuild"
	"github.com/simon04/aur-out-of-date/config"
	"github.com/simon04/aur-out-of-date/pkg"
	"github.com/simon04/aur-out-of-date/status"
	"github.com/simon04/aur-out-of-date/upstream"
	xdgbasedir "github.com/zchee/go-xdgbasedir"
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
	upstreamCompleteVersion, err := pkgbuild.NewCompleteVersion(upstreamVersion.String())
	if err != nil {
		s.Status = status.Unknown
		s.Message = fmt.Sprintf("Failed to parse upstream version: %v", err)
		statistics.Unknown++
		return s
	}
	s.Upstream = upstreamVersion

	newer := upstreamCompleteVersion.Newer(pkgVersion)
	ignored := conf.IsIgnored(pkg.Name(), s.Upstream)
	if pkg.OutOfDate() {
		s.Status = status.FlaggedOutOfDate
		s.Message = fmt.Sprintf("has been flagged out-of-date and should be updated to %v", upstreamVersion)
		statistics.FlaggedOutOfDate++
	} else if newer && ignored {
		s.Status = status.Unknown
		s.Message = fmt.Sprintf("ignoring package upgrade to %v", upstreamVersion)
		statistics.OutOfDate++
	} else if newer {
		s.Status = status.OutOfDate
		s.Message = fmt.Sprintf("should be updated to %v", upstreamVersion)
		statistics.OutOfDate++
	} else {
		s.Status = status.UpToDate
		s.Message = fmt.Sprintf("matches upstream version %v", upstreamVersion)
		statistics.UpToDate++
	}
	return s
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
			s := handlePackage(pkg)
			if commandline.printJSON {
				s.PrintJSONTextSequence()
			} else {
				s.Print()
			}
			if s.Status == status.OutOfDate {
				flagOnAur(pkg, s.Upstream)
			}
		}
	}
}

func main() {
	defaultConfigFile := path.Join(xdgbasedir.ConfigHome(), "aur-out-of-date", "config.json")
	flag.StringVar(&commandline.user, "user", "", "AUR username")
	flag.StringVar(&commandline.config, "config", defaultConfigFile, "Config file")
	flag.BoolVar(&commandline.remote, "pkg", false, "AUR package name(s)")
	flag.BoolVar(&commandline.local, "local", false, "Local .SRCINFO files")
	flag.BoolVar(&commandline.includeVcsPkgs, "devel", false, "Check -git/-svn/-hg packages")
	flag.BoolVar(&commandline.printStatistics, "statistics", false, "Print summary statistics")
	flag.BoolVar(&commandline.flagOnAur, "flag", false, "Flag out-of-date on AUR")
	flag.BoolVar(&commandline.printJSON, "json", false, "Generate JSON Text Sequences (RFC 7464)")
	flag.Parse()

	// cache HTTP requests (RFC 7234)
	cacheDir := path.Join(xdgbasedir.CacheHome(), "aur-out-of-date")
	http.DefaultClient = httpcache.NewTransport(diskcache.New(cacheDir)).Client()

	if c, err := config.FromFile(commandline.config); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read config:", err)
		os.Exit(1)
	} else {
		conf = c
	}

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
	if commandline.printStatistics && commandline.printJSON {
		statistics.PrintJSONTextSequence()
	} else if commandline.printStatistics {
		statistics.Print()
	}
	if statistics.OutOfDate > 0 {
		os.Exit(4)
	}
}
