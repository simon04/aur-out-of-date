package upstream

import (
	"net/http"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockGitHubAtom() *gock.Response {
	return gock.New("https://github.com/").
		Get("/leaflet/leaflet/releases.atom").
		Reply(http.StatusOK).
		BodyString(`
			<?xml version="1.0" encoding="UTF-8"?>
			<feed>
				<title>Release notes from Leaflet</title>
				<updated>2019-11-17T21:04:59Z</updated>
				<entry>
					<id>tag:github.com,2008:Repository/931135/v1.6.0</id>
					<updated>2019-11-17T21:12:44Z</updated>
					<link rel="alternate" type="text/html" href="https://github.com/Leaflet/Leaflet/releases/tag/v1.6.0"/>
					<title>v1.6.0</title>
					<author>
						<name>cherniavskii</name>
					</author>
				</entry>
			</feed>`)
}

func TestLeafletGitHubAtom(t *testing.T) {
	defer gock.Off()
	mockGitHubAtom()

	os.Setenv("GITHUB_ATOM", "1")
	p := pkg.New("leaflet", "0", "https://github.com/leaflet/leaflet")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version.String() != "1.6.0" {
		t.Errorf("Expecting version 1.6.0, but got %v", version)
	}
	os.Unsetenv("GITHUB_ATOM")
}
