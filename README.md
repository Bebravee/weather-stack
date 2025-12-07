# Weather Stack

Набор сервисов для регистрации пользователей, выдачи погоды/новостей и рекомендаций по одежде. В репозитории есть три части: `backend` (основное API), `web` (тонкий HTTP‑прокси/агрегатор поверх backend) и `tg-bot`.

## Как запустить
- Docker: `docker compose up --build` (поднимет PostgreSQL, backend на `:80`, web-прокси на `:8080`, бота).
- Переменные окружения backend лежат в `backend/configs/.env`; web прокидывает `REMOTE_HOST` и ходит на backend.
- Swagger из репозитория доступен по `http://localhost:80/docs` (backend) и `http://localhost:8080/docs` (web).

## Базовые адреса API
- Backend: `http://localhost:80/api/v1/...`
- Web-прокси (те же ручки, но без префикса `/api/v1`): `http://localhost:8080/...`

## Ручки backend (`/api/v1`)

### Пользователи / профиль
- `POST /auth/register` — регистрация пользователя. Тело:
  ```json
  {
    "name": "Иван",
    "sex": "male",
    "age": 27,
    "city_n": "Moscow",
    "city_w": "Moscow",
    "drop_time": "08:30",
    "t_comfort": 20,
    "t_tol": 15,
    "t_puh": 5,
    "temp1": 0,
    "telegram_id": 123456789,
    "password": "optional"
  }
  ```
  Ответ `200 OK`:
  ```json
  { "id": "UUID" }
  ```
  Ошибки: `422` (некорректное тело), `500`.

- `POST /auth/login` — заглушка, всегда `200 OK` без тела.

- `GET /profile/by-id/{id}` — получение профиля по UUID. Ответ `200 OK`:
  ```json
  {
    "id": "UUID",
    "name": "Иван",
    "sex": "male",
    "age": 27,
    "city_n": "Moscow",
    "city_w": "Moscow",
    "drop_time": "08:30",
    "t_comfort": 20,
    "t_tol": 15,
    "t_puh": 5,
    "temp1": 0,
    "telegram_id": 123456789,
    "created_at": "2025-01-01T12:00:00Z"
  }
  ```
  Ошибки: `404` (нет пользователя), `422` (параметр), `500`.

- `GET /profile/by-telegram-id/{telegram_id}` — профиль по Telegram ID. Ответ как выше. Ошибки: `404`, `422`, `500`.

- `PATCH /profile` — заглушка, сейчас возвращает `200 OK` без изменений.

### Погода
- `GET /weather/{city}` и `GET /weather/city/{city}` — актуальная погода по городу (берётся из кеша в БД, при промахе запрашивается у OpenWeather и сохраняется). Ответ `200 OK`:
  ```json
  {
    "city": "Moscow",
    "temperature": 14.2,
    "feels": 12.8,
    "description": "пасмурно",
    "humidity": 80,
    "pressure": 1008,
    "wind_speed": 4.2,
    "created_at": "2025-01-01T12:00:00Z"
  }
  ```
  Ошибки: `404` (нет записи и не удалось получить снаружи), `422`, `500`.

- `GET /weather/by-telegram-id/{telegram_id}` — погода по городу пользователя (ищет пользователя по Telegram ID, берёт `city_w`). Ответ как выше. Ошибки: `404` (нет пользователя или погоды), `422`, `500`.

- `GET /weather/clothes/{telegram_id}` — рекомендации по одежде для пользователя. При включённом LLM берутся подсказки из модели, иначе — stub. Пример ответа:
  ```json
  {
    "stub": true,
    "code": "01020304",
    "message": "Текст рекомендации",
    "user_city": "Moscow",
    "user_temps": { "comf": 20, "tol": 15, "puh": 5 },
    "weather_used": {
      "city": "Moscow",
      "temperature": 14.2,
      "feels": 12.8,
      "description": "пасмурно",
      "humidity": 80,
      "pressure": 1008,
      "wind_speed": 4.2,
      "created_at": "2025-01-01T12:00:00Z"
    },
    "created_at": "2025-01-01T12:00:00Z"
  }
  ```
  Ошибки: `404` (нет пользователя), `500`.

- `POST /weather` — сохранение записи о погоде (без обращения к внешнему API). Тело:
  ```json
  {
    "city": "Moscow",
    "temperature": 14.2,
    "feels_like": 12.8,
    "description": "пасмурно",
    "humidity": 80,
    "pressure": 1008,
    "wind_speed": 4.2
  }
  ```
  Ответ `200 OK`: `{ "city": "Moscow" }`. Ошибки: `422`, `500`.

### Новости
- `GET /news/{city_id}` и `GET /news/city/{city_id}` — новости по городу (кеш в БД + сохранение JSON на диск; при промахе запрос в NewsAPI). Ответ — JSON, возвращаемый NewsAPI, например:
  ```json
  {
    "status": "ok",
    "totalResults": 20,
    "articles": [ { "...": "..." } ]
  }
  ```
  Ошибки: `404` (нет новостей), `422`, `500` или `500` с телом `{ "status": "error", ... }` при ошибке внешнего API.

- `GET /news/by-telegram-id/{telegram_id}` — новости по городу пользователя (`city_n`). Ответ как выше. Ошибки: `404` (нет пользователя или новостей), `422`, `500`.

## Ручки web-прокси
Web-сервис повторяет те же маршруты, но без префикса `/api/v1` и проксирует запросы на backend, сохраняя схемы запросов/ответов:
- `/profile/register`, `/profile/{id}`, `/profile/by-telegram-id/{telegram_id}`
- `/weather/city/{city}`, `/weather/by-telegram-id/{telegram_id}`, `/weather/clothes/{telegram_id}`
- `/news/city/{city}`, `/news/by-telegram-id/{telegram_id}`

Если прокси запущен, им можно пользоваться вместо прямых обращений к backend.
