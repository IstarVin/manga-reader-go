package server

import (
	"github.com/IstarVin/manga-reader-go/controller"
	"github.com/IstarVin/manga-reader-go/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func categoryAPI(api *gin.RouterGroup) {
	category := api.Group("/categories")
	{
		category.GET("/", controller.SendCategoryList)
	}
}

func categoryIDValidator(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}
	categoryObj := database.CategoryDB.FindCategoryWithID(categoryID)
	if categoryObj == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	c.Set("category", categoryObj)

	c.Next()
}
