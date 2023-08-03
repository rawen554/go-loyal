# go-loyal

## Запуск через docker-copmpose

- `docker compose up -d`

## Makefile

### Конфигурация сервиса

- адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
  `example: :8080`
- Required!: адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
  `example: postgres://gophermart:P@ssw0rd@localhost:5432/gophermart?sslmode=disable`
- Required!: адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r.
  `example: :8081`

Перед запуском необходимо убедиться:

- что база данных работает на `localhost:5432`. Запуск БД - `make pg`.
- запущен экземпляр Accrual на `localhost:5432`.

Запуск с конфигурацией по умолчанию - `make run`.

Сборка приложения - `make build`.
