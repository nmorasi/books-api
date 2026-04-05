package character

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Character struct {
	ID        int64     `db:"id"         json:"id"`
	BookID    int64     `db:"book_id"    json:"book_id"`
	UserID    int64     `db:"user_id"    json:"user_id"`
	Name      string    `db:"name"       json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, bookID, userID int64) ([]Character, error) {
	var characters []Character
	err := r.db.SelectContext(ctx, &characters,
		`SELECT * FROM characters WHERE book_id = $1 AND user_id = $2 ORDER BY name ASC`,
		bookID, userID)
	return characters, err
}

func (r *Repository) Create(ctx context.Context, c *Character) error {
	rows, err := r.db.NamedQueryContext(ctx,
		`INSERT INTO characters (book_id, user_id, name)
		 VALUES (:book_id, :user_id, :name)
		 RETURNING id, created_at`, c)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&c.ID, &c.CreatedAt)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM characters WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}
