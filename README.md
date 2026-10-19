# Сокращатель ссылок (Go)

[![hexlet-check](https://github.com/i1yas/go-from-scratch-project-278/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/i1yas/go-from-scratch-project-278/actions)
[![Linter and tests](https://github.com/i1yas/go-from-scratch-project-278/actions/workflows/ci.yml/badge.svg)](https://github.com/i1yas/go-from-scratch-project-278/actions/workflows/ci.yml)

Приложение формирует короткие ссылки и собирает статистику перехода по ним.

Демо: [go-from-scratch-project-278-uejj.onrender.com](https://go-from-scratch-project-278-uejj.onrender.com/)

## Стек

- Go - основная логика приложения
- Caddy - reverse-proxy сервер: раздача статики для веб-интерфейса и проксирование на API приложения
- Postgres - база данных: хранение ссылок и информации о посещениях

## Установка и запуск

### Склонируйте проект

```bash
git clone https://github.com/i1yas/go-from-scratch-project-278.git
cd go-from-scratch-project-278
```

### Соберите приложение с помощью docker

```bash
docker build -t shorturl:latest .
```

### Запустите базу данных и подготовьте переменные окружения

Запустите базу данных и выставите [переменные окружения](#переменные-окружения)

### Запустите приложение

На данном этапе ожидается, что Postgres запущен

```bash
docker run \
    -e ENV=production \
    -e BASE_URL=http://example.com \
    -e DATABASE_URL=postgres://postgres:postgres@localhost:5432/appdb \
    -p 127.0.0.1:80:80 \
    shorturl:latest
```

Приложение будет доступно по адресу `http://localhost`, откроется веб-интерфейс приложения.

## Использование

Работать с приложением можно через веб-интерфейс либо через HTTP API.

Примеры работы с API, замените `host` на свой адрес:

```bash
# Создание короткой ссылки
curl -X POST host/api/links \
    -d '{"original_url":"http://mywebsite.com","short_name":"mywebsite"}'
# ответ: {"id":1, ..., "short_url": ".../r/mywebsite"}

# Информация о ссылке по id
curl host/api/links/1
# ответ: {"id":1, ..., "short_url": ".../r/mywebsite"}

# Первые 10 ссылок
curl --get host/api/links \
    -d 'range=[0,9]'
# ответ: [{"id":1, ..., "short_url": ".../r/mywebsite"}, ...]

# 3 последние ссылки
curl --get host/api/links \
    -d 'range=[0,3]' \
    -d 'sort=["id","DESC"]'
```

## Разработка

Для запуска сервера потребуется база данных Postgres. В корне проекта приложен `docker-compose.yml` с БД. Описание `.env.dev` в разделе ["переменные окружения"](#переменные-окружения)

### Запуск сервера отдельно

```
source ./bin/read-dotenv.sh .env.dev  # Выставить переменные окружения
make db-up                            # Запустить базу данных в compose
make db-migrate                       # Прогнать миграции
make run                              # Запустить приложение

# опционально
make db-seed                          # Заполнить БД данными
```

- Приложение запущено на порте `8080`

### Запуск сервера вместе с фронтендом

Для запуска фронтенда потребуется установленный `node`. Фронтенд поставляется в пакете `npm`.

```
make frontend-setup                   # Установить зависимости фронтенда
source ./bin/read-dotenv.sh .env.dev  # Выставить переменные окружения
make db-up                            # Запустить базу данных в compose
make db-migrate                       # Прогнать миграции
make dev                              # Запустить сервер и фронтенд

# опционально
make db-seed                          # Заполнить БД данными
```

- Фронтенд сервер запущен на порте `5173`, на нем доступен веб-интерфейс
- Приложение запущено на порте `8080`

### Запуск сервера в контейнере

Dockerfile содержит полноценную конфигурацию для продакшена. Его можно запустить вместе тестовой базой с помощью `docker-compose.yml`.

```
source ./bin/read-dotenv.sh .env.dev  # Выставить переменные окружения
make up                               # Запустить контейнер сервера и базу данных в compose

# опционально
make db-seed                          # Заполнить БД данными

# после изменений в коде приложения нужно делать ребилд
make rebuild
```

- Приложение запущено на порте `8080`: Caddy проксирует запросы к серверу, а статику веб-интерфейса отдает напрямую.

## Переменные окружения

Приложение использует переменные окружения для конфигурации. Если в корне проекта лежит файл `.env`, приложение вычитает переменные из файла.

| Название       | Обязательно | По умолчанию  | Описание                             |
| -------------- | ----------- | ------------- | ------------------------------------ |
| `ENV`          | Нет         | `development` | Определяет режим работы приложения   |
| `BASE_URL`     | Да          |               | URL для формирования короткой ссылки |
| `DATABASE_URL` | Да          |               | URL для подключения к Postgres       |
| `HTTP_ADDR`    | Нет         | `:8080`       | Адрес HTTP API в форме `HOST:PORT`   |
| `SENTRY_DSN`   | Нет         |               | DSN для подключения мониторинга      |

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
