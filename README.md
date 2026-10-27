# Сокращатель ссылок (Go)

[![hexlet-check](https://github.com/i1yas/go-from-scratch-project-278/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/i1yas/go-from-scratch-project-278/actions)
[![Linter and tests](https://github.com/i1yas/go-from-scratch-project-278/actions/workflows/ci.yml/badge.svg)](https://github.com/i1yas/go-from-scratch-project-278/actions/workflows/ci.yml)

Приложение формирует короткие ссылки и собирает статистику перехода по ним.

Демо: [go-from-scratch-project-278-uejj.onrender.com](https://go-from-scratch-project-278-uejj.onrender.com/)

## Стек

- Go - основная логика приложения
- Caddy - reverse-proxy сервер: раздача статики для веб-интерфейса и проксирование на API приложения
- Postgres - база данных: хранение ссылок и информации о посещениях
- Docker - контейнеризация приложения

## Установка и запуск

### Требования

- bash
- git
- GNU Make
- Docker CLI 18.09+ (наличие команды `compose`)

### Склонируйте проект

```bash
git clone https://github.com/i1yas/go-from-scratch-project-278.git
cd go-from-scratch-project-278
```

### Запустите приложение

```bash
source ./bin/read-dotenv.sh .env.dev  # Выставить переменные окружения
make up                               # Запустить контейнер сервера и базу данных в compose
```

Приложение с веб-интерфейсом и HTTP API будет доступно по адресу `localhost:8080`

Приложение запускается внутри контейнера, через `docker-compose.yml`. Используются [переменные окружения](#переменные-окружения) `ENV` и `BASE_URL`. Для подключения к базе данных используется внутренний адрес compose, а адрес HTTP сервера переназначен на `localhost:8080` для локальной разработки.

## Использование

Работать с приложением можно через веб-интерфейс либо через HTTP API.
Ниже примеры работы с HTTP API.

### Создание короткой ссылки

```bash
curl -i -X POST http://localhost:8080/api/links \
    -H 'Content-Type: application/json' \
    -d '{"original_url":"https://google.com","short_name":"google"}'
```

Ответ

```
HTTP/1.1 201 Created
Content-Type: application/json
```

```json
{
    "id": 1,
    "original_url": "https://google.com",
    "short_name": "google",
    "short_url": "http://localhost:8080/r/google"
}
```

Теперь можно открыть ссылку [http://localhost:8080/r/google](http://localhost:8080/r/google) в браузере, приложение перенаправит на `https://google.com`

### Изменение ссылки

Установите переменную окружения `ID`, например `1`:

```bash
export ID=1
curl -i -X PUT http://localhost:8080/api/links/$ID \
    -H 'Content-Type: application/json' \
    -d '{"original_url":"https://google.com/some-page","short_name":"google"}'
```

Ответ

```
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
    "id": 1,
    "original_url": "https://google.com/some-page",
    "short_name": "google",
    "short_url": "http://localhost:8080/r/google"
}
```

### Посмотреть информацию о ссылке по ID

```bash
export ID=1
curl -i http://localhost:8080/api/links/$ID
```

Ответ

```
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
    "id": 1,
    "original_url": "https://google.com/some-page",
    "short_name": "google",
    "short_url": "http://localhost:8080/r/google"
}
```

### Посмотреть список ссылок
Параметр `range` (в формате [Range](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Range)), выставляем `[0,99]` чтобы получить первые 100 ссылок:
```bash
curl -i --get http://localhost:8080/api/links -d 'range=[0,99]'
```

Ответ

```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Range: links 0-0/1
```

```json
[
    {
        "id": 1,
        "original_url": "https://google.com/some-page",
        "short_name": "google",
        "short_url": "http://localhost:8080/r/google"
    }
]
```

### Посмотреть список переходов по ссылкам
Параметр `range` (в формате [Range](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Range)), выставляем `[0,99]` чтобы получить первые 100 переходов:
```bash
curl -i --get http://localhost:8080/api/link_visits -d 'range=[0,99]'
```

Ответ

```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Range: link_visits 0-1/2
```

```json
[
    {
        "id": 1,
        "link_id": 1,
        "ip": "127.0.0.1",
        "referer": "",
        "status": 302,
        "user_agent": "curl/7.81.0",
        "created_at": "2026-10-11T17:22:22Z"
    },
    {
        "id": 2,
        "link_id": 1,
        "ip": "::1",
        "referer": "http://mysite",
        "status": 302,
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36",
        "created_at": "2026-10-25:13:10:07Z"
    }
]
```

### Удаление ссылки

```bash
export ID=1
curl -i -X DELETE http://localhost:8080/api/links/$ID
```

Ответ

```
HTTP/1.1 204 No Content
```

Переход по ссылке [http://localhost:8080/r/google](http://localhost:8080/r/google) вернет ответ `404`

## Разработка

Для запуска сервера потребуется база данных Postgres. В корне проекта приложен `docker-compose.yml` с БД. Описание `.env.dev` в разделе ["переменные окружения"](#переменные-окружения)

### Требования

- bash
- git
- GNU Make
- Docker CLI 18.09+ (наличие команды `compose`)
- Go 1.26+
- node 20+

### Запуск приложения вне контейнера

```bash
source ./bin/read-dotenv.sh .env.dev  # Выставить переменные окружения
make db-up                            # Запустить базу данных в compose
make db-migrate                       # Прогнать миграции
make run                              # Запустить приложение
```

Приложение запущено на порте `8080`

Опционально можно заполнить БД данными

```bash
make db-seed
```

### Запуск вместе с фронтендом

Для запуска фронтенда потребуется установленный `node`. Фронтенд поставляется в пакете `npm`.

```bash
make frontend-setup                   # Установить зависимости фронтенда
source ./bin/read-dotenv.sh .env.dev  # Выставить переменные окружения
make db-up                            # Запустить базу данных в compose
make db-migrate                       # Прогнать миграции
make dev                              # Запустить сервер и фронтенд
```

- Фронтенд сервер запущен на порте `5173`, на нем доступен веб-интерфейс
- Приложение запущено на порте `8080`

## Переменные окружения

Приложение использует переменные окружения для конфигурации. Если в корне проекта лежит файл `.env`, приложение вычитает переменные из файла.

| Название       | Обязательно | По умолчанию | Описание                             |
| -------------- | ----------- | ------------ | ------------------------------------ |
| `ENV`          | Нет         | `production` | Определяет режим работы приложения   |
| `BASE_URL`     | Да          |              | URL для формирования короткой ссылки |
| `DATABASE_URL` | Да          |              | URL для подключения к Postgres       |
| `HTTP_ADDR`    | Нет         | `:8080`      | Адрес HTTP API в форме `HOST:PORT`   |
| `SENTRY_DSN`   | Нет         |              | DSN для подключения мониторинга      |

Возможные значения `ENV`:

- `production` - боевой режим
- `development` - режим разработки

В файле `.env.dev` переменные для разработки. В этом файле `DATABASE_URL` совпадает с конфигурацией базы данных в docker-compose.yml.

## Описание Makefile команд

| Команда                                  | Что делает                                               |
| ---------------------------------------- | -------------------------------------------------------- |
| `make build`                             | Собирает приложение в исполняемый файл bin/urlshort      |
| `make run`                               | Запускает приложение                                     |
| `make clean`                             | Удаляет артефакты сборки                                 |
| `make lint`                              | Проверяет код линтером                                   |
| `make lint-fix`                          | Исправляет ошибки линтера                                |
| `make format`                            | Форматирует код                                          |
| `make test`                              | Запускает тесты                                          |
| `make coverage`                          | Показывает покрытие тестов в терминале                   |
| `make coverage-html`                     | Показывает покрытие тестов в виде HTML-отчёта            |
| `make frontend-setup`                    | Устанавливает зависимости фронтенда                      |
| `make dev`                               | Запускает сервер и фронтенд                              |
| `make dev-frontend`                      | Запускает фронтенд отдельно                              |
| `make up`                                | Запускает контейнер с приложением и базу данных          |
| `make rebuild`                           | Пересобирает контейнер с приложением                     |
| `make down`                              | Останавливает контейнер с приложением и базу данных      |
| `make db-up`                             | Запускает базу данных                                    |
| `make db-down`                           | Останавливает базу данных                                |
| `make db-remove`                         | Удаляет базу данных                                      |
| `make db-shell`                          | Подключается к базе данных с помощью psql                |
| `make db-migrate`                        | Прогоняет миграции                                       |
| `make db-rollback`                       | Откатывает миграцию                                      |
| `make db-seed`                           | Заполняет базу данными из db/seeds/development.sql       |
| `make migration-add NAME=migration-name` | Создает миграцию                                         |
| `make sqlc-generate`                     | Генерирует код для работы с базой данных из sql-запросов |
