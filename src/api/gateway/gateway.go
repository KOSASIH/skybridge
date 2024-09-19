package gateway

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/skybridge/api/types"
	"github.com/skybridge/crypto/encryption"
	"github.com/skybridge/crypto/signature"
	"github.com/skybridge/lib/errors"
	"github.com/skybridge/lib/logging"
	"github.com/skybridge/lib/utils"
)

// Gateway represents the API gateway for the skybridge network.
type Gateway struct {
	// Config is the configuration for the gateway.
	Config *Config

	// Router is the HTTP router for the gateway.
	Router *mux.Router

	// Server is the HTTP server for the gateway.
	Server *http.Server

	// WebSocketServer is the WebSocket server for the gateway.
	WebSocketServer *websocket.Server

	// TLSConfig is the TLS configuration for the gateway.
	TLSConfig *tls.Config

	// Logger is the logger for the gateway.
	Logger *logging.Logger

	// ErrorReporter is the error reporter for the gateway.
	ErrorReporter *errors.ErrorReporter

	// sync.RWMutex is used to protect access to the gateway's state.
	sync.RWMutex
}

// Config represents the configuration for the gateway.
type Config struct {
	// Address is the address to listen on.
	Address string

	// Port is the port to listen on.
	Port int

	// TLSKeyFile is the path to the TLS key file.
	TLSKeyFile string

	// TLSCertFile is the path to the TLS certificate file.
	TLSCertFile string

	// WebSocketAddress is the address to listen on for WebSocket connections.
	WebSocketAddress string

	// WebSocketPort is the port to listen on for WebSocket connections.
	WebSocketPort int

	// MaxConnections is the maximum number of connections allowed.
	MaxConnections int

	// MaxMessageSize is the maximum size of a message allowed.
	MaxMessageSize int

	// ReadTimeout is the timeout for reading from the connection.
	ReadTimeout time.Duration

	// WriteTimeout is the timeout for writing to the connection.
	WriteTimeout time.Duration
}

// NewGateway returns a new gateway instance.
func NewGateway(config *Config) (*Gateway, error) {
	gateway := &Gateway{
		Config: config,
		Router: mux.NewRouter(),
		Server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.Address, config.Port),
			Handler:      nil,
			TLSConfig:    nil,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
		WebSocketServer: &websocket.Server{
			Addr:         fmt.Sprintf("%s:%d", config.WebSocketAddress, config.WebSocketPort),
			Handler:      nil,
			TLSConfig:    nil,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{},
			RootCAs:      x509.NewCertPool(),
		},
		Logger:        logging.NewLogger(),
		ErrorReporter: errors.NewErrorReporter(),
	}

	gateway.Router.HandleFunc("/api/v1/nodes", gateway.handleGetNodes).Methods("GET")
	gateway.Router.HandleFunc("/api/v1/nodes/{id}", gateway.handleGetNode).Methods("GET")
	gateway.Router.HandleFunc("/api/v1/nodes", gateway.handleCreateNode).Methods("POST")
	gateway.Router.HandleFunc("/api/v1/nodes/{id}", gateway.handleUpdateNode).Methods("PUT")
	gateway.Router.HandleFunc("/api/v1/nodes/{id}", gateway.handleDeleteNode).Methods("DELETE")

	gateway.WebSocketServer.HandleFunc("/ws", gateway.handleWebSocketConnection)

	return gateway, nil
}

// Start starts the gateway.
func (g *Gateway) Start() error {
	g.Logger.Info("Starting gateway...")

	// Load TLS certificates
	cert, err := tls.LoadX509KeyPair(g.Config.TLSCertFile, g.Config.TLSKeyFile)
	if err != nil {
		return err
	}
	g.TLSConfig.Certificates = []tls.Certificate{cert}

	// Start HTTP server
	g.Server.Handler = g.Router
	g.Server.TLSConfig = g.TLSConfig
	go func() {
		g.Logger.Info("Listening on %s:%d", g.Config.Address, g.Config.Port)
		if err := g.Server.ListenAndServeTLS(g.Config.TLSCertFile, g.Config.TLSKeyFile); err != nil {
			g.Logger.Fatal("Failed to start HTTP server: %v", err)
		}
	}()

	// Start WebSocket server
	g.WebSocketServer.Handler = g.WebSocketServer
	g.WebSocketServer.TLSConfig = g.TLSConfig
	go func() {
		g.Logger.Info("Listening on %s:%d", g.Config.WebSocketAddress, g.Config.WebSocketPort)
		if err := g.WebSocketServer.ListenAndServeTLS(g.Config.TLSCertFile, g.Config.TLSKeyFile); err != nil {
			g.Logger.Fatal("Failed to start WebSocket server: %v", err)
		}
	}()

	return nil
}

// Stop stops the gateway.
func (g *Gateway) Stop() error {
	g.Logger.Info("Stopping gateway...")

	// Stop HTTP server
	if err := g.Server.Close(); err != nil {
		return err
	}

	// Stop WebSocket server
	if err := g.WebSocketServer.Close(); err != nil {
		return err
	}

	return nil
}

// handleGetNodes handles GET requests to /api/v1/nodes.
func (g *Gateway) handleGetNodes(w http.ResponseWriter, r *http.Request) {
	g.Logger.Info("Handling GET request to /api/v1/nodes")

	// Authenticate request
	if err := g.authenticateRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get nodes
	nodes, err := g.getNodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return nodes
	json.NewEncoder(w).Encode(nodes)
}

// handleGetNode handles GET requests to /api/v1/nodes/{id}.
func (g *Gateway) handleGetNode(w http.ResponseWriter, r *http.Request) {
	g.Logger.Info("Handling GET request to /api/v1/nodes/{id}")

	// Authenticate request
	if err := g.authenticateRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get node ID
	nodeID := mux.Vars(r)["id"]

	// Get node
	node, err := g.getNode(nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return node
	json.NewEncoder(w).Encode(node)
}

// handleCreateNode handles POST requests to /api/v1/nodes.
func (g *Gateway) handleCreateNode(w http.ResponseWriter, r *http.Request) {
	g.Logger.Info("Handling POST request to /api/v1/nodes")

	// Authenticate request
	if err := g.authenticateRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Decode request body
	var node types.Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create node
	if err := g.createNode(node); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created node
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(node)
}

// handleUpdateNode handles PUT requests to /api/v1/nodes/{id}.
func (g *Gateway) handleUpdateNode(w http.ResponseWriter, r *http.Request) {
	g.Logger.Info("Handling PUT request to /api/v1/nodes/{id}")

	// Authenticate request
	if err := g.authenticateRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get node ID
	nodeID := mux.Vars(r)["id"]

	// Decode request body
	var node types.Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update node
	if err := g.updateNode(nodeID, node); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated node
	json.NewEncoder(w).Encode(node)
}

// handleDeleteNode handles DELETE requests to /api/v1/nodes/{id}.
func (g *Gateway) handleDeleteNode(w http.ResponseWriter, r *http.Request) {
	g.Logger.Info("Handling DELETE request to /api/v1/nodes/{id}")

	// Authenticate request
	if err := g.authenticateRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get node ID
	nodeID := mux.Vars(r)["id"]

	// Delete node
	if err := g.deleteNode(nodeID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}

// handleWebSocketConnection handles WebSocket connections.
func (g *Gateway) handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	g.Logger.Info("Handling WebSocket connection")

	// Upgrade to WebSocket connection
	conn, err := g.WebSocketServer.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		g.Logger.Error("Failed to upgrade to WebSocket connection: %v", err)
		return
	}

	// Handle WebSocket messages
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				g.Logger.Error("Failed to read WebSocket message: %v", err)
				return
			}

			// Handle message
			g.handleWebSocketMessage(message)
		}
	}()
}

// handleWebSocketMessage handles WebSocket messages.
func (g *Gateway) handleWebSocketMessage(message []byte) {
	g.Logger.Info("Handling WebSocket message")

	// Decode message
	var msg types.WebSocketMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		g.Logger.Error("Failed to decode WebSocket message: %v", err)
		return
	}

	// Handle message
	switch msg.Type {
	case types.WebSocketMessageType_NODE_CREATED:
		g.handleNodeCreated(msg.Data)
	case types.WebSocketMessageType_NODE_UPDATED:
		g.handleNodeUpdated(msg.Data)
	case types.WebSocketMessageType_NODE_DELETED:
		g.handleNodeDeleted(msg.Data)
	default:
		g.Logger.Error("Unknown WebSocket message type: %v", msg.Type)
	}
}

// authenticateRequest authenticates a request.
func (g *Gateway) authenticateRequest(r *http.Request) error {
	// Get authentication token from request header
	token := r.Header.Get("Authorization")

	// Verify token
	if err := g.verifyToken(token); err != nil {
		return err
	}

	return nil
}

// verifyToken verifies an authentication token.
func (g *Gateway) verifyToken(token string) error {
	// Decode token
	var claims types.Claims
	if err := jwt.Unmarshal(token, &claims); err != nil {
		return err
	}

	// Verify claims
	if err := g.verifyClaims(claims); err != nil {
		return err
	}

	return nil
}

// verifyClaims verifies authentication claims.
func (g *Gateway) verifyClaims(claims types.Claims) error {
	// Verify issuer
	if claims.Issuer != g.Config.Issuer {
		return errors.New("Invalid issuer")
	}

	// Verify audience
	if claims.Audience != g.Config.Audience {
		return errors.New("Invalid audience")
	}

	// Verify expiration
	if claims.ExpiresAt < time.Now().Unix() {
		return errors.New("Token has expired")
	}

	return nil
}
