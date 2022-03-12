package libs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type Cryptor struct {
	key []byte
}

func MkCryptor(key string) *Cryptor {
	return &Cryptor{
		key: []byte(key),
	}
}

func (cryptor *Cryptor) Encrypt(plainText []byte) ([]byte, error) {
	c, err := aes.NewCipher(cryptor.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plainText, nil), nil
}

func (cryptor *Cryptor) Decrypt(cipherText []byte) ([]byte, error) {
	c, err := aes.NewCipher(cryptor.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: %v", cipherText)
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	return gcm.Open(nil, nonce, cipherText, nil)
}
