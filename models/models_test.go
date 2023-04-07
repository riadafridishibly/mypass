package models

import (
	"encoding/json"
	"testing"

	"filippo.io/age"
	"github.com/spf13/viper"
)

func TestAsymSecretStrJSON(t *testing.T) {
	i, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal("Failed to create age x25519 identity:", err)
	}
	s := AsymSecretStr("hello world")
	viper.Set("public_keys", []string{i.Recipient().String()})
	viper.Set("private_keys", []string{i.String()})
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Failed to marshal AsymSecretStr:", err)
	}
	t.Log("Marshalled: ", string(data))
	var v AsymSecretStr
	err = json.Unmarshal(data, &v)
	if err != nil {
		t.Fatal("Failed to unmarshal AsymSecretStr:", err)
	}
	t.Log("Unmarshalled:", string(v))
	if v != s {
		t.Fatal("Unmarshalled and original data are not same")
	}
}

func TestSymSecretStrJSON(t *testing.T) {
	s := SymSecretStr("hello world")
	viper.Set("password", "pa$$w0rd")
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Failed to marshal SymSecretStr:", err)
	}
	t.Log("Marshalled: ", string(data))
	var v SymSecretStr
	err = json.Unmarshal(data, &v)
	if err != nil {
		t.Fatal("Failed to unmarshal SymSecretStr:", err)
	}
	t.Log("Unmarshalled:", string(v))
	if v != s {
		t.Fatal("Unmarshalled and original data are not same")
	}
}
