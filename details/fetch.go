package details

import (
	"github.com/IstarVin/manga-reader-go/models"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
)

type ScrapeArgs struct {
	URL string

	MangaSelector string

	DetailsSelector string

	NameSelector        string
	CoverSelector       string
	AuthorSelector      string
	ArtistSelector      string
	DescriptionSelector string
	GenreSelector       string
	StatusSelector      string
	ChaptersSelector    string

	MangaParser colly.HTMLCallback

	DetailsParser colly.HTMLCallback

	NameParser        colly.HTMLCallback
	CoverParser       colly.HTMLCallback
	AuthorParser      colly.HTMLCallback
	ArtistParser      colly.HTMLCallback
	DescriptionParser colly.HTMLCallback
	GenreParser       colly.HTMLCallback
	StatusParser      colly.HTMLCallback
	ChaptersParser    colly.HTMLCallback
}

var AvailableSources = []func(string) (*models.MangaModel, []*models.ChapterModel){
	Mangakakalot,
}

func GetDetails(name string) (*models.MangaModel, []*models.ChapterModel) {
	name = filter(name)

	var manga *models.MangaModel
	var chapters []*models.ChapterModel

	for _, source := range AvailableSources {
		manga, chapters = source(name)
		if manga != nil && chapters != nil {
			return manga, chapters
		}
	}

	return nil, nil
}

func Scrape(args *ScrapeArgs) *colly.Collector {
	collector := colly.NewCollector()

	// Manga Search
	collector.OnHTML(args.MangaSelector, args.MangaParser)

	// Details
	collector.OnHTML(args.DetailsSelector, args.DetailsParser)

	// Cover
	collector.OnHTML(args.CoverSelector, args.CoverParser)

	// Name
	collector.OnHTML(args.NameSelector, args.NameParser)

	// Author
	collector.OnHTML(args.AuthorSelector, args.AuthorParser)

	// Artist
	collector.OnHTML(args.ArtistSelector, args.ArtistParser)

	// Description
	collector.OnHTML(args.DescriptionSelector, args.DescriptionParser)

	// Genre
	collector.OnHTML(args.GenreSelector, args.GenreParser)

	// Status
	collector.OnHTML(args.StatusSelector, args.StatusParser)

	// Chapters
	collector.OnHTML(args.ChaptersSelector, args.ChaptersParser)

	err := collector.Visit(args.URL)
	if err != nil {
		log.Fatalln(err)
	}

	return collector
}

func DownloadFile(url string) *[]byte {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	// Write the body to file
	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	err = resp.Body.Close()
	if err != nil {
		return nil
	}
	return &file
}
