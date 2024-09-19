package node

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/skybridge/api/types"
	"github.com/skybridge/lib/errors"
	"github.com/skybridge/lib/logging"
)

// Node represents a node in the network
type Node struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IPAddress string    `json:"ip_address"`
	Port      int       `json:"port"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	mu sync.RWMutex
}

// NewNode creates a new node
func NewNode(id, name, ipAddress string, port int) *Node {
	return &Node{
		ID:        id,
		Name:      name,
		IPAddress: ipAddress,
		Port:      port,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GetID returns the node's ID
func (n *Node) GetID() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.ID
}

// GetName returns the node's name
func (n *Node) GetName() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Name
}

// GetIPAddress returns the node's IP address
func (n *Node) GetIPAddress() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.IPAddress
}

// GetPort returns the node's port
func (n *Node) GetPort() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Port
}

// GetCreatedAt returns the node's creation time
func (n *Node) GetCreatedAt() time.Time {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.CreatedAt
}

// GetUpdatedAt returns the node's update time
func (n *Node) GetUpdatedAt() time.Time {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.UpdatedAt
}

// Update updates the node's information
func (n *Node) Update(name, ipAddress string, port int) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Name = name
	n.IPAddress = ipAddress
	n.Port = port
	n.UpdatedAt = time.Now()
	return nil
}

// MarshalJSON marshals the node to JSON
func (n *Node) MarshalJSON() ([]byte, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return json.Marshal(struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		IPAddress string    `json:"ip_address"`
		Port      int       `json:"port"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        n.ID,
		Name:      n.Name,
		IPAddress: n.IPAddress,
		Port:      n.Port,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	})
}

// UnmarshalJSON unmarshals JSON to a node
func (n *Node) UnmarshalJSON(data []byte) error {
	var node struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		IPAddress string    `json:"ip_address"`
		Port      int       `json:"port"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	if err := json.Unmarshal(data, &node); err != nil {
		return err
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.ID = node.ID
	n.Name = node.Name
	n.IPAddress = node.IPAddress
	n.Port = node.Port
	n.CreatedAt = node.CreatedAt
	n.UpdatedAt = node.UpdatedAt
	return nil
}

// NodeService represents a node service
type NodeService interface {
	GetNode(ctx context.Context, id string) (*Node, error)
	CreateNode(ctx context.Context, node *Node) error
	UpdateNode(ctx context.Context, node *Node) error
	DeleteNode(ctx context.Context, id string) error
}

// NewNodeService creates a new node service
func NewNodeService(logger logging.Logger) NodeService {
	return &nodeService{
		logger: logger,
		nodes:  make(map[string]*Node),
	}
}

type nodeService struct {
	logger logging.Logger
	nodes  map[string]*Node
	mu     sync.RWMutex
}

func (s *nodeService) GetNode(ctx context.Context, id string) (*Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	node, ok := s.nodes[id]
	if !ok {
		return nil, errors.NewNotFoundError("node not found")
	}
	return node, nil
}

func (s *nodeService) CreateNode(ctx context.Context, node *Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes[node.ID] = node
	return nil
}

func (s *nodeService) UpdateNode(ctx context.Context, node *Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes[node.ID] = node
	return nil
}

func (s *nodeService) DeleteNode(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.nodes, id)
	return nil
}

func (s *nodeService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	switch r.Method {
	case http.MethodGet:
		node, err := s.GetNode(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(node)
	case http.MethodPost:
		var node Node
		if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.CreateNode(r.Context(), &node); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodPut:
		var node Node
		if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.UpdateNode(r.Context(), &node); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		if err := s.DeleteNode(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
