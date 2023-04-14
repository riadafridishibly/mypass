package backend

import (
	"github.com/riadafridishibly/mypass/config"
	"github.com/riadafridishibly/mypass/models"
	"github.com/spf13/viper"
)

type Backend interface {
	// Will be called to initialize the instance
	Init(cfg *config.Config) error

	CreateItem(i *models.Item) (*models.Item, error)
	ListAllItems() ([]*models.Item, error)
	GetItemByID(id int) (*models.Item, error)
	UpdateItemByID(id int, i *models.Item) (*models.Item, error)
	RemoveItemByID(id int) (*models.Item, error)

	PublicKeys() ([]string, error)
	AddPublicKeys(pubKeys ...string) error
	RemovePublicKeys(pubKeys ...string) error

	Flush() error
}

const (
	backendType   = "backend"
	BackendJSON   = "json"
	BackendSqlite = "sqlite3"
)

func Get() (Backend, error) {
	if b := viper.Get("__backend_object"); b != nil {
		return b.(Backend), nil
	}
	bknd, err := func() (Backend, error) {
		if viper.GetString(backendType) == BackendJSON {
			return newJSONBackend()
		}
		return newSqliteBackend()
	}()
	if err != nil {
		return nil, err
	}
	viper.Set("__backend_object", bknd)
	return bknd, nil
}
