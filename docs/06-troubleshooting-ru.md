# Troubleshooting / Устранение неполадок

## Русский

### Web UI не открывается

```powershell
docker compose ps
docker compose logs web --tail=200
```

### API возвращает 401

- Повторно выполните вход
- проверьте `FLOWSCOPE_AUTH_DISABLED`
- проверьте единый `FLOWSCOPE_JWT_SECRET`

### Нет данных

```powershell
docker compose logs seed --tail=200
docker compose logs flowgen --tail=200
docker compose logs worker --tail=200
```

### Полная перезагрузка

```powershell
docker compose down -v
docker compose up --build -d
```

## English

### Web UI is unavailable

```powershell
docker compose ps
docker compose logs web --tail=200
```

### API returns 401

- log in again
- verify `FLOWSCOPE_AUTH_DISABLED`
- ensure consistent `FLOWSCOPE_JWT_SECRET`

### No data in dashboards

```powershell
docker compose logs seed --tail=200
docker compose logs flowgen --tail=200
docker compose logs worker --tail=200
```

### Full restart

```powershell
docker compose down -v
docker compose up --build -d
```
