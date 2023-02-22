# tm

An opinionated manager for `tmux` sessions. Use `tm` for managing groups of
`tmux` sessions.

A `tmux` session might contain a combination of multiple windows and/or multiple
splits (panes). With a project containing one or more sessions.

Some inspiration was taken from [smug](https://github.com/ivaaaan/smug), however
it lacked the ability start multiple sessions, which is part of my personal
workflow.

## Installation

You can install via [Homebrew](https://brew.sh) from the `tombell/formulae` tap.

    brew install tombell/formulae/tm

Alternatively, you can install the latest version if you have `go` installed:

    go install github.com/tombell/tm/cmd/tm@latest

## Configuration

Configuration files for projects are located in `~/.config/tm` (currently this
cannot be configured itself, it may be possible in the future). Each YAML
configuration file represents a "project". The name of the file is used as the
project name (for example `~/.config/tm/awesome-web-app.yml`, the project name
will be **awesome-web-app**.

Below is all the available fields for the YAML file.

```yaml
# the initial root for the tmux sessions, ~ will be resolved to the home directory
root: ~/Code
# a list of commands to run in the root before starting all tmux sessions
before_start:
    - echo "Hello world"
    - echo "Starting..."
# a list of commands to run in the root after stopping all tmux sessions
after_stop:
    - echo "Stopping..."
    - echo "Stopped"
# a list of tmux sessions to start
sessions:
    # the name of the sessions
    - name: frontend
      # an absolute path or relative to the top level root
      root: ./frontend
      # a list of windows to create in the tmux session
      windows:
        # the name of the window
        - name: server
          # an absolute path or relative to the session root
          root: ./server
          # a list of commands to run in the window
          commands:
            - echo "This is the server"
          # the layout for the panes/splits in the window: main-horizontal|main-vertical|even-horizontal|even-vertical|tiled
          layout: main-horizontal
          # a list of panes/splits to create in the window
          panes:
            # the type of pane/split: horizontal|vertical
            - type: horizontal
              # an absolute path or relative to the window root
              root: ./tests
              # a list of commands to run in the pane/split
              commands:
                - npm test
            # another pane/split in the window
            - type: horizontal
              root: ./stories
              commands:
                - npm run storybook
        # another window in the session
        - name: editor
          root: ./server
          commands:
            - nvim .
    # another session in the project
    - name: backend
      root: ./backend
      windows:
        - name: server
          root: ./server
          commands:
            - bundle exec rails s
        - name: editor
          root: ./server
          commands:
            - nvim .
```

You can create as many project files as you like for your different projects.
You can even start multiple projects with `tm`.

## Usage

Below is the documentation for each available command.

If you would like to see debug information logged out when running a command you
can specify the `-d`/`--debug` flag before the command argument when running
`tm` (for example: `tm --debug start my-project`).

### List Projects

The `list` command will list out all the available projects based on the project
configuration files in `~/.config/tm`.

    $ tm list
    memoir
    residently
    stark
    testing

### Start Project

The `start` command will start the `tmux` sessions for the given project name.

    $ tm start my-project

Currently you can only specify a single project name, but support for multiple
will be added.

### Stop Project

The `stop` command will stop the `tmux` sessions for the given project name.

    $ tm stop my-project

Currently you can only specify a single project name, but support for multiple
will be added.
