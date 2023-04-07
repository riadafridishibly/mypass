package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"filippo.io/age"
	"github.com/riadafridishibly/mypass/config"
	"github.com/riadafridishibly/mypass/models"
	"github.com/spf13/viper"
)

func JSONTestPrivateKeys() {
	viper.Set("password", "passw0rd")
	privKeys := models.PrivateKeys{
		Meta: models.Meta{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now().Add(24 * time.Hour),
		},
	}
	for i := 0; i < 5; i++ {
		i, err := age.GenerateX25519Identity()
		if err != nil {
			continue
		}
		privKeys.Keys = append(privKeys.Keys, models.SymSecretStr(i.String()))
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "   ")
	err := e.Encode(privKeys)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	i, err := age.GenerateX25519Identity()
	if err != nil {
		log.Fatal("Failed to create age x25519 identity:", err)
	}
	viper.Set("public_keys", []string{i.Recipient().String()})
}

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("Err:", err)
	}
	fmt.Println("Password:", viper.GetString("password"))
	fmt.Println("Private Keys:", viper.GetStringSlice("private_keys"))
	var db models.Database
	db.PublicKeys = viper.GetStringSlice("public_keys")
	err = db.AddItem("work", &models.Item{
		ID: 7,
		Title: "Hello World",
		Password: &models.PasswordItem{
			Username: "jon",
			SiteName: "code.google.com",
			URL:      "https://code.google.com/login",
			Password: "12345598",
		},
	})
	if err != nil {
		log.Fatal("Err:", err)
	}
	err = db.AddItem("work", &models.Item{
		ID: 7,
		Title: "Hello World",
		Password: &models.PasswordItem{
			Username: "jon",
			SiteName: "code.google.com",
			URL:      "https://code.google.com/login",
			Password: "12345598",
		},
	})
	if err != nil {
		log.Fatal("Err:", err)
	}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "   ")
	err = e.Encode(db)
	if err != nil {
		log.Fatal("Err:", err)
	}
	// defer profile.Start().Stop()
	// JSONTestPrivateKeys()
	// viper.AddConfigPath(".")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	log.Fatal("Failed to read config file:", err)
	// }
	// viper.Set("Verbose", true)

	// m := models.Namespaces{
	// 	"work": &models.Namespce{
	// 		Items: map[string]*models.Item{
	// 			"P-134": {
	// 				ID:    "P-123",
	// 				Title: "Prod 01 Server Fr",
	// 				Type:  models.ItemSSH,
	// 				SSH: &models.SSHItem{
	// 					Host:     "prod01.example.com",
	// 					Port:     "4428",
	// 					Username: "work",
	// 					Password: []byte("asdfxyg"),
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// db := models.Database{
	// 	PublicKeys: []string{
	// 		"PUB-KEY-1",
	// 		"PUB-KEY-2",
	// 	},
	// 	Namespaces: m,
	// }

	// data, err := json.MarshalIndent(db, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(data))

	// var mm models.Database
	// err = json.Unmarshal(data, &mm)
	// if err != nil {
	// 	panic(err)
	// }
	// pp.Println(mm)
}
