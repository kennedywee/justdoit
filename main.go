package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// initialModel creates and initializes the application model
func initialModel() model {
	homeDir, _ := os.UserHomeDir()
	todoDir := filepath.Join(homeDir, ".tui_todos")
	archiveDir := filepath.Join(homeDir, ".tui_todos", "archive")

	// Create directories if they don't exist
	os.MkdirAll(todoDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	// Load list of todo files
	files := loadTodoFiles(todoDir)
	archivedFiles := loadTodoFiles(archiveDir)

	var currentFile string
	var todoList *TodoList

	if len(files) > 0 {
		currentFile = files[0]
		todoList = NewTodoList(filepath.Join(todoDir, currentFile))
	} else {
		// Create default file if none exist
		currentFile = "default.json"
		todoList = NewTodoList(filepath.Join(todoDir, currentFile))
		files = []string{currentFile}
	}

	return model{
		todoList:       todoList,
		activePanel:    FilePanel,
		fileCursor:     0,
		todoCursor:     0,
		mode:           NormalMode,
		editingIndex:   -1,
		files:          files,
		archivedFiles:  archivedFiles,
		todoDir:        todoDir,
		archiveDir:     archiveDir,
		currentFile:    currentFile,
		showingArchive: false,
		styles:         NewStyles(),
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
