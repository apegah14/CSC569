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

// function to fail increment HeartbeatCount every 20 ms
func (node *Node) incrementHeartbeatCountFail() {
	fmt.Println("Starting heartbeat count for node", node.ID)
	for {
		time.Sleep(20 * time.Millisecond)
		node.hbCounter++

		//fmt.Println("Heartbeat Count for node", node.ID, ": ", node.hbCounter)
	}
}

// function to update heartbeat table every 200 ms
func (node *Node) updateHeartbeatTable(recChannel1 chan []Heartbeat, recChannel2 chan []Heartbeat, sendChannel1 chan []Heartbeat, sendChannel2 chan []Heartbeat) {

	//failedNeighbor := 3
	//deleteFlag := 0
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

	i := 0
	deadFlag := 0
	failTime = 0
	// loop forever
	for {
		i++
		// sleep for 200 ms
		time.Sleep(200 * time.Millisecond)

		// update current node's entry in heartbeat table
		node.HeartbeatTable[node.ID].NeighborID = node.ID
		node.HeartbeatTable[node.ID].hbCounter = node.hbCounter
		node.HeartbeatTable[node.ID].Time = time.Now().Unix()

		//fmt.Println("current node heartbeat table: ", node.HeartbeatTable)
		// simulate node failing
		if node.ID == 3 {
			switch i {
			case 10:
				// close channel to stop sending data to neighbors
				fmt.Println("Node", node.ID, "failed")
				close(sendChannel1)
				close(sendChannel2)
				deadFlag = 1
			default:
				// do nothing
			}
		}

		if node.ID == 4 {
			switch i {
			case 100:
				// close channel to stop sending data to neighbors
				fmt.Println("Node", node.ID, "failed")
				close(sendChannel1)
				close(sendChannel2)
				deadFlag = 1
			default:
				// do nothing
			}
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

				for i := 0; i < len(tempNeighbor1); i++ {
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
				fmt.Println(failTime)
				//deleteFlag = 1
				// if neighbor has had fail flag for more than 3 seconds, remove it from heartbeat table
			} else if time.Now().Unix()-failTime > 1 && node.HeartbeatTable[node.Neighbors[i]].Fail == true && failTime != 0 {
				fmt.Println(time.Now().Unix())
				node.HeartbeatTable[node.Neighbors[i]] = Heartbeat{}
				node.HeartbeatTable[node.Neighbors[i]].NeighborID = node.Neighbors[i]
				node.HeartbeatTable[node.Neighbors[i]].Fail = true
				failTime = 0
			}
		}

		/*
			// send current heartbeat table to send channels
			// checks to see if node has failed
			if failedNeighbor == 1 {
				sendChannel1 <- node.HeartbeatTable
			} else if failedNeighbor == 0 {
				sendChannel2 <- node.HeartbeatTable
			} else if failedNeighbor == 2 {
				// do nothing if both neighbors have failed
			} else {
				sendChannel1 <- node.HeartbeatTable
				sendChannel2 <- node.HeartbeatTable
			}

			// grab neighbor's heartbeat table from receive channels
			// checks to see if node has failed
			if failedNeighbor == 1 {
				tempNeighbor0 = <-recChannel1
			} else if failedNeighbor == 0 {
				tempNeighbor1 = <-recChannel2
			} else if failedNeighbor == 2 {
				// do nothing if both neighbors have failed
			} else {
				tempNeighbor0 = <-recChannel1
				tempNeighbor0 = <-recChannel2
			}
		*/
		//fmt.Println("tempNeighbor0 for node", node.ID, ": ", tempNeighbor0)
		//fmt.Println("tempNeighbor1 for node", node.ID, ": ", tempNeighbor1)
		/*
			// checks to see if either of the neighbor's heartbeat counts are the same as last cycle
			if deleteFlag == 0 && len(tempNeighbor0) != 0 {
				if node.HeartbeatTable[node.Neighbors[0]].hbCounter == tempNeighbor0[node.Neighbors[0]].hbCounter || node.HeartbeatTable[node.Neighbors[1]].hbCounter == tempNeighbor0[node.Neighbors[1]].hbCounter {

					// if neighbor's heartbeat hasn't changed in 10 seconds, remove from heartbeat table
					for i := 0; i < len(node.Neighbors); i++ {
						// checks to see if a specific amount of time has passed
						if node.HeartbeatTable[node.Neighbors[i]].Time < time.Now().Unix()-10 && node.HeartbeatTable[node.Neighbors[i]].Time != 0 {
							deleteFlag = 1
							// case if both neighbors have failed
							if failedNeighbor == 0 || failedNeighbor == 1 {
								failedNeighbor = 2
							} else {
								failedNeighbor = i
							}
						}
					}
				}
			} else if deleteFlag == 1 {
				// remove neighbor from heartbeat table
				node.HeartbeatTable[node.Neighbors[failedNeighbor]] = Heartbeat{}
				deleteFlag = 0
			}

			// close channel to failed neighbor
			if failedNeighbor == 0 {
				close(recChannel1)
			} else if failedNeighbor == 1 {
				close(recChannel2)
			}
		*/
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
	for i := 0; i < 7; i++ {
		wg.Add(1)
		go nodes[i].incrementHeartbeatCount()
		//go updateHeartbeatTable(&nodes[i], channels[i-1], channels[i+1])
	}
	go nodes[7].incrementHeartbeatCountFail()
	go nodes[0].updateHeartbeatTable(rightChannels[7], leftChannels[0], leftChannels[7], rightChannels[0])
	go nodes[1].updateHeartbeatTable(rightChannels[0], leftChannels[1], leftChannels[0], rightChannels[1])
	go nodes[2].updateHeartbeatTable(rightChannels[1], leftChannels[2], leftChannels[1], rightChannels[2])
	go nodes[3].updateHeartbeatTable(rightChannels[2], leftChannels[3], leftChannels[2], rightChannels[3])
	go nodes[4].updateHeartbeatTable(rightChannels[3], leftChannels[4], leftChannels[3], rightChannels[4])
	go nodes[5].updateHeartbeatTable(rightChannels[4], leftChannels[5], leftChannels[4], rightChannels[5])
	go nodes[6].updateHeartbeatTable(rightChannels[5], leftChannels[6], leftChannels[5], rightChannels[6])
	go nodes[7].updateHeartbeatTable(rightChannels[6], leftChannels[7], leftChannels[6], rightChannels[7])

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
