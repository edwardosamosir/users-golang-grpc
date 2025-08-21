package connection

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() *gorm.DB {
	// Resolve absolute path to .env (relative to this file)
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../../")
	envPath := filepath.Join(projectRoot, ".env")

	// Load env variables
	envMap, err := godotenv.Read(envPath)
	if err != nil {
		log.Fatalf("[fatal] could not load .env from %s: %v", envPath, err)
	}
	log.Printf("[info] loaded .env from %s", envPath)

	// Read required vars
	host := envMap["DB_HOST"]
	port := envMap["DB_PORT"]
	user := envMap["DB_USER"]
	pass := envMap["DB_PASS"]
	dbname := envMap["DB_NAME"]
	ssl := envMap["DB_SSLMODE"]

	log.Printf("[env] DB config → host=%s port=%s user=%s db=%s sslmode=%s", host, port, user, dbname, ssl)

	if host == "" || port == "" || user == "" || dbname == "" {
		log.Fatalf("[fatal] missing required database configuration")
	}

	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s TimeZone=UTC", host, port, user, dbname, ssl)
	if pass != "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC", host, port, user, pass, dbname, ssl)
	}

	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("[info] connected to PostgreSQL successfully")
	return connection
}
