package action

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/simon04/aur-out-of-date/pkg"
	"github.com/simon04/aur-out-of-date/upstream"
)

// UpdatePKGBUILD updates pkgver/pkgrel in local PKGBUILD files after prompting the user
func UpdatePKGBUILD(pkg pkg.Pkg, upstreamVersion upstream.Version) {
	file := pkg.LocalPKGBUILD()
	if file == "" {
		return
	}
	input, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("updatePKGBUILD: failed to read file %s: %v\n", file, err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "pkgver=") {
			lineUpdate := strings.Replace(line, string(pkg.Version().Version), upstreamVersion.String(), 1)
			fmt.Printf("--- a/%s\n", file)
			fmt.Printf("+++ b/%s\n", file)
			fmt.Printf("-%s\n", line)
			fmt.Printf("+%s\n", lineUpdate)
			fmt.Printf("Should the package %s be updated to version %s? [y/N] ", pkg.Name(), upstreamVersion)
			if !promptYesNo() {
				return
			}
			lines[i] = lineUpdate
		} else if strings.HasPrefix(line, "pkgrel=") {
			lines[i] = "pkgrel=1"
		}
	}
	err = ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		fmt.Printf("updatePKGBUILD: failed to write file %s: %v\n", file, err)
	}
}
