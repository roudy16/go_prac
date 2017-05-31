package pwlb

import (
	"strconv"
	"strings"
)

type TaskId_t int
type WorkstationId_t int

type Task struct {
	id       TaskId_t
	cost     float64    // time in seconds to complete task
	prereqs  []TaskId_t // other tasks that must be completed prior this one being started
	assigned bool
}

type Workstation struct {
	id    WorkstationId_t
	tasks []TaskId_t // tasks assigned to this workstation
}

func NewTask(_id TaskId_t, _c float64) *Task {
	return &Task{id: _id, cost: _c}
}

// Returns a task as a string like "2,20.2,nil" or "4,3.3,2 3"
func (t *Task) ToStr() string {
	str := strconv.Itoa(int(t.id))
	str += "," + strconv.FormatFloat(float64(t.cost), 'f', 1, 64) + ","
	str += taskIdsToStr(t.prereqs)
	return strings.TrimSpace(str)
}

func NewWorkstation(_id WorkstationId_t) *Workstation {
	return &Workstation{id: _id}
}

// Returns the sum cost of all tasks assigned to workstation
func (ws *Workstation) GetCost() float64 {
	tc := GetTaskContainer()

	cost := 0.0
	for _, taskid := range ws.tasks {
		cost += tc.tasks[tc.taskMapping[taskid]].cost
	}
	return cost
}

func (ws *Workstation) Contains(tid TaskId_t) bool {
	for _, taskid := range ws.tasks {
		if taskid == tid {
			return true
		}
	}
	return false
}

// Returns a workstation as a string like "1 45.4:0 2 3" or "0 0.0:nil"
func (ws *Workstation) ToStr() string {
	str := strconv.Itoa(int(ws.id))
	str += " " + strconv.FormatFloat(float64(ws.GetCost()), 'f', 1, 64) + ":"
	str += taskIdsToStr(ws.tasks)
	return strings.TrimSpace(str)
}

///////////////////////////////////
// Extra for alternate heuristic
///////////////////////////////////

type PostReqNode struct {
	task_id  TaskId_t
	branches []*PostReqNode
}

// Represents a graph of follow on dependencies, a DAG
type PostReqGraph []*PostReqNode

func (node *PostReqNode) IsLeaf() bool {
	return len(node.branches) == 0
}

func (node *PostReqNode) HeightWithPath() (int, []TaskId_t) {
	this_path := []TaskId_t{node.task_id}
	if node.IsLeaf() {
		return 0, this_path
	}

	max_branch_height := 0
	max_branch_path := []TaskId_t{node.branches[0].task_id}
	for i, _ := range node.branches {
		branch_height, branch_path := node.branches[i].HeightWithPath()
		if branch_height > max_branch_height {
			max_branch_height = branch_height
			max_branch_path = branch_path
		}
	}

	return (1 + max_branch_height), (append(this_path, max_branch_path...))
}

func (node *PostReqNode) contains(query TaskId_t) bool {
	if node.task_id == query {
		return true
	}

	for _, subnode := range node.branches {
		if subnode.contains(query) {
			return true
		}
	}
	return false
}

//func (g *PostReqGraph) MaxHeight()
