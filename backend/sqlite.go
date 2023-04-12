package backend

import (
	"errors"
	"time"

	"github.com/riadafridishibly/mypass/config"
	"github.com/riadafridishibly/mypass/models"
	"github.com/riadafridishibly/mypass/vkeys"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

type SqliteBackend struct {
	engine *xorm.Engine
}

// AddPublicKeys implements Backend
func (b *SqliteBackend) AddPublicKeys(pubKeys ...string) error {
	if len(pubKeys) == 0 {
		return nil
	}
	var v []PublicKey
	for _, pubKey := range pubKeys {
		v = append(v, PublicKey{Key: pubKey})
	}
	affected, err := b.engine.Insert(&v)
	if err != nil {
		return err
	}
	if affected != int64(len(pubKeys)) {
		// Some are not inserted!
		return errors.New("some keys are not inserted")
	}
	return nil
}

// RemovePublicKeys implements Backend
func (*SqliteBackend) RemovePublicKeys(pubKeys ...string) error {
	panic("unimplemented")
}

type PublicKey struct {
	ID        int64 `xorm:"pk autoincr 'id'"`
	Key       string
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

// PublicKeys implements Backend
func (b *SqliteBackend) PublicKeys() ([]string, error) {
	var v []PublicKey
	err := b.engine.Find(&v)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(v))
	for i := range out {
		out[i] = v[i].Key
	}
	return out, nil
}

func newSqliteBackend() (Backend, error) {
	var v SqliteBackend
	err := v.Init(&config.Config{
		DatabasePath: viper.GetString(vkeys.DatabasePath),
	})
	return &v, err
}

// CreateItem implements Backend
func (b *SqliteBackend) CreateItem(i *models.Item) (*models.Item, error) {
	affected, err := b.engine.Insert(i)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, errors.New("failed to insert data")
	}
	return b.GetItemByID(i.ID)
}

// Flush implements Backend
func (b *SqliteBackend) Flush() error {
	return b.engine.Close()
}

// GetItemByID implements Backend
func (b *SqliteBackend) GetItemByID(id int) (*models.Item, error) {
	var i models.Item
	i.ID = id
	found, err := b.engine.Get(&i)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("item not found")
	}
	return &i, nil
}

// Init implements Backend
func (b *SqliteBackend) Init(cfg *config.Config) error {
	// FIXME: update the config struct to support different backend
	e, err := xorm.NewEngine("sqlite3", cfg.DatabasePath+".sqlite")
	if err != nil {
		return err
	}
	b.engine = e
	return b.engine.Sync(new(models.Item), new(PublicKey))
}

// ListAllItems implements Backend
func (b *SqliteBackend) ListAllItems() ([]*models.Item, error) {
	var out []*models.Item
	err := b.engine.Find(&out)
	return out, err
}

// RemoveItemByID implements Backend
func (*SqliteBackend) RemoveItemByID(id int) (*models.Item, error) {
	panic("unimplemented")
}

// UpdateItemByID implements Backend
func (*SqliteBackend) UpdateItemByID(id int, i *models.Item) (*models.Item, error) {
	panic("unimplemented")
}

var _ Backend = (*SqliteBackend)(nil)
