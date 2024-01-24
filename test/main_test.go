package main_test

import (
	"bytes"
	"crypto/rand"
	"fileserver/internal/fcrypt"
	"log"
	"os"
	"testing"
)

// Generate a random AES key for testing
func generateRandomAESKey(t *testing.T) []byte {
	key := make([]byte, 32) // AES-256
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}
	log.Printf("Random key: %x", key)
	return key
}

func TestEncryptDecrypt(t *testing.T) {
	originalContent := []byte("Hello, World!")

	// Create a temporary file for testing
	srcFile, err := os.CreateTemp("", "test-encrypt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(srcFile.Name()) // Clean up

	// Write original content to source file
	if _, err := srcFile.Write(originalContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := srcFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Re-open the file for reading
	srcFile, err = os.Open(srcFile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer srcFile.Close()

	// Create a temporary file for encrypted data
	encryptedFile, err := os.CreateTemp("", "test-encrypted")
	if err != nil {
		t.Fatalf("Failed to create temp file for encrypted data: %v", err)
	}
	defer os.Remove(encryptedFile.Name()) // Clean up

	key := generateRandomAESKey(t)

	// Encrypt the file
	if err := fcrypt.EncryptWithGCM(srcFile, encryptedFile, key); err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	if err := encryptedFile.Close(); err != nil {
		t.Fatalf("Failed to close encrypted temp file: %v", err)
	}

	// Decrypt the file
	decryptedFile, err := os.CreateTemp("", "test-decrypted")
	if err != nil {
		t.Fatalf("Failed to create temp file for decrypted data: %v", err)
	}
	defer os.Remove(decryptedFile.Name()) // Clean up

	if err := fcrypt.DecryptWithGCM(encryptedFile.Name(), decryptedFile.Name(), key); err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	if err := decryptedFile.Close(); err != nil {
		t.Fatalf("Failed to close decrypted temp file: %v", err)
	}

	// Read the decrypted content
	decryptedContent, err := os.ReadFile(decryptedFile.Name())
	if err != nil {
		t.Fatalf("Failed to read decrypted file: %v", err)
	}

	// Compare original and decrypted content
	if !bytes.Equal(originalContent, decryptedContent) {
		t.Errorf("Original and decrypted content do not match")
	}
}

// Include your EncryptWithGCM and DecryptWithGCM functions here
