package upstream

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockGitHub() *gock.Response {
	return gock.New("https://github.com/").
		Get("/gogits/gogs/releases.atom").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/atom+xml").
		BodyString(`
			<?xml version="1.0" encoding="UTF-8"?>
			<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
				<id>tag:github.com,2008:https://github.com/gogits/gogs/releases</id>
				<link type="text/html" rel="alternate" href="https://github.com/gogits/gogs/releases"/>
				<link type="application/atom+xml" rel="self" href="https://github.com/gogits/gogs/releases.atom"/>
				<title>Release notes from gogs</title>
				<updated>2017-11-22T20:46:14+01:00</updated>
				<entry>
					<id>tag:github.com,2008:Repository/16752620/v0.11.34</id>
					<updated>2017-11-22T20:52:48+01:00</updated>
					<link rel="alternate" type="text/html" href="/gogits/gogs/releases/tag/v0.11.34"/>
					<title>0.11.34</title>
					<content type="html"></content>
					<author>
						<name>Unknwon</name>
					</author>
					<media:thumbnail height="30" width="30" url="https://avatars0.githubusercontent.com/u/2946214?s=60&amp;v=4"/>
				</entry>
			</feed>
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
