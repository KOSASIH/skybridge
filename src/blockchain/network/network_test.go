package network

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/skybridge/lib/errors"
)

func TestNetwork_NewNetwork(t *testing.T) {
	config := Config{
		Port:     8080,
		Timeout:  10 * time.Second,
		TLS: TLSConfig{
			Cert: "cert",
			Key:  "key",
		},
	}
	network := NewNetwork(config)
	if network == nil {
		t.Errorf("Expected NewNetwork to return a non-nil network")
	}
}

func TestNetwork_Start(t *testing.T) {
	config := Config{
		Port:     8080,
		Timeout:  10 * time.Second,
		TLS: TLSConfig{
			Cert: "cert",
			Key:  "key",
		},
	}
	network := NewNetwork(config)
	err := network.Start()
	if err != nil {
		t.Errorf("Expected Start to return a non-nil error")
	}
}

func TestNetwork_handleConn(t *testing.T) {
	config := Config{
		Port:     8080,
		Timeout:  10 * time.Second,
		TLS: TLSConfig{
			Cert: "cert",
			Key:  "key",
		},
	}
	network := NewNetwork(config)
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		t.Errorf("Expected to be able to dial the network")
	}
	defer conn.Close()

	peer := Peer{
		ID:      "peer-id",
		Address: "peer-address",
	}
	err = json.NewEncoder(conn).Encode(peer)
	if err != nil {
		t.Errorf("Expected to be able to encode the peer")
	}

	// Check that the peer was added to the list of peers
	network.Mutex.RLock()
	defer network.Mutex.RUnlock()
	for _, p := range network.Peers {
		if p.ID == peer.ID && p.Address == peer.Address {
			return
		}
	}
	t.Errorf("Expected the peer to be added to the list of peers")
}

func TestNetwork_Dial(t *testing.T) {
	config := Config{
		Port:     8080,
		Timeout:  10 * time.Second,
		TLS: TLSConfig{
			Cert: "cert",
			Key:  "key",
		},
	}
	network := NewNetwork(config)
	err := network.Start()
	if err != nil {
		t.Errorf("Expected Start to return a non-nil error")
	}

	peer := Peer{
		ID:      "peer-id",
		Address: "peer-address",
	}
	conn, err := network.Dial(peer)
	if err != nil {
		t.Errorf("Expected to be able to dial the peer")
	}
	defer conn.Close()

	// Check that the peer's ID and address were sent
	var receivedPeer Peer
	err = json.NewDecoder(conn).Decode(&receivedPeer)
	if err != nil {
		t.Errorf("Expected to be able to decode the peer")
	}
	if receivedPeer.ID != peer.ID || receivedPeer.Address != peer.Address {
		t.Errorf("Expected the peer's ID and address to be sent")
	}
}

func TestNetwork_MarshalJSON(t *testing.T) {
	config := Config{
		Port:     8080,
		Timeout:  10 * time.Second,
		TLS: TLSConfig{
			Cert: "cert",
			Key:  "key",
		},
	}
	network := NewNetwork(config)
	data, err := network.MarshalJSON()
	if err != nil {
		t.Errorf("Expected MarshalJSON to return a non-nil error")
	}

	var networkJSON struct {
		Config Config `json:"config"`
		Peers  []Peer `json:"peers"`
	}
	err = json.Unmarshal(data, &networkJSON)
	if err != nil {
		t.Errorf("Expected UnmarshalJSON to return a non-nil error")
	}
	if networkJSON.Config.Port != config.Port || networkJSON.Config.Timeout != config.Timeout {
		t.Errorf("Expected the config to be marshaled correctly")
	}
}

func TestNetwork_UnmarshalJSON(t *testing.T) {
	config := Config{
		Port:     8080,
		Timeout:  10 * time.Second,
		TLS: TLSConfig{
			Cert: "cert",
			Key:  "key",
		},
	}
	network := NewNetwork(config)
	data, err := json.Marshal(struct {
		Config Config `json:"config"`
		Peers  []Peer `json:"peers"`
	}{
		Config: config,
		Peers:  make([]Peer, 0),
	})
	if err != nil {
		t.Errorf("Expected json.Marshal to return a non-nil error")
	}
	err = network.UnmarshalJSON(data)
	if err != nil {
		t.Errorf("Expected UnmarshalJSON to return a non-nil error")
	}
	if network.Config.Port != config.Port || network.Config.Timeout != config.Timeout {
		t.Errorf("Expected the config to be unmarshaled correctly")
	}
}

func TestNetwork_GenerateTLS(t *testing.T) {
	cert, key, err := GenerateTLS()
	if err != nil {
		t.Errorf("Expected GenerateTLS to return a non-nil error")
	}
	if cert == "" || key == "" {
		t.Errorf("Expected GenerateTLS to return a non-empty cert and key")
	}
}
