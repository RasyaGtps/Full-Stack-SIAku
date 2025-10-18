package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config struct untuk simpan env
type Config struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	JWTSecret           string
	ServerPort          string
	WhatsAppServiceURL  string
}

var AppConfig Config
var DB *gorm.DB

// Load .env ke struct
func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	AppConfig = Config{
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		ServerPort:         os.Getenv("SERVER_PORT"),
		WhatsAppServiceURL: os.Getenv("WHATSAPP_SERVICE_URL"),
	}
	return nil
}

// custom logger sederhana
type customLogger struct{}

func (l *customLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}
func (l *customLogger) Info(ctx context.Context, msg string, data ...interface{})  {}
func (l *customLogger) Warn(ctx context.Context, msg string, data ...interface{})  {}
func (l *customLogger) Error(ctx context.Context, msg string, data ...interface{}) {}
func (l *customLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		fmt.Printf("❌ Error: %v | SQL: %s\n", err, sql)
		return
	}
	fmt.Printf("✅ Query OK | %s | Rows: %d | Time: %v\n", sql, rows, elapsed)
}

// auto create DB jika belum ada
func createDatabase() error {
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable",
		AppConfig.DBHost,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBPort,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error connect postgres: %w", err)
	}
	defer db.Close()

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", AppConfig.DBName)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking db existence: %w", err)
	}

	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", AppConfig.DBName))
		if err != nil {
			return fmt.Errorf("error creating db: %w", err)
		}
		fmt.Printf("🎉 Database '%s' created successfully\n", AppConfig.DBName)
	} else {
		fmt.Printf("ℹ️ Database '%s' already exists\n", AppConfig.DBName)
	}

	return nil
}

// InitDB koneksi ke postgres via GORM
func InitDB() (*gorm.DB, error) {
	if err := createDatabase(); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		AppConfig.DBHost,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBName,
		AppConfig.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &customLogger{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed connect db: %w", err)
	}

	fmt.Println("✅ Connected to DB")
	DB = db
	return db, nil
}
