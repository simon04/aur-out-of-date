package upstream

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockGitLab() *gock.Response {
	return gock.New("https://gitlab.com").
		Get("/api/v4/projects/gitlab-org/gitlab-ce/repository/tags").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		// Since it always returns a taglist, added 2 to the return
		BodyString(`
			[
				{
				"name": "v11.0.0-rc13",
				"message": "Version v11.0.0-rc13",
				"target": "162a86504a9c56dc142aab8bb6beeeaca29a6b86",
				"commit": {
					"id": "0a9f160f39d47948a0a011907a7d0e989df3cd69",
					"short_id": "0a9f160f",
					"title": "Update VERSION to 11.0.0-rc13",
					"created_at": "2018-06-18T11:55:51.000Z",
					"parent_ids": [
						"b020eb1e15bc8323affe372cf9f713e3b3ae27ad"
					],
					"message": "Update VERSION to 11.0.0-rc13\n",
					"author_name": "GitLab Release Tools Bot",
					"author_email": "robert+release-tools@gitlab.com",
					"authored_date": "2018-06-18T11:55:51.000Z",
					"committer_name": "GitLab Release Tools Bot",
					"committer_email": "robert+release-tools@gitlab.com",
					"committed_date": "2018-06-18T11:55:51.000Z"
				},
				"release": null
			},
			{
				"name": "v11.0.0-rc12",
				"message": "Version v11.0.0-rc12",
				"target": "6dc273b4f9881e0c90524bd478dad3bba126ab86",
				"commit": {
					"id": "44f330568c3fee5c28224d0eb4d6e4e0fe46be53",
					"short_id": "44f33056",
					"title": "Update VERSION to 11.0.0-rc12",
					"created_at": "2018-06-14T12:56:56.000Z",
					"parent_ids": [
						"99dfb12a912c5d38074a5c352d5ff76781e6f2cc"
					],
					"message": "Update VERSION to 11.0.0-rc12\n",
					"author_name": "GitLab Release Tools Bot",
					"author_email": "robert+release-tools@gitlab.com",
					"authored_date": "2018-06-14T12:56:56.000Z",
					"committer_name": "GitLab Release Tools Bot",
					"committer_email": "robert+release-tools@gitlab.com",
					"committed_date": "2018-06-14T12:56:56.000Z"
				},
				"release": null
			}
			]
		`)
}

func TestGitlabceGitLabSource(t *testing.T) {
	defer gock.Off()
	mockGitLab()

	p := pkg.New("gitlab-ce", "0", "https://gitlab.com/gitlab-org/gitlab-ce", "https://gitlab.com/gitlab-org/gitlab-ce/-/archive/v11.0.0-rc13/gitlab-ce-v11.0.0-rc13.tar.gz")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version.String() != "11.0.0-rc13" {
		t.Errorf("Expecting version 11.0.0-rc13, but got %v", version)
	}
}
