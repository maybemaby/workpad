# workpad

## Setup

```bash
task build-ci --concurrency 1
./
```

### Migrations

```bash
DB_URL=postgres://postgres:postgres@localhost:5432/workpadpg make migration-up
DB_URL=postgres://postgres:postgres@localhost:5432/workpadpg make migration-down
```
