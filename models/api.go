package models

type CategoryAPIModel struct {
	Id      int    `json:"id"`
	Order   int    `json:"order"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type MangaAPIModel struct {
	Url          string   `json:"url"`
	ThumbnailUrl string   `json:"thumbnailUrl"`
	Title        string   `json:"title"`
	Artist       string   `json:"artist"`
	Author       string   `json:"author"`
	Description  string   `json:"description"`
	Genre        []string `json:"genre"`
	Status       string   `json:"status"`
}

type ChapterAPIModel struct {
	Url           string  `json:"url"`
	Name          string  `json:"name"`
	UploadDate    int64   `json:"uploadDate"`
	ChapterNumber float32 `json:"chapterNumber"`
	Scanlator     string  `json:"scanlator"`
	PageCount     int     `json:"pageCount"`
}
