# Issue Flow

A CLI tool for managing GitHub issues, git worktrees, and development workflows across multiple projects.

## Features

- ✅ Multi-project support
- ✅ Issue templates and management
- ✅ Git worktree automation
- ✅ GitHub integration
- ✅ OpenCode integration
- ✅ Cross-platform (macOS, Linux, Windows)

## Installation

```bash
# Via Go
go install github.com/paolorechia/issue-flow@latest

# Or download binary from releases
```

## Quick Start

```bash
# Initialize in your project
issue-flow init

# Create an issue
issue-flow issue create --type feature --title "Add new feature"

# Start working on an issue
issue-flow start 123

# List worktrees
issue-flow worktree list
```

## Documentation

See [docs/](docs/) for complete documentation:

- [Technical Design](docs/ISSUE_FLOW_DESIGN.md)
- [Developer Bootstrap Guide](docs/ISSUE_FLOW_BOOTSTRAP.md)
- [Quick Reference](docs/ISSUE_FLOW_QUICKREF.md)
- [Implementation Plan](docs/ISSUE_FLOW_IMPLEMENTATION_PLAN.md)

## Development

```bash
# Build
make build

# Test
make test

# Run
go run main.go version
```

## License

MIT
