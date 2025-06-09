# VK Backend Internship: Dynamic Worker Pool

**Тема:** Примитивный динамический worker-pool на Go с возможностью добавления и удаления воркеров в рантайме.  
Входные данные — строки в канал, воркеры обрабатывают их (выводят номер воркера и саму строку).

## Структура проекта
```
├── README.md
├── Taskfile.yml
├── cmd
│   └── main.go                — точка входа CLI-приложения, демонстрирующего работу пула.
├── go.mod
├── go.sum
└── internal
    └── workerpool
        ├── interfaces.go
        ├── pool
        │   ├── pool.go        — реализация `Pool`.
        │   └── pool_test.go   — unit-тесты для основных сценариев.
        └── worker
            ├── worker.go      — реализация `Worker`.
            └── worker_test.go — unit-тесты для основных сценариев.
```

## Установка и запуск

1. **Склонируйте репозиторий**  
   ```bash
   git clone https://github.com/EgorikA4/go-dev-vk-backend-challenge.git
   cd go-dev-vk-backend-challenge
   ```

2. **Соберите двоичный файл**

   ```bash
   go mod tidy
   go build -o workerpool ./cmd
   ```

3. **Запустите CLI-демо**

   ```bash
   ./bin/workerpool-cli
   ```

   Программа автоматически:

   * Запускает 4 воркера.
   * Отправляет несколько заданий.
   * Динамически добавляет и удаляет воркеров.
   * Выводит в консоль, какой воркер и какую задачу обрабатывает.
   * Выполняет `Shutdown`.

---

## Тестирование

Все unit-тесты лежат в директориях реализации интерфейсов.

```bash
go test ./...
```

Процент покрытия тестами:

```bash
go test ./... -cover -coverprofile=coverage.out && go tool cover -html=coverage.out
```
