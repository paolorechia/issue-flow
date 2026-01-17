package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

type Project struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	GitHubOwner string    `db:"github_owner"`
	GitHubRepo  string    `db:"github_repo"`
	LocalPath   string    `db:"local_path"`
	WorktreeDir string    `db:"worktree_dir"`
	Config      string    `db:"config"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Worktree struct {
	ID          string    `db:"id"`
	ProjectID   string    `db:"project_id"`
	IssueNumber int       `db:"issue_number"`
	Path        string    `db:"path"`
	Branch      string    `db:"branch"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}

type IssueCache struct {
	ID          int       `db:"id"`
	ProjectID   string    `db:"project_id"`
	IssueNumber int       `db:"issue_number"`
	Title       string    `db:"title"`
	Type        string    `db:"type"`
	Priority    string    `db:"priority"`
	Status      string    `db:"status"`
	CachedAt    time.Time `db:"cached_at"`
}

func New() (*Database, error) {
	return NewWithDBPath("")
}

func NewWithDBPath(dbPath string) (*Database, error) {
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}

		dbDir := filepath.Join(home, ".issue-flow")
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}

		dbPath = filepath.Join(dbDir, "database.db")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	d := &Database{db: db}
	if err := d.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
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

	CREATE INDEX IF NOT EXISTS idx_worktrees_project_id ON worktrees(project_id);
	CREATE INDEX IF NOT EXISTS idx_worktrees_issue_number ON worktrees(issue_number);

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

	CREATE INDEX IF NOT EXISTS idx_issue_cache_project ON issue_cache(project_id);
	`

	_, err := d.db.Exec(schema)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateProject(p *Project) error {
	query := `
	INSERT INTO projects (id, name, github_owner, github_repo, local_path, worktree_dir, config)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, p.ID, p.Name, p.GitHubOwner, p.GitHubRepo, p.LocalPath, p.WorktreeDir, p.Config)
	return err
}

func (d *Database) GetProject(id string) (*Project, error) {
	query := `SELECT id, name, github_owner, github_repo, local_path, worktree_dir, config, created_at, updated_at FROM projects WHERE id = ?`

	row := d.db.QueryRow(query, id)
	var p Project
	err := row.Scan(&p.ID, &p.Name, &p.GitHubOwner, &p.GitHubRepo, &p.LocalPath, &p.WorktreeDir, &p.Config, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (d *Database) ListProjects() ([]Project, error) {
	query := `SELECT id, name, github_owner, github_repo, local_path, worktree_dir, config, created_at, updated_at FROM projects ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.GitHubOwner, &p.GitHubRepo, &p.LocalPath, &p.WorktreeDir, &p.Config, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (d *Database) DeleteProject(id string) error {
	query := `DELETE FROM projects WHERE id = ?`
	_, err := d.db.Exec(query, id)
	return err
}

func (d *Database) CreateWorktree(w *Worktree) error {
	query := `
	INSERT INTO worktrees (id, project_id, issue_number, path, branch, status)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, w.ID, w.ProjectID, w.IssueNumber, w.Path, w.Branch, w.Status)
	return err
}

func (d *Database) GetWorktree(id string) (*Worktree, error) {
	query := `SELECT id, project_id, issue_number, path, branch, status, created_at FROM worktrees WHERE id = ?`

	row := d.db.QueryRow(query, id)
	var w Worktree
	err := row.Scan(&w.ID, &w.ProjectID, &w.IssueNumber, &w.Path, &w.Branch, &w.Status, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (d *Database) ListWorktrees() ([]Worktree, error) {
	query := `SELECT id, project_id, issue_number, path, branch, status, created_at FROM worktrees ORDER BY created_at`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var worktrees []Worktree
	for rows.Next() {
		var w Worktree
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.IssueNumber, &w.Path, &w.Branch, &w.Status, &w.CreatedAt); err != nil {
			return nil, err
		}
		worktrees = append(worktrees, w)
	}

	return worktrees, nil
}

func (d *Database) ListWorktreesByProject(projectID string) ([]Worktree, error) {
	query := `SELECT id, project_id, issue_number, path, branch, status, created_at FROM worktrees WHERE project_id = ? ORDER BY created_at`

	rows, err := d.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var worktrees []Worktree
	for rows.Next() {
		var w Worktree
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.IssueNumber, &w.Path, &w.Branch, &w.Status, &w.CreatedAt); err != nil {
			return nil, err
		}
		worktrees = append(worktrees, w)
	}

	return worktrees, nil
}

func (d *Database) DeleteWorktree(id string) error {
	query := `DELETE FROM worktrees WHERE id = ?`
	_, err := d.db.Exec(query, id)
	return err
}

func (d *Database) CacheIssue(c *IssueCache) error {
	query := `
	INSERT OR REPLACE INTO issue_cache (project_id, issue_number, title, type, priority, status, cached_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, c.ProjectID, c.IssueNumber, c.Title, c.Type, c.Priority, c.Status, c.CachedAt)
	return err
}

func (d *Database) GetIssueCache(projectID string, issueNumber int) (*IssueCache, error) {
	query := `SELECT id, project_id, issue_number, title, type, priority, status, cached_at FROM issue_cache WHERE project_id = ? AND issue_number = ?`

	row := d.db.QueryRow(query, projectID, issueNumber)
	var c IssueCache
	err := row.Scan(&c.ID, &c.ProjectID, &c.IssueNumber, &c.Title, &c.Type, &c.Priority, &c.Status, &c.CachedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Database) ListIssueCache(projectID string) ([]IssueCache, error) {
	query := `SELECT id, project_id, issue_number, title, type, priority, status, cached_at FROM issue_cache WHERE project_id = ? ORDER BY issue_number`

	rows, err := d.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []IssueCache
	for rows.Next() {
		var c IssueCache
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.IssueNumber, &c.Title, &c.Type, &c.Priority, &c.Status, &c.CachedAt); err != nil {
			return nil, err
		}
		issues = append(issues, c)
	}

	return issues, nil
}

func (d *Database) ClearIssueCache(projectID string) error {
	query := `DELETE FROM issue_cache WHERE project_id = ?`
	_, err := d.db.Exec(query, projectID)
	return err
}
