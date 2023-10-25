package repository

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	_ "github.com/lib/pq"

	"tgtrello/config"
)

func NewDB(logger *zap.Logger, cfg *config.Config) *sql.DB {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.DBName)

	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		logger.Fatal("open database connection", zap.Error(err))
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("database connection timeout", zap.Error(err))
	}

	return db
}
