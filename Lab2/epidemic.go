package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// create heartbeat table struct containing neighbor id, heartbeat count, and time
type Heartbeat struct {
	NeighborID int
	hbCounter  int
	Time       int64
	Fail       bool
}

// create node struct
type Node struct {
	ID             int
	hbCounter      int
	Neighbors      []int
	HeartbeatTable []Heartbeat
}

// function to increment HeartbeatCount every 20 ms
func (node *Node) incrementHeartbeatCount() {
	fmt.Println("Starting heartbeat count for node", node.ID)
	for {
		time.Sleep(20 * time.Millisecond)
		node.hbCounter++
		node.HeartbeatTable[node.ID].Time = time.Now().Unix()
		//fmt.Println("Heartbeat Count for node", node.ID, ": ", node.hbCounter)
	}
}

// function to simulate node failing
func killNode(failChannels []chan int) {
	for {
		time.Sleep(10 * time.Second)
		failChannels[rand.Intn(8)] <- 1
	}
}

// function to update heartbeat table every 200 ms
func (node *Node) updateHeartbeatTable(recChannel1 chan []Heartbeat, recChannel2 chan []Heartbeat, sendChannel1 chan []Heartbeat, sendChannel2 chan []Heartbeat, failChannel chan int) {
	var tempNeighbor0 []Heartbeat
	var tempNeighbor1 []Heartbeat
	var failTime int64

	// starting current node heartbeat values
	node.HeartbeatTable[node.ID].NeighborID = node.ID
	node.HeartbeatTable[node.ID].hbCounter = node.hbCounter
	node.HeartbeatTable[node.ID].Time = time.Now().Unix()

	// sending initial data to neighbors
	sendChannel1 <- node.HeartbeatTable
	sendChannel2 <- node.HeartbeatTable

	deadFlag := 0
	failTime = 0
	// loop forever
	for {
		//i++
		// sleep for 200 ms
		time.Sleep(200 * time.Millisecond)

		// update current node's entry in heartbeat table
		node.HeartbeatTable[node.ID].NeighborID = node.ID
		node.HeartbeatTable[node.ID].hbCounter = node.hbCounter
		node.HeartbeatTable[node.ID].Time = time.Now().Unix()

		// simulate node failing
		select {
		case deadFlag = <-failChannel:
			// close channel to stop sending data to neighbors
			fmt.Println("Node", node.ID, "failed")
			close(sendChannel1)
			close(sendChannel2)
		default:
			// do nothing
		}
		// break out of goroutine if node is dead
		if deadFlag == 1 {
			break
		}

		// check if channel has new data
		select {
		case tempNeighbor0 = <-recChannel1:
			select {
			case sendChannel1 <- node.HeartbeatTable:
				// message sent
				// check to see if heartbeat is greater than current time for each node in heartbeat table
				for i := 0; i < len(node.HeartbeatTable); i++ {
					// don't replace dead node with bad info
					if node.HeartbeatTable[i].hbCounter == 0 && node.HeartbeatTable[i].Fail == true {
						// do nothing
					} else if (tempNeighbor0[i].hbCounter > node.HeartbeatTable[i].hbCounter && tempNeighbor0[i].Fail == false) || tempNeighbor0[i].Fail == true {
						node.HeartbeatTable[i] = tempNeighbor0[i]
					}
				}
			default:
				// do nothing
			}

		default:
			// do nothing
		}

		select {
		case tempNeighbor1 = <-recChannel2:
			select {
			case sendChannel2 <- node.HeartbeatTable:
				// message sent

				for i := 0; i < len(node.HeartbeatTable); i++ {
					if node.HeartbeatTable[i].hbCounter == 0 && node.HeartbeatTable[i].Fail == true {
						// do nothing
					} else if (tempNeighbor1[i].hbCounter > node.HeartbeatTable[i].hbCounter && tempNeighbor1[i].Fail == false) || tempNeighbor1[i].Fail == true {
						node.HeartbeatTable[i] = tempNeighbor1[i]
					}
				}
			default:
				// do nothing
			}
		default:
			// do nothing
		}

		// if neighbor hasn't sent heartbeat for more than 3 seconds, flag it for removal from heartbeat table
		for i := 0; i < len(node.Neighbors); i++ {
			if time.Now().Unix()-node.HeartbeatTable[node.Neighbors[i]].Time > 3 && node.HeartbeatTable[node.Neighbors[i]].Time != 0 && node.HeartbeatTable[node.Neighbors[i]].Fail == false {
				//fmt.Println("Node", node.ID, "removed neighbor", node.Neighbors[i])
				node.HeartbeatTable[node.Neighbors[i]].Fail = true
				failTime = time.Now().Unix()
				// if neighbor has had fail flag for more than 3 seconds, remove it from heartbeat table
			} else if time.Now().Unix()-failTime > 1 && node.HeartbeatTable[node.Neighbors[i]].Fail == true && failTime != 0 {
				node.HeartbeatTable[node.Neighbors[i]] = Heartbeat{}
				node.HeartbeatTable[node.Neighbors[i]].NeighborID = node.Neighbors[i]
				node.HeartbeatTable[node.Neighbors[i]].Fail = true
				failTime = 0
			}
		}
		fmt.Println("heartbeat table for node", node.ID, ":", node.HeartbeatTable)

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

	// create 16 channels to send/receive heartbeat messages
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

	// channels for simulating a node failing
	failChannels := []chan int{
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
	}

	// 8 go routines to incrememnt heartbeat count
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go nodes[i].incrementHeartbeatCount()
		//go updateHeartbeatTable(&nodes[i], channels[i-1], channels[i+1])
	}

	go nodes[0].updateHeartbeatTable(rightChannels[7], leftChannels[0], leftChannels[7], rightChannels[0], failChannels[0])
	go nodes[1].updateHeartbeatTable(rightChannels[0], leftChannels[1], leftChannels[0], rightChannels[1], failChannels[1])
	go nodes[2].updateHeartbeatTable(rightChannels[1], leftChannels[2], leftChannels[1], rightChannels[2], failChannels[2])
	go nodes[3].updateHeartbeatTable(rightChannels[2], leftChannels[3], leftChannels[2], rightChannels[3], failChannels[3])
	go nodes[4].updateHeartbeatTable(rightChannels[3], leftChannels[4], leftChannels[3], rightChannels[4], failChannels[4])
	go nodes[5].updateHeartbeatTable(rightChannels[4], leftChannels[5], leftChannels[4], rightChannels[5], failChannels[5])
	go nodes[6].updateHeartbeatTable(rightChannels[5], leftChannels[6], leftChannels[5], rightChannels[6], failChannels[6])
	go nodes[7].updateHeartbeatTable(rightChannels[6], leftChannels[7], leftChannels[6], rightChannels[7], failChannels[7])

	go killNode(failChannels)

	// run forever
	wg.Wait()

}
