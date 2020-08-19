package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func filterOlderBranches(branchArr []string, oldChan chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	app := "git"
	ret := []string{}
	for _, branchName := range branchArr {
		if len(branchName) == 0 {
			continue
		}

		strippedBranchName := branchName[1 : len(branchName)-1]
		args := []string{"log", "-1", "--since=\"1 year ago\"", strippedBranchName}
		cmd := exec.Command(app, args...)
		stdOut, err := cmd.Output()

		if err != nil {
			panic(err)
		}

		if len(stdOut) == 0 {
			ret = append(ret, strippedBranchName)
		}
	}

	oldChan <- ret
}

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

func fetchOldBranches(outArr []string, wg *sync.WaitGroup) (olderBranches []string) {
	oldChan := make(chan []string)

	for i := 0; i < len(outArr); i += 10 {
		wg.Add(1)
		go filterOlderBranches(outArr[i:min(i+10, len(outArr))], oldChan, wg)
	}

	go func() {
		for res := range oldChan {
			olderBranches = append(olderBranches, res...)
		}
	}()

	wg.Wait()
	close(oldChan)

	return
}

func removeStaleBranches(olderBranches []string, wg *sync.WaitGroup) int {
	removeChan := make(chan int)

	for i := 0; i < len(olderBranches); i += 10 {
		wg.Add(1)
		go removeBranches(olderBranches[i:min(i+10, len(olderBranches))], removeChan, wg)
	}

	staleBranches := 0
	go func() {
		for cnt := range removeChan {
			staleBranches += cnt
		}
	}()

	wg.Wait()
	close(removeChan)

	return staleBranches
}

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

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
