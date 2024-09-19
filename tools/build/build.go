package build

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-username/satellite"
	"github.com/your-username/node"
)

// Builder represents a builder for satellite and node configurations
type Builder struct {
	satelliteConfig *satellite.Satellite
	nodeConfig     *node.Node
}

// NewBuilder returns a new Builder instance
func NewBuilder(satelliteConfig *satellite.Satellite, nodeConfig *node.Node) *Builder {
	return &Builder{
		satelliteConfig: satelliteConfig,
		nodeConfig:     nodeConfig,
	}
}

// BuildSatelliteConfig builds the satellite configuration
func (b *Builder) BuildSatelliteConfig() error {
	log.Println("Building satellite configuration")
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
	_, err = file.WriteString(fmt.Sprintf("id=%s\n", b.satelliteConfig.id))
	if err != nil {
		return err
	}

	// Write the ground stations to the file
	for _, groundStation := range b.satelliteConfig.groundStations {
		_, err = file.WriteString(fmt.Sprintf("ground_station=%s\n", groundStation.id))
		if err != nil {
			return err
		}
	}

	// Write the transceivers to the file
	for _, transceiver := range b.satelliteConfig.transceivers {
		_, err = file.WriteString(fmt.Sprintf("transceiver=%s\n", transceiver.id))
		if err != nil {
			return err
		}
	}

	return nil
}

// BuildNodeConfig builds the node configuration
func (b *Builder) BuildNodeConfig() error {
	log.Println("Building node configuration")
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
	_, err = file.WriteString(fmt.Sprintf("id=%s\n", b.nodeConfig.id))
	if err != nil {
		return err
	}

	// Write the neighbors to the file
	for _, neighbor := range b.nodeConfig.neighbors {
		_, err = file.WriteString(fmt.Sprintf("neighbor=%s\n", neighbor.id))
		if err != nil {
			return err
		}
	}

	return nil
}

// Build builds both the satellite and node configurations
func (b *Builder) Build() error {
	err := b.BuildSatelliteConfig()
	if err != nil {
		return err
	}

	err = b.BuildNodeConfig()
	if err != nil {
		return err
	}

	return nil
}
