package main

import (
	"fmt"
	"sync"
)

func main() {
	outArr := fetchBranches()
	var wg sync.WaitGroup
	olderBranches := fetchOldBranches(outArr, &wg)
	staleBranches := removeStaleBranches(olderBranches, &wg)
	fmt.Printf("Stale branches %d\n", staleBranches)
}
