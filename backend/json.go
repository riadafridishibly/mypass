package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/riadafridishibly/mypass/config"
	"github.com/riadafridishibly/mypass/models"
	"github.com/riadafridishibly/mypass/vkeys"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type JSONBackend struct {
	file string
	db   *models.Database
}

// AddPublicKeys implements Backend
func (*JSONBackend) AddPublicKeys(pubKeys ...string) error {
	panic("unimplemented")
}

// RemovePublicKeys implements Backend
func (*JSONBackend) RemovePublicKeys(pubKeys ...string) error {
	panic("unimplemented")
}

// PublicKeys implements Backend
func (jb *JSONBackend) PublicKeys() ([]string, error) {
	return jb.db.PublicKeys, nil
}

func newJSONBackend() (Backend, error) {
	var b JSONBackend
	err := b.Init(&config.Config{
		DatabasePath: viper.GetString(vkeys.DatabasePath),
	})
	return &b, err
}

// Flush implements Backend
func (jb *JSONBackend) Flush() error {
	return jb.save()
}

// CreateItem implements Backend
func (jb *JSONBackend) CreateItem(i *models.Item) (*models.Item, error) {
	return jb.db.AddItem(i)
}

// GetItemByID implements Backend
func (jb *JSONBackend) GetItemByID(id int) (*models.Item, error) {
	return jb.db.FindItemByID(id)
}

// Init implements Backend
func (jb *JSONBackend) Init(cfg *config.Config) error {
	jb.file = cfg.DatabasePath
	currentData, err := os.ReadFile(cfg.DatabasePath)
	if err != nil {
		jww.ERROR.Println("Failed to open database file")
		return err
	}
	// FIXME: we may lose data if we fail to write again
	// Create a backup for now
	dir, _ := filepath.Split(cfg.DatabasePath)
	backupDb := filepath.Join(dir, fmt.Sprintf("db-%d", time.Now().Unix()))
	err = os.WriteFile(backupDb, currentData, 0600)
	if err != nil {
		return err
	}
	jb.db, err = config.LoadDatabase()
	if err != nil {
		return err
	}
	return err
}

// ListAllItems implements Backend
func (jb *JSONBackend) ListAllItems() ([]*models.Item, error) {
	return jb.db.Items, nil
}

// RemoveItemByID implements Backend
func (jb *JSONBackend) RemoveItemByID(id int) (*models.Item, error) {
	return jb.db.RemoveItem(id)
}

// UpdateItemByID implements Backend
func (jb *JSONBackend) UpdateItemByID(id int, i *models.Item) (*models.Item, error) {
	return jb.db.UpdateItem(id, i)
}

var _ Backend = (*JSONBackend)(nil)

func (jb *JSONBackend) save() error {
	data, err := json.Marshal(jb.db)
	if err != nil {
		return err
	}
	// maybe rewind here?
	return os.WriteFile(jb.file, data, 0600)
}
