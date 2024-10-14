package repositories

import (
	"context"
	"lab1_isbn/internal/core/domain"
)

type TaskRepository interface {
	Get() []domain.Task
	Create(ctx context.Context, record domain.Task) error
	Update(ctx context.Context, record domain.Task) error
	Delete(number uint) error
}
