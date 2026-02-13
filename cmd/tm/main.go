package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tombell/tm/internal/cmd"
	"github.com/tombell/tm/internal/config"
	"github.com/tombell/tm/internal/manager"
	"github.com/tombell/tm/internal/tmux"
)

const helpText = `usage: tm [<flags>] <command>

Commands:

  list          List all projects
  start         Start a tmux project
  stop          Stop a tmux project

Special options:

  -d/--debug    Show debug logging
  -v/--version  Show the version number, then exit
  --help        Show this message, then exit
`

var (
	debug   bool
	version bool
)

const projectsDir = "~/.config/tm"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, helpText)
		os.Exit(2)
	}

	flag.BoolVar(&debug, "debug", false, "")
	flag.BoolVar(&debug, "d", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&version, "v", false, "")

	flag.Parse()

	if version {
		fmt.Fprintf(os.Stdout, "tm %s (%s)\n", Version, Commit)
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
	}

	var logger *log.Logger
	if debug {
		logger = log.New(os.Stderr, "", 0)
	}

	switch args[0] {
	case "list":
		if err := list(); err != nil {
			log.Fatal(err)
		}
	case "start":
		if err := start(logger, args[1:]); err != nil {
			log.Fatal(err)
		}
	case "stop":
		if err := stop(logger, args[1:]); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintf(os.Stderr, "error: %q is not a known command\n", os.Args[1])
		flag.Usage()
	}
}

func usageText(text string) func() {
	return func() {
		fmt.Fprintln(os.Stderr, text)
		os.Exit(2)
	}
}

func list() error {
	files, err := os.ReadDir(manager.ExpandPath(projectsDir))
	if err != nil {
		return fmt.Errorf("failed to read projects directory: %w", err)
	}

	for _, file := range files {
		filename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		fmt.Println(filename)
	}

	return nil
}

func runCommand(logger *log.Logger, args []string, cmdName string, action func(manager.Manager, *config.Config) error) error {
	flagSet := flag.NewFlagSet(cmdName, flag.ExitOnError)
	flagSet.Usage = usageText(fmt.Sprintf("usage: tm %s <project names...>", cmdName))

	var help bool
	flagSet.BoolVar(&help, "help", false, "")
	flagSet.BoolVar(&help, "h", false, "")

	flagSet.Parse(args)

	if help {
		flagSet.Usage()
		os.Exit(0)
	}

	subArgs := flagSet.Args()

	if len(subArgs) < 1 {
		return fmt.Errorf("usage: tm %s <project names...>", cmdName)
	}

	var invalidProjects []string
	var configs []*config.Config

	for _, project := range subArgs {
		projectPath := fmt.Sprintf("%s/%s.yml", projectsDir, project)

		cfg, err := config.Load(projectPath)
		if err != nil {
			invalidProjects = append(invalidProjects, project)
			continue
		}

		configs = append(configs, cfg)
	}

	if len(invalidProjects) > 0 {
		return fmt.Errorf("the following projects do not exist or are invalid: %s", strings.Join(invalidProjects, ", "))
	}

	for i, cfg := range configs {
		project := subArgs[i]

		c := cmd.NewDefaultCmd(logger)
		t := tmux.New(c)
		m := manager.New(t, c)

		if err := action(m, cfg); err != nil {
			return fmt.Errorf("failed to %s project %s: %w", cmdName, project, err)
		}
	}

	return nil
}

func start(logger *log.Logger, args []string) error {
	return runCommand(logger, args, "start", func(m manager.Manager, cfg *config.Config) error {
		return m.Start(cfg, manager.CreateContext())
	})
}

func stop(logger *log.Logger, args []string) error {
	return runCommand(logger, args, "stop", func(m manager.Manager, cfg *config.Config) error {
		return m.Stop(cfg, manager.CreateContext())
	})
}
