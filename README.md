# TUI Todo

A terminal user interface (TUI) todo list manager built with Go and Bubble Tea.

## Features

- Multiple todo lists with file management
- Archive completed lists
- Keyboard-driven navigation
- Mouse support
- Clean, modern UI with dual-panel layout

## Build

```bash
go build -o tui_todo
```

## Run

```bash
./tui_todo
```

## Usage

### File Panel (Left)
- `j/k` or `↑/↓`: Navigate files
- `Enter` or `Space`: Open file
- `a`: Create new file
- `d`: Delete file
- `A` (Shift+A): Archive file
- `z`: Toggle archived files view
- `h/l` or `←/→`: Switch panels
- `Tab`: Switch panels

### Todo Panel (Right)
- `j/k` or `↑/↓`: Navigate todos
- `a`: Add new todo
- `i`: Edit todo
- `d`: Delete todo
- `x` or `Space`: Toggle completion
- `h/l` or `←/→`: Switch panels
- `Tab`: Switch panels

### General
- `q` or `Ctrl+C`: Quit
- `Esc`: Cancel operation or return to file panel

## Data Storage

Todo files are stored in `~/.tui_todos/`
Archived files are stored in `~/.tui_todos/archive/`
