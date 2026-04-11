# Quick Start / Быстрый старт

## Русский

### Требования

- Docker 24+
- Docker Compose v2
- свободные порты:
  - `5173` (web)
  - `8088` (api)
  - `18123`, `19000` (clickhouse)
  - `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp` (flow listeners)

### Запуск

```bash
docker compose up --build -d
```

### Проверка

```bash
docker compose ps
curl http://localhost:8088/api/health
```

Ожидаемый ответ:

```json
{"status":"ok"}
```

### Доступ

- `http://localhost:5173`
- User: `admin`
- Password: `admin123`

### Остановка

```bash
docker compose down
```

Удаление volumes:

```bash
docker compose down -v
```

### Что дальше

1. Откройте `Overview` и проверьте топовые сущности.
2. Откройте `Interaction Map` и выберите ребро с большим объемом.
3. Провалитесь в raw flows для проверки.

## English

### Prerequisites

- Docker 24+
- Docker Compose v2
- free ports:
  - `5173` (web)
  - `8088` (api)
  - `18123`, `19000` (clickhouse)
  - `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp` (flow listeners)

### Run

```bash
docker compose up --build -d
```

### Verify

```bash
docker compose ps
curl http://localhost:8088/api/health
```

Expected response:

```json
{"status":"ok"}
```

### Open UI

- `http://localhost:5173`
- User: `admin`
- Password: `admin123`

### Stop

```bash
docker compose down
```

Delete data volumes:

```bash
docker compose down -v
```

### Next Steps

1. Open `Overview` and check top entities.
2. Open `Interaction Map` and select a high-volume edge.
3. Drill down to raw flows for validation.
