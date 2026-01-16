# Issue Flow: Developer Bootstrap Guide

**Quick start guide for developers implementing the Issue Flow CLI tool**

---

## Prerequisites

Before starting, ensure you have:

- **Go 1.21+** installed (`go version`)
- **Git** installed and configured
- **GitHub CLI** (`gh`) installed and authenticated
- **Code editor** with Go support (VS Code + Go extension recommended)

---

## Quick Start (5 minutes)

```bash
# 1. Create repository
mkdir issue-flow && cd issue-flow

# 2. Initialize Go module
go mod init github.com/whisper-notes/issue-flow

# 3. Install core dependencies
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/manifoldco/promptui@latest
go get github.com/fatih/color@latest

# 4. Create basic structure
mkdir -p cmd internal/{config,project,issue,worktree,github,git,storage,ui} pkg/templates

# 5. Create main.go
cat > main.go << 'EOF'
package main

import "github.com/whisper-notes/issue-flow/cmd"

func main() {
    cmd.Execute()
}
EOF

# 6. Create root command
cat > cmd/root.go << 'EOF'
package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "issue-flow",
    Short: "Multi-project workflow management tool",
    Long:  "A CLI tool for managing GitHub issues, git worktrees, and workflows across multiple projects.",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("issue-flow v0.1.0")
    },
}
EOF

# 7. Build and test
go mod tidy
go build -o issue-flow
./issue-flow version
```

---

## Project Structure Setup

```bash
# Create all directories
mkdir -p {cmd,internal/{config,project,issue,worktree,github,git,storage,ui},pkg/templates,docs}

# Create placeholder files
touch cmd/{root,project,issue,worktree,config}.go
touch internal/config/{config,schema}.go
touch internal/project/{manager,types}.go
touch internal/issue/{manager,template,types}.go
touch internal/worktree/{manager,types}.go
touch internal/github/{client,auth}.go
touch internal/git/operations.go
touch internal/storage/{db,models}.go
touch internal/ui/{prompts,table}.go
```

Final structure:
```
issue-flow/
├── main.go
├── go.mod
├── go.sum
├── cmd/
│   ├── root.go       ✓ Create first
│   ├── project.go
│   ├── issue.go
│   ├── worktree.go
│   └── config.go
├── internal/
│   ├── config/
│   ├── project/
│   ├── issue/
│   ├── worktree/
│   ├── github/
│   ├── git/
│   ├── storage/
│   └── ui/
└── pkg/
    └── templates/
```

---

## Core Dependencies

Add to `go.mod`:

```bash
# CLI framework
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest

# Interactive UI
go get github.com/manifoldco/promptui@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/olekukonko/tablewriter@latest
go get github.com/fatih/color@latest

# Git operations
go get github.com/go-git/go-git/v5@latest

# GitHub API
go get github.com/google/go-github/v58/github@latest
go get golang.org/x/oauth2@latest

# Storage
go get github.com/mattn/go-sqlite3@latest

# YAML/JSON
go get gopkg.in/yaml.v3@latest

# Testing
go get github.com/stretchr/testify@latest
```

---

## Development Workflow

### Phase 1: Core CLI (Week 1)

**Goal**: Basic CLI with commands structure and configuration loading.

**Tasks**:
1. Implement `cmd/root.go` with cobra
2. Implement `cmd/project.go` with subcommands:
   - `issue-flow project list`
   - `issue-flow project add`
   - `issue-flow project use`
3. Implement `internal/config/config.go` with viper
4. Implement `internal/storage/db.go` with SQLite

**Test**:
```bash
go build
./issue-flow project list
./issue-flow project add
```

---

### Phase 2: Issue Management (Week 2)

**Goal**: Create and list GitHub issues.

**Tasks**:
1. Implement `internal/github/client.go`
2. Implement `internal/issue/manager.go`
3. Implement `cmd/issue.go` with:
   - `issue-flow issue create`
   - `issue-flow issue list`
4. Implement template engine in `internal/issue/template.go`

**Test**:
```bash
./issue-flow issue create --type feature --title "Test issue"
./issue-flow issue list
```

---

### Phase 3: Worktree Management (Week 3)

**Goal**: Create and manage git worktrees.

**Tasks**:
1. Implement `internal/git/operations.go`
2. Implement `internal/worktree/manager.go`
3. Implement `cmd/worktree.go`
4. Add `issue-flow start` command

**Test**:
```bash
./issue-flow start 123
cd ~/issue-worktrees/test-project/issue-123
git status
```

---

## Code Examples

### Example 1: Root Command (`cmd/root.go`)

```go
package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgFile string
    verbose bool
)

var rootCmd = &cobra.Command{
    Use:   "issue-flow",
    Short: "Multi-project workflow management",
    Long: `Issue Flow is a CLI tool for managing GitHub issues,
git worktrees, and development workflows across multiple projects.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.issue-flow/config.yaml)")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
    
    viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)
        
        viper.AddConfigPath(home + "/.issue-flow")
        viper.SetConfigType("yaml")
        viper.SetConfigName("config")
    }
    
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err == nil && verbose {
        fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    }
}
```

---

### Example 2: Project Types (`internal/project/types.go`)

```go
package project

import "time"

// Project represents a GitHub repository with workflow config
type Project struct {
    ID          string        `json:"id" yaml:"id"`
    Name        string        `json:"name" yaml:"name"`
    GitHubOwner string        `json:"github_owner" yaml:"github_owner"`
    GitHubRepo  string        `json:"github_repo" yaml:"github_repo"`
    LocalPath   string        `json:"local_path" yaml:"local_path"`
    WorktreeDir string        `json:"worktree_dir" yaml:"worktree_dir"`
    Config      ProjectConfig `json:"config" yaml:"config"`
    CreatedAt   time.Time     `json:"created_at" yaml:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at" yaml:"updated_at"`
}

// ProjectConfig holds project-specific workflow configuration
type ProjectConfig struct {
    IssueTypes   []IssueType   `json:"issue_types" yaml:"issue_types"`
    BranchConfig BranchConfig  `json:"branch_config" yaml:"branch_config"`
    OpenCode     OpenCodeConfig `json:"opencode" yaml:"opencode"`
}

// IssueType defines a category of issues
type IssueType struct {
    Name         string   `json:"name" yaml:"name"`
    Label        string   `json:"label" yaml:"label"`
    Priority     []string `json:"priority" yaml:"priority"`
    BranchPrefix string   `json:"branch_prefix" yaml:"branch_prefix"`
    Template     string   `json:"template" yaml:"template"`
    GuidesDir    string   `json:"guides_dir" yaml:"guides_dir"`
}

// BranchConfig defines branch naming patterns
type BranchConfig struct {
    Pattern      string `json:"pattern" yaml:"pattern"`
    MaxSlugLength int   `json:"max_slug_length" yaml:"max_slug_length"`
}

// OpenCodeConfig defines OpenCode integration settings
type OpenCodeConfig struct {
    Enabled         bool   `json:"enabled" yaml:"enabled"`
    AutoLaunch      bool   `json:"auto_launch" yaml:"auto_launch"`
    ContextFile     string `json:"context_file" yaml:"context_file"`
    ContextTemplate string `json:"context_template" yaml:"context_template"`
}
```

---

### Example 3: SQLite Storage (`internal/storage/db.go`)

```go
package storage

import (
    "database/sql"
    "path/filepath"
    "os"
    _ "github.com/mattn/go-sqlite3"
)

type Database struct {
    db *sql.DB
}

// New creates or opens the SQLite database
func New() (*Database, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    
    dbDir := filepath.Join(home, ".issue-flow")
    if err := os.MkdirAll(dbDir, 0755); err != nil {
        return nil, err
    }
    
    dbPath := filepath.Join(dbDir, "database.db")
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    d := &Database{db: db}
    if err := d.initSchema(); err != nil {
        return nil, err
    }
    
    return d, nil
}

func (d *Database) initSchema() error {
    schema := `
    CREATE TABLE IF NOT EXISTS projects (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        github_owner TEXT NOT NULL,
        github_repo TEXT NOT NULL,
        local_path TEXT NOT NULL,
        worktree_dir TEXT NOT NULL,
        config TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TABLE IF NOT EXISTS worktrees (
        id TEXT PRIMARY KEY,
        project_id TEXT NOT NULL,
        issue_number INTEGER NOT NULL,
        path TEXT NOT NULL,
        branch TEXT NOT NULL,
        status TEXT NOT NULL DEFAULT 'active',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (project_id) REFERENCES projects(id)
    );
    
    CREATE TABLE IF NOT EXISTS issue_cache (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        project_id TEXT NOT NULL,
        issue_number INTEGER NOT NULL,
        title TEXT NOT NULL,
        type TEXT,
        priority TEXT,
        status TEXT,
        cached_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (project_id) REFERENCES projects(id),
        UNIQUE(project_id, issue_number)
    );
    `
    
    _, err := d.db.Exec(schema)
    return err
}

func (d *Database) Close() error {
    return d.db.Close()
}
```

---

## Testing Strategy

### Unit Tests

Create `*_test.go` files alongside code:

```go
// internal/project/manager_test.go
package project

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestProjectValidation(t *testing.T) {
    p := Project{
        ID: "test-project",
        Name: "Test Project",
        GitHubOwner: "test-owner",
        GitHubRepo: "test-repo",
    }
    
    err := p.Validate()
    assert.NoError(t, err)
}
```

Run tests:
```bash
go test ./...
go test -v ./internal/project
go test -cover ./...
```

---

### Integration Tests

Test against real GitHub API (use test repositories):

```go
// internal/github/client_test.go
//go:build integration
package github

import (
    "testing"
    "os"
)

func TestCreateIssue(t *testing.T) {
    if os.Getenv("GITHUB_TOKEN") == "" {
        t.Skip("Skipping integration test: GITHUB_TOKEN not set")
    }
    
    client := NewClient(os.Getenv("GITHUB_TOKEN"))
    // Test issue creation
}
```

Run integration tests:
```bash
GITHUB_TOKEN=ghp_xxx go test -tags=integration ./...
```

---

## Build & Release

### Local Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install
```

### Makefile

```makefile
.PHONY: build build-all install test clean

VERSION := 0.1.0
BINARY := issue-flow

build:
	go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY) main.go

build-all:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY)-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY)-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY)-windows-amd64.exe main.go

install:
	go install -ldflags="-X main.version=$(VERSION)"

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/
	rm -f coverage.out
```

---

## Debugging Tips

### Enable Verbose Logging

```go
// internal/ui/logger.go
package ui

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/spf13/viper"
)

func Info(msg string, args ...interface{}) {
    fmt.Printf("ℹ️  "+msg+"\n", args...)
}

func Success(msg string, args ...interface{}) {
    color.Green("✓ " + msg, args...)
    fmt.Println()
}

func Error(msg string, args ...interface{}) {
    color.Red("✗ " + msg, args...)
    fmt.Println()
}

func Debug(msg string, args ...interface{}) {
    if viper.GetBool("verbose") {
        color.Cyan("[DEBUG] " + msg, args...)
        fmt.Println()
    }
}
```

Usage:
```bash
issue-flow --verbose project list
```

---

## Common Pitfalls

### 1. CGO_ENABLED for SQLite

SQLite requires CGO. When cross-compiling:

```bash
# macOS → Linux requires cross-compiler
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build
```

**Solution**: Use pure-Go alternatives like `modernc.org/sqlite` or build on native platform.

### 2. Path Handling

Always use `filepath.Join()` for cross-platform paths:

```go
// ✓ Good
path := filepath.Join(home, ".issue-flow", "config.yaml")

// ✗ Bad (breaks on Windows)
path := home + "/.issue-flow/config.yaml"
```

### 3. Error Handling

Always wrap errors with context:

```go
if err != nil {
    return fmt.Errorf("failed to create worktree: %w", err)
}
```

---

## Resources

- **Go CLI Tutorial**: https://cobra.dev
- **Viper Config**: https://github.com/spf13/viper
- **Go-Git**: https://github.com/go-git/go-git
- **GitHub API**: https://docs.github.com/en/rest
- **SQLite**: https://github.com/mattn/go-sqlite3
- **Table Writer**: https://github.com/olekukonko/tablewriter

---

## Getting Help

- **Design Doc**: See `docs/ISSUE_FLOW_DESIGN.md`
- **Architecture**: See design doc sections
- **Examples**: See `examples/` directory (create these)
- **Issues**: Open GitHub issue

---

**Ready to start? Begin with Phase 1 and build incrementally!**
