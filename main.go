package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
type model struct {
	todoList      *TodoList
	activePanel   Panel
	fileCursor    int
	todoCursor    int
	mode          Mode
	inputText     string
	editingIndex  int // -1 means adding new, >= 0 means editing existing
	width         int
	height        int
	statusMessage string
	files         []string
	todoDir       string
	currentFile   string
}

func initialModel() model {
	homeDir, _ := os.UserHomeDir()
	todoDir := filepath.Join(homeDir, ".tui_todos")

	// Create todo directory if it doesn't exist
	os.MkdirAll(todoDir, 0755)

	// Load list of todo files
	files := loadTodoFiles(todoDir)

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
		todoList:     todoList,
		activePanel:  FilePanel,
		fileCursor:   0,
		todoCursor:   0,
		mode:         NormalMode,
		editingIndex: -1,
		files:        files,
		todoDir:      todoDir,
		currentFile:  currentFile,
	}
}

func loadTodoFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			files = append(files, entry.Name())
		}
	}
	return files
}

func (m model) Init() tea.Cmd {
	return nil
}

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

func (m model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		// Go back to file panel from todo panel
		if m.activePanel == TodoPanel {
			m.activePanel = FilePanel
		}

	case "tab":
		// Switch between file and todo panel
		if m.activePanel == FilePanel {
			m.activePanel = TodoPanel
		} else {
			m.activePanel = FilePanel
		}

	case "j", "down":
		if m.activePanel == FilePanel {
			if m.fileCursor < len(m.files)-1 {
				m.fileCursor++
			}
		} else {
			if m.todoCursor < len(m.todoList.Todos)-1 {
				m.todoCursor++
			}
		}

	case "k", "up":
		if m.activePanel == FilePanel {
			if m.fileCursor > 0 {
				m.fileCursor--
			}
		} else {
			if m.todoCursor > 0 {
				m.todoCursor--
			}
		}

	case "enter":
		// Select file from file panel
		if m.activePanel == FilePanel && m.fileCursor < len(m.files) {
			m.currentFile = m.files[m.fileCursor]
			m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
			m.activePanel = TodoPanel
			m.todoCursor = 0
			m.statusMessage = fmt.Sprintf("Opened: %s", m.currentFile)
		}

	case "n":
		// Create new file (only in file panel)
		if m.activePanel == FilePanel {
			m.mode = EditMode
			m.editingIndex = -2 // Special value for new file
			m.inputText = ""
			m.statusMessage = "Enter filename (without .json)"
		}

	case "a":
		// Add new todo (only in todo panel)
		if m.activePanel == TodoPanel {
			m.mode = EditMode
			m.editingIndex = -1
			m.inputText = ""
			m.todoCursor = 0
			m.statusMessage = "Adding new todo (Enter to save, Esc to cancel)"
		}

	case "i":
		// Edit current todo (only in todo panel)
		if m.activePanel == TodoPanel && m.todoCursor < len(m.todoList.Todos) {
			m.mode = EditMode
			m.editingIndex = m.todoCursor
			m.inputText = m.todoList.Todos[m.todoCursor].Title
			m.statusMessage = "Editing todo (Enter to save, Esc to cancel)"
		}

	case "d":
		// Delete current todo (only in todo panel)
		if m.activePanel == TodoPanel && m.todoCursor < len(m.todoList.Todos) {
			m.todoList.Delete(m.todoCursor)
			if m.todoCursor >= len(m.todoList.Todos) && m.todoCursor > 0 {
				m.todoCursor--
			}
			m.statusMessage = "Deleted todo"
		}

	case "x", " ":
		// Toggle completion (only in todo panel)
		if m.activePanel == TodoPanel && m.todoCursor < len(m.todoList.Todos) {
			wasCompleted := m.todoList.Todos[m.todoCursor].Completed
			m.todoList.Toggle(m.todoCursor)

			if !wasCompleted {
				if m.todoCursor >= len(m.todoList.Todos) {
					m.todoCursor = len(m.todoList.Todos) - 1
				}
			} else {
				if m.todoCursor >= len(m.todoList.Todos) {
					m.todoCursor = len(m.todoList.Todos) - 1
				}
			}
			m.statusMessage = "Toggled todo status"
		}
	}

	return m, nil
}

func (m model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = NormalMode
		m.statusMessage = "Cancelled"
		return m, nil

	case "enter":
		if m.inputText != "" {
			if m.editingIndex == -2 {
				// Creating new file
				filename := m.inputText + ".json"
				newPath := filepath.Join(m.todoDir, filename)
				m.todoList = NewTodoList(newPath)
				m.todoList.Save() // Force save to create the file
				m.currentFile = filename
				m.files = loadTodoFiles(m.todoDir) // Reload file list after save

				// Find index of new file
				for i, f := range m.files {
					if f == filename {
						m.fileCursor = i
						break
					}
				}

				m.activePanel = TodoPanel
				m.todoCursor = 0
				m.statusMessage = fmt.Sprintf("Created: %s", filename)
			} else if m.editingIndex == -1 {
				// Adding new todo at top
				m.todoList.Insert(m.todoCursor, m.inputText)
				m.todoCursor = 0
			} else {
				// Editing existing todo
				m.todoList.Update(m.editingIndex, m.inputText)
			}
			m.mode = NormalMode
			if m.editingIndex >= 0 {
				m.statusMessage = "Saved"
			}
		} else {
			m.statusMessage = "Cannot be empty"
		}
		return m, nil

	case "backspace":
		if len(m.inputText) > 0 {
			m.inputText = m.inputText[:len(m.inputText)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.inputText += msg.String()
		}
	}

	return m, nil
}

func (m model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action != tea.MouseActionPress && msg.Action != tea.MouseActionRelease {
		return m, nil
	}

	x, y := msg.X, msg.Y

	// Calculate panel boundaries
	leftWidth := m.width / 4
	leftPanelEnd := leftWidth + 2

	// Click in left panel (files)
	if x >= 0 && x < leftPanelEnd {
		m.activePanel = FilePanel
		clickedLine := y - 3
		if clickedLine >= 0 && clickedLine < len(m.files) {
			m.fileCursor = clickedLine
		}
		return m, nil
	}

	// Click in right panel (todos)
	if x >= leftPanelEnd && x < m.width {
		m.activePanel = TodoPanel
		clickedLine := y - 3
		if clickedLine >= 0 && clickedLine < len(m.todoList.Todos) {
			m.todoCursor = clickedLine
			m.statusMessage = fmt.Sprintf("Selected: %s", m.todoList.Todos[clickedLine].Title)
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Define styles
	var (
		selectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7dcfff")).
				Background(lipgloss.Color("#283457"))

		borderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#565f89"))

		activeBorderStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#7aa2f7"))

		titleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#bb9af7")).
				Bold(true)

		completedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#565f89")).
				Strikethrough(true)

		hintStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9ece6a"))

		editStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f7768e"))
	)

	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - 4

	// Left Panel - Files
	leftContent := ""
	if m.mode == EditMode && m.editingIndex == -2 {
		leftContent = editStyle.Render("▶ " + m.inputText + "█.json") + "\n"
		for _, file := range m.files {
			leftContent += "  " + file + "\n"
		}
	} else {
		for i, file := range m.files {
			if m.activePanel == FilePanel && i == m.fileCursor {
				leftContent += selectedStyle.Render("▶ " + file) + "\n"
			} else if file == m.currentFile {
				leftContent += lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Render("• " + file) + "\n"
			} else {
				leftContent += "  " + file + "\n"
			}
		}
	}

	leftBorder := borderStyle
	if m.activePanel == FilePanel {
		leftBorder = activeBorderStyle
	}

	leftPanel := leftBorder.
		Width(leftWidth).
		Height(m.height - 4).
		Render(titleStyle.Render("Files") + "\n\n" + leftContent)

	// Right Panel - Todos
	rightContent := ""

	if len(m.todoList.Todos) == 0 {
		rightContent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Render("No todos yet. Press 'a' to add one.")
	} else {
		for i, todo := range m.todoList.Todos {
			checkbox := "[ ]"
			if todo.Completed {
				checkbox = "[✓]"
			}

			line := fmt.Sprintf("%s %s", checkbox, todo.Title)

			if todo.Completed {
				line = completedStyle.Render(line)
			}

			if m.mode == EditMode && m.editingIndex == i {
				line = editStyle.Render("✎ " + m.inputText + "█")
			} else if m.activePanel == TodoPanel && i == m.todoCursor {
				line = selectedStyle.Render("▶ " + line)
			} else {
				line = "  " + line
			}

			rightContent += line + "\n"
		}
	}

	// Show new todo input inline
	if m.mode == EditMode && m.editingIndex == -1 {
		lines := ""
		lines += editStyle.Render("▶ [ ] " + m.inputText + "█") + "\n"
		for i := 0; i < len(m.todoList.Todos); i++ {
			checkbox := "[ ]"
			if m.todoList.Todos[i].Completed {
				checkbox = "[✓]"
			}
			line := fmt.Sprintf("  %s %s", checkbox, m.todoList.Todos[i].Title)
			if m.todoList.Todos[i].Completed {
				line = completedStyle.Render(line)
			}
			lines += line + "\n"
		}
		rightContent = lines
	}

	rightBorder := borderStyle
	if m.activePanel == TodoPanel {
		rightBorder = activeBorderStyle
	}

	rightPanel := rightBorder.
		Width(rightWidth).
		Height(m.height - 4).
		Render(titleStyle.Render(fmt.Sprintf("Todos: %s", m.currentFile)) + "\n\n" + rightContent)

	// Combine panels
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Hints bar
	hints := ""
	if m.mode == EditMode {
		if m.editingIndex == -2 {
			hints = hintStyle.Render("Enter: create file | Esc: cancel")
		} else {
			hints = hintStyle.Render("Enter: save | Esc: cancel")
		}
	} else if m.activePanel == FilePanel {
		hints = hintStyle.Render("j/k: navigate | n: new file | Enter: open | Tab: switch panel | q: quit")
	} else {
		hints = hintStyle.Render("j/k: navigate | a: add | i: edit | d: delete | x/space: toggle | Tab: switch | q: quit")
	}

	statusBar := ""
	if m.statusMessage != "" {
		statusBar = "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ece6a")).
			Render(m.statusMessage)
	}

	return panels + "\n" + hints + statusBar
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
