package pkg

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/mikkeloscar/aur"
)

func mockInfo() *gock.Response {
	g := gock.New("https://aur.archlinux.org/")
	g.URLStruct.RawPath = "/rpc.php?arg[]=python2-mwclient&type=info&v=5"
	return g.
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		BodyString(`
		{
			"version": 5,
			"type": "multiinfo",
			"resultcount": 1,
			"results": [
				{
					"ID": 533001,
					"Name": "python2-mwclient",
					"PackageBaseID": 109452,
					"PackageBase": "python-mwclient",
					"Version": "0.9.1-1",
					"Description": "A Python framework to interface with the MediaWiki API",
					"URL": "https://github.com/mwclient/mwclient",
					"NumVotes": 6,
					"Popularity": 0.455733,
					"OutOfDate": null,
					"Maintainer": "simon04",
					"FirstSubmitted": 1459241039,
					"LastModified": 1533710052,
					"URLPath": "/cgit/aur.git/snapshot/python-mwclient.tar.gz",
					"Depends": ["python2", "python2-requests-oauthlib"],
					"License": ["MIT"],
					"Keywords": []
				}
			]
		}`)
}

func mockSRCINFO() *gock.Response {
	g := gock.New("https://aur.archlinux.org/")
	g.URLStruct.RawPath = "/cgit/aur.git/plain/.SRCINFO?h=python-mwclient"
	return g.
		Reply(http.StatusOK).
		BodyString("pkgbase = python-mwclient\n\tpkgdesc = A Python framework to interface with the MediaWiki API\n\tpkgver = 0.9.1\n\tpkgrel = 1\n\turl = https://github.com/mwclient/mwclient\n\tarch = any\n\tlicense = MIT\n\tsource = python-mwclient-0.9.1.tar.gz::https://github.com/mwclient/mwclient/archive/v0.9.1.tar.gz\n\tsha512sums = e2c8d720bc583f2cf0de2bdfaab3dfce9f23ed541c34fa8d164d35e9c134e39110d1f9b791daf4a4cf79f18084052ec644ba96980d2037a06b2d0a7851af5ed4\n\npkgname = python-mwclient\n\tdepends = python\n\tdepends = python-requests-oauthlib\n\npkgname = python2-mwclient\n\tdepends = python2\n\tdepends = python2-requests-oauthlib\n")
}

func TestSplitPkg(t *testing.T) {
	defer gock.Off()
	mockInfo()
	mockSRCINFO()

	info, err := aur.Info([]string{"python2-mwclient"})
	if err != nil {
		t.Error(err)
	}
	pkgs := NewRemotePkgs(info)
	if len(pkgs) != 1 {
		t.Errorf("Found %d pkgs!", len(pkgs))
	}
	sources, err := pkgs[0].Sources()
	if err != nil {
		t.Error(err)
	}
	if len(sources) == 0 || sources[0] != "python-mwclient-0.9.1.tar.gz::https://github.com/mwclient/mwclient/archive/v0.9.1.tar.gz" {
		t.Errorf("Unexpected sources %v", sources)
	}
}
