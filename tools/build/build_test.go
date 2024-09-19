package build

import (
	"testing"
	"github.com/your-username/satellite"
	"github.com/your-username/node"
)

func TestBuilder(t *testing.T) {
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

	builder := NewBuilder(satelliteConfig, nodeConfig)

	err := builder.Build()
	if err != nil {
		t.Errorf("Error building configurations: %v", err)
	}
}
