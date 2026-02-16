package database

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	dbURL := strings.TrimSpace(firstNonEmpty(
		"DATABASE_URL",
		"DATABASE_PRIVATE_URL",
		"DATABASE_PUBLIC_URL",
		"POSTGRES_URL",
		"POSTGRESQL_URL",
		"PG_URL",
	))

	if dbURL == "" {
		host := strings.TrimSpace(firstNonEmpty("DB_HOST", "PGHOST", "POSTGRES_HOST"))
		port := strings.TrimSpace(firstNonEmpty("DB_PORT", "PGPORT", "POSTGRES_PORT"))
		user := strings.TrimSpace(firstNonEmpty("DB_USER", "PGUSER", "POSTGRES_USER"))
		password := firstNonEmpty("DB_PASSWORD", "PGPASSWORD", "POSTGRES_PASSWORD")
		dbname := strings.TrimSpace(firstNonEmpty("DB_NAME", "PGDATABASE", "POSTGRES_DB"))
		sslmode := strings.TrimSpace(firstNonEmpty("DB_SSLMODE", "PGSSLMODE"))

		if sslmode == "" {
			if strings.EqualFold(os.Getenv("APP_ENV"), "production") {
				sslmode = "require"
			} else {
				sslmode = "disable"
			}
		}

		missing := make([]string, 0)
		if host == "" {
			missing = append(missing, "DB_HOST/PGHOST")
		}
		if port == "" {
			missing = append(missing, "DB_PORT/PGPORT")
		}
		if user == "" {
			missing = append(missing, "DB_USER/PGUSER")
		}
		if password == "" {
			missing = append(missing, "DB_PASSWORD/PGPASSWORD")
		}
		if dbname == "" {
			missing = append(missing, "DB_NAME/PGDATABASE")
		}

		if len(missing) > 0 {
			return fmt.Errorf("database config incomplete: missing %s. Set DATABASE_URL or all DB_* vars", strings.Join(missing, ", "))
		}

		q := url.Values{}
		q.Set("sslmode", sslmode)

		u := &url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(user, password),
			Host:     net.JoinHostPort(host, port),
			Path:     dbname,
			RawQuery: q.Encode(),
		}
		dbURL = u.String()
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
	return nil
}

func firstNonEmpty(keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); strings.TrimSpace(val) != "" {
			return val
		}
	}
	return ""
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
	return nil
}
