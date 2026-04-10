# Quick Start / Быстрый старт

## Русский

### Требования

- Docker Desktop / Docker Engine + Compose v2
- минимум 4 vCPU / 6 GB RAM для комфортной демо-нагрузки
- свободные порты: `5173`, `8088`, `18123`, `19000`, `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp`

### Первый запуск

```powershell
cd c:\Users\bgrx4\OneDrive\Документы\git\flowscope
docker compose up --build -d
```

### Проверка

```powershell
docker compose ps
Invoke-RestMethod -Uri http://localhost:8088/api/health
```

Ожидаемо:
- `clickhouse` -> `healthy`
- `seed` -> `Exited (0)`
- `api`, `collector`, `worker`, `flowgen`, `web` -> `Up`

### Вход

- Web UI: `http://localhost:5173`
- username: `admin`
- password: `admin123`

### Остановка

```powershell
docker compose down
# полная очистка:
docker compose down -v
```

### Повторный импорт seed

```powershell
docker compose run --rm seed
```

## English

### Requirements

- Docker Desktop / Docker Engine + Compose v2
- at least 4 vCPU / 6 GB RAM for local demo load
- free ports: `5173`, `8088`, `18123`, `19000`, `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp`

### First run

```powershell
cd c:\Users\bgrx4\OneDrive\Документы\git\flowscope
docker compose up --build -d
```

### Verification

```powershell
docker compose ps
Invoke-RestMethod -Uri http://localhost:8088/api/health
```

Expected:
- `clickhouse` -> `healthy`
- `seed` -> `Exited (0)`
- `api`, `collector`, `worker`, `flowgen`, `web` -> `Up`

### Login

- Web UI: `http://localhost:5173`
- username: `admin`
- password: `admin123`

### Stop

```powershell
docker compose down
# full cleanup:
docker compose down -v
```

### Re-run seed import

```powershell
docker compose run --rm seed
```
