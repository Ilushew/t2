# Структура проекта

```
t2/
├── .env                          # Переменные окружения (не коммитить)
├── .env.example                  # Шаблон переменных окружения
├── .gitignore                    # Git ignore rules
├── docker-compose.yml            # Docker Compose конфигурация
├── Taskfile.yml                  # Taskfile команды (task run build/test/etc)
├── README.md                     # Документация проекта
├── STRUCTURE.md                  # Этот файл
│
├── infra/                        # Инфраструктура
│   └── nginx/
│       └── nginx.conf            # Nginx конфигурация
│
├── docs/                         # Документация
│
├── ml_service/                   # ML сервис (отдельный проект)
│   └── Dockerfile
│
└── backend/                      # Go бэкенд
    ├── go.mod                    # Go модуль зависимости
    ├── go.sum                    # Go зависимости (auto-generated)
    ├── Dockerfile                # Docker образ бэкенда
    ├── .golangci.yml             # GolangCI-Lint конфигурация
    │
    ├── cmd/                      # Точки входа приложения
    │   └── main.go               # Главный файл приложения
    │
    ├── internal/                 # Внутренний код приложения (private)
    │   ├── app/                  # Инициализация приложения
    │   │   └── app.go
    │   │
    │   ├── handlers/             # HTTP хендлеры (контроллеры)
    │   │   └── routes.go         # Маршруты и хендлеры
    │   │
    │   ├── migrations/           # SQL миграции (goose)
    │   │   ├── 0001_create_user_table.sql
    │   │   └── migrations.go     # embed.FS для миграций
    │   │
    │   ├── models/               # Модели данных
    │   │   └── user.go           # Модель пользователя
    │   │
    │   ├── repository/           # Работа с БД (SQL запросы)
    │   │   ├── user_repository.go
    │   │   └── user_repository_test.go
    │   │
    │   ├── services/             # Бизнес-логика
    │   │   ├── email.go          # Email сервис
    │   │   └── ml_mock.go        # Mock ML сервиса
    │   │
    │   └── ml_client/            # ML клиент (внешний сервис)
    │
    ├── pkg/                      # Публичные пакеты (shared code)
    │   ├── config/
    │   │   └── config.go         # Конфигурация (env variables)
    │   │
    │   ├── migrator/
    │   │   └── migrator.go       # Goose миграции
    │   │
    │   └── postgres/
    │       └── pool.go           # PostgreSQL connection pool
    │
    ├── static/                   # Статические файлы (CSS, JS, images)
    │
    └── templates/                # HTML шаблоны
        ├── base.html
        ├── index.html
        └── partials/
            └── route-result.html
```

## 📁 Описание папок

| Папка | Назначение |
|-------|------------|
| `cmd/` | Точки входа (main.go) |
| `internal/` | Приватный код приложения (нельзя импортировать извне) |
| `internal/handlers/` | HTTP обработчики запросов |
| `internal/models/` | Структуры данных (DTO, domain models) |
| `internal/repository/` | Доступ к данным (SQL запросы через Squirrel) |
| `internal/services/` | Бизнес-логика |
| `internal/migrations/` | SQL миграции базы данных |
| `pkg/` | Публичные пакеты (можно переиспользовать) |
| `pkg/postgres/` | Подключение к PostgreSQL |
| `pkg/migrator/` | Запуск миграций |
| `pkg/config/` | Работа с конфигурацией |
| `static/` | Статика (CSS, JS, изображения) |
| `templates/` | HTML шаблоны для Gin |

## 🚀 Основные команды

```bash
task up              # Запустить всё (Docker Compose)
task down            # Остановить Docker Compose
task build           # Собрать бинарник
task test            # Запустить тесты
task test-repo       # Тесты репозитория
task lint            # Линтер
task db-up           # Запустить PostgreSQL
task db-migrate      # Применить миграции
task logs            # Логи Docker
```

## 📦 Зависимости

### Go модули:
- `github.com/gin-gonic/gin` — HTTP фреймворк
- `github.com/jackc/pgx/v5` — PostgreSQL драйвер
- `github.com/Masterminds/squirrel` — SQL query builder
- `github.com/pressly/goose/v3` — Миграции БД
- `github.com/google/uuid` — UUID генерация

### Docker:
- `postgres:16-alpine` — База данных
- `nginx:alpine` — Reverse proxy

## 🔐 Авторизация (Passwordless)

```
POST /auth/request-code  → Отправка кода на email
POST /auth/verify        → Проверка кода → JWT-токен
GET  /auth/me            → Получить текущий профиль (JWT)
POST /auth/logout        → Выйти (invalidate JWT)
```

## 📊 База данных

**Таблицы:**
- `users` — пользователи (email, verification code, is_verified)
- `schema_migrations` — применённые миграции (goose)

**Миграции:**
- `0001_create_user_table.sql` — создание таблицы users
