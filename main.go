package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	cursor        int
	mode          Mode
	inputText     string
	editingIndex  int // -1 means adding new, >= 0 means editing existing
	width         int
	height        int
	statusMessage string
}

func initialModel() model {
	homeDir, _ := os.UserHomeDir()
	todoPath := filepath.Join(homeDir, ".tui_todo.json")

	return model{
		todoList:     NewTodoList(todoPath),
		cursor:       0,
		mode:         NormalMode,
		editingIndex: -1,
	}
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

	case "j", "down":
		if m.cursor < len(m.todoList.Todos)-1 {
			m.cursor++
		}

	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}

	case "a":
		// Add new todo below current cursor
		m.mode = EditMode
		m.editingIndex = -1
		m.inputText = ""
		m.statusMessage = "Adding new todo (Enter to save, Esc to cancel)"

	case "i":
		// Edit current todo
		if m.cursor < len(m.todoList.Todos) {
			m.mode = EditMode
			m.editingIndex = m.cursor
			m.inputText = m.todoList.Todos[m.cursor].Title
			m.statusMessage = "Editing todo (Enter to save, Esc to cancel)"
		}

	case "d":
		// Delete current todo
		if m.cursor < len(m.todoList.Todos) {
			m.todoList.Delete(m.cursor)
			if m.cursor >= len(m.todoList.Todos) && m.cursor > 0 {
				m.cursor--
			}
			m.statusMessage = "Deleted todo"
		}

	case "x", " ":
		// Toggle completion
		if m.cursor < len(m.todoList.Todos) {
			m.todoList.Toggle(m.cursor)
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
			if m.editingIndex == -1 {
				// Adding new todo
				m.todoList.Insert(m.cursor, m.inputText)
				m.cursor++ // Move cursor to the new item
			} else {
				// Editing existing todo
				m.todoList.Update(m.editingIndex, m.inputText)
			}
			m.mode = NormalMode
			m.statusMessage = "Saved"
		} else {
			m.statusMessage = "Todo cannot be empty"
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

	// Click anywhere in the list to select
	// Border top (0), Title (1), blank line (2), todos start at (3)
	clickedLine := y - 3

	if clickedLine >= 0 && clickedLine < len(m.todoList.Todos) {
		m.cursor = clickedLine
		m.statusMessage = fmt.Sprintf("Selected: %s", m.todoList.Todos[clickedLine].Title)
	}

	// Suppress unused variable warning
	_ = x

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

	// Todo list content
	content := ""

	if len(m.todoList.Todos) == 0 {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Render("No todos yet. Press 'a' to add one.")
	} else {
		for i, todo := range m.todoList.Todos {
			checkbox := "[ ]"
			if todo.Completed {
				checkbox = "[✓]"
			}

			line := fmt.Sprintf("%s %s", checkbox, todo.Title)

			// Apply completed styling
			if todo.Completed {
				line = completedStyle.Render(line)
			}

			// Show edit mode inline
			if m.mode == EditMode && m.editingIndex == i {
				line = editStyle.Render("✎ " + m.inputText + "█")
			} else if i == m.cursor {
				line = selectedStyle.Render("▶ " + line)
			} else {
				line = "  " + line
			}

			content += line + "\n"
		}
	}

	// Show new todo input inline at cursor position
	if m.mode == EditMode && m.editingIndex == -1 {
		lines := ""
		for i := 0; i <= m.cursor && i < len(m.todoList.Todos); i++ {
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
		lines += editStyle.Render("▶ [ ] " + m.inputText + "█") + "\n"
		for i := m.cursor + 1; i < len(m.todoList.Todos); i++ {
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
		content = lines
	}

	panel := borderStyle.
		Width(m.width - 4).
		Height(m.height - 4).
		Render(titleStyle.Render("Todo List") + "\n\n" + content)

	// Hints bar
	hints := ""
	if m.mode == EditMode {
		hints = hintStyle.Render("Enter: save | Esc: cancel")
	} else {
		hints = hintStyle.Render("j/k: navigate | a: add | i: edit | d: delete | x/space: toggle | q: quit")
	}

	statusBar := ""
	if m.statusMessage != "" {
		statusBar = "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ece6a")).
			Render(m.statusMessage)
	}

	return panel + "\n" + hints + statusBar
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
