package book

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type CreateRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}

type UpdateRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateRequest) (*Book, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errors.New("title is required")
	}
	if strings.TrimSpace(req.Author) == "" {
		return nil, errors.New("author is required")
	}

	b := &Book{
		UserID:      userID,
		Title:       strings.TrimSpace(req.Title),
		Author:      strings.TrimSpace(req.Author),
		Description: strings.TrimSpace(req.Description),
		Year:        req.Year,
	}

	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) List(ctx context.Context, userID int64) ([]Book, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id, userID int64) (*Book, error) {
	b, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return b, nil
}

func (s *Service) Update(ctx context.Context, id, userID int64, req UpdateRequest) (*Book, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errors.New("title is required")
	}
	if strings.TrimSpace(req.Author) == "" {
		return nil, errors.New("author is required")
	}

	b, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	b.Title = strings.TrimSpace(req.Title)
	b.Author = strings.TrimSpace(req.Author)
	b.Description = strings.TrimSpace(req.Description)
	b.Year = req.Year

	if err := s.repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) Delete(ctx context.Context, id, userID int64) error {
	_, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("book not found")
		}
		return err
	}
	return s.repo.Delete(ctx, id, userID)
}
