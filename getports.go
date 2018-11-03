package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func runCmd(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()
	_ = err

	return string(out)
}

func getAvailablePort(start, max int) int {
	cmd := "lsof -i"
	o := runCmd(cmd)
	


	return 0
}

func main() {
	getAvailablePort(0, 0)
}