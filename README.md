# Book API (pet project)

Simple learning API on Go for books, authors, and categories.
Data is stored in memory (map/slice), so after restart everything is reset.

## Features

- CRUD for books
- list and create authors
- list and create categories
- book filters by `author_id` and `category_id`
- pagination for `GET /books`
- basic input validation

## Stack

- Go
- `gorilla/mux`
- `net/http`

## How to run

```powershell
cd C:\coding\godev\Book-API_Golang
go run .
```

Server starts at `http://localhost:8080`.

## Endpoints

### Books

- `GET /books`
- `POST /books`
- `GET /books/{id}`
- `PUT /books/{id}`
- `DELETE /books/{id}`

Query params for `GET /books`:

- `page` (default 1)
- `page_size` (default 10)
- `author_id`
- `category_id`

### Authors

- `GET /authors`
- `POST /authors`

### Categories

- `GET /categories`
- `POST /categories`

## Validation

- author/category `name` must not be empty
- book `title` is required
- book `price` cannot be less than 0
- `author_id` and `category_id` must exist

## What can be improved later

- add tests
- add sorting
- move storage into a separate layer
- connect a real DB (PostgreSQL/MySQL)
