package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/hisnameisivan/demo_url_short/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("[%s] conn err: %w", fn, err)
	}

	// migrations
	stmt, err := db.Prepare(createTableUrlQuery)
	if err != nil {
		return nil, fmt.Errorf("[%s] prepare err: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("[%s] exec err: %w", fn, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const fn = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare(insertUrlQuery)
	if err != nil {
		return 0, fmt.Errorf("[%s] prepare err: %w", fn, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok &&
			errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			// sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("[%s] exec err: %w", fn, storage.ErrUrlExists)
		}

		return 0, fmt.Errorf("[%s] exec err: %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[%s] get last id err: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const fn = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare(getUrlQuery)
	if err != nil {
		return "", fmt.Errorf("[%s] prepare err: %w", fn, err)
	}

	var resultUrl string
	err = stmt.QueryRow(alias).Scan(&resultUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("[%s] query err: %w", fn, storage.ErrUrlNotFound)
		}

		return "", fmt.Errorf("[%s] query err: %w", fn, err)
	}

	return resultUrl, nil
}
