package pkg

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

var (
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
)

func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	DBUser = os.Getenv("MYSQL_USER")
	DBPass = os.Getenv("MYSQL_PASSWORD")
	DBHost = os.Getenv("MYSQL_HOST")
	DBPort = os.Getenv("MYSQL_PORT")
	DBName = os.Getenv("MYSQL_DATABASE")
	return nil
}

func ConnectDB() (*sql.DB, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPass, DBHost, DBPort, DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if ping := db.Ping(); ping != nil {
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).Error(ping.Error())
		return nil, ping
	}
	return db, nil
}
