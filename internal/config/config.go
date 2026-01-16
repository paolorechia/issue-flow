package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var cfgFile string
var v *viper.Viper

func init() {
	v = viper.New()
}

type Config struct {
	Version  string       `mapstructure:"version"`
	Settings Settings     `mapstructure:"settings"`
	GitHub   GitHubConfig `mapstructure:"github"`
	Projects []ProjectRef `mapstructure:"projects"`
}

type Settings struct {
	Editor          string `mapstructure:"editor"`
	OpenCodeEnabled bool   `mapstructure:"opencode_enabled"`
	WorktreeBase    string `mapstructure:"worktree_base"`
	Verbose         bool   `mapstructure:"verbose"`
}

type GitHubConfig struct {
	AuthMethod string `mapstructure:"auth_method"`
	Token      string `mapstructure:"token"`
}

type ProjectRef struct {
	ID          string `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	GitHubOwner string `mapstructure:"github_owner"`
	GitHubRepo  string `mapstructure:"github_repo"`
	LocalPath   string `mapstructure:"local_path"`
	WorktreeDir string `mapstructure:"worktree_dir"`
}

func Load() (*Config, error) {
	v.SetDefault("settings.editor", "code")
	v.SetDefault("settings.opencode_enabled", true)
	v.SetDefault("settings.worktree_base", filepath.Join(homeDir(), "issue-worktrees"))
	v.SetDefault("github.auth_method", "gh_cli")

	v.SetEnvPrefix("ISSUE_FLOW")
	v.AutomaticEnv()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home := homeDir()
		configDir := filepath.Join(home, ".issue-flow")

		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create config directory: %w", err)
			}
		}

		v.AddConfigPath(configDir)
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func SetConfigFile(file string) {
	cfgFile = file
}

func GetConfigFile() string {
	if cfgFile != "" {
		return cfgFile
	}
	home := homeDir()
	return filepath.Join(home, ".issue-flow", "config.yaml")
}

func GetConfigPath() string {
	home := homeDir()
	return filepath.Join(home, ".issue-flow")
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}
