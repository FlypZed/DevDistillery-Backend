package board

import (
	"context"
	"fmt"
	"func/internal/domain"
	repository "func/internal/repository/board"
)

type boardServiceImpl struct {
	repo repository.BoardRepository
}

func NewBoardService(repo repository.BoardRepository) Service {
	return &boardServiceImpl{repo: repo}
}

func (s *boardServiceImpl) GetBoard(ctx context.Context, id string) (*domain.Board, error) {
	if id == "" {
		return nil, fmt.Errorf("board ID cannot be empty")
	}

	board, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	return board, nil
}

func (s *boardServiceImpl) SaveBoard(ctx context.Context, board *domain.Board) error {
	if board == nil {
		return fmt.Errorf("board cannot be nil")
	}

	if board.Title == "" {
		return fmt.Errorf("board title cannot be empty")
	}

	// Validar nodos y edges si es necesario
	for _, node := range board.Nodes {
		if node.ID == "" {
			return fmt.Errorf("node ID cannot be empty")
		}
		if node.Type == "" {
			return fmt.Errorf("node type cannot be empty")
		}
	}

	for _, edge := range board.Edges {
		if edge.ID == "" {
			return fmt.Errorf("edge ID cannot be empty")
		}
		if edge.Source == "" || edge.Target == "" {
			return fmt.Errorf("edge source and target cannot be empty")
		}
	}

	err := s.repo.Save(ctx, board)
	if err != nil {
		return fmt.Errorf("failed to save board: %w", err)
	}

	return nil
}

func (s *boardServiceImpl) UpdateBoard(ctx context.Context, board *domain.Board) error {
	if board == nil {
		return fmt.Errorf("board cannot be nil")
	}

	if board.ID == "" {
		return fmt.Errorf("board ID cannot be empty")
	}

	if board.Title == "" {
		return fmt.Errorf("board title cannot be empty")
	}

	for _, node := range board.Nodes {
		if node.ID == "" {
			return fmt.Errorf("node ID cannot be empty")
		}
		if node.Type == "" {
			return fmt.Errorf("node type cannot be empty")
		}
		if node.BoardID != board.ID {
			return fmt.Errorf("node does not belong to this board")
		}
	}

	for _, edge := range board.Edges {
		if edge.ID == "" {
			return fmt.Errorf("edge ID cannot be empty")
		}
		if edge.Source == "" || edge.Target == "" {
			return fmt.Errorf("edge source and target cannot be empty")
		}
		if edge.BoardID != board.ID {
			return fmt.Errorf("edge does not belong to this board")
		}
	}

	err := s.repo.Update(ctx, board)
	if err != nil {
		return fmt.Errorf("failed to update board: %w", err)
	}

	return nil
}
