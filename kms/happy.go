package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/scrypt"
)

type Key struct {
	Version string
	Salt    []byte
	Algo    string
	Key     [32]byte
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func generateKey(passphrase string, salt []byte) ([32]byte, error) {
	keyBytes, err := scrypt.Key([]byte(passphrase), salt, 32768, 8, 1, 32)
	if err != nil {
		return [32]byte{}, err
	}
	var key [32]byte
	copy(key[:], keyBytes)
	return key, nil
}

func rotateKey(passphrase string, store map[string]Key) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	key, err := generateKey(passphrase, salt)
	if err != nil {
		return "", err
	}

	version := time.Now().Format(time.RFC3339)
	store[version] = Key{
		Version: version,
		Salt:    salt,
		Algo:    "scrypt",
		Key:     key,
	}

	return version, nil
}

func generateGCM(key []byte) (gcm cipher.AEAD, block cipher.Block, err error) {
	block, err = aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	return gcm, block, nil
}

func encrypt(data []byte, key []byte) ([]byte, error) {
	gcm, _, err := generateGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
	gcm, _, err := generateGCM(key)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func reEncrypt(data []byte, oldKey []byte, newKey []byte) ([]byte, error) {
	plaintext, err := decrypt(data, oldKey)
	if err != nil {
		return nil, err
	}

	return encrypt(plaintext, newKey)
}

// StreamEncrypt encrypts data using AES in GCM mode with streaming
func StreamEncrypt(data io.Reader, key []byte) (io.Reader, error) {
	gcm, block, err := generateGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	return cipher.StreamReader{
		S: cipher.NewCTR(block, nonce),
		R: data,
	}, nil
}

// StreamDecrypt decrypts data using AES in GCM mode with streaming
func StreamDecrypt(data io.Reader, key []byte) (io.Reader, error) {
	gcm, block, err := generateGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(data, nonce); err != nil {
		return nil, err
	}

	return cipher.StreamReader{
		S: cipher.NewCTR(block, nonce),
		R: data,
	}, nil
}

func testStreamEncryptionRotation(passphrase string, keyStore map[string]Key) {

	version, err := rotateKey(passphrase, keyStore)
	if err != nil {
		log.Fatalf("Failed to rotate key: %v", err)
	}

	oldKey, exists := keyStore[version]
	if !exists {
		log.Fatalf("Key not found for version: %s", version)
	}

	// Replace this with a reader for your large data (e.g., file, network stream)
	data, err := os.ReadFile("./test.txt")
	if err != nil {
		log.Fatalf("Failed to read data: %v", err)
	}
	dataReader := bytes.NewReader(data)

	// Encrypt data with the old key
	encryptedStream, err := StreamEncrypt(dataReader, oldKey.Key[:])
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}

	// Read the encrypted data from the stream (for demonstration)
	encryptedData, err := io.ReadAll(encryptedStream)
	if err != nil {
		log.Fatalf("Failed to read encrypted data: %v", err)
	}

	newVersion, err := rotateKey(passphrase, keyStore)
	if err != nil {
		log.Fatalf("Failed to rotate key: %v", err)
	}

	newKey, exists := keyStore[newVersion]
	if !exists {
		log.Fatalf("Key not found for version: %s", newVersion)
	}

	// Decrypt the data with the new key
	decryptedStream, err := StreamDecrypt(bytes.NewReader(encryptedData), newKey.Key[:])
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}

	// Read the decrypted data from the stream
	decryptedData, err := io.ReadAll(decryptedStream)
	if err != nil {
		log.Fatalf("Failed to read decrypted data: %v", err)
	}

	log.Printf("Original data: %s", data)
	log.Printf("Decrypted data: %s", decryptedData)
}

func testEncryptionRotation(passphrase string, keyStore map[string]Key) {
	version, err := rotateKey(passphrase, keyStore)
	if err != nil {
		log.Fatalf("Failed to rotate key: %v", err)
	}

	oldKey, exists := keyStore[version]
	if !exists {
		log.Fatalf("Key not found for version: %s", version)
	}

	data := []byte("Sensitive data here")
	encryptedData, err := encrypt(data, oldKey.Key[:])
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}

	newVersion, err := rotateKey(passphrase, keyStore)
	if err != nil {
		log.Fatalf("Failed to rotate key: %v", err)
	}

	newKey, exists := keyStore[newVersion]
	if !exists {
		log.Fatalf("Key not found for version: %s", newVersion)
	}

	// Re-encrypt data with new key
	encryptedData, err = reEncrypt(encryptedData, oldKey.Key[:], newKey.Key[:])
	if err != nil {
		log.Fatalf("Failed to re-encrypt data: %v", err)
	}

	decryptedData, err := decrypt(encryptedData, newKey.Key[:])
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}

	log.Printf("Original data: %s", data)
	log.Printf("Decrypted data: %s", decryptedData)
}
