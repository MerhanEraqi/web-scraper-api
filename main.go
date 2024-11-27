package main

import (
	"log"
	"os"
	"web-scraper-api/db"
	"web-scraper-api/models"
	"web-scraper-api/routes"
	"web-scraper-api/server"

	"github.com/gin-gonic/gin"
)

func main() {
	// Host news html pages locally 
	go server.HostStaticPages()
	
	// Initialize database
    db.ConnectPostgresDB()
    defer db.PostgresDB.Close()

	err := db.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
    
    // Initialize the logger
	logsCollection := db.MongoClient.Database("logging_db").Collection("logs")
	logger := models.InitLogger(logsCollection)

	// Call SetupLoggingToFile and handle any errors and Ensure the file is closed when the program exits
    logFile, err := SetupLoggingToFile("app.log")
    if err != nil {
        log.Println("Error setting up log file:", err)
    }
    defer logFile.Close()

    // Fetch articles from local HTML page
	// List of URLs to scrape
    urls := []string{
        "http://localhost:8081/news.html",
        "http://localhost:8081/news2.html",
    }

	go models.StartPeriodicScraping(urls)

	// Initialize the Gin router
    router := gin.Default()
	router.Use(logger.MiddlewareLogger())
	routes.RegisterRoutes(router)

    // Start the server on port 8080
    log.Println("Server running on port 8080...")
    log.Fatal(router.Run(":8080"))
}


// SetupLoggingToFile sets up the log output to a specified file.
func SetupLoggingToFile(logFilePath string) (*os.File, error) {
    // Open (or create) the log file
    file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err // Return the error if file opening fails
    }

    // Set log output to the file
    log.SetOutput(file)

    return file, nil
}