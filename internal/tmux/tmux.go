package tmux

import (
	"os/exec"
	"strings"

	"github.com/tombell/tm/internal/cmd"
)

type TmuxSession struct {
	ID   string
	Name string
	Root string
}

type TmuxWindow struct {
	ID     string
	Name   string
	Root   string
	Layout string
}

type Tmux struct {
	cmd cmd.Cmd
}

func New(cmd cmd.Cmd) Tmux {
	return Tmux{cmd}
}

func (t Tmux) NewSession(name, root, windowName string) (string, error) {
	cmd := exec.Command("tmux", "new-session", "-Pd", "-s", name, "-n", windowName, "-c", root)
	return t.cmd.Exec(cmd)
}

func (t Tmux) KillSession(target string) (string, error) {
	cmd := exec.Command("tmux", "kill-session", "-t", target)
	return t.cmd.Exec(cmd)
}

func (t Tmux) ListSessions() ([]TmuxSession, error) {
	var sessions []TmuxSession

	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_id};#{session_name};#{session_path}")
	out, err := t.cmd.Exec(cmd)
	if err != nil {
		return sessions, err
	}

	sessionList := strings.Split(out, "\n")

	for _, s := range sessionList {
		sessionInfo := strings.Split(s, ";")
		session := TmuxSession{
			ID:   sessionInfo[0],
			Name: sessionInfo[1],
			Root: sessionInfo[2],
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (t Tmux) NewWindow(target, name, root string) (string, error) {
	cmd := exec.Command("tmux", "new-window", "-Pd", "-t", target, "-n", name, "-c", root, "-F", "#{window_id}")
	return t.cmd.Exec(cmd)
}

func (t Tmux) KillWindow(target string) (string, error) {
	cmd := exec.Command("tmux", "kill-window", "-t", target)
	return t.cmd.Exec(cmd)
}

func (t Tmux) ListWindows(target string) ([]TmuxWindow, error) {
	var windows []TmuxWindow

	cmd := exec.Command("tmux", "list-windows", "-t", target, "-F", "#{window_id};#{window_name};#{pane_current_path};#{window_layout}")
	out, err := t.cmd.Exec(cmd)
	if err != nil {
		return windows, err
	}

	windowList := strings.Split(out, "\n")

	for _, w := range windowList {
		windowInfo := strings.Split(w, ";")
		window := TmuxWindow{
			ID:     windowInfo[0],
			Name:   windowInfo[1],
			Root:   windowInfo[2],
			Layout: windowInfo[3],
		}

		windows = append(windows, window)
	}

	return windows, nil
}
