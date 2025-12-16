package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"simpus/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Error loading config or .env file: %v", err)
		// Proceeding as config.Load has defaults
	}

	// Connect to Database (without DB name first to create it if needed)
	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	dbRoot, err := sql.Open("mysql", dsnRoot)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL root: %v", err)
	}
	defer dbRoot.Close()

	// Create Database if not exists
	_, err = dbRoot.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.Database.Name))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", cfg.Database.Name, err)
	}
	fmt.Printf("Database %s checked/created.\n", cfg.Database.Name)

	// Connect to specific database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Read migration file
	migrationFile := filepath.Join("database", "migrations", "001_init.sql")
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		log.Fatalf("Failed to read migration file %s: %v", migrationFile, err)
	}

	// Execute migration
	// Note: multiStatements=true is required in DSN for this to work with multiple queries in one Execute,
	// or we split them.
	// Since 001_init.sql might contain DELIMITER or complex statements, simple splitting might be risky.
	// However, the Go mysql driver supports executing multiple statements if enabled.

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	fmt.Println("Migration completed successfully!")
}
