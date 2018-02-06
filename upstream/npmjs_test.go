package upstream

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockNpm(version string) *gock.Response {
	return gock.New("https://registry.npmjs.org/").
		Get("/-/package/webpack/dist-tags").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		BodyString(fmt.Sprintf(`{"latest":"%v","legacy":"1.15.0"}`, version))
}

func TestWebpackNpmUrl(t *testing.T) {
	defer gock.Off()
	mockNpm("3.9.0")

	p := pkg.New("webpack", "3.6.0", "https://www.npmjs.com/package/webpack")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "3.9.0" {
		t.Errorf("Expecting version 3.9.0, but got %v", version)
	}
}

func TestWebpackNpmSource(t *testing.T) {
	defer gock.Off()
	mockNpm("4.0.0")

	p := pkg.New("webpack", "3.6.0", "https://webpack.js.org/", "http://registry.npmjs.org/webpack/-/webpack-3.6.0.tgz")
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "4.0.0" {
		t.Errorf("Expecting version 4.0.0, but got %v", version)
	}
}

func TestWebpackNoSource(t *testing.T) {
	p := pkg.New("webpack", "3.6.0", "https://webpack.js.org/")
	_, err := VersionForPkg(p)
	if err == nil {
		t.Error("Expecting an error, but got none")
	}
}

func TestWebpackUnknownSource(t *testing.T) {
	p := pkg.New("webpack", "3.6.0", "https://webpack.js.org/", "http://webpack.js.org/webpack-3.6.0.tgz")
	_, err := VersionForPkg(p)
	if err == nil {
		t.Error("Expecting an error, but got none")
	}
}
