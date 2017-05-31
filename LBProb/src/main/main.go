package main

import (
	"fmt"
	"os"
	"pwlb"
	"strconv"
)

func main() {
	prog_args := os.Args[1:]

	// Read from file if there is an arg passed to program otherwise read from stdin
	var input_task_strings []string
	if len(prog_args) != 0 {
		specs_path := prog_args[0]
		input_task_strings = pwlb.ReadSpecFromPath(specs_path)
	} else {
		input_task_strings = pwlb.ReadSpecFromStdin()
	}

	// Create and initialize global task container
	tc := pwlb.GetTaskContainer()
	tc.FillFrom(input_task_strings)
	defer tc.Clear()

	// Create and initialize global workstation container
	wsc := pwlb.GetWorkstationContainer()
	wsc.FillFrom(tc)
	defer wsc.Clear()

	// Perform task assignments using SST algorithm
	sol := pwlb.ComputeSolutionSST()

	// Report results
	fmt.Println("theoretical_min=" + strconv.Itoa(pwlb.GetTheoreticalMin()))
	fmt.Println("measured_min=" + strconv.Itoa(sol.GetMeasuredMin()))
	fmt.Println("line_efficiency=" + sol.GetLineEfficiencyStr())
	fmt.Println("smoothness_index=" + sol.GetSmoothnessIndexStr())
	fmt.Print("\n\n")

	fmt.Print(pwlb.PrettySolutionStr(sol))
}
