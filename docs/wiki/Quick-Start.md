# Quick Start / Быстрый старт

## Prerequisites / Требования

- Docker 24+
- Docker Compose v2
- free ports / свободные порты:
  - `5173` (web)
  - `8088` (api)
  - `18123`, `19000` (clickhouse)
  - `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp` (flow listeners)

## Run / Запуск

```bash
docker compose up --build -d
```

## Verify / Проверка

```bash
docker compose ps
curl http://localhost:8088/api/health
```

Expected response / Ожидаемый ответ:

```json
{"status":"ok"}
```

## Open UI / Открыть интерфейс

- `http://localhost:5173`
- User: `admin`
- Password: `admin123`

## Stop / Остановка

```bash
docker compose down
```

Delete data volumes / удалить данные:

```bash
docker compose down -v
```

## Next Steps / Что дальше

1. Open **Overview** and check top entities.  
Откройте **Overview** и проверьте топовые сущности.
2. Open **Interaction Map** and select a high-volume edge.  
Откройте **Interaction Map** и выберите ребро с большим объемом.
3. Drill down to raw flows for validation.  
Провалитесь в сырые потоки для проверки.
