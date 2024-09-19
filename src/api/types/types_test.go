package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/skybridge/lib/errors"
)

func TestNodeMetadata_MarshalJSON(t *testing.T) {
	m := &NodeMetadata{
		ID:        "node-1",
		Name:      "Node 1",
		Type:      NodeTypeServer,
		Status:    NodeStatusOnline,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	var node struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Type      NodeType  `json:"type"`
		Status    NodeStatus `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	err = json.Unmarshal(jsonBytes, &node)
	if err != nil {
		t.Fatal(err)
	}

	if node.ID != m.ID {
		t.Errorf("Expected ID %s, got %s", m.ID, node.ID)
	}
	if node.Name != m.Name {
		t.Errorf("Expected Name %s, got %s", m.Name, node.Name)
	}
	if node.Type != m.Type {
		t.Errorf("Expected Type %s, got %s", m.Type, node.Type)
	}
	if node.Status != m.Status {
		t.Errorf("Expected Status %s, got %s", m.Status, node.Status)
	}
	if !node.CreatedAt.Equal(m.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", m.CreatedAt, node.CreatedAt)
	}
	if !node.UpdatedAt.Equal(m.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", m.UpdatedAt, node.UpdatedAt)
	}
}

func TestNodeMetadata_UnmarshalJSON(t *testing.T) {
	jsonBytes := []byte(`{"id": "node-1", "name": "Node 1", "type": "server", "status": "online", "created_at": "2022-01-01T00:00:00Z", "updated_at": "2022-01-01T00:00:00Z"}`)

	m := &NodeMetadata{}
	err := json.Unmarshal(jsonBytes, m)
	if err != nil {
		t.Fatal(err)
	}

	if m.ID != "node-1" {
		t.Errorf("Expected ID %s, got %s", "node-1", m.ID)
	}
	if m.Name != "Node 1" {
		t.Errorf("Expected Name %s, got %s", "Node 1", m.Name)
	}
	if m.Type != NodeTypeServer {
		t.Errorf("Expected Type %s, got %s", NodeTypeServer, m.Type)
	}
	if m.Status != NodeStatusOnline {
		t.Errorf("Expected Status %s, got %s", NodeStatusOnline, m.Status)
	}
	if !m.CreatedAt.Equal(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected CreatedAt %v, got %v", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), m.CreatedAt)
	}
	if !m.UpdatedAt.Equal(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected UpdatedAt %v, got %v", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), m.UpdatedAt)
	}
}

func TestNodeMetadata_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *NodeMetadata
		wantErr bool
	}{
		{
			name: "valid",
			m: &NodeMetadata{
				ID:        "node-1",
				Name:      "Node 1",
				Type:      NodeTypeServer,
				Status:    NodeStatusOnline,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			m: &NodeMetadata{
				ID:        "",
				Name:      "Node 1",
				Type:      NodeTypeServer,
				Status:    NodeStatusOnline,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty Name",
			m: &NodeMetadata{
				ID:        "node-1",
				Name:      "",
				Type:      NodeTypeServer,
				Status:    NodeStatusOnline,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty Type",
			m: &NodeMetadata{
				ID:        "node-1",
				Name:      "Node 1",
				Type:      "",
				Status:    NodeStatusOnline,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty Status",
			m: &NodeMetadata{
				ID:        "node-1",
				Name:      "Node 1",
				Type:      NodeTypeServer,
				Status:    "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeMetadata_String(t *testing.T) {
	m := &NodeMetadata{
		ID:        "node-1",
		Name:      "Node 1",
		Type:      NodeTypeServer,
		Status:    NodeStatusOnline,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s := m.String()
	if !strings.Contains(s, "NodeMetadata") {
		t.Errorf("Expected string to contain 'NodeMetadata', got %s", s)
	}
	if !strings.Contains(s, "ID: node-1") {
		t.Errorf("Expected string to contain 'ID: node-1', got %s", s)
	}
	if !strings.Contains(s, "Name: Node 1") {
		t.Errorf("Expected string to contain 'Name: Node 1', got %s", s)
	}
	if !strings.Contains(s, "Type: server") {
		t.Errorf("Expected string to contain 'Type: server', got %s", s)
	}
	if !strings.Contains(s, "Status: online") {
		t.Errorf("Expected string to contain 'Status: online', got %s", s)
	}
	if !strings.Contains(s, "CreatedAt: ") {
		t.Errorf("Expected string to contain 'CreatedAt: ', got %s", s)
	}
	if !strings.Contains(s, "UpdatedAt: ") {
		t.Errorf("Expected string to contain 'UpdatedAt: ', got %s", s)
	}
}
