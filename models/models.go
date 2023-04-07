package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/riadafridishibly/mypass/encryption"
	"github.com/spf13/viper"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type AsymSecretStr string

var (
	_ json.Marshaler   = (*AsymSecretStr)(nil)
	_ json.Unmarshaler = (*AsymSecretStr)(nil)
)

func (asc AsymSecretStr) MarshalJSON() ([]byte, error) {
	keys := viper.GetStringSlice("public_keys")
	data, err := encryption.Encrypt([]byte(asc), keys...)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func (asc *AsymSecretStr) UnmarshalJSON(data []byte) error {
	var b []byte
	err := json.Unmarshal(data, &b)
	if err != nil {
		return err
	}
	keys := viper.GetStringSlice("private_keys")
	data, err = encryption.Decrypt(b, keys...)
	if err != nil {
		return err
	}
	*asc = AsymSecretStr(data)
	return nil
}

type SymSecretStr string

var (
	_ json.Marshaler   = (*SymSecretStr)(nil)
	_ json.Unmarshaler = (*SymSecretStr)(nil)
)

func (asc SymSecretStr) MarshalJSON() ([]byte, error) {
	password := viper.GetString("password")
	data, err := encryption.EncryptWithPassword([]byte(asc), password)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func (asc *SymSecretStr) UnmarshalJSON(data []byte) error {
	var b []byte
	err := json.Unmarshal(data, &b)
	if err != nil {
		return err
	}
	password := viper.GetString("password")
	data, err = encryption.DecryptWithPassword(b, password)
	if err != nil {
		return err
	}
	*asc = SymSecretStr(data)
	return nil
}

type ItemType string

const (
	ItemPassword ItemType = "password"
	ItemSSH      ItemType = "ssh"
)

type Meta struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type PrivateKeys struct {
	Meta Meta           `json:"meta,omitempty"`
	Keys []SymSecretStr `json:"keys,omitempty"`
}

type Database struct {
	PublicKeys []string   `json:"public_keys,omitempty"`
	Namespaces Namespaces `json:"namespaces,omitempty"`
}

func (db *Database) AddItem(namespace string, i *Item) error {
	if namespace == "" {
		return errors.New("namespace can't be empty")
	}
	if i.Title == "" {
		return errors.New("title can't be empty")
	}
	if db.Namespaces == nil {
		db.Namespaces = Namespaces{}
	}
	ns := db.Namespaces[namespace]
	if ns == nil {
		db.Namespaces[namespace] = &Namespace{}
		ns = db.Namespaces[namespace]
	}
	if i.ID == 0 {
		i.ID = len(ns.Items) + 1
	}
	ns.Items = append(ns.Items, i)
	return nil
}

// TODO: Change Item to some struct with pointer to detect
// which fields to update
func (db *Database) UpdateItem(namespace string, id int, i *Item) error {
	panic("not implemented")
}

func (db *Database) FindItemByID(namespace string, id int) (*Item, error) {
	n, ok := db.Namespaces[namespace]
	if !ok {
		return nil, fmt.Errorf("%w: namespace %q not found", ErrItemNotFound, namespace)
	}
	for _, i := range n.Items {
		if i.ID == id {
			return i, nil
		}
	}
	return nil, fmt.Errorf("%w: id=%d namespce %v", ErrItemNotFound, id, namespace)
}

func (db *Database) RemoveItem(namespace string, id int) error {
	n, ok := db.Namespaces[namespace]
	if !ok {
		return fmt.Errorf("namespace %q not found", namespace)
	}
	for idx, i := range n.Items {
		if i.ID == id {
			// Maybe prompt before deleting?
			n.Items = append(n.Items[:idx], n.Items[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("item ID(%d) not found in namespce %q", id, namespace)
}

type Namespaces map[string]*Namespace

func (n Namespaces) Keys() []string {
	if n == nil {
		return nil
	}
	l := make([]string, 0, len(n))
	for k := range n {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}

type Namespace struct {
	Meta  Meta    `json:"meta,omitempty"`
	Items []*Item `json:"items,omitempty"`
}

func (n *Namespace) reID() {
	for idx, i := range n.Items {
		i.ID = idx + 1
	}
}

func (n *Namespace) MarshalJSON() ([]byte, error) {
	n.reID()
	// To avoid recursive calling MarshalJSON
	type nsAlias Namespace
	var v nsAlias = nsAlias(*n)
	return json.Marshal(v)
}

type Item struct {
	ID       int           `json:"id,string,omitempty"`
	Title    string        `json:"title,omitempty"`
	Type     ItemType      `json:"type,omitempty"`
	Password *PasswordItem `json:"password,omitempty"`
	SSH      *SSHItem      `json:"ssh,omitempty"`
}

func (i Item) String() string {
	common := fmt.Sprintf("id=%d title=%q ", i.ID, i.Title)
	var args []any
	args = append(args, common)
	if i.Password != nil {
		args = append(args, i.Password)
	}
	if i.SSH != nil {
		args = append(args, i.SSH)
	}
	return fmt.Sprint(args...)
}

type PasswordItem struct {
	Username string        `json:"username,omitempty"`
	SiteName string        `json:"site_name,omitempty"`
	URL      string        `json:"url,omitempty"`
	Password AsymSecretStr `json:"password,omitempty"`
}

func (p PasswordItem) String() string {
	return fmt.Sprintf("username=%s site=%s", p.Username, p.SiteName)
}

type SSHItem struct {
	Host     string        `json:"host,omitempty"`
	Port     uint16        `json:"port,omitempty"`
	Username string        `json:"username,omitempty"`
	Password AsymSecretStr `json:"password,omitempty"`
}

func (p SSHItem) String() string {
	return fmt.Sprintf("ssh -p %d %s@%s", p.Port, p.Username, p.Host)
}
