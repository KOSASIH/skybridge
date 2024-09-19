package node

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/skybridge/api/types"
	"github.com/skybridge/lib/errors"
	"github.com/skybridge/lib/logging"
)

func TestNodeService(t *testing.T) {
	logger := logging.NewLogger("node_test")
	nodeService := NewNodeService(logger)

	t.Run("GetNode", func(t *testing.T) {
		node := NewNode("node-1", "Node 1", "192.168.1.1", 8080)
		nodeService.CreateNode(context.Background(), node)

		req, err := http.NewRequest("GET", "/api/v1/nodes/node-1", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		nodeService.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var gotNode Node
		err = json.NewDecoder(w.Body).Decode(&gotNode)
		if err != nil {
			t.Fatal(err)
		}

		if gotNode.ID != node.ID {
			t.Errorf("Expected node ID %s, got %s", node.ID, gotNode.ID)
		}
	})

	t.Run("CreateNode", func(t *testing.T) {
		node := NewNode("node-2", "Node 2", "192.168.1.2", 8081)

		req, err := http.NewRequest("POST", "/api/v1/nodes", strings.NewReader(`{"name": "Node 2", "ip_address": "192.168.1.2", "port": 8081}`))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		nodeService.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		gotNode, err := nodeService.GetNode(context.Background(), node.ID)
		if err != nil {
			t.Fatal(err)
		}

		if gotNode.ID != node.ID {
			t.Errorf("Expected node ID %s, got %s", node.ID, gotNode.ID)
		}
	})

	t.Run("UpdateNode", func(t *testing.T) {
		node := NewNode("node-3", "Node 3", "192.168.1.3", 8082)
		nodeService.CreateNode(context.Background(), node)

		req, err := http.NewRequest("PUT", "/api/v1/nodes/node-3", strings.NewReader(`{"name": "Node 3 Updated", "ip_address": "192.168.1.3", "port": 8082}`))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		nodeService.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		gotNode, err := nodeService.GetNode(context.Background(), node.ID)
		if err != nil {
			t.Fatal(err)
		}

		if gotNode.Name != "Node 3 Updated" {
			t.Errorf("Expected node name %s, got %s", "Node 3 Updated", gotNode.Name)
		}
	})

	t.Run("DeleteNode", func(t *testing.T) {
		node := NewNode("node-4", "Node 4", "192.168.1.4", 8083)
		nodeService.CreateNode(context.Background(), node)

		req, err := http.NewRequest("DELETE", "/api/v1/nodes/node-4", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		nodeService.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		}

		_, err = nodeService.GetNode(context.Background(), node.ID)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
