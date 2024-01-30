package types

import "sync"

type Node struct {
	NodeID           int
	Priority         int
	CurrentLoad      string
	Neigbhours       []*Node
	UnavailableLoads []string
}

func NewNode(nodeid int, priority int) *Node {
	return &Node{
		NodeID:           nodeid,
		Priority:         priority,
		CurrentLoad:      "",
		Neigbhours:       nil,
		UnavailableLoads: nil,
	}
}

func (node *Node) AddNeigbhour(neig []*Node) {
	node.Neigbhours = append(node.Neigbhours, neig...)
}

func (node *Node) NotifyNeigbhours(tasks []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, neigh := range node.Neigbhours {
		if node.CurrentLoad == neigh.CurrentLoad && node.Priority < neigh.Priority {
			neigh.UnavailableLoads = append(neigh.UnavailableLoads, node.CurrentLoad)
			neigh.ChangeTask(tasks)
		}
	}
}

func (node *Node) CheckAvailability(task string) bool {
	for _, t := range node.UnavailableLoads {
		if task == t {
			return false
		}
	}
	return true
}

func (node *Node) ChangeTask(tasks []string) {
	for _, task := range tasks {
		if node.CheckAvailability(task) {
			node.CurrentLoad = task
			break
		} else {
			node.CurrentLoad = ""
		}
	}
}

func (node *Node) CheckNeigbhours(result chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, neigh := range node.Neigbhours {
		if node.CurrentLoad == neigh.CurrentLoad || node.CurrentLoad == "" {
			result <- false
			return
		}
	}
	result <- true
}
