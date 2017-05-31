package pwlb

import (
	"log"
	"strconv"
	"strings"
	"sync"
)

///////////////////////////////////
// TaskContainer Singleton
///////////////////////////////////

type TaskContainer struct {
	tasks       []Task
	taskMapping map[TaskId_t]int

	// Extra for alternate heuristic
	postreqs map[TaskId_t]PostReqGraph
}

var task_container_instance *TaskContainer
var task_container_once sync.Once
var task_container_mutex sync.RWMutex

func GetTaskContainer() *TaskContainer {
	task_container_once.Do(func() {
		task_container_instance = &TaskContainer{}
		task_container_instance.taskMapping = make(map[TaskId_t]int)

		// Extra for alternate heuristic
		task_container_instance.postreqs = make(map[TaskId_t]PostReqGraph)
	})
	return task_container_instance
}

func (tc *TaskContainer) GetTaskReadOnly(id TaskId_t) Task {
	some_task := tc.tasks[tc.taskMapping[id]]
	return some_task
}

func (tc *TaskContainer) FillFrom(task_list []string) {
	// Each element in the passed in string sequence should represent a task.
	// We parse each element and construct a Task object to add to the container
	for _, ele := range task_list {
		// Trim any leading or trailing whitespace from each task string
		// then split into a sequence of substrings that represent task fields
		trimmed_ele := strings.TrimSpace(ele)
		task_fields := strings.Split(trimmed_ele, ",")
		if len(task_fields) != 3 {
			log.Fatal("Improper task format in task: " + ele)
		}

		// The first to fields are id and cost of the task
		id, _ := strconv.Atoi(task_fields[0])
		cost, _ := strconv.ParseFloat(task_fields[1], 64)

		// Remaining field is whitespace seperated list of prereq tasks, we
		// must parse those prereqs
		var prereqs []TaskId_t
		for _, tidstr := range strings.Split(strings.TrimSpace(task_fields[2]), " ") {
			tid, err := strconv.Atoi(tidstr)

			// If err found we assume no more prereqs and exit the loop
			if err != nil {
				break
			}

			// Otherwise append new prereq to the list
			prereqs = append(prereqs, TaskId_t(tid))
		}

		// Initialize a new Task from the parsed string
		new_task := NewTask(TaskId_t(id), cost)
		new_task.prereqs = append(new_task.prereqs, prereqs...)

		// Add the new task to the container
		task_idx := len(tc.tasks)
		tc.tasks = append(tc.tasks, *new_task)
		tc.taskMapping[new_task.id] = task_idx
	}
}

func (tc *TaskContainer) Clear() {
	tc.tasks = nil
	tc.taskMapping = make(map[TaskId_t]int)

	// Extra for alternate heuristic
	tc.postreqs = make(map[TaskId_t]PostReqGraph)
}

func (tc *TaskContainer) ToStrArr() []string {
	var strs []string
	for _, ele := range tc.tasks {
		strs = append(strs, ele.ToStr())
	}
	return strs
}

///////////////////////////////////
// WorkstationContainer Singleton
///////////////////////////////////

type WorkstationContainer struct {
	workstations map[WorkstationId_t]*Workstation
}

var ws_container_instance *WorkstationContainer
var ws_container_once sync.Once
var ws_container_mutex sync.RWMutex

func GetWorkstationContainer() *WorkstationContainer {
	ws_container_once.Do(func() {
		ws_container_instance = &WorkstationContainer{}
		ws_container_instance.workstations = make(map[WorkstationId_t]*Workstation)
	})
	return ws_container_instance
}

func (wsc *WorkstationContainer) GetWorkstationReadOnly(id WorkstationId_t) Workstation {
	some_ws := wsc.workstations[id]
	return *some_ws
}

func (wsc *WorkstationContainer) FillFrom(tc *TaskContainer) {
	for i, _ := range tc.tasks {
		new_id := WorkstationId_t(tc.tasks[i].id)
		wsc.workstations[new_id] = NewWorkstation(new_id)
	}
}

func (wsc *WorkstationContainer) Clear() {
	wsc.workstations = make(map[WorkstationId_t]*Workstation)
}

func (wsc *WorkstationContainer) clearStationTasks(id WorkstationId_t) {
	wsc.workstations[id].tasks = nil
}

func (wsc *WorkstationContainer) clearAllStationsTasks() {
	for k, _ := range wsc.workstations {
		wsc.clearStationTasks(k)
	}
}

func (wsc *WorkstationContainer) ToStrArr() []string {
	var strs []string
	for _, ele := range wsc.workstations {
		strs = append(strs, ele.ToStr())
	}
	return strs
}

///////////////////////////////////
// TaskContainer Extra for alternate heuristic
///////////////////////////////////

func buildPostReqGraphs() {
	tc := GetTaskContainer()

	// Iterating front to back we establish the first link in all postreq chains
	for i, _ := range tc.tasks {
		for _, prereqId := range tc.tasks[i].prereqs {
			new_node := PostReqNode{tc.tasks[i].id, nil}
			tc.postreqs[prereqId] = append(tc.postreqs[prereqId], &new_node)
		}
	}

	// Iterating back to front we append current task's postreq nodes onto any task
	// that has the current task as the first postreq
	num_tasks := len(tc.tasks)
	for i := num_tasks - 1; i >= 0; i-- {
		cur_task := &tc.tasks[i]

		// Get all the prereq node for the current node, the current should be
		// in the postreqs for each prereq node
		for _, cur_prereqId := range cur_task.prereqs {
			// Get the postreq nodes to append
			for j, _ := range tc.postreqs[cur_prereqId] {
				cur_prereq_postreq_node := tc.postreqs[cur_prereqId][j]

				// skip nodes that don't map to the current task
				if cur_prereq_postreq_node.task_id != cur_task.id {
					continue
				}

				// append the current tasks postreq graph onto the node in the prereq's
				// postreq graph representing the current node.
				for k, _ := range tc.postreqs[cur_task.id] {
					cur_prereq_postreq_node.branches = append(cur_prereq_postreq_node.branches,
						(tc.postreqs[cur_task.id])[k])
				}
			}
		}
	}
}
