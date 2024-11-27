package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB

// ConnectPostgresDB initializes the PostgreSQL database connection and runs migrations.
func ConnectPostgresDB() error {
	// Load environment variables
    err := godotenv.Load()
	if err != nil {
		log.Println("Warning: failed to load .env file. Falling back to environment variables.")
	}

	// Build connection string from environment variables
	connStr := buildConnectionString()
	
	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	// Verify connection
    errDB := db.Ping()
	if errDB != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	PostgresDB = db
	log.Println("Connected to PostgreSQL database")

	// Run migrations
    errMigration := runMigrations(db);
	if errMigration != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// buildConnectionString constructs the PostgreSQL connection string.
func buildConnectionString() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		user, password, dbName, host, port, sslMode,
	)
}

// runMigrations ensures required database tables are created.
func runMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS articles (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		link TEXT NOT NULL,
		timestamp TIMESTAMPTZ NOT NULL
	);`

    _, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create articles table: %w", err)
	}

	log.Println("Migrations completed: articles table ensured")
	return nil
}
