package services

import (
	"context"
	"lab1_isbn/internal/core/domain"
)

type RecordService interface {
	GetRecords(ctx context.Context) ([]domain.Task, error)
	CreateRecord(ctx context.Context, task domain.Task) error
	UpdateRecord(ctx context.Context, task domain.Task) error
	DeleteRecord(ctx context.Context, number uint) error
}
