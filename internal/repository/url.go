package repository

import (
	"database/sql"
	"errors"

	"github.com/davidjchavez/url-shortener/internal/model"
)

type URLRepository struct {
	db *sql.DB
}

var ErrNotFound = errors.New("url not found")

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		db,
	}
}

func (r *URLRepository) Create(url *model.URL) error {
	sqlStatement := `INSERT INTO url (code, original_url)
	VALUES ($1, $2)
	RETURNING id_url`
	err := r.db.QueryRow(sqlStatement, url.Code, url.OriginalURL).Scan(&url.ID)
	return err
}

func (r *URLRepository) GetByCode(code string) (*model.URL, error) {
	sqlStatement := `SELECT id_url, code, original_url, clicks, created_at FROM url WHERE code = $1`
	var url model.URL
	err := r.db.QueryRow(sqlStatement, code).Scan(&url.ID, &url.Code, &url.OriginalURL, &url.Clicks, &url.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) IncrementClicks(code string) error {
	sqlStatement := `UPDATE url SET clicks = clicks + 1 WHERE code = $1`
	_, err := r.db.Exec(sqlStatement, code)
	return err
}

func (r *URLRepository) CodeExists(code string) (bool, error) {
	sqlStatement := `SELECT EXISTS(SELECT 1 FROM url WHERE code = $1)`
	var exists bool
	err := r.db.QueryRow(sqlStatement, code).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
