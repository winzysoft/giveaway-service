package giveaway

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"giveaway-service/internal/domain"
	"log/slog"
)

const TableName = "giveaways"

type Repository struct {
	log  *slog.Logger
	psql *squirrel.StatementBuilderType
}

func (r *Repository) Find(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Giveaway, error) {
	const op = "repository.giveaway.Find"

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	sql, args, err := r.psql.Select("*").From(TableName).Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		r.log.Error(fmt.Sprintf("%s %s", op, err))
		return nil, err
	}

	var giveaway domain.Giveaway
	// Выполняем запрос
	err = tx.QueryRow(ctx, sql, args...).Scan(
		&giveaway.Id,
		&giveaway.Status,
		&giveaway.ChatId,
		&giveaway.MessageId,
		&giveaway.ShouldEditMessage,
		&giveaway.ShouldSendNew,
		&giveaway.WinCount,
		&giveaway.CreatedAt,
		&giveaway.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.log.Info(fmt.Sprintf("%s not found giveaway with id: %s", op, id))
		return nil, err
	}
	if err != nil {
		r.log.Error(fmt.Sprintf("%s %s", op, err))
		return nil, err
	}
	return &giveaway, nil
}

func (r *Repository) Save(
	ctx context.Context,
	giveaway domain.Giveaway,
) (*domain.Giveaway, error) {
	const op = "repository.giveaway.Save"

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	insertSql, args, err := r.psql.Insert(TableName).
		Columns("id",
			"status",
			"chat_id",
			"message_id",
			"should_edit_message",
			"should_send_new",
			"win_count",
			"created_at",
			"updated_at").
		Suffix("RETURNING id").
		Values(giveaway.Id,
			giveaway.Status,
			giveaway.ChatId,
			giveaway.MessageId,
			giveaway.ShouldEditMessage,
			giveaway.ShouldSendNew,
			giveaway.WinCount,
			giveaway.CreatedAt,
			giveaway.UpdatedAt).
		ToSql()

	if err != nil {
		r.log.Error(fmt.Sprintf("%s %s", op, err))
		return nil, err
	}

	// Выполняем запрос
	_, err = tx.Exec(ctx, insertSql, args...)
	if err != nil {
		r.log.Error(fmt.Sprintf("%s %s", op, err))
		return nil, err
	}

	return &giveaway, nil
}

func (r *Repository) Update(
	ctx context.Context,
) (*domain.Giveaway, error) {
	// todo:
	return nil, nil
}

func (r *Repository) Delete(
	ctx context.Context,
) error {
	// todo:
	return nil
}

func New(
	log *slog.Logger,
	psql *squirrel.StatementBuilderType,
) *Repository {
	return &Repository{
		log:  log,
		psql: psql,
	}
}
