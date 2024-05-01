package database

import (
	"aut_reg/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func InitDB(storagePath string) (*Storage, error) {
	var err error
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, err
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email, username string, password []byte) (int64, error) {
	prep, err := s.db.Prepare("INSERT INTO users(email, username, pass_hash) VALUES (?, ?, ?)")
	if err != nil {
		//TODO
		return 0, fmt.Errorf("Error saving user1: %d", err)
	}

	res, err := prep.ExecContext(ctx, email, username, password)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("Error: %w", ErrUserExists)
		}
		return 0, fmt.Errorf("Error: %w", err)
	}

	id, err := res.LastInsertId()

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, username string) (models.User, error) {
	prep, err := s.db.Prepare("SELECT id, email, username, pass_hash FROM users WHERE username = ?")
	if err != nil {
		//TODO
		return models.User{}, fmt.Errorf("Error geting user")
	}
	row := prep.QueryRowContext(ctx, username)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("Error: %w", ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("Error: %w", err)
	}
	return user, nil
}
