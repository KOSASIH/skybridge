package unit

import (
	"testing"
	"github.com/your-username/satellite"
	"github.com/your-username/node"
)

func TestSatelliteUnit(t *testing.T) {
	// Create a satellite instance
	satellite := satellite.NewSatellite("sat1")

	// Verify that the satellite's ID is set correctly
	if satellite.id != "sat1" {
		t.Errorf("Satellite ID is not set correctly")
	}

	// Verify that the satellite's ground stations are empty
	if len(satellite.groundStations) != 0 {
		t.Errorf("Satellite's ground stations are not empty")
	}

	// Verify that the satellite's transceivers are empty
	if len(satellite.transceivers) != 0 {
		t.Errorf("Satellite's transceivers are not empty")
	}
}

func TestNodeUnit(t *testing.T) {
	// Create a node instance
	node := node.NewNode("node1")

	// Verify that the node's ID is set correctly
	if node.id != "node1" {
		t.Errorf("Node ID is not set correctly")
	}

	// Verify that the node's neighbors are empty
	if len(node.neighbors) != 0 {
		t.Errorf("Node's neighbors are not empty")
	}
}
