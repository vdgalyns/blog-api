# Development Guide

Как работать с этим проектом.

## 🏗️ Архитектура

```
internal/
├── domain/          # Модели и интерфейсы
│   ├── post.go
│   ├── comment.go
│   └── error.go
├── usecase/         # Бизнес-логика
│   ├── post.go
│   └── comment.go
├── repository/      # Работа с БД
│   └── postgres/
│       ├── post.go
│       └── comment.go
└── delivery/http/   # HTTP обработчики
    ├── post.go
    ├── comment.go
    ├── response.go
    └── router.go
```

Слои независимы друг от друга через интерфейсы.

## 📝 Добавление нового endpoint'а

Пример: добавить удаление комментария

### 1. Добавить метод в usecase

`internal/usecase/comment.go`:
```go
func (c *CommentUseCase) Delete(ctx context.Context, id int) error {
    if id <= 0 {
        return domain.ErrInvalidID
    }
    return c.repository.Delete(ctx, id)
}
```

### 2. Добавить метод в repository

`internal/repository/postgres/comment.go`:
```go
func (r *CommentRepository) Delete(ctx context.Context, id int) error {
    _, err := r.db.Exec(ctx, 
        "DELETE FROM comments WHERE id = $1", id)
    return err
}
```

### 3. Добавить handler

`internal/delivery/http/comment.go`:
```go
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid id")
        return
    }
    
    err = h.useCase.Delete(r.Context(), id)
    if err != nil {
        respondError(w, http.StatusNotFound, "comment not found")
        return
    }
    
    respondJSON(w, http.StatusOK, map[string]string{
        "message": "deleted",
    })
}
```

### 4. Добавить маршрут

`internal/delivery/http/router.go`:
```go
router.Delete("/{id}", commentHandler.Delete)
```

### 5. Написать тесты

`tests/usecase/comment_test.go`:
```go
func TestCommentUseCase_Delete_Success(t *testing.T) {
    repo := &MockCommentRepository{
        deleteFunc: func(ctx context.Context, id int) error {
            return nil
        },
    }
    
    uc := usecase.NewCommentUseCase(repo)
    err := uc.Delete(context.Background(), 1)
    
    if err != nil {
        t.Errorf("Expected nil, got %v", err)
    }
}
```

## 🧪 Тестирование

### Написать тест

```go
func TestSomething(t *testing.T) {
    // Arrange: подготовка
    mock := &MockRepository{}
    uc := usecase.NewUseCase(mock)
    
    // Act: действие
    result, err := uc.SomeMethod(context.Background())
    
    // Assert: проверка
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if result == nil {
        t.Error("Expected result, got nil")
    }
}
```

### Запустить

```bash
# Все тесты
go test ./...

# Конкретный тест
go test -v -run TestSomething ./tests/usecase

# С покрытием
go test -cover ./...
```

## 🔍 Валидация

Добавлять в usecase методы вроде `validatePost()`:

```go
func (uc *PostUseCase) validatePost(p *domain.Post) error {
    if len(p.Title) < 3 {
        return domain.ErrTitleTooShort
    }
    if len(p.Title) > 255 {
        return domain.ErrTitleTooLong
    }
    if len(p.Content) < 1 {
        return domain.ErrContentTooShort
    }
    return nil
}
```

## 🗄️ Работа с БД

Используем pgx v5 для работы с PostgreSQL:

```go
// Запрос
rows, err := r.db.Query(ctx, "SELECT * FROM posts WHERE id = $1", id)
defer rows.Close()

// Сканирование
if rows.Next() {
    post := &domain.Post{}
    err := rows.Scan(&post.ID, &post.Title, ...)
}

// Выполнение
_, err := r.db.Exec(ctx, "INSERT INTO posts (...) VALUES (...)", ...)
```

## 📦 Структура пакета domain

```go
// Модель
type Post struct {
    ID        int
    Title     string
    Content   string
    CreatedAt time.Time
}

// Интерфейсы
type PostUseCase interface {
    Create(ctx context.Context, title, content string) (*Post, error)
    GetByID(ctx context.Context, id int) (*Post, error)
    Update(ctx context.Context, id int, title, content string) (*Post, error)
    Delete(ctx context.Context, id int) error
    GetAll(ctx context.Context) ([]*Post, error)
}

type PostRepository interface {
    Create(ctx context.Context, p *Post) (*Post, error)
    GetByID(ctx context.Context, id int) (*Post, error)
    // ...
}

// Ошибки
var (
    ErrNotFound = errors.New("not found")
    ErrInvalidID = errors.New("invalid id")
)
```

## 🔧 Конфигурация

В `.env`:
```env
POSTGRES_DSN=postgres://user:pass@localhost:5432/db?sslmode=disable
HTTP_PORT=8080
```

В коде:
```go
cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}
// cfg.Postgres.DSN
// cfg.HTTP.Port
```

## 💡 Советы

1. **Используй интерфейсы** - так легче тестировать и менять реализацию
2. **Валидация в usecase** - бизнес-правила должны быть там
3. **Mock объекты** - создавай mock'и для тестов с нужным поведением
4. **Логирование** - добавляй логи для отладки
5. **Обработка ошибок** - не игнорируй ошибки

## 📖 Полезные команды

```bash
make run                 # Локальный запуск
make test                # Все тесты
make docker-up          # Docker старт
make docker-logs        # Логи
make fmt                # Форматирование кода
```

---

Всё готово для разработки! 🚀
