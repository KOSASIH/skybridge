package node

import (
	"testing"
)

func TestNode(t *testing.T) {
	node1 := NewNode("1")
	node2 := NewNode("2")
	node3 := NewNode("3")
	node4 := NewNode("4")

	node1.AddNeighbor(node2)
	node1.AddNeighbor(node3)
	node2.AddNeighbor(node4)

	visited := make(map[string]bool)
	node1.Traverse(visited)

	if len(visited) != 4 {
		t.Errorf("Expected 4 nodes to be visited, got %d", len(visited))
	}

	for _, id := range []string{"1", "2", "3", "4"} {
		if !visited[id] {
			t.Errorf("Node %s was not visited", id)
		}
	}

	var wg sync.WaitGroup
	visited = make(map[string]bool)
	wg.Add(1)
	go func() {
		node1.TraverseConcurrently(visited, &wg)
	}()

	wg.Wait()

	if len(visited) != 4 {
		t.Errorf("Expected 4 nodes to be visited, got %d", len(visited))
	}

	for _, id := range []string{"1", "2", "3", "4"} {
		if !visited[id] {
			t.Errorf("Node %s was not visited", id)
		}
	}
}
