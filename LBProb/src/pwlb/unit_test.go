package pwlb

import (
	//"log"
	//"os"
	"path/filepath"
	//"strings"
	"testing"
)

var tc = GetTaskContainer()
var wsc = GetWorkstationContainer()

func getSpecFilesDir() string {
	exe_path := getRuntimeExePath()
	spec_path, _ := filepath.Abs(exe_path + "/../specs/")
	return spec_path + "/"
}

func TestAreStringsSame(t *testing.T) {
	lhs0 := []string{"apple", "monkey"}
	rhs0 := []string{"apple", "monkey"}
	res0 := AreStringsSame(lhs0, rhs0)

	if res0 != true {
		t.Error("Expected []string objects to evaluate as same. LHS: ",
			lhs0, " RHS: ", rhs0)
	}

	lhs1 := []string{}
	rhs1 := []string{}
	res1 := AreStringsSame(lhs1, rhs1)

	if res1 != true {
		t.Error("Expected []string objects to evaluate as same. LHS: ",
			lhs1, " RHS: ", rhs1)
	}

	lhs2 := []string{"Data"}
	rhs2 := []string{}
	res2 := AreStringsSame(lhs2, rhs2)

	if res2 != false {
		t.Error("Expected []string objects to evaluate as different. LHS: ",
			lhs2, " RHS: ", rhs2)
	}

	lhs3 := []string{"Data"}
	rhs3 := []string{"Data", "ExtraData"}
	res3 := AreStringsSame(lhs3, rhs3)

	if res3 != false {
		t.Error("Expected []string objects to evaluate as different. LHS: ",
			lhs3, " RHS: ", rhs3)
	}

	lhs4 := []string{"ExtraData", "Data"}
	rhs4 := []string{"Data", "ExtraData"}
	res4 := AreStringsSame(lhs4, rhs4)

	if res4 != false {
		t.Error("Expected []string objects to evaluate as different. LHS: ",
			lhs4, " RHS: ", rhs4)
	}
}

var test_spec1 = []string{"0,12.5,nil",
	"1,20.0,0",
	"2,2.3,nil",
	"3,9.5,nil",
	"4,25.3,1 3"}

func TestSpec1Parse(t *testing.T) {
	// ##########
	file_tasks := test_spec1
	tc.FillFrom(file_tasks)
	defer tc.Clear()

	tc_tasks := tc.ToStrArr()

	res := AreStringsSame(file_tasks, tc_tasks)
	if res != true {
		t.Error("Tasks in container did not match tasks from file",
			file_tasks, tc_tasks)
	}
}

func TestTheoreticalMin(t *testing.T) {
	// ##########
	file_tasks := test_spec1
	tc.FillFrom(file_tasks)
	defer tc.Clear()

	min := GetTheoreticalMin()
	expected := 2

	if min != expected {
		t.Error("Expected min: ", expected, " got: ", min)
	}
}

func TestSolutionValidation(t *testing.T) {
	// ##########
	file_tasks := test_spec1
	tc.FillFrom(file_tasks)
	defer tc.Clear()
	wsc.FillFrom(tc)
	defer wsc.Clear()

	ta0_0 := TaskAssignment{1, 0}
	ta0_1 := TaskAssignment{2, 0}
	ta0_2 := TaskAssignment{3, 0}
	ta0_3 := TaskAssignment{0, 1}
	ta0_4 := TaskAssignment{4, 1}
	bad_sol0 := PartialSolution{}
	bad_sol0.assignments = append(bad_sol0.assignments, ta0_0)
	bad_sol0.assignments = append(bad_sol0.assignments, ta0_1)
	bad_sol0.assignments = append(bad_sol0.assignments, ta0_2)
	bad_sol0.assignments = append(bad_sol0.assignments, ta0_3)
	bad_sol0.assignments = append(bad_sol0.assignments, ta0_4)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 1)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 2)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 3)
	wsc.workstations[1].tasks = append(wsc.workstations[1].tasks, 0)
	wsc.workstations[1].tasks = append(wsc.workstations[1].tasks, 4)

	isValid := IsSolutionValid(&bad_sol0)
	if isValid {
		t.Error("Solution expected to fail validation but didn't  " + bad_sol0.ToStr())
	}

	wsc.workstations[0].tasks = nil
	wsc.workstations[1].tasks = nil

	ta1_0 := TaskAssignment{0, 0}
	ta1_1 := TaskAssignment{1, 0}
	ta1_2 := TaskAssignment{2, 0}
	ta1_3 := TaskAssignment{3, 0}
	ta1_4 := TaskAssignment{4, 1}
	bad_sol1 := PartialSolution{}
	bad_sol1.assignments = append(bad_sol1.assignments, ta1_0)
	bad_sol1.assignments = append(bad_sol1.assignments, ta1_1)
	bad_sol1.assignments = append(bad_sol1.assignments, ta1_2)
	bad_sol1.assignments = append(bad_sol1.assignments, ta1_3)
	bad_sol1.assignments = append(bad_sol1.assignments, ta1_4)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 0)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 1)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 2)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 3)
	wsc.workstations[1].tasks = append(wsc.workstations[1].tasks, 4)

	isValid = IsSolutionValid(&bad_sol1)
	if !isValid {
		t.Error("Solution failed validation  " + bad_sol1.ToStr())
	}

	wsc.workstations[0].tasks = nil
	wsc.workstations[1].tasks = nil

	ta2_0 := TaskAssignment{0, 0}
	ta2_1 := TaskAssignment{2, 0}
	ta2_2 := TaskAssignment{1, 1}
	ta2_3 := TaskAssignment{3, 1}
	ta2_4 := TaskAssignment{4, 1}
	bad_sol2 := PartialSolution{}
	bad_sol2.assignments = append(bad_sol2.assignments, ta2_0)
	bad_sol2.assignments = append(bad_sol2.assignments, ta2_1)
	bad_sol2.assignments = append(bad_sol2.assignments, ta2_2)
	bad_sol2.assignments = append(bad_sol2.assignments, ta2_3)
	bad_sol2.assignments = append(bad_sol2.assignments, ta2_4)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 0)
	wsc.workstations[0].tasks = append(wsc.workstations[0].tasks, 2)
	wsc.workstations[1].tasks = append(wsc.workstations[1].tasks, 1)
	wsc.workstations[1].tasks = append(wsc.workstations[1].tasks, 3)
	wsc.workstations[1].tasks = append(wsc.workstations[1].tasks, 4)

	isValid = IsSolutionValid(&bad_sol2)
	if isValid {
		t.Error("Expected solution to fail validation due to  " + bad_sol2.ToStr())
	}

}

///////////////////////////////////
// Testing Extra for alternate heuristic
///////////////////////////////////

func TestPostReqNode(t *testing.T) {
	// ##########
	// Setup a tiny graph to test finding height and longest path
	node5 := PostReqNode{5, nil}
	node4 := PostReqNode{4, []*PostReqNode{&node5}}
	node3 := PostReqNode{3, nil}
	node2 := PostReqNode{2, []*PostReqNode{&node4}}
	node1 := PostReqNode{1, []*PostReqNode{&node3}}
	node0 := PostReqNode{0, []*PostReqNode{&node1, &node2}}

	height, path := node0.HeightWithPath()
	if height != 3 {
		t.Error("Node height wrong, expected: ", 3, " got: ", height)
	}

	if !AreTaskIdsSame(path, []TaskId_t{0, 2, 4, 5}) {
		t.Error("Expected path 0, 2, 4, 5 got: " + taskIdsToStr(path))
	}
}

func TestNextPermutation(t *testing.T) {
	// ##########
	// Test finding next permutation span starting at index 1 to end of slice
	ids0 := []TaskId_t{1, 3, 4, 2, 0}
	expect0 := []TaskId_t{1, 4, 0, 2, 3}

	hasNextPerm := npNextPerm(&ids0, 1)
	if !hasNextPerm {
		t.Error("Expected another permutation of " + taskIdsToStr(ids0) + " to be found")
	}

	if !AreTaskIdsSame(ids0, expect0) {
		t.Error("Expected next perm to be: " + taskIdsToStr(expect0) + "got: " + taskIdsToStr(ids0))
	}

	// ##########
	// Test sorting of behavior when no next permutation is possible
	ids1 := []TaskId_t{1, 4, 3, 2, 0}
	expect1 := []TaskId_t{1, 4, 3, 2, 0}

	hasNextPerm = npNextPerm(&ids1, 1)
	if hasNextPerm {
		t.Error("Expected another permutation of " + taskIdsToStr(ids1) + " NOT to be found")
	}

	if !AreTaskIdsSame(ids1, expect1) {
		t.Error("Expected next perm to be: " + taskIdsToStr(expect1) + "got: " + taskIdsToStr(ids1))
	}
}
