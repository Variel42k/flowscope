# Быстрый старт (RU)

## Требования

- Docker Desktop / Docker Engine + Compose v2
- минимум 4 vCPU / 6 GB RAM для комфортной локальной демо-нагрузки
- свободные порты:
  - `5173` (Web)
  - `8088` (API)
  - `18123` (ClickHouse HTTP)
  - `19000` (ClickHouse native)
  - `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp` (flow listeners)

## Первый запуск

```powershell
cd c:\Users\bgrx4\OneDrive\Документы\git\flowscope
docker compose up --build -d
```

Проверка статуса:

```powershell
docker compose ps
```

Ожидаемо:
- `clickhouse` в статусе `healthy`
- `seed` завершился `Exited (0)`
- `api`, `collector`, `worker`, `flowgen`, `web` в статусе `Up`

## Вход в систему

- Web UI: `http://localhost:5173`
- Логин по умолчанию:
  - username: `admin`
  - password: `admin123`

## Быстрая проверка API

```powershell
Invoke-RestMethod -Uri http://localhost:8088/api/health
```

Ожидаемый ответ содержит `"status": "ok"`.

## Остановка

```powershell
docker compose down
```

Полная очистка (включая данные ClickHouse):

```powershell
docker compose down -v
```

## Повторный импорт демо-данных

```powershell
docker compose run --rm seed
```
