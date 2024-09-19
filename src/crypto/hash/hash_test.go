package hash

import (
	"testing"
)

func TestHash_Hash(t *testing.T) {
	hash, err := NewHash("my_secret_key")
	if err != nil {
		t.Errorf("Expected NewHash to return a non-nil error")
	}
	data := []byte("Hello, World!")
	hashed, err := hash.Hash(data)
	if err != nil {
		t.Errorf("Expected Hash to return a non-nil error")
	}
	if len(hashed) != 32 {
		t.Errorf("Expected hashed data to be 32 bytes long")
	}
}

func TestHash_HashBase64(t *testing.T) {
	hash, err := NewHash("my_secret_key")
	if err != nil {
		t.Errorf("Expected NewHash to return a non-nil error")
	}
	data := []byte("Hello, World!")
	hashed, err := hash.HashBase64(data)
	if err != nil {
		t.Errorf("Expected HashBase64 to return a non-nil error")
	}
	if len(hashed) != 44 {
		t.Errorf("Expected base64-encoded hashed data to be 44 characters long")
	}
}

func TestHash_Verify(t *testing.T) {
	hash, err := NewHash("my_secret_key")
	if err != nil {
		t.Errorf("Expected NewHash to return a non-nil error")
	}
	data := []byte("Hello, World!")
	hashed, err := hash.Hash(data)
	if err != nil {
		t.Errorf("Expected Hash to return a non-nil error")
	}
	if !hash.Verify(data, hashed) {
		t.Errorf("Expected Verify to return true")
	}
}

func TestHash_VerifyBase64(t *testing.T) {
	hash, err := NewHash("my_secret_key")
	if err != nil {
		t.Errorf("Expected NewHash to return a non-nil error")
	}
	data := []byte("Hello, World!")
	hashed, err := hash.HashBase64(data)
	if err != nil {
		t.Errorf("Expected HashBase64 to return a non-nil error")
	}
	if !hash.VerifyBase64(data, hashed) {
		t.Errorf("Expected VerifyBase64 to return true")
	}
}

func TestHash_InvalidKey(t *testing.T) {
	_, err := NewHash("short")
	if err == nil {
		t.Errorf("Expected NewHash to return a non-nil error for short key")
	}
}

func TestHash_Random(t *testing.T) {
	hash, err := NewHash("my_secret_key")
	if err != nil {
		t.Errorf("Expected NewHash to return a non-nil error")
	}
	for i := 0; i < 100; i++ {
		data := fmt.Sprintf("Hello, World! %d", i)
		hashed, err := hash.HashBase64([]byte(data))
		if err != nil {
			t.Errorf("Expected HashBase64 to return a non-nil error")
		}
		if !hash.VerifyBase64([]byte(data), hashed) {
			t.Errorf("Expected VerifyBase64 to return true")
		}
	}
}
