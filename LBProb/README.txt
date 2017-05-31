The line balancing problem is a hard problem similar to the bin packing problem but with an additional constraint. In the line balancing problem each workstation can be viewed as a bin with a capacity of 50 units, each task is an object that must be packed (or assigned) into a bin (or workstation). The added constraint is that there is an ordering the tasks must follow, in other words some tasks have prerequisite tasks that must be completed before the task in question can be started.

My goal for the line balancing problem is to implement an algorithm that analyzes a given set of tasks and their prerequisites to provide an optimized distribution of the tasks among a minimal number of workstations. Each possible distribution of tasks among workstations is a potential solution. Potential solutions are evaluated by two metrics: line efficiency and smoothness index.

Build:

   make runall

The above command builds my program, runs it with the provided specs and stores the output files.

Test:

   make test-pwlb

The above command runs some Go unit tests I made.

Design Overview:

There are two main containers, one holds objects representing tasks and is used as a global source of the task data. The other container holds objects representing the workstations. There are also data structures representing the following:

	Tasks - Have Id and cost
	Workstations - Have Id and container of Tasks
	TaskAssignments - A pairing of a Task and a Workstation
	Solutions - A container of task assignments for some solution

The general flow of what happens in my program is:
	1. Read all tasks and populate global Tasks container
	2. Assign tasks to workstations using SST heuristic
		a. Fill lowest Id workstations first
	3. Store assignments in a Solution object
	4. Evaluate finished Solution for qualities described in write-up
	5. Print analysis and found solution data

One problem was that starting indices for tasks did not have to be 0. Because of this I decided to maintain a mapping between the index read from input (TaskId) and the index of that task in my global container of tasks (an array or slice). This was actually not as much trouble as I thought it might be, the ease of using the mapping is a direct result of creating an Id Type that could not implicitely used as an int.

Disclaimer:

This is the first program I've written in Go beyond a "Hello World"-like program. I spent a great deal of time learning how to program in Go. Although a lot of time was spent learning Go I've found that I really like that it comes as a feature-rich environment, not just a language. So if anything looks really weird it may be because I'm really new to Go.

I also started on an implementation to find the best solution but I realized it would likely run forever given the complexity of the problem so I didn't finish it.
