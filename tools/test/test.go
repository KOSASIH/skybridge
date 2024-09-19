package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/your-username/satellite"
	"github.com/your-username/node"
)

// Tester represents a tester for satellite and node configurations
type Tester struct {
	satelliteConfig *satellite.Satellite
	nodeConfig     *node.Node
}

// NewTester returns a new Tester instance
func NewTester(satelliteConfig *satellite.Satellite, nodeConfig *node.Node) *Tester {
	return &Tester{
		satelliteConfig: satelliteConfig,
		nodeConfig:     nodeConfig,
	}
}

// TestSatelliteConfig tests the satellite configuration
func (t *Tester) TestSatelliteConfig() error {
	log.Println("Testing satellite configuration")
	// Perform tests on the satellite configuration
	if t.satelliteConfig.id == "" {
		return fmt.Errorf("satellite ID is empty")
	}

	for _, groundStation := range t.satelliteConfig.groundStations {
		if groundStation.id == "" {
			return fmt.Errorf("ground station ID is empty")
		}
	}

	for _, transceiver := range t.satelliteConfig.transceivers {
		if transceiver.id == "" {
			return fmt.Errorf("transceiver ID is empty")
		}
	}

	return nil
}

// TestNodeConfig tests the node configuration
func (t *Tester) TestNodeConfig() error {
	log.Println("Testing node configuration")
	// Perform tests on the node configuration
	if t.nodeConfig.id == "" {
		return fmt.Errorf("node ID is empty")
	}

	for _, neighbor := range t.nodeConfig.neighbors {
		if neighbor.id == "" {
			return fmt.Errorf("neighbor ID is empty")
		}
	}

	return nil
}

// Test tests both the satellite and node configurations
func (t *Tester) Test() error {
	err := t.TestSatelliteConfig()
	if err != nil {
		return err
	}

	err = t.TestNodeConfig()
	if err != nil {
		return err
	}

	return nil
}
