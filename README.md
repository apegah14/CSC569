# CSC 569 Lab 2: Gossip Heartbeat Protocol

## Implementation Details
The following code is an implementation of a gossip heartbeat protocol in GoLang. Each node has its own set ID, two neighbors, and a local heartbeat table. The heartbeat table holds the ID of the neighbor, the heartbeat counter, time, and a flag to mark if the node is known failed. Each node is a seperate Goroutine and the communication between nodes are channels. Each node has two receive and two send channels.

This distributed system has fault detection. If a node has detected its adjacent neighbor has failed, usually in the form the heartbeat count remaining constant for a prolonged period of time, it will mark that node as failed. If the node has been marked as failed for a certain amount of time it will clear that entry in the heartbeat table and gossip the info to other nodes. If additional nodes were to be added into the system they would take the place of the failed nodes first (not implemented here)