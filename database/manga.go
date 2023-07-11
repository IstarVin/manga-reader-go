package database

import (
	"archive/zip"
	"fmt"
	"github.com/IstarVin/manga-reader-go/global"
	"github.com/IstarVin/manga-reader-go/models"
	json "github.com/json-iterator/go"
	"log"
	"os"
	"path/filepath"
)

var MangaDB MangaDatabase

type MangaDatabase struct {
	Database []*models.MangaModel
}

func (m *MangaDatabase) FindMangaWithID(id int) *models.MangaModel {
	for _, manga := range m.Database {
		if manga.ID == id {
			return manga
		}
	}
	return nil
}

func (m *MangaDatabase) FindChapterWithIndex(index int, manga *models.MangaModel) *models.ChapterModel {
	for _, chapter := range manga.Chapters {
		if chapter.Index == index {
			return chapter
		}
	}
	return nil
}

func (m *MangaDatabase) Update() {
	mangas, err := os.ReadDir(global.MangasDirectory)
	if err != nil {
		log.Fatal("Error reading the manga directory")
	}

	for mangaIndex, manga := range mangas {
		if !manga.IsDir() && manga.Type().String()[0] != 'L' {
			continue
		}

		var mangaModel models.MangaModel

		detailsFile, err := os.ReadFile(filepath.Join(global.MangasDirectory, manga.Name(), "details.json"))
		if err != nil {
			if os.IsExist(err) {
				log.Fatal("Error reading the manga details file")
			}
			mangaModel.Title = manga.Name()
		} else {
			err = json.Unmarshal(detailsFile, &mangaModel)
			if err != nil {
				log.Fatal("Error unmarshalling the manga details")
			}
		}

		mangaModel.ID = mangaIndex
		mangaModel.Url = fmt.Sprintf("/api/v1/manga/%d", mangaIndex)
		mangaModel.ThumbnailUrl = fmt.Sprintf("/api/v1/manga/%d/thumbnail", mangaIndex)

		var chaptersDetails []models.ChapterModel

		chapterDetailsFile, err := os.ReadFile(filepath.Join(global.MangasDirectory, manga.Name(), "chapters.json"))
		if err != nil {
			if os.IsExist(err) {
				log.Fatal("Error reading the manga chapters file")
			}
		} else {
			err = json.Unmarshal(chapterDetailsFile, &chaptersDetails)
			if err != nil {
				log.Fatal("Error unmarshalling the manga chapters")
			}
		}

		chapters, err := os.ReadDir(filepath.Join(global.MangasDirectory, manga.Name()))
		if err != nil {
			return
		}

		var chapterIndex int
		for _, chapter := range chapters {
			// Skip if chapter is not cbz
			if filepath.Ext(chapter.Name()) != ".cbz" {
				continue
			}

			var chapterModel models.ChapterModel

			if len(chaptersDetails) > chapterIndex {
				chapterModel = chaptersDetails[chapterIndex]
			} else {
				chapterModel.Name = chapter.Name()
			}

			chapterModel.Index = chapterIndex
			chapterModel.Url = fmt.Sprintf("/api/v1/manga/%d/chapter/%d", mangaIndex, chapterIndex)

			chapterCBZ, err := zip.OpenReader(filepath.Join(global.MangasDirectory, manga.Name(), chapter.Name()))
			if err != nil {
				log.Fatal("Error opening the chapter cbz file")
			}

			chapterModel.PageCount = len(chapterCBZ.File)

			mangaModel.Chapters = append(mangaModel.Chapters, &chapterModel)

			chapterIndex++
		}

		m.Database = append(m.Database, &mangaModel)
	}
}
