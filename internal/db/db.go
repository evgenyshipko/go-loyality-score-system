package db

import (
	"database/sql"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	_ "github.com/jackc/pgx/stdlib"
)

func ConnectToDB(serverDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", serverDSN)
	if err != nil {
		logger.Instance.Warnw("ConnectToDB", "Не удалось подключиться к базе данных", err)
		return nil, err
	}

	return db, nil
}
