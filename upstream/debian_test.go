package upstream

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockDebian() *gock.Response {
	return gock.New("https://sources.debian.org").
		Get("/api/src/babeltrace/").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		BodyString(`
			{
				"package": "babeltrace",
				"path": "babeltrace",
				"pathl": [
					[
						"babeltrace",
						"/src/babeltrace/"
					]
				],
				"suite": "",
				"type": "package",
				"versions": [
					{
						"area": "main",
						"suites": [
							"buster",
							"sid"
						],
						"version": "1.5.6-1"
					},
					{
						"area": "main",
						"suites": [
							"stretch"
						],
						"version": "1.5.1-1"
					},
					{
						"area": "main",
						"suites": [
							"jessie",
							"jessie-kfreebsd"
						],
						"version": "1.2.3-2"
					}
				]
			}`)
}

func TestDebianSource1(t *testing.T) {
	defer gock.Off()
	mockDebian()

	p := pkg.New("babeltrace", "0", "", "http://debian.backend.mirrors.debian.org/debian/pool/main/b/babeltrace/python3-babeltrace_1.5.6-1_hurd-i386.deb")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "1.5.6" {
		t.Errorf("Expecting version 1.5.6, but got %v", version)
	}
}
