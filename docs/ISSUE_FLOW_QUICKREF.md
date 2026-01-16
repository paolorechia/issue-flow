# Issue Flow - Quick Reference Card

**Command structure, file locations, and key patterns**

---

## Command Structure

```
issue-flow
├── project          # Manage projects
│   ├── add          # Add new project
│   ├── list         # List all projects
│   ├── use          # Switch active project
│   ├── info         # Show project details
│   └── remove       # Remove project
├── issue            # Manage issues
│   ├── create       # Create new issue
│   ├── list         # List issues
│   └── show         # Show issue details
├── start [number]   # Start work on issue (creates worktree)
├── worktree         # Manage worktrees
│   ├── list         # List worktrees
│   └── remove       # Remove worktree
├── cleanup          # Clean up worktrees
├── status           # Show status
├── config           # Manage configuration
│   ├── get          # Get config value
│   └── set          # Set config value
└── version          # Show version
```

---

## File Locations

| File | Location | Purpose |
|------|----------|---------|
| Global config | `~/.issue-flow/config.yaml` | User settings, projects list |
| Database | `~/.issue-flow/database.db` | SQLite database |
| Project config | `<repo>/.issue-flow.yaml` | Project-specific settings |
| Templates | `<repo>/templates/` | Issue templates |
| Guides | `<repo>/<guides_dir>/` | Implementation guides |
| Worktrees | `~/issue-worktrees/<project>/issue-<N>/` | Git worktrees |

---

## Common Commands

```bash
# Initialize in new project
cd ~/dev/my-project
issue-flow init

# Create issue with worktree
issue-flow issue create --type feature --title "New feature" --worktree

# Start work on existing issue
issue-flow start 123

# List all worktrees
issue-flow worktree list

# Show status
issue-flow status 123

# Clean up completed work
issue-flow cleanup --completed

# Switch projects
issue-flow project use my-other-project
```

---

## Configuration Examples

### Global Config (`~/.issue-flow/config.yaml`)

```yaml
version: "1.0"
settings:
  editor: "code"
  opencode_enabled: true
  worktree_base: "~/issue-worktrees"
github:
  auth_method: "gh_cli"
projects:
  - id: "my-project"
    name: "My Project"
    github_owner: "myorg"
    github_repo: "my-repo"
    local_path: "~/dev/my-project"
```

### Project Config (`.issue-flow.yaml`)

```yaml
version: "1.0"
issue_types:
  - name: "feature"
    branch_prefix: "feature"
    template: "templates/feature.md"
    labels: ["enhancement"]
branch:
  pattern: "{prefix}/{issue-number}-{slug}"
opencode:
  enabled: true
  auto_launch: false
```

---

## Go Code Patterns

### Command Structure

```go
// cmd/mycommand.go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Short description",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
    myCmd.Flags().StringP("flag", "f", "", "Description")
}
```

### Database Query

```go
// internal/storage/projects.go
func (d *Database) GetProject(id string) (*Project, error) {
    var p Project
    row := d.db.QueryRow("SELECT * FROM projects WHERE id = ?", id)
    err := row.Scan(&p.ID, &p.Name, &p.GitHubOwner, &p.GitHubRepo, &p.LocalPath, &p.WorktreeDir, &p.Config, &p.CreatedAt, &p.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("project not found: %w", err)
    }
    return &p, nil
}
```

### Interactive Prompt

```go
// internal/ui/prompts.go
func PromptSelect(label string, options []string) (string, error) {
    prompt := promptui.Select{
        Label: label,
        Items: options,
    }
    _, result, err := prompt.Run()
    return result, err
}
```

---

## SQL Schema

```sql
-- Projects
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    github_owner TEXT NOT NULL,
    github_repo TEXT NOT NULL,
    local_path TEXT NOT NULL,
    worktree_dir TEXT NOT NULL,
    config TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Worktrees
CREATE TABLE worktrees (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    issue_number INTEGER NOT NULL,
    path TEXT NOT NULL,
    branch TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- Issue Cache
CREATE TABLE issue_cache (
    id INTEGER PRIMARY KEY,
    project_id TEXT NOT NULL,
    issue_number INTEGER NOT NULL,
    title TEXT NOT NULL,
    type TEXT,
    priority TEXT,
    status TEXT,
    cached_at TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    UNIQUE(project_id, issue_number)
);
```

---

## Build Commands

```bash
# Development
go run main.go version
go run main.go project list

# Build
make build                # Current platform
make build-all            # All platforms
make install              # Install to $GOPATH/bin

# Test
make test                 # Run tests
make test-coverage        # Coverage report

# Clean
make clean                # Remove binaries
```

---

## Dependencies

```bash
# Core CLI
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest

# UI
go get github.com/manifoldco/promptui@latest
go get github.com/fatih/color@latest
go get github.com/olekukonko/tablewriter@latest

# Git & GitHub
go get github.com/go-git/go-git/v5@latest
go get github.com/google/go-github/v58/github@latest

# Storage
go get github.com/mattn/go-sqlite3@latest

# Config
go get gopkg.in/yaml.v3@latest
```

---

## Branch Naming Pattern

Default: `{prefix}/{issue-number}-{slug}`

Examples:
- `feature/123-add-oauth-support`
- `fix/456-crash-on-startup`
- `compliance/789-gdpr-export`

Customize in `.issue-flow.yaml`:
```yaml
branch:
  pattern: "{prefix}/{slug}-{issue-number}"
  max_slug_length: 50
```

---

## Template Variables

Available in issue templates:

```markdown
{{.Title}}          - Issue title
{{.Description}}    - Issue description
{{.Priority}}       - Priority level
{{.IssueNumber}}    - GitHub issue number
{{.BranchName}}     - Generated branch name
{{.EstimatedTime}}  - Estimated completion time
{{.ProjectName}}    - Project name

{{range .Requirements}}
- [ ] {{.}}
{{end}}
```

---

## Error Handling Pattern

```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create worktree: %w", err)
}

// Log errors before returning
if err != nil {
    ui.Error("Failed to create issue: %s", err)
    return err
}
```

---

## Testing Patterns

```go
// Unit test
func TestProjectValidation(t *testing.T) {
    p := &Project{ID: "test"}
    assert.NoError(t, p.Validate())
}

// Integration test (requires build tag)
//go:build integration
func TestGitHubAPI(t *testing.T) {
    // Test real API
}
```

Run:
```bash
go test ./...                        # Unit tests only
go test -tags=integration ./...     # Include integration tests
```

---

## Useful Git Commands

```bash
# List worktrees
git worktree list

# Add worktree
git worktree add <path> -b <branch>

# Remove worktree
git worktree remove <path>

# Prune deleted worktrees
git worktree prune
```

---

## Debugging

```bash
# Verbose mode
issue-flow --verbose project list

# Check config
issue-flow config get github.auth_method

# Check database
sqlite3 ~/.issue-flow/database.db "SELECT * FROM projects;"
```

---

## Release Checklist

- [ ] Update version in `main.go`
- [ ] Run tests: `make test`
- [ ] Build all platforms: `make build-all`
- [ ] Test binaries on each platform
- [ ] Create Git tag: `git tag v0.1.0`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] Create GitHub release
- [ ] Upload binaries to release
- [ ] Update Homebrew formula (if applicable)

---

## Support

- Design: `docs/ISSUE_FLOW_DESIGN.md`
- Bootstrap: `docs/ISSUE_FLOW_BOOTSTRAP.md`
- GitHub: https://github.com/whisper-notes/issue-flow
