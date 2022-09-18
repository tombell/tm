package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Cmd interface {
	Exec(cmd *exec.Cmd) (string, error)
	ExecSilent(cmd *exec.Cmd) error
}

type DefaultCmd struct {
	logger *log.Logger
}

func (c DefaultCmd) Exec(cmd *exec.Cmd) (string, error) {
	if c.logger != nil {
		c.logger.Println(strings.Join(cmd.Args, ""))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if c.logger != nil {
			c.logger.Println(err, string(output))
		}

		return "", fmt.Errorf("executing command failed: %w", err)
	}

	return strings.TrimSuffix(string(output), "\n"), nil
}

func (c DefaultCmd) ExecSilent(cmd *exec.Cmd) error {
	if c.logger != nil {
		c.logger.Println(strings.Join(cmd.Args, ""))
	}

	if err := cmd.Run(); err != nil {
		if c.logger != nil {
			c.logger.Println(err)
		}

		return fmt.Errorf("executing command failed: %w", err)
	}

	return nil
}
