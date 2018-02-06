package upstream

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockPerl() *gock.Response {
	return gock.New("https://fastapi.metacpan.org/").
		Get("/v1/release/Perl-Critic").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
			"main_module" : "Perl::Critic",
			"date" : "2017-07-21T04:28:55",
			"status" : "latest",
			"author" : "PETDANCE",
			"distribution" : "Perl-Critic",
			"version" : "1.130",
			"changes_file" : "Changes",
			"name" : "Perl-Critic-1.130"
		}`)
}

func TestPerlSource1(t *testing.T) {
	defer gock.Off()
	mockPerl()

	p := pkg.New("perl-critic", "0", "", "http://search.cpan.org/CPAN/authors/id/T/TH/THALJEF/Perl-Critic-1.126.tar.gz")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "1.130" {
		t.Errorf("Expecting version 1.130, but got %v", version)
	}
}

func TestPerlSource2(t *testing.T) {
	defer gock.Off()
	mockPerl()

	p := pkg.New("perl-critic", "0", "", "https://cpan.metacpan.org/authors/id/P/PE/PETDANCE/Perl-Critic-1.126.tar.gz")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "1.130" {
		t.Errorf("Expecting version 1.130, but got %v", version)
	}
}
