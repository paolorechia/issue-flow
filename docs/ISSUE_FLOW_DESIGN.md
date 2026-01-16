# Issue Flow: Multi-Project Workflow Management Tool
## Language Selection & Technical Design Document

**Version**: 1.0  
**Date**: January 15, 2026  
**Status**: Design Phase

---

## Executive Summary

A standalone CLI tool for managing GitHub issues, git worktrees, and implementation workflows across multiple projects. Generalizes the compliance workflow pattern into a universal developer productivity tool.

---

## Language Comparison & Selection

### Option 1: **Go** ⭐ RECOMMENDED

**Pros**:
- ✅ Single binary distribution (no runtime dependencies)
- ✅ Excellent CLI library ecosystem (cobra, viper, bubbletea)
- ✅ Fast compilation and execution
- ✅ Cross-platform builds (macOS, Linux, Windows)
- ✅ Strong concurrency for parallel operations
- ✅ Static typing with excellent tooling
- ✅ Small binary size (~10MB)
- ✅ Great for system tools (git operations, file I/O)

**Cons**:
- ❌ Steeper learning curve if unfamiliar
- ❌ Verbose error handling
- ❌ Less flexible than dynamic languages

**Best For**: Production-grade CLI tools, system utilities

**Examples**: Docker, Kubernetes, GitHub CLI, Terraform

---

### Option 2: **Rust**

**Pros**:
- ✅ Single binary distribution
- ✅ Maximum performance
- ✅ Memory safety guarantees
- ✅ Excellent CLI libraries (clap, tokio)
- ✅ Growing ecosystem

**Cons**:
- ❌ Steepest learning curve
- ❌ Slower compilation
- ❌ Smaller ecosystem than Go for CLI tools
- ❌ Overkill for this use case

**Best For**: Performance-critical tools, systems programming

**Examples**: ripgrep, bat, fd

---

### Option 3: **TypeScript/Node.js**

**Pros**:
- ✅ You already know it well
- ✅ Rich ecosystem (commander, inquirer, chalk)
- ✅ Fast development
- ✅ Easy to prototype

**Cons**:
- ❌ Requires Node.js runtime (not truly standalone)
- ❌ Slower startup time (~100-200ms)
- ❌ Larger distribution size
- ❌ pkg/nexe for binaries is less reliable

**Best For**: Rapid prototyping, JavaScript ecosystems

**Examples**: npm, yarn, prettier

---

### Option 4: **Python**

**Pros**:
- ✅ Rapid development
- ✅ Excellent libraries (click, rich, typer)
- ✅ Easy to read/write

**Cons**:
- ❌ Requires Python runtime
- ❌ Distribution complexity (PyInstaller/cx_Freeze)
- ❌ Slower than compiled languages
- ❌ Dependency management headaches

**Best For**: Data tools, automation scripts

---

### Option 5: **Deno/Bun** (TypeScript)

**Pros**:
- ✅ Single binary output
- ✅ TypeScript native
- ✅ Fast runtime
- ✅ Modern tooling

**Cons**:
- ❌ Newer ecosystem (less mature)
- ❌ Still evolving
- ❌ Smaller community

**Best For**: Modern TypeScript projects

---

## 🏆 Final Recommendation: **Go**

### Why Go Wins

1. **Distribution**: Single binary, no dependencies
2. **Performance**: Fast startup, efficient execution
3. **Tooling**: `cobra` (CLI), `viper` (config), `bubbletea` (TUI)
4. **Ecosystem**: Proven for dev tools (gh, docker, terraform)
5. **Cross-platform**: Easy builds for all platforms
6. **Production-ready**: Battle-tested in similar tools

### Go Package Ecosystem

```go
// CLI framework
"github.com/spf13/cobra"      // Command structure
"github.com/spf13/viper"      // Configuration
"github.com/manifoldco/promptui" // Interactive prompts
"github.com/charmbracelet/bubbletea" // Terminal UI (optional)

// Git operations
"github.com/go-git/go-git/v5" // Pure Go git implementation

// GitHub API
"github.com/google/go-github/v58/github"
"golang.org/x/oauth2"

// Utilities
"github.com/fatih/color"      // Colored output
"github.com/olekukonko/tablewriter" // Tables
"gopkg.in/yaml.v3"           // YAML config
```

---

## Architecture Design

### High-Level Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Issue Flow CLI                      │
├─────────────────────────────────────────────────────┤
│                                                      │
│  ┌────────────┐  ┌────────────┐  ┌──────────────┐ │
│  │  Commands  │  │   Config   │  │   Storage    │ │
│  └────────────┘  └────────────┘  └──────────────┘ │
│         │              │                  │         │
│  ┌──────▼──────────────▼──────────────────▼─────┐ │
│  │              Core Services                     │ │
│  │  - Project Manager                             │ │
│  │  - Issue Manager                               │ │
│  │  - Worktree Manager                            │ │
│  │  - Template Engine                             │ │
│  └────────────────────────────────────────────────┘ │
│         │              │                  │         │
│  ┌──────▼──────┐ ┌────▼─────┐  ┌────────▼──────┐ │
│  │  GitHub API │ │   Git    │  │  OpenCode     │ │
│  └─────────────┘ └──────────┘  └───────────────┘ │
└─────────────────────────────────────────────────────┘
```

### Project Structure

```
issue-flow/
├── cmd/
│   ├── root.go              # Root command
│   ├── project.go           # Project commands
│   ├── issue.go             # Issue commands
│   ├── worktree.go          # Worktree commands
│   └── config.go            # Config commands
├── internal/
│   ├── config/
│   │   ├── config.go        # Config management
│   │   └── schema.go        # Config schema
│   ├── project/
│   │   ├── manager.go       # Project operations
│   │   └── types.go         # Project types
│   ├── issue/
│   │   ├── manager.go       # Issue operations
│   │   ├── template.go      # Template engine
│   │   └── types.go         # Issue types
│   ├── worktree/
│   │   ├── manager.go       # Worktree operations
│   │   └── types.go         # Worktree types
│   ├── github/
│   │   ├── client.go        # GitHub API client
│   │   └── auth.go          # Authentication
│   ├── git/
│   │   └── operations.go    # Git operations
│   ├── storage/
│   │   ├── db.go            # Local state storage
│   │   └── models.go        # Data models
│   └── ui/
│       ├── prompts.go       # Interactive prompts
│       └── table.go         # Table rendering
├── pkg/
│   └── templates/           # Built-in templates
├── docs/
│   ├── README.md
│   └── DESIGN.md
├── go.mod
├── go.sum
├── Makefile
└── main.go
```

---

## Data Model

### Multi-Project Management

```go
// Project represents a GitHub repository with workflow config
type Project struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    GitHubOwner string    `json:"github_owner"`
    GitHubRepo  string    `json:"github_repo"`
    LocalPath   string    `json:"local_path"`
    WorktreeDir string    `json:"worktree_dir"`
    Config      ProjectConfig `json:"config"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectConfig holds project-specific workflow configuration
type ProjectConfig struct {
    IssueTypes   []IssueType   `json:"issue_types"`
    Labels       LabelConfig   `json:"labels"`
    BranchConfig BranchConfig  `json:"branch_config"`
    OpenCode     OpenCodeConfig `json:"opencode"`
}

// IssueType defines a category of issues with templates
type IssueType struct {
    Name         string   `json:"name"`
    Label        string   `json:"label"`
    Priority     []string `json:"priority"`
    BranchPrefix string   `json:"branch_prefix"`
    Template     string   `json:"template"`
    GuidesDir    string   `json:"guides_dir"`
}

// Worktree represents an active git worktree
type Worktree struct {
    ID          string    `json:"id"`
    ProjectID   string    `json:"project_id"`
    IssueNumber int       `json:"issue_number"`
    Path        string    `json:"path"`
    Branch      string    `json:"branch"`
    CreatedAt   time.Time `json:"created_at"`
    Status      string    `json:"status"` // active, completed, abandoned
}
```

---

## Configuration Files

### Global Config: `~/.issue-flow/config.yaml`

```yaml
version: "1.0"

# Global settings
settings:
  editor: "code"              # Default editor
  opencode_enabled: true
  worktree_base: "~/issue-worktrees"

# GitHub authentication
github:
  auth_method: "gh_cli"      # gh_cli, token, oauth
  token: ""                   # Optional personal access token

# Projects registry
projects:
  - id: "whisper-web-server"
    name: "Whisper Web Server"
    github_owner: "Whisper-Notes"
    github_repo: "whisper-web-server"
    local_path: "~/dev/whisper-web-server"
    worktree_dir: "~/whisper-compliance-work"
  
  - id: "whisper-desktop"
    name: "Whisper Desktop App"
    github_owner: "Whisper-Notes"
    github_repo: "whisper-desktop"
    local_path: "~/dev/whisper-desktop"
    worktree_dir: "~/whisper-desktop-work"
```

### Project Config: `<repo>/.issue-flow.yaml`

```yaml
version: "1.0"

# Issue type definitions
issue_types:
  - name: "compliance"
    label: "Compliance Issue"
    priority: ["critical", "high", "medium"]
    branch_prefix: "compliance"
    template: "templates/compliance.md"
    guides_dir: "compliance-prompts"
    labels: ["compliance", "legal-risk"]
  
  - name: "feature"
    label: "Feature"
    priority: ["high", "medium", "low"]
    branch_prefix: "feature"
    template: "templates/feature.md"
    guides_dir: "feature-guides"
    labels: ["enhancement"]
  
  - name: "bug"
    label: "Bug Fix"
    priority: ["critical", "high", "medium", "low"]
    branch_prefix: "fix"
    template: "templates/bug.md"
    guides_dir: "bug-guides"
    labels: ["bug"]

# Branch naming
branch:
  pattern: "{prefix}/{issue-number}-{slug}"
  max_slug_length: 50

# Worktree settings
worktree:
  base_dir: "~/issue-worktrees/{project-id}"
  copy_guides: true
  create_context: true

# OpenCode integration
opencode:
  enabled: true
  auto_launch: false
  context_file: ".opencode-context"
  context_template: |
    # Issue #{issue_number}: {title}
    
    This worktree is for implementing {issue_type} issue #{issue_number}.
    
    ## Quick Start
    1. Review IMPLEMENTATION_GUIDE.md
    2. Follow step-by-step tasks
    3. Run tests as you implement
    
    ## Context
    - Project: {project_name}
    - Branch: {branch}
    - Priority: {priority}
```

---

## Command Design

### Project Management

```bash
# Initialize issue-flow in current directory
$ issue-flow init
? Project name: Whisper Web Server
? GitHub owner: Whisper-Notes
? GitHub repo: whisper-web-server
✓ Created .issue-flow.yaml
✓ Created templates/ directory
✓ Registered project globally

# Add existing project
$ issue-flow project add
? Project name: Whisper Desktop
? GitHub repo: Whisper-Notes/whisper-desktop
? Local path: ~/dev/whisper-desktop
✓ Added project: whisper-desktop

# List all projects
$ issue-flow project list
┌──────────────────────┬─────────────────────────────┬────────────┐
│ Project              │ Repository                  │ Worktrees  │
├──────────────────────┼─────────────────────────────┼────────────┤
│ whisper-web-server   │ Whisper-Notes/whisper-web-… │ 3 active   │
│ whisper-desktop      │ Whisper-Notes/whisper-desk… │ 1 active   │
└──────────────────────┴─────────────────────────────┴────────────┘

# Switch active project
$ issue-flow project use whisper-desktop
✓ Switched to project: whisper-desktop

# Show project details
$ issue-flow project info
Project: whisper-web-server
Repository: Whisper-Notes/whisper-web-server
Local Path: ~/dev/whisper-web-server
Worktree Dir: ~/whisper-compliance-work

Issue Types:
  - compliance (3 templates)
  - feature (2 templates)
  - bug (1 template)

Active Worktrees: 3
  - #126: compliance/critical-stripe-cancellation
  - #127: compliance/critical-data-export
  - #134: feature/134-oauth-desktop
```

---

### Issue Management

```bash
# Create issue (interactive)
$ issue-flow issue create
? Select project: whisper-web-server
? Issue type: compliance
? Title: GDPR data export endpoint
? Priority: critical
? Description: [opens editor]
? Create worktree now? Yes
✓ Created issue #127
✓ Created worktree at ~/whisper-compliance-work/issue-127

# Create issue (flags)
$ issue-flow issue create \
    --project whisper-web-server \
    --type compliance \
    --title "GDPR data export" \
    --priority critical \
    --worktree

# Create from template file
$ issue-flow issue create --from template.yaml

# List issues
$ issue-flow issue list
$ issue-flow issue list --project whisper-web-server
$ issue-flow issue list --type compliance --priority critical

# Show issue details
$ issue-flow issue show 127
```

---

### Worktree Management

```bash
# Start work on existing issue
$ issue-flow start 127
$ issue-flow start 127 --project whisper-web-server

# List worktrees
$ issue-flow worktree list
$ issue-flow worktree list --project whisper-web-server

# Show worktree status
$ issue-flow status 127

# Clean up worktree
$ issue-flow cleanup 127
$ issue-flow cleanup --all --completed

# Switch to worktree
$ issue-flow worktree switch 127
# (Changes directory in shell - needs shell integration)
```

---

### Global Commands

```bash
# Show overall status
$ issue-flow status

# Search across all projects
$ issue-flow search "oauth"

# Show statistics
$ issue-flow stats

# Update configuration
$ issue-flow config set github.token "ghp_xxxx"
$ issue-flow config get github.auth_method
```

---

## Storage Strategy

### Option 1: SQLite (Recommended)

**Pros**:
- ✅ Embedded database (no server)
- ✅ SQL queries for complex filtering
- ✅ Reliable and fast
- ✅ Easy backup (single file)

**Schema**:
```sql
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    github_owner TEXT NOT NULL,
    github_repo TEXT NOT NULL,
    local_path TEXT NOT NULL,
    worktree_dir TEXT NOT NULL,
    config TEXT NOT NULL, -- JSON
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

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

**Location**: `~/.issue-flow/database.db`

---

### Option 2: JSON Files

**Pros**:
- ✅ Simple implementation
- ✅ Human-readable
- ✅ Easy to edit manually

**Cons**:
- ❌ Slower for large datasets
- ❌ No complex queries
- ❌ Race conditions possible

**Structure**:
```
~/.issue-flow/
├── config.yaml           # Global config
├── projects/
│   ├── whisper-web-server.json
│   └── whisper-desktop.json
└── worktrees/
    ├── whisper-web-server/
    │   ├── 126.json
    │   └── 127.json
    └── whisper-desktop/
        └── 45.json
```

---

## Template System

### Template Format (Go text/template)

```markdown
# [{{.Priority}}] {{.Title}}

## Problem Statement

{{.Description}}

## Requirements

{{range .Requirements}}
- [ ] {{.}}
{{end}}

## Implementation Tasks

{{range .Tasks}}
### {{.Title}}
{{.Description}}

{{end}}

---

**Estimated Time**: {{.EstimatedTime}}
**Priority**: {{.Priority}}
**Branch**: {{.BranchName}}
**Issue**: #{{.IssueNumber}}
```

### Template Metadata (YAML)

```yaml
# templates/compliance.meta.yaml
name: "Compliance Issue"
description: "Template for GDPR/LGPD/CCPA compliance issues"

fields:
  - name: title
    type: string
    required: true
    prompt: "Issue title"
  
  - name: priority
    type: select
    options: ["CRITICAL", "HIGH", "MEDIUM", "LOW"]
    required: true
    default: "HIGH"
  
  - name: description
    type: text
    required: true
    prompt: "Detailed description"
  
  - name: requirements
    type: array
    prompt: "Add requirement (empty to finish)"
  
  - name: estimated_time
    type: string
    default: "2-4 hours"
    prompt: "Estimated time"

labels: ["compliance", "legal-risk"]
```

---

## Build & Distribution

### Build System (Makefile)

```makefile
.PHONY: build install test clean

# Build for current platform
build:
	go build -o bin/issue-flow main.go

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o bin/issue-flow-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/issue-flow-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/issue-flow-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/issue-flow-windows-amd64.exe main.go

# Install locally
install:
	go install

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
```

### Installation Methods

**Homebrew** (macOS/Linux):
```bash
brew tap whisper-notes/issue-flow
brew install issue-flow
```

**Direct Download**:
```bash
curl -L https://github.com/whisper-notes/issue-flow/releases/latest/download/issue-flow-$(uname -s)-$(uname -m) \
  -o /usr/local/bin/issue-flow
chmod +x /usr/local/bin/issue-flow
```

**Go Install**:
```bash
go install github.com/whisper-notes/issue-flow@latest
```

---

## Development Roadmap

### Phase 1: Core Foundation (Week 1)
- [ ] Set up Go project structure
- [ ] Implement CLI framework with cobra
- [ ] Implement configuration loading (viper)
- [ ] Create SQLite storage layer
- [ ] Implement project management commands

### Phase 2: Issue Management (Week 2)
- [ ] GitHub API integration
- [ ] Issue creation with templates
- [ ] Template engine implementation
- [ ] Issue listing and filtering
- [ ] Template metadata support

### Phase 3: Worktree Management (Week 3)
- [ ] Git worktree operations
- [ ] Branch creation logic
- [ ] Implementation guide copying
- [ ] Worktree status tracking
- [ ] Cleanup functionality

### Phase 4: Polish & Integration (Week 4)
- [ ] OpenCode integration
- [ ] Interactive prompts (promptui)
- [ ] Rich terminal output
- [ ] Error handling improvements
- [ ] Comprehensive testing

### Phase 5: Distribution (Week 5)
- [ ] Build pipeline (Makefile)
- [ ] Cross-platform builds
- [ ] Homebrew formula
- [ ] Documentation
- [ ] Release automation

---

## Key Features Summary

✅ **Multi-Project Support** - Manage multiple repositories  
✅ **Issue Templates** - Customizable issue types  
✅ **Git Worktrees** - Isolated development environments  
✅ **GitHub Integration** - Create/list/manage issues  
✅ **OpenCode Integration** - Launch AI coding sessions  
✅ **Cross-Platform** - macOS, Linux, Windows  
✅ **Single Binary** - No dependencies  
✅ **Fast** - Go performance  
✅ **Type-Safe** - Go static typing  

---

## Next Steps

1. **Review this design** - Approve language choice and architecture
2. **Bootstrap repository** - Create new GitHub repo
3. **Set up Go project** - Initialize with go.mod
4. **Start Phase 1** - Core foundation implementation

---

**Questions?**
- Is Go the right choice for you?
- Any architecture changes needed?
- Ready to bootstrap the new repository?

**Estimated Development Time**: 4-5 weeks (full-featured v1.0)
