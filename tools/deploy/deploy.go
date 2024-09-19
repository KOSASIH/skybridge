package deploy

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/your -username/satellite"
	"github.com/your-username/node"
)

// Deployer represents a deployer for satellite and node configurations
type Deployer struct {
	satelliteConfig *satellite.Satellite
	nodeConfig     *node.Node
}

// NewDeployer returns a new Deployer instance
func NewDeployer(satelliteConfig *satellite.Satellite, nodeConfig *node.Node) *Deployer {
	return &Deployer{
		satelliteConfig: satelliteConfig,
		nodeConfig:     nodeConfig,
	}
}

// DeploySatelliteConfig deploys the satellite configuration
func (d *Deployer) DeploySatelliteConfig() error {
	log.Println("Deploying satellite configuration")
	// Create a directory for the satellite configuration
	dir, err := os.MkdirTemp("", "satellite-config-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	// Create a file for the satellite configuration
	file, err := os.Create(filepath.Join(dir, "satellite.config"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the satellite configuration to the file
	_, err = file.WriteString(fmt.Sprintf("id=%s\n", d.satelliteConfig.id))
	if err != nil {
		return err
	}

	// Write the ground stations to the file
	for _, groundStation := range d.satelliteConfig.groundStations {
		_, err = file.WriteString(fmt.Sprintf("ground_station=%s\n", groundStation.id))
		if err != nil {
			return err
		}
	}

	// Write the transceivers to the file
	for _, transceiver := range d.satelliteConfig.transceivers {
		_, err = file.WriteString(fmt.Sprintf("transceiver=%s\n", transceiver.id))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeployNodeConfig deploys the node configuration
func (d *Deployer) DeployNodeConfig() error {
	log.Println("Deploying node configuration")
	// Create a directory for the node configuration
	dir, err := os.MkdirTemp("", "node-config-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	// Create a file for the node configuration
	file, err := os.Create(filepath.Join(dir, "node.config"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the node configuration to the file
	_, err = file.WriteString(fmt.Sprintf("id=%s\n", d.nodeConfig.id))
	if err != nil {
		return err
	}

	// Write the neighbors to the file
	for _, neighbor := range d.nodeConfig.neighbors {
		_, err = file.WriteString(fmt.Sprintf("neighbor=%s\n", neighbor.id))
		if err != nil {
			return err
		}
	}

	return nil
}

// Deploy deploys both the satellite and node configurations
func (d *Deployer) Deploy() error {
	err := d.DeploySatelliteConfig()
	if err != nil {
		return err
	}

	err = d.DeployNodeConfig()
	if err != nil {
		return err
	}

	return nil
}
