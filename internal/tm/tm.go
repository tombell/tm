package tm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tombell/tm/internal/cmd"
	"github.com/tombell/tm/internal/config"
	"github.com/tombell/tm/internal/tmux"
)

const (
	defaultWindowName = "default_win_name"
)

type Tm struct {
	tmux tmux.Tmux
	cmd  cmd.Cmd
}

func New(tmux tmux.Tmux, cmd cmd.Cmd) Tm {
	return Tm{tmux, cmd}
}

func (tm Tm) Start(cfg *config.Config, ctx Context) error {
	root := expandPath(cfg.Root)

	for _, s := range cfg.Sessions {
		if tm.tmux.SessionExists(s.Name) {
			continue
		}

		sessionRoot := resolvePath(root, s.Root)
		if _, err := tm.tmux.NewSession(s.Name, sessionRoot, defaultWindowName); err != nil {
			return err
		}

		for _, w := range s.Windows {
			windowRoot := resolvePath(sessionRoot, w.Root)
			if _, err := tm.tmux.NewWindow(s.Name, w.Name, windowRoot); err != nil {
				return err
			}
		}

		if _, err := tm.tmux.KillWindow(defaultWindowName); err != nil {
			return err
		}
	}

	return nil
}

func (tm Tm) Stop(cfg *config.Config) error {
	for _, s := range cfg.Sessions {
		if _, err := tm.tmux.KillSession(s.Name); err != nil {
			return err
		}
	}

	return nil
}

func expandPath(name string) string {
	if strings.HasPrefix(name, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return name
		}

		return strings.Replace(name, "~", homeDir, 1)
	}

	return name
}

func resolvePath(root, name string) string {
	baseRoot := expandPath(name)
	if baseRoot == "" || !filepath.IsAbs(baseRoot) {
		baseRoot = filepath.Join(root, name)
	}
	return baseRoot
}
