package app

import (
	"fmt"

	"github.com/riadafridishibly/mypass/config"
	"github.com/riadafridishibly/mypass/models"
)

type App struct {
	DB *models.Database
}

func (a *App) AddItem(i *models.Item) error {
	return a.DB.AddItem(i)
}

func (a *App) ListAllItems() []*models.Item {
	return a.DB.Items
}

func (a *App) ListNamespaces() []string {
	return a.DB.Namespaces()
}

func (a *App) ListAllItemsStrings() []string {
	var l []string
	for _, i := range a.DB.Items {
		l = append(l, fmt.Sprintf("ns=%s %s", i.Namespace, i.String()))
	}
	return l
}

func NewApp() (*App, error) {
	db, err := config.LoadDatabase()
	if err != nil {
		return nil, err
	}
	return &App{
		DB: db,
	}, nil
}
