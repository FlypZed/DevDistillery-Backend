package board

import (
	"context"
	"func/internal/domain"
)

type BoardRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Board, error)
	Save(ctx context.Context, board *domain.Board) error
	Update(ctx context.Context, board *domain.Board) error
}
