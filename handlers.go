package main

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// handleNormalMode handles keyboard input in normal mode
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

	case "a":
		if m.activePanel == FilePanel {
			// Create new file (only in file panel, not in archive view)
			if !m.showingArchive {
				m.mode = EditMode
				m.editingIndex = -2 // Special value for new file
				m.inputText = ""
				m.statusMessage = "Enter filename (without .json)"
			}
		} else if m.activePanel == TodoPanel {
			// Add new todo (only in todo panel)
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
		if m.activePanel == FilePanel {
			// Delete file (only in file panel, not in archive view)
			if !m.showingArchive && m.fileCursor < len(m.files) {
				m.mode = EditMode
				m.editingIndex = -4 // Special value for delete file confirmation
				m.statusMessage = "Delete this file? (y/n)"
			}
		} else if m.activePanel == TodoPanel && m.todoCursor < len(m.todoList.Todos) {
			// Delete current todo (only in todo panel)
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
			m.toggleTodoWithArchivePrompt()
		}

	case "x":
		// Toggle completion (only in todo panel)
		if m.activePanel == TodoPanel && m.todoCursor < len(m.todoList.Todos) {
			m.toggleTodoWithArchivePrompt()
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

// toggleTodoWithArchivePrompt toggles a todo and prompts for archiving if all are complete
func (m *model) toggleTodoWithArchivePrompt() {
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

// handleEditMode handles keyboard input in edit mode
func (m model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle delete file prompt (y/n)
	if m.editingIndex == -4 {
		switch msg.String() {
		case "y", "Y":
			m.deleteCurrentFile()
			m.mode = NormalMode
			m.activePanel = FilePanel
			m.statusMessage = "File deleted!"
			return m, nil
		case "n", "N", "esc":
			m.mode = NormalMode
			m.statusMessage = "Cancelled"
			return m, nil
		}
		return m, nil
	}

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

// handleMouse handles mouse input
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
