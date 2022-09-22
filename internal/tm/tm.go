package tm

import (
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

func (tm Tm) Start(cfg *config.Config) error {
	for _, s := range cfg.Sessions {
		if _, err := tm.tmux.NewSession(s.Name, s.Root, defaultWindowName); err != nil {
			return err
		}

		for _, w := range s.Windows {
			if _, err := tm.tmux.NewWindow(s.Name, w.Name, "."); err != nil {
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
