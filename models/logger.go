package models

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RequestLogEntry represents the structure of a request log entry.
type RequestLogEntry struct {
	Method     string    `bson:"method"`
	Path       string    `bson:"path"`
	StatusCode int       `bson:"status_code"`
	IP         string    `bson:"ip"`
	UserAgent  string    `bson:"user_agent"`
	Duration   int64     `bson:"duration_ms"`
	Timestamp  time.Time `bson:"timestamp"`
}

// Logger handles logging to MongoDB.
type Logger struct {
	collection *mongo.Collection
}

// InitLogger initializes a new Logger instance.
func InitLogger(collection *mongo.Collection) *Logger {
	return &Logger{collection: collection}
}

// MiddlewareLogger is the middleware function that logs requests.
func (l *Logger) MiddlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		// Collect request and response details
		entry := RequestLogEntry{
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Duration:   time.Since(startTime).Milliseconds(),
			Timestamp:  time.Now(),
		}

		// Log the entry to MongoDB
		if l.collection != nil {
			_, err := l.collection.InsertOne(c.Request.Context(), entry)
			if err != nil {
				log.Printf("Failed to log request: %v", err)
			}
		}
	}
}
