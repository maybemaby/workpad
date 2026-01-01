# workpad

## Overview

This project is a general purpose API template.

- Tech Stack: Go, OpenAPI, Postgres

## Development

- Use `task` to run commands, see `Taskfile.yml` for available tasks.
- See the `/migrations` folder for the database schema. Use tasks `migrate-up`, `migrate-down`, and `migrate-add` to manage migrations. Default to postgres.
- OpenAPI spec is created declaratively when mounting routes in `api/router.go`. Uses github.com/oaswrap/spec.
- Use go-spec skill to understand how to define routes with OpenAPI spec.


## Code Structure

- Use standard net/http handlers.
- Define api handlers in the `/api` folder.
- Domain specific logic can go in individual folders in `/api`.
- Use modern stdlib packages like `slices` `strings` and `maps` for common operations.

## Code Style
- Mark all required fields in request structs with `required:"true"` tag for OpenAPI specification generation.

