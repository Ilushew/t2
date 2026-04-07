package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ilushew/udmurtia-trip/backend/internal/models"
)

func (r *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	query, args, err := r.psq.
		Select("id", "email", "is_verified", "is_admin").
		From("users").
		OrderBy("id DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	query, args, err := r.psq.
		Update("users").
		SetMap(updates).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteUser удаляет пользователя по ID (для админки)
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.psq.
		Delete("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, query, args...)
	return err
}
