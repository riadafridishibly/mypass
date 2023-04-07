package app

import (
	"errors"
	"fmt"

	"github.com/riadafridishibly/mypass/config"
	"github.com/riadafridishibly/mypass/models"
)

type App struct {
	db *models.Database
}

func (a *App) LoadConfigs() error {
	return nil
}

func (a *App) ListNamespaces() ([]string, error) {
	if a.db == nil {
		return nil, errors.New("database not initialized. (did you initialize app?)")
	}
	return a.db.Namespaces.Keys(), nil
}

func (a *App) ListAllItems() []string {
	var l []string
	for k, ns := range a.db.Namespaces {
		for _, i := range ns.Items {
			l = append(l, fmt.Sprintf("ns=%s %s", k, i.String()))
		}
	}
	return l
}

func NewApp() (*App, error) {
	config.Init()
	db, err := config.LoadDatabase()
	if err != nil {
		return nil, err
	}
	return &App{
		db: db,
	}, nil
}
