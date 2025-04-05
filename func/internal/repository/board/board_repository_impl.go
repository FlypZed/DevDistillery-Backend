package board

import (
	"context"
	"errors"
	"fmt"
	"time"

	"func/internal/domain"
	"gorm.io/gorm"
)

type boardRepositoryImpl struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepositoryImpl{db: db}
}

func (r *boardRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.Board, error) {
	var board domain.Board

	result := r.db.WithContext(ctx).
		Preload("Nodes").
		Preload("Edges").
		First(&board, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("board not found")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get board: %w", result.Error)
	}

	return &board, nil
}

func (r *boardRepositoryImpl) Save(ctx context.Context, board *domain.Board) error {
	now := time.Now()
	board.CreatedAt = now
	board.UpdatedAt = now

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(board).Error; err != nil {
			return fmt.Errorf("failed to create board: %w", err)
		}

		for i := range board.Nodes {
			board.Nodes[i].BoardID = board.ID
			if err := tx.Create(&board.Nodes[i]).Error; err != nil {
				return fmt.Errorf("failed to create node: %w", err)
			}
		}

		for i := range board.Edges {
			board.Edges[i].BoardID = board.ID
			if err := tx.Create(&board.Edges[i]).Error; err != nil {
				return fmt.Errorf("failed to create edge: %w", err)
			}
		}

		return nil
	})
}

func (r *boardRepositoryImpl) Update(ctx context.Context, board *domain.Board) error {
	board.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(board).Updates(map[string]interface{}{
			"title":       board.Title,
			"description": board.Description,
			"updated_at":  board.UpdatedAt,
		}).Error; err != nil {
			return fmt.Errorf("failed to update board: %w", err)
		}

		if err := tx.Where("board_id = ?", board.ID).Delete(&domain.Node{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing nodes: %w", err)
		}
		if err := tx.Where("board_id = ?", board.ID).Delete(&domain.Edge{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing edges: %w", err)
		}

		for i := range board.Nodes {
			board.Nodes[i].BoardID = board.ID
			if err := tx.Create(&board.Nodes[i]).Error; err != nil {
				return fmt.Errorf("failed to create node: %w", err)
			}
		}

		for i := range board.Edges {
			board.Edges[i].BoardID = board.ID
			if err := tx.Create(&board.Edges[i]).Error; err != nil {
				return fmt.Errorf("failed to create edge: %w", err)
			}
		}

		return nil
	})
}
