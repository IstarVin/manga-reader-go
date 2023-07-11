package intializer

import (
	"github.com/IstarVin/manga-reader-go/database"
	"github.com/IstarVin/manga-reader-go/global"
	"github.com/IstarVin/manga-reader-go/models"
	json "github.com/json-iterator/go"
	"log"
	"os"
)

func LoadDatabase() {
	loadMangaDatabase()
	loadCategoryDatabase()
}

func loadMangaDatabase() {
	mangaDBFile, err := os.ReadFile(global.MangaDatabasePath)
	if err != nil {
		if os.IsExist(err) {
			log.Fatal("Error reading the manga database file", err)
		}
		database.MangaDB.Update()
		mangaDBFile, err = json.Marshal(database.MangaDB.Database)
		if err != nil {
			log.Fatal("Error marshalling the manga database file")
		}
		err = os.WriteFile(global.MangaDatabasePath, mangaDBFile, 0644)
		if err != nil {
			log.Fatal("Error writing the manga database file")
		}

	} else {
		err = json.Unmarshal(mangaDBFile, &database.MangaDB.Database)
		if err != nil {
			log.Fatal("Error unmarshalling the manga database file")
		}
	}
}

func loadCategoryDatabase() {
	categoryDBFile, err := os.ReadFile(global.CategoryDatabasePath)
	if err != nil {
		if os.IsExist(err) {
			log.Fatal("Error reading the category database file")
		}

		database.CategoryDB = database.CategoryDatabase{
			Database: []*models.CategoryModel{
				{CategoryAPIModel: database.DefaultCategory, Mangas: database.MangaDB.Database},
			},
		}

		categoryDBFile, err = json.Marshal(database.CategoryDB.Database)
		if err != nil {
			log.Fatal("Error marshalling the category database file")
		}
		err = os.WriteFile(global.CategoryDatabasePath, categoryDBFile, 0644)
		if err != nil {
			log.Fatal("Error writing the category database file")
		}
	} else {
		err = json.Unmarshal(categoryDBFile, &database.CategoryDB.Database)
		if err != nil {
			log.Fatal("Error unmarshalling the category database file")
		}
	}
}
