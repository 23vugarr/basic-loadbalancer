package main

import (
	"loadbalancer/types"
)

func main() {
	node1 := types.NewNode(1, 1)
	node2 := types.NewNode(2, 2)
	node3 := types.NewNode(3, 3)
	node4 := types.NewNode(4, 4)
	node5 := types.NewNode(5, 5)
	node6 := types.NewNode(6, 6)

	node1.AddNeigbhour([]*types.Node{node2, node3, node4, node5, node6})
	node2.AddNeigbhour([]*types.Node{node1, node3, node4, node5, node6})
	node3.AddNeigbhour([]*types.Node{node2, node1, node4, node5, node6})
	node4.AddNeigbhour([]*types.Node{node2, node3, node1, node5, node6})
	node5.AddNeigbhour([]*types.Node{node1, node2, node3, node4, node6})
	node6.AddNeigbhour([]*types.Node{node1, node2, node3, node4, node5})

	loadbalancer := types.NewLoadBalancer([]*types.Node{node3, node5, node1, node6, node2, node4}, []string{"task1", "task2", "task6", "task5", "task3", "task4"})
	loadbalancer.InitBalancer()
	loadbalancer.GetResult()
}
