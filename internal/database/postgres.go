package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
    // Get database URL from environment
    dbURL := os.Getenv("DATABASE_URL")
    
    // If DATABASE_URL not set, build from individual vars
    if dbURL == "" {
        host := os.Getenv("DB_HOST")
        port := os.Getenv("DB_PORT")
        user := os.Getenv("DB_USER")
        password := os.Getenv("DB_PASSWORD")
        dbname := os.Getenv("DB_NAME")
        sslmode := os.Getenv("DB_SSLMODE")
        
        if sslmode == "" {
            sslmode = "disable"
        }
        
        dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
            host, port, user, password, dbname, sslmode)
    }

    var err error
    DB, err = sql.Open("postgres", dbURL)
    if err != nil {
        return fmt.Errorf("error opening database: %v", err)
    }

    // Test connection
    if err = DB.Ping(); err != nil {
        return fmt.Errorf("error connecting to database: %v", err)
    }

    log.Println("✓ Database connected successfully")
    return nil
}

func RunMigrations() error {
    // Read migration file
    migration, err := os.ReadFile("internal/database/migrations/001_create_devices_table.sql")
    if err != nil {
        return fmt.Errorf("error reading migration: %v", err)
    }

    // Execute migration
    _, err = DB.Exec(string(migration))
    if err != nil {
        return fmt.Errorf("error running migration: %v", err)
    }

    log.Println("✓ Database migrations completed")
    return nil
}
