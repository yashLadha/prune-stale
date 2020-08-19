package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	app := "git"

	args := []string{"branch", "-r", "--format=\"%(refname:short)\""}

	cmd := exec.Command(app, args...)
	stdOut, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	strOutput := string(stdOut)
	outArr := strings.Split(strOutput, "\n")

	var wg sync.WaitGroup
	olderBranches := fetchOldBranches(outArr, &wg)
	staleBranches := removeStaleBranches(olderBranches, &wg)
	fmt.Printf("Stale branches %d\n", staleBranches)
}
