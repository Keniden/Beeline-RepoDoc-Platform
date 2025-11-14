# Beeline RepoDoc Platform (Go)

Бэкенд уровня enterprise для аналитики репозиториев, извлечения AST, анализа истории Git, построения графов зависимостей и генерации документации через LLM (YandexGPT) на Go.

## Архитектура

- **многослойные сервисы**: ingestion, AST-парсер, git-анализатор, графостроитель, LLM-пайплайн, генератор документации, адаптеры хранилищ, Kafka-события.
- **Хранилища**: PostgreSQL (SQLC), ClickHouse, Redis-кеш, MinIO артефакты.
- **API**: Gin REST-эндпоинты и gRPC-сервисы, определённые в `proto/repodoc.proto`.
- **LLM**: клиент YandexGPT с батчингом, ретраями и выводом Markdown.
- **Граф**: DAG архитектуры экспортируется в Graphviz `.dot` и adjacency list для UI.
- **Обсервабилити**: OpenTelemetry-трейсы, метрики Prometheus, дашборды Grafana через compose.
- **Инфраструктура**: Docker + docker-compose поднимают API, worker, Postgres, ClickHouse, Kafka/Zookeeper, Redis, MinIO, Prometheus, Grafana.

## Быстрый старт

1. Скопируйте `configs/config.yaml` и переопределите секреты через переменные окружения.
2. Запустите миграции: `psql ... -f migrations/001_init.sql`.
3. Поднимите стек: `docker-compose up --build`.
4. REST API доступно на `localhost:8080`, gRPC — на `localhost:9090`.
5. Запускайте анализ: `POST /repos/upload` → `POST /repos/{id}/analyze`.

## Тесты и качество

- `make lint` (golangci-lint)
- `make test` (тесты через go test, Testify можно добавлять по мере роста модулей)
- `make build` / `make docker`

## CI/CD

См. `.gitlab-ci.yml` — пайплайн проверяет код (golangci-lint), запускает юнит-тесты и собирает Docker-образы.

## Обсервабилити

- Конфиг Prometheus для скрейпа в `scripts/prometheus.yml`.
- Grafana подключается через docker-compose.
- OpenTelemetry трассирует каждый этап (AST → Git → Graph → LLM → Storage).

## Protobuf и gRPC

- Определения в `proto/repodoc.proto` и сгенерированная `proto/repodoc.pb.go`.
- Сервисы: `RepoManager` (upload/analyze/status/graph/docs) и `DocService` (LLM-генерация документации).

## Графы и документация

- `internal/graph` строит DAG и отдает `.dot`, JSON-списки смежности.
- Документация сохраняется в `docs/{repo_id}` и доступна через REST/gRPC.

## Дополнительно

- SQLC-конфигурация — в `sqlc/`, миграции — в `migrations/`.
- Kafka-события публикуются на этапах (`file_analyzed`, `module_completed`, `doc_generated`), воркер (`cmd/worker`) их потребляет для обновления агрегатов.
- LLM-пайплайн устойчив: ограничение параллельных вызовов, логирование, батчинг и запись Markdown.
