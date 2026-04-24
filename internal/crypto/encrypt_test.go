package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	passphrase := "super-secret-passphrase"
	plaintext := []byte("MY_API_KEY=abc123\nDB_PASSWORD=hunter2")

	ciphertext, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := Decrypt(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptProducesUniqueOutput(t *testing.T) {
	passphrase := "passphrase"
	plaintext := []byte("same plaintext")

	c1, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("first Encrypt failed: %v", err)
	}
	c2, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("second Encrypt failed: %v", err)
	}

	if bytes.Equal(c1, c2) {
		t.Fatal("two encryptions of the same plaintext should differ due to random salt/nonce")
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	ciphertext, err := Encrypt([]byte("secret value"), "correct-passphrase")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(ciphertext, "wrong-passphrase")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestDecryptTruncatedData(t *testing.T) {
	_, err := Decrypt([]byte("short"), "passphrase")
	if err == nil {
		t.Fatal("expected error for truncated data")
	}
}

func TestDeriveKeyDeterministic(t *testing.T) {
	salt := []byte("fixed-salt-value")
	k1 := DeriveKey("mypassphrase", salt)
	k2 := DeriveKey("mypassphrase", salt)
	if !bytes.Equal(k1, k2) {
		t.Fatal("DeriveKey should be deterministic for same inputs")
	}
	if len(k1) != KeySize {
		t.Fatalf("expected key size %d, got %d", KeySize, len(k1))
	}
}
