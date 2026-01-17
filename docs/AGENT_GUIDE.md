# Agent Guide for Issue Flow Development

**Best practices and conventions for AI agents contributing to Issue Flow**

---

## Overview

Issue Flow is a Go CLI tool for managing GitHub issues, git worktrees, and development workflows across multiple projects. This guide ensures AI agents follow consistent patterns, maintain code quality, and adhere to the design documented in `docs/`.

**Key Technologies:**
- **Language**: Go 1.25.6
- **CLI Framework**: github.com/spf13/cobra
- **Database**: SQLite (github.com/mattn/go-sqlite3)
- **Testing**: github.com/stretchr/testify
- **Project Structure**: Standard Go project with `cmd/`, `internal/`, `testutil/`

---

## Quick Start Checklist

When working on Issue Flow, always:

1. ✅ **Read design docs** in `docs/` before implementing features
2. ✅ **Follow existing patterns** in `cmd/` and `internal/`
3. ✅ **Write CLI-level tests** with in-memory databases
4. ✅ **Format code**: `make fmt` or `gofmt -w .`
5. ✅ **Run linting**: `make vet` or `go vet ./...`
6. ✅ **Run tests**: `make test` or `go test -v ./...`
7. ✅ **Test on real CLI**: Build and verify functionality
8. ✅ **Commit with clear messages** following git log style

---

## Project Structure

```
issue-flow/
├── cmd/                          # CLI commands
│   ├── root.go                   # Root command, DB injection
│   ├── root_test.go              # Version command test
│   ├── project.go                # Project commands
│   ├── project_test.go           # Project CLI integration tests
│   └── [new commands].go        # Add new commands here
├── internal/                     # Private packages
│   ├── storage/
│   │   └── db.go               # Database layer (SQL + models)
│   ├── project/
│   │   ├── manager.go            # Project business logic
│   │   └── types.go             # Project types
│   ├── config/
│   │   └── config.go            # Configuration management
│   └── [future packages]/
├── testutil/                     # Test utilities (DO NOT import in production code)
│   └── testutil.go              # DB helpers, assertions, test data
├── pkg/                          # Public packages (future use)
├── docs/                         # Design and reference docs
│   ├── ISSUE_FLOW_DESIGN.md      # Technical design
│   ├── ISSUE_FLOW_QUICKREF.md    # Quick reference
│   ├── ISSUE_FLOW_IMPLEMENTATION_PLAN.md
│   └── AGENT_GUIDE.md           # This document
├── Makefile                      # Build and test commands
├── go.mod                        # Go module definition
└── main.go                       # Application entry point
```

---

## Code Conventions

### 1. Go Code Style

Follow standard Go conventions:
- Use `gofmt` for formatting (automated via `make fmt`)
- Exported functions: `PascalCase`, internal: `camelCase`
- Errors: Always wrap with context using `fmt.Errorf`: `fmt.Errorf("failed to create project: %w", err)`
- Comments: Write godoc for exported functions
- Packages: One type per file when possible

### 2. Command Structure (Cobra)

```go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var myFlag string

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "One-line description",
    Long:  "Longer description with more details.",
    Args:  cobra.ExactArgs(1), // or NoArgs, MinimumNArgs, etc.
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
        fmt.Printf("Result: %s\n", args[0])
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
    myCmd.Flags().StringVarP(&myFlag, "flag", "f", "default", "Flag description")
}
```

**Key Points:**
- Always use `getDB()` instead of `storage.New()` to support test DB injection
- Always check `shouldCloseDB()` before `defer db.Close()` to keep test DB alive for assertions
- Write errors to `os.Stderr`, exit with `os.Exit(1)` on failures
- Use `fmt.Printf()` for success output to `os.Stdout`

### 3. Database Layer

All database operations go through `internal/storage/db.go`:

```go
// Add new CRUD method following existing patterns
func (d *Database) GetSomething(id string) (*Something, error) {
    query := `SELECT id, field1, field2 FROM somethings WHERE id = ?`
    row := d.db.QueryRow(query, id)

    var s Something
    err := row.Scan(&s.ID, &s.Field1, &s.Field2)
    if err != nil {
        return nil, err
    }
    return &s, nil
}

func (d *Database) ListSomethings() ([]Something, error) {
    query := `SELECT id, field1, field2 FROM somethings ORDER BY field1`
    rows, err := d.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var somethings []Something
    for rows.Next() {
        var s Something
        if err := rows.Scan(&s.ID, &s.Field1, &s.Field2); err != nil {
            return nil, err
        }
        somethings = append(somethings, s)
    }
    return somethings, nil
}
```

**Important:**
- Use parameterized queries (not string interpolation)
- Always check errors on `Scan()`
- Always `defer rows.Close()` when using `Query()`
- Add new tables to `initSchema()` function

### 4. Manager Layer (Business Logic)

Place business logic in `internal/*` packages:

```go
package project

import (
    "fmt"

    "github.com/paolorechia/issue-flow/internal/storage"
)

type Manager struct {
    db *storage.Database
}

func NewManager(db *storage.Database) *Manager {
    return &Manager{db: db}
}

func (m *Manager) Add(p *Project) error {
    if err := p.Validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    sp := &storage.Project{
        ID:          p.ID,
        Name:        p.Name,
        // ... map fields
    }

    return m.db.CreateProject(sp)
}

func (m *Manager) Get(id string) (*Project, error) {
    sp, err := m.db.GetProject(id)
    if err != nil {
        return nil, err
    }

    return &Project{
        ID:          sp.ID,
        Name:        sp.Name,
        // ... map fields back
    }, nil
}
```

---

## Testing Requirements

### 1. CLI-Level Integration Tests (Required)

All CLI commands must have integration tests in `cmd/*_test.go`:

```go
package cmd

import (
    "bytes"
    "testing"

    "github.com/paolorechia/issue-flow/testutil"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMyCommand(t *testing.T) {
    // Setup: Create fresh in-memory DB
    db := testutil.NewTestDB(t)
    testutil.AssertDBEmpty(t, db)

    // Setup: Configure command
    myFlag = "test-value"
    buf := new(bytes.Buffer)
    rootCmd.SetOut(buf)
    rootCmd.SetErr(buf)
    rootCmd.SetArgs([]string{"mycommand", "arg1"})

    // Inject test DB
    testDB = db
    t.Cleanup(func() { testDB = nil })

    // Execute command
    err := rootCmd.Execute()
    require.NoError(t, err, "Command should succeed")

    // Verify DB state
    testutil.AssertProjectCount(t, db, 1)
    project := testutil.AssertProjectExists(t, db, "project-id")
    assert.Equal(t, "expected-name", project.Name)

    // Optionally verify output
    assert.Contains(t, buf.String(), "Expected output")
}
```

**Key Requirements:**
- Each test starts with empty in-memory database
- Inject test DB via `testDB = db` before `rootCmd.Execute()`
- Always cleanup with `t.Cleanup(func() { testDB = nil })`
- Verify BOTH CLI output AND database state
- Use `require.NoError()` for setup failures
- Use `assert.*` for verification failures

### 2. Using Test Utilities

Available helpers from `testutil/testutil.go`:

```go
// Database helpers
db := testutil.NewTestDB(t)                    // Create in-memory DB
testutil.AssertDBEmpty(t, db)                  // Verify empty DB
testutil.AssertProjectCount(t, db, 2)          // Count projects
project := testutil.AssertProjectExists(t, db, id) // Get and verify project
testutil.AssertProjectNotExists(t, db, id)      // Verify doesn't exist
testutil.AssertWorktreeCount(t, db, 1)         // Count worktrees
worktree := testutil.AssertWorktreeExists(t, db, id) // Get worktree

// Test data helpers
project := testutil.CreateTestProject(t, db)            // Create test project
worktree := testutil.CreateTestWorktree(t, db, projectID, 123) // Create worktree

// Table parsing (for CLI output)
table := testutil.ParseTableOutput(t, buf.String())
testutil.AssertTableRow(t, table, 0, []string{"ID", "Name", "Repo"})
```

### 3. Testing Database Injection Pattern

Commands that exit on error cannot test error paths:

```go
// Test that calls os.Exit(1)
func TestCommandError(t *testing.T) {
    t.Skip("Skipped: command calls os.Exit(1) which terminates test process")
}

// Alternative: Test the function logic directly in internal package tests
```

### 4. Storage Layer Tests

Add unit tests for storage layer methods:

```go
package storage

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestDatabase_CreateProject(t *testing.T) {
    db, err := NewWithDBPath(":memory:")
    require.NoError(t, err)
    defer db.Close()

    project := &Project{
        ID:   "test",
        Name: "Test Project",
        // ... other fields
    }

    err = db.CreateProject(project)
    require.NoError(t, err)

    retrieved, err := db.GetProject("test")
    require.NoError(t, err)
    assert.Equal(t, "Test Project", retrieved.Name)
}
```

---

## Development Workflow

### Before Writing Code

1. **Read design docs**: Check `docs/ISSUE_FLOW_DESIGN.md` and `docs/ISSUE_FLOW_QUICKREF.md`
2. **Understand existing patterns**: Look at similar commands in `cmd/`
3. **Plan database changes**: If adding new features, check if tables need changes
4. **Plan tests**: Consider test cases before implementing

### During Implementation

1. **Write tests first**: Create test file with failing tests
2. **Implement minimal code**: Make tests pass
3. **Format code**: Run `make fmt` after changes
4. **Run tests frequently**: `make test` after each feature
5. **Run vet**: `make vet` to catch issues

### Before Committing

```bash
# 1. Format code
make fmt

# 2. Run linting
make vet

# 3. Run tests
make test

# 4. Build and test manually
make build
./bin/issue-flow project list

# 5. Check git status
git status

# 6. Review diff
git diff

# 7. Commit with clear message
git add .
git commit -m "Add feature: description of change"
```

---

## Adding New Features

### Adding a New Command

1. Create `cmd/mycommand.go`
2. Implement command following structure above
3. Add command to root in `init()`
4. Create `cmd/mycommand_test.go` with integration tests
5. Update `docs/ISSUE_FLOW_QUICKREF.md` with command details

### Adding Database Operations

1. Add table schema to `internal/storage/db.go` in `initSchema()`
2. Add struct for model (or use existing)
3. Implement CRUD methods in `internal/storage/db.go`
4. Add tests for new storage methods
5. Add manager methods in `internal/*` package if needed

### Adding Configuration

1. Update `internal/config/config.go` with new config fields
2. Add validation logic if needed
3. Document in `docs/ISSUE_FLOW_QUICKREF.md`
4. Update example config in docs

### Adding Templates

1. Create template file in appropriate location
2. Document template variables
3. Update template documentation in `docs/`

---

## Common Patterns

### Error Handling

```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create worktree: %w", err)
}

// Validate before operations
if projectID == "" {
    return fmt.Errorf("project ID is required")
}

// Use custom errors for better messages
var ErrProjectNotFound = errors.New("project not found")
```

### Command Flag Patterns

```go
// Required flags
var projectID string
var projectName string

cmd.Flags().StringVarP(&projectID, "id", "i", "", "Project ID (required)")
cmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name (required)")

// Validate in Run()
if projectID == "" || projectName == "" {
    fmt.Fprintln(os.Stderr, "Error: --id and --name are required")
    os.Exit(1)
}

// Optional flags
var verbose bool
cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
```

### Output Formatting

```go
// Simple output
fmt.Printf("✓ Created project: %s\n", project.Name)

// Tabular output
w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
fmt.Fprintln(w, "ID\tNAME\tREPOSITORY")
for _, p := range projects {
    fmt.Fprintf(w, "%s\t%s\t%s\n", p.ID, p.Name, p.Repo)
}
w.Flush()
```

---

## Design Document Adherence

### Must Follow

1. **Architecture**: Layered architecture (cmd → manager → storage)
2. **Storage**: SQLite with proper schema and migrations
3. **CLI**: Cobra command structure with proper flags
4. **Testing**: CLI-level integration tests with in-memory DBs
5. **Error Handling**: Wrap errors with context
6. **Code Style**: Standard Go conventions with gofmt

### Design Document References

- **Overall Design**: `docs/ISSUE_FLOW_DESIGN.md`
- **Implementation Plan**: `docs/ISSUE_FLOW_IMPLEMENTATION_PLAN.md`
- **Quick Reference**: `docs/ISSUE_FLOW_QUICKREF.md`
- **Bootstrap Guide**: `docs/ISSUE_FLOW_BOOTSTRAP.md`

When implementing features:
1. Check if feature is documented
2. Follow documented patterns
3. Update docs if adding new patterns

---

## Quality Assurance

### Pre-Commit Checklist

- [ ] Code formatted with `make fmt`
- [ ] No warnings from `make vet`
- [ ] All tests pass with `make test`
- [ ] New features have integration tests
- [ ] Documentation updated if needed
- [ ] Commit message clear and descriptive

### Test Coverage

- Aim for >80% coverage on new code
- All CLI commands must have integration tests
- Critical business logic should have unit tests

### Code Review Criteria

- Follows project conventions
- Tests verify both success and error paths
- Error messages are helpful
- Code is readable and maintainable
- No unnecessary dependencies added

---

## Troubleshooting

### Common Issues

**Issue**: Tests fail with "database is closed"
- **Fix**: Check `shouldCloseDB()` usage - don't close test DB

**Issue**: Commands don't find test data
- **Fix**: Ensure `testDB = db` is set before `rootCmd.Execute()`

**Issue**: fmt/vet errors
- **Fix**: Run `make fmt` and `make vet` before committing

**Issue**: Tests pass but CLI doesn't work
- **Fix**: Build and test manually with `make build && ./bin/issue-flow`

### Debugging

```bash
# Run specific test with verbose output
go test -v ./cmd -run TestMyCommand

# Run all tests with coverage
make test-coverage

# Check database directly
sqlite3 ~/.issue-flow/database.db "SELECT * FROM projects;"

# Run with verbose logging (when implemented)
issue-flow --verbose command args
```

---

## Commands Reference

### Makefile Commands

```bash
make build              # Build for current platform
make build-all          # Build for all platforms
make install            # Install to $GOPATH/bin
make test               # Run all tests
make test-coverage      # Run tests with coverage report
make clean              # Remove build artifacts
make fmt               # Format code (check only)
make fmt-fix           # Format code (fix issues)
make vet                # Run go vet
make lint               # Run fmt + vet + test
make ci                 # Run fmt + vet + test (CI)
```

### Go Commands

```bash
go run main.go version       # Run without building
go test -v ./...           # Run tests verbose
go test ./cmd -run TestProject  # Run specific test
go vet ./...               # Run static analysis
gofmt -w .                # Format all Go files
```

---

## FAQ

**Q: Should I write unit tests or integration tests?**
A: For CLI commands, always write integration tests with in-memory databases. For internal packages, write unit tests.

**Q: How do I handle commands that call os.Exit(1)?**
A: Skip error path tests or test the underlying logic in the manager/storage layer instead.

**Q: Can I add new dependencies?**
A: Only if necessary. Check if standard library or existing dependencies can solve the problem first.

**Q: How do I update the database schema?**
A: Add migration logic. For now, the app recreates schema on startup - add the table to `initSchema()`.

**Q: Where should I put my tests?**
A: CLI command tests go in `cmd/*_test.go`. Internal package tests go in `internal/*/*_test.go`.

**Q: What if the design docs conflict with existing code?**
A: Existing code takes precedence. Update docs to reflect current implementation, or discuss with maintainers.

---

## Getting Help

- **Design Documentation**: See `docs/` directory
- **Code Examples**: Look at existing commands in `cmd/`
- **Test Examples**: See `cmd/project_test.go` and `testutil/testutil.go`
- **Quick Reference**: `docs/ISSUE_FLOW_QUICKREF.md`
- **Issues**: https://github.com/paolorechia/issue-flow/issues

---

**Remember**: Quality code is tested code. Always ensure tests pass before committing!
