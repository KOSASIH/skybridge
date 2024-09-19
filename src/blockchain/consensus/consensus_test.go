package consensus

import (
	"encoding/json"
	"testing"
	"time"

func TestConsensus_NewConsensus(t *testing.T) {
	config := Config{
		RoundDuration: 10 * time.Second,
		BlockTimeout:  5 * time.Second,
		VoteTimeout:   2 * time.Second,
	}
	consensus := NewConsensus(config)
	if consensus == nil {
		t.Errorf("Expected NewConsensus to return a non-nil consensus")
	}
}

func TestConsensus_Start(t *testing.T) {
	config := Config{
		RoundDuration: 10 * time.Second,
		BlockTimeout:  5 * time.Second,
		VoteTimeout:   2 * time.Second,
	}
	consensus := NewConsensus(config)
	consensus.Start()
	if consensus.State.Round != 0 {
		t.Errorf("Expected Start to set Round to 0, got %d", consensus.State.Round)
	}
}

func TestConsensus_Round(t *testing.T) {
	config := Config{
		RoundDuration: 10 * time.Second,
		BlockTimeout:  5 * time.Second,
		VoteTimeout:   2 * time.Second,
	}
	consensus := NewConsensus(config)
	consensus.Start()
	consensus.Round()
	if consensus.State.Round != 1 {
		t.Errorf("Expected Round to increment Round to 1, got %d", consensus.State.Round)
	}
}

func TestConsensus_ProposeBlock(t *testing.T) {
	config := Config{
		RoundDuration: 10 * time.Second,
		BlockTimeout:  5 * time.Second,
		VoteTimeout:   2 * time.Second,
	}
	consensus := NewConsensus(config)
	consensus.Start()
	consensus.ProposeBlock()
	if consensus.Block == nil {
		t.Errorf("Expected ProposeBlock to set Block to a non-nil value")
	}
}

func TestConsensus_Vote(t *testing.T) {
	config := Config{
		RoundDuration: 10 * time.Second,
		BlockTimeout:  5 * time.Second,
		VoteTimeout:   2 * time.Second,
	}
	consensus := NewConsensus(config)
	consensus.Start()
	consensus.ProposeBlock()
	consensus.Vote()
	if len(consensus.Votes) != 1 {
		t.Errorf("Expected Vote to add a vote to the Votes map")
	}
}

func TestConsensus_IsValid(t *testing.T) {
	config := Config{
		RoundDuration: 10 * time.Second,
		BlockTimeout:  5 * time.Second,
		VoteTimeout:   2 * time.Second,
	}
	consensus := NewConsensus(config)
	consensus.Start()
	consensus.ProposeBlock()
	block := &Block{
		Hash: "valid-hash",
		Transactions: []Transaction{
			{
				From: "from",
				To:   "to",
				Amount: 1,
			},
		},
	}
	consensus.Block = block
	if !consensus.IsValid() {
		t.Errorf("Expected IsValid to return true for a valid block")
	}
}

func TestConsensus_MarshalJSON(t *testing.T) {
	config := Config{
		RoundDuration: 10 * time.Second,
		BlockTimeout:  5 * time.Second,
		VoteTimeout:   2 * time.Second,
	}
	consensus := NewConsensus(config)
	consensus.Start()
	data, err := consensus.MarshalJSON()
	if err != nil {
		t.Errorf("Expected MarshalJSON to return a non-nil error")
	}
	var consensusJSON struct {
		Config Config `json:"config"`
		State State   `json:"state"`
		Peers []Peer `json:"peers"`
		Block *Block `json:"block"`
		Votes map[string]bool `json:"votes"`
	}
	err = json.Unmarshal(data, &consensusJSON)
	if err != nil {
		t.Errorf("Expected UnmarshalJSON to return a non-nil error")
	}
}

func TestConsensus_UnmarshalJSON(t *testing.T) {
	config := Config{
		RoundDuration: 10 * time.Second,
		BlockTimeout:  5 * time.Second,
		VoteTimeout:   2 * time.Second,
	}
	consensus := NewConsensus(config)
	consensus.Start()
	data, err := json.Marshal(struct {
		Config Config `json:"config"`
		State State   `json:"state"`
		Peers []Peer `json:"peers"`
		Block *Block `json:"block"`
		Votes map[string]bool `json:"votes"`
	}{
		Config: config,
		State: State{
			Round: 0,
		},
		Peers: make([]Peer, 0),
		Block: &Block{
			Hash: "",
			Transactions: make([]Transaction, 0),
		},
		Votes: make(map[string]bool),
	})
	if err != nil {
		t.Errorf("Expected json.Marshal to return a non-nil error")
	}
	err = consensus.UnmarshalJSON(data)
	if err != nil {
		t.Errorf("Expected UnmarshalJSON to return a non-nil error")
	}
}	
