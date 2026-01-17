package testutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/paolorechia/issue-flow/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CLITest struct {
	T      *testing.T
	DB     *storage.Database
	Stdout strings.Builder
	Stderr strings.Builder
}

func NewTestDB(t *testing.T) *storage.Database {
	db, err := storage.NewWithDBPath(":memory:")
	require.NoError(t, err, "Failed to create in-memory database")

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func ExecCLI(t *testing.T, args ...string) *CLITest {
	return ExecCLIWithDB(t, nil, args...)
}

func ExecCLIWithDB(t *testing.T, db *storage.Database, args ...string) *CLITest {
	if db == nil {
		db = NewTestDB(t)
	}

	return &CLITest{
		T:  t,
		DB: db,
	}
}

func (ct *CLITest) AssertSuccess() *CLITest {
	require.Empty(ct.T, ct.Stderr.String(), "Expected no stderr output")
	return ct
}

func (ct *CLITest) AssertContains(substr string) *CLITest {
	assert.Contains(ct.T, ct.Stdout.String(), substr, "Output should contain substring")
	return ct
}

func (ct *CLITest) AssertNotContains(substr string) *CLITest {
	assert.NotContains(ct.T, ct.Stdout.String(), substr, "Output should not contain substring")
	return ct
}

func (ct *CLITest) AssertExitCode(code int) *CLITest {
	assert.Equal(ct.T, code, 0, "Expected exit code %d", code)
	return ct
}

func (ct *CLITest) Output() string {
	return ct.Stdout.String()
}

func (ct *CLITest) ErrOutput() string {
	return ct.Stderr.String()
}

func AssertProjectCount(t *testing.T, db *storage.Database, expected int) {
	projects, err := db.ListProjects()
	require.NoError(t, err)
	assert.Len(t, projects, expected, "Expected %d projects", expected)
}

func AssertProjectExists(t *testing.T, db *storage.Database, id string) *storage.Project {
	project, err := db.GetProject(id)
	require.NoError(t, err, "Expected project %s to exist", id)
	return project
}

func AssertProjectNotExists(t *testing.T, db *storage.Database, id string) {
	_, err := db.GetProject(id)
	assert.Error(t, err, "Expected project %s to not exist", id)
}

func AssertWorktreeCount(t *testing.T, db *storage.Database, expected int) {
	worktrees, err := db.ListWorktrees()
	require.NoError(t, err)
	assert.Len(t, worktrees, expected, "Expected %d worktrees", expected)
}

func AssertWorktreeExists(t *testing.T, db *storage.Database, id string) *storage.Worktree {
	worktree, err := db.GetWorktree(id)
	require.NoError(t, err, "Expected worktree %s to exist", id)
	return worktree
}

func AssertIssueCacheCount(t *testing.T, db *storage.Database, projectID string, expected int) {
	issues, err := db.ListIssueCache(projectID)
	require.NoError(t, err)
	assert.Len(t, issues, expected, "Expected %d cached issues for project %s", expected, projectID)
}

func AssertDBEmpty(t *testing.T, db *storage.Database) {
	projects, err := db.ListProjects()
	require.NoError(t, err)
	assert.Empty(t, projects, "Expected empty projects table")

	worktrees, err := db.ListWorktrees()
	require.NoError(t, err)
	assert.Empty(t, worktrees, "Expected empty worktrees table")
}

func CreateTestProject(t *testing.T, db *storage.Database) *storage.Project {
	project := &storage.Project{
		ID:          "test-project",
		Name:        "Test Project",
		GitHubOwner: "testowner",
		GitHubRepo:  "testrepo",
		LocalPath:   "/tmp/test-project",
		WorktreeDir: "/tmp/test-worktrees",
		Config:      `{"issue_types":[],"branch_config":{"pattern":"{prefix}/{issue-number}-{slug}","max_slug_length":50},"opencode":{"enabled":true,"auto_launch":false,"context_file":".opencode-context"}}`,
	}

	err := db.CreateProject(project)
	require.NoError(t, err, "Failed to create test project")
	return project
}

func CreateTestWorktree(t *testing.T, db *storage.Database, projectID string, issueNumber int) *storage.Worktree {
	worktree := &storage.Worktree{
		ID:          fmt.Sprintf("wt-%s-%d", projectID, issueNumber),
		ProjectID:   projectID,
		IssueNumber: issueNumber,
		Path:        fmt.Sprintf("/tmp/worktrees/%s/%d", projectID, issueNumber),
		Branch:      fmt.Sprintf("feature/%d-test-issue", issueNumber),
		Status:      "active",
	}

	err := db.CreateWorktree(worktree)
	require.NoError(t, err, "Failed to create test worktree")
	return worktree
}

func TableOutput() string {
	var buf strings.Builder
	return buf.String()
}

func ParseTableOutput(t *testing.T, output string) [][]string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return [][]string{}
	}

	var result [][]string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			result = append(result, fields)
		}
	}

	return result
}

func AssertTableRow(t *testing.T, table [][]string, row int, expected []string) {
	require.GreaterOrEqual(t, len(table), row+1, "Expected at least %d rows in table", row+1)
	require.Len(t, table[row], len(expected), "Expected row %d to have %d columns", row, len(expected))

	for i, expectedCell := range expected {
		assert.Equal(t, expectedCell, table[row][i], "Row %d column %d mismatch", row, i)
	}
}
