package upstream

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/simon04/aur-out-of-date/pkg"
)

func mockPython() *gock.Response {
	return gock.New("https://pypi.python.org/").
		Get("/pypi/httpie/json").
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
			"info": {
				"package_url": "http://pypi.python.org/pypi/httpie",
				"download_url": "https://github.com/jkbrzt/httpie",
				"platform": "UNKNOWN",
				"version": "0.9.9",
				"release_url": "http://pypi.python.org/pypi/httpie/0.9.9"
			}
		}`)
}

func TestPythonHttpieSource1(t *testing.T) {
	testPythonHttpie(t, "https://pypi.python.org/packages/source/h/httpie/httpie-0.9.9.tar.gz")
}

func TestPythonHttpieSource2(t *testing.T) {

}

func TestPythonHttpieSource3(t *testing.T) {
	// URL from https://pypi.org/project/httpie/#files
	testPythonHttpie(t, "https://files.pythonhosted.org/packages/28/93/4ebf2de4bc74bd517a27a600b2b23a5254a20f28e6e36fc876fd98f7a51b/httpie-0.9.9.tar.gz")
}

func testPythonHttpie(t *testing.T, url string) {
	defer gock.Off()
	mockPython()

	p := pkg.New("httpie", "0", "", url)
	version, err := VersionForPkg(p)
	if err != nil {
		t.Error(err)
	}
	if version != "0.9.9" {
		t.Errorf("Expecting version 0.9.9, but got %v", version)
	}
}
