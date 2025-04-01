package db

import (
	"database/sql"
	"embed"
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
)

func ConnectToDB(serverDSN string, autoMigrations bool) (*sql.DB, error) {
	db, err := sql.Open("pgx", serverDSN)
	if err != nil {
		logger.Instance.Warnw("ConnectToDB", "Не удалось подключиться к базе данных", err)
		return nil, err
	}

	if autoMigrations {
		err = RunMigrations(db)
		if err != nil {
			logger.Instance.Warnw("RunMigrations", "Ошибка проката миграций", err)
			return nil, err
		}
	}

	return db, nil
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(db *sql.DB) error {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	logger.Instance.Info("Миграции сработали успешно")

	return nil
}
