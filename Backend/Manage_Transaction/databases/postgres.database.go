package databases

import (
	"chatify/configs"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectWithRetry continuously tries to connect to PostgreSQL using
// exponential backoff with jitter. It never gives up until context is cancelled.
func ConnectPostgresWithRetry(ctx context.Context, delay time.Duration) (*gorm.DB, error) {
	var dsn string = GetPostgresConnection()

	var (
		db      *gorm.DB
		sqlDB   *sql.DB
		err     error
		attempt int
	)

	for {
		attempt++
		log.Printf("🔌 Connecting to PostgreSQL (attempt %d)", attempt)

		// open connection
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err = db.DB()
		}

		if err == nil {
			err = sqlDB.PingContext(ctx)
		}

		if err == nil {
			// configure pooling
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(151)
			sqlDB.SetConnMaxLifetime(10 * time.Minute)

			DB = db
			log.Println("🟢 PostgreSQL connected successfully!")

			go MonitorPostgresConnection(ctx, sqlDB, delay)
			return db, nil
		}

		log.Printf("❌ Attempt %d failed: %v", attempt, err)

		// exponential backoff with jitter
		var sleep time.Duration = delay * time.Duration(1<<uint(attempt-1))
		var jitter time.Duration = time.Duration(rand.Int63n(int64(sleep / 2)))
		sleep = sleep/2 + jitter

		log.Printf("⏳ retrying in %s...", sleep)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(sleep):
		}
	}
}

// monitor pings the DB at the given interval until context is cancelled.
func MonitorPostgresConnection(ctx context.Context, sqlDB *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🚪 Stopping PostgreSQL Monitor")
			return
		case <-ticker.C:
			if err := sqlDB.PingContext(ctx); err != nil {

				log.Printf("🛑 Lost DB connection: %v", err)

				// trigger reconnection and update DB in background
				go func(oldDB *sql.DB) {
					log.Println("🔄 Attempting to reconnect to PostgreSQL...")
					NewDB, err := ConnectPostgresWithRetry(ctx, interval)

					if err != nil {
						log.Printf("❌ Reconnect failed: %v", err)
						return
					}

					NewSQLDB, err := NewDB.DB()
					if err != nil {
						log.Printf("❌ Could not get new sql.DB: %v", err)
						return
					}

					// Swap global DB and sqlDB reference
					DB = NewDB
					sqlDB = NewSQLDB
					log.Println("🟢 DB connection restored by monitor")
				}(sqlDB)
				return
			}
		}
	}
}

func GetPostgresConnection() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=%s",
		configs.ENV.DatabaseSetting.User,
		configs.ENV.DatabaseSetting.Password,
		configs.ENV.DatabaseSetting.DatabaseName,
		configs.ENV.DatabaseSetting.Host,
		configs.ENV.DatabaseSetting.DatabasePort,
		configs.ENV.DatabaseSetting.Timezone,
	)
}
