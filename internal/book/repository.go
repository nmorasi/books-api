package book

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Book struct {
	ID          int64     `db:"id"          json:"id"`
	UserID      int64     `db:"user_id"     json:"user_id"`
	Title       string    `db:"title"       json:"title"`
	Author      string    `db:"author"      json:"author"`
	Description string    `db:"description" json:"description"`
	Year        int       `db:"year"        json:"year"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, b *Book) error {
	query := `
		INSERT INTO books (user_id, title, author, description, year)
		VALUES (:user_id, :title, :author, :description, :year)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, b)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	}
	return nil
}

func (r *Repository) List(ctx context.Context, userID int64) ([]Book, error) {
	var books []Book
	err := r.db.SelectContext(ctx, &books,
		`SELECT * FROM books WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (r *Repository) GetByID(ctx context.Context, id, userID int64) (*Book, error) {
	var b Book
	err := r.db.GetContext(ctx, &b,
		`SELECT * FROM books WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) Update(ctx context.Context, b *Book) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE books SET title=$1, author=$2, description=$3, year=$4, updated_at=NOW()
		WHERE id=$5 AND user_id=$6`,
		b.Title, b.Author, b.Description, b.Year, b.ID, b.UserID)
	return err
}

func (r *Repository) Delete(ctx context.Context, id, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM books WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}
