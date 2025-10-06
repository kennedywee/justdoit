package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"justdoit/todo"
	"justdoit/ui"
)

// initialModel creates and initializes the application model
func initialModel() ui.Model {
	homeDir, _ := os.UserHomeDir()
	todoDir := filepath.Join(homeDir, ".tui_todos")
	archiveDir := filepath.Join(homeDir, ".tui_todos", "archive")

	// Create directories if they don't exist
	os.MkdirAll(todoDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	// Load list of todo files
	files := ui.LoadTodoFiles(todoDir)
	archivedFiles := ui.LoadTodoFiles(archiveDir)

	var currentFile string
	var todoList *todo.TodoList

	if len(files) > 0 {
		currentFile = files[0]
		todoList = todo.NewTodoList(filepath.Join(todoDir, currentFile))
	} else {
		// Create default file if none exist
		currentFile = "default.json"
		todoList = todo.NewTodoList(filepath.Join(todoDir, currentFile))
		files = []string{currentFile}
	}

	return ui.Model{
		TodoList:       todoList,
		ActivePanel:    ui.FilePanel,
		FileCursor:     0,
		TodoCursor:     0,
		Mode:           ui.NormalMode,
		EditingIndex:   -1,
		Files:          files,
		ArchivedFiles:  archivedFiles,
		TodoDir:        todoDir,
		ArchiveDir:     archiveDir,
		CurrentFile:    currentFile,
		ShowingArchive: false,
		Styles:         ui.NewStyles(),
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
