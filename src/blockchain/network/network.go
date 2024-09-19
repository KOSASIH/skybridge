package network

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/skybridge/lib/errors"
)

// Network represents a network
type Network struct {
	// Config is the configuration for the network
	Config Config

	// Peers is a list of peers in the network
	Peers []Peer

	// Mutex is a mutex to protect access to the peers
	Mutex sync.RWMutex

	// Listener is the listener for incoming connections
	Listener net.Listener

	// Dialer is the dialer for outgoing connections
	Dialer net.Dialer
}

// Config represents the configuration for the network
type Config struct {
	// Port is the port to listen on
	Port int

	// Timeout is the timeout for connections
	Timeout time.Duration

	// TLS is the TLS configuration
	TLS TLSConfig
}

// TLSConfig represents the TLS configuration
type TLSConfig struct {
	// Cert is the TLS certificate
	Cert string

	// Key is the TLS private key
	Key string
}

// Peer represents a peer in the network
type Peer struct {
	// ID is the ID of the peer
	ID string

	// Address is the address of the peer
	Address string

	// Conn is the connection to the peer
	Conn net.Conn
}

// NewNetwork returns a new network
func NewNetwork(config Config) *Network {
	return &Network{
		Config: config,
		Peers:  make([]Peer, 0),
	}
}

// Start starts the network
func (n *Network) Start() error {
	// Create a new listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", n.Config.Port))
	if err != nil {
		return err
	}
	n.Listener = listener

	// Create a new dialer
	dialer := net.Dialer{
		Timeout: n.Config.Timeout,
	}
	n.Dialer = dialer

	// Start listening for incoming connections
	go func() {
		for {
			conn, err := n.Listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go n.handleConn(conn)
		}
	}()

	return nil
}

// handleConn handles an incoming connection
func (n *Network) handleConn(conn net.Conn) {
	// Read the peer's ID and address
	var peer Peer
	err := json.NewDecoder(conn).Decode(&peer)
	if err != nil {
		log.Println(err)
		return
	}

	// Add the peer to the list of peers
	n.Mutex.Lock()
	n.Peers = append(n.Peers, peer)
	n.Mutex.Unlock()

	// Send a response to the peer
	err = json.NewEncoder(conn).Encode(struct {
		ID string `json:"id"`
	}{
		ID: n.Config.TLS.Cert,
	})
	if err != nil {
		log.Println(err)
		return
	}

	// Close the connection
	conn.Close()
}

// Dial dials a peer
func (n *Network) Dial(peer Peer) (net.Conn, error) {
	// Dial the peer
	conn, err := n.Dialer.Dial("tcp", peer.Address)
	if err != nil {
		return nil, err
	}

	// Send the peer's ID and address
	err = json.NewEncoder(conn).Encode(peer)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// MarshalJSON marshals the network to JSON
func (n *Network) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Config Config `json:"config"`
		Peers  []Peer `json:"peers"`
	}{
		Config: n.Config,
		Peers:  n.Peers,
	})
}

// UnmarshalJSON unmarshals JSON to the network
func (n *Network) UnmarshalJSON(data []byte) error {
	var network struct {
		Config Config `json:"config"`
		Peers  []Peer `json:"peers"`
	}
	err := json.Unmarshal(data, &network)
	if err != nil {
		return err
	}
	n.Config = network.Config
	n.Peers = network.Peers
	return nil
}

// GenerateTLS generates a TLS certificate and private key
func GenerateTLS() (string, string, error) {
	// Generate a new private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Generate a new certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Skybridge"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),

		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	// Encode the certificate and private key
	cert := pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: derBytes},
	)
	key := pem.EncodeToMemory(
		&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)},
	)

	return string(cert), string(key), nil
}
