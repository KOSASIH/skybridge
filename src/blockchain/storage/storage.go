package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/skybridge/lib/errors"
)

// Storage represents a storage system
type Storage struct {
	// Config is the configuration for the storage
	Config Config

	// Mutex is a mutex to protect access to the storage
	Mutex sync.RWMutex

	// Data is the data stored in the storage
	Data map[string][]byte
}

// Config represents the configuration for the storage
type Config struct {
	// Path is the path to the storage directory
	Path string

	// Timeout is the timeout for storage operations
	Timeout time.Duration
}

// NewStorage returns a new storage system
func NewStorage(config Config) *Storage {
	return &Storage{
		Config: config,
		Data:   make(map[string][]byte),
	}
}

// Put puts a value into the storage
func (s *Storage) Put(key string, value []byte) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Data[key] = value
	return nil
}

// Get gets a value from the storage
func (s *Storage) Get(key string) ([]byte, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	value, ok := s.Data[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return value, nil
}

// Delete deletes a value from the storage
func (s *Storage) Delete(key string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Data, key)
	return nil
}

// MarshalJSON marshals the storage to JSON
func (s *Storage) MarshalJSON() ([]byte, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	data := make(map[string]string)
	for key, value := range s.Data {
		data[key] = string(value)
	}
	return json.Marshal(data)
}

// UnmarshalJSON unmarshals JSON to the storage
func (s *Storage) UnmarshalJSON(data []byte) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	var storageData map[string]string
	err := json.Unmarshal(data, &storageData)
	if err != nil {
		return err
	}
	s.Data = make(map[string][]byte)
	for key, value := range storageData {
		s.Data[key] = []byte(value)
	}
	return nil
}

// Persist persists the storage to disk
func (s *Storage) Persist() error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	data, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(s.Config.Path, "storage.json"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Load loads the storage from disk
func (s *Storage) Load() error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	data, err := ioutil.ReadFile(filepath.Join(s.Config.Path, "storage.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &s.Data)
	if err != nil {
		return err
	}
	return nil
}
