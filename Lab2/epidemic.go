package main

import (
	"fmt"
	"sync"
	"time"
)

// create heartbeat table struct containing neighbor id, heartbeat count, and time
type Heartbeat struct {
	NeighborID []int
	hbCounter  int
	Time       int64
}

// create node struct
type Node struct {
	ID             int
	hbCounter      int
	Neighbors      []int
	HeartbeatTable []Heartbeat
}

// function to increment HeartbeatCount every 20 ms
func incrementHeartbeatCount(node *Node) {
	fmt.Println("Starting heartbeat count for node", node.ID)
	for {
		time.Sleep(20 * time.Millisecond)
		node.hbCounter++
		//fmt.Println("Heartbeat Count for node", node.ID, ": ", node.hbCounter)
	}
}

// function to update heartbeat table every 100 ms
func updateHeartbeatTable(node *Node, channel1 chan []Heartbeat, channel2 chan []Heartbeat) {
	//var tempNeighbor0 []Heartbeat
	//var tempNeighbor1 []Heartbeat

	// update node's entry in heartbeat table
	node.HeartbeatTable[node.ID].NeighborID = node.Neighbors
	node.HeartbeatTable[node.ID].hbCounter = node.hbCounter
	node.HeartbeatTable[node.ID].Time = time.Now().Unix()

	//fmt.Println("current node heartbeat table: ", node.HeartbeatTable)

	channel1 <- node.HeartbeatTable
	//channel2 <- node.HeartbeatTable

	fmt.Println("did this finish")

	// update node's neighbors' entries in heartbeat table using channels from other go routines

	//tempNeighbor0 = <-channel1
	tempNeighbor1 := <-channel2

	//fmt.Println("tempNeighbor0: ", tempNeighbor0)
	fmt.Println("tempNeighbor1: ", tempNeighbor1)

	// check if tempNeighbor0 is

	//node.HeartbeatTable[neighbor] := receiveChannels(node.ID - 1)

	/*
		for {
			time.Sleep(100 * time.Millisecond)
			for i := 0; i < len(node.HeartbeatTable); i++ {
				if node.HeartbeatTable[i].NeighborID == node.ID {
					node.HeartbeatTable[i].hbCounter = node.hbCounter
					node.HeartbeatTable[i].Time = time.Now().Unix()
				}
			}
		}
	*/
}

func main() {
	// shit to just make sure it runs forever
	var wg sync.WaitGroup

	// create 8 nodes with id 0-7 and neighbors id - 1 and id + 1
	nodes := []Node{
		Node{
			ID:             0,
			hbCounter:      0,
			Neighbors:      []int{7, 1},
			HeartbeatTable: make([]Heartbeat, 8),
		},
		Node{
			ID:             1,
			hbCounter:      0,
			Neighbors:      []int{0, 2},
			HeartbeatTable: make([]Heartbeat, 8),
		},
		Node{
			ID:             2,
			hbCounter:      0,
			Neighbors:      []int{1, 3},
			HeartbeatTable: make([]Heartbeat, 8),
		},
		Node{
			ID:             3,
			hbCounter:      0,
			Neighbors:      []int{2, 4},
			HeartbeatTable: make([]Heartbeat, 8),
		},
		Node{
			ID:             4,
			hbCounter:      0,
			Neighbors:      []int{3, 5},
			HeartbeatTable: make([]Heartbeat, 8),
		},
		Node{
			ID:             5,
			hbCounter:      0,
			Neighbors:      []int{4, 6},
			HeartbeatTable: make([]Heartbeat, 8),
		},
		Node{
			ID:             6,
			hbCounter:      0,
			Neighbors:      []int{5, 7},
			HeartbeatTable: make([]Heartbeat, 8),
		},
		Node{
			ID:             7,
			hbCounter:      0,
			Neighbors:      []int{6, 0},
			HeartbeatTable: make([]Heartbeat, 8),
		},
	}

	// create 8 channels to send heartbeat messages
	channels := []chan []Heartbeat{
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
	}

	// 8 go routines to incrememnt heartbeat count
	for i := 0; i < 8; i++ {
		wg.Add(1)
		//go incrementHeartbeatCount(&nodes[i])
		//go updateHeartbeatTable(&nodes[i], channels[i-1], channels[i+1])
	}
	go updateHeartbeatTable(&nodes[0], channels[7], channels[0])
	go updateHeartbeatTable(&nodes[1], channels[0], channels[1])
	go updateHeartbeatTable(&nodes[2], channels[1], channels[2])
	go updateHeartbeatTable(&nodes[3], channels[2], channels[3])
	go updateHeartbeatTable(&nodes[4], channels[3], channels[4])
	go updateHeartbeatTable(&nodes[5], channels[4], channels[5])
	go updateHeartbeatTable(&nodes[6], channels[5], channels[6])
	go updateHeartbeatTable(&nodes[7], channels[6], channels[7])

	time.Sleep(2 * time.Second)
	/*
		fmt.Println("channel 0: ", <-channels[0])
		fmt.Println("channel 1: ", <-channels[1])
		fmt.Println("channel 2: ", <-channels[2])
		fmt.Println("channel 3: ", <-channels[3])
		fmt.Println("channel 4: ", <-channels[4])
		fmt.Println("channel 5: ", <-channels[5])
		fmt.Println("channel 6: ", <-channels[6])
		fmt.Println("channel 7: ", <-channels[7])
	*/
	//wg.Wait()

	// create 8 go routines to send heartbeat messages

}
