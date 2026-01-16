# Issue Flow Documentation

**Complete documentation for the Issue Flow multi-project workflow management tool**

---

## 📚 Documentation Overview

This directory contains complete design and implementation documentation for **Issue Flow** - a CLI tool for managing GitHub issues, git worktrees, and development workflows across multiple projects.

---

## 📄 Documents

### 1. [ISSUE_FLOW_DESIGN.md](./ISSUE_FLOW_DESIGN.md) - Technical Design Document

**For**: Architects, Senior Developers, Decision Makers

**Contents**:
- Language comparison and selection (Go vs Rust vs TypeScript vs Python)
- High-level architecture diagram
- Complete data models (Go structs)
- Configuration schema (YAML examples)
- Full command design with examples
- Storage strategy (SQLite schema)
- Template system design
- Build & distribution strategy
- 5-week development roadmap

**Use when**: Making architectural decisions, understanding system design

---

### 2. [ISSUE_FLOW_BOOTSTRAP.md](./ISSUE_FLOW_BOOTSTRAP.md) - Developer Bootstrap Guide

**For**: Developers implementing the tool

**Contents**:
- 5-minute quick start guide
- Complete project structure setup
- All dependencies with install commands
- Phase-by-phase development workflow
- Complete code examples (root command, types, database)
- Testing strategy (unit + integration)
- Build & release process
- Common pitfalls and solutions
- Debugging tips

**Use when**: Starting development, implementing features

---

### 3. [ISSUE_FLOW_QUICKREF.md](./ISSUE_FLOW_QUICKREF.md) - Quick Reference Card

**For**: Developers during daily development

**Contents**:
- Command structure diagram
- File locations reference table
- Common command examples
- Configuration examples (global + project)
- Go code patterns (commands, database, prompts)
- SQL schema reference
- Build commands
- Template variables
- Error handling patterns
- Release checklist

**Use when**: Quick lookups, daily development, troubleshooting

---

## 🚀 Getting Started

### If you're building Issue Flow:

1. **Start here**: [ISSUE_FLOW_BOOTSTRAP.md](./ISSUE_FLOW_BOOTSTRAP.md)
   - Follow the 5-minute quick start
   - Set up your development environment
   - Build the initial version

2. **Reference**: [ISSUE_FLOW_DESIGN.md](./ISSUE_FLOW_DESIGN.md)
   - Understand the architecture
   - Review data models
   - Check command specifications

3. **Daily use**: [ISSUE_FLOW_QUICKREF.md](./ISSUE_FLOW_QUICKREF.md)
   - Command syntax
   - Code patterns
   - Configuration examples

---

## 🎯 Key Decisions

### Language: **Go**
- Single binary distribution (no runtime dependencies)
- Fast startup and execution
- Proven for CLI tools (GitHub CLI, Docker, Terraform)
- Cross-platform builds

### Storage: **SQLite**
- Embedded database (no server)
- Fast queries for filtering/searching
- Single file backup
- Reliable and battle-tested

### Architecture: **Cobra + Viper**
- Industry-standard CLI framework
- Configuration management with viper
- Rich ecosystem of UI libraries

---

## 📊 Project Structure

```
issue-flow/                           # New repository
├── docs/
│   ├── ISSUE_FLOW_DESIGN.md         # ← Technical design
│   ├── ISSUE_FLOW_BOOTSTRAP.md      # ← Developer guide
│   └── ISSUE_FLOW_QUICKREF.md       # ← Quick reference
├── cmd/                              # CLI commands
│   ├── root.go                       # Main entry point
│   ├── project.go                    # Project management
│   ├── issue.go                      # Issue operations
│   ├── worktree.go                   # Worktree management
│   └── config.go                     # Configuration
├── internal/                         # Core implementation
│   ├── config/                       # Config loading
│   ├── project/                      # Project manager
│   ├── issue/                        # Issue manager + templates
│   ├── worktree/                     # Worktree manager
│   ├── github/                       # GitHub API client
│   ├── git/                          # Git operations
│   ├── storage/                      # SQLite database
│   └── ui/                           # Terminal UI (prompts, tables)
├── pkg/
│   └── templates/                    # Built-in templates
├── main.go                           # Entry point
├── go.mod
├── go.sum
└── Makefile                          # Build system
```

---

## 🛠️ Quick Commands

```bash
# Bootstrap new project
mkdir issue-flow && cd issue-flow
go mod init github.com/whisper-notes/issue-flow
# ... follow ISSUE_FLOW_BOOTSTRAP.md

# Build
make build

# Test
make test

# Install locally
make install

# Use
issue-flow project add
issue-flow issue create --type feature --worktree
issue-flow start 123
```

---

## 📅 Development Timeline

| Phase | Duration | Focus | Document Section |
|-------|----------|-------|------------------|
| Phase 1 | Week 1 | Core CLI, config, storage | Bootstrap → Phase 1 |
| Phase 2 | Week 2 | Issue management, templates | Bootstrap → Phase 2 |
| Phase 3 | Week 3 | Worktree operations | Bootstrap → Phase 3 |
| Phase 4 | Week 4 | Polish, OpenCode integration | Design → Roadmap |
| Phase 5 | Week 5 | Distribution, docs | Design → Distribution |

**Total**: 4-5 weeks for v1.0

---

## ✅ What's Included

All documents provide:

- ✅ **Complete working examples** - Copy-paste ready code
- ✅ **Step-by-step instructions** - No guesswork
- ✅ **Best practices** - Proven Go patterns
- ✅ **Error handling** - Robust error patterns
- ✅ **Testing strategy** - Unit + integration tests
- ✅ **Build automation** - Makefile with all targets
- ✅ **Cross-platform** - macOS, Linux, Windows
- ✅ **Distribution** - Homebrew, direct download, go install

---

## 🎓 Learning Path

### For Backend Developers (TypeScript/Node.js background):

1. **Quick Go Primer**: https://go.dev/tour/
2. **Read**: ISSUE_FLOW_BOOTSTRAP.md (familiar patterns)
3. **Start**: Implement Phase 1 (core CLI)
4. **Reference**: Code examples in Bootstrap doc

### For Go Developers:

1. **Read**: ISSUE_FLOW_DESIGN.md (architecture)
2. **Skim**: ISSUE_FLOW_BOOTSTRAP.md (setup)
3. **Use**: ISSUE_FLOW_QUICKREF.md (daily reference)
4. **Start**: Implement Phase 1

---

## 🤔 FAQ

**Q: Why Go instead of TypeScript?**  
A: Single binary distribution, no runtime dependencies, proven for CLI tools.

**Q: Why SQLite instead of JSON files?**  
A: Fast queries, reliable, easy backup, handles complex filtering.

**Q: Can I use this for non-GitHub projects?**  
A: Design supports it, but v1.0 focuses on GitHub. GitLab/Bitbucket support in v2.0.

**Q: How big will the binary be?**  
A: ~10-15MB (Go binary + SQLite driver)

**Q: Do I need to know Go?**  
A: Helpful but not required. Bootstrap doc provides all code examples.

---

## 📞 Support

- **Design questions**: See [ISSUE_FLOW_DESIGN.md](./ISSUE_FLOW_DESIGN.md)
- **Implementation help**: See [ISSUE_FLOW_BOOTSTRAP.md](./ISSUE_FLOW_BOOTSTRAP.md)
- **Quick lookups**: See [ISSUE_FLOW_QUICKREF.md](./ISSUE_FLOW_QUICKREF.md)
- **Issues**: Open GitHub issue (once repo created)

---

## 🎯 Ready to Build?

1. ✅ Review design: [ISSUE_FLOW_DESIGN.md](./ISSUE_FLOW_DESIGN.md)
2. ✅ Follow bootstrap: [ISSUE_FLOW_BOOTSTRAP.md](./ISSUE_FLOW_BOOTSTRAP.md)
3. ✅ Reference quickref: [ISSUE_FLOW_QUICKREF.md](./ISSUE_FLOW_QUICKREF.md)
4. 🚀 Start coding!

---

**Questions?** Read the docs above or open an issue.

**Ready?** Start with the Bootstrap guide and build incrementally!
