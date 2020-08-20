package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func removeBranches(branches []string, removeChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	cnt := 0
	app := "git"
	for _, branchName := range branches {
		separatorArr := strings.Split(branchName, "/")
		if len(separatorArr) == 0 {
			continue
		}

		cnt++
		remoteName := separatorArr[0]
		remoteBranchName := strings.Join(separatorArr[1:], "/")
		args := []string{"push", remoteName, "--delete", remoteBranchName}
		cmd := exec.Command(app, args...)
		_, err := cmd.Output()

		fmt.Printf("Processing branch: %s\n", remoteBranchName)

		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	removeChan <- cnt
}

func removeStaleBranches(olderBranches []string, wg *sync.WaitGroup) int {
	removalLimt := 5
	removeChan := make(chan int, 20)

	for i := 0; i < len(olderBranches); i += removalLimt {
		wg.Add(1)
		go removeBranches(olderBranches[i:min(i+removalLimt, len(olderBranches))], removeChan, wg)
	}

	wg.Wait()
	close(removeChan)

	staleBranches := 0
	for cnt := range removeChan {
		staleBranches += cnt
	}

	return staleBranches
}
