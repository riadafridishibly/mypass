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
	backendType   = "backend-type"
	BackendJSON   = "backend-json"
	BackendSqlite = "backend-sqlite"
)

func New() (Backend, error) {
	if viper.GetString(backendType) == BackendJSON {
		return newJSONBackend()
	}
	return newSqliteBackend()
}
