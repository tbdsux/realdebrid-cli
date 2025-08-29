# realdebrid-cli

> [!NOTE]
> On-going development.

Manage your real-debrid account via a cli.

```sh
go install github.com/tbdsux/realdebrid-cli/rd@latest
```

## Usage

```sh
❯ rd

Real-Debrid CLI

Manage your Real-Debrid files and account.

Usage:
  rd [command]

Available Commands:
  account     Show account information
  completion  Generate the autocompletion script for the specified shell
  config      Manage configuration settings
  help        Help about any command
  magnet      Upload a magnet link
  torrent     Upload a torrent file

Flags:
      --config string   config file (default is $HOME/.realdebrid-cli.yaml)
  -d, --debug           Enable debug mode
  -h, --help            help for rd

Use "rd [command] --help" for more information about a command.

```

## Development

CLI uses [Bubble Tea](https://github.com/charmbracelet/bubbletea/) for power for the commands.

> This is not a TUI app though.
