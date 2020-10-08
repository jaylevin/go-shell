package manager

import "container/list"

const BLOCKED = 0

type Resource struct {
	state int // number of available units

	//	r.waitlist contains pairs (i, k) where i is the waiting process and k is the number of requested units
	waitList []*ProcessWaitingTuple
}

type Process struct {
	priority int
	state    int
	parent   int

	resources []*ResourceOwnedTuple
	children  *list.List
}

// public class Manager {
//     static final int MAX_PROCESS = 16;
//     static final int RL_LEVELS = 3;
//     Process[] PCB;
//     LinkedList<Integer>[] RL;
// ​
//     public Manager(){
//         PCB = new Process[MAX_PROCESS] ;
//         RL = new LinkedList[RL_LEVELS];
// 	}
//}
const MAX_PROCESS = 16
const RL_LEVELS = 3
const MAX_RESOURCES = 4

type Manager struct {
	rcb [MAX_RESOURCES]*Resource
	pcb [MAX_PROCESS]*Process
	rl  [RL_LEVELS]*list.List
}

// Tuples
type ProcessWaitingTuple struct {
	proc         int // the process that is waiting
	unitsWaiting int // number of units the process is waiting on
}

type ResourceOwnedTuple struct {
	res       int // index of resource that is mapped to RCB list
	unitsHeld int // number of units occupied
}
