package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sh a256"
	"encoding/base64"
	"errors"
)

type Signature struct {
	privateKey *rsa.PrivateKey
}

func NewSignature(privateKey *rsa.PrivateKey) *Signature {
	return &Signature{privateKey}
}

func (s *Signature) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (s *Signature) SignBase64(data []byte) (string, error) {
	signature, err := s.Sign(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (s *Signature) Verify(data []byte, signature []byte) bool {
	hash := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(&rsa.PublicKey{N: s.privateKey.N, E: s.privateKey.E}, crypto.SHA256, hash[:], signature)
	return err == nil
}

func (s *Signature) VerifyBase64(data []byte, signature string) bool {
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	return s.Verify(data, signatureBytes)
}
