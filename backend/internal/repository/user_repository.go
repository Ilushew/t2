package repository

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository struct {
	pool *pgxpool.Pool
	psq  squirrel.StatementBuilderType
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
		psq:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// rowScanner интерфейс для сканирования строки
type rowScanner interface {
	Scan(dest ...any) error
}

// scanUser сканирует результат запроса в структуру User
func scanUser(row rowScanner) (*models.User, error) {
	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.IsVerified,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, email string) (*models.User, error) {
	id := uuid.New()

	query, args, err := r.psq.
		Insert("users").
		Columns("id", "email", "is_verified").
		Values(id, email, false).
		Suffix("RETURNING id, email, is_verified").
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.IsVerified,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query, args, err := r.psq.
		Select("id", "email", "is_verified").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	user, err := scanUser(r.pool.QueryRow(ctx, query, args...))
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query, args, err := r.psq.
		Select("id", "email", "is_verified").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	user, err := scanUser(r.pool.QueryRow(ctx, query, args...))
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
