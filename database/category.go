package database

import (
	"github.com/IstarVin/manga-reader-go/models"
)

var defaultCategory = models.CategoryAPIModel{
	Id:      0,
	Order:   1,
	Name:    "All",
	Default: true,
}

var CategoryDB CategoryDatabase = CategoryDatabase{
	Database: []*models.CategoryModel{
		{CategoryAPIModel: defaultCategory, Mangas: MangaDB.Database},
	},
}

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
