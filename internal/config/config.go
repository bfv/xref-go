package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bfv/xref/internal/models"
)

const configDirName = ".xrefcli"
const configFileName = "xrefconfig.json"

// Repo represents a configured xref repository.
type Repo struct {
	Name   string `json:"name"`
	Dir    string `json:"dir"`
	SrcDir string `json:"srcdir"`
}

// ConfigData holds the persisted configuration.
type ConfigData struct {
	Current string `json:"current"`
	Repos   []Repo `json:"repos"`
}

// Config manages application configuration and repo data.
type Config struct {
	Data       ConfigData
	configDir  string
	configFile string
}

// NewConfig creates and initializes a Config instance.
func NewConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}

	c := &Config{
		configDir:  filepath.Join(homeDir, configDirName),
		configFile: filepath.Join(homeDir, configDirName, configFileName),
	}

	if err := c.ensureConfigDir(); err != nil {
		return nil, err
	}

	if err := c.load(); err != nil {
		// First run — initialize empty config
		c.Data = ConfigData{Repos: []Repo{}}
	}

	return c, nil
}

func (c *Config) ensureConfigDir() error {
	return os.MkdirAll(c.configDir, 0755)
}

func (c *Config) load() error {
	data, err := os.ReadFile(c.configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.Data)
}

// Save persists the configuration to disk.
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal config: %w", err)
	}
	return os.WriteFile(c.configFile, data, 0644)
}

// GetRepo returns the repo with the given name, or an error if not found.
func (c *Config) GetRepo(name string) (*Repo, error) {
	for i := range c.Data.Repos {
		if c.Data.Repos[i].Name == name {
			return &c.Data.Repos[i], nil
		}
	}
	return nil, fmt.Errorf("repo '%s' not found", name)
}

// AddRepo adds a new repo to the configuration.
func (c *Config) AddRepo(name, dir, srcdir string) error {
	for _, r := range c.Data.Repos {
		if r.Name == name {
			return fmt.Errorf("repo '%s' already exists", name)
		}
	}
	c.Data.Repos = append(c.Data.Repos, Repo{
		Name:   name,
		Dir:    dir,
		SrcDir: srcdir,
	})
	return nil
}

// RemoveRepo removes a repo by name.
func (c *Config) RemoveRepo(name string) error {
	for i, r := range c.Data.Repos {
		if r.Name == name {
			c.Data.Repos = append(c.Data.Repos[:i], c.Data.Repos[i+1:]...)
			if c.Data.Current == name {
				c.Data.Current = ""
			}
			return nil
		}
	}
	return fmt.Errorf("repo '%s' not found", name)
}

// SetCurrent sets the current active repo.
func (c *Config) SetCurrent(name string) error {
	for _, r := range c.Data.Repos {
		if r.Name == name {
			c.Data.Current = name
			return nil
		}
	}
	return fmt.Errorf("repo '%s' not found", name)
}

// GetCurrentRepo returns the current active repo.
func (c *Config) GetCurrentRepo() (*Repo, error) {
	if c.Data.Current == "" {
		return nil, fmt.Errorf("no current repo set")
	}
	return c.GetRepo(c.Data.Current)
}

// RepoDataPath returns the path to the JSON data file for a repo.
func (c *Config) RepoDataPath(name string) string {
	return filepath.Join(c.configDir, name+".json")
}

// WriteRepoData writes parsed xref data to disk for a repo.
func (c *Config) WriteRepoData(name string, xreffiles []*models.XrefFile) error {
	data, err := json.MarshalIndent(xreffiles, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal repo data: %w", err)
	}
	return os.WriteFile(c.RepoDataPath(name), data, 0644)
}

// ReadRepoData reads parsed xref data from disk for a repo.
func (c *Config) ReadRepoData(name string) ([]*models.XrefFile, error) {
	data, err := os.ReadFile(c.RepoDataPath(name))
	if err != nil {
		return nil, fmt.Errorf("cannot read repo data for '%s': %w", name, err)
	}

	var xreffiles []*models.XrefFile
	if err := json.Unmarshal(data, &xreffiles); err != nil {
		return nil, fmt.Errorf("cannot parse repo data for '%s': %w", name, err)
	}
	return xreffiles, nil
}
