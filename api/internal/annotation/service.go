package annotation

import (
	"context"
	"errors"
	"strings"
)

type CreateRequest struct {
	Body string `json:"body"`
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, bookID, userID int64, req CreateRequest) (*Annotation, error) {
	if strings.TrimSpace(req.Body) == "" {
		return nil, errors.New("body is required")
	}

	a := &Annotation{
		BookID: bookID,
		UserID: userID,
		Body:   strings.TrimSpace(req.Body),
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) ListByBook(ctx context.Context, bookID, userID int64) ([]Annotation, error) {
	return s.repo.ListByBook(ctx, bookID, userID)
}

func (s *Service) Delete(ctx context.Context, id, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}
