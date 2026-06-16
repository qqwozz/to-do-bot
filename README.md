# To-Do Bot

Telegram бот для управления задачами с микросервисной архитектурой.

```
┌──────────────────┐     HTTP      ┌──────────────────┐     SQL      ┌──────────┐
│   Telegram Bot   │ ──────────── │   Backend API    │ ──────────── │  SQLite  │
│      (Go)        │              │    (C++/Crow)    │              │          │
└──────────────────┘              └──────────────────┘              └──────────┘
```

## Быстрый старт

```bash
cp .env.example .env   # добавь BOT_TOKEN
make docker-up          # запуск
```

## Команды бота

| Команда | Описание |
|---------|----------|
| `/start` | Главное меню |
| `/today` | Планы на сегодня |
| `/tomorrow` | Планы на завтра |
| `/week` | Планы на неделю |
| `/nextweek` | Планы на следующую неделю |

## API

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/plans` | Создать план |
| `GET` | `/plans?date=YYYY-MM-DD` | Планы на дату |
| `GET` | `/plans/range?start=&end=` | Планы за период |
| `DELETE` | `/plans/:id` | Удалить план |
| `GET` | `/health` | Проверка здоровья |

## Команды

```bash
make help            # все команды
make docker-up       # запуск через Docker
make docker-down     # остановка
make test-go         # Go тесты
make test-cpp-run    # C++ тесты
make bot-run         # запуск бота локально
make backend-run     # запуск backend локально
```

## Структура

```
├── cmd/bot/main.go              # точка входа Go
├── internal/
│   ├── bot/                     # логика бота
│   └── config/                  # конфигурация
├── cpp-backend/
│   ├── include/                 # заголовки
│   ├── src/                     # реализация
│   └── tests/                   # C++ тесты
├── .github/workflows/           # CI/CD
└── docs/                        # документация
```

## Тесты

```bash
make test       # все тесты
```

## Технологии

- **Go** — Telegram бот
- **C++** — Backend API (Crow)
- **SQLite** — база данных
- **Docker** — контейнеризация
- **GitHub Actions** — CI/CD
