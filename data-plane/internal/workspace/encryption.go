package workspace

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func Encrypt(data []byte) ([]byte, error) {
	gcm, err := gcmFromEnv()
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

func Decrypt(data []byte) ([]byte, error) {
	gcm, err := gcmFromEnv()
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func gcmFromEnv() (cipher.AEAD, error) {
	keyB64 := os.Getenv("WORKSPACE_KEY_B64")
	if keyB64 == "" {
		return nil, errors.New("WORKSPACE_KEY_B64 not set")
	}
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("WORKSPACE_KEY_B64 must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
