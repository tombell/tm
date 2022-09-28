package tm

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

type Tm struct {
	tmux tmux.Tmux
	cmd  cmd.Cmd
}

func New(tmux tmux.Tmux, cmd cmd.Cmd) Tm {
	return Tm{tmux, cmd}
}

func (tm Tm) Start(cfg *config.Config, ctx Context) error {
	root := expandPath(cfg.Root)

	if err := tm.execShellCommands(cfg.BeforeStart, root); err != nil {
		return err
	}

	for _, s := range cfg.Sessions {
		if tm.tmux.SessionExists(s.Name) {
			continue
		}

		sessionRoot := resolvePath(root, s.Root)

		if err := tm.execShellCommands(s.Commands, sessionRoot); err != nil {
			return err
		}

		if _, err := tm.tmux.NewSession(s.Name, sessionRoot, defaultWindowName); err != nil {
			return err
		}

		for _, w := range s.Windows {
			windowRoot := resolvePath(sessionRoot, w.Root)

			if _, err := tm.tmux.NewWindow(s.Name, w.Name, windowRoot); err != nil {
				return err
			}

			for _, c := range w.Commands {
				if err := tm.tmux.SendKeys(w.Name, c); err != nil {
					return err
				}
			}

			for i, p := range w.Panes {
				paneRoot := resolvePath(windowRoot, p.Root)

				pane, err := tm.tmux.SplitWindow(w.Name, p.Type, paneRoot)
				if err != nil {
					return err
				}

				if i%2 == 0 {
					if _, err := tm.tmux.SelectLayout(w.Name, tmux.Tiled); err != nil {
						return err
					}
				}

				for _, c := range p.Commands {
					if err := tm.tmux.SendKeys(w.Name+"."+pane, c); err != nil {
						return err
					}
				}
			}

			layout := w.Layout
			if layout == "" {
				layout = tmux.EvenVertical
			}

			if _, err := tm.tmux.SelectLayout(w.Name, layout); err != nil {
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

func (tm Tm) execShellCommands(commands []string, path string) error {
	for _, c := range commands {
		cmd := exec.Command("/bin/sh", "-c", c)
		cmd.Dir = path

		if _, err := tm.cmd.Exec(cmd); err != nil {
			return err
		}
	}

	return nil
}
