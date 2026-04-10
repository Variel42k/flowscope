# Admin Operations / Руководство администратора

## Русский

### Базовые команды

```powershell
docker compose up --build -d
docker compose ps
docker compose logs -f --tail=200
docker compose down
```

### Компоненты

- `collector` — ingestion + normalization
- `api` — REST + auth
- `worker` — rollups
- `flowgen` — synthetic flow generator
- `seed` — demo dataset import
- `clickhouse` — storage
- `web` — frontend

### Минимальный production checklist

1. Сменить дефолтные креды и секреты.
2. Ограничить внешние порты ClickHouse и UDP listeners.
3. Включить TLS на ingress/reverse proxy.
4. Убедиться, что `FLOWSCOPE_AUTH_DISABLED=false`.

## English

### Basic commands

```powershell
docker compose up --build -d
docker compose ps
docker compose logs -f --tail=200
docker compose down
```

### Components

- `collector` — ingestion + normalization
- `api` — REST + auth
- `worker` — rollups
- `flowgen` — synthetic flow generator
- `seed` — demo dataset import
- `clickhouse` — storage
- `web` — frontend

### Minimal production checklist

1. Rotate default credentials and secrets.
2. Restrict external exposure of ClickHouse and UDP listeners.
3. Enable TLS at ingress/reverse proxy.
4. Ensure `FLOWSCOPE_AUTH_DISABLED=false`.
