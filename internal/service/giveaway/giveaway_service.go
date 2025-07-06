package giveaway

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"giveaway-service/internal/domain"
	"giveaway-service/internal/dto"
	"log/slog"
	"time"
)

type Finder interface {
	Find(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Giveaway, error)
}

type Saver interface {
	Save(
		ctx context.Context,
		dto domain.Giveaway,
	) (*domain.Giveaway, error)
}

type Updater interface {
	Update(
		ctx context.Context,
	) (*domain.Giveaway, error)
}

type Deleter interface {
	Delete(
		ctx context.Context,
	) error
}

type Service struct {
	log             *slog.Logger
	pool            *pgxpool.Pool
	giveawayFinder  Finder
	giveawaySaver   Saver
	giveawayUpdater Updater
	giveawayDeleter Deleter
}

func New(
	log *slog.Logger,
	pool *pgxpool.Pool,
	finder Finder,
	saver Saver,
	updater Updater,
	deleter Deleter,
) *Service {
	return &Service{
		log:             log,
		pool:            pool,
		giveawayFinder:  finder,
		giveawaySaver:   saver,
		giveawayUpdater: updater,
		giveawayDeleter: deleter,
	}
}

func (s *Service) Find(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Giveaway, error) {
	const op = "service.giveaway.Find"
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s %s", op, err))
		return nil, err
	}
	// Безопасный откат при любом исходе
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			if err == nil {
				err = fmt.Errorf("rollback failed: %w", rollbackErr)
			} else {
				err = errors.Join(err, fmt.Errorf("rollback failed: %w", rollbackErr))
			}
			s.log.Error(fmt.Sprintf("%s %s", op, err))
		}
	}()

	// Передаём транзакцию в контексте
	txCtx := context.WithValue(ctx, "tx", tx)
	found, err := s.giveawayFinder.Find(txCtx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFoundError("giveaway not found")
	} else if err != nil {
		return nil, domain.NewInternalServerError("failed to get giveaway", err)
	}
	return found, err
}

func (s *Service) Save(
	ctx context.Context,
	giveawaySaveRequest dto.GiveawaySaveRequest,
) (*dto.GiveawaySaveResponse, error) {
	const op = "service.giveaway.Save"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s %s", op, err))
		return nil, err
	}

	// Безопасный откат при любом исходе
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			if err == nil {
				err = fmt.Errorf("rollback failed: %w", rollbackErr)
			} else {
				err = errors.Join(err, fmt.Errorf("rollback failed: %w", rollbackErr))
			}
			s.log.Error(fmt.Sprintf("%s %s", op, err))
		}
	}()

	// Передаём транзакцию в контексте
	txCtx := context.WithValue(ctx, "tx", tx)
	saved, err := s.giveawaySaver.Save(txCtx, domain.Giveaway{
		Id:                uuid.New(),
		Status:            giveawaySaveRequest.Status,
		ChatId:            giveawaySaveRequest.ChatId,
		MessageId:         giveawaySaveRequest.MessageId,
		ShouldEditMessage: giveawaySaveRequest.ShouldEditMessage,
		ShouldSendNew:     giveawaySaveRequest.ShouldSendNew,
		WinCount:          giveawaySaveRequest.WinCount,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	})

	if err != nil {
		s.log.Error(fmt.Sprintf("%s %s", op, err))
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s %s", op, err))
		// todo: вернуть соотв ошибки
		return nil, err
	}

	return &dto.GiveawaySaveResponse{
		Id:                saved.Id,
		Status:            saved.Status,
		ChatId:            saved.ChatId,
		MessageId:         saved.MessageId,
		ShouldEditMessage: saved.ShouldEditMessage,
		ShouldSendNew:     saved.ShouldSendNew,
		WinCount:          saved.WinCount,
		CreatedAt:         saved.CreatedAt,
		UpdatedAt:         saved.UpdatedAt,
	}, nil
}
