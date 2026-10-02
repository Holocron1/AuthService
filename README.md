# AuthService

Сервис аутентификации на Go. Проверяет логин и пароль, выдаёт JWT и продлевает его. Работает по двум протоколам одновременно: HTTP (BasicAuth / Bearer) и gRPC. Пользователи хранятся в PostgreSQL, пароли — в виде bcrypt-хешей.

> Эндпоинт регистрации пока не реализован: пользователей нужно добавлять в БД вручную (см. [Тестовый пользователь](#тестовый-пользователь)).

## Возможности

| Операция | HTTP | gRPC |
|----------|------|------|
| Login — проверка логина и пароля, выдача JWT | `/login`, BasicAuth. Токен в заголовке `Authorization: Bearer <jwt>` | `auth.AuthService/Login` |
| Verify — проверка токена и выдача нового с продлённым сроком | `/verify`, заголовок `Authorization: Bearer <jwt>`. Новый токен в заголовке ответа | `auth.AuthService/Verify` |

Токен подписывается алгоритмом HS256 и живёт `JWT_TTL` (в примерах — 60 минут). Содержимое токена не шифруется, подпись защищает его от подделки.

## Технологии

- Go 1.27
- PostgreSQL 16, драйвер [pgx/v5](https://github.com/jackc/pgx)
- [gRPC](https://grpc.io/) и Protocol Buffers
- [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) — генерация и проверка JWT
- `golang.org/x/crypto/bcrypt` — хеширование паролей
- [viper](https://github.com/spf13/viper) — чтение конфигурации из переменных окружения, [godotenv](https://github.com/joho/godotenv) — загрузка `.env` при локальном запуске
- [go.uber.org/mock](https://github.com/uber-go/mock) — моки для тестов
- Docker и Docker Compose

## Архитектура

```
HTTP-хендлеры (internal/handler) ──┐
                                   ├──► UserService (internal/service) ──► UserStore ──► PostgreSQL
gRPC-сервер  (internal/grpc)   ────┘
```

- Вся бизнес-логика (проверка пароля, выдача и обновление токена) находится в `internal/service`. HTTP и gRPC — тонкие слои: достают данные из запроса, вызывают сервис и переводят результат в ответ или код ошибки.
- Интерфейсы объявлены на стороне потребителя: хендлеры и gRPC-сервер описывают нужный им `UserService`, а сервис — нужный ему `UserStore`.
- Все зависимости собираются один раз в `cmd/auth/main.go` через конструкторы. HTTP и gRPC используют один и тот же экземпляр сервиса.
- gRPC-сервер запускается в отдельной горутине рядом с HTTP-сервером, у них разные порты.
- Контекст запроса передаётся от транспорта через сервис до запроса в БД.

## Структура проекта

```
cmd/auth/            точка входа (main.go)
configs/             загрузка конфигурации
internal/domain/     доменная модель (User)
internal/service/    бизнес-логика (UserServiceImpl)
internal/store/      хранилища пользователей (PostgresStore; InMemoryStore остался от ранних этапов)
internal/handler/    HTTP-хендлеры (/login, /verify), моки и тесты
internal/grpc/       gRPC-сервер (Login, Verify)
proto/               описание gRPC-контракта (auth.proto)
pkg/authpb/          код, сгенерированный из auth.proto
migrations/          SQL-миграции
Dockerfile           multi-stage сборка образа
docker-compose.yml   сервис и PostgreSQL
```

## Переменные окружения

| Переменная     | Описание                                    | Пример                                                       |
|----------------|---------------------------------------------|--------------------------------------------------------------|
| `HTTP_PORT`    | Порт HTTP-сервера                           | `8080`                                                       |
| `GRPC_PORT`    | Порт gRPC-сервера                           | `50051`                                                      |
| `DATABASE_URL` | Строка подключения к PostgreSQL             | `postgres://user:pass@localhost:5432/authdb?sslmode=disable` |
| `JWT_SECRET`   | Секретный ключ для подписи JWT              | `change-me`                                                  |
| `JWT_TTL`      | Время жизни токена (формат `time.Duration`) | `60m`                                                        |

Пример заполнен в `.env.example`. При локальном запуске значения читаются из файла `.env`, в Docker Compose они передаются через `environment`. Реальный `.env` в репозиторий не коммитится.

## Запуск через Docker Compose

```
docker compose up --build
```

Поднимутся два контейнера: `postgres` (с healthcheck) и `auth-service` (стартует только когда база готова принимать соединения). Порты: HTTP — `8080`, gRPC — `50051`.

Остановка с удалением данных БД:

```
docker compose down -v
```

После `down -v` база пустая, миграцию и тестового пользователя нужно создать заново.

### Миграция

Миграция лежит в `migrations/001_create_users.sql` и применяется вручную. Откройте psql внутри контейнера с базой:

```
docker compose exec postgres psql -U firstUser -d firstDatabase
```

и выполните содержимое файла миграции.

### Тестовый пользователь

В psql после создания таблицы:

```sql
INSERT INTO users (username, password)
VALUES ('testguy1', '$2b$10$fWRciNlAZsXefCmp5HijdetCWRkWVfNP2WUib8g9E/Ha.7n5vx4uG');
```

Это bcrypt-хеш пароля `secret password`.

## Локальный запуск

1. Скопировать `.env.example` в `.env` и заполнить значения (`DATABASE_URL` должен указывать на доступный PostgreSQL).
2. Применить миграцию и создать пользователя, как описано выше.
3. Запустить сервис:

```
go run cmd/auth/main.go
```

## Примеры запросов

### HTTP

```
curl -i -u testguy1:"secret password" http://localhost:8080/login
```

```
HTTP/1.1 200 OK
Authorization: Bearer <jwt>
```

```
curl -i -H "Authorization: Bearer <jwt>" http://localhost:8080/verify
```

```
HTTP/1.1 200 OK
Authorization: Bearer <новый_jwt>
```

В PowerShell вызывайте `curl.exe`: просто `curl` там — алиас на `Invoke-WebRequest`.

### gRPC

Нужен [grpcurl](https://github.com/fullstorydev/grpcurl). Сервер не публикует описание методов, поэтому контракт передаётся флагами `-import-path` и `-proto`.

PowerShell:

```
'{"username":"testguy1","password":"secret password"}' | grpcurl -plaintext -import-path proto -proto auth.proto -d '@' localhost:50051 auth.AuthService/Login
```

```
'{"jwtToken":"<jwt>"}' | grpcurl -plaintext -import-path proto -proto auth.proto -d '@' localhost:50051 auth.AuthService/Verify
```

bash:

```
grpcurl -plaintext -import-path proto -proto auth.proto -d '{"username":"testguy1","password":"secret password"}' localhost:50051 auth.AuthService/Login
```

Успешный ответ:

```json
{
  "jwtToken": "<jwt>"
}
```

## Коды ответов

### HTTP

| Код | Когда возвращается |
|-----|--------------------|
| 200 | Успешный login или verify |
| 401 | login: нет BasicAuth или неверный пароль. verify: нет заголовка, неверная схема, пустой, поддельный или просроченный токен |
| 500 | Внутренняя ошибка (например, недоступна БД) |

### gRPC

| Код | Когда возвращается |
|-----|--------------------|
| `OK` | Успешный Login или Verify |
| `Unauthenticated` | Неверный пароль или недействительный токен |
| `Internal` | Внутренняя ошибка |

**Известное ограничение:** если пользователя нет в базе, сервис сейчас отвечает 500 / `Internal`, а не 401 / `Unauthenticated`: хранилище не отличает «пользователь не найден» от сбоя БД.

## Тесты

```
go test ./...
```

Табличные юнит-тесты HTTP-хендлеров (`internal/handler`) используют мок `UserService` и `net/http/httptest`, обращения к БД нет. Тесты gRPC-слоя пока не написаны.

## Генерация кода

Go-код из `proto/auth.proto` генерируется утилитой `protoc`. Нужны `protoc` ([релизы protobuf](https://github.com/protocolbuffers/protobuf/releases)) и два плагина:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Генерация (из корня проекта, результат попадает в `pkg/authpb/`):

```
protoc --go_out=. --go_opt=module=github.com/Holocron1/authservice --go-grpc_out=. --go-grpc_opt=module=github.com/Holocron1/authservice proto/auth.proto
```

Моки генерируются через `mockgen` (`go install go.uber.org/mock/mockgen@latest`):

```
mockgen -source internal/handler/handler.go -destination internal/handler/mocks/user_service_mock.go -package mocks
mockgen -source internal/service/user.go -destination internal/service/mocks/user_mock.go -package mocks
```

Сгенерированные файлы лежат в репозитории. После изменения интерфейса или `.proto` их нужно перегенерировать.

## Планы

- Мониторинг и метрики (Prometheus, Grafana)
- Эндпоинт регистрации
- Тесты gRPC-слоя