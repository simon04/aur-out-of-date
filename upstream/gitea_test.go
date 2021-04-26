package upstream

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockCodeberg() *gock.Response {
	return gock.New("https://codeberg.org").
		Get("/api/v1/repos/Anoxinon_e.V./xmppc/releases").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		BodyString(`
			[
				{
					"id": 409858,
					"tag_name": "0.1.0",
					"target_commitish": "master",
					"name": "xmppc Version 0.1.0",
					"body": "* Support Account configuration with and without pwd (ask user to enter pwd)\r\n* Request roster list\r\n* Show MAM - XEP-0313: Message Archive Management\r\n* Show Bookmarks - XEP-0048: Bookmarks\r\n* Request Service Discovery items and info - XEP-0030: Service Discovery\r\n* Display OMEMO device list and fingerprints (URI format)\r\n* Send chat message (unencrypted) \r\n* Send chat message signcrypted - XEP-0373: OpenPGP for XMPP\r\n* Send chat mesage pgp - XEP-0027\r\n* Monitor XMPP stanza \r\n",
					"url": "https://codeberg.org/api/v1/repos/Anoxinon_e.V./xmppc/releases/409858",
					"html_url": "https://codeberg.org/Anoxinon_e.V./xmppc/releases/tag/0.1.0",
					"tarball_url": "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.1.0.tar.gz",
					"zipball_url": "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.1.0.zip",
					"draft": true,
					"prerelease": false,
					"created_at": "2020-05-16T18:19:05+02:00",
					"published_at": "2020-05-16T18:19:05+02:00",
					"author": {
						"id": 558,
						"login": "DebXWoody",
						"full_name": "",
						"email": "debxwoody@noreply.codeberg.org",
						"avatar_url": "https://codeberg.org/user/avatar/DebXWoody/-1",
						"language": "en-US",
						"is_admin": false,
						"last_login": "2021-04-24T10:32:42+02:00",
						"created": "2019-04-01T08:11:02+02:00",
						"username": "DebXWoody"
					},
					"assets": []
				},
				{
					"id": 409512,
					"tag_name": "0.0.6",
					"target_commitish": "master",
					"name": "xmppc Version 0.0.6",
					"body": "BugFix  - OpenPGP\r\n\r\n* https://codeberg.org/Anoxinon_e.V./xmppc/issues/7\r\n* https://codeberg.org/Anoxinon_e.V./xmppc/issues/8",
					"url": "https://codeberg.org/api/v1/repos/Anoxinon_e.V./xmppc/releases/409512",
					"html_url": "https://codeberg.org/Anoxinon_e.V./xmppc/releases/tag/0.0.6",
					"tarball_url": "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.0.6.tar.gz",
					"zipball_url": "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.0.6.zip",
					"draft": false,
					"prerelease": true,
					"created_at": "2020-05-02T16:05:33+02:00",
					"published_at": "2020-05-02T16:05:33+02:00",
					"author": {
						"id": 558,
						"login": "DebXWoody",
						"full_name": "",
						"email": "debxwoody@noreply.codeberg.org",
						"avatar_url": "https://codeberg.org/user/avatar/DebXWoody/-1",
						"language": "en-US",
						"is_admin": false,
						"last_login": "2021-04-24T10:32:42+02:00",
						"created": "2019-04-01T08:11:02+02:00",
						"username": "DebXWoody"
					},
					"assets": []
				},
				{
					"id": 409241,
					"tag_name": "0.0.5",
					"target_commitish": "master",
					"name": "xmppc Version 0.0.5",
					"body": "xmppc Version 0.0.5\r\n\r\n * XEP-0313: Message Archive Management\r\n * XEP-0048: Bookmarks\r\n * XEP-0030: Service Discovery",
					"url": "https://codeberg.org/api/v1/repos/Anoxinon_e.V./xmppc/releases/409241",
					"html_url": "https://codeberg.org/Anoxinon_e.V./xmppc/releases/tag/0.0.5",
					"tarball_url": "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.0.5.tar.gz",
					"zipball_url": "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.0.5.zip",
					"draft": false,
					"prerelease": false,
					"created_at": "2020-04-25T06:25:32+02:00",
					"published_at": "2020-04-25T06:25:32+02:00",
					"author": {
						"id": 558,
						"login": "DebXWoody",
						"full_name": "",
						"email": "debxwoody@noreply.codeberg.org",
						"avatar_url": "https://codeberg.org/user/avatar/DebXWoody/-1",
						"language": "en-US",
						"is_admin": false,
						"last_login": "2021-04-24T10:32:42+02:00",
						"created": "2019-04-01T08:11:02+02:00",
						"username": "DebXWoody"
					},
					"assets": []
				},
				{
					"id": 409065,
					"tag_name": "0.0.4",
					"target_commitish": "master",
					"name": "xmppc Version 0.0.4",
					"body": "* Config file for accounts\r\n* Changed output format of omemo list to URL format\r\n* Bugfixes for OpenPGP / PGP Key lookup\r\n",
					"url": "https://codeberg.org/api/v1/repos/Anoxinon_e.V./xmppc/releases/409065",
					"html_url": "https://codeberg.org/Anoxinon_e.V./xmppc/releases/tag/0.0.4",
					"tarball_url": "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.0.4.tar.gz",
					"zipball_url": "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.0.4.zip",
					"draft": false,
					"prerelease": false,
					"created_at": "2020-04-21T18:23:40+02:00",
					"published_at": "2020-04-21T18:23:40+02:00",
					"author": {
						"id": 558,
						"login": "DebXWoody",
						"full_name": "",
						"email": "debxwoody@noreply.codeberg.org",
						"avatar_url": "https://codeberg.org/user/avatar/DebXWoody/-1",
						"language": "en-US",
						"is_admin": false,
						"last_login": "2021-04-24T10:32:42+02:00",
						"created": "2019-04-01T08:11:02+02:00",
						"username": "DebXWoody"
					},
					"assets": []
				}
			]
		`)
}

func TestGiteaSource(t *testing.T) {
	defer gock.Off()
	mockCodeberg()

	p := pkg.New("xmppc", "0", "https://codeberg.org/Anoxinon_e.V./xmppc", "https://codeberg.org/Anoxinon_e.V./xmppc/archive/0.1.0.tar.gz")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version.String() != "0.0.5" {
		t.Errorf("Expecting version 0.0.5, but got %v", version)
	}
}
