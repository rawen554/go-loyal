# go-loyal

## Конфигурирование сервиса накопительной системы лояльности:

- адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
  `example: :8080`
- Required!: адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
  `example: postgres://gophermart:P@ssw0rd@localhost:5432/gophermart?sslmode=disable`
- Required!: адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r.
  `example: :8081`

## Запуск через docker-copmpose

- `docker compose up -d`

## Makefile

Для чистой сборки приложения и запуска используйте `make clean-run`.
