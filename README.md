# Blog Backend

REST API для моего блога на Go с PostgreSQL.

## 🎯 Что это?

Полнофункциональный REST API для блога:
- Управление постами (CRUD)
- Комментарии к постам
- PostgreSQL база
- Полное тестирование (36 тестов)

## 📋 Функциональность

### Реализовано

✅ **Посты (полный CRUD)**
- `GET /api/posts` - получить все посты
- `POST /api/posts` - создать новый пост
- `GET /api/posts/{id}` - получить пост по ID
- `PUT /api/posts/{id}` - обновить пост
- `DELETE /api/posts/{id}` - удалить пост

✅ **Комментарии**
- `POST /api/posts/{postID}/comments` - создать комментарий для поста
- `GET /api/posts/{postID}/comments` - получить все комментарии поста

✅ **Валидация**
- Проверка минимальной/максимальной длины полей
- Детальные сообщения об ошибках
- Обработка ошибок БД

✅ **Тестирование**
- Unit тесты для UseCase слоя
- Integration тесты для HTTP обработчиков
- Mock объекты для изоляции компонентов

✅ **Конфигурация**
- Поддержка переменных окружения
- Гибкая конфигурация БД и сервера

## 🚀 Быстрый старт

```bash
go mod download
docker-compose -f infra/docker-compose.dev.yaml up -d
```

Подробнее: [QUICKSTART.md](docs/QUICKSTART.md)

## 🧪 Тесты

```bash
go test ./...
go test -v ./tests/usecase ./tests/delivery/http
```
### Запуск локально

```bash
```

### Запуск с Docker Compose

```bash
docker-compose -f infra/docker-compose.dev.yaml up -d
```

## 📚 API примеры

### Создание поста

```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is my first post content"
  }'
```

**Ответ (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "My First Post",
    "content": "This is my first post content",
    "created_at": "2025-12-16T10:30:00Z"
  }
}
```

### Получение всех постов

```bash
curl http://localhost:8080/api/posts
```

**Ответ (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "My First Post",
      "content": "This is my first post content",
      "created_at": "2025-12-16T10:30:00Z"
    }
  ]
}
```

### Получение поста по ID

```bash
curl http://localhost:8080/api/posts/1
```

### Обновление поста

```bash
curl -X PUT http://localhost:8080/api/posts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content"
  }'
```

### Удаление поста

```bash
curl -X DELETE http://localhost:8080/api/posts/1
```

### Добавление комментария

```bash
curl -X POST http://localhost:8080/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -d '{"content": "Great post!"}'
```

### Получение комментариев поста

```bash
curl http://localhost:8080/api/posts/1/comments
```

## 📊 Структура проекта

```
backend/
├── cmd/
│   └── app/
│       └── main.go                 # Точка входа приложения
├── config/
│   └── config.go                   # Конфигурация приложения
├── internal/
│   ├── app/
│   │   └── app.go                  # Инициализация приложения
│   ├── delivery/http/
│   │   ├── router.go               # Маршруты
│   │   ├── post_handler.go         # Обработчик постов
│   │   ├── post_handler_test.go    # Тесты обработчика постов
│   │   ├── comment_handler.go      # Обработчик комментариев
│   │   ├── comment_handler_test.go # Тесты обработчика комментариев
│   │   ├── response.go             # Структуры ответов
│   │   └── ...                     # Другие обработчики
│   ├── domain/
│   │   ├── post.go                 # Модель и интерфейсы Post
│   │   ├── comment.go              # Модель и интерфейсы Comment
│   │   ├── error.go                # Пользовательские ошибки
│   │   └── ...                     # Другие модели
│   ├── repository/
│   │   └── postgres/
│   │       ├── post.go             # Repository Post
│   │       ├── comment.go          # Repository Comment
│   │       └── ...                 # Другие репозитории
│   └── usecase/
│       ├── post.go                 # Use case постов
│       ├── post_test.go            # Тесты use case постов
│       ├── comment.go              # Use case комментариев
│       ├── comment_test.go         # Тесты use case комментариев
│       └── ...                     # Другие use cases
├── pkg/
│   ├── httpserver/                 # HTTP сервер
│   ├── postgres/                   # Подключение к PostgreSQL
│   └── query/                      # Вспомогательные функции
├── sql/
│   ├── init.sql                    # Инициализация БД
│   └── blog_schema.sql             # Схема для Posts и Comments
├── infra/
│   ├── docker-compose.yaml         # Docker Compose конфигурация
│   └── angie/
│       └── angie.conf              # Конфигурация реверс-прокси
├── go.mod                          # Go модули
├── .env.example                    # Пример переменных окружения
└── README.md                       # Этот файл
```

## 🏗️ Архитектура

Проект использует **Clean Architecture** паттерн:

- **Delivery Layer** (HTTP обработчики) - взаимодействие с клиентом
- **UseCase Layer** - бизнес логика и валидация
- **Repository Layer** - работа с БД (PostgreSQL)
- **Domain Layer** - модели данных и интерфейсы

## ✔️ Валидация

### Посты
- **Title**: минимум 3 символа, максимум 255 символов
- **Content**: должен быть не пустым

### Комментарии
- **Content**: минимум 1 символ, максимум 1000 символов
- **PostID**: должен быть положительным числом и существующим постом

## 🔒 Обработка ошибок

Все ошибки возвращаются в единообразном формате:

```json
{
  "success": false,
  "error": "title must be at least 3 characters long"
}
```

## 📝 Типы ошибок

| Ошибка | Описание |
|--------|---------|
| `invalid input` | Невалидные входные данные |
| `title must be at least 3 characters long` | Заголовок слишком короткий |
| `title must not exceed 255 characters` | Заголовок слишком длинный |
| `content must be at least 1 character long` | Контент пустой |
| `comment must be at least 1 character long` | Комментарий пустой |
| `comment must not exceed 1000 characters` | Комментарий слишком длинный |
| `not found` | Ресурс не найден |

## 🔄 Workflow Git

1. Создайте ветку для функции: `git checkout -b feature/feature-name`
2. Делайте коммиты: `git commit -m "описание изменений"`
3. Отправляйте изменения: `git push origin feature/feature-name`
4. Создавайте Pull Request

## 📈 Планы на будущее

- [ ] Добавить аутентификацию и авторизацию
- [ ] Реализовать пагинацию
- [ ] Добавить фильтрацию и сортировку
- [ ] Реализовать кэширование
- [ ] Добавить логирование
- [ ] Интеграционные тесты с реальной БД
- [ ] Метрики и мониторинг
- [ ] API документация (Swagger)

## 🤝 Вклад

Приветствуются pull requests и issues. Для больших изменений сначала откройте issue для обсуждения.

## 📄 Лицензия

MIT License - смотрите файл LICENSE для деталей.

## 📧 Контакты

- Issues: [Создать issue](https://github.com/vdgalyns/blog-api/issues)

---

**Последнее обновление:** 16 декабря 2025 г.