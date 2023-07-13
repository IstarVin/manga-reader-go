package details

import (
	"github.com/IstarVin/manga-reader-go/models"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"strings"
	"time"
)

func Mangakakalot(name string) (*models.MangaModel, []*models.ChapterModel) {
	var manga *models.MangaModel
	var chapters []*models.ChapterModel

	var scrapeArgs ScrapeArgs

	scrapeArgs.URL = "https://mangakakalot.com/search/story/" + strings.ReplaceAll(name, " ", "_")
	scrapeArgs.MangaSelector = "div.story_item:nth-child(1) .story_name > a"
	scrapeArgs.MangaParser = func(e *colly.HTMLElement) {
		manga = &models.MangaModel{}

		resp, err := http.Head(e.Attr("href"))
		if err != nil {
			log.Println(err)
		}

		err = e.Request.Visit(resp.Request.URL.String())
		if err != nil {
			log.Println(err)
		}
	}

	scrapeArgs.NameSelector = ".story-info-right > h1"
	scrapeArgs.CoverSelector = ".info-image > img"
	scrapeArgs.DescriptionSelector = "#panel-story-info-description"

	scrapeArgs.DetailsSelector = ".variations-tableInfo > tbody"

	scrapeArgs.ChaptersSelector = ".row-content-chapter"

	scrapeArgs.NameParser = func(e *colly.HTMLElement) {
		manga.Title = strings.TrimSpace(e.Text)
	}
	scrapeArgs.CoverParser = func(e *colly.HTMLElement) {
		manga.CoverURL = strings.TrimSpace(e.Attr("src"))
	}
	scrapeArgs.DescriptionParser = func(e *colly.HTMLElement) {
		manga.Description = strings.ReplaceAll(e.Text, "Description :", "")
		manga.Description = strings.TrimSpace(manga.Description)
	}

	scrapeArgs.DetailsParser = func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			detailName := element.ChildAttr("td > i", "class")
			switch detailName {
			case "info-author":
				manga.Author = element.ChildText("td > a")
			case "info-status":
				manga.Status = element.ChildText("td.table-value")
			case "info-genres":
				element.ForEach("td > a", func(i int, genreElement *colly.HTMLElement) {
					manga.Genre = append(manga.Genre, genreElement.Text)
				})
			}
		})
	}

	scrapeArgs.ChaptersParser = func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, element *colly.HTMLElement) {
			var chapter models.ChapterAPIModel

			chapter.Name = element.ChildText("a")
			uploadDate := element.ChildAttr("span.chapter-time", "title")

			parseDate, err := time.Parse("Jan 2,2006 15:04", uploadDate)
			if err != nil {
				log.Println(err)
				parseDate = time.Now()
			}

			chapter.UploadDate = parseDate.UnixMilli()

			chapters = append(chapters, &models.ChapterModel{ChapterAPIModel: chapter, Index: i})
		})
	}

	Scrape(&scrapeArgs).Wait()

	return manga, chapters
}
