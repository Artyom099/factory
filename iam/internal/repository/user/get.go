package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/iam/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuidOrLogin string) (model.User, error) {
	dbUser, err := getUser(ctx, r.pool, uuidOrLogin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	dbMethods, err := getNotificationMethods(ctx, r.pool, dbUser.ID)
	if err != nil {
		return model.User{}, err
	}

	dbUser.NotificationMethods = dbMethods

	return converter.ToModelUser(dbUser), nil
}

func getUser(ctx context.Context, pool *pgxpool.Pool, uuidOrLogin string) (repoModel.RepoUser, error) {
	builderUser := sq.
		Select("id", "login", "email", "hash", "created_at", "updated_at").
		From("users").
		Where(
			sq.Or{
				sq.Eq{"id": uuidOrLogin}, // todo - здесь ошибка - invalid input syntax for type uuid - задать вопрос в чате
				sq.Eq{"login": uuidOrLogin},
			},
		).
		PlaceholderFormat(sq.Dollar)

	queryUser, argsUser, err := builderUser.ToSql()
	if err != nil {
		return repoModel.RepoUser{}, err
	}

	rowsUser, err := pool.Query(ctx, queryUser, argsUser...)
	if err != nil {
		return repoModel.RepoUser{}, err
	}
	defer rowsUser.Close()

	dbUser, err := pgx.CollectOneRow(rowsUser, pgx.RowToStructByName[repoModel.RepoUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoModel.RepoUser{}, model.ErrUserNotFound
		}
		return repoModel.RepoUser{}, err
	}

	return dbUser, nil
}

func getNotificationMethods(ctx context.Context, pool *pgxpool.Pool, userID string) ([]repoModel.RepoNotificationMethod, error) {
	builderMethods := sq.
		Select("id", "user_id", "provider_name", "target").
		From("notification_methods").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	queryNm, argsNm, err := builderMethods.ToSql()
	if err != nil {
		return []repoModel.RepoNotificationMethod{}, err
	}

	rowsNm, err := pool.Query(ctx, queryNm, argsNm...)
	if err != nil {
		return []repoModel.RepoNotificationMethod{}, err
	}
	defer rowsNm.Close()

	dbMethods, err := pgx.CollectRows(rowsNm, pgx.RowToStructByName[repoModel.RepoNotificationMethod])
	if err != nil {
		return []repoModel.RepoNotificationMethod{}, err
	}

	return dbMethods, nil
}
