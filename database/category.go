package database

import (
	"github.com/IstarVin/manga-reader-go/global"
	"github.com/IstarVin/manga-reader-go/models"
	json "github.com/json-iterator/go"
	"log"
	"os"
)

var DefaultCategory = models.CategoryAPIModel{
	Id:      0,
	Order:   1,
	Name:    "All",
	Default: true,
}

var CategoryDB CategoryDatabase

type CategoryDatabase struct {
	Database []*models.CategoryModel `json:"database"`
}

func (c *CategoryDatabase) FindCategoryWithID(id int) *models.CategoryModel {
	for _, category := range c.Database {
		if category.ID == id {
			return category
		}
	}
	return nil
}

func (c *CategoryDatabase) Save() {
	categoryDBFile, err := json.Marshal(c.Database)
	if err != nil {
		log.Fatal("Error marshalling the category database file")
	}

	err = os.WriteFile(global.CategoryDatabasePath, categoryDBFile, 0644)
	if err != nil {
		log.Fatal("Error writing the category database file")
	}
}
