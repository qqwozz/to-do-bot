# Contributing

## Требования

- Go 1.21+
- C++17 компилятор
- CMake 3.14+
- SQLite3
- Docker (опционально)

## Развитие

### Go Bot

```bash
cd internal/bot
go test -v ./...
```

### C++ Backend

```bash
cd cpp-backend
mkdir build && cd build
cmake .. && make
```

### Docker

```bash
make docker-up
make docker-logs
```

## Тесты

```bash
make test       # все тесты
make test-go    # Go тесты
make test-cpp-run  # C++ тесты
```

## Структура коммитов

```
feat: add new feature
fix: bug fix
docs: documentation
test: add tests
refactor: code refactor
```
