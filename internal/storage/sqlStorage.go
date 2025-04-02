package storage

import (
	"database/sql"
	"fmt"
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"sync"
	"time"
)

type SQLStorage struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
	mu         sync.RWMutex //ЗАПОМНИТЬ: мапа не потокобезопасна, поэтому при конкуррентном чтении/записи могут возникать ошибки конкуррентного доступа к данным
}

func NewSQLStorage(db *sql.DB) *SQLStorage {
	return &SQLStorage{
		db:         db,
		statements: map[string]*sql.Stmt{},
		mu:         sync.RWMutex{},
	}
}

func (storage *SQLStorage) prepareStmt(query string) (*sql.Stmt, error) {
	storage.mu.RLock() // Блокируем только для чтения
	stmt, exists := storage.statements[query]
	storage.mu.RUnlock() // Разблокируем чтение
	if exists {
		return stmt, nil
	}

	storage.mu.Lock()         // Блокируем на запись
	defer storage.mu.Unlock() // Разблокируем запись

	stmt, err := storage.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	storage.statements[query] = stmt
	return stmt, nil
}

func (storage *SQLStorage) InsertUser(login string, hashedPassword string) error {
	logger.Instance.Debugw("InsertUser", "login", login, "hashedPassword", hashedPassword)

	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2);`

	stmt, err := storage.prepareStmt(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(login, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	Id             string
	Login          string
	HashedPassword string
	CreatedAt      time.Time
}
type UserNotFoundError struct {
	Login string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user with login %s not found", e.Login)
}

func (storage *SQLStorage) GetUser(login string) (*User, error) {
	logger.Instance.Debugw("GetUser", "login", login)

	query := `SELECT * FROM users WHERE login = $1;`

	stmt, err := storage.prepareStmt(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(login)

	var user User

	err = row.Scan(&user.Id, &user.Login, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Instance.Debugw("User not found", "login", login)
			return nil, &UserNotFoundError{login}
		}

		return nil, err
	}

	return &user, nil
}

func (storage *SQLStorage) SaveUserTokens(userId string, access string, refresh string) error {
	logger.Instance.Info("SaveUserTokens")

	query := `INSERT INTO sessions (user_id, access_token, refresh_token) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE 
    SET access_token = $2, refresh_token = $3;`

	stmt, err := storage.prepareStmt(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userId, access, refresh)
	if err != nil {
		return err
	}
	return nil
}

func (storage *SQLStorage) DropUserTokens(userId string) error {
	logger.Instance.Info("DropUserTokens")

	query := `DELETE FROM sessions WHERE user_id = $1;`

	stmt, err := storage.prepareStmt(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userId)
	if err != nil {
		return err
	}
	return nil
}
