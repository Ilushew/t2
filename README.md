# 🗺️ Udmurtia Trip — AI-маршрутизатор путешествий

> Интеллектуальная система генерации персональных туристических маршрутов по Удмуртии с использованием ML-рекомендаций и passwordless-авторизацией.

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.12-008080?logo=gin)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker)](https://www.docker.com/)

---

##  О проекте

**Udmurtia Trip** — это микросервисное веб-приложение, которое помогает пользователям открывать для себя Удмуртию через персонализированные маршруты. Система учитывает предпочтения путешественника (природа, культура, гастрономия, семейный отдых) и формирует оптимальный маршрут с помощью ML-модели.

### ✨ Ключевые возможности

-  **AI-генерация маршрутов** — интеллектуальный подбор мест на основе предпочтений
-  **Passwordless-авторизация** — вход по коду подтверждения через email без паролей
-  **Email-верификация** — безопасное подтверждение почты через одноразовые коды
-  **Современный UI** — адаптивный интерфейс с HTMX для динамической подгрузки контента
-  **Чистая архитектура** — разделение на слои (handlers → services → repository)
- **Docker-first** — полный запуск одной командой

---

##  Архитектура

### Микросервисная топология

```
┌──────────────┐
│   Client     │
│  (Browser)   │
└──────┬───────┘
       │ HTTP
       ▼
┌─────────────────────────────────────────────┐
│              Nginx (Reverse Proxy)          │
│  • Маршрутизация запросов                   │
│  • Раздача статики (кэш 30 дней)            │
│  • SSL termination (готов к production)     │
└──────┬──────────────────┬───────────────────┘
       │                  │
       ▼                  ▼
┌──────────────┐   ┌──────────────┐
│   Backend    │   │   Static     │
│   (Go/Gin)   │   │   Files      │
│   :8080      │   │  (CSS, JS)   │
└──────┬───────┘   └──────────────┘
       │
       ▼
┌──────────────────────────────────────────────┐
│            Data Layer                        │
│                                              │
│  ┌──────────────┐      ┌──────────────────┐  │
│  │  PostgreSQL  │      │      Redis       │  │
│  │  • Users     │      │  • Session data  │  │
│  │  • Migrations│      │  • Auth codes    │  │
│  └──────────────┘      └──────────────────┘  │
└──────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────┐
│         ML Service                            │
│  • Рекомендации мест                          │
│  • Фильтрация по категориям                   │
│  • Расширение: реальная ML-модель             │
└──────────────────────────────────────────────┘
```

### Структура проекта

```
t2/
├── backend/                     # Go-бэкенд (основное приложение)
│   ├── cmd/main.go              # Точка входа
│   ├── internal/                # Приватный код
│   │   ├── app/                 # Инициализация приложения
│   │   ├── handlers/            # HTTP-обработчики (controllers)
│   │   ├── services/            # Бизнес-логика (email, ML mock)
│   │   ├── repository/          # Работа с БД (SQL через Squirrel)
│   │   ├── models/              # Модели данных (User, Place)
│   │   ├── migrations/          # SQL миграции (goose embed.FS)
│   │   └── ml_client/           # Клиент для ML-сервиса
│   ├── pkg/                     # Публичные пакеты
│   │   ├── config/              # Загрузка .env конфигурации
│   │   ├── postgres/            # PostgreSQL connection pool
│   │   └── migrator/            # Запуск goose миграций
│   ├── templates/               # HTML шаблоны (Gin multitemplate)
│   ├── static/                  # CSS, JS, изображения
│   ├── Dockerfile               # Multi-stage сборка
│   └── go.mod                   # Go зависимости
├── ml_service/                  # ML-сервис
├── infra/
│   ├── nginx/nginx.conf         # Конфигурация reverse proxy
│   └── db/                      # Резерв для БД
├── docs/                        # Документация
├── docker-compose.yml           # Оркестрация 4 сервисов
├── taskfile.yml                 # Task-команды (аналог Makefile)
└── .env.example                 # Шаблон переменных окружения
```

---

## Технологический стек

### Backend (Go 1.25)

| Технология | Назначение |
|------------|------------|
| **[Gin](https://gin-gonic.com/)** | Высокопроизводительный HTTP-фреймворк |
| **[pgx/v5](https://github.com/jackc/pgx)** | PostgreSQL драйвер + connection pool |
| **[Squirrel](https://github.com/Masterminds/squirrel)** | SQL query builder для безопасных запросов |
| **[Goose](https://github.com/pressly/goose)** | Управление миграциями БД (embed.FS) |
| **[go-redis/v9](https://github.com/redis/go-redis)** | Redis клиент для сессий и кодов |
| **[go-mail](https://github.com/wneessen/go-mail)** | Отправка email через SMTP |
| **[HTMX](https://htmx.org/)** | Динамическая подгрузка контента без JS |
| **[Sessions](https://github.com/gin-contrib/sessions)** | Управление сессиями (securecookie) |
| **[UUID](https://github.com/google/uuid)** | Генерация уникальных ID |

### Инфраструктура

| Сервис | Версия | Назначение |
|--------|--------|------------|
| **PostgreSQL** | 16 Alpine | Основная БД (users, verification codes) |
| **Redis** | 7 Alpine | Хранение auth кодов, сессии |
| **Nginx** | Alpine | Reverse proxy, раздача статики |

### DevOps & Tooling

- **Docker Compose** — оркестрация всех сервисов
- **Taskfile** — удобные команды сборки (`task build`, `task test`, `task lint`)
- **golangci-lint** — линтинг и автофикс кода
- **Multi-stage Docker build** — оптимизированный размер образа

---

## Система авторизации

Приложение использует **passwordless** подход — никаких паролей, только email-коды:

```
1. POST /auth/request-code  → Введите email → Получите код
2. POST /auth/verify        → Введите код   → Получите сессию
3. GET  /auth/me            → Проверка профиля (защищённый роут)
4. POST /auth/logout        → Завершение сессии
```

**Безопасность:**
- 🔒 Одноразовые коды с TTL (хранятся в Redis)
- 🔒 Secure cookie для сессий
- 🔒 Email верификация обязательна для доступа

---

## База данных

### Схема

```sql
-- Пользователи
users (
    id UUID PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    is_verified BOOLEAN DEFAULT false
)

-- Миграции отслеживаются через goose
schema_migrations (version, is_applied, ...)
```

### Миграции

Автоматический запуск при старте приложения через `embed.FS`:

```go
// internal/migrations/migrations.go
//go:embed *.sql
var migrationFS embed.FS
```

---

## Быстрый старт

### Требования

| Зависимость | Версия |
|-------------|--------|
| [Docker](https://www.docker.com/get-started) | ≥ 20.x |
| [Docker Compose](https://docs.docker.com/compose/) | ≥ 2.x |
| [Task](https://taskfile.dev/) (опционально) | ≥ 3.x |

### 1️⃣ Клонирование репозитория

```bash
git clone https://github.com/Ilushew/t2
cd t2
```

### Настройка окружения

```bash
# Скопируйте шаблон
cp .env.example .env

# Отредактируйте .env под вашу конфигурацию
# Обязательно укажите SMTP данные для email-сервиса
```

**Пример `.env`:**

```env
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=udmurtia_trip

# Email (SMTP)
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=465
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_FROM=your-email@gmail.com

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=redis-secret-password
REDIS_DB=0

# Sessions
SESSION_SECRET=your-super-secret-key
```

### Запуск через Docker Compose

```bash
# Сборка и запуск всех сервисов
docker compose up -d --build

# Проверка статуса
docker compose ps

# Логи в реальном времени
docker compose logs -f
```

 **Готово!** Приложение доступно по адресу: **http://localhost**

### Альтернатива: Запуск через Task

```bash
task up          # Аналог docker compose up --build -d
task logs        # Просмотр логов
task down        # Остановка
```

---

## Локальная разработка

### Сборка бинарника (Windows)

```bash
task build
# или вручную:
cd backend
go build -o backend.exe ./cmd/main.go
```

### Запуск тестов

```bash
task test              # Все тесты
task test-repo         # Тесты репозитория (требует БД)
task test-cover        # С отчётом покрытия
```

### Линтинг и форматирование

```bash
task fmt               # go fmt ./...
task lint              # golangci-lint run
task lint-fix          # Автофикс проблем
task check             # Полный чек: fmt + lint + test
```

### Работа с БД

```bash
task db-up             # Запуск только PostgreSQL
task db-migrate        # Применение миграций
task db-reset          # Полный сброс БД
```

### Горячая перезагрузка

```bash
# Установите air: go install github.com/air-verse/air@latest
task dev               # Запуск с hot reload
```

---

## API Reference

### Публичные эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/` | Главная страница с формой генерации |
| `POST` | `/generate` | Генерация маршрута (HTMX) |

### Авторизация

| Метод | Путь | Описание | Тело запроса |
|-------|------|----------|--------------|
| `POST` | `/auth/request-code` | Запрос кода на email | `email=user@example.com` |
| `POST` | `/auth/verify` | Верификация кода | `email=...&code=123456` |
| `GET` | `/auth/me` | Получить профиль | Требуется сессия |
| `POST` | `/auth/logout` | Выход из системы | Требуется сессия |


## Тестирование

### Запуск тестовой БД

```bash
task test-db-up        # Создание тестовой базы
task test-repo         # Тесты репозитория
```

### Покрытие кода

```bash
task test-cover

# Пример вывода:
# ok      github.com/ilushew/udmurtia-trip/backend/internal/repository    0.45s
# coverage: 85.2% of statements
```

---

## Docker Production Ready

### Multi-stage сборка

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/bin/server ./cmd/main.go

# Stage 2: Run (minimal image)
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/server .
EXPOSE 8080
CMD ["./server"]
```

**Результат:** ~15MB вместо ~500MB

### Health Checks

- **PostgreSQL:** `pg_isready` каждые 5 секунд
- **Redis:** `redis-cli ping` каждые 5 секунд
- **Backend:** graceful shutdown при SIGINT/SIGTERM

---