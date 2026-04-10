# Architecture / Архитектура

## English

### Service topology
- `collector`: UDP ingestion, decode, normalize, enrich, raw insert
- `worker`: 1m/1h rollups + node dimension refresh
- `api`: query/filter endpoints for tables/charts/map
- `web`: React UI using REST APIs

### Data pipeline
1. UDP packet received.
2. Decoder maps source protocol into canonical flow schema.
3. Normalization derives subnet/private/time/hash fields.
4. Enrichment adds optional metadata.
5. Data stored in ClickHouse raw table.
6. Worker generates rollup and edge tables.
7. UI/API query rollups by default and raw flows for drill-down.

## Русский

### Топология сервисов
- `collector`: UDP ingest, декодирование, нормализация, enrichment, запись raw
- `worker`: rollup таблицы 1m/1h + обновление node dimension
- `api`: query/filter endpoints для таблиц/графиков/карты
- `web`: React UI, работающий через REST API

### Пайплайн данных
1. Принимается UDP пакет.
2. Декодер приводит запись к канонической схеме flow.
3. Нормализация вычисляет subnet/private/time/hash поля.
4. Enrichment добавляет метаданные.
5. Данные пишутся в raw таблицу ClickHouse.
6. Worker строит rollup и edge таблицы.
7. UI/API используют rollups по умолчанию и raw flows для drill-down.
