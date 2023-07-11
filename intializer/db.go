package intializer

import (
	"github.com/IstarVin/manga-reader-go/database"
	"github.com/IstarVin/manga-reader-go/global"
	json "github.com/json-iterator/go"
	"log"
	"os"
)

func LoadDatabase() {
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
