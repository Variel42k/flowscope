# Security and Hardening / Безопасность и hardening

## Security Baseline / Базовый профиль безопасности

- Rotate default admin credentials before public deployment.  
Смените дефолтные учетные данные администратора до публикации.
- Use strong JWT secret and keep it outside repository history.  
Используйте сильный JWT secret и не храните его в истории репозитория.
- Expose only required ports.  
Публикуйте только нужные порты.
- Keep ClickHouse private to internal network segment.  
Держите ClickHouse во внутреннем сегменте сети.
- Deploy behind TLS reverse proxy for web/api.  
Размещайте web/api за TLS reverse proxy.

## Auth and Access / Доступ и авторизация

- RBAC roles: `admin`, `viewer`
- write operations require `admin`
- optional OIDC SSO can replace local login flow

## Sensitive Data / Чувствительные данные

- Treat inventory files as sensitive metadata.
- Restrict file permissions for `.env` and runtime secrets.
- Avoid logging full authentication tokens and passwords.

## Defensive Scope / Defensive-рамки

FlowScope is designed for observability and defensive monitoring only.  
FlowScope предназначен только для наблюдаемости и defensive-мониторинга.

It does not include offensive functionality (scanning, exploitation, persistence, remote access tooling).  
Проект не включает offensive-функциональность (сканирование, эксплуатация, persistence, средства удаленного доступа).
