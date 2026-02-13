package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.yml")

	yamlContent := `
root: ~/test
sessions:
  - name: test-session
    root: ./session
    windows:
      - name: editor
        commands:
          - vim .
`

	err := os.WriteFile(tempFile, []byte(yamlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "valid absolute path",
			filePath: tempFile,
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filePath: "/non/existent/file.yml",
			wantErr:  true,
		},
		{
			name: "invalid YAML",
			filePath: func() string {
				invalidFile := filepath.Join(tempDir, "invalid.yml")
				os.WriteFile(invalidFile, []byte("invalid: yaml: content: ["), 0o644)
				return invalidFile
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Load(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cfg == nil {
				t.Errorf("Load() returned nil config")
			}
		})
	}
}

func TestLoadWithTildeExpansion(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	homeFile := filepath.Join(homeDir, "test_config.yml")
	yamlContent := `
root: /tmp
sessions:
  - name: home-session
`
	defer os.Remove(homeFile)

	err = os.WriteFile(homeFile, []byte(yamlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to write home test file: %v", err)
	}

	expandedPath := strings.Replace("~/test_config.yml", "~", homeDir, 1)
	cfg, err := Load(expandedPath)
	if err != nil {
		t.Errorf("Load() with tilde expansion failed: %v", err)
	}
	if cfg == nil {
		t.Errorf("Load() returned nil config with tilde")
	}
}
