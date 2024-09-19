package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type Hash struct {
	key []byte
}

func NewHash(key string) (*Hash, error) {
	if len(key) < 16 {
		return nil, errors.New("key must be at least 16 characters long")
	}
	return &Hash{[]byte(key)}, nil
}

func (h *Hash) Hash(data []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, h.key)
	_, err := mac.Write(data)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func (h *Hash) HashBase64(data []byte) (string, error) {
	hash, err := h.Hash(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

func (h *Hash) Verify(data []byte, hash []byte) bool {
	mac := hmac.New(sha256.New, h.key)
	_, err := mac.Write(data)
	if err != nil {
		return false
	}
	expected := mac.Sum(nil)
	return hmac.Equal(expected, hash)
}

func (h *Hash) VerifyBase64(data []byte, hash string) bool {
	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}
	return h.Verify(data, hashBytes)
}
