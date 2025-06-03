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

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=mail
SMTP_PASS=password

```


2. Собрать проект:

```
go mod tidy
go run main.go
```

---

## Документация

```
https://documenter.getpostman.com/view/27528778/2sB2qi8HvC#3e0f4e11-21fe-44c4-8bad-0841e89aafae
```

## WebSocket

Я установил расширение Simple WebSocket Client и указал в поле Server Location адрес:
ws://localhost:3000/api/ws?token={token}.
После изменения статуса статьи в message log отображается её ID.
Также, если статьи были опубликованы до того, как пользователь подключился,
в логах будут выведены и эти статьи.
