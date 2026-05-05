package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize   = 16
	KeySize    = 32
	Iterations = 100_000
)

// DeriveKey derives a 256-bit AES key from a passphrase and salt using PBKDF2.
func DeriveKey(passphrase string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, Iterations, KeySize, sha256.New)
}

// Encrypt encrypts plaintext using AES-256-GCM with a key derived from passphrase.
// Output format: [salt (16 bytes)] [nonce (12 bytes)] [ciphertext+tag]
func Encrypt(plaintext []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := DeriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	result := make([]byte, 0, SaltSize+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)
	return result, nil
}

// Decrypt decrypts data previously encrypted with Encrypt.
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	if len(data) < SaltSize+12 {
		return nil, errors.New("crypto: data too short")
	}

	salt := data[:SaltSize]
	key := DeriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < SaltSize+nonceSize {
		return nil, errors.New("crypto: data too short for nonce")
	}

	nonce := data[SaltSize : SaltSize+nonceSize]
	ciphertext := data[SaltSize+nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("crypto: decryption failed — wrong passphrase or corrupted data")
	}
	return plaintext, nil
}

// ReEncrypt decrypts data with oldPassphrase and re-encrypts it with newPassphrase.
// This is useful for passphrase rotation without exposing the plaintext to the caller.
func ReEncrypt(data []byte, oldPassphrase, newPassphrase string) ([]byte, error) {
	plaintext, err := Decrypt(data, oldPassphrase)
	if err != nil {
		return nil, err
	}
	return Encrypt(plaintext, newPassphrase)
}
