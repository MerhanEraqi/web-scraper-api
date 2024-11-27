package models

import (
	"fmt"
	"time"
	"web-scraper-api/db"
)

// Article represents a news article with a title, link, and timestamp.
type Article struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Link      string    `json:"link"`
    Timestamp time.Time `json:"timestamp"`
}

// Check if an article exists by link
func ArticleExists(link string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS (SELECT 1 FROM articles WHERE link=$1)"
    err := db.PostgresDB.QueryRow(query, link).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check if article exists: %w", err)
    }
    return exists, nil
}

// saveArticle saves a single article to the database
func SaveArticle(article Article) (int, error) {
    var articleID int

    query := `INSERT INTO articles (title, link, timestamp) VALUES ($1, $2, $3) RETURNING id`
    err := db.PostgresDB.QueryRow(query, article.Title, article.Link, article.Timestamp).Scan(&articleID)
    if err != nil {
        return 0, fmt.Errorf("failed to insert article: %w", err)
    }

    return articleID, nil
}

// GetAllArticles retrieves all articles from the database.
func GetAllArticles() ([]Article, error) {
    query := `SELECT id, title, link, timestamp FROM articles ORDER BY timestamp DESC;`
    rows, err := db.PostgresDB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve articles: %w", err)
    }
    defer rows.Close()

    var articles []Article

    for rows.Next() {
        var article Article
        err := rows.Scan(&article.ID, &article.Title, &article.Link, &article.Timestamp)

        if err != nil {
            return nil, fmt.Errorf("failed to scan article row: %w", err)
        }
        articles = append(articles, article)
    }

    return articles, nil
}

type PaginatedArticles struct {
    Articles        []Article `json:"articles"`
    RemainingCount  int       `json:"remaining_count"`
    RemainingPages  int       `json:"remaining_pages"`
    CurrentPage     int       `json:"current_page"`
    TotalPages      int       `json:"total_pages"`
}

// GetArticlesByPage retrieves a paginated list of articles and returns the articles,
// along with the remaining article count, remaining pages, current page, and total pages.
func GetArticlesByPage(limit, offset int) (PaginatedArticles, error) {
    var paginatedArticles PaginatedArticles

    // Query to get paginated articles and total count using COUNT(*) OVER()
    query := `
        SELECT id, title, link, timestamp, COUNT(*) OVER() AS total_count
        FROM articles
        ORDER BY timestamp DESC
        LIMIT $1 OFFSET $2;
    `
    
    rows, err := db.PostgresDB.Query(query, limit, offset)
    if err != nil {
        return paginatedArticles, fmt.Errorf("failed to retrieve paginated articles: %w", err)
    }
    defer rows.Close()

    var articles []Article
    var totalCount int // We will get the total count from the first row

    for rows.Next() {
        var article Article
        var count int

        err := rows.Scan(&article.ID, &article.Title, &article.Link, &article.Timestamp, &count)

        if err != nil {
            return paginatedArticles, fmt.Errorf("failed to scan article row: %w", err)
        }

        if totalCount == 0 { // The first row will contain the total count
            totalCount = count
        }
        articles = append(articles, article)
    }

    // Calculate pagination values correctly
    // Current page based on offset
    currentPage := (offset / limit) + 1             

    // Calculate total pages
    totalPages := (totalCount + limit - 1) / limit

    // Check if requested page exceeds available pages
    if currentPage > totalPages || currentPage == 0 {
        return paginatedArticles, fmt.Errorf("requested page exceeds available pages")
    }

    remainingCount := totalCount - (offset + len(articles)) // Remaining articles after this page
    if remainingCount < 0 {
        remainingCount = 0
    }
    remainingPages := (remainingCount + limit - 1) / limit // Remaining pages after this page

    // Set the paginated articles struct
    paginatedArticles.Articles = articles
    paginatedArticles.RemainingCount = remainingCount
    paginatedArticles.RemainingPages = remainingPages
    paginatedArticles.CurrentPage = currentPage
    paginatedArticles.TotalPages = totalPages

    return paginatedArticles, nil
}


