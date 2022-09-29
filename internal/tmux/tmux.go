package tmux

import (
	"os/exec"
	"strings"

	"github.com/tombell/tm/internal/cmd"
)

const (
	VerticalSplit   = "vertical"
	HorizontalSplit = "horizontal"
)

const (
	MainHorizontal = "main-horizontal"
	MainVertical   = "main-vertical"
	EvenHorizontal = "even-horizontal"
	EvenVertical   = "even-vertical"
	Tiled          = "tiled"
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

func (tmux Tmux) NewSession(name, root, windowName string) (string, error) {
	cmd := exec.Command("tmux", "new-session", "-Pd", "-s", name, "-n", windowName, "-c", root)
	return tmux.cmd.Exec(cmd)
}

func (tmux Tmux) KillSession(target string) (string, error) {
	cmd := exec.Command("tmux", "kill-session", "-t", target)
	return tmux.cmd.Exec(cmd)
}

func (tmux Tmux) ListSessions() ([]TmuxSession, error) {
	var sessions []TmuxSession

	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_id};#{session_name};#{session_path}")
	out, err := tmux.cmd.Exec(cmd)
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

func (tmux Tmux) SessionExists(target string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", target)
	out, err := tmux.cmd.Exec(cmd)
	return err == nil && out == ""
}

func (tmux Tmux) NewWindow(target, name, root string) (string, error) {
	cmd := exec.Command("tmux", "new-window", "-Pd", "-t", target, "-n", name, "-c", root, "-F", "#{window_id}")
	return tmux.cmd.Exec(cmd)
}

func (tmux Tmux) KillWindow(target string) (string, error) {
	cmd := exec.Command("tmux", "kill-window", "-t", target)
	return tmux.cmd.Exec(cmd)
}

func (tmux Tmux) ListWindows(target string) ([]TmuxWindow, error) {
	var windows []TmuxWindow

	cmd := exec.Command("tmux", "list-windows", "-t", target, "-F", "#{window_id};#{window_name};#{pane_current_path};#{window_layout}")
	out, err := tmux.cmd.Exec(cmd)
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

func (tmux Tmux) SplitWindow(target, splitType, root string) (string, error) {
	args := []string{"split-window", "-Pd"}

	switch splitType {
	case VerticalSplit:
		args = append(args, "-v")
	case HorizontalSplit:
		args = append(args, "-h")
	}

	args = append(args, []string{"-t", target, "-c", root, "-F", "#{pane_id}"}...)
	cmd := exec.Command("tmux", args...)

	return tmux.cmd.Exec(cmd)
}

func (tmux Tmux) RenumberWindows(target string) error {
	cmd := exec.Command("tmux", "move-window", "-r", "-s", target, "-t", target)
	_, err := tmux.cmd.Exec(cmd)
	return err
}

func (tmux Tmux) SelectLayout(target, layoutType string) (string, error) {
	cmd := exec.Command("tmux", "select-layout", "-t", target, layoutType)
	return tmux.cmd.Exec(cmd)
}

func (tmux Tmux) SendKeys(target, command string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", target, command, "Enter")
	return tmux.cmd.ExecSilent(cmd)
}
