package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tombell/tm/internal/cmd"
	"github.com/tombell/tm/internal/config"
	"github.com/tombell/tm/internal/manager"
	"github.com/tombell/tm/internal/tmux"
)

const helpText = `usage: tm [<flags>] <command>

Commands:

  start         Start a tmux project
  stop          Stop a tmux project

Special options:

  -d/--debug    SHow debug logging
  -v/--version  SHow the version number, then exit
  --help        Show this message, then exit
`

var (
	debug   bool
	version bool
)

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
		fmt.Fprintf(os.Stdout, "tm 0.0.0 (0)\n")
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
	case "start":
		start(logger, args[1:])
	case "stop":
		fmt.Println("stopping...")
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

func start(logger *log.Logger, args []string) {
	flagSet := flag.NewFlagSet("start", flag.ExitOnError)
	flagSet.Usage = usageText("usage: tm start <project name>")
	flagSet.Parse(args)

	subArgs := flagSet.Args()

	if len(subArgs) != 1 {
		flagSet.Usage()
	}

	projectPath := fmt.Sprintf("~/.config/tm/%s.yml", subArgs[0])

	cfg, err := config.Load(projectPath)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(3)
	}

	c := cmd.NewDefaultCmd(logger)
	t := tmux.New(c)
	m := manager.New(t, c)

	if err := m.Start(cfg, manager.CreateContext()); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(3)
	}
}
