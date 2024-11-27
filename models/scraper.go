package models

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
)

// Periodically scrape articles every 5 minutes
func StartPeriodicScraping(urls []string) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        log.Println("Starting new scraping cycle...")
        FetchArticlesFromMultiUrls(urls)
        log.Println("Scraping cycle completed.")

        <-ticker.C // Wait for the next tick
    }
}


// Fetch articles concurrently from multiple URLs
func FetchArticlesFromMultiUrls(urls []string) {
    // to handle multi goroutines
    var wg sync.WaitGroup
    // ensure safe access to the DB
    var mutex sync.Mutex

    for _, url := range urls {
        wg.Add(1)
        go func(url string) {
            defer wg.Done()

            articles, err := fetchArticlesFromURL(url)
            if err != nil {
                log.Printf("Failed to fetch from %s: %v\n", url, err)
                return
            }

            // Insert into the database
            mutex.Lock()
            for _, article := range articles {
                 exists, err := ArticleExists(article.Link)
                if err != nil {
                    log.Printf("Failed to check if article exists: %v\n", err)
                    continue
                }

                if exists {
                    log.Printf("Article with link %s already exists, skipping...\n", article.Link)
                    continue
                }

                _, errSave := SaveArticle(article)
                if errSave != nil {
                    log.Printf("Failed to save article: %v", err)
                } else {
                    log.Printf("Saved article: %s", article.Title)
                }
            }
            mutex.Unlock()
        }(url)
    }

    wg.Wait()
}

// fetchArticles fetches and parses articles from a given URL, then saves them to the DB.
func fetchArticlesFromURL(url string) ([]Article, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch URL: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("non-200 HTTP status: %d", resp.StatusCode)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse HTML: %v", err)
    }

    var articles []Article
    doc.Find("article").Each(func(i int, s *goquery.Selection) {
        title := s.Find("h2 a").Text()
        link, _ := s.Find("h2 a").Attr("href")
        timestamp, _ := time.Parse(time.RFC3339, s.Find("time").AttrOr("datetime", ""))

        article := Article{
            Title:     title,
            Link:      link,
            Timestamp: timestamp,
        }
        articles = append(articles, article)
    })

    return articles, nil
}

