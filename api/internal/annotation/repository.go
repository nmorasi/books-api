package annotation

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Annotation struct {
	ID        int64     `db:"id"         json:"id"`
	BookID    int64     `db:"book_id"    json:"book_id"`
	UserID    int64     `db:"user_id"    json:"user_id"`
	Body      string    `db:"body"       json:"body"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, a *Annotation) error {
	query := `
		INSERT INTO annotations (book_id, user_id, body)
		VALUES (:book_id, :user_id, :body)
		RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, a)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&a.ID, &a.CreatedAt)
	}
	return nil
}

func (r *Repository) ListByBook(ctx context.Context, bookID, userID int64) ([]Annotation, error) {
	var annotations []Annotation
	err := r.db.SelectContext(ctx, &annotations,
		`SELECT * FROM annotations WHERE book_id = $1 AND user_id = $2 ORDER BY created_at DESC`,
		bookID, userID)
	if err != nil {
		return nil, err
	}
	return annotations, nil
}

func (r *Repository) Delete(ctx context.Context, id, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM annotations WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}
