# Troubleshooting (RU)

## Симптом: Web UI не открывается

Проверьте:

```powershell
docker compose ps
docker compose logs web --tail=200
```

Ожидаемо `web` в статусе `healthy`.

## Симптом: API `401 unauthorized`

1. Выполните логин в UI.
2. Проверьте, что `FLOWSCOPE_AUTH_DISABLED=false` (или корректно настроен auth).
3. Убедитесь, что `FLOWSCOPE_JWT_SECRET` одинаковый для API сервиса.

## Симптом: Нет данных в дашбордах

1. Проверьте `seed` и `flowgen`:

```powershell
docker compose logs seed --tail=200
docker compose logs flowgen --tail=200
```

2. Проверьте worker rollups:

```powershell
docker compose logs worker --tail=200
```

## Симптом: Collector не принимает exporter трафик

1. Проверьте проброс UDP портов.
2. Убедитесь, что exporter отправляет в правильный host/port.
3. Посмотрите `collector` logs на decode errors.

## Симптом: Пустая карта Interaction Map

1. Расширьте time-range (`From/To`).
2. Уменьшите `Min Bytes` и `Min Flows`.
3. Выключите `group_by` (поставьте `none`) и проверьте снова.

## Полная перезагрузка локального стенда

```powershell
docker compose down -v
docker compose up --build -d
```
