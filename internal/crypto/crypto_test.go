package crypto_test

import (
	"bytes"
	"testing"

	"github.com/cargohold/cargohold/internal/crypto"
)

func TestDeriveKey(t *testing.T) {
	key := crypto.DeriveKey("mysecretpassphrase")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}

	key2 := crypto.DeriveKey("mysecretpassphrase")
	if !bytes.Equal(key, key2) {
		t.Fatal("expected same passphrase to produce same key")
	}

	key3 := crypto.DeriveKey("differentpassphrase")
	if bytes.Equal(key, key3) {
		t.Fatal("expected different passphrases to produce different keys")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := crypto.DeriveKey("testpassphrase")
	plaintext := []byte("hello, cargohold!")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should differ from plaintext")
	}

	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key := crypto.DeriveKey("correctpassphrase")
	wrongKey := crypto.DeriveKey("wrongpassphrase")
	plaintext := []byte("sensitive data")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = crypto.Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong key")
	}
}

func TestDecryptCorruptedData(t *testing.T) {
	key := crypto.DeriveKey("testpassphrase")
	_, err := crypto.Decrypt(key, []byte("tooshort"))
	if err == nil {
		t.Fatal("expected error for corrupted/short ciphertext")
	}
}

func TestEncryptProducesUniqueCiphertexts(t *testing.T) {
	key := crypto.DeriveKey("testpassphrase")
	plaintext := []byte("same plaintext")

	c1, _ := crypto.Encrypt(key, plaintext)
	c2, _ := crypto.Encrypt(key, plaintext)

	if bytes.Equal(c1, c2) {
		t.Fatal("expected unique ciphertexts due to random nonce")
	}
}
