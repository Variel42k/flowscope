# Security Policy

## Scope

FlowScope is a defensive observability project.
Security issues in ingestion, authentication, data access, and deployment defaults are in scope.

## Supported Versions

Only the latest `main` state is actively maintained in this MVP stage.

## Reporting a Vulnerability

Please do not open public issues for critical vulnerabilities before maintainers can triage.

Include:

1. Affected version/commit
2. Impact assessment
3. Reproduction steps (safe, non-destructive)
4. Suggested remediation (if available)

## Response Expectations

- Initial triage target: 3 business days
- Severity classification: critical/high/medium/low
- Fix ETA depends on impact and complexity

## Hardening Baseline for Deployments

1. Rotate all default credentials and secrets.
2. Disable `FLOWSCOPE_AUTH_DISABLED` outside local development.
3. Restrict network exposure of ClickHouse and UDP listeners.
4. Enable TLS termination at ingress/reverse proxy.
5. Use allowlists/firewalls for exporter source addresses.
6. Keep dependencies and base images patched.

## Security Testing Recommendations

- `go test ./cmd/... ./internal/...`
- frontend tests/lint
- dependency checks (`govulncheck`, `npm audit`)
- container image scanning in CI/CD
