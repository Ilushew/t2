# Udmurtia Trip — AI-маршрутизатор

Система для генерации маршрутов путешествий с использованием ML-моделей. Микросервисная архитектура на базе Go и nginx.

## Структура проекта

```
t2/
├── backend/                 # Go-бэкенд (Gin framework)
│   ├── cmd/
│   │   └── main.go          # Точка входа приложения
│   ├── internal/
│   │   ├── handlers/        # HTTP-обработчики
│   │   ├── ml_client/       # Клиент для взаимодействия с ML-сервисом
│   │   ├── models/          # Модели данных
│   │   ├── repository/      # Слой доступа к данным
│   │   └── services/        # Бизнес-логика
│   ├── templates/           # HTML-шаблоны (base, index, partials)
│   ├── static/              # Статические файлы (CSS, JS, изображения)
│   ├── Dockerfile           # Multi-stage сборка Go-приложения
│   ├── go.mod               # Зависимости Go
│   └── go.sum
├── ml_service/              # ML-сервис (в разработке)
│   └── Dockerfile
├── infra/
│   ├── nginx/
│   │   └── nginx.conf       # Конфигурация reverse proxy
│   └── db/                  # Резервировано для БД
├── docs/                    # Документация
├── docker-compose.yml       # Оркестрация сервисов
└── README.md
```

## Компоненты

### Backend (Go 1.25)
- **Framework:** Gin
- **Порт:** 8080 (внутри контейнера)
- **Функционал:**
  - Главная страница с формой генерации маршрута
  - POST `/generate` — генерация маршрута через ML-сервис
  - Раздача статики через `/static/`

### Nginx (Alpine)
- **Порт:** 80 (хост)
- **Функционал:**
  - Reverse proxy на backend
  - Раздача статических файлов с кэшированием (30 дней)
  - Логирование access/error

### ML Service
- Зарезервирован для будущего ML-сервиса

## Быстрый старт

### Требования
- Docker ≥ 20.x
- Docker Compose ≥ 1.29

### Запуск

```bash
# Сборка и запуск всех сервисов
docker-compose up -d --build

# Проверка статуса
docker-compose ps

# Логи в реальном времени
docker-compose logs -f

# Остановка
docker-compose down
```

После запуска приложение доступно по адресу: **http://localhost**

## Разработка

### Пересборка после изменений

```bash
docker-compose up -d --build
```

### Логи отдельных сервисов

```bash
docker-compose logs -f backend
docker-compose logs -f nginx
```

### Очистка

```bash
# Остановка + удаление контейнеров и сетей
docker-compose down --volumes

# Полная очистка (включая образы)
docker-compose down --volumes --rmi all
```

## Архитектура

```
┌─────────┐     ┌─────────┐     ┌─────────────┐
│  Client │ ──> │  Nginx  │ ──> │   Backend   │
│  :80    │     │  :80    │     │    :8080    │
└─────────┘     └─────────┘     └─────────────┘
                                       │
                                       ▼
                                 ┌─────────────┐
                                 │ ML Service  │
                                 │   (TODO)    │
                                 └─────────────┘
```

## API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/` | Главная страница |
| POST | `/generate` | Генерация маршрута |