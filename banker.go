// Banker's Algorithm implementation
//
// For more info, see:
// - [EN] https://en.wikipedia.org/wiki/Banker%27s_algorithm
// - [IT] https://it.wikipedia.org/wiki/Algoritmo_del_banchiere
//
// Made with ♥ by Simone (https://github.com/BuriedInTheGround)

package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	nOfResources = 20 // Number of resources
	minResQty    = 5  // Minimum quantity for every resource
	maxResQty    = 10 // Maximum quantity for every resource

	nOfProcesses = 20 // Number of processes

	minAssignedRes = 0 // Minimum quantity of resources assigned per-process
	maxAssignedRes = 2 // Maximum quantity of resources assigned per-process

	// WARNING: minTotalRes MUST BE at least equals to maxAssignedRes
	minTotalRes = 4 // Minimum quantity of resources used by a process
	maxTotalRes = 8 // Maximum quantity of resources used by a process
)

var (
	available    []int    // Number of available resources for every resource
	assigned     [][]int  // Number of assigned resource per-process
	total        [][]int  // Number of total resource usage per-process
	necessary    [][]int  // Number of remaining necessary resources per-process
	safeSequence []string // Safe sequence
)

// Flags
var (
	debugOn bool
	runs    int
)

func init() {
	// Setup flags
	flag.BoolVar(&debugOn, "debug", false, "Activate debug messages.")
	flag.IntVar(&runs, "runs", 1, "Number of runs.")
}

func initRun() {
	// Initialize variables
	available = make([]int, nOfResources)
	assigned = make([][]int, nOfProcesses)
	total = make([][]int, nOfProcesses)
	necessary = make([][]int, nOfProcesses)
	safeSequence = make([]string, 0)

	// Reset random seed
	rand.Seed(time.Now().UnixNano())

	// Generate resources quantities
	for i := range available {
		available[i] = (rand.Int() % (maxResQty - minResQty + 1)) + minResQty
	}

	// Generate assigned resources for every process
	for i := range assigned {
		assigned[i] = make([]int, nOfResources)
		for j := range assigned[i] {
			assigned[i][j] = (rand.Int() % (maxAssignedRes - minAssignedRes + 1)) + minAssignedRes
		}
	}

	if minTotalRes < maxAssignedRes {
		fmt.Println("Error: minTotalRes MUST BE at least equals to maxAssignedRes")
		os.Exit(1)
	}

	// Generate total number of resources for every process
	for i := range total {
		total[i] = make([]int, nOfResources)
		for j := range total[i] {
			total[i][j] = (rand.Int() % (maxTotalRes - minTotalRes + 1)) + minTotalRes
		}
	}

	// Generate necessary remaining resources
	for i := range necessary {
		necessary[i] = make([]int, nOfResources)
		for j := range necessary[i] {
			necessary[i][j] = total[i][j] - assigned[i][j]
		}
	}
}

func main() {
	flag.Parse()

	for i := 0; i < runs; i++ {
		initRun()

		fmt.Printf("Run %v\n", i)
		fmt.Println("-----")

		safe, err := safeState(i)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if safe {
			fmt.Println("Safe sequence found!")
			fmt.Printf("Safe sequence: <")
			for i, p := range safeSequence {
				fmt.Printf("%v", p)
				if i < nOfProcesses-1 {
					fmt.Printf(" ")
				}
			}
			fmt.Println(">")
		} else {
			fmt.Println("Safe sequence NOT found! :(")
		}

		fmt.Printf("\n")
	}
}

func safeState(run int) (bool, error) {
	if debugOn {
		defer elapsed(fmt.Sprintf("safeState run %v", run))()
	}

	var (
		isSufficient = false
		need         = make([]int, nOfResources)
		done         = make([]bool, nOfProcesses)
		doneCounter  = 0
		failCounter  = 0
	)

	for p := 0; ; p++ {
		need = necessary[p]
		isSufficient = true

		// Check if resources are sufficient to end this process
		for r := range need {
			if need[r] > available[r] {
				isSufficient = false
			}
		}

		if !done[p] && isSufficient {
			if debugOn {
				fmt.Printf("[DEBUG] resources are sufficient: freeing resource (P%v)\n", p)
			}
			for rindex, rqty := range assigned[p] {
				available[rindex] += rqty
			}
			done[p] = true
			doneCounter++
			safeSequence = append(safeSequence, fmt.Sprintf("P%v", p))
		} else if doneCounter >= nOfProcesses {
			return true, nil
		} else if failCounter >= nOfProcesses-1 {
			return false, nil
		} else {
			if debugOn {
				if !done[p] {
					fmt.Printf("[DEBUG] resources are not sufficient (P%v)\n", p)
				} else {
					fmt.Printf("[DEBUG] process already done (P%v)\n", p)
				}
			}
			failCounter++
		}

		// Restart if was the last iteration to redo checks
		if p == len(necessary)-1 {
			p = -1
			failCounter = 0
		}
	}

	return false, errors.New("something went wrong")
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("[DEBUG] %s took %v\n", what, time.Since(start))
	}
}
