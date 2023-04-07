package encryption

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"filippo.io/age"
	"github.com/spf13/viper"
)

func pubKeys2recipients(pubKeys ...string) ([]age.Recipient, error) {
	var recipients []age.Recipient
	for _, pubKey := range pubKeys {
		recipient, err := age.ParseX25519Recipient(pubKey)
		if err != nil {
			if viper.GetBool("verbose") {
				log.Printf("Failed to parse public key: %s, err: %v", pubKey, err)
			}
			return nil, fmt.Errorf("failed to parse public key %s: %w", pubKey, err)
		}
		recipients = append(recipients, recipient)
	}
	return recipients, nil
}

func Encrypt(plaintext []byte, pubKeys ...string) ([]byte, error) {
	recipients, err := pubKeys2recipients(pubKeys...)
	if err != nil {
		return nil, err
	}
	return encrypt(plaintext, recipients...)
}

func encrypt(plaintext []byte, r ...age.Recipient) ([]byte, error) {
	out := &bytes.Buffer{}
	w, err := age.Encrypt(out, r...)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}
	if _, err := w.Write(plaintext); err != nil {
		return nil, fmt.Errorf("failed to write encrypt data: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}
	return out.Bytes(), nil
}

func privKeys2identities(privKeys ...string) ([]age.Identity, error) {
	var identities []age.Identity
	for _, privKey := range privKeys {
		identity, err := age.ParseX25519Identity(privKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		identities = append(identities, identity)
	}
	return identities, nil
}

func Decrypt(ciphertext []byte, privKeys ...string) ([]byte, error) {
	identities, err := privKeys2identities(privKeys...)
	if err != nil {
		return nil, err
	}
	return decrypt(ciphertext, identities...)
}

func decrypt(ciphertext []byte, i ...age.Identity) ([]byte, error) {
	r, err := age.Decrypt(bytes.NewReader(ciphertext), i...)
	if err != nil {
		return nil, fmt.Errorf("failed to open encrypted data: %w", err)
	}
	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, fmt.Errorf("failed to read encrypted data: %w", err)
	}
	return out.Bytes(), nil
}

func EncryptWithPassword(plaintext []byte, password string) ([]byte, error) {
	r, err := age.NewScryptRecipient(password)
	if err != nil {
		return nil, err
	}
	return encrypt(plaintext, r)
}

func DecryptWithPassword(ciphertext []byte, password string) ([]byte, error) {
	i, err := age.NewScryptIdentity(password)
	if err != nil {
		return nil, err
	}
	return decrypt(ciphertext, i)
}
