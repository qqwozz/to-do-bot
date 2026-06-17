<div align="center">

# To-Do Bot

### Telegram бот для управления задачами

[![CI](https://github.com/qqwozz/to-do-bot/actions/workflows/ci.yml/badge.svg)](https://github.com/qqwozz/to-do-bot/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://go.dev/)
[![C++](https://img.shields.io/badge/C++-17-00599C?logo=cplusplus)](https://isocpp.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

```
┌──────────────────┐     HTTP      ┌──────────────────┐     SQL      ┌──────────┐
│   Telegram Bot   │ ──────────── │   Backend API    │ ──────────── │  SQLite  │
│      (Go)        │              │    (C++/Crow)    │              │          │
└──────────────────┘              └──────────────────┘              └──────────┘
```

</div>

---

## Быстрый старт

```bash
git clone https://github.com/qqwozz/to-do-bot.git
cd to-do-bot
cp .env.example .env
# добавь BOT_TOKEN в .env
make docker-up
```

## Интерфейс

Бот работает полностью через **inline-кнопки** — никаких команд вводить не нужно.

### Главное меню

```
👋 To-Do бот

Выберите действие:

[➕ Создать план]  [📋 Мои планы]
```

### Просмотр планов

```
📅 Выберите период:

[📅 Сегодня]    [📅 Завтра]
[📆 Эта неделя] [📆 След. неделя]
[◀️ Назад]
```

### Планы на день

```
📋 Сегодня

1. Встреча
   Обсудить проект
   🕐 14:00

2. Дедлайн
   Сдать отчёт
   ⏰ Весь день

[➕ Создать план]  [◀️ Назад]
```

### Создание плана (пошагово)

```
📝 Название:        → Встреча
📝 Описание:        → Обсудить проект
🗓 Дата:            → [Сегодня] [Завтра] [Через неделю] [Ввести дату]
⏰ Время:           → [🌅 Утро] [☀️ День] [🌆 Вечер] [⏰ Весь день]

✅ План создан!
```

## Архитектура

### Структура проекта

```
.
├── app/                          # исходный код
│   ├── cmd/bot/main.go           # точка входа
│   ├── internal/
│   │   ├── bot/
│   │   │   ├── bot.go            # Bot, New(), Run()
│   │   │   ├── handlers.go       # роутинг callback'ов
│   │   │   ├── keyboards.go      # inline-кнопки
│   │   │   ├── state.go          # state machine диалогов
│   │   │   ├── client.go         # HTTP к backend
│   │   │   ├── models.go         # Plan, PlanRequest
│   │   │   └── format.go         # форматирование
│   │   └── config/
│   │       └── config.go         # конфигурация
│   ├── cpp-backend/              # C++ backend
│   │   ├── include/              # заголовки
│   │   ├── src/                  # реализация
│   │   └── tests/                # тесты
│   ├── go.mod
│   └── go.sum
├── resources/                    # скриншоты, диаграммы
├── docs/                         # документация
├── .github/workflows/            # CI/CD
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

### Callback data

| Префикс | Описание |
|---------|----------|
| `menu:main` | главное меню |
| `plans:show` | выбор периода |
| `plans:today/tomorrow/week/nextweek` | показ планов |
| `plan:create` | начало создания |
| `plan:date:*` | выбор даты |
| `plan:time:*` | выбор времени |
| `plan:cancel` | отмена |

## API

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/plans` | Создать план |
| `GET` | `/plans?date=YYYY-MM-DD` | Планы на дату |
| `GET` | `/plans/range?start=&end=` | Планы за период |
| `DELETE` | `/plans/:id` | Удалить план |
| `GET` | `/health` | Проверка здоровья |

## Тесты

```bash
# Go тесты
make test-go

# C++ тесты
make test-cpp
make test-cpp-run

# Все тесты
make test
```

При запуске бота через `make bot-run` или Docker автоматически проверяются все Go тесты.

## Команды Makefile

```bash
make help            # список команд
make docker-up       # запуск через Docker
make docker-down     # остановка
make docker-logs     # логи
make bot-run         # запуск бота
make backend-run     # запуск backend
make backend-build   # сборка backend
make test            # все тесты
make clean           # очистка
```

## Технологии

| Компонент | Технология |
|-----------|------------|
| Telegram бот | Go 1.21 + [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) |
| Backend API | C++17 + [Crow](https://github.com/CrowCpp/Crow) |
| База данных | SQLite (WAL mode) |
| Контейнеры | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Тесты | Go testing + Google Test |

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `BOT_TOKEN` | Токен Telegram бота | — (обязательно) |
| `BACKEND_URL` | URL backend API | `http://localhost:8081` |
| `PORT` | Порт бота | `8080` |

## Лицензия

[MIT](LICENSE) — Dima Kiselev
