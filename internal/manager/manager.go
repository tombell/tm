package manager

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tombell/tm/internal/cmd"
	"github.com/tombell/tm/internal/config"
	"github.com/tombell/tm/internal/tmux"
)

const (
	defaultWindowName = "default_win_name"
)

type Manager struct {
	tmux tmux.Tmux
	cmd  cmd.Cmd
}

func New(tmux tmux.Tmux, cmd cmd.Cmd) Manager {
	return Manager{tmux, cmd}
}

func (m Manager) Start(cfg *config.Config, ctx Context) error {
	root := expandPath(cfg.Root)

	if err := m.execShellCommands(cfg.BeforeStart, root); err != nil {
		return err
	}

	for _, s := range cfg.Sessions {
		if m.tmux.SessionExists(s.Name) {
			continue
		}

		sessionRoot := resolvePath(root, s.Root)

		if err := m.execShellCommands(s.Commands, sessionRoot); err != nil {
			return err
		}

		if _, err := m.tmux.NewSession(s.Name, sessionRoot, defaultWindowName); err != nil {
			return err
		}

		for _, w := range s.Windows {
			windowRoot := resolvePath(sessionRoot, w.Root)

			if _, err := m.tmux.NewWindow(s.Name, w.Name, windowRoot); err != nil {
				return err
			}

			for _, c := range w.Commands {
				if err := m.tmux.SendKeys(w.Name, c); err != nil {
					return err
				}
			}

			for i, p := range w.Panes {
				paneRoot := resolvePath(windowRoot, p.Root)

				pane, err := m.tmux.SplitWindow(w.Name, p.Type, paneRoot)
				if err != nil {
					return err
				}

				if i%2 == 0 {
					if _, err := m.tmux.SelectLayout(w.Name, tmux.Tiled); err != nil {
						return err
					}
				}

				for _, c := range p.Commands {
					if err := m.tmux.SendKeys(w.Name+"."+pane, c); err != nil {
						return err
					}
				}
			}

			layout := w.Layout
			if layout == "" {
				layout = tmux.EvenVertical
			}

			if _, err := m.tmux.SelectLayout(w.Name, layout); err != nil {
				return err
			}
		}

		if _, err := m.tmux.KillWindow(defaultWindowName); err != nil {
			return err
		}
	}

	return nil
}

func (m Manager) Stop(cfg *config.Config) error {
	for _, s := range cfg.Sessions {
		if _, err := m.tmux.KillSession(s.Name); err != nil {
			return err
		}
	}

	return nil
}

func (m Manager) execShellCommands(commands []string, path string) error {
	for _, c := range commands {
		cmd := exec.Command("/bin/sh", "-c", c)
		cmd.Dir = path

		if _, err := m.cmd.Exec(cmd); err != nil {
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