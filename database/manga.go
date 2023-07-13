package database

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/IstarVin/manga-reader-go/details"
	"github.com/IstarVin/manga-reader-go/global"
	"github.com/IstarVin/manga-reader-go/models"
	"github.com/IstarVin/manga-reader-go/syncmanager"
	json "github.com/json-iterator/go"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func (m *MangaDatabase) GetIDAll() []int {
	var ids []int
	for _, manga := range m.Database {
		ids = append(ids, manga.ID)
	}
	return ids
}

func (m *MangaDatabase) Update() {
	log.Println("Updating the manga database...")
	println("this may take a while on first run")
	mangas, err := os.ReadDir(global.MangasDirectory)
	if err != nil {
		log.Fatal("Error reading the manga directory")
	}

	for mangaIndex, manga := range mangas {
		syncmanager.SyncManager.AddQueue(func(args ...any) {
			_mangaIndex := args[0].(int)
			_manga := args[1].(os.DirEntry)

			if !_manga.IsDir() && _manga.Type().String()[0] != 'L' {
				return
			}

			var mangaModel models.MangaModel
			var chaptersDetails []*models.ChapterModel

			mangaModel.Title = _manga.Name()

			detailsFile, err := os.ReadFile(filepath.Join(global.MangasDirectory, _manga.Name(), "details.json"))
			if err != nil {
				if os.IsExist(err) {
					log.Fatal("Error reading the manga details file")
				}
				mangaModel.Genre = []string{}

				fetchedMangaDetails, fetchedChaptersDetails := details.GetDetails(_manga.Name())
				if fetchedMangaDetails == nil {
					log.Println("No fetched manga details for", _manga.Name())
				} else {
					mangaModel.MangaAPIModel = fetchedMangaDetails.MangaAPIModel
					mangaModel.CoverURL = fetchedMangaDetails.CoverURL
					chaptersDetails = fetchedChaptersDetails

					log.Println("Fetched manga details for", _manga.Name())

					// Save manga details
					toSaveDetailsFile, err := json.MarshalIndent(mangaModel.MangaModelServer, "", "  ")
					if err != nil {
						log.Fatal("Error marshalling the manga details file:", err)
					}
					err = os.WriteFile(filepath.Join(global.MangasDirectory, _manga.Name(), "details.json"), toSaveDetailsFile, 0644)
					if err != nil {
						log.Fatal("Error writing the manga details file")
					}

					// Save chapters details
					toSaveChaptersDetailsFile, err := json.MarshalIndent(chaptersDetails, "", "  ")
					if err != nil {
						log.Fatal("Error marshalling the manga chapters file")
					}
					err = os.WriteFile(filepath.Join(global.MangasDirectory, _manga.Name(), "chapters.json"), toSaveChaptersDetailsFile, 0644)
					if err != nil {
						log.Fatal("Error writing the manga chapters file")
					}
				}
			} else {
				err = json.Unmarshal(detailsFile, &mangaModel)
				if err != nil {
					log.Fatal("Error unmarshalling the manga details file")
				}
			}

			// Check if cover photo exists
			_, err = os.Stat(filepath.Join(global.MangasDirectory, _manga.Name(), "cover.jpg"))
			if err != nil {
				if os.IsNotExist(err) && mangaModel.CoverURL != "" {
					// Download cover photo
					cover := details.DownloadFile(mangaModel.CoverURL)
					// Save cover photo
					err = os.WriteFile(filepath.Join(global.MangasDirectory, _manga.Name(), "cover.jpg"), *cover, 0644)
					if err != nil {
						log.Println("Error saving the cover photo")
					}
				}
			}

			mangaModel.ID = _mangaIndex
			mangaModel.PathName = _manga.Name()
			mangaModel.Url = fmt.Sprintf("/api/v1/manga/%d", _mangaIndex)
			mangaModel.ThumbnailUrl = fmt.Sprintf("/api/v1/manga/%d/thumbnail", _mangaIndex)

			chapterDetailsFile, err := os.ReadFile(filepath.Join(global.MangasDirectory, _manga.Name(), "chapters.json"))
			if err != nil {
				if os.IsExist(err) {
					log.Fatal("Error reading the manga chapters file")
				}
			} else {
				err = json.Unmarshal(chapterDetailsFile, &chaptersDetails)
				if err != nil {
					log.Fatal("Error unmarshalling the manga chapters file")
				}
			}

			chapters, err := os.ReadDir(filepath.Join(global.MangasDirectory, _manga.Name()))
			if err != nil {
				log.Fatal("Error reading the manga directory")
			}

			var chapterIndex int

			// Sort the chapters in descending order
			sort.Slice(chapters, func(i, j int) bool {
				if filepath.Ext(chapters[i].Name()) != ".cbz" {
					return true
				}
				if filepath.Ext(chapters[j].Name()) != ".cbz" {
					return false
				}

				iNumber, err1 := parseChapterNumber(chapters[i].Name())
				jNumber, err2 := parseChapterNumber(chapters[j].Name())

				if err1 != nil || err2 != nil {
					return false
				}

				return iNumber > jNumber
			})

			for _, chapter := range chapters {
				// Skip if chapter is not cbz
				if filepath.Ext(chapter.Name()) != ".cbz" {
					continue
				}

				chapterModel := &models.ChapterModel{}

				if len(chaptersDetails) > chapterIndex {
					for _, detail := range chaptersDetails {
						chapterNumber, err1 := parseChapterNumber(chapter.Name())
						chapterDetailNumber, err2 := parseChapterNumber(detail.Name)
						if err1 != nil || err2 != nil {
							log.Println("Error parsing the chapter number")
							chapterModel.Name = strings.TrimSuffix(chapter.Name(), filepath.Ext(chapter.Name()))
							break
						}

						if chapterNumber == chapterDetailNumber {
							chapterModel = detail
							break
						}
					}
				} else {
					chapterModel.Name = strings.TrimSuffix(chapter.Name(), filepath.Ext(chapter.Name()))
				}

				chapterModel.Index = chapterIndex
				chapterModel.Path = chapter.Name()
				chapterModel.Url = fmt.Sprintf("/api/v1/manga/%d/chapter/%d", mangaIndex, chapterIndex)

				chapterCBZ, err := zip.OpenReader(filepath.Join(global.MangasDirectory, _manga.Name(), chapter.Name()))
				if err != nil {
					log.Fatal("Error opening the chapter cbz file", err)
				}

				chapterModel.PageCount = len(chapterCBZ.File)

				mangaModel.Chapters = append(mangaModel.Chapters, chapterModel)
				chapterIndex++
			}

			m.Database = append(m.Database, &mangaModel)
		}, mangaIndex, manga)
	}

	syncmanager.SyncManager.WaitFinish()
}

func (m *MangaDatabase) Save() {
	mangasJSON, err := json.MarshalIndent(m.Database, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling the manga database")
	}

	err = os.WriteFile(global.MangaDatabasePath, mangasJSON, 0644)
	if err != nil {
		log.Fatal("Error writing the manga database")
	}
}

func parseChapterNumber(chapterName string) (int, error) {
	numberRegex, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Fatal("Error compiling the number regex")
		return 0, err
	}

	chapterNumberString := numberRegex.FindString(chapterName)
	if chapterNumberString == "" {
		log.Println("No number found in chapter:", chapterName)
		return 0, errors.New("no number found in chapter")
	}

	chapterNumber, err := strconv.Atoi(numberRegex.FindString(chapterName))
	if err != nil {
		log.Fatal("Error converting the number regex: ", chapterName)
		return 0, err
	}

	return chapterNumber, nil
}
