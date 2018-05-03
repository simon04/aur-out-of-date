package upstream

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockRubyGems() *gock.Response {
	return gock.New("https://rubygems.org").
		Get("/api/v1/versions/htmlbeautifier.json").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		BodyString(`[{
			"description": "A normaliser/beautifier for HTML that also understands embedded Ruby.",
			"number": "1.3.1",
			"summary": "HTML/ERB beautifier",
			"licenses": ["MIT"],
			"sha": "1af1b96b60969ad4721abe925620baa5aa68a6a77db71af8fe33e77e862b019c"
		}]`)
}

func TestRubyGemsSource1(t *testing.T) {
	defer gock.Off()
	mockRubyGems()

	p := pkg.New("ruby-htmlbeautifier", "0", "", "https://rubygems.org/downloads/htmlbeautifier-1.3.1.gem")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "1.3.1" {
		t.Errorf("Expecting version 1.3.1, but got %v", version)
	}
}
