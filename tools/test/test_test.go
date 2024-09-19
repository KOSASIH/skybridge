package test

import (
	"testing"
	"github.com/your-username/satellite"
	"github.com/your-username/node"
)

func TestTester(t *testing.T) {
	satelliteConfig := &satellite.Satellite{
		id: "sat1",
		groundStations: []*satellite.GroundStation{
			&satellite.GroundStation{id: "gs1"},
			&satellite.GroundStation{id: "gs2"},
		},
		transceivers: []*satellite.Transceiver{
			&satellite.Transceiver{id: "tx1"},
			&satellite.Transceiver{id: "tx2"},
		},
	}

	nodeConfig := &node.Node{
		id: "node1",
		neighbors: []*node.Node{
			&node.Node{id: "node2"},
			&node.Node{id: "node3"},
		},
	}

	tester := NewTester(satelliteConfig, nodeConfig)

	err := tester.Test()
	if err != nil {
		t.Errorf("Error testing configurations: %v", err)
	}
}
