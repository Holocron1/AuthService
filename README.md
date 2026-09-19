# AuthService

Сервис аутентификации на Go: регистрация пользователей (in-memory на ранних этапах, PostgreSQL — в финальной версии), логин через BasicAuth, выдача и обновление JWT-токенов.

## Возможности

- **POST /login** — принимает логин и пароль через BasicAuth, при успешной проверке возвращает `200 OK` и JWT (жизненный цикл 60 минут) в заголовке `Authorization: Bearer <token>`.
- **POST /verify** — принимает access-токен в заголовке `Authorization: Bearer <token>`, проверяет его валидность и, если токен действителен, возвращает `200 OK` с новым токеном и продлённым сроком действия.

## Технологии

- Go 1.27
- PostgreSQL 16 (хранение пользователей)
- [pgx/v5](https://github.com/jackc/pgx) — драйвер PostgreSQL
- [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) — генерация и проверка JWT
- golang.org/x/crypto/bcrypt — хеширование паролей
- [go.uber.org/mock](https://github.com/uber-go/mock) — моки для юнит-тестов
- Docker / Docker Compose

## Структура проекта

```
cmd/auth/            точка входа (main.go)
configs/             загрузка конфигурации из переменных окружения
internal/domain/     доменная модель (User)
internal/service/    бизнес-логика (UserService: валидация, генерация/обновление токена)
internal/store/      хранилища пользователей (InMemoryStore, PostgresStore)
internal/handler/    HTTP-хендлеры (/login, /verify) и их тесты
migrations/          SQL-миграции
Dockerfile           multi-stage сборка образа сервиса
docker-compose.yml   сервис + PostgreSQL
```

## Переменные окружения

| Переменная     | Описание                                    | Пример                                                            |
|----------------|----------------------------------------------|--------------------------------------------------------------------|
| `HTTP_PORT`    | Порт, на котором слушает сервис             | `8080`                                                              |
| `DATABASE_URL` | Строка подключения к PostgreSQL             | `postgres://user:pass@localhost:5432/authdb?sslmode=disable`       |
| `JWT_SECRET`   | Секретный ключ для подписи JWT              | `change-me`                                                         |
| `JWT_TTL`      | Время жизни токена (формат `time.Duration`) | `60m`                                                               |

Пример заполнен в `.env.example`. Реальный `.env` в репозиторий не коммитится.

## Локальный запуск

1. Скопировать `.env.example` в `.env` и заполнить значения.
2. Поднять локальный PostgreSQL и применить миграцию из `migrations/001_create_users.sql` (`CREATE TABLE users ...`).
3. Запустить сервис:
   ```
   go run cmd/auth/main.go
   ```

## Запуск через Docker Compose

```
docker compose up --build
```

Поднимет два контейнера: `postgres` (с healthcheck) и `auth-service` (стартует только после готовности базы). Сервис будет доступен на `http://localhost:8080`.

Полный сброс (включая данные БД):
```
docker compose down -v
```

## Миграции

SQL-миграция лежит в `migrations/001_create_users.sql`. Применяется вручную, например через `psql` внутри контейнера с базой:

```
docker exec -it <имя_контейнера_postgres> psql -U <user> -d <database>
```

и далее выполнить содержимое файла миграции.

## Запуск тестов

```
go test ./...
```

Юнит-тесты для `/login` и `/verify` используют сгенерированные моки `UserService` (`go.uber.org/mock`) и `net/http/httptest`, без обращения к реальной БД.

## Примеры запросов

### POST /login

```
curl -i -u <username>:<password> http://localhost:8080/login
```

Успешный ответ:
```
HTTP/1.1 200 OK
Authorization: Bearer <jwt>
```

### POST /verify

```
curl -i -H "Authorization: Bearer <jwt>" http://localhost:8080/verify
```

Успешный ответ:
```
HTTP/1.1 200 OK
Authorization: Bearer <новый_jwt>
```

## Ожидаемые HTTP-коды

| Код | Когда возвращается                                                        |
|-----|-----------------------------------------------------------------------------|
| 200 | Успешный login / verify                                                    |
| 401 | Нет BasicAuth, неверный пароль, несуществующий пользователь (login); нет заголовка Authorization, неверная схема, пустой/поддельный/просроченный токен (verify) |