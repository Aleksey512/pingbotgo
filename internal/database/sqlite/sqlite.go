package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"misbotgo/internal/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(cfg *config.Settings) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", cfg.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS chat_ids (
			chat_id TEXT PRIMARY KEY
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SQLiteStorage) AddChatID(ctx context.Context, chatID string) error {
	if s.db == nil {
		return errors.New("SQLite database is not initialized")
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO chat_ids (chat_id) VALUES (?)
	`, chatID)
	return err
}

func (s *SQLiteStorage) RemoveChatID(ctx context.Context, chatID string) error {
	if s.db == nil {
		return errors.New("SQLite database is not initialized")
	}

	_, err := s.db.ExecContext(ctx, `
		DELETE FROM chat_ids WHERE chat_id = ?
	`, chatID)
	return err
}

func (s *SQLiteStorage) GetChatIDs(ctx context.Context) ([]string, error) {
	if s.db == nil {
		return nil, errors.New("SQLite database is not initialized")
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT chat_id FROM chat_ids
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatIDs []string
	for rows.Next() {
		var chatID string
		if err := rows.Scan(&chatID); err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatIDs, nil
}
