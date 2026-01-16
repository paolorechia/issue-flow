package project

import "time"

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

type ProjectConfig struct {
	IssueTypes   []IssueType    `json:"issue_types" yaml:"issue_types"`
	BranchConfig BranchConfig   `json:"branch_config" yaml:"branch_config"`
	OpenCode     OpenCodeConfig `json:"opencode" yaml:"opencode"`
}

type IssueType struct {
	Name         string   `json:"name" yaml:"name"`
	Label        string   `json:"label" yaml:"label"`
	Priority     []string `json:"priority" yaml:"priority"`
	BranchPrefix string   `json:"branch_prefix" yaml:"branch_prefix"`
	Template     string   `json:"template" yaml:"template"`
	GuidesDir    string   `json:"guides_dir" yaml:"guides_dir"`
}

type BranchConfig struct {
	Pattern       string `json:"pattern" yaml:"pattern"`
	MaxSlugLength int    `json:"max_slug_length" yaml:"max_slug_length"`
}

type OpenCodeConfig struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`
	AutoLaunch      bool   `json:"auto_launch" yaml:"auto_launch"`
	ContextFile     string `json:"context_file" yaml:"context_file"`
	ContextTemplate string `json:"context_template" yaml:"context_template"`
}
