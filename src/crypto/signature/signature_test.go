package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestSignature_Sign(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("Expected GenerateKey to return a non-nil error")
	}
	signature := NewSignature(privateKey)
	data := []byte("Hello, World!")
	signed, err := signature.Sign(data)
	if err != nil {
		t.Errorf("Expected Sign to return a non-nil error")
	}
	if len(signed) != 256 {
		t.Errorf("Expected signed data to be 256 bytes long")
	}
}

func TestSignature_SignBase64(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("Expected GenerateKey to return a non-nil error")
	}
	signature := NewSignature(privateKey)
	data := []byte("Hello, World!")
	signed, err := signature.SignBase64(data)
	if err != nil {
		t.Errorf("Expected SignBase64 to return a non-nil error")
	}
	if len(signed) != 344 {
		t.Errorf("Expected base64-encoded signed data to be 344 characters long")
	}
}

func TestSignature_Verify(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("Expected GenerateKey to return a non-nil error")
	}
	signature := NewSignature(privateKey)
	data := []byte("Hello, World!")
	signed, err := signature.Sign(data)
	if err != nil {
		t.Errorf("Expected Sign to return a non-nil error")
	}
	if !signature.Verify(data, signed) {
		t.Errorf("Expected Verify to return true")
	}
}

func TestSignature_VerifyBase64(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("Expected GenerateKey to return a non-nil error")
	}
	signature := NewSignature(privateKey)
	data := []byte("Hello, World!")
	signed, err := signature.SignBase64(data)
	if err != nil {
		t.Errorf("Expected SignBase64 to return a non-nil error")
	}
	if !signature.VerifyBase64(data, signed) {
		t.Errorf("Expected VerifyBase64 to return true")
	}
}

func TestSignature_InvalidKey(t *testing.T) {
	_, err := rsa.GenerateKey(rand.Reader, 1024)
	if err == nil {
		t.Errorf("Expected GenerateKey to return a non-nil error for small key size")
	}
}

func TestSignature_Random(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("Expected GenerateKey to return a non-nil error")
	}
	signature := NewSignature(privateKey)
	for i := 0; i < 100; i++ {
		data := fmt.Sprintf("Hello, World! %d", i)
		signed, err := signature.SignBase64([]byte(data))
		if err != nil {
			t.Errorf("Expected SignBase64 to return a non-nil error")
		}
		if !signature.VerifyBase64([]byte(data), signed) {
			t.Errorf("Expected VerifyBase64 to return true")
		}
	}
}
