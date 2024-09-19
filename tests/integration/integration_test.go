package integration

import (
	"testing"
	"github.com/your-username/satellite"
	"github.com/your-username/node"
)

func TestIntegration(t *testing.T) {
	// Create a satellite instance
	satellite := satellite.NewSatellite("sat1")

	// Create a node instance
	node := node.NewNode("node1")

	// Add the node as a ground station to the satellite
	satellite.AddGroundStation(node)

	// Start the satellite's orbit
	go satellite.StartOrbit()

	// Wait for the satellite to communicate with the node
	time.Sleep(30 * time.Second)

	// Verify that the node received data from the satellite
	if len(node.GetNeighbors()) == 0 {
		t.Errorf("Node did not receive data from satellite")
	}
}

func TestConcurrentIntegration(t *testing.T) {
	// Create multiple satellite instances
	satellites := []*satellite.Satellite{
		satellite.NewSatellite("sat1"),
		satellite.NewSatellite("sat2"),
		satellite.NewSatellite("sat3"),
	}

	// Create multiple node instances
	nodes := []*node.Node{
		node.NewNode("node1"),
		node.NewNode("node2"),
		node.NewNode("node3"),
	}

	// Add nodes as ground stations to satellites
	for i, satellite := range satellites {
		satellite.AddGroundStation(nodes[i])
	}

	// Start the satellites' orbits concurrently
	var wg sync.WaitGroup
	for _, satellite := range satellites {
		wg.Add(1)
		go func(satellite *satellite.Satellite) {
			defer wg.Done()
			satellite.StartOrbit()
		}(satellite)
	}

	wg.Wait()

	// Wait for the satellites to communicate with the nodes
	time.Sleep(30 * time.Second)

	// Verify that the nodes received data from the satellites
	for _, node := range nodes {
		if len(node.GetNeighbors()) == 0 {
			t.Errorf("Node did not receive data from satellite")
		}
	}
}
