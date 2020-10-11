package manager

import (
	"container/list"
	"errors"
	"fmt"

	"github.com/apex/log"
)

/*

//•	A process descriptor array PCB[16]
//•	A resource descriptor array RCB[4] with multiunit resources
//•	RCB[0] and RCB[1] have 1 unit each; RCB[2] has 2 units; RCB[3] has 3 units
//•	A 3-level RL
​
​
//For scheduling purposes, the manager maintains all PCBs on one of several lists.
// Blocked processes are kept on waiting lists associated with the relevant resource, as explained in Section 2.4.
// All ready processes are kept on a Ready List (RL). The basic version of the manager treats all processes as
// having the same priority and thus the RL is organized as a single linked list of PCB indices.
​
//When the system is initialized, process 0 is created automatically and it becomes the first running process.
// All other processes are created and destroyed dynamically using the following operations.
​
​
import java.util.LinkedList;
​
//Examples of conditions that must be detected and reported include:
    //Creating more than n processes
    //Destroying a process that is not a child of the current process
    //Requesting a nonexistent resource
    //Requesting a resource the process is already holding
    //Releasing a resource the process is not holding
    //Process 0 should be prevented from requesting any resource to avoid the possibility of a deadlock where no process is on the RL.
public class Manager {
    static final int MAX_PROCESS = 16;
    static final int RL_LEVELS = 3;
    Process[] PCB;
    LinkedList<Integer>[] RL;
​
    public Manager(){
        PCB = new Process[MAX_PROCESS] ;
        RL = new LinkedList[RL_LEVELS];
    }
​
​
​
    //The currently running process, i, can create a new child process, j, by invoking the create function:
    public void create(int priority){
        //allocate new PCB[j]
        //state = ready
        //insert j into children of i
        //parent = i
        //children = NULL
        //resources = NULL
        // plvel = priority
        //insert j into RL
        //display: “process j created”
    }
​
​
    //The currently running process, i, can destroy one of its children, j, including all of j’s descendants,
    // by invoking the destroy function. The reason for destroying the entire subtree staring with process j is
    // to avoid the creation of orphan processes, that is, processes without a parent. The process i can also
    // destroy itself and all of its descendants by invoking destroy(i).
    public void destroy() {
        // for all k in children of j destroy(k)
            //remove j from parent's list
            //remove j from RL or waiting list
            //release all resources of j
            //free PCB of j
            //display: “n processes destroyed”
​
    }
​
    //The currently running process, i, may request any of the resources, r, at any time using the request function:
    public void request(int resource){
        //if state of r is free
            //state of r = allocated
            //insert r into list of resources of process i
            //display: “resource r allocated”
        //else
            //state of i = blocked
            //move i from RL to waitlist of r
            //display: “process i blocked”
            //scheduler()
​
    }
​
    //The currently running process, i, may release any of the resources, r, it is holding using the release function:
    public void release(){
        //remove r from resources list of process i
        //if waitlist of r is empty
            //state of r = free
        //else
            //move process j from the head of waitlist of r to RL
            //state of j = ready
            //insert r into resources list of process j
        //display: “resource r released”
    }
​
    //The system mimics preemptive scheduling by a function timeout(), which moves the process, i,
    // currently at the head of the RL, to the end of the RL, and calls the scheduler to perform the context switch.
    public void timeout(){
        //move process i from head of RL to end of RL
        //scheduler()
    }
​
    //The task of the scheduler function is to perform the context switch from the currently running process i
    // to the new process j. The scheduler must be called whenever process i blocks on a resource and is removed
    // from the RL, and whenever the timeout function moves the process to the end of the RL.'
    //Starting with the highest priority level, 2, the scheduler finds the first non-empty list.
    // The head of the list is the highest priority ready process j.
    public void scheduler(){
        // find highest priority ready process j
        //display: “process i running”
    }
​
    //All PCB entries are initialized to free except PCB[0].
    //PCB[0] is initialized to be a running process with no parent, no children, and no resources.
    //All RCB entries are initialized to free.
    //RL contains process 0
    public void init() {
        //•	Erase all previous contents of the data structures PCB, RCB, RL
        //•	Create a single running process at PCB[0] with priority 0
        //•	Enter the process into the RL at the lowest-priority level 0
    }
​
}
*/

/*
    //The currently running process, i, may request any of the resources, r, at any time using the request function:
    public void request(int resource){
        //if state of r is free
            //state of r = allocated
            //insert r into list of resources of process i
            //display: “resource r allocated”
        //else
            //state of i = blocked
            //move i from RL to waitlist of r
            //display: “process i blocked”
            //scheduler()
​
    }

*/

func (m *Manager) removeProcessFromRL(process *Process) error {
	// Remove the process from the manager's RL
	p := m.rl[process.priority].Front()
	if p == nil {
		return fmt.Errorf("Could not find process in RL: %v", *process)
	}
	m.rl[process.priority].Remove(p)

	return nil
}

func (m *Manager) Request(resIndx int, unitsRequested int) {
	pIndx := m.getRunningProcessIndex()
	p := m.pcb[pIndx]
	res := m.rcb[resIndx]

	if res.state >= unitsRequested {
		res.state -= unitsRequested

		// insert (res, unitsRequested) into p.resources
		p.resources = append(p.resources, &ResourceOwnedTuple{
			resource:  resIndx,
			unitsHeld: unitsRequested,
		})
	} else {
		p.state = BLOCKED

		//    remove p from RL
		if err := m.removeProcessFromRL(p); err != nil {
			log.Debugf("Could not find process in RL: %s", err.Error())
		}

		//    insert (p, k) into res.waitlist
		res.waitList = append(res.waitList, &ProcessWaitingTuple{
			proc:         pIndx,
			unitsWaiting: unitsRequested,
		})
		m.scheduler()
	}
}

//All PCB entries are initialized to free except PCB[0].
//PCB[0] is initialized to be a running process with no parent, no children, and no resources.
//All RCB entries are initialized to free.
//RL contains process 0
//•	Erase all previous contents of the data structures PCB, RCB, RL
//•	Create a single running process at PCB[0] with priority 0
//•	Enter the process into the RL at the lowest-priority level 0
func (m *Manager) Init() {
	m.pcb = []*Process{
		{
			parent:    -1,
			resources: []*ResourceOwnedTuple{},
			children:  list.New(),
			priority:  0,
			state:     RUNNING,
		},
	}

	m.rcb = []*Resource{
		{state: 1},
		{state: 1},
		{state: 2},
		{state: 3},
	}

	m.rl = [RL_LEVELS]*list.List{list.New(), list.New(), list.New()}
	m.rl[0].PushFront(0)
}

func (m *Manager) GetReadyList() [RL_LEVELS]*list.List {
	return m.rl
}

//The task of the scheduler function is to perform the context switch from the currently running process i
// to the new process j. The scheduler must be called whenever process i blocks on a resource and is removed
// from the RL, and whenever the timeout function moves the process to the end of the RL.'
//Starting with the highest priority level, 2, the scheduler finds the first non-empty list.
// The head of the list is the highest priority ready process j.
// public void scheduler(){
//     // find highest priority ready process j
//     //display: “process i running”
// }
func (m *Manager) scheduler() {
	for _, list := range m.rl {
		if list.Front() != nil {
			fmt.Println("Scheduler found the first non empty ready list: " + list.Front().Value.(string))
		}
	}
}

func (m *Manager) release(process *Process, resourceIndex int) {
	// // remove (r, k) from i.resources
	// res.state += r.state + k
	// while (r.waitlist != empty && r.state > 0)
	//    get next (j, k) from r.waitlist
	//    if (r.state >= k)
	// 	  r.state = r.state - k
	// 	  insert (r, k) into j.resources
	// 	  j.state = ready
	// 	  remove (j, k) from r.waitlist
	// 	  insert j into RL
	//    else break

	// Remove (r, k) from process.resources
	resourceOwnership := process.resources[resourceIndex]
	resource := m.rcb[resourceIndex]

	// Resource.state += r.state + k
	resource.state += resourceOwnership.unitsHeld
	process.resources[resourceIndex] = nil

	for index, waiter := range resource.waitList {
		if resource.state >= waiter.unitsRequested {

			process := m.pcb[waiter.proc]

			// Insert (r, k) into j.resources
			process.resources = append(process.resources, &ResourceOwnedTuple{
				resource:  resourceIndex,
				unitsHeld: waiter.unitsRequested,
			})

			// j.state = ready
			process.state = READY

			// remove (j, k) from r.waitlist
			resource.waitList[index] = nil

			// insert j into RL
			m.rl[process.priority].PushBack(waiter.proc)

			// Subtract units available from resource
			resource.state -= waiter.unitsRequested
		}

	}

	m.scheduler()
}

// //The currently running process, i, can create a new child process, j, by invoking the create function:
// public void create(int priority){
//     //allocate new PCB[j]
//     //state = ready
//     //insert j into children of i
//     //parent = i
//     //children = NULL
//     //resources = NULL
//     // plvel = priority
//     //insert j into RL
//     //display: “process j created”
// }
func (m *Manager) Create(priority int) error {
	if priority < 0 || priority > 2 {
		return errors.New("Priority level must be in the range of [0, 2].")
	}

	i := m.getRunningProcessIndex()
	processI := m.pcb[i]

	// Insert j into children of process i
	j := len(m.pcb)
	processI.children.PushBack(j)

	// Allocate new PCB[j]
	m.pcb = append(m.pcb, &Process{
		priority:  priority,
		parent:    i,
		state:     READY,
		children:  list.New(),
		resources: []*ResourceOwnedTuple{},
	})

	// Insert j into RL
	m.rl[priority].PushBack(j)

	fmt.Println("Process j created")
	return nil
}

func (m *Manager) getRunningProcessIndex() int {
	for _, priorityLevel := range m.rl {
		fmt.Println("plevel: ", priorityLevel)
		elem := priorityLevel.Front()
		if elem != nil {
			return elem.Value.(int)
		}
	}

	return -1
}

//The currently running process, i, can destroy one of its children, j, including all of j’s descendants,
// by invoking the destroy function. The reason for destroying the entire subtree staring with process j is
// to avoid the creation of orphan processes, that is, processes without a parent. The process i can also
// destroy itself and all of its descendants by invoking destroy(i).
// public void destroy() {
// for all k in children of j destroy(k)
//remove j from parent's list
//remove j from RL or waiting list
//release all resources of j
//free PCB of j
//display: “n processes destroyed”
// }
func (m *Manager) Destroy(index int) (int, error) {
	if index > len(m.pcb)-1 {
		return 0, fmt.Errorf("PCB entry for process %v does not exist.", index)
	}

	numDeleted := 1
	j := m.pcb[index]

	for child := j.children.Front(); child != nil; child = child.Next() {
		childIndex := child.Value.(int)

		n, err := m.Destroy(childIndex)
		if err != nil {
			return 0, err
		}
		numDeleted += n

		if j.parent >= 0 {
			//remove j from parent's list
			parent := m.pcb[j.parent]
			for c := parent.children.Front(); c != nil; child = c.Next() {
				if c.Value.(int) == index {
					fmt.Printf("Removing %v from parent's children", index)
					parent.children.Remove(c)
				}
			}
		}
	}

	//remove j from RL or waiting list
	for _, level := range m.rl {
		for p := level.Front(); p != nil; p = p.Next() {
			i := p.Value.(int)
			if index == i {
				level.Remove(p)
			}
		}
	}

	for _, resource := range j.resources {

	}

	// free PCB of j
	m.pcb[index] = nil

	return numDeleted, nil
}
