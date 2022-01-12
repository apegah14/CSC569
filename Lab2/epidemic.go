package main

import (
	"fmt"
	"sync"
	"time"
)

// create heartbeat table struct containing neighbor id, heartbeat count, and time
type Heartbeat struct {
	NeighborID int
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

// function to update heartbeat table every 200 ms
func updateHeartbeatTable(node *Node, recChannel1 chan []Heartbeat, recChannel2 chan []Heartbeat, sendChannel1 chan []Heartbeat, sendChannel2 chan []Heartbeat) {
	// loop forever
	for {

		// sleep for 200 ms
		time.Sleep(200 * time.Millisecond)

		// update current node's entry in heartbeat table
		node.HeartbeatTable[node.ID].NeighborID = node.ID
		node.HeartbeatTable[node.ID].hbCounter = node.hbCounter
		node.HeartbeatTable[node.ID].Time = time.Now().Unix()

		//fmt.Println("current node heartbeat table: ", node.HeartbeatTable)

		// send current heartbeat table to send channels
		sendChannel1 <- node.HeartbeatTable
		sendChannel2 <- node.HeartbeatTable

		// update node's neighbors' entries in heartbeat table using channels from other go routines

		// grab neighbor's heartbeat table from receive channels
		tempNeighbor0 := <-recChannel1
		tempNeighbor1 := <-recChannel2

		//fmt.Println("tempNeighbor0 for node", node.ID, ": ", tempNeighbor0)
		//fmt.Println("tempNeighbor1 for node", node.ID, ": ", tempNeighbor1)

		// check to see if heartbeat is greater than current time for each node in heartbeat table
		for i := 0; i < len(tempNeighbor0); i++ {
			if tempNeighbor0[i].hbCounter > node.HeartbeatTable[i].hbCounter {
				node.HeartbeatTable[i] = tempNeighbor0[i]
			}
		}

		for i := 0; i < len(tempNeighbor1); i++ {
			if tempNeighbor1[i].hbCounter > node.HeartbeatTable[i].hbCounter {
				node.HeartbeatTable[i] = tempNeighbor1[i]
			}
		}

		// if node in heartbeat table hasn't been updated in 10 seconds, remove it
		for i := 0; i < len(node.HeartbeatTable); i++ {
			if time.Now().Unix()-node.HeartbeatTable[i].Time > 10 && node.HeartbeatTable[i].Time != 0 {
				node.HeartbeatTable[i] = node.HeartbeatTable[len(node.HeartbeatTable)-1]
				node.HeartbeatTable = node.HeartbeatTable[:len(node.HeartbeatTable)-1]
			}
		}

		fmt.Println("current node heartbeat table for node ", node.ID, ": ", node.HeartbeatTable)

	}

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
	// counter clockwise channels in ring
	leftChannels := []chan []Heartbeat{
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
		make(chan []Heartbeat, 2),
	}

	// clockwise channels in ring
	rightChannels := []chan []Heartbeat{
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
		go incrementHeartbeatCount(&nodes[i])
		//go updateHeartbeatTable(&nodes[i], channels[i-1], channels[i+1])
	}
	go updateHeartbeatTable(&nodes[0], rightChannels[7], leftChannels[0], leftChannels[7], rightChannels[0])
	go updateHeartbeatTable(&nodes[1], rightChannels[0], leftChannels[1], leftChannels[0], rightChannels[1])
	go updateHeartbeatTable(&nodes[2], rightChannels[1], leftChannels[2], leftChannels[1], rightChannels[2])
	go updateHeartbeatTable(&nodes[3], rightChannels[2], leftChannels[3], leftChannels[2], rightChannels[3])
	go updateHeartbeatTable(&nodes[4], rightChannels[3], leftChannels[4], leftChannels[3], rightChannels[4])
	go updateHeartbeatTable(&nodes[5], rightChannels[4], leftChannels[5], leftChannels[4], rightChannels[5])
	go updateHeartbeatTable(&nodes[6], rightChannels[5], leftChannels[6], leftChannels[5], rightChannels[6])
	go updateHeartbeatTable(&nodes[7], rightChannels[6], leftChannels[7], leftChannels[6], rightChannels[7])

	//time.Sleep(2 * time.Second)
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
	wg.Wait()

	// create 8 go routines to send heartbeat messages

}
