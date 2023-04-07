package encryption

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptWithPassword(t *testing.T) {
	plaintext := []byte("the quick brown!! fox jumped@@ over lazzzyy dog!")
	password := "123@@passw04rd"
	ciphertext, err := EncryptWithPassword(plaintext, password)
	if err != nil {
		t.Fatal("Failed to encrypt:", err)
	}
	decryptedText, err := DecryptWithPassword(ciphertext, password)
	if err != nil {
		t.Fatal("Failed to decrypt:", err)
	}
	if !bytes.Equal(plaintext, decryptedText) {
		t.Fatal("Decrypted text is not same as origianl text!")
	}
	_, err = DecryptWithPassword(ciphertext, password[1:])
	if err == nil {
		t.Fatal("Expected error, but found nil")
	}
}
