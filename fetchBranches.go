package main

import (
	"os/exec"
	"strings"
)

func fetchBranches() (outArr []string) {
	app := "git"
	args := []string{"branch", "-r", "--format=\"%(refname:short)\""}

	cmd := exec.Command(app, args...)
	stdOut, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	strOutput := string(stdOut)
	outArr = strings.Split(strOutput, "\n")
	return
}
