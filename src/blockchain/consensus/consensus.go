package consensus

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/skybridge/lib/errors"
)

// Consensus represents a consensus algorithm
type Consensus struct {
	// Config is the configuration for the consensus algorithm
	Config Config

	// State is the current state of the consensus algorithm
	State State

	// Mutex is a mutex to protect access to the state
	Mutex sync.RWMutex

	// Peers is a list of peers in the consensus algorithm
	Peers []Peer

	// Block is the current block being proposed
	Block *Block

	// Votes is a map of votes for the current block
	Votes map[string]bool

	// Timer is a timer to trigger the next round of the consensus algorithm
	Timer *time.Timer
}

// Config represents the configuration for the consensus algorithm
type Config struct {
	// RoundDuration is the duration of each round of the consensus algorithm
	RoundDuration time.Duration

	// BlockTimeout is the timeout for proposing a block
	BlockTimeout time.Duration

	// VoteTimeout is the timeout for voting on a block
	VoteTimeout time.Duration
}

// State represents the current state of the consensus algorithm
type State struct {
	// Round is the current round of the consensus algorithm
	Round uint64

	// Leader is the leader of the current round
	Leader string

	// Block is the current block being proposed
	Block *Block
}

// Peer represents a peer in the consensus algorithm
type Peer struct {
	// ID is the ID of the peer
	ID string

	// Address is the address of the peer
	Address string
}

// Block represents a block in the consensus algorithm
type Block struct {
	// Hash is the hash of the block
	Hash string

	// Transactions is a list of transactions in the block
	Transactions []Transaction
}

// Transaction represents a transaction in the consensus algorithm
type Transaction struct {
	// From is the sender of the transaction
	From string

	// To is the recipient of the transaction
	To string

	// Amount is the amount of the transaction
	Amount uint64
}

// NewConsensus returns a new consensus algorithm
func NewConsensus(config Config) *Consensus {
	return &Consensus{
		Config: config,
		State: State{
			Round: 0,
		},
		Peers: make([]Peer, 0),
		Votes: make(map[string]bool),
	}
}

// Start starts the consensus algorithm
func (c *Consensus) Start() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.State.Round = 0
	c.State.Leader = c.Peers[0].ID
	c.State.Block = &Block{
		Hash: "",
		Transactions: make([]Transaction, 0),
	}

	c.Timer = time.AfterFunc(c.Config.RoundDuration, c.Round)
}

// Round is the main loop of the consensus algorithm
func (c *Consensus) Round() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// Propose a new block
	c.ProposeBlock()

	// Vote on the block
	c.Vote()

	// Check if the block is valid
	if c.IsValid() {
		// Add the block to the blockchain
		c.AddBlock()

		// Move to the next round
		c.State.Round++
		c.State.Leader = c.Peers[(c.State.Round)%len(c.Peers)].ID
		c.State.Block = &Block{
			Hash: "",
			Transactions: make([]Transaction, 0),
		}
	} else {
		// Restart the round
		c.Round()
	}
}

// ProposeBlock proposes a new block
func (c *Consensus) ProposeBlock() {
	// Generate a random hash for the block
	hash := fmt.Sprintf("%x", sha256.Sum256(rand.Reader))

	// Create a new block
	block := &Block{
		Hash: hash,
		Transactions: make([]Transaction, 0),
	}

	// Add transactions to the block
	for i := 0; i < 10; i++ {
		transaction := Transaction{
			From: c.Peers[i%len(c.Peers)].ID,
			To:   c.Peers[(i+1)%len(c.Peers)].ID,
			Amount: uint64(i),
		}
		block.Transactions = append(block.Transactions, transaction)
	}

	// Set the block as the current block
	c.Block = block
}

// Vote votes on the current block
func (c *Consensus) Vote() {
	// Vote for the block
	c.Votes[c.State.Leader] = true

	// Check if all peers have voted
	if len(c.Votes) == len(c.Peers) {
		// Move to the next round
		c.Round()
	} else {
		// Wait for the next peer to vote
		time.Sleep(c.Config.VoteTimeout)
		c.Vote()
	}
}

// IsValid checks if the current block is valid
func (c *Consensus) IsValid() bool {
	// Check if the block has a valid hash
	if c.Block.Hash == "" {
		return false
	}

	// Check if all transactions are valid
	for _, transaction := range c.Block.Transactions {
		if transaction.From == "" || transaction.To == "" || transaction.Amount == 0 {
			return false
		}
	}

	return true
}

// AddBlock adds the current block to the blockchain
func (c *Consensus) AddBlock() {
	// Add the block to the blockchain
	// ...

	// Update the state
	c.State.Block = c.Block
}

// MarshalJSON marshals the consensus algorithm to JSON
func (c *Consensus) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Config Config `json:"config"`
		State State   `json:"state"`
		Peers []Peer `json:"peers"`
		Block *Block `json:"block"`
		Votes map[string]bool `json:"votes"`
	}{
		Config: c.Config,
		State:  c.State,
		Peers:  c.Peers,
		Block:  c.Block,
		Votes:  c.Votes,
	})
}

// UnmarshalJSON unmarshals JSON to the consensus algorithm
func (c *Consensus) UnmarshalJSON(data []byte) error {
	var consensus struct {
		Config Config `json:"config"`
		State State   `json:"state"`
		Peers []Peer `json:"peers"`
		Block *Block `json:"block"`
		Votes map[string]bool `json:"votes"`
	}
	if err := json.Unmarshal(data, &consensus); err != nil {
		return err
	}
	c.Config = consensus.Config
	c.State = consensus.State
	c.Peers = consensus.Peers
	c.Block = consensus.Block
	c.Votes = consensus.Votes
	return nil
}
