package storage

import (
	"context"
	"lab1_isbn/internal/core/domain"

	"github.com/uptrace/bun"
)

type Storage struct {
	db *bun.DB
}

func NewStorage(db *bun.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Get() []domain.Task {
	records := []domain.Task{}

	return records
}

func (s *Storage) Update(ctx context.Context, record domain.Task) error {
	_, err := s.db.NewUpdate().Model(&record).Where("name = ?", record.Name).Exec(ctx)
	return err
}

func (s *Storage) Create(ctx context.Context, record domain.Task) error {
	_, err := s.db.NewInsert().Model(&record).Exec(ctx)
	return err
}

func (s *Storage) Delete(number uint) error {
	return nil
}
