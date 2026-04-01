package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ilushew/udmurtia-trip/backend/internal/repository"
	"github.com/ilushew/udmurtia-trip/backend/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()
	cfg := postgres.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Database: getEnv("DB_NAME", "udmurtia_trip_test"),
	}

	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		t.Skipf("Skipping test, cannot connect to database: %v", err)
	}

	return pool
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func TestUserRepository_CreateUser(t *testing.T) {
	pool := setupTestPool(t)
	defer pool.Close()

	ctx := context.Background()
	repo := repository.NewUserRepository(pool)

	email := "test@example.com"
	passwordHash := "hashed_password_123"

	user, err := repo.CreateUser(ctx, email, passwordHash)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	if user.IsVerified {
		t.Error("Expected IsVerified to be false")
	}

	if user.PasswordHash != passwordHash {
		t.Errorf("Expected password hash %s, got %s", passwordHash, user.PasswordHash)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	pool := setupTestPool(t)
	defer pool.Close()

	ctx := context.Background()
	repo := repository.NewUserRepository(pool)

	email := "findbyme@example.com"
	passwordHash := "hashed_password_456"

	// Создаём пользователя
	_, err := repo.CreateUser(ctx, email, passwordHash)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Ищем по email
	user, err := repo.FindByEmail(ctx, email)
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	pool := setupTestPool(t)
	defer pool.Close()

	ctx := context.Background()
	repo := repository.NewUserRepository(pool)

	_, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	if err != repository.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	pool := setupTestPool(t)
	defer pool.Close()

	ctx := context.Background()
	repo := repository.NewUserRepository(pool)

	email := "findbyid@example.com"
	passwordHash := "hashed_password_789"

	// Создаём пользователя
	created, err := repo.CreateUser(ctx, email, passwordHash)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Ищем по ID
	user, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if user.ID != created.ID {
		t.Errorf("Expected ID %v, got %v", created.ID, user.ID)
	}
}

func TestUserRepository_UpdateVerificationStatus(t *testing.T) {
	pool := setupTestPool(t)
	defer pool.Close()

	ctx := context.Background()
	repo := repository.NewUserRepository(pool)

	email := "verify@example.com"
	passwordHash := "hashed_password"

	// Создаём пользователя
	created, err := repo.CreateUser(ctx, email, passwordHash)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Верифицируем
	err = repo.UpdateVerificationStatus(ctx, created.ID, true)
	if err != nil {
		t.Fatalf("UpdateVerificationStatus failed: %v", err)
	}

	// Проверяем
	user, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if !user.IsVerified {
		t.Error("Expected IsVerified to be true")
	}
}

func TestUserRepository_SetVerificationCode(t *testing.T) {
	pool := setupTestPool(t)
	defer pool.Close()

	ctx := context.Background()
	repo := repository.NewUserRepository(pool)

	email := "verifocode@example.com"
	passwordHash := "hashed_password"

	// Создаём пользователя
	created, err := repo.CreateUser(ctx, email, passwordHash)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Устанавливаем код
	code := "123456"
	expiresAt := time.Now().Add(15 * time.Minute)
	err = repo.SetVerificationCode(ctx, created.ID, code, expiresAt)
	if err != nil {
		t.Fatalf("SetVerificationCode failed: %v", err)
	}

	// Проверяем
	user, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if user.EmailVerificationCode == nil {
		t.Fatal("Expected EmailVerificationCode to be set")
	}

	if *user.EmailVerificationCode != code {
		t.Errorf("Expected code %s, got %s", code, *user.EmailVerificationCode)
	}
}
