# API Practice

## Как запустить локально

1. Изменить `.env`:

```
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=your_database
DB_PORT=3306

SECRET=your_secret

```


2. Собрать проект:

```
go mod tidy
go run main.go
```

---

## Эндпоинты

```
http://localhost:3000/api/swagger/index.html
```