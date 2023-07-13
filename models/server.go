package models

type CategoryModel struct {
	ID int `json:"id"`

	CategoryAPIModel

	Mangas []int
}

type MangaModel struct {
	MangaModelServer
	Sync chan string
}

type MangaModelServer struct {
	PathName string `json:"pathName"`

	MangaAPIModel

	CoverURL string `json:"coverUrl"`
	Chapters []*ChapterModel
}

type ChapterModel struct {
	Index int    `json:"index"`
	Path  string `json:"path"`

	ChapterAPIModel
}
