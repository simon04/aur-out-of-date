package action

import "fmt"

func promptYesNo() bool {
	var response string
	chars, err := fmt.Scanln(&response)
	if err != nil || chars == 0 {
		return false
	}
	if response != "y" && response != "Y" {
		return false
	}
	return true
}
