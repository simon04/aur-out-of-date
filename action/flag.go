package action

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/simon04/aur-out-of-date/pkg"
	"github.com/simon04/aur-out-of-date/upstream"
)

// FlagOnAur flags the package out-of-date after prompting the user
func FlagOnAur(pkg pkg.Pkg, upstreamVersion upstream.Version) {
	fmt.Printf("Should the package %s be flagged out-of-date? [y/N] ", pkg.Name())
	if !promptYesNo() {
		return
	}
	fmt.Printf("Flagging package %s out-of-date ...\n", pkg.Name())
	comment := fmt.Sprintf("Version %s is out. #simon04/aur-out-of-date", upstreamVersion)
	cmd := exec.Command("ssh", "aur@aur.archlinux.org", "flag", pkg.Name(), "\""+comment+"\"")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to flag out-of-date (running \"%v\"): %v\n%s\n", strings.Join(cmd.Args, "\" \""), err, output)
	} else {
		fmt.Printf("%s", output)
	}
}
