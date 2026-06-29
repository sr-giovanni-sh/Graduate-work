## README.md

---

# Scheduler — планировщик задач на Go
# Использую GitHub Container Registry, т.к. Docker Hub недоступен

Проект представляет собой веб‑сервис для управления задачами (планировщик) на языке Go. Реализованы API для CRUD‑операций над задачами, поддержка повторяющихся задач с гибкой логикой расчёта следующей даты, аутентификация через JWT, хранение данных в SQLite.

Проект подходит как пример для выпускной работы (Graduate work), демонстрирует работу с HTTP‑сервером, базой данных, middleware, валидацией и обработкой ошибок на Go.

---

## Возможности

- **Управление задачами**: создание, чтение, обновление, удаление.
- **Повторяющиеся задачи**: поддержка правил повторения:
  - `d <интервал>` — ежедневно с заданным интервалом (1–400 дней);
  - `y` — ежегодно;
  - `w <дни>` — по выбранным дням недели (1–7, через запятую);
  - `m <дни>[,<месяцы>]` — по дням месяца (включая последний и предпоследний день), опционально по месяцам.
- **Авторасчёт следующей даты**: при завершении повторяющейся задачи автоматически вычисляется и сохраняется следующая дата по правилу.
- **Аутентификация**: вход по паролю, выдача JWT‑токена, проверка токена через middleware. Если пароль не задан, авторизация отключается.
- **Хранение данных**: SQLite в одном файле (`scheduler.db`).
- **Статические файлы**: сервер отдаёт фронтенд из папки `web`.

# Все задания со звёздочкой были выполнены
---

## Установка и запуск

### Требования

- Go 1.21+
- SQLite (драйвер встроен)

### Сборка и запуск

```bash
git clone <url-репозитория>
cd Graduate-work
go build -o scheduler
./scheduler
```

Сервер запустится на порту `7540` по умолчанию.

### Переменные окружения

| Переменная | Описание | По умолчанию |
| --- | --- | --- |
| `TODO_PORT` | Порт сервера | `7540` |
| `TODO_DBFILE` | Путь к файлу БД | `scheduler.db` |
| `TODO_PASSWORD` | Пароль для входа (если не задан — авторизация отключена) | — |

Пример запуска с настройками:

```bash
export TODO_PORT=8080
export TODO_DBFILE=./data/scheduler.db
export TODO_PASSWORD=my_secure_password
./scheduler
```

---

## API

### Публичные эндпоинты

- `POST /api/signin` — вход, в теле JSON: `{"password": "…"}`. В ответе: `{"token": "…"}`.
- `GET /api/nextdate?repeat=d 7&now=20251010` — расчёт следующей даты по правилу повторения. Обязателен параметр `repeat`, опционально `date` и `now`.

### Защищённые эндпоинты (требуется JWT‑токен в cookie `token`)

- `GET /api/tasks?search=…` — список задач, поддерживает пустой поиск, поиск по дате (формат `DD.MM.YYYY`) и по тексту (в заголовке или комментарии). Лимит — 50 задач.
- `POST /api/task` — создать задачу. Тело: `{"title": "…", "date": "20251001", "comment": "…", "repeat": "d 7"}`.
- `GET /api/task?id=1` — получить задачу по ID.
- `PUT /api/task` — обновить задачу. Обязательно указать `id` и `title`, остальные поля опциональны.
- `DELETE /api/task` — удалить задачу. В запросе указать `?id=…`.
- `POST /api/task/done` — отметить задачу выполненной. В запросе `?id=…`. Если у задачи есть правило повторения, будет рассчитана и установлена следующая дата; если нет — задача удаляется.

---

## Структура проекта

- `pkg/db` — работа с базой данных: инициализация, CRUD, структура `Task`.
- `pkg/server` — запуск HTTP‑сервера, настройка роутера, отдача статики.
- `pkg/api` — регистрация маршрутов, обработка запросов, логика расчёта дат, middleware авторизации.
- `main.go` — точка входа: инициализирует БД, запускает сервер.

---

## Примеры использования

**Создание задачи с ежедневным повторением:**

```bash
curl -X POST http://localhost:7540/api/task \
  -H "Content-Type: application/json" \
  -d '{"title":"Ежедневная задача","date":"20251001","repeat":"d 1"}'
```

**Завершение задачи — пересчёт следующей даты:**

```bash
curl -X POST http://localhost:7540/api/task/done?id=1 \
  -b "token=<ваш_JWT_токен>"
```

**Поиск задач по тексту:**

```bash
curl "http://localhost:7540/api/tasks?search=задача" \
  -b "token=<ваш_JWT_токен>"
```

---

## Тестирование

Для запуска тестов:

```bash
go test ./...
```
---

## Особенности реализации

- Все SQL‑запросы используют параметризацию, чтобы избежать SQL‑инъекций.
- Логика расчёта следующей даты защищена от бесконечных циклов лимитом итераций (`maxIterations = 2000`).
- При отсутствии файла БД он создаётся автоматически, таблица `scheduler` инициализируется.
- Папка со статикой ищется сначала в `./web`, затем в `../web`.

---

## Лицензия

Проект распространяется под лицензией MIT (если не указано иное).

---

## Контакты и вклад

Вклад приветствуется: форкайте репозиторий, создавайте PR с тестами и документацией.

---

---

# Scheduler — Task Scheduler in Go

A Go‑based web service for task management (scheduler) with support for recurring tasks, JWT authentication, and SQLite storage. Suitable as a graduate work example, demonstrating HTTP server, database access, middleware, validation, and error handling in Go.

---

## Features

- **Task management**: create, read, update, delete tasks.
- **Recurring tasks**: supports multiple repeat rules:
  - `d <interval>` — daily with interval (1–400 days);
  - `y` — yearly;
  - `w <days>` — specific weekdays (1–7, comma‑separated);
  - `m <days>[,<months>]` — specific days of the month (including last and second‑to‑last day), optional specific months.
- **Next date calculation**: when a recurring task is marked as done, the next occurrence date is automatically calculated and stored.
- **Authentication**: password‑based sign‑in, JWT token issuance, JWT validation via middleware. If no password is set, auth is disabled.
- **Data storage**: SQLite single‑file database (`scheduler.db`).
- **Static files**: serves frontend from the `web` directory.

---

## Installation and Running

### Requirements

- Go 1.21+
- SQLite (driver included)

### Build and Run

```bash
git clone <repository-url>
cd Graduate-work
go build -o scheduler
./scheduler
```

The server starts on port `7540` by default.

### Environment Variables

| Variable | Description | Default |
| --- | --- | --- |
| `TODO_PORT` | Server port | `7540` |
| `TODO_DBFILE` | Database file path | `scheduler.db` |
| `TODO_PASSWORD` | Login password (if unset, auth is disabled) | — |

Example with custom settings:

```bash
export TODO_PORT=8080
export TODO_DBFILE=./data/scheduler.db
export TODO_PASSWORD=my_secure_password
./scheduler
```

---

## API

### Public Endpoints

- `POST /api/signin` — sign in, JSON body: `{"password": "…"}`. Response: `{"token": "…"}`.
- `GET /api/nextdate?repeat=d 7&now=20251010` — calculate next date based on repeat rule. `repeat` is required; `date` and `now` are optional.

### Protected Endpoints (requires JWT token in `token` cookie)

- `GET /api/tasks?search=…` — list tasks. Supports empty search, date search (`DD.MM.YYYY`), and text search (in title or comment). Limit: 50 tasks.
- `POST /api/task` — create a task. Body: `{"title": "…", "date": "20251001", "comment": "…", "repeat": "d 7"}`.
- `GET /api/task?id=1` — get task by ID.
- `PUT /api/task` — update task. `id` and `title` are required; other fields are optional.
- `DELETE /api/task` — delete task. Specify `?id=…` in the request.
- `POST /api/task/done` — mark task as done. Specify `?id=…`. If the task has a repeat rule, the next date is calculated and stored; otherwise, the task is deleted.

---

## Project Structure

- `pkg/db` — database operations: initialization, CRUD, `Task` struct.
- `pkg/server` — HTTP server startup, router setup, static file serving.
- `pkg/api` — route registration, request handling, date logic, auth middleware.
- `main.go` — entry point: initializes DB, starts server.

---

## Usage Examples

**Create a daily recurring task:**

```bash
curl -X POST http://localhost:7540/api/task \
  -H "Content-Type: application/json" \
  -d '{"title":"Daily task","date":"20251001","repeat":"d 1"}'
```

**Mark task as done — recalculate next date:**

```bash
curl -X POST http://localhost:7540/api/task/done?id=1 \
  -b "token=<your_jwt_token>"
```

**Search tasks by text:**

```bash
curl "http://localhost:7540/api/tasks?search=task" \
  -b "token=<your_jwt_token>"
```

---

## Testing

To run tests:

```bash
go test ./...
```
---

## Implementation Notes

- All SQL queries use parameterization to prevent SQL injection.
- Next date calculation logic is protected against infinite loops with an iteration limit (`maxIterations = 2000`).
- If the database file is missing, it’s created automatically, and the `scheduler` table is initialized.
- Static files directory is searched first in `./web`, then in `../web`.

---

## License

This project is licensed under the MIT license (unless otherwise specified).

---

## Contacts and Contributions

Contributions are welcome: fork the repo, submit PRs with tests and documentation.