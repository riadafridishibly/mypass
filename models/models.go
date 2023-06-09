package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
	"unicode"

	"github.com/riadafridishibly/mypass/encryption"
	"github.com/riadafridishibly/mypass/vkeys"
	"github.com/spf13/viper"

	"xorm.io/xorm/convert"
)

// TODO: these functionalities will be replaced by backend

var (
	ErrItemNotFound = errors.New("item not found")
)

type AsymSecretStr string

var (
	_ json.Marshaler   = (*AsymSecretStr)(nil)
	_ json.Unmarshaler = (*AsymSecretStr)(nil)
)

func (asc AsymSecretStr) MarshalJSON() ([]byte, error) {
	keys := viper.GetStringSlice(vkeys.PublicKeys)
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
	keys := viper.GetStringSlice(vkeys.PrivateKeys)
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
	password := viper.GetString(vkeys.Password)
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
	password := viper.GetString(vkeys.Password)
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
	PublicKeys []string `json:"public_keys,omitempty"`
	Items      []*Item  `json:"items,omitempty"`
}

func (db *Database) AddItem(i *Item) (*Item, error) {
	if i.Namespace == "" {
		return i, errors.New("namespace can't be empty")
	}
	if i.Title == "" {
		return i, errors.New("title can't be empty")
	}
	if i.ID == 0 {
		i.ID = len(db.Items) + 1
	}
	i.Meta.CreatedAt = time.Now()
	i.Meta.UpdatedAt = time.Now()
	db.Items = append(db.Items, i)
	return i, nil
}

// TODO: Change Item to some struct with pointer to detect
// which fields to update
func (db *Database) UpdateItem(id int, i *Item) (*Item, error) {
	panic("not implemented")
}

func (db *Database) FindItemByID(id int) (*Item, error) {
	for _, i := range db.Items {
		if i.ID == id {
			return i, nil
		}
	}
	return nil, fmt.Errorf("%w: id=%d", ErrItemNotFound, id)
}

func (db *Database) RemoveItem(id int) (*Item, error) {
	for idx, i := range db.Items {
		if i.ID == id {
			// Maybe prompt before deleting?
			removed := db.Items[idx]
			db.Items = append(db.Items[:idx], db.Items[idx+1:]...)
			return removed, nil
		}
	}
	return nil, fmt.Errorf("%w: id=%d", ErrItemNotFound, id)
}

func (db *Database) Namespaces() []string {
	l := make(map[string]struct{}, len(db.Items))
	uniq := make([]string, 0, len(db.Items))
	for _, i := range db.Items {
		if _, ok := l[i.Namespace]; ok {
			continue
		}
		uniq = append(uniq, i.Namespace)
		l[i.Namespace] = struct{}{}
	}
	sort.Strings(uniq)
	return uniq
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
	ID        int           `xorm:"pk autoincr 'id'" json:"id,string,omitempty"`
	Title     string        `json:"title,omitempty"`
	Namespace string        `json:"namespace,omitempty"`
	Type      ItemType      `json:"type,omitempty"`
	Meta      Meta          `json:"meta,omitempty" xorm:"-"`
	Password  *PasswordItem `xorm:"text 'password'" json:"password,omitempty"`
	SSH       *SSHItem      `xorm:"text 'ssh'" json:"ssh,omitempty"`
}

func (i *Item) InnerItemString() string {
	if i.Password != nil {
		return i.Password.String()
	}
	if i.SSH != nil {
		return i.SSH.String()
	}
	panic("internal error: both are null")
}

func (i *Item) GetPassword() (string, error) {
	if i == nil {
		return "", errors.New("item is nil")
	}
	if i.Password != nil {
		return string(i.Password.Password), nil
	}
	if i.SSH != nil {
		return string(i.SSH.Password), nil
	}
	return "", errors.New("not a password or ssh item")
}

func EllipticalTruncate(text string, maxLen int) string {
	lastSpaceIx := -1
	len := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len >= maxLen {
			if lastSpaceIx != -1 {
				return text[:lastSpaceIx] + "..."
			}
			// If here, string is longer than max, but has no spaces
		}
	}
	// If here, string is shorter than max
	return text
}

func (i Item) String() string {
	common := fmt.Sprintf("%-3d %-28q ", i.ID, EllipticalTruncate(i.Title, 25))
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

// FromDB implements convert.Conversion
func (p *PasswordItem) FromDB(data []byte) error {
	var v PasswordItem
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*p = v
	return nil
}

// ToDB implements convert.Conversion
func (p *PasswordItem) ToDB() ([]byte, error) {
	return json.Marshal(p)
}

var _ convert.Conversion = (*PasswordItem)(nil)

func (p PasswordItem) String() string {
	return fmt.Sprintf("user=%s site=%s", p.Username, p.SiteName)
}

type SSHItem struct {
	Host     string        `json:"host,omitempty"`
	Port     uint16        `json:"port,omitempty"`
	Username string        `json:"username,omitempty"`
	Password AsymSecretStr `json:"password,omitempty"`
}

// FromDB implements convert.Conversion
func (s *SSHItem) FromDB(data []byte) error {
	var v SSHItem
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// ToDB implements convert.Conversion
func (s *SSHItem) ToDB() ([]byte, error) {
	return json.Marshal(s)
}

var _ convert.Conversion = (*SSHItem)(nil)

func (p SSHItem) String() string {
	return fmt.Sprintf("ssh -p %d %s@%s", p.Port, p.Username, p.Host)
}
