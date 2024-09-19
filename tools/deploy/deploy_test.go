package deploy

import (
	"testing"
	"github.com/your-username/satellite"
	"github.com/your-username/node"
)

func TestDeployer(t *testing.T) {
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

	deployer := NewDeployer(satelliteConfig, nodeConfig)

	err := deployer.Deploy()
	if err != nil {
		t.Errorf("Error deploying configurations: %v", err)
	}
}
