# Архитектура

## Обзор

Микросервисная архитектура:

- **Go Bot** — UI через inline-кнопки, state machine диалогов
- **C++ Backend** — REST API, CRUD, SQLite
- **SQLite** — персистентное хранение

```
User → Telegram → Go Bot → HTTP → C++ Backend → SQLite
                      ← JSON ←
```

## Структура

```
app/
├── cmd/bot/main.go           # точка входа
├── internal/
│   ├── bot/
│   │   ├── bot.go            # Bot, New(), Run()
│   │   ├── handlers.go       # роутинг callback'ов
│   │   ├── keyboards.go      # inline-кнопки
│   │   ├── state.go          # conversation state
│   │   ├── client.go         # HTTP клиент
│   │   ├── models.go         # типы
│   │   └── format.go         # форматирование
│   └── config/
│       └── config.go         # конфигурация
└── cpp-backend/
    ├── include/              # заголовки
    ├── src/                  # реализация
    └── tests/                # тесты
```

## Go Bot

State machine для пошаговых диалогов. Callback data с префиксным роутингом.

| Префикс | Описание |
|---------|----------|
| `menu:main` | главное меню |
| `plans:*` | просмотр планов |
| `plan:create` | создание |
| `plan:date:*` | выбор даты |
| `plan:time:*` | выбор времени |
| `plan:cancel` | отмена |

## C++ Backend

REST API на Crow, SQLite с WAL mode.

### Модули

- `models.h` — Plan, PlanRequest
- `database.h/cpp` — DAO
- `handlers.h/cpp` — HTTP обработчики

### Схема БД

```sql
CREATE TABLE plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    date TEXT NOT NULL,
    time TEXT,
    is_all_day INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_plans_date ON plans(date);
```

## Деплой

```bash
docker-compose up -d
```
