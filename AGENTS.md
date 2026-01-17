# Issue Flow - AI Agent Context

## Project Overview

Issue Flow is a Go CLI tool for managing GitHub issues, git worktrees, and development workflows across multiple projects.

## Architecture

**Language**: Go 1.25.6
**CLI Framework**: github.com/spf13/cobra
**Database**: SQLite (github.com/mattn/go-sqlite3)
**Testing**: github.com/stretchr/testify

### Directory Structure

```
issue-flow/
├── cmd/                    # CLI commands
├── internal/               # Private packages
│   ├── storage/            # Database layer
│   ├── project/           # Project business logic
│   └── config/            # Configuration
├── testutil/              # Test utilities (DO NOT import in production)
└── docs/                  # Design and reference docs
```

## Key Patterns

### 1. CLI Commands (cmd/*.go)

- Use Cobra framework
- Commands call `getDB()` for database access (injected for testing)
- Always check `shouldCloseDB()` before `defer db.Close()`
- Write errors to stderr, exit with code 1 on failures
- Use `fmt.Printf()` for success output

```go
func myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    Run: func(cmd *cobra.Command, args []string) {
        db, err := getDB()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
        if shouldCloseDB(db) {
            defer db.Close()
        }
        // Implementation
    },
}
```

### 2. Database Layer (internal/storage/db.go)

- Use parameterized SQL queries
- Always check errors on `Scan()`
- Always `defer rows.Close()` when using `Query()`
- Wrap errors with context: `fmt.Errorf("failed to ...: %w", err)`

### 3. Manager Layer (internal/*/manager.go)

- Business logic between commands and storage
- Convert internal types to storage types and back
- Validate before database operations
- Return wrapped errors with context

### 4. Testing (cmd/*_test.go)

- **CLI-level integration tests required for all commands**
- Use in-memory databases via `testutil.NewTestDB(t)`
- Inject test DB: `testDB = db`
- Cleanup: `t.Cleanup(func() { testDB = nil })`
- Verify both CLI output AND database state

```go
func TestMyCommand(t *testing.T) {
    db := testutil.NewTestDB(t)
    testDB = db
    t.Cleanup(func() { testDB = nil })

    rootCmd.SetArgs([]string{"command", "args"})
    err := rootCmd.Execute()
    require.NoError(t, err)

    testutil.AssertProjectCount(t, db, 1)
}
```

## Code Quality Requirements

- **Format**: Always run `make fmt` (or `gofmt -w .`)
- **Lint**: Always run `make vet` (or `go vet ./...`)
- **Test**: All tests must pass: `make test`
- **Pre-commit checklist**: Format → Vet → Test → Commit

## Design Adherence

All changes must follow documented patterns in `docs/`:
- **AGENT_GUIDE.md**: Best practices for AI agents contributing
- **ISSUE_FLOW_QUICKREF.md**: Command structure, file locations, patterns
- **ISSUE_FLOW_DESIGN.md**: Technical design and architecture
- **ISSUE_FLOW_IMPLEMENTATION_PLAN.md**: Implementation phases and roadmap

## Database Schema

- **projects**: ID, name, GitHub owner/repo, local path, worktree dir, config
- **worktrees**: ID, project ID, issue number, path, branch, status
- **issue_cache**: Project-scoped issue metadata (cached from GitHub)

## Storage Location

- **Database**: `~/.issue-flow/database.db` (SQLite)
- **Config**: `~/.issue-flow/config.yaml` (global), `<repo>/.issue-flow.yaml` (project)

## OpenCode Integration

When working on Issue Flow:
- Use in-memory databases (`:memory:`) for testing
- Database injection via `testDB` global variable in cmd/root.go
- CLI tests verify commands produce correct DB state
- See `testutil/testutil.go` for test helpers

## Common Commands

```bash
make build              # Build CLI
make test               # Run all tests
make fmt                # Format code
make vet                # Run static analysis
go run main.go version # Run without building
```
