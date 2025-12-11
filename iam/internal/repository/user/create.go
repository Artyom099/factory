package user

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/iam/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/iam/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, dto model.User) (string, error) {
	user := converter.ToRepoUser(dto)

	userID, err := insertUser(ctx, r.pool, user)
	if err != nil {
		return "", err
	}

	err = insertNotificationMethods(ctx, r.pool, userID, user.NotificationMethods)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func insertUser(ctx context.Context, pool *pgxpool.Pool, dto repoModel.RepoUser) (string, error) {
	builderInsertUser := sq.
		Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("login", "email", "hash").
		Values(dto.Login, dto.Email, dto.Hash).
		Suffix("RETURNING id")

	query, args, err := builderInsertUser.ToSql()
	if err != nil {
		log.Printf("failed to build insert user query: %v\n", err)
		return "", err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("failed to insert user: %v\n", err)
		return "", err
	}
	defer rows.Close()

	userID, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	if err != nil {
		log.Printf("failed to collect inserted user id: %v\n", err)
		return "", err
	}

	return userID, nil
}

func insertNotificationMethods(
	ctx context.Context,
	pool *pgxpool.Pool,
	userID string,
	methods []repoModel.RepoNotificationMethod,
) error {
	if len(methods) == 0 {
		return nil
	}

	builder := sq.
		Insert("notification_methods").
		Columns("user_id", "provider_name", "target").
		PlaceholderFormat(sq.Dollar)

	for _, nm := range methods {
		builder = builder.Values(userID, nm.ProviderName, nm.Target)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build insert notification_methods query: %v\n", err)
		return err
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to insert notification_methods: %v\n", err)
		return err
	}

	return nil
}
