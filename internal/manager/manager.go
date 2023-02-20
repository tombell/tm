package manager

import (
	"os/exec"

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

	return m.createSessions(cfg.Sessions, root)
}

func (m Manager) Stop(cfg *config.Config, ctx Context) error {
	root := expandPath(cfg.Root)

	return m.killSessions(cfg.Sessions, root)

}

func (m Manager) createSessions(sessions []config.Session, root string) error {
	for _, s := range sessions {
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

		if err := m.createWindows(s.Windows, s.Name, sessionRoot); err != nil {
			return err
		}

		if _, err := m.tmux.KillWindow(defaultWindowName); err != nil {
			return err
		}

		if err := m.tmux.RenumberWindows(s.Name); err != nil {
			return err
		}
	}

	return nil
}

func (m Manager) killSessions(sessions []config.Session, root string) error {
	for _, s := range sessions {
		if !m.tmux.SessionExists(s.Name) {
			continue
		}

		if _, err := m.tmux.KillSession(s.Name); err != nil {
			return err
		}
	}

	return nil
}

func (m Manager) createWindows(windows []config.Window, sessionName, sessionRoot string) error {
	for _, w := range windows {
		windowRoot := resolvePath(sessionRoot, w.Root)

		if _, err := m.tmux.NewWindow(sessionName, w.Name, windowRoot); err != nil {
			return err
		}

		if err := m.sendCommands(w.Commands, w.Name); err != nil {
			return err
		}

		if err := m.createPanes(w.Panes, w.Name, windowRoot); err != nil {
			return err
		}

		layout := w.Layout
		if layout == "" {
			layout = tmux.EvenVertical
		}

		if _, err := m.tmux.SelectLayout(w.Name, layout); err != nil {
			return err
		}
	}

	return nil
}

func (m Manager) createPanes(panes []config.Pane, windowName, windowRoot string) error {
	for i, p := range panes {
		paneRoot := resolvePath(windowRoot, p.Root)

		pane, err := m.tmux.SplitWindow(windowName, p.Type, paneRoot)
		if err != nil {
			return err
		}

		if i%2 == 0 {
			if _, err := m.tmux.SelectLayout(windowName, tmux.Tiled); err != nil {
				return err
			}
		}

		if err := m.sendCommands(p.Commands, windowName+"."+pane); err != nil {
			return err
		}
	}

	return nil
}

func (m Manager) sendCommands(commands []string, target string) error {
	for _, command := range commands {
		if err := m.tmux.SendKeys(target, command); err != nil {
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
