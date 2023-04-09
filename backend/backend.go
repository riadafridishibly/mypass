package backend

import "github.com/riadafridishibly/mypass/models"

type Backend interface {
	CreateItem(i *models.Item) (*models.Item, error)
	ListAllItems() ([]*models.Item, error)
	GetItemByID(id int) (*models.Item, error)
	UpdateItemByID(id int, i *models.Item) (*models.Item, error)
	RemoveItemByID(id int) (*models.Item, error)
}
