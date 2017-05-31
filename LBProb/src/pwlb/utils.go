package pwlb

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
func (this *Task) HasPriorityOver(that *Task) bool {

}
*/

func getRuntimeExePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("Failed to determine path to runtime executable")
	}

	return dir
}

func readSpecFromFileBuffer(f *os.File) []string {
	var task_strings []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Read each line of text and append it to the list of task strings
		line := scanner.Text()

		// Stop reading if an empty line is read
		if line == "" {
			break
		}

		task_strings = append(task_strings, scanner.Text())
	}

	return task_strings
}

func ReadSpecFromStdin() []string {
	return readSpecFromFileBuffer(os.Stdin)
}

func ReadSpecFromPath(spec_path string) []string {
	specs_file, err := os.Open(spec_path)
	if err != nil {
		log.Fatal("Failed to read file at " + spec_path)
	}

	return readSpecFromFileBuffer(specs_file)
}

func AreStringsSame(lhs, rhs []string) bool {
	if lhs == nil && rhs == nil {
		return true
	}

	if lhs == nil || rhs == nil {
		return false
	}

	if len(lhs) != len(rhs) {
		return false
	}

	for i, ele := range lhs {
		if ele != rhs[i] {
			return false
		}
	}

	return true
}

func AreTaskIdsSame(lhs, rhs []TaskId_t) bool {
	if lhs == nil && rhs == nil {
		return true
	}

	if lhs == nil || rhs == nil {
		return false
	}

	if len(lhs) != len(rhs) {
		return false
	}

	for i, ele := range lhs {
		if ele != rhs[i] {
			return false
		}
	}

	return true
}

func taskIdsToStr(ids []TaskId_t) string {
	str := ""
	if len(ids) == 0 {
		str += "nil"
		return str
	}

	for _, ele := range ids {
		str += strconv.Itoa(int(ele)) + " "
	}

	return strings.TrimSpace(str)
}

// ###########
// Functions for getting permutations of "working set" of task ids
// ###########

func npSwap(taskIds *[]TaskId_t, l, r int) {
	if l == r {
		return
	}
	(*taskIds)[l] ^= (*taskIds)[r]
	(*taskIds)[r] ^= (*taskIds)[l]
	(*taskIds)[l] ^= (*taskIds)[r]
}

func npReverse(taskIds *[]TaskId_t, l, r int) {
	for l < r {
		npSwap(taskIds, l, r)
		l++
		r--
	}
}

func npBSearch(taskIds *[]TaskId_t, query TaskId_t, l, r int) int {
	i := -1
	for l <= r {
		mid := l + (r-1)/2
		if (*taskIds)[mid] <= query {
			r = mid - 1
		} else {
			l = mid + 1
			if i == -1 || (*taskIds)[i] <= (*taskIds)[mid] {
				i = mid
			}
		}
	}
	return i
}

// Find the next highest permutation in a slice of task ids from a given index
// Returns false if there is no next permutation. Boundary indices are inclusive
func npNextPerm(taskIds *[]TaskId_t, left_bound_idx int) bool {
	right_bound_idx := len(*taskIds) - 1
	cur_idx := right_bound_idx - 1

	// If at last digit we cannot make another permutation
	if left_bound_idx >= right_bound_idx {
		return false
	}

	// Find index of first id that is less than an id to the right
	for cur_idx >= left_bound_idx && (*taskIds)[cur_idx] >= (*taskIds)[cur_idx+1] {
		cur_idx--
	}

	// If no such id was found then there is no next permutation
	if cur_idx < left_bound_idx {
		return false
	}

	// Find the index of id with the least value among the ids greater than the one we found
	swap_idx := npBSearch(taskIds, (*taskIds)[cur_idx], cur_idx+1, right_bound_idx)

	// Swap the id at the current index and the one with the next higher value
	npSwap(taskIds, cur_idx, swap_idx)

	// Reverse the sequence of ids after the current index
	npReverse(taskIds, cur_idx+1, right_bound_idx)

	return true
}

// ###########
// End permutations section
// ###########
