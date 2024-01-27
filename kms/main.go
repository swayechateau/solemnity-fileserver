package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	// Setting up the logger to write to standard output
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// func main() {
// 	keyStore := make(map[string]Key)
// 	passphrase := "your-secure-passphrase"

// 	testEncryptionRotation(passphrase, keyStore)
// 	testStreamEncryptionRotation(passphrase, keyStore)
// }

func encryptChunk(block cipher.Block, plaintext []byte, nonce []byte) ([]byte, error) {
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Println("Failed to create new GCM:", err)
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, plaintext, nil), nil
}

func encryptLargeDataToFile(data io.Reader, key []byte, chunkSize int, filePath string) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Println("Failed to create new cipher:", err)
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		logger.Println("Failed to create file:", err)
		return err
	}
	defer file.Close()

	chunk := make([]byte, chunkSize)
	nonce := make([]byte, 12) // 12 bytes nonce for GCM

	for {
		n, err := data.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Println("Failed to read data:", err)
			return err
		}

		if _, err := rand.Read(nonce); err != nil {
			logger.Println("Failed to read nonce:", err)
			return err
		}

		encryptedChunk, err := encryptChunk(block, chunk[:n], nonce)
		if err != nil {
			logger.Println("Failed to encrypt chunk:", err)
			return err
		}

		// Write nonce and encrypted chunk to file
		if _, err := file.Write(nonce); err != nil {
			logger.Println("Failed to write nonce:", err)
			return err
		}
		if _, err := file.Write(encryptedChunk); err != nil {
			logger.Println("Failed to write encrypted chunk:", err)
			return err
		}
	}

	return nil
}

func decryptChunk(block cipher.Block, encryptedChunk []byte, nonce []byte) ([]byte, error) {
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Println("Failed to create new GCM:", err)
		return nil, err
	}

	return aesgcm.Open(nil, nonce, encryptedChunk, nil)
}

func decryptFileToFile(encryptedFilePath, decryptedFilePath string, key []byte, chunkSize int) error {
	encryptedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		logger.Println("Failed to open encrypted file:", err)
		return err
	}
	defer encryptedFile.Close()

	decryptedFile, err := os.Create(decryptedFilePath)
	if err != nil {
		logger.Println("Failed to create decrypted file:", err)
		return err
	}
	defer decryptedFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Println("Failed to create new cipher:", err)
		return err
	}

	nonceSize := 12 // 12 bytes nonce for GCM
	chunk := make([]byte, chunkSize+nonceSize)

	for {
		n, err := encryptedFile.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Println("Failed to read encrypted file:", err)
			return err
		}

		nonce := chunk[:nonceSize]
		encryptedChunk := chunk[nonceSize:n]

		decryptedChunk, err := decryptChunk(block, encryptedChunk, nonce)
		if err != nil {
			logger.Println("Failed to decrypt chunk:", err)
			return err
		}

		if _, err := decryptedFile.Write(decryptedChunk); err != nil {
			logger.Println("Failed to write decrypted chunk:", err)
			return err
		}
	}

	return nil
}

func main() {
	// Open the test.txt file
	inputFile, err := os.Open("test.jpeg")
	if err != nil {
		logger.Fatalf("Failed to open input file: %v", err)
	}
	defer inputFile.Close()

	// Generate a key (32 bytes for AES-256)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		logger.Fatal(err)
	}

	// File path where encrypted data will be saved
	outputFilePath := "encrypted_data.bin"

	// Encrypt data and save to file
	if err := encryptLargeDataToFile(inputFile, key, 1024*1024, outputFilePath); err != nil {
		logger.Fatalf("Failed to encrypt and save data: %v", err)
	}

	log.Println("Encrypted data saved to:", outputFilePath)
	// decrypt data

	decryptedFilePath := "decrypted_test.jpeg"

	// Decrypt the file
	if err := decryptFileToFile(outputFilePath, decryptedFilePath, key, 1024*1024); err != nil {
		logger.Fatalf("Failed to decrypt file: %v", err)
	}

	logger.Println("Decrypted data saved to:", decryptedFilePath)
}
