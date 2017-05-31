package pwlb

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	//"strings"
)

const k_cycle_time = float64(50.0)

func GetTheoreticalMin() int {
	tc := GetTaskContainer()
	task_time_sum := float64(0.0)

	for i, _ := range tc.tasks {
		task_time_sum += tc.tasks[i].cost
	}

	min := int(math.Ceil(task_time_sum / k_cycle_time))
	return min
}

func PrettySolutionStr(sol *PartialSolution) string {
	wsc := GetWorkstationContainer()
	tc := GetTaskContainer()

	station_id := WorkstationId_t(tc.tasks[0].id)
	num_stations := len(wsc.workstations)

	str := ""
	for i := 0; i < num_stations; i++ {
		this_ws := wsc.workstations[station_id]

		// Skip an empty workstation
		if len(this_ws.tasks) == 0 {
			break
		}

		station_cost := this_ws.GetCost()

		// Sort unassigned tasks ids on task time
		sort.Slice(this_ws.tasks, func(a, b int) bool {
			lhs := this_ws.tasks[a]
			rhs := this_ws.tasks[b]
			return lhs < rhs
		})

		str += "Station " + strconv.Itoa(int(this_ws.id)) + ":      "
		str += "TaskTime " + strconv.FormatFloat(station_cost, 'f', 2, 64) + "   Tasks"
		for j, _ := range this_ws.tasks {
			str += " " + strconv.Itoa(int(this_ws.tasks[j]))
		}

		str += "\n"
		station_id++
	}

	return str
}

type TaskAssignment struct {
	i TaskId_t
	j WorkstationId_t
}

type PartialSolution struct {
	assignments []TaskAssignment
}

func (sol *PartialSolution) ToStr() string {
	str := "Solution Assignments:\n"
	for _, asg := range sol.assignments {
		str += "(" + strconv.Itoa(int(asg.i)) + "," + strconv.Itoa(int(asg.j)) + ")\n"
	}
	return str
}

func (sol *PartialSolution) TaskPrereqsMet(task *Task, ws_max WorkstationId_t) bool {
	// Check that every prereq for the task is in the solution
	for _, prereqid := range task.prereqs {
		found_prereq := false

		// Search all assignments for prereq
		for i, _ := range sol.assignments {
			// Limit search up to certain workstation
			if sol.assignments[i].j > ws_max {
				break
			}

			if sol.assignments[i].i == prereqid {
				found_prereq = true
				break
			}
		}

		if !found_prereq {
			return false
		}
	}

	return true
}

func (sol *PartialSolution) getNumActiveWorkstations() int {
	workstation_counter := make(map[WorkstationId_t]byte)
	for i, _ := range sol.assignments {
		workstation_counter[sol.assignments[i].j] = '0'
	}

	return len(workstation_counter)
}

func (sol *PartialSolution) GetMeasuredMin() int {
	return sol.getNumActiveWorkstations()
}

func (sol *PartialSolution) GetLineEfficiency() float64 {
	tc := GetTaskContainer()

	num_workstations := sol.getNumActiveWorkstations()
	if num_workstations == 0 {
		return 0.0
	}

	denom := k_cycle_time * float64(num_workstations)

	task_time_sum := 0.0
	for _, task_asg := range sol.assignments {
		task_time_sum += tc.GetTaskReadOnly(task_asg.i).cost
	}

	return task_time_sum / denom
}

func (sol *PartialSolution) GetLineEfficiencyStr() string {
	eff := sol.GetLineEfficiency() * 100.0
	str := strconv.FormatFloat(eff, 'f', 1, 64) + "%"
	return str
}

func (sol *PartialSolution) GetSmoothnessIndex() float64 {
	tc := GetTaskContainer()

	// Make a map to hold the sum task cost in a workstation and initialize it
	// Only workstations that have task assignments are initialized
	workstation_costs := make(map[WorkstationId_t]float64)
	for i, _ := range sol.assignments {
		workstation_costs[sol.assignments[i].j] = 0.0
	}

	// Calculate the sum of task costs at each used workstation
	for i, _ := range sol.assignments {
		task_cost := tc.tasks[tc.taskMapping[sol.assignments[i].i]].cost
		workstation_costs[sol.assignments[i].j] += task_cost
	}

	smoothness_acc := 0.0
	for _, val := range workstation_costs {
		smoothness_acc += math.Pow(k_cycle_time-val, 2.0)
	}

	smoothness_acc = math.Sqrt(smoothness_acc)
	return smoothness_acc
}

func (sol *PartialSolution) GetSmoothnessIndexStr() string {
	smooth := sol.GetSmoothnessIndex()
	str := strconv.FormatFloat(smooth, 'f', 1, 64)
	return str
}

func workstationCapacityRemaining(ws *Workstation) float64 {
	tc := GetTaskContainer()

	total_cost := 0.0
	for _, taskid := range ws.tasks {
		task_idx := tc.taskMapping[taskid]
		total_cost += tc.tasks[task_idx].cost
	}
	return k_cycle_time - total_cost
}

func IsSolutionValid(sol *PartialSolution) bool {
	tc := GetTaskContainer()
	wsc := GetWorkstationContainer()
	is_valid := true

	// Id of first workstation
	ws_start_id := WorkstationId_t(tc.tasks[0].id)

	// find id of first workstation with no tasks
	ws_end_id := ws_start_id
	for len(wsc.workstations[ws_end_id].tasks) != 0 {
		ws_end_id++
	}
	if ws_end_id == ws_start_id {
		fmt.Println("WARNING: Workstation ids not as expected in solution validation")
		return false
	}

	for i := ws_start_id; i < ws_end_id; i++ {
		if len(wsc.workstations[i].tasks) == 0 {
			fmt.Println("There should be no workstations in the solution range with no tasks")
			return false
		}

		if wsc.workstations[i].GetCost() > k_cycle_time {
			fmt.Println("Cycle time exceed on workstation " + strconv.Itoa(int(i)))
			return false
		}
	}

	// Verify prereqs met
	for i, _ := range sol.assignments {
		taskid := sol.assignments[i].i
		wsid := sol.assignments[i].j
		_task := tc.tasks[tc.taskMapping[taskid]]

		for _, prereqid := range _task.prereqs {
			found_prereq := false
			for cur_ws_id := ws_start_id; cur_ws_id <= wsid; cur_ws_id++ {
				if wsc.workstations[cur_ws_id].Contains(prereqid) {
					found_prereq = true
					break
				}
			}
			if !found_prereq {
				fmt.Println("Not all prereqs found for task " + _task.ToStr())
				return false
			}
		}
	}

	return is_valid
}

func ComputeSolutionSST() *PartialSolution {
	tc := GetTaskContainer()
	wsc := GetWorkstationContainer()
	unassigned := []*Task{}
	sol := PartialSolution{}

	// Fill container of unassigned task ids
	for i, _ := range tc.tasks {
		unassigned = append(unassigned, &tc.tasks[i])
	}

	// Sort unassigned tasks ids on task time
	sort.Slice(unassigned, func(i, j int) bool {
		lhs_cost := unassigned[i].cost
		rhs_cost := unassigned[j].cost
		return lhs_cost < rhs_cost
	})

	// ############################################
	// Assign tasks to stations using SST heuristic

	// Problem writeup states that workstations and tasks start with common index
	// The lowest id of workstation that can fit some unassigned task
	min_ws_id := WorkstationId_t(tc.tasks[0].id)
	// The highest id of workstations to check for met prereqs
	max_ws_id := WorkstationId_t(tc.tasks[0].id)
	// The current id of workstations targetted to fill
	cur_ws_id := WorkstationId_t(tc.tasks[0].id)

	// Assign tasks until there are none left unassigned
	for len(unassigned) != 0 {
		min_cost := unassigned[0].cost

		// Increment the lowest workstation to check, if needed
		for min_cost > workstationCapacityRemaining(wsc.workstations[min_ws_id]) {
			min_ws_id++
			if min_ws_id > max_ws_id {
				max_ws_id++
			}
		}
		cur_ws_id = min_ws_id

		// There is at least one task that can be assigned in the workstation range
		// given min_cost <= capacity(max_ws_id) enforced above

		// Attempt to fill lowest workstations first
		task_assigned := false
		grew_workstations := false
		for cur_ws_id <= max_ws_id {
			cur_cap := workstationCapacityRemaining(wsc.workstations[cur_ws_id])

			// Assign first task found that fits in the cur workstation if possible
			for i, task := range unassigned {
				// Assign task if conditions met
				if task.cost <= cur_cap && sol.TaskPrereqsMet(task, cur_ws_id) {
					// Add assignment to solution
					assignment := TaskAssignment{task.id, cur_ws_id}
					sol.assignments = append(sol.assignments, assignment)

					// Add task to workstation
					wstasks := &wsc.workstations[cur_ws_id].tasks
					*wstasks = append(*wstasks, task.id)

					// remove this taskid from unassigned, linear time, sad day.
					unassigned = append(unassigned[:i], unassigned[i+1:]...)

					task_assigned = true
					break
				}
			}

			if task_assigned {
				break
			}

			cur_ws_id++

			if cur_ws_id > max_ws_id && !grew_workstations {
				max_ws_id++
				grew_workstations = true
			}
		}
		if !grew_workstations && !task_assigned {
			panic("There is was no task assigned, this should never happen")
		}
	}

	// Validate solution found
	if !IsSolutionValid(&sol) {
		fmt.Println("WARNING Invalid solution detected")
	}

	// Perform qualitative analysis on solution found

	return &sol
}
