# 🚀 Quick Start

Быстрый старт для моего блог-проекта на Go.

## ⚡ 5 минут до первого запроса

### 1. Подготовка (1 минута)

```bash
cd backend
go mod download
```

### 2. Запуск (1 минута)

Самый быстрый способ - Docker:

```bash
docker-compose -f infra/docker-compose.dev.yaml up -d
```

Готово! Сервер на `http://localhost:8080`, БД на `localhost:5432`

Остановить:
```bash
docker-compose -f infra/docker-compose.dev.yaml down
```

### 3. Локально (если нужно)

```bash
cp .env.example .env

# Только БД в Docker
docker-compose -f infra/docker-compose.dev.yaml up -d postgres

make run
```

### 4. Тестирование

```bash
# Создать пост
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Post","content":"Hello World!"}'

# Получить все посты
curl http://localhost:8080/api/posts

# Создать комментарий
curl -X POST http://localhost:8080/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -d '{"content":"Great post!"}'

# Получить комментарии поста
curl http://localhost:8080/api/posts/1/comments
```

---

## 📚 Команды

```bash
# Docker
docker-compose -f infra/docker-compose.dev.yaml up -d    # Старт
docker-compose -f infra/docker-compose.dev.yaml logs -f  # Логи
docker-compose -f infra/docker-compose.dev.yaml down     # Стоп

# Разработка
make run                 # Запустить
make test                # Тесты
make docker-up/down      # Docker управление
make help                # Все команды
```

## 📖 Ещё документация

- [README.md](../README.md) - описание проекта
- [API.md](API.md) - API документация
- [DEVELOPMENT.md](DEVELOPMENT.md) - разработка
- [COMPLETION_REPORT.md](COMPLETION_REPORT.md) - детали реализации

## 🐛 Проблемы

**Порт 8080 занят?**
```bash
lsof -i :8080  # Узнать кто занял
# или просто измени PORT в .env
```

**Ошибка подключения к БД?**
```bash
docker-compose -f infra/docker-compose.dev.yaml logs postgres
```

**Нужно пересоздать БД?**
```bash
docker-compose -f infra/docker-compose.dev.yaml down -v
docker-compose -f infra/docker-compose.dev.yaml up -d
```

---

## 📝 Примеры API

### Posts

```bash
# Создать
POST /api/posts
{"title":"Hello","content":"World"}

# Получить все
GET /api/posts

# Получить один
GET /api/posts/1

# Обновить
PUT /api/posts/1
{"title":"Updated","content":"New content"}

# Удалить
DELETE /api/posts/1
```

### Comments

```bash
# Создать
POST /api/posts/1/comments
{"content":"Great!"}

# Получить
GET /api/posts/1/comments
```

---

## 🆘 Помощь

- [API.md](API.md) - список всех endpoint'ов
- [DEVELOPMENT.md](DEVELOPMENT.md) - как писать код

---

**Запусти и используй! 🚀**
