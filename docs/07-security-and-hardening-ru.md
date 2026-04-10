# Security and Hardening / Безопасность и hardening

## Русский

Минимум перед production:
1. Сменить demo-креды и JWT secret.
2. Ограничить доступ к портам ClickHouse.
3. Ограничить UDP listeners по сети/ACL.
4. Включить TLS на ingress.
5. Не использовать `FLOWSCOPE_AUTH_DISABLED` вне dev.

## English

Minimum before production:
1. Rotate demo credentials and JWT secret.
2. Restrict ClickHouse port exposure.
3. Limit UDP listeners by ACL/firewall.
4. Enable TLS at ingress.
5. Do not use `FLOWSCOPE_AUTH_DISABLED` outside dev.
