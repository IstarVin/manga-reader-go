package server

import (
	"github.com/IstarVin/manga-reader-go/controller"
	"github.com/IstarVin/manga-reader-go/database"
	"github.com/IstarVin/manga-reader-go/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func mangaAPI(api *gin.RouterGroup) {
	manga := api.Group("/manga/:mangaID", mangaIDValidator)
	{
		manga.GET("", controller.SendMangaDetails)
		manga.GET("/thumbnail", controller.SendMangaThumbnail)
		manga.GET("/chapters", controller.SendChapterList)
		manga.GET("/chapter/:chapterIndex", chapterIndexValidator, controller.SendChapterDetails)
		manga.GET("/chapter/:chapterIndex/page/:pageIndex", chapterIndexValidator, controller.SendPage)
	}
}

func mangaIDValidator(c *gin.Context) {
	mangaID, err := strconv.Atoi(c.Param("mangaID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid manga ID",
		})
		return
	}
	mangaObj := database.MangaDB.FindMangaWithID(mangaID)
	if mangaObj == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Manga not found",
		})
		return
	}

	c.Set("manga", mangaObj)

	c.Next()
}

func chapterIndexValidator(c *gin.Context) {
	chapterIndex, err := strconv.Atoi(c.Param("chapterIndex"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid chapter index",
		})
		return
	}

	mangaObj := c.MustGet("manga").(*models.MangaModel)
	chapterObj := database.MangaDB.FindChapterWithIndex(chapterIndex, mangaObj)
	if chapterObj == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Chapter not found",
		})
		return
	}

	c.Set("chapter", chapterObj)

	c.Next()
}
