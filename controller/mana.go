package controller

import (
	"archive/zip"
	"fmt"
	"github.com/IstarVin/manga-reader-go/global"
	"github.com/IstarVin/manga-reader-go/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func SendMangaDetails(c *gin.Context) {
	manga := c.MustGet("manga").(*models.MangaModel)

	c.JSON(http.StatusOK, manga.MangaAPIModel)
}

func SendMangaThumbnail(c *gin.Context) {
	manga := c.MustGet("manga").(*models.MangaModel)

	thumbnail, err := os.ReadFile(filepath.Join(global.MangasDirectory, manga.PathName, "cover.jpg"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Thumbnail not found",
		})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", thumbnail)
}

func SendChapterList(c *gin.Context) {
	manga := c.MustGet("manga").(*models.MangaModel)
	var chaptersAPI []models.ChapterAPIModel
	for _, chapter := range manga.Chapters {
		chapter.Url = fmt.Sprintf("/api/v1/manga/%d/chapter/%d", manga.ID, chapter.Index)
		chaptersAPI = append(chaptersAPI, chapter.ChapterAPIModel)
	}

	c.JSON(http.StatusOK, chaptersAPI)
}

func SendChapterDetails(c *gin.Context) {
	chapter := c.MustGet("chapter").(*models.ChapterModel)

	c.JSON(http.StatusOK, chapter.ChapterAPIModel)
}

func SendPage(c *gin.Context) {
	manga := c.MustGet("manga").(*models.MangaModel)
	pageIndex, err := strconv.Atoi(c.Param("pageIndex"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid page index",
		})
		return
	}

	chapter := c.MustGet("chapter").(*models.ChapterModel)

	page, err := zip.OpenReader(filepath.Join(global.MangasDirectory, manga.PathName, chapter.Path))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open chapter file",
		})
		return
	}

	if pageIndex < 0 || pageIndex >= len(page.File) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid page index",
		})
		return
	}

	pageFile, err := page.File[pageIndex].Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open page file",
		})
		return
	}

	pageImg, err := io.ReadAll(pageFile)
	err = page.Close()
	if err != nil {
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read page file",
		})
		return
	}

	var contentType string
	switch filepath.Ext(page.File[pageIndex].Name) {
	case ".jpg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".webp":
		contentType = "image/webp"
	default:
		contentType = "image/jpeg"
	}

	c.Data(http.StatusOK, contentType, pageImg)
}
