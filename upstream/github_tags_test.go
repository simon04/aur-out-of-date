package upstream

import (
	"net/http"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockGitHubTags() *gock.Response {
	return gock.New("https://api.github.com/").
		Get("/repos/gogits/gogs/tags").
		Reply(http.StatusOK).
		BodyString(`
			[
				{
					"name": "v0.11.34",
					"zipball_url": "https://api.github.com/repos/gogs/gogs/zipball/v0.11.34",
					"tarball_url": "https://api.github.com/repos/gogs/gogs/tarball/v0.11.34",
					"commit": {
						"sha": "6f2347fc71f17b5703a9b1f383a2d3451f88b741",
						"url": "https://api.github.com/repos/gogs/gogs/commits/6f2347fc71f17b5703a9b1f383a2d3451f88b741"
					},
					"node_id": "MDM6UmVmMTY3NTI2MjA6djAuMTEuMzQ="
				},
				{
					"name": "v0.11.33",
					"zipball_url": "https://api.github.com/repos/gogs/gogs/zipball/v0.11.33",
					"tarball_url": "https://api.github.com/repos/gogs/gogs/tarball/v0.11.33",
					"commit": {
						"sha": "b752fe680811119954ccef051e6f3b3e2a04c2e8",
						"url": "https://api.github.com/repos/gogs/gogs/commits/b752fe680811119954ccef051e6f3b3e2a04c2e8"
					},
					"node_id": "MDM6UmVmMTY3NTI2MjA6djAuMTEuMzM="
				}
			]`)
}

func TestGogsGitHubTags(t *testing.T) {
	defer gock.Off()
	mockGitHubTags()

	os.Setenv("GITHUB_TAGS", "1")
	p := pkg.New("gogs", "0", "https://github.com/gogits/gogs")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version.String() != "0.11.34" {
		t.Errorf("Expecting version 0.11.34, but got %v", version)
	}
	os.Unsetenv("GITHUB_TAGS")
}
