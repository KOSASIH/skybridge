package node

import (
	"fmt"
	"log"
	"sync"

	"github.com/KOSASIH/skybridge/utils"
)

// Node represents a node in a graph
type Node struct {
	id        string
	neighbors []*Node
}

// NewNode returns a new Node instance
func NewNode(id string) *Node {
	return &Node{
		id:        id,
		neighbors: make([]*Node, 0),
	}
}

// AddNeighbor adds a neighbor to the node
func (n *Node) AddNeighbor(neighbor *Node) {
	n.neighbors = append(n.neighbors, neighbor)
}

// GetNeighbors returns the node's neighbors
func (n *Node) GetNeighbors() []*Node {
	return n.neighbors
}

// Traverse performs a depth-first traversal of the graph starting from the node
func (n *Node) Traverse(visited map[string]bool) {
	if visited[n.id] {
		return
	}

	visited[n.id] = true
	log.Printf("Visiting node %s", n.id)

	for _, neighbor := range n.neighbors {
		neighbor.Traverse(visited)
	}
}

// TraverseConcurrently performs a concurrent depth-first traversal of the graph starting from the node
func (n *Node) TraverseConcurrently(visited map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if visited[n.id] {
		return
	}

	visited[n.id] = true
	log.Printf("Visiting node %s", n.id)

	var neighborWgs sync.WaitGroup
	for _, neighbor := range n.neighbors {
		neighborWgs.Add(1)
		go func(neighbor *Node) {
			neighbor.TraverseConcurrently(visited, &neighborWgs)
		}(neighbor)
	}

	neighborWgs.Wait()
}
