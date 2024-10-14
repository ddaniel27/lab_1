package record

import (
	"context"
	"lab1_isbn/internal/core/domain"
	"lab1_isbn/internal/core/ports/repositories"
)

var format = map[bool]int8{
	false: 0,
	true:  1,
}

type Service struct {
	repo repositories.TaskRepository
}

func NewService(repo repositories.TaskRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetRecords(_ context.Context) ([]domain.Task, error) {
	return s.repo.Get(), nil
}

func (s *Service) CreateRecord(ctx context.Context, task domain.Task) error {
	return s.repo.Create(ctx, task)
}

func (s *Service) UpdateRecord(ctx context.Context, task domain.Task) error {
	return s.repo.Update(ctx, task)
}

func (s *Service) DeleteRecord(_ context.Context, number uint) error {
	return s.repo.Delete(number)
}
