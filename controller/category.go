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

	c.JSON(http.StatusOK, categoryList)
}

func SendCategoryMangas(c *gin.Context) {
	category := c.MustGet("category").(*models.CategoryModel)

	var mangaList []models.MangaAPIModel
	for _, mangaID := range category.Mangas {
		mangaList = append(mangaList, database.MangaDB.FindMangaWithID(mangaID).MangaAPIModel)
	}

	c.JSON(http.StatusOK, mangaList)
}
