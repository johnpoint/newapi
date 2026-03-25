---
name: project-init
description: Guide for initializing a new Go API project using the newapi init command. Use this when asked to create a new project, scaffold a project structure, or set up a Go API project from scratch using newapi, gin, and go-bootstrap.
---

## Using `newapi init` to Initialize a Project

The `newapi init` command scaffolds a complete Go API project structure in the current directory.

### Command Format

```bash
newapi init [module-name]
```

- If `module-name` is provided, `go mod init <module-name>` is run automatically.
- If `go.mod` already exists, the module name is read from it and no argument is needed.

### Generated Project Layout

```
<project>/
├── cmd/
│   └── server/
│       └── main.go              ← server entry point
├── internal/
│   ├── endpoints/               ← generated endpoint files (newapi new/gen)
│   ├── services/                ← generated service files (newapi new/gen)
│   └── schemas/                 ← generated schema files (newapi new/gen)
└── go.mod
```

### Steps

1. Create and enter a new directory:
   ```bash
   mkdir myapp && cd myapp
   ```

2. Run `newapi init` with your module name:
   ```bash
   newapi init github.com/yourname/myapp
   ```

3. Generate your first API:
   ```bash
   newapi new -n Ping -m get -p /ping
   ```

4. Install dependencies and run:
   ```bash
   go mod tidy
   go run ./cmd/server
   ```

### What Gets Created

| Path | Description |
|------|-------------|
| `cmd/server/main.go` | Server entry point; imports `_ internal/endpoints` to register routes |
| `internal/endpoints/` | Directory for generated endpoint files |
| `internal/services/` | Directory for generated service files |
| `internal/schemas/` | Directory for generated schema files |
| `go.mod` | Created automatically if not present |

### Notes

- Running `newapi init` in an existing project (with `go.mod`) is safe — it only creates missing directories and skips existing files.
- The `main.go` entry point uses a side-effect import (`_ "module/internal/endpoints"`) so all endpoint `init()` functions run and register routes automatically.
- After initialization, use `newapi new` or `newapi gen` to add API endpoints.
