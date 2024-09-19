package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/skybridge/lib/errors"
)

// NodeType represents a node type
type NodeType string

const (
	// NodeTypeServer represents a server node
	NodeTypeServer NodeType = "server"
	// NodeTypeClient represents a client node
	NodeTypeClient NodeType = "client"
)

// NodeStatus represents a node status
type NodeStatus string

const (
	// NodeStatusOnline represents an online node
	NodeStatusOnline NodeStatus = "online"
	// NodeStatusOffline represents an offline node
	NodeStatusOffline NodeStatus = "offline"
)

// NodeMetadata represents node metadata
type NodeMetadata struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      NodeType  `json:"type"`
	Status    NodeStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MarshalJSON marshals the node metadata to JSON
func (m *NodeMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Type      NodeType  `json:"type"`
		Status    NodeStatus `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        m.ID,
		Name:      m.Name,
		Type:      m.Type,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	})
}

// UnmarshalJSON unmarshals JSON to node metadata
func (m *NodeMetadata) UnmarshalJSON(data []byte) error {
	var node struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Type      NodeType  `json:"type"`
		Status    NodeStatus `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	if err := json.Unmarshal(data, &node); err != nil {
		return err
	}
	m.ID = node.ID
	m.Name = node.Name
	m.Type = node.Type
	m.Status = node.Status
	m.CreatedAt = node.CreatedAt
	m.UpdatedAt = node.UpdatedAt
	return nil
}

// Validate validates the node metadata
func (m *NodeMetadata) Validate() error {
	if m.ID == "" {
		return errors.NewValidationError("id", "cannot be empty")
	}
	if m.Name == "" {
		return errors.NewValidationError("name", "cannot be empty")
	}
	if m.Type == "" {
		return errors.NewValidationError("type", "cannot be empty")
	}
	if m.Status == "" {
		return errors.NewValidationError("status", "cannot be empty")
	}
	return nil
}

// String returns a string representation of the node metadata
func (m *NodeMetadata) String() string {
	return fmt.Sprintf("NodeMetadata{ID: %s, Name: %s, Type: %s, Status: %s, CreatedAt: %v, UpdatedAt: %v}",
		m.ID, m.Name, m.Type, m.Status, m.CreatedAt, m.UpdatedAt)
}
