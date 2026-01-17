package cmd

import (
	"fmt"
	"os"

	"github.com/paolorechia/issue-flow/internal/storage"
	"github.com/spf13/cobra"
)

var testDB *storage.Database

var rootCmd = &cobra.Command{
	Use:   "issue-flow",
	Short: "Multi-project workflow management tool",
	Long:  "A CLI tool for managing GitHub issues, git worktrees, and workflows across multiple projects.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getDB() (*storage.Database, error) {
	if testDB != nil {
		return testDB, nil
	}
	return storage.New()
}

func shouldCloseDB(db *storage.Database) bool {
	return db != testDB
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "issue-flow v0.1.0")
	},
}
