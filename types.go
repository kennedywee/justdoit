package main

import tea "github.com/charmbracelet/bubbletea"

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

// model holds the application state
type model struct {
	todoList       *TodoList
	activePanel    Panel
	fileCursor     int
	todoCursor     int
	mode           Mode
	inputText      string
	editingIndex   int // -1 means adding new, >= 0 means editing existing, -2 means new file, -3 means archive prompt, -4 means delete prompt
	width          int
	height         int
	statusMessage  string
	files          []string
	archivedFiles  []string
	todoDir        string
	archiveDir     string
	currentFile    string
	showingArchive bool
	styles         Styles
}

// Init initializes the model (Bubble Tea interface)
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model (Bubble Tea interface)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.MouseMsg:
		if m.mode == NormalMode {
			return m.handleMouse(msg)
		}

	case tea.KeyMsg:
		if m.mode == EditMode {
			return m.handleEditMode(msg)
		}
		return m.handleNormalMode(msg)
	}

	return m, nil
}
