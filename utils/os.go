package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func RootPath() string {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s", strings.TrimSuffix(string(cmdOut), "\n"))
}
