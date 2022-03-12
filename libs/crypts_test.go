package libs

import (
	"fmt"
	"testing"
)

func TestCrypt(t *testing.T) {
	key := "11111111111111111111111111111111"
	cryptor := MkCryptor(key)
	plaintext := "123123"
	ciphertext, err := cryptor.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatal("Unable to encrypt plaintext, error: ", err)
	}
	fmt.Printf("ciphertext: %x\n", ciphertext)
	cleartext, err := cryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatal("Unable to decrypt ciphertext, error: ", err)
	}
	fmt.Printf("cleatext: %s\n", string(cleartext))
}
