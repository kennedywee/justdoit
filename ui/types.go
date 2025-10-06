// Package ui implements the terminal user interface using Bubble Tea.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"justdoit/todo"
)

// Panel represents which panel is active
type Panel int

const (
	FilePanel Panel = iota
	TodoPanel
)

// Mode represents the current input mode
type Mode int

const (
	NormalMode Mode = iota
	EditMode
)

// Model holds the application state
type Model struct {
	TodoList       *todo.TodoList
	ActivePanel    Panel
	FileCursor     int
	TodoCursor     int
	Mode           Mode
	InputText      string
	EditingIndex   int // -1 means adding new, >= 0 means editing existing, -2 means new file, -3 means archive prompt, -4 means delete prompt
	Width          int
	Height         int
	StatusMessage  string
	Files          []string
	ArchivedFiles  []string
	TodoDir        string
	ArchiveDir     string
	CurrentFile    string
	ShowingArchive bool
	Styles         Styles
}

// Init initializes the model (Bubble Tea interface)
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model (Bubble Tea interface)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.MouseMsg:
		if m.Mode == NormalMode {
			return m.handleMouse(msg)
		}

	case tea.KeyMsg:
		if m.Mode == EditMode {
			return m.handleEditMode(msg)
		}
		return m.handleNormalMode(msg)
	}

	return m, nil
}
