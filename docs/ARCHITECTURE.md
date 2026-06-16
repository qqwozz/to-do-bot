# Архитектура проекта

## Обзор

Проект построен по принципу микросервисной архитектуры с разделением ответственности:

- **Telegram Bot (Go)** — обработка пользовательских команд, клавиатуры, взаимодействие с Telegram API
- **Backend API (C++/Crow)** — CRUD операции, бизнес-логика, работа с базой данных
- **SQLite** — персистентное хранение данных

## Схема взаимодействия

```
User ──> Telegram ──> Go Bot ──HTTP──> C++ Backend ──> SQLite
                           <──JSON──<
```

## Микросервисы

### Go Bot (порт по умолчанию: Telegram webhook)

**Ответственность:**
- Обработка команд (/start, /today, etc.)
- Формирование inline клавиатур
- Валидация пользовательского ввода
- HTTP запросы к Backend API

**Пакеты:**
- `cmd/bot` — точка входа
- `internal/bot` — логика бота
- `internal/config` — конфигурация

### C++ Backend (порт 8081)

**Ответственность:**
- REST API для CRUD операций
- Работа с SQLite
- Валидация данных

**Модули:**
- `models.h` — структуры данных
- `database.h/cpp` — DAO слой
- `handlers.h/cpp` — HTTP обработчики

## База данных

Таблица `plans`:
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
```

Индексы:
- `idx_plans_date` — ускорение поиска по дате

## Деплой

### Docker Compose
```bash
docker-compose up -d
```

### Переменные окружения
- `BOT_TOKEN` — токен Telegram бота (обязательно)
- `BACKEND_URL` — URL backend (по умолчанию http://localhost:8081)
