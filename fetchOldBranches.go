package main

import (
	"os/exec"
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
		args := []string{"log", "-1", "--since=\"1 months ago\"", strippedBranchName}
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

func fetchOldBranches(outArr []string, wg *sync.WaitGroup) (olderBranches []string) {
	oldChan := make(chan []string, 20)
	fetchOlderBranchesLimit := 10

	for i := 0; i < len(outArr); i += fetchOlderBranchesLimit {
		wg.Add(1)
		go filterOlderBranches(outArr[i:min(i+fetchOlderBranchesLimit, len(outArr))], oldChan, wg)
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
