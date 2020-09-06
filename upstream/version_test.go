package upstream

import (
	"testing"
)

func TestScript1(t *testing.T) {
	version, err := VersionForScript("echo 42")
	if err != nil {
		t.Error(err)
	}
	if version != "42" {
		t.Errorf("Expecting version 42, but got %v", version)
	}
}
