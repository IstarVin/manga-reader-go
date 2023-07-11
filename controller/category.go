package controller

import (
	"github.com/IstarVin/manga-reader-go/database"
	"github.com/IstarVin/manga-reader-go/models"
	"github.com/gin-gonic/gin"
)

func SendCategoryList(c *gin.Context) {
	var categoryList []models.CategoryAPIModel

	for _, category := range database.CategoryDB.Database {
		categoryList = append(categoryList, category.CategoryAPIModel)
	}

	c.JSON(200, categoryList)
}
