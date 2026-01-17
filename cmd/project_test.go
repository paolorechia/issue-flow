package cmd

import (
	"bytes"
	"testing"

	"github.com/paolorechia/issue-flow/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectAddCommand(t *testing.T) {
	db := testutil.NewTestDB(t)
	testutil.AssertDBEmpty(t, db)

	projectID := "test-project"
	projectName := "Test Project"
	githubOwner := "testowner"
	githubRepo := "testrepo"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"project", "add",
		"--id", projectID,
		"--name", projectName,
		"--owner", githubOwner,
		"--repo", githubRepo,
	})

	testDB = db
	t.Cleanup(func() { testDB = nil })

	err := rootCmd.Execute()

	require.NoError(t, err, "Command should succeed")

	testutil.AssertProjectCount(t, db, 1)

	project := testutil.AssertProjectExists(t, db, projectID)
	assert.Equal(t, projectID, project.ID)
	assert.Equal(t, projectName, project.Name)
	assert.Equal(t, githubOwner, project.GitHubOwner)
	assert.Equal(t, githubRepo, project.GitHubRepo)
}

func TestProjectAddCommandMissingRequiredFlags(t *testing.T) {
	t.Skip("Skipped: command calls os.Exit(1) which terminates test process")
}

func TestProjectListCommand_Empty(t *testing.T) {
	db := testutil.NewTestDB(t)
	testutil.AssertDBEmpty(t, db)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"project", "list"})

	testDB = db
	t.Cleanup(func() { testDB = nil })

	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestProjectListCommand_WithProjects(t *testing.T) {
	db := testutil.NewTestDB(t)

	projectID = "test-project-1"
	projectName = "Test Project 1"
	githubOwner = "testowner1"
	githubRepo = "testrepo1"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"project", "add",
		"--id", projectID,
		"--name", projectName,
		"--owner", githubOwner,
		"--repo", githubRepo,
	})

	testDB = db
	err := rootCmd.Execute()
	require.NoError(t, err)

	projectID = "test-project-2"
	projectName = "Test Project 2"
	githubOwner = "testowner2"
	githubRepo = "testrepo2"

	buf = new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"project", "add",
		"--id", projectID,
		"--name", projectName,
		"--owner", githubOwner,
		"--repo", githubRepo,
	})

	err = rootCmd.Execute()
	require.NoError(t, err)

	testutil.AssertProjectCount(t, db, 2)

	buf = new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"project", "list"})

	t.Cleanup(func() { testDB = nil })

	err = rootCmd.Execute()
	require.NoError(t, err)
}

func TestProjectShowCommand(t *testing.T) {
	db := testutil.NewTestDB(t)

	project := testutil.CreateTestProject(t, db)
	projectID := project.ID

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"project", "show", projectID})

	testDB = db
	t.Cleanup(func() { testDB = nil })

	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestProjectShowCommand_NotFound(t *testing.T) {
	t.Skip("Skipped: command calls os.Exit(1) which terminates test process")
}

func TestProjectWorkflow(t *testing.T) {
	db := testutil.NewTestDB(t)
	testutil.AssertDBEmpty(t, db)

	projectID := "workflow-test"
	projectName := "Workflow Test"
	githubOwner := "workflowowner"
	githubRepo := "workflowrepo"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"project", "add",
		"--id", projectID,
		"--name", projectName,
		"--owner", githubOwner,
		"--repo", githubRepo,
	})

	testDB = db
	t.Cleanup(func() { testDB = nil })

	err := rootCmd.Execute()
	require.NoError(t, err)

	testutil.AssertProjectCount(t, db, 1)

	project := testutil.AssertProjectExists(t, db, projectID)
	assert.Equal(t, projectID, project.ID)
	assert.Equal(t, projectName, project.Name)

	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"project", "show", projectID})
	err = rootCmd.Execute()
	require.NoError(t, err)

	worktree := testutil.CreateTestWorktree(t, db, projectID, 123)
	assert.NotNil(t, worktree)

	testutil.AssertWorktreeCount(t, db, 1)
	testutil.AssertWorktreeExists(t, db, worktree.ID)
}
