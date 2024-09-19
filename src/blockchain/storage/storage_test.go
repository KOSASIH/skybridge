package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestStorage_Put(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key", []byte("value"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
}

func TestStorage_Get(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key", []byte("value"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	value, err := storage.Get("key")
	if err != nil {
		t.Errorf("Expected Get to return a non-nil error")
	}
	if string(value) != "value" {
		t.Errorf("Expected Get to return the correct value")
	}
}

func TestStorage_Delete(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key", []byte("value"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	err = storage.Delete("key")
	if err != nil {
		t.Errorf("Expected Delete to return a non-nil error")
	}
	_, err = storage.Get("key")
	if err == nil {
		t.Errorf("Expected Get to return a non-nil error after Delete")
	}
}

func TestStorage_MarshalJSON(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key", []byte("value"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	data, err := json.Marshal(storage)
	if err != nil {
		t.Errorf("Expected MarshalJSON to return a non-nil error")
	}
	var unmarshaled map[string][]byte
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Expected UnmarshalJSON to return a non-nil error")
	}
	if string(unmarshaled["key"]) != "value" {
		t.Errorf("Expected MarshalJSON and UnmarshalJSON to work correctly")
	}
}

func TestStorage_UnmarshalJSON(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	data, err := json.Marshal(map[string][]byte{"key": []byte("value")})
	if err != nil {
		t.Errorf("Expected MarshalJSON to return a non-nil error")
	}
	err = json.Unmarshal(data, &storage)
	if err != nil {
		t.Errorf("Expected UnmarshalJSON to return a non-nil error")
	}
	value, err := storage.Get("key")
	if err != nil {
		t.Errorf("Expected Get to return a non-nil error")
	}
	if string(value) != "value" {
		t.Errorf("Expected UnmarshalJSON to work correctly")
	}
}

func TestStorage_Persist(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key", []byte("value"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	err = storage.Persist()
	if err != nil {
		t.Errorf("Expected Persist to return a non-nil error")
	}
	storage2 := NewStorage(config)
	err = storage2.Load()
	if err != nil {
		t.Errorf("Expected Load to return a non-nil error")
	}
	value, err := storage2.Get("key")
	if err != nil {
		t.Errorf("Expected Get to return a non-nil error")
	}
	if string(value) != "value" {
		t.Errorf("Expected Persist and Load to work correctly")
	}
}

func TestStorage_Load(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key", []byte("value"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	err = storage.Persist()
	if err != nil {
		t.Errorf("Expected Persist to return a non-nil error")
	}
	storage2 := NewStorage(config)
	err = storage2.Load()
	if err != nil {
		t.Errorf("Expected Load to return a non-nil error")
	}
	value, err := storage2.Get("key")
	if err != nil {
		t.Errorf("Expected Get to return a non-nil error")
	}
	if string(value) != "value" {
		t.Errorf("Expected Load to work correctly")
	}
}

func TestStorage_Close(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Close()
	if err != nil {
		t.Errorf("Expected Close to return a non-nil error")
	}
}

func TestStorage_GetSize(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key1", []byte("value1"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	err = storage.Put("key2", []byte("value2"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	size := storage.GetSize()
	if size != 2 {
		t.Errorf("Expected GetSize to return the correct size")
	}
}

func TestStorage_GetKeys(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key1", []byte("value1"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	err = storage.Put("key2", []byte("value2"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	keys := storage.GetKeys()
	if len(keys) != 2 || keys[0] != "key1" || keys[1] != "key2" {
		t.Errorf("Expected GetKeys to return the correct keys")
	}
}

func TestStorage_GetValues(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key1", []byte("value1"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	err = storage.Put("key2", []byte("value2"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	values := storage.GetValues()
	if len(values) != 2 || string(values[0]) != "value1" || string(values[1]) != "value2" {
		t.Errorf("Expected GetValues to return the correct values")
	}
}

func TestStorage_Iterate(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key1", []byte("value1"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	err = storage.Put("key2", []byte("value2"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	i := 0
	storage.Iterate(func(key string, value []byte) {
		i++
		if i == 1 {
			if key != "key1" || string(value) != "value1" {
				t.Errorf("Expected Iterate to return the correct key and value")
			}
		} else if i == 2 {
			if key != "key2" || string(value) != "value2" {
				t.Errorf("Expected Iterate to return the correct key and value")
			}
		}
	})
	if i != 2 {
		t.Errorf("Expected Iterate to return the correct number of items")
	}
}

func TestStorage_Clear(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	err := storage.Put("key1", []byte("value1"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	err = storage.Put("key2", []byte("value2"))
	if err != nil {
		t.Errorf("Expected Put to return a non-nil error")
	}
	storage.Clear()
	size := storage.GetSize()
	if size != 0 {
		t.Errorf("Expected Clear to remove all items")
	}
}

func TestStorage_Random(t *testing.T) {
	config := Config{
		Path:     "storage",
		Timeout:  10 * time.Second,
	}
	storage := NewStorage(config)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		err := storage.Put(key, []byte(value))
		if err != nil {
			t.Errorf("Expected Put to return a non-nil error")
		}
	}
	size := storage.GetSize()
	if size != 100 {
		t.Errorf("Expected GetSize to return the correct size")
	}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		value, err := storage.Get(key)
		if err != nil {
			t.Errorf("Expected Get to return a non-nil error")
		}
		if string(value) != fmt.Sprintf("value%d", i) {
			t.Errorf("Expected Get to return the correct value")
		}
	}
}
