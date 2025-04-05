package board

import (
	"context"
	"func/internal/domain"
)

type Service interface {
	GetBoard(ctx context.Context, id string) (*domain.Board, error)
	SaveBoard(ctx context.Context, board *domain.Board) error
	UpdateBoard(ctx context.Context, board *domain.Board) error
}
