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

type Tmux struct {
	cmd.Cmd
}

func (t Tmux) NewSession(name, root, windowName string) (string, error) {
	cmd := exec.Command("tmux", "new-session", "-Pd", "-s", name, "-n", windowName, "-c", root)
	return t.Exec(cmd)
}

func (t Tmux) KillSession(target string) (string, error) {
	cmd := exec.Command("tmux", "kill-session", "-t", target)
	return t.Exec(cmd)
}

func (t Tmux) ListSessions() ([]TmuxSession, error) {
	var sessions []TmuxSession

	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_id};#{session_name};#{session_path}")
	out, err := t.Exec(cmd)
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
