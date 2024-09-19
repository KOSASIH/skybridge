package encryption

import (
	"testing"
)

func TestEncryption_EncryptDecrypt(t *testing.T) {
	encryption, err := NewEncryption("my_secret_key")
	if err != nil {
		t.Errorf("Expected NewEncryption to return a non-nil error")
	}
	plaintext := "Hello, World!"
	ciphertext, err := encryption.Encrypt([]byte(plaintext))
	if err != nil {
		t.Errorf("Expected Encrypt to return a non-nil error")
	}
	decrypted, err := encryption.Decrypt(ciphertext)
	if err != nil {
		t.Errorf("Expected Decrypt to return a non-nil error")
	}
	if string(decrypted) != plaintext {
		t.Errorf("Expected decrypted text to match original plaintext")
	}
}

func TestEncryption_EncryptBase64DecryptBase64(t *testing.T) {
	encryption, err := NewEncryption("my_secret_key")
	if err != nil {
		t.Errorf("Expected NewEncryption to return a non-nil error")
	}
	plaintext := "Hello, World!"
	ciphertext, err := encryption.EncryptBase64(plaintext)
	if err != nil {
		t.Errorf("Expected EncryptBase64 to return a non-nil error")
	}
	decrypted, err := encryption.DecryptBase64(ciphertext)
	if err != nil {
		t.Errorf("Expected DecryptBase64 to return a non-nil error")
	}
	if decrypted != plaintext {
		t.Errorf("Expected decrypted text to match original plaintext")
	}
}

func TestEncryption_InvalidKey(t *testing.T) {
	_, err := NewEncryption("short")
	if err == nil {
		t.Errorf("Expected NewEncryption to return a non-nil error for short key")
	}
}

func TestEncryption_Random(t *testing.T) {
	encryption, err := NewEncryption("my_secret_key")
	if err != nil {
		t.Errorf("Expected NewEncryption to return a non-nil error")
	}
	for i := 0; i < 100; i++ {
		plaintext := fmt.Sprintf("Hello, World! %d", i)
		ciphertext, err := encryption.EncryptBase64(plaintext)
		if err != nil {
			t.Errorf("Expected EncryptBase64 to return a non-nil error")
		}
		decrypted, err := encryption.DecryptBase64(ciphertext)
		if err != nil {
			t.Errorf("Expected DecryptBase64 to return a non-nil error")
		}
		if decrypted != plaintext {
			t.Errorf("Expected decrypted text to match original plaintext")
		}
	}
}
