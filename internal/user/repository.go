package user

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           int64     `db:"id"            json:"id"`
	Name         string    `db:"name"          json:"name"`
	Email        string    `db:"email"         json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES (:name, :email, :password_hash)
		RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, u)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&u.ID, &u.CreatedAt)
	}
	return nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
