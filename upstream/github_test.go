package upstream

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockGitHub() *gock.Response {
	return gock.New("https://api.github.com/").
		Get("/repos/gogits/gogs/releases/latest").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/atom+xml").
		BodyString(`
			{
				"url": "https://api.github.com/repos/gogits/gogs/releases/8625798",
				"id": 8625798,
				"tag_name": "v0.11.34",
				"target_commitish": "master",
				"name": "0.11.34",
				"draft": false,
				"prerelease": false,
				"created_at": "2017-11-22T19:46:14Z",
				"published_at": "2017-11-22T19:52:48Z"
			}			
			`)
}

func TestGogsGitHubUrl(t *testing.T) {
	defer gock.Off()
	mockGitHub()

	p := pkg.New("gogs", "0", "https://github.com/gogits/gogs")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "0.11.34" {
		t.Errorf("Expecting version 0.11.34, but got %v", version)
	}
}

func TestGogsGitHubSource(t *testing.T) {
	defer gock.Off()
	mockGitHub()

	p := pkg.New("gogs", "0", "https://gogs.io/", "https://github.com/gogits/gogs/archive/v0.11.34.tar.gz")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "0.11.34" {
		t.Errorf("Expecting version 0.11.34, but got %v", version)
	}
}
