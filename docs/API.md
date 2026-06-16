# API Documentation

## Base URL

```
http://localhost:8081
```

## Endpoints

### POST /plans

Создать новый план.

**Request:**
```json
{
  "title": "Встреча",
  "description": "Обсудить проект",
  "date": "2024-12-25",
  "time": "14:00",
  "is_all_day": false
}
```

**Response (201):**
```json
{
  "id": 1,
  "message": "Plan created"
}
```

**Errors:**
- `400` — невалидный JSON или отсутствуют обязательные поля
- `500` — ошибка сервера

---

### GET /plans

Получить планы. Без параметров — все планы.

**Query Parameters:**
- `date` (optional) — фильтр по дате (YYYY-MM-DD)

**Response (200):**
```json
[
  {
    "id": 1,
    "title": "Встреча",
    "description": "Обсудить проект",
    "date": "2024-12-25",
    "time": "14:00",
    "is_all_day": false,
    "created_at": "2024-12-25 10:30:00"
  }
]
```

---

### GET /plans/range

Получить планы за период.

**Query Parameters:**
- `start` (required) — начальная дата (YYYY-MM-DD)
- `end` (required) — конечная дата (YYYY-MM-DD)

**Response (200):**
```json
[
  {
    "id": 1,
    "title": "Встреча",
    "description": "Обсудить проект",
    "date": "2024-12-25",
    "time": "14:00",
    "is_all_day": false,
    "created_at": "2024-12-25 10:30:00"
  }
]
```

**Errors:**
- `400` — отсутствуют параметры start/end

---

### DELETE /plans/:id

Удалить план по ID.

**Response (200):**
```json
{
  "message": "Plan deleted"
}
```

**Errors:**
- `404` — план не найден

---

### GET /health

Проверка здоровья сервиса.

**Response (200):**
```json
{
  "status": "ok",
  "service": "todo-backend"
}
```

## Models

### Plan
```json
{
  "id": 1,
  "title": "string",
  "description": "string",
  "date": "YYYY-MM-DD",
  "time": "HH:MM",
  "is_all_day": true,
  "created_at": "datetime"
}
```

### PlanRequest
```json
{
  "title": "string",
  "description": "string",
  "date": "YYYY-MM-DD",
  "time": "HH:MM",
  "is_all_day": true
}
```
