package routes

import (
	"log"
	"net/http"
	"strconv"
	"web-scraper-api/models"

	"github.com/gin-gonic/gin"
)

func GetArticlesHandler(context *gin.Context) {
	articles, err := models.GetAllArticles()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch articles. Try again later."})
		log.Panic(err)
		return
	}
	context.JSON(http.StatusOK, articles)
}

func GetArticlesInPageHandler(context *gin.Context) {
    // Default values
    const pageSize = 10
    page := -1 // Default page number

    // Parse page from query parameter
    pageStr, _ := context.GetQuery("page")
    if pageStr != "" {
        parsedPage, err := strconv.Atoi(pageStr)
        if err == nil && parsedPage > 0 {
            page = parsedPage
        } else {
            context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page number"})
            return
        }
    }

    // Calculate offset based on page
    offset := (page - 1) * pageSize

    // Fetch articles using the limit and offset
    articles, err := models.GetArticlesByPage(pageSize, offset)
    if err != nil {
        log.Printf("Error fetching articles: %v", err)
        context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

    // Return the articles as JSON
    context.JSON(http.StatusOK, articles)
    log.Printf("Success: Api runs successfully")
}



func createArticle(context *gin.Context) {
	var article models.Article
    
	err := context.ShouldBindJSON(&article)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		log.Panic(err)
		return
	}

	articleId, errSave := models.SaveArticle(article)
    article.ID = articleId

	if errSave != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create article. Try again later."})
		log.Panic(errSave)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Article created!", "article": article})
}