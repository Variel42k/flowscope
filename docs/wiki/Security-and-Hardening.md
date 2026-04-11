# Security and Hardening / Безопасность и hardening

## Русский

### Базовый профиль безопасности

- Смените дефолтные учетные данные администратора перед публичным запуском.
- Используйте сильный `FLOWSCOPE_JWT_SECRET` и не храните его в истории репозитория.
- Публикуйте только необходимые порты.
- Держите ClickHouse во внутреннем сегменте сети.
- Размещайте web/api за TLS reverse proxy.

### Доступ и авторизация

- RBAC roles: `admin`, `viewer`
- write-операции требуют роль `admin`
- optional OIDC SSO может заменить локальный login-flow

### Чувствительные данные

- Inventory-файлы рассматривайте как чувствительные метаданные.
- Ограничьте права доступа для `.env` и runtime secrets.
- Не логируйте целиком токены аутентификации и пароли.

### Defensive-рамки

FlowScope предназначен только для observability и defensive-мониторинга. Проект не включает offensive-функциональность (сканирование, эксплуатация, persistence, инструменты удаленного доступа).

## English

### Security Baseline

- Rotate default admin credentials before public deployment.
- Use a strong `FLOWSCOPE_JWT_SECRET` and keep it out of repository history.
- Expose only required ports.
- Keep ClickHouse private to an internal network segment.
- Place web/api behind a TLS reverse proxy.

### Auth and Access

- RBAC roles: `admin`, `viewer`
- write operations require the `admin` role
- optional OIDC SSO can replace local login flow

### Sensitive Data

- Treat inventory files as sensitive metadata.
- Restrict file permissions for `.env` and runtime secrets.
- Avoid logging full authentication tokens and passwords.

### Defensive Scope

FlowScope is designed for observability and defensive monitoring only. It does not include offensive functionality (scanning, exploitation, persistence, remote access tooling).
