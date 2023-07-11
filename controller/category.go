package controller

import (
	"github.com/IstarVin/manga-reader-go/database"
	"github.com/IstarVin/manga-reader-go/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendCategoryList(c *gin.Context) {
	var categoryList []models.CategoryAPIModel

	for _, category := range database.CategoryDB.Database {
		categoryList = append(categoryList, category.CategoryAPIModel)
	}

	c.JSON(200, categoryList)
}

func SendCategoryMangas(c *gin.Context) {
	category := c.MustGet("category").(*models.CategoryModel)

	var mangaList []models.MangaAPIModel
	for i, mangas := range category.Mangas {
		if i == 0 {
			continue
		}
		mangaList = append(mangaList, mangas.MangaAPIModel)
	}

	c.JSON(http.StatusOK, mangaList)
}
