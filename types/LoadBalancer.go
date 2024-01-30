package types

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

type LoadBalancer struct {
	Nodes     []*Node
	TaskLoads []string
	cycle     int
}

func NewLoadBalancer(nodes []*Node, tasks []string) *LoadBalancer {
	return &LoadBalancer{
		Nodes:     nodes,
		TaskLoads: tasks,
		cycle:     0,
	}
}

func (balancer *LoadBalancer) InitBalancer() {
	var wg sync.WaitGroup
	wg.Add(len(balancer.Nodes))
	balancer.Prioritize()

	for _, node := range balancer.Nodes {
		go balancer.RandomAssigner(node, &wg)
	}
	wg.Wait()

	balancer.AsyncBacktracking()
}

func (balancer *LoadBalancer) AsyncBacktracking() (bool, error) {
	if balancer.cycle >= len(balancer.TaskLoads) {
		return false, errors.New("maximum cycles reached")
	}

	solution := balancer.CheckSolution()

	if solution {
		return true, nil
	}
	
	balancer.GetResult()
	fmt.Println("------------------------------------------------")

	balancer.StartMessaging()
	balancer.cycle++

	return balancer.AsyncBacktracking()

}

func (balancer *LoadBalancer) StartMessaging() {
	var wg sync.WaitGroup
	wg.Add(len(balancer.Nodes))
	for _, node := range balancer.Nodes {
		go node.NotifyNeigbhours(balancer.TaskLoads, &wg)
	}
	wg.Wait()
}

func (balancer *LoadBalancer) CheckSolution() bool {
	var wg sync.WaitGroup
	wg.Add(len(balancer.Nodes))
	result := make(chan bool, len(balancer.Nodes))

	for _, node := range balancer.Nodes {
		go node.CheckNeigbhours(result, &wg)
	}

	wg.Wait()
	close(result)

	for res := range result {
		if !res {
			return false
		}
	}

	return true
}

func (balancer *LoadBalancer) RandomAssigner(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()
	randomIndex := rand.Intn(len(balancer.TaskLoads))
	node.CurrentLoad = balancer.TaskLoads[randomIndex]
}

func (balancer *LoadBalancer) Prioritize() {
	sort.Slice(balancer.Nodes, func(i, j int) bool {
		return balancer.Nodes[i].Priority < balancer.Nodes[j].Priority
	})
}

func (balancer *LoadBalancer) GetResult() {
	for _, value := range balancer.Nodes {
		fmt.Println("nodeID", value.NodeID, "current load", value.CurrentLoad)
	}
	fmt.Println("SOLUTION DONE: ", balancer.CheckSolution())
}
