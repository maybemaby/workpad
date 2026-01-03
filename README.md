# workpad

## Setup

```bash
task build-ci --concurrency 1
./
```

## Local Build

```bash
task build-ci
```

Manual
```bash
(cd frontend && pnpm install --frozen-lockfile)
(cd frontend && pnpm build)
go build -o workpad ./cmd/server/main.go
```

## Migrations

Up migrations are run on application start. To run migrations manually, use:

```bash
task migrate-up
task migrate-down
```

Adding new migrations:

```bash
task migrate-add
```
