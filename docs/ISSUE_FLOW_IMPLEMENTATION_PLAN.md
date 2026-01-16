# Implementation Plan: Issue-Worktree CLI Tool

## Overview

A generalized CLI tool for creating GitHub issues and spawning git worktrees with implementation guides, inspired by the compliance workflow automation.

**Tool Name**: `issue-flow` (or `iflow`)

---

## Goals

1. **Generalize** the compliance workflow pattern for any project
2. **Simplify** issue creation with templates
3. **Automate** git worktree creation and management
4. **Integrate** with OpenCode for AI-assisted development
5. **Provide** a flexible, configuration-driven workflow

---

## Architecture

### Technology Stack

**Option A: Node.js/TypeScript** (Recommended)
- ✅ Already used in this project
- ✅ Rich ecosystem (commander.js, inquirer, chalk)
- ✅ Easy distribution via npm
- ✅ TypeScript for type safety

**Option B: Bash** (Current approach)
- ✅ No dependencies
- ❌ Limited for complex logic
- ❌ Harder to maintain/extend

**Decision**: Use **Node.js/TypeScript** for better maintainability and features

### Core Components

```
issue-flow/
├── src/
│   ├── cli.ts                 # Main CLI entry point
│   ├── commands/
│   │   ├── create.ts          # Create issue + worktree
│   │   ├── start.ts           # Start work on existing issue
│   │   ├── list.ts            # List issues/worktrees
│   │   ├── cleanup.ts         # Clean up worktrees
│   │   └── status.ts          # Show workflow status
│   ├── core/
│   │   ├── config.ts          # Configuration management
│   │   ├── github.ts          # GitHub API integration
│   │   ├── git.ts             # Git worktree operations
│   │   ├── template.ts        # Issue template rendering
│   │   └── opencode.ts        # OpenCode integration
│   ├── types/
│   │   └── index.ts           # TypeScript types
│   └── utils/
│       ├── logger.ts          # Colored logging
│       ├── prompt.ts          # Interactive prompts
│       └── validation.ts      # Input validation
├── templates/                 # Issue templates directory
│   ├── compliance/            # Compliance issue templates
│   ├── feature/               # Feature templates
│   └── bug/                   # Bug templates
├── .issue-flow.config.json    # Configuration file
├── package.json
├── tsconfig.json
└── README.md
```

---

## Configuration Schema

### `.issue-flow.config.json`

```json
{
  "$schema": "https://issue-flow.dev/schema.json",
  "version": "1.0",
  
  "github": {
    "owner": "Whisper-Notes",
    "repo": "whisper-web-server",
    "labels": {
      "default": ["enhancement"],
      "presets": {
        "compliance": ["compliance", "legal-risk"],
        "feature": ["enhancement"],
        "bug": ["bug"],
        "security": ["security", "high-priority"]
      }
    }
  },
  
  "worktree": {
    "baseDir": "~/issue-worktrees",
    "branchPrefix": "feature",
    "branchPattern": "{prefix}/{issue-number}-{slug}"
  },
  
  "templates": {
    "directory": "./templates",
    "guidesDirectory": "./guides"
  },
  
  "opencode": {
    "enabled": true,
    "autoLaunch": false,
    "contextFile": ".opencode-context"
  },
  
  "issueTypes": [
    {
      "name": "compliance",
      "label": "Compliance Issue",
      "priority": ["critical", "high", "medium"],
      "branchPrefix": "compliance",
      "template": "templates/compliance/template.md",
      "guidesDir": "compliance-prompts"
    },
    {
      "name": "feature",
      "label": "Feature",
      "priority": ["high", "medium", "low"],
      "branchPrefix": "feature",
      "template": "templates/feature/template.md",
      "guidesDir": "feature-guides"
    },
    {
      "name": "bug",
      "label": "Bug Fix",
      "priority": ["critical", "high", "medium", "low"],
      "branchPrefix": "fix",
      "template": "templates/bug/template.md",
      "guidesDir": "bug-guides"
    }
  ]
}
```

---

## CLI Commands

### 1. `issue-flow create`

Create a new issue with optional worktree.

```bash
# Interactive mode
issue-flow create

# With flags
issue-flow create \
  --type compliance \
  --title "GDPR data export endpoint" \
  --priority critical \
  --labels compliance,legal-risk \
  --worktree \
  --guide ./guides/data-export.md

# From template file
issue-flow create --template ./templates/compliance/data-export.json
```

**Flow**:
1. Prompt for issue type (compliance, feature, bug, custom)
2. Prompt for title, description, priority, labels
3. Render issue template with variables
4. Create GitHub issue via `gh` CLI
5. Optionally create git worktree
6. Copy implementation guide to worktree
7. Optionally launch OpenCode

---

### 2. `issue-flow start`

Start work on an existing issue.

```bash
# Interactive selection
issue-flow start

# Direct issue number
issue-flow start 126

# With options
issue-flow start 126 --no-opencode --guide ./guides/custom.md
```

**Flow**:
1. Fetch issue details from GitHub
2. Create git worktree
3. Create feature branch
4. Copy implementation guide
5. Create `.opencode-context` file
6. Optionally launch OpenCode

---

### 3. `issue-flow list`

List issues and worktrees.

```bash
# List all open issues
issue-flow list issues

# List all worktrees
issue-flow list worktrees

# List by type
issue-flow list issues --type compliance --priority critical

# List combined view
issue-flow list --combined
```

**Output**:
```
┌─────────┬──────────────────────────────────────┬──────────┬──────────┬────────────┐
│ Issue # │ Title                                │ Type     │ Priority │ Worktree   │
├─────────┼──────────────────────────────────────┼──────────┼──────────┼────────────┤
│ #126    │ Cancel Stripe subscription           │ Compliance│ CRITICAL │ ✓ Active   │
│ #127    │ User data export endpoint            │ Compliance│ CRITICAL │ ✗ None     │
│ #134    │ Update OAuth desktop clients         │ Feature  │ HIGH     │ ✓ Active   │
└─────────┴──────────────────────────────────────┴──────────┴──────────┴────────────┘
```

---

### 4. `issue-flow cleanup`

Clean up worktrees.

```bash
# Interactive cleanup
issue-flow cleanup

# Clean specific issue
issue-flow cleanup 126

# Clean all completed issues
issue-flow cleanup --completed

# Clean all
issue-flow cleanup --all
```

**Flow**:
1. List existing worktrees
2. Check issue status (open/closed)
3. Prompt for confirmation
4. Remove worktree
5. Optionally delete branch

---

### 5. `issue-flow status`

Show workflow status.

```bash
# Overall status
issue-flow status

# Issue-specific status
issue-flow status 126
```

**Output**:
```
Issue Flow Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Issue #126: Cancel Stripe subscription
  Status:     Open
  Type:       Compliance
  Priority:   CRITICAL
  Branch:     compliance/critical-stripe-cancellation
  Worktree:   ~/issue-worktrees/issue-126
  Guide:      IMPLEMENTATION_GUIDE.md
  Created:    2 days ago
  Updated:    1 hour ago

Git Status:
  ✓ Worktree exists
  ✓ Branch exists
  ⚠ Uncommitted changes (3 files)
  ✗ Not pushed to remote

Next Steps:
  1. Commit your changes
  2. Push to remote
  3. Create PR: gh pr create
```

---

### 6. `issue-flow init`

Initialize issue-flow in a new repository.

```bash
# Interactive setup
issue-flow init

# With options
issue-flow init --owner Whisper-Notes --repo my-repo
```

**Flow**:
1. Detect git repository
2. Prompt for GitHub owner/repo
3. Create `.issue-flow.config.json`
4. Create `templates/` directory structure
5. Create example templates

---

## Issue Template System

### Template Variables

Templates support Handlebars-style variables:

```markdown
# [{{priority}}] {{title}}

## Problem Statement

{{description}}

## Requirements

{{#each requirements}}
- [ ] {{this}}
{{/each}}

## Implementation Tasks

{{#each tasks}}
### {{this.title}}
{{this.description}}
{{/each}}

---

**Estimated Time**: {{estimatedTime}}
**Priority**: {{priority}}
**Branch**: {{branchName}}
```

### Template Metadata

Each template can have a companion `.meta.json` file:

```json
{
  "name": "Compliance Issue",
  "description": "Template for GDPR/LGPD/CCPA compliance issues",
  "fields": [
    {
      "name": "title",
      "type": "string",
      "required": true,
      "prompt": "Issue title"
    },
    {
      "name": "priority",
      "type": "select",
      "options": ["CRITICAL", "HIGH", "MEDIUM", "LOW"],
      "required": true
    },
    {
      "name": "estimatedTime",
      "type": "string",
      "default": "2-4 hours"
    },
    {
      "name": "requirements",
      "type": "array",
      "prompt": "Add requirement (leave empty to finish)"
    }
  ],
  "labels": ["compliance", "legal-risk"],
  "branchPrefix": "compliance"
}
```

---

## Implementation Phases

### Phase 1: Core CLI Framework (Week 1)

**Tasks**:
1. Set up TypeScript project with CLI framework (commander.js)
2. Implement configuration loading/validation
3. Implement basic logging and prompts
4. Create project structure

**Deliverables**:
- Basic CLI with `--help` and `--version`
- Configuration file loading
- Logger utility

---

### Phase 2: Issue Management (Week 2)

**Tasks**:
1. Integrate GitHub API (via `gh` CLI or Octokit)
2. Implement `issue-flow create` command
3. Implement template rendering system
4. Implement `issue-flow list issues` command

**Deliverables**:
- Create issues from templates
- List issues with filtering
- Template variable substitution

---

### Phase 3: Worktree Management (Week 2-3)

**Tasks**:
1. Implement git worktree operations
2. Implement `issue-flow start` command
3. Implement branch creation logic
4. Implement guide copying

**Deliverables**:
- Create worktrees for issues
- Automatic branch naming
- Implementation guide copying

---

### Phase 4: OpenCode Integration (Week 3)

**Tasks**:
1. Implement `.opencode-context` generation
2. Implement OpenCode launch integration
3. Add context prompt templates

**Deliverables**:
- OpenCode session launching
- Context file generation

---

### Phase 5: Cleanup & Status (Week 4)

**Tasks**:
1. Implement `issue-flow cleanup` command
2. Implement `issue-flow status` command
3. Implement `issue-flow list worktrees` command

**Deliverables**:
- Worktree cleanup with safety checks
- Status reporting
- Worktree listing

---

### Phase 6: Polish & Documentation (Week 4)

**Tasks**:
1. Add comprehensive error handling
2. Write user documentation
3. Create example templates
4. Add unit tests
5. Create demo video

**Deliverables**:
- Complete README
- Example templates for common use cases
- Test coverage
- Installation guide

---

## Example Usage Scenarios

### Scenario 1: Compliance Issue Workflow

```bash
# Create compliance issue
$ issue-flow create --type compliance

? Issue title: GDPR data export endpoint
? Priority: CRITICAL
? Description: Implement user data export per GDPR Art. 20
? Estimated time: 4-6 hours

✓ Created issue #127: GDPR data export endpoint
✓ Created worktree at ~/issue-worktrees/issue-127
✓ Created branch compliance/critical-data-export
✓ Copied guide to IMPLEMENTATION_GUIDE.md
? Launch OpenCode now? Yes

🚀 Starting OpenCode session...
```

---

### Scenario 2: Feature Development

```bash
# Start work on existing issue
$ issue-flow start 134

Issue #134: Update OAuth desktop clients
Type: Feature | Priority: HIGH

✓ Created worktree at ~/issue-worktrees/issue-134
✓ Created branch feature/134-oauth-desktop
✓ Copied guide oauth-guides/high-oauth-desktop-update.md
? Launch OpenCode now? No

📁 Worktree ready at: ~/issue-worktrees/issue-134
💡 To start coding: cd ~/issue-worktrees/issue-134
```

---

### Scenario 3: Cleanup After Completion

```bash
# After merging PR
$ issue-flow cleanup

Found 3 worktrees:
  ✓ #126 - Merged 2 days ago
  ✓ #127 - Merged 1 day ago
  ⚠ #134 - Still open

? Clean up merged worktrees? Yes
? Delete branches too? No

✓ Removed worktree for #126
✓ Removed worktree for #127
⏭ Skipped #134 (still open)

✨ Cleanup complete!
```

---

## Distribution

### NPM Package

```bash
# Install globally
npm install -g @whisper-notes/issue-flow

# Or use via npx
npx @whisper-notes/issue-flow create
```

### Standalone Binary

Build standalone binaries for macOS/Linux/Windows using `pkg`:

```bash
# Download binary
curl -L https://github.com/Whisper-Notes/issue-flow/releases/latest/download/issue-flow-macos -o issue-flow
chmod +x issue-flow
mv issue-flow /usr/local/bin/
```

---

## Success Metrics

1. **Reduces workflow time** by 50% (from manual git commands + issue creation)
2. **Eliminates mistakes** (wrong branch names, missing guides, etc.)
3. **Reusable** across multiple projects
4. **Easy to extend** with new issue types
5. **Well-documented** with examples

---

## Future Enhancements

### Phase 7+ (Future)

1. **PR Management**
   - `issue-flow pr create` - Create PR from worktree
   - `issue-flow pr sync` - Sync worktree with upstream

2. **Team Collaboration**
   - Shared templates repository
   - Template marketplace

3. **AI Integration**
   - Auto-generate issue descriptions from code
   - Suggest implementation steps

4. **IDE Integration**
   - VS Code extension
   - JetBrains plugin

5. **Analytics**
   - Track time spent per issue
   - Workflow efficiency metrics

---

## Technical Decisions

### Why TypeScript?
- Type safety prevents runtime errors
- Better IDE support
- Self-documenting code
- Easier refactoring

### Why Commander.js?
- Industry standard for Node.js CLIs
- Great documentation
- Supports subcommands well
- Active maintenance

### Why Handlebars for Templates?
- Simple syntax
- Widely used
- Supports helpers and partials
- Logic-less (separates presentation from logic)

### Why `gh` CLI over Octokit?
- Users already have `gh` installed
- Handles authentication automatically
- Simpler than managing API tokens
- Can fall back to Octokit if needed

---

## Dependencies

```json
{
  "dependencies": {
    "commander": "^11.0.0",
    "inquirer": "^9.0.0",
    "chalk": "^5.0.0",
    "handlebars": "^4.7.8",
    "cli-table3": "^0.6.3",
    "ora": "^7.0.0",
    "execa": "^8.0.0",
    "cosmiconfig": "^9.0.0",
    "zod": "^3.22.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0",
    "tsx": "^4.0.0",
    "vitest": "^1.0.0"
  }
}
```

---

## File Structure After Build

```
issue-flow/
├── dist/                      # Compiled JavaScript
│   ├── cli.js                # Main entry point
│   └── ...
├── src/                       # TypeScript source
├── templates/                 # Default templates
│   ├── compliance/
│   ├── feature/
│   └── bug/
├── bin/
│   └── issue-flow            # Executable script
├── package.json
├── tsconfig.json
└── README.md
```

---

## Next Steps

1. **Review this plan** and provide feedback
2. **Approve architecture** decisions
3. **Start Phase 1**: Set up TypeScript CLI project
4. **Create basic commands** structure
5. **Test with compliance workflow** use case

---

## Questions for Discussion

1. Should we build this as a standalone tool or integrate into existing codebase?
2. Do you want to support multiple Git providers (GitLab, Bitbucket)?
3. Should templates be in Markdown, JSON, or both?
4. Do you want built-in templates or user-provided only?
5. Should we support issue dependencies (e.g., #135 depends on #134)?

---

**Estimated Total Development Time**: 3-4 weeks (1 developer)

**Next Action**: Get your approval to proceed with Phase 1 implementation.
