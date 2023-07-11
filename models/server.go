package models

type CategoryModel struct {
	ID int `json:"id"`

	CategoryAPIModel

	Mangas []*MangaModel
}

type MangaModel struct {
	ID int `json:"id"`

	MangaAPIModel

	Chapters []*ChapterModel
}

type ChapterModel struct {
	Index int `json:"index"`

	ChapterAPIModel
}
