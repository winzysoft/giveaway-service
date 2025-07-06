package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"giveaway-service/internal/config"
	"log/slog"
)

func MustInitDB(
	ctx context.Context,
	log *slog.Logger,
	pgCfg config.PostgresDBConfig,
) *pgxpool.Pool {
	pool, err := initDB(ctx, log, pgCfg)
	if err != nil {
		panic(err)
	}

	log.Info(fmt.Sprintf("Successfully connected to PostgreSQL"))
	return pool
}

func initDB(
	ctx context.Context,
	log *slog.Logger,
	pgCfg config.PostgresDBConfig,
) (*pgxpool.Pool, error) {
	const op = "repository.InitDB"
	log = log.With(
		slog.String("op", op),
	)

	//poolCfg, err := pgxpool.ParseConfig(fmt.Sprintf(
	//	"postgres://%s:%s@%s:%d/%s",
	//	pgCfg.UserName, pgCfg.Password, pgCfg.Host, pgCfg.Port, pgCfg.Database,
	//))
	//
	//poolCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
	//	conn.TypeMap().RegisterType(&pgtype.Type{
	//		Name:  "uuid",
	//		OID:   2950,
	//		Codec: &pgtype.UUIDCodec{},
	//	})
	//	return true
	//}
	//
	//poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	//	// Регистрируем тип UUID для автоматического распознавания
	//	conn.TypeMap().RegisterType(&pgtype.Type{
	//		Name:  "uuid",
	//		OID:   2950,
	//		Codec: &pgtype.UUIDCodec{},
	//	})
	//	return nil
	//}
	//
	//pool, err := pgxpool.NewWithConfig(ctx, poolCfg)

	// 1. Конфигурация пула с явной регистрацией UUID

	poolCfg, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s%s",
		pgCfg.UserName, pgCfg.Password, pgCfg.Host, pgCfg.Port, pgCfg.Database, "?sslmode=disable"))
	if err != nil {
		log.Error(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	//// 2. Важная настройка BeforeAcquire
	//poolCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
	//	conn.TypeMap().RegisterType(&pgtype.Type{
	//		Name:  "uuid",
	//		OID:   2950,
	//		Codec: &pgtype.UUIDCodec{},
	//	})
	//	pgxUUID.Register(conn.TypeMap())
	//	return true
	//}
	//
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "uuid",
			OID:   2950,
			Codec: &pgtype.UUIDCodec{},
		})
		return nil
	}

	// 3. Создаем пул
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)

	if err != nil {
		log.Error(fmt.Sprintf("Unable to connect to database: %v\n", err))
		return nil, err
	}

	err = pool.Ping(ctx)

	if err != nil {
		log.Error(fmt.Sprintf("Unable to connect to database: %v\n", err))
		return nil, err
	}

	return pool, nil
}
