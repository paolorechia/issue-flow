package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/paolorechia/issue-flow/internal/project"
	"github.com/spf13/cobra"
)

var (
	projectID     string
	projectName   string
	githubOwner   string
	githubRepo    string
	localPath     string
	worktreeDir   string
	verboseOutput bool
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  "Add, list, and manage projects tracked by issue-flow.",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := getDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
			os.Exit(1)
		}
		if shouldCloseDB(db) {
			defer db.Close()
		}

		manager := project.NewManager(db)
		projects, err := manager.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing projects: %v\n", err)
			os.Exit(1)
		}

		if len(projects) == 0 {
			fmt.Println("No projects found. Use 'issue-flow project add' to add a project.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tREPOSITORY")
		for _, p := range projects {
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.ID, p.Name, p.GitHubFullName())
		}
		w.Flush()
	},
}

var projectAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new project",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := getDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
			os.Exit(1)
		}
		if shouldCloseDB(db) {
			defer db.Close()
		}

		if projectID == "" || projectName == "" || githubOwner == "" || githubRepo == "" {
			fmt.Fprintln(os.Stderr, "Error: --id, --name, --owner, and --repo are required")
			os.Exit(1)
		}

		p := &project.Project{
			ID:          projectID,
			Name:        projectName,
			GitHubOwner: githubOwner,
			GitHubRepo:  githubRepo,
			LocalPath:   localPath,
			WorktreeDir: worktreeDir,
			Config: project.ProjectConfig{
				BranchConfig: project.BranchConfig{
					Pattern:       "{prefix}/{issue-number}-{slug}",
					MaxSlugLength: 50,
				},
				OpenCode: project.OpenCodeConfig{
					Enabled:     true,
					AutoLaunch:  false,
					ContextFile: ".opencode-context",
				},
			},
		}

		manager := project.NewManager(db)
		if err := manager.Add(p); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Added project: %s (%s)\n", p.ID, p.GitHubFullName())
	},
}

var projectShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show project details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		db, err := getDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
			os.Exit(1)
		}
		if shouldCloseDB(db) {
			defer db.Close()
		}

		manager := project.NewManager(db)
		p, err := manager.Get(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Project: %s\n", p.ID)
		fmt.Printf("  Name: %s\n", p.Name)
		fmt.Printf("  Repository: %s\n", p.GitHubFullName())
		fmt.Printf("  Local Path: %s\n", p.LocalPath)
		fmt.Printf("  Worktree Dir: %s\n", p.WorktreeDir)
		fmt.Printf("  Created: %s\n", p.CreatedAt.Format("2006-01-02"))
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectAddCmd)
	projectCmd.AddCommand(projectShowCmd)

	projectAddCmd.Flags().StringVarP(&projectID, "id", "i", "", "Project ID (required)")
	projectAddCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name (required)")
	projectAddCmd.Flags().StringVarP(&githubOwner, "owner", "o", "", "GitHub owner (required)")
	projectAddCmd.Flags().StringVarP(&githubRepo, "repo", "r", "", "GitHub repo (required)")
	projectAddCmd.Flags().StringVarP(&localPath, "path", "p", "", "Local path (optional)")
	projectAddCmd.Flags().StringVar(&worktreeDir, "worktree-dir", "", "Worktree directory (optional)")
}
