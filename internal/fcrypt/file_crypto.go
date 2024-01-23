package fcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"mime/multipart"
	"os"

	"golang.org/x/crypto/blake2b"
)

func HashWithBlake2(file multipart.File) ([]byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func EncryptWithGCM(file multipart.File, dst *os.File, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Write the nonce to the file (it will be needed for decryption)
	if _, err = dst.Write(nonce); err != nil {
		return err
	}

	// Read the file content and encrypt it
	buffer := make([]byte, 4096) // Adjust the buffer size according to your needs
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n > 0 {
			encrypted := gcm.Seal(nil, nonce, buffer[:n], nil)
			if _, err = dst.Write(encrypted); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}

func DecryptWithGCM(srcPath string, dstPath string, key []byte) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// The block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// The GCM mode AEAD
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Read the nonce from the file
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(srcFile, nonce); err != nil {
		return err
	}

	// Create a buffer to read the encrypted data
	buffer := make([]byte, 4096) // This should match the size used during encryption
	for {
		n, err := srcFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n > 0 {
			decrypted, err := gcm.Open(nil, nonce, buffer[:n], nil)
			if err != nil {
				return err
			}
			if _, err = dstFile.Write(decrypted); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}
