package project

import (
	"encoding/json"
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
	if p.ID == "" {
		return fmt.Errorf("project ID is required")
	}
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if p.GitHubOwner == "" {
		return fmt.Errorf("GitHub owner is required")
	}
	if p.GitHubRepo == "" {
		return fmt.Errorf("GitHub repo is required")
	}

	configJSON, err := json.Marshal(p.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	sp := &storage.Project{
		ID:          p.ID,
		Name:        p.Name,
		GitHubOwner: p.GitHubOwner,
		GitHubRepo:  p.GitHubRepo,
		LocalPath:   p.LocalPath,
		WorktreeDir: p.WorktreeDir,
		Config:      string(configJSON),
	}

	return m.db.CreateProject(sp)
}

func (m *Manager) Get(id string) (*Project, error) {
	sp, err := m.db.GetProject(id)
	if err != nil {
		return nil, err
	}

	var config ProjectConfig
	if err := json.Unmarshal([]byte(sp.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &Project{
		ID:          sp.ID,
		Name:        sp.Name,
		GitHubOwner: sp.GitHubOwner,
		GitHubRepo:  sp.GitHubRepo,
		LocalPath:   sp.LocalPath,
		WorktreeDir: sp.WorktreeDir,
		Config:      config,
		CreatedAt:   sp.CreatedAt,
		UpdatedAt:   sp.UpdatedAt,
	}, nil
}

func (m *Manager) List() ([]Project, error) {
	projects, err := m.db.ListProjects()
	if err != nil {
		return nil, err
	}

	result := make([]Project, len(projects))
	for i, sp := range projects {
		var config ProjectConfig
		if err := json.Unmarshal([]byte(sp.Config), &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		result[i] = Project{
			ID:          sp.ID,
			Name:        sp.Name,
			GitHubOwner: sp.GitHubOwner,
			GitHubRepo:  sp.GitHubRepo,
			LocalPath:   sp.LocalPath,
			WorktreeDir: sp.WorktreeDir,
			Config:      config,
			CreatedAt:   sp.CreatedAt,
			UpdatedAt:   sp.UpdatedAt,
		}
	}

	return result, nil
}

func (m *Manager) Delete(id string) error {
	return m.db.DeleteProject(id)
}

func (p *Project) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("project ID is required")
	}
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if p.GitHubOwner == "" {
		return fmt.Errorf("GitHub owner is required")
	}
	if p.GitHubRepo == "" {
		return fmt.Errorf("GitHub repo is required")
	}
	return nil
}

func (p *Project) GitHubFullName() string {
	return p.GitHubOwner + "/" + p.GitHubRepo
}
