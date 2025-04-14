package storage

import (
	"database/sql"
	"fmt"
	c "github.com/evgenyshipko/go-rag-chat-helper/internal/const"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"strings"
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

type Chunk struct {
	Rank float64
	Text string
	Id   int
}

func (storage *SQLStorage) SearchChunks(keywords []string) ([]Chunk, error) {
	var parts []string
	for _, kw := range keywords {
		cleaned := strings.Join(strings.Fields(kw), " & ")
		parts = append(parts, cleaned)
	}
	tsquery := strings.Join(parts, " | ")

	query := `SELECT ts_rank(text_tsvector, to_tsquery('russian', $1)) AS rank, text, id
	FROM doc_chunks
	WHERE text_tsvector @@ to_tsquery('russian', $1)
	ORDER BY rank DESC`

	stmt, err := storage.prepareStmt(query)
	if err != nil {
		return make([]Chunk, 0), err
	}

	rows, err := stmt.Query(tsquery)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var results []Chunk
	for rows.Next() {
		var chunk Chunk
		if err := rows.Scan(&chunk.Rank, &chunk.Text, &chunk.Id); err != nil {
			return nil, fmt.Errorf("row scanning failed: %w", err)
		}
		results = append(results, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return results, nil
}

func (storage *SQLStorage) SaveChunks(chunks []c.DocumentChunk) error {
	tx, err := storage.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}

	stmt, err := tx.Prepare(`
        INSERT INTO doc_chunks (text, text_tsvector)
        VALUES ($1, to_tsvector('russian', $1))
    `)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при подготовке запроса: %w", err)
	}
	defer stmt.Close()

	for _, chunk := range chunks {
		_, err := stmt.Exec(chunk.Text)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка при вставке чанка: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}
