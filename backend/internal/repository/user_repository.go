package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCode       = errors.New("invalid verification code")
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
	var code sql.NullString
	var expiresAt sql.NullTime

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&code,
		&expiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if code.Valid {
		user.EmailVerificationCode = &code.String
	} else {
		user.EmailVerificationCode = nil
	}

	if expiresAt.Valid {
		user.EmailVerificationExpiresAt = &expiresAt.Time
	} else {
		user.EmailVerificationExpiresAt = nil
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	id := uuid.New()
	now := time.Now()

	query, args, err := r.psq.
		Insert("users").
		Columns("id", "email", "password_hash", "is_verified", "created_at", "updated_at").
		Values(id, email, passwordHash, false, now, now).
		Suffix("RETURNING id, email, password_hash, is_verified, email_verification_code, email_verification_expires_at, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.EmailVerificationCode,
		&user.EmailVerificationExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
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
		Select("id", "email", "password_hash", "is_verified", "email_verification_code", "email_verification_expires_at", "created_at", "updated_at").
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
		Select("id", "email", "password_hash", "is_verified", "email_verification_code", "email_verification_expires_at", "created_at", "updated_at").
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

func (r *UserRepository) VerifyCode(ctx context.Context, id uuid.UUID, code string) error {
	query, args, err := r.psq.
		Select("COUNT(1)").
		From("users").
		Where(squirrel.And{
			squirrel.Eq{"id": id},
			squirrel.Eq{"email_verification_code": code},
			squirrel.GtOrEq{"email_verification_expires_at": time.Now()},
		}).ToSql()
	if err != nil {
		return err
	}
	var count int
	err = r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidCode
	}
	return nil
}

// UpdateVerificationStatus обновляет статус верификации email
func (r *UserRepository) UpdateVerificationStatus(ctx context.Context, id uuid.UUID, isVerified bool) error {
	query, args, err := r.psq.
		Update("users").
		Set("is_verified", isVerified).
		Set("updated_at", time.Now()).
		Set("email_verification_code", nil).
		Set("email_verification_expires_at", nil).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}

// SetVerificationCode устанавливает код верификации и время его истечения
func (r *UserRepository) SetVerificationCode(ctx context.Context, id uuid.UUID, code string, expiresAt time.Time) error {
	query, args, err := r.psq.
		Update("users").
		Set("email_verification_code", code).
		Set("email_verification_expires_at", expiresAt).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}
