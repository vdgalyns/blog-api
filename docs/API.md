# API Documentation

## Base URL
```
http://localhost:8080/api
```

## Health Check

### Check server status
```http
GET /health
```

**Response:**
```
200 OK
OK
```

---

## Posts API

### Get All Posts

```http
GET /posts
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "First Post",
      "content": "Content here",
      "created_at": "2025-12-16T10:00:00Z"
    },
    {
      "id": 2,
      "title": "Second Post",
      "content": "Another content",
      "created_at": "2025-12-16T11:00:00Z"
    }
  ]
}
```

---

### Get Post by ID

```http
GET /posts/{id}
```

**Path Parameters:**
- `id` (integer, required) - Post ID

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "First Post",
    "content": "Content here",
    "created_at": "2025-12-16T10:00:00Z"
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "error": "post not found"
}
```

---

### Create Post

```http
POST /posts
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "My New Post",
  "content": "This is the content of my new post"
}
```

**Validation Rules:**
- `title` - Required, min 3 chars, max 255 chars
- `content` - Required, min 1 char

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 3,
    "title": "My New Post",
    "content": "This is the content of my new post",
    "created_at": "2025-12-16T12:00:00Z"
  }
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": "title must be at least 3 characters long"
}
```

---

### Update Post

```http
PUT /posts/{id}
Content-Type: application/json
```

**Path Parameters:**
- `id` (integer, required) - Post ID

**Request Body:**
```json
{
  "title": "Updated Title",
  "content": "Updated content"
}
```

**Validation Rules:**
- Same as Create Post

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Updated Title",
    "content": "Updated content",
    "created_at": "2025-12-16T10:00:00Z"
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "error": "post not found"
}
```

---

### Delete Post

```http
DELETE /posts/{id}
```

**Path Parameters:**
- `id` (integer, required) - Post ID

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "post deleted"
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "error": "post not found"
}
```

---

## Comments API

### Get Comments for Post

```http
GET /posts/{postID}/comments
```

**Path Parameters:**
- `postID` (integer, required) - Post ID

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "post_id": 1,
      "content": "Great post!",
      "created_at": "2025-12-16T10:30:00Z"
    },
    {
      "id": 2,
      "post_id": 1,
      "content": "Thanks for sharing",
      "created_at": "2025-12-16T10:45:00Z"
    }
  ]
}
```

**Response (200 OK, empty list):**
```json
{
  "success": true,
  "data": []
}
```

---

### Create Comment

```http
POST /posts/{postID}/comments
Content-Type: application/json
```

**Path Parameters:**
- `postID` (integer, required) - Post ID

**Request Body:**
```json
{
  "content": "This is a great comment!"
}
```

**Validation Rules:**
- `content` - Required, min 1 char, max 1000 chars
- `postID` - Required, must be positive integer

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 3,
    "post_id": 1,
    "content": "This is a great comment!",
    "created_at": "2025-12-16T13:00:00Z"
  }
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": "comment must not exceed 1000 characters"
}
```

**Response (400 Bad Request, invalid post ID):**
```json
{
  "success": false,
  "error": "invalid post id"
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "success": false,
  "error": "error message"
}
```

### Common Error Messages

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `invalid request body` | JSON is malformed |
| 400 | `invalid post id` | Post ID is not a valid integer |
| 400 | `title must be at least 3 characters long` | Title is too short |
| 400 | `title must not exceed 255 characters` | Title is too long |
| 400 | `invalid input` | Required field is empty |
| 400 | `content must be at least 1 character long` | Content is empty |
| 400 | `comment must be at least 1 character long` | Comment is empty |
| 400 | `comment must not exceed 1000 characters` | Comment is too long |
| 404 | `post not found` | Post with given ID doesn't exist |
| 500 | `failed to create post` | Database error |
| 500 | `failed to get post` | Database error |
| 500 | `failed to update post` | Database error |
| 500 | `failed to delete post` | Database error |
| 500 | `failed to create comment` | Database error |
| 500 | `failed to get comments` | Database error |

---

## Examples with cURL

### Create a post
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post"
  }'
```

### Get all posts
```bash
curl http://localhost:8080/api/posts
```

### Get specific post
```bash
curl http://localhost:8080/api/posts/1
```

### Update post
```bash
curl -X PUT http://localhost:8080/api/posts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Post",
    "content": "Updated content here"
  }'
```

### Delete post
```bash
curl -X DELETE http://localhost:8080/api/posts/1
```

### Add comment to post
```bash
curl -X POST http://localhost:8080/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -d '{"content": "Great post!"}'
```

### Get comments for post
```bash
curl http://localhost:8080/api/posts/1/comments
```

---

## Response Format

All responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "data": { /* response data */ },
  "pagination": null
}
```

### Error Response
```json
{
  "success": false,
  "error": "error message"
}
```

### Pagination (for future use)
```json
{
  "success": true,
  "data": [],
  "pagination": {
    "limit": 10,
    "offset": 0
  }
}
```

---

## Rate Limiting

Currently not implemented. Future versions may include rate limiting.

## Authentication

Currently not implemented. Future versions may include JWT authentication.

## Versioning

API version: **v1** (implicitly, through `/api` prefix)

Future versions will use `/api/v2`, etc.

---

**Last updated:** 16 December 2025
