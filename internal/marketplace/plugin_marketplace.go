package marketplace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type PluginEntry struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Type        string            `json:"type"`
	Tags        []string          `json:"tags"`
	Downloads   int               `json:"downloads"`
	Installed   bool              `json:"installed,omitempty"`
	Config      map[string]any    `json:"config,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type PluginMarketplace struct {
	cacheDir    string
	installDir  string
}

func NewPluginMarketplace(cacheDir, installDir string) *PluginMarketplace {
	return &PluginMarketplace{
		cacheDir:   cacheDir,
		installDir: installDir,
	}
}

func (m *PluginMarketplace) Publish(entry PluginEntry) error {
	entries, err := m.loadPlugins()
	if err != nil {
		entries = []PluginEntry{}
	}

	entry.UpdatedAt = time.Now()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}

	for i, e := range entries {
		if e.Name == entry.Name {
			entries[i] = entry
			return m.savePlugins(entries)
		}
	}

	entries = append(entries, entry)
	return m.savePlugins(entries)
}

func (m *PluginMarketplace) Get(name string) (*PluginEntry, error) {
	entries, err := m.loadPlugins()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Name == name {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("plugin %s not found", name)
}

func (m *PluginMarketplace) Search(query string, tags []string) ([]PluginEntry, error) {
	entries, err := m.loadPlugins()
	if err != nil {
		return nil, err
	}

	var results []PluginEntry
	for _, entry := range entries {
		if query != "" {
			matched := false
			if containsStr(entry.Name, query) || containsStr(entry.Description, query) {
				matched = true
			}
			for _, tag := range entry.Tags {
				if containsStr(tag, query) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		if len(tags) > 0 {
			hasTag := false
			for _, filterTag := range tags {
				for _, entryTag := range entry.Tags {
					if entryTag == filterTag {
						hasTag = true
						break
					}
				}
				if hasTag {
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		results = append(results, entry)
	}

	return results, nil
}

func (m *PluginMarketplace) List() ([]PluginEntry, error) {
	return m.loadPlugins()
}

func (m *PluginMarketplace) Install(name string) error {
	entry, err := m.Get(name)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(m.installDir, 0o755); err != nil {
		return err
	}

	pluginDir := filepath.Join(m.installDir, name)
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		return err
	}

	configFile := filepath.Join(pluginDir, "plugin.json")
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, data, 0o600); err != nil {
		return err
	}

	entries, err := m.loadPlugins()
	if err != nil {
		return err
	}

	for i, e := range entries {
		if e.Name == name {
			entries[i].Installed = true
			entries[i].Downloads++
			entries[i].UpdatedAt = time.Now()
			return m.savePlugins(entries)
		}
	}

	return nil
}

func (m *PluginMarketplace) Uninstall(name string) error {
	pluginDir := filepath.Join(m.installDir, name)
	if err := os.RemoveAll(pluginDir); err != nil {
		return err
	}

	entries, err := m.loadPlugins()
	if err != nil {
		return err
	}

	for i, entry := range entries {
		if entry.Name == name {
			entries[i].Installed = false
			entries[i].UpdatedAt = time.Now()
			return m.savePlugins(entries)
		}
	}

	return nil
}

func (m *PluginMarketplace) IsInstalled(name string) bool {
	pluginDir := filepath.Join(m.installDir, name)
	_, err := os.Stat(pluginDir)
	return err == nil
}

func (m *PluginMarketplace) ListInstalled() ([]PluginEntry, error) {
	entries, err := m.loadPlugins()
	if err != nil {
		return nil, err
	}

	var installed []PluginEntry
	for _, entry := range entries {
		if m.IsInstalled(entry.Name) {
			entry.Installed = true
			installed = append(installed, entry)
		}
	}

	return installed, nil
}

func (m *PluginMarketplace) loadPlugins() ([]PluginEntry, error) {
	path := filepath.Join(m.cacheDir, "plugins.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return m.defaultPlugins(), nil
		}
		return nil, err
	}

	var entries []PluginEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func (m *PluginMarketplace) savePlugins(entries []PluginEntry) error {
	if err := os.MkdirAll(m.cacheDir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(m.cacheDir, "plugins.json"), data, 0o600)
}

func (m *PluginMarketplace) defaultPlugins() []PluginEntry {
	return []PluginEntry{
		{
			Name:        "naeos-lint",
			Version:     "1.0.0",
			Description: "Advanced linting rules for NAEOS specs",
			Author:      "naeos",
			Type:        "lint",
			Tags:        []string{"lint", "validation", "quality"},
			Downloads:   500,
		},
		{
			Name:        "naeos-security",
			Version:     "1.0.0",
			Description: "Security audit plugin for specifications",
			Author:      "naeos",
			Type:        "security",
			Tags:        []string{"security", "audit", "compliance"},
			Downloads:   350,
		},
		{
			Name:        "naeos-docs",
			Version:     "1.0.0",
			Description: "Auto-generate documentation from specs",
			Author:      "naeos",
			Type:        "documentation",
			Tags:        []string{"docs", "documentation", "generation"},
			Downloads:   280,
		},
		{
			Name:        "naeos-test",
			Version:     "1.0.0",
			Description: "Test generation and execution plugin",
			Author:      "naeos",
			Type:        "testing",
			Tags:        []string{"test", "testing", "coverage"},
			Downloads:   220,
		},
	}
}
