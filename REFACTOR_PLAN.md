# План рефакторинга backend

> Создан: 3 апреля 2026
> Статус: Ожидает утверждения приоритетов

---

## 🔴 КРИТИЧНО — Безопасность и стабильность

### 1. Хардкод session secret key
**Файл:** `cmd/main.go` (строка ~72)

```go
// БЫЛО (неправильно):
secret := "your-super-secret-key-min-32-characters-long"
if secret == "" {
    log.Fatal("SESSION_SECRET environment variable is required")
}

// СТАНЕТ:
secret := os.Getenv("SESSION_SECRET")
if secret == "" {
    log.Fatal("SESSION_SECRET environment variable is required")
}
```

**Риск:** Злоумышленник может подделать cookie-сессию.

---

### 2. SSL отключён для PostgreSQL
**Файл:** `pkg/postgres/pool.go` (строка ~17)

```go
// БЫЛО:
"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"

// СТАНЕТ:
"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"
```

Добавить поле `SSLMode` в `postgres.Config`. По умолчанию `prefer` для dev, `require` для prod.

---

### 3. Redis без пароля и настройки
**Файл:** `pkg/redis/client.go`, `cmd/main.go`

```go
// БЫЛО:
redisAddr := "redis:6379"
redis.NewClient(&redis.Options{Addr: addr})

// СТАНЕТ:
redisOptions := redisPkg.Options{
    Addr:     os.Getenv("REDIS_ADDR"),
    Password: os.Getenv("REDIS_PASSWORD"),
    DB:       config.GetInt("REDIS_DB"),
}
```

---

### 4. Нет graceful shutdown
**Файл:** `cmd/main.go`

```go
// БЫЛО:
r.Run(":8080")

// СТАНЕТ:
srv := &http.Server{Addr: ":8080", Handler: r}
go func() { srv.ListenAndServe() }()
// await SIGINT/SIGTERM
srv.Shutdown(ctx)
```

---

### 5. Redis клиент не закрывается
**Файл:** `cmd/main.go`

Добавить `defer redisClient.Close()` после создания.

---

## 🟡 АРХИТЕКТУРА

### 1. Вынести инициализацию из `main()`

**Файл:** `cmd/main.go`

Создать `internal/app/app.go`:
```go
type App struct {
    Router     *gin.Engine
    UserRepo   UserRepository
    AuthSvc    AuthService
    CodeSvc    CodeService
    EmailSvc   EmailService
    DB         *pgxpool.Pool
    Redis      *redis.Client
}

func Initialize(cfg *Config) (*App, error) { ... }
```

`main()` станет тонкой обёрткой:
```go
func main() {
    app, err := app.Initialize(cfg)
    if err != nil { log.Fatal(err) }
    app.Run()
}
```

---

### 2. Ввести интерфейсы для зависимостей

**Новый файл:** `internal/handlers/interfaces.go` (или в каждом пакете свой интерфейс)

```go
type UserRepository interface {
    Create(ctx context.Context, email string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type CodeService interface {
    Generate(ctx context.Context, email string) (string, error)
    Verify(ctx context.Context, email, code string) (bool, error)
}

type EmailService interface {
    SendVerificationCode(ctx context.Context, email, code string) error
}
```

Handlers будут зависеть от интерфейсов, а не от конкретных типов.

---

### 3. TripHandler — вынести ML-сервис

**Файл:** `internal/handlers/routes.go`

```go
// БЫЛО:
func NewTripHandler() *TripHandler {
    return &TripHandler{mlService: services.NewMockMLService()}
}

// СТАНЕТ:
func NewTripHandler(mlService MLService) *TripHandler {
    return &TripHandler{mlService: mlService}
}
```

В `main.go`: передавать зависимость явно.

---

### 4. `pkg/config` — использовать или удалить

Сейчас пакет импортируется, но не используется. Варианты:

- **Использовать:** заменить все `os.Getenv()` на `config.Get()` / `config.MustGet()`
- **Удалить:** если `os.Getenv()` достаточно для проекта

Также исправить баг в `MustGet`:
```go
// БЫЛО:
panic(fmt.Sprintf("config.MustGet: %s reqiuired variable is not set", val))

// СТАНЕТ:
panic(fmt.Sprintf("config.MustGet: %s required variable is not set", key))
```

---

### 5. Глобальное состояние в migrator

**Файл:** `pkg/migrator/migrator.go`

`goose.SetBaseFS`, `goose.SetTableName`, `goose.SetDialect` — глобальные вызовы.
При параллельных вызовах — race condition.

Решение: создавать экземпляр goose с параметрами или использовать мьютекс.

---

## 🟢 КАЧЕСТВО КОДА

### 1. `scanUser` скрывает ошибки БД

**Файл:** `internal/repository/user_repository.go`

```go
// БЫЛО:
if err != nil {
    return nil, ErrUserNotFound
}

// СТАНЕТ:
if errors.Is(err, pgx.ErrNoRows) {
    return nil, ErrUserNotFound
}
if err != nil {
    return nil, fmt.Errorf("scanUser: %w", err)
}
```

---

### 2. Email в горутине — добавить контекст и retry

**Файл:** `internal/handlers/auth.go`

```go
// БЫЛО:
go func() {
    err = h.emailSvc.SendVerificationCode(email, code)
    ...
}()

// СТАНЕТ:
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    for attempt := 0; attempt < 3; attempt++ {
        err := h.emailSvc.SendVerificationCode(ctx, email, code)
        if err == nil { return }
        log.Printf("Email send attempt %d failed: %v", attempt+1, err)
        time.Sleep(time.Second * time.Duration(attempt))
    }
}()
```

---

### 3. Race condition в `Register`

**Файл:** `internal/handlers/auth.go`

```sql
-- Вместо SELECT + INSERT:
INSERT INTO users (email, created_at, updated_at)
VALUES ($1, NOW(), NOW())
ON CONFLICT (email) DO NOTHING
RETURNING id, email, created_at, updated_at;
```

---

### 4. Валидация email

**Файл:** `internal/handlers/auth.go`

Добавить проверку формата email перед обработкой.
Можно использовать `net/mail` пакет или регулярное выражение.

---

### 5. `crypto/rand` для кодов

**Файл:** `internal/handlers/auth.go`

```go
// БЫЛО:
code := fmt.Sprintf("%06d", rand.IntN(1000000))

// СТАНЕТ:
func generateCode() (string, error) {
    b := make([]byte, 3)
    _, err := crypto/rand.Read(b)
    if err != nil { return "", err }
    return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2]) % 1000000), nil
}
```

---

### 6. SMTP — постоянное соединение

**Файл:** `internal/services/email.go`

`DialAndSendWithContext` создаёт новое соединение каждый раз.
Варианты:
- Использовать `mail.Dialer` с `Client()` для повторного использования
- Пул SMTP-соединений
- Перейти на сторонний email-сервис (SendGrid, Resend и т.д.)

---

## 📊 Сводка по файлам

| Файл | Критичность | Кол-во проблем |
|------|-------------|----------------|
| `cmd/main.go` | 🔴 | 5 |
| `pkg/postgres/pool.go` | 🔴 | 1 |
| `pkg/redis/client.go` | 🔴 | 2 |
| `internal/handlers/auth.go` | 🟡 | 4 |
| `internal/handlers/routes.go` | 🟡 | 2 |
| `internal/repository/user_repository.go` | 🟢 | 3 |
| `internal/services/email.go` | 🟢 | 2 |
| `internal/services/code_service.go` | 🟢 | 2 |
| `pkg/config/config.go` | 🟡 | 2 |
| `pkg/migrator/migrator.go` | 🟡 | 2 |

---

## ✅ Следующие шаги

1. Утвердить приоритет (A — критичные, B — полный рефакторинг, C — выборочно)
2. Создать ветку `refactor/cleanup-and-security`
3. Начать с 🔴 критичных проблем
4. После каждого пункта — запуск тестов и линтера
