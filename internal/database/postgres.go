package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var lastResetTime time.Time

const MaxDevices = 25
const ResetInterval = 1 * time.Hour

func Connect() error {
	dbURL := os.Getenv("DATABASE_URL")
	
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

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("✓ Database connected successfully")
	
	// Start auto-reset background task
	go autoResetDatabase()
	
	return nil
}

func RunMigrations() error {
	migration, err := os.ReadFile("internal/database/migrations/001_create_devices_table.sql")
	if err != nil {
		return fmt.Errorf("error reading migration: %v", err)
	}

	_, err = DB.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("error running migration: %v", err)
	}

	log.Println("✓ Database migrations completed")
	lastResetTime = time.Now()
	return nil
}

// Check if device limit reached
func IsDeviceLimitReached() (bool, int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM devices").Scan(&count)
	if err != nil {
		return false, 0, err
	}
	
	return count >= MaxDevices, count, nil
}

// Reset database to initial state
func ResetDatabase() error {
	log.Println("🔄 Resetting database to demo state...")
	
	// Delete all devices
	_, err := DB.Exec("DELETE FROM devices")
	if err != nil {
		return fmt.Errorf("error deleting devices: %v", err)
	}
	
	// Re-insert dummy data
	_, err = DB.Exec(`
		INSERT INTO devices (name, ip_address, location, status, last_seen) VALUES
			('Router-BDG-01', '192.168.100.11', 'Bandung', 'online', NOW()),
			('Router-JKT-01', '192.168.100.12', 'Jakarta', 'online', NOW()),
			('Router-SBY-01', '192.168.100.13', 'Surabaya', 'offline', NOW() - INTERVAL '1 hour')
		ON CONFLICT (name) DO NOTHING
	`)
	
	if err != nil {
		return fmt.Errorf("error inserting dummy data: %v", err)
	}
	
	lastResetTime = time.Now()
	log.Println("✓ Database reset completed - 3 demo devices restored")
	return nil
}

// Auto-reset background task
func autoResetDatabase() {
	ticker := time.NewTicker(10 * time.Minute) // Check every 10 minutes
	defer ticker.Stop()
	
	for range ticker.C {
		// Check if reset interval passed
		if time.Since(lastResetTime) >= ResetInterval {
			// Check device count
			var count int
			err := DB.QueryRow("SELECT COUNT(*) FROM devices").Scan(&count)
			if err != nil {
				log.Printf("Error checking device count: %v", err)
				continue
			}
			
			// Reset if limit was reached
			if count >= MaxDevices {
				log.Printf("🔄 Auto-reset triggered: %d devices (limit: %d)", count, MaxDevices)
				if err := ResetDatabase(); err != nil {
					log.Printf("Error during auto-reset: %v", err)
				}
			}
		}
	}
}

// Get time until next reset
func GetTimeUntilReset() time.Duration {
	elapsed := time.Since(lastResetTime)
	remaining := ResetInterval - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}
