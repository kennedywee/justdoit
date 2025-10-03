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
	todoList       *TodoList
	activePanel    Panel
	fileCursor     int
	todoCursor     int
	mode           Mode
	inputText      string
	editingIndex   int // -1 means adding new, >= 0 means editing existing, -2 means new file, -3 means archive prompt
	width          int
	height         int
	statusMessage  string
	files          []string
	archivedFiles  []string
	todoDir        string
	archiveDir     string
	currentFile    string
	showingArchive bool
}

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

func (m *model) archiveCurrentFile() {
	srcPath := filepath.Join(m.todoDir, m.currentFile)
	dstPath := filepath.Join(m.archiveDir, m.currentFile)

	os.Rename(srcPath, dstPath)

	// Reload file lists
	m.files = loadTodoFiles(m.todoDir)
	m.archivedFiles = loadTodoFiles(m.archiveDir)

	// Load next file or create default
	if len(m.files) > 0 {
		m.fileCursor = 0
		m.currentFile = m.files[0]
		m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
	} else {
		m.currentFile = "default.json"
		m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
		m.todoList.Save()
		m.files = loadTodoFiles(m.todoDir)
	}
	m.todoCursor = 0
}

func (m *model) unarchiveFile(filename string) {
	srcPath := filepath.Join(m.archiveDir, filename)
	dstPath := filepath.Join(m.todoDir, filename)

	os.Rename(srcPath, dstPath)

	// Reload file lists
	m.files = loadTodoFiles(m.todoDir)
	m.archivedFiles = loadTodoFiles(m.archiveDir)

	// Switch to the unarchived file
	m.currentFile = filename
	m.todoList = NewTodoList(dstPath)
	m.showingArchive = false

	// Find cursor position
	for i, f := range m.files {
		if f == filename {
			m.fileCursor = i
			break
		}
	}
}

func (m *model) allTodosCompleted() bool {
	if len(m.todoList.Todos) == 0 {
		return false
	}
	for _, todo := range m.todoList.Todos {
		if !todo.Completed {
			return false
		}
	}
	return true
}

func (m *model) previewFile() {
	var filename string
	var dir string

	if m.showingArchive {
		if m.fileCursor >= len(m.archivedFiles) {
			return
		}
		filename = m.archivedFiles[m.fileCursor]
		dir = m.archiveDir
	} else {
		if m.fileCursor >= len(m.files) {
			return
		}
		filename = m.files[m.fileCursor]
		dir = m.todoDir
	}

	// Load the file for preview (without switching activePanel)
	previewPath := filepath.Join(dir, filename)
	m.todoList = NewTodoList(previewPath)
	m.todoCursor = 0
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

	case "h", "left":
		// Go to left panel (file panel)
		m.activePanel = FilePanel
		// Preview current file when entering file panel
		m.previewFile()

	case "l", "right":
		// Go to right panel (todo panel)
		m.activePanel = TodoPanel

	case "tab":
		// Switch between file and todo panel
		if m.activePanel == FilePanel {
			m.activePanel = TodoPanel
		} else {
			m.activePanel = FilePanel
			// Preview current file when entering file panel
			m.previewFile()
		}

	case "j", "down":
		if m.activePanel == FilePanel {
			maxFiles := len(m.files)
			if m.showingArchive {
				maxFiles = len(m.archivedFiles)
			}
			if m.fileCursor < maxFiles-1 {
				m.fileCursor++
				// Preview file on cursor move
				m.previewFile()
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
				// Preview file on cursor move
				m.previewFile()
			}
		} else {
			if m.todoCursor > 0 {
				m.todoCursor--
			}
		}

	case "enter":
		// Select file from file panel
		if m.activePanel == FilePanel {
			if m.showingArchive && m.fileCursor < len(m.archivedFiles) {
				// Unarchive the selected file
				m.unarchiveFile(m.archivedFiles[m.fileCursor])
				m.activePanel = TodoPanel
				m.statusMessage = fmt.Sprintf("Unarchived: %s", m.currentFile)
			} else if !m.showingArchive && m.fileCursor < len(m.files) {
				// Open selected file
				m.currentFile = m.files[m.fileCursor]
				m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
				m.activePanel = TodoPanel
				m.todoCursor = 0
				m.statusMessage = fmt.Sprintf("Opened: %s", m.currentFile)
			}
		}

	case "z":
		// Toggle archive view (only in file panel)
		if m.activePanel == FilePanel {
			m.showingArchive = !m.showingArchive
			m.fileCursor = 0
			if m.showingArchive {
				m.statusMessage = "Showing archived files"
			} else {
				m.statusMessage = "Showing active files"
			}
		}

	case "n":
		// Create new file (only in file panel, not in archive view)
		if m.activePanel == FilePanel && !m.showingArchive {
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

	case " ":
		// Space key - toggle in todo panel, open file in file panel
		if m.activePanel == FilePanel {
			// Same as Enter in file panel
			if m.showingArchive && m.fileCursor < len(m.archivedFiles) {
				m.unarchiveFile(m.archivedFiles[m.fileCursor])
				m.activePanel = TodoPanel
				m.statusMessage = fmt.Sprintf("Unarchived: %s", m.currentFile)
			} else if !m.showingArchive && m.fileCursor < len(m.files) {
				m.currentFile = m.files[m.fileCursor]
				m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
				m.activePanel = TodoPanel
				m.todoCursor = 0
				m.statusMessage = fmt.Sprintf("Opened: %s", m.currentFile)
			}
		} else if m.activePanel == TodoPanel && m.todoCursor < len(m.todoList.Todos) {
			// Toggle completion in todo panel
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

			// Check if all todos are completed
			if m.allTodosCompleted() {
				m.mode = EditMode
				m.editingIndex = -3
				m.statusMessage = "All complete! Archive this list? (y/n)"
			} else {
				m.statusMessage = "Toggled todo status"
			}
		}

	case "x":
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

			// Check if all todos are completed
			if m.allTodosCompleted() {
				m.mode = EditMode
				m.editingIndex = -3 // Archive prompt
				m.statusMessage = "All complete! Archive this list? (y/n)"
			} else {
				m.statusMessage = "Toggled todo status"
			}
		}

	case "A":
		// Manual archive (shift+a, only in file panel)
		if m.activePanel == FilePanel && !m.showingArchive {
			m.mode = EditMode
			m.editingIndex = -3
			m.statusMessage = "Archive this file? (y/n)"
		}
	}

	return m, nil
}

func (m model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle archive prompt (y/n)
	if m.editingIndex == -3 {
		switch msg.String() {
		case "y", "Y":
			m.archiveCurrentFile()
			m.mode = NormalMode
			m.activePanel = FilePanel // Go back to file panel
			m.statusMessage = "File archived!"
			return m, nil
		case "n", "N", "esc":
			m.mode = NormalMode
			m.statusMessage = "Cancelled"
			return m, nil
		}
		return m, nil
	}

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

	// Left Panel - Files and Archives
	leftContent := ""
	if m.mode == EditMode && m.editingIndex == -2 {
		leftContent = editStyle.Render("▶ " + m.inputText + "█.json") + "\n"
		for _, file := range m.files {
			leftContent += "  " + file + "\n"
		}
	} else if m.showingArchive {
		// Show archived files
		leftContent += lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89")).Render("--- Archived ---") + "\n"
		for i, file := range m.archivedFiles {
			if m.activePanel == FilePanel && i == m.fileCursor {
				leftContent += selectedStyle.Render("▶ " + file) + "\n"
			} else {
				leftContent += "  " + file + "\n"
			}
		}
	} else {
		// Show active files
		for i, file := range m.files {
			if m.activePanel == FilePanel && i == m.fileCursor && !m.showingArchive {
				leftContent += selectedStyle.Render("▶ " + file) + "\n"
			} else if file == m.currentFile {
				leftContent += lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Render("• " + file) + "\n"
			} else {
				leftContent += "  " + file + "\n"
			}
		}

		// Show archive section
		if len(m.archivedFiles) > 0 {
			leftContent += "\n"
			leftContent += lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89")).Render("--- Archived ---") + "\n"
			leftContent += lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89")).Render(fmt.Sprintf("  (%d files)", len(m.archivedFiles))) + "\n"
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
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Archive confirmation - clean full screen
	if m.mode == EditMode && m.editingIndex == -3 {
		confirmStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#f7768e")).
			Padding(2, 4).
			Align(lipgloss.Center)

		title := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f7768e")).
			Bold(true).
			Render("Archive Confirmation")

		filename := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7aa2f7")).
			Render(fmt.Sprintf("'%s'", m.currentFile))

		question := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c0caf5")).
			Render("Archive this file?")

		options := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ece6a")).
			Render("[y] Yes, archive    [n] No, cancel")

		confirmContent := lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			"",
			filename,
			question,
			"",
			options,
		)

		confirmBox := confirmStyle.Render(confirmContent)

		panels := lipgloss.Place(
			m.width,
			m.height-4,
			lipgloss.Center,
			lipgloss.Center,
			confirmBox,
		)

		return panels + "\n" + hintStyle.Render("y: archive | n: cancel")
	}

	panels := mainView

	// Hints bar
	hints := ""
	if m.mode == EditMode {
		if m.editingIndex == -2 {
			hints = hintStyle.Render("Enter: create file | Esc: cancel")
		} else if m.editingIndex == -3 {
			hints = hintStyle.Render("y: archive | n: cancel")
		} else {
			hints = hintStyle.Render("Enter: save | Esc: cancel")
		}
	} else if m.activePanel == FilePanel {
		if m.showingArchive {
			hints = hintStyle.Render("j/k: navigate | Enter: unarchive | z: show active | h/l: switch panel | q: quit")
		} else {
			hints = hintStyle.Render("j/k: navigate | n: new | Enter: open | A: archive | z: show archived | h/l: switch | q: quit")
		}
	} else {
		hints = hintStyle.Render("j/k: navigate | a: add | i: edit | d: delete | x/space: toggle | h/l: switch panel | q: quit")
	}

	statusBar := ""
	if m.statusMessage != "" && m.editingIndex != -3 {
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
