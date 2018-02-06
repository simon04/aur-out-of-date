package pkg

import (
	"strings"
	"testing"

	"github.com/mikkeloscar/aur"
)

func TestSplitPkg(t *testing.T) {
	info, err := aur.Info([]string{"argtable-docs"})
	if err != nil {
		t.Error(err)
	}
	pkgs := NewRemotePkgs(info)
	sources, err := pkgs[0].Sources()
	if err != nil {
		t.Error(err)
	}
	if len(sources) == 0 || !strings.Contains(sources[0], "sourceforge.net") {
		t.Errorf("Unexpected sources %v", sources)
	}
}
