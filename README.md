# Taskflow

Taskflow — pet-проект на Go: HTTP API для управления задачами с использованием PostgreSQL, Redis и Docker.  
Проект реализует авторизацию по API Key и rate limiting для создания задач.

## Стек
- Go 
- PostgreSQL
- Redis
- Docker / Docker Compose

## Возможности
- CRUD для задач
- Хранение данных в PostgreSQL
- Rate limit на `POST /tasks` (5 запросов за 10 секунд)
- Авторизация через `X-API-Key`
- Миграции БД
- Graceful shutdown
- Логирование HTTP-запросов

## Структура проекта

- `cmd/api` — точка входа приложения  
- `internal/api` — HTTP handlers  
- `internal/service` — бизнес-логика  
- `internal/repository` — доступ к данным  
- `internal/middleware` — middleware (auth, rate limit)  
- `migrations` — SQL миграции

## Переменные окружения

Пример (`.env.example`):

```env
DATABASE_URL=postgres://postgres:postgres@db:5432/taskflow?sslmode=disable
REDIS_ADDR=redis:6379
API_KEY=dev-secret-key
```
## Авторизация
Все запросы требуют заголовок:  
X-API-Key: dev-secret-key  
При отсутствии или неверном ключе сервер вернёт 401 Unauthorized.

## Примеры запросов

### Получить список задач
```bash
curl http://localhost:8080/tasks \
  -H "X-API-Key: dev-secret-key"
```

### Создать задачу
```bash
curl -X POST http://localhost:8080/tasks \
  -H "X-API-Key: dev-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"title":"my first task"}'
```

### Получить определенную задачу
```bash
curl -X GET http://localhost:8080/tasks/1 \
  -H "X-API-Key: dev-secret-key" 
```

### Обновить определенную задачу
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "X-API-Key: dev-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"title":"my first task_updated","completed":true}'
```

### Удалить определенную задачу
```bash
curl -X DELETE http://localhost:8080/tasks/1 \
  -H "X-API-Key: dev-secret-key" 
```

## Rate limit
Применяется только к POST /tasks  
Лимит: 5 запросов за 10 секунд  
При превышении возвращается:  
429 Too Many Requests  
