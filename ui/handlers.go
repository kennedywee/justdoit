package ui

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"justdoit/todo"
)

// handleNormalMode handles keyboard input in normal mode
func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		// Go back to file panel from todo panel
		if m.ActivePanel == TodoPanel {
			m.ActivePanel = FilePanel
		}

	case "h", "left":
		// Go to left panel (file panel)
		m.ActivePanel = FilePanel
		// Preview current file when entering file panel
		m.previewFile()

	case "l", "right":
		// Go to right panel (todo panel)
		m.ActivePanel = TodoPanel

	case "tab":
		// Switch between file and todo panel
		if m.ActivePanel == FilePanel {
			m.ActivePanel = TodoPanel
		} else {
			m.ActivePanel = FilePanel
			// Preview current file when entering file panel
			m.previewFile()
		}

	case "j", "down":
		if m.ActivePanel == FilePanel {
			maxFiles := len(m.Files)
			if m.ShowingArchive {
				maxFiles = len(m.ArchivedFiles)
			}
			if m.FileCursor < maxFiles-1 {
				m.FileCursor++
				// Preview file on cursor move
				m.previewFile()
			}
		} else {
			if m.TodoCursor < len(m.TodoList.Todos)-1 {
				m.TodoCursor++
			}
		}

	case "k", "up":
		if m.ActivePanel == FilePanel {
			if m.FileCursor > 0 {
				m.FileCursor--
				// Preview file on cursor move
				m.previewFile()
			}
		} else {
			if m.TodoCursor > 0 {
				m.TodoCursor--
			}
		}

	case "enter":
		// Select file from file panel
		if m.ActivePanel == FilePanel {
			if m.ShowingArchive && m.FileCursor < len(m.ArchivedFiles) {
				// Unarchive the selected file
				m.unarchiveFile(m.ArchivedFiles[m.FileCursor])
				m.ActivePanel = TodoPanel
				m.StatusMessage = fmt.Sprintf("Unarchived: %s", m.CurrentFile)
			} else if !m.ShowingArchive && m.FileCursor < len(m.Files) {
				// Open selected file
				m.CurrentFile = m.Files[m.FileCursor]
				m.TodoList = todo.NewTodoList(filepath.Join(m.TodoDir, m.CurrentFile))
				m.ActivePanel = TodoPanel
				m.TodoCursor = 0
				m.StatusMessage = fmt.Sprintf("Opened: %s", m.CurrentFile)
			}
		}

	case "z":
		// Toggle archive view (only in file panel)
		if m.ActivePanel == FilePanel {
			m.ShowingArchive = !m.ShowingArchive
			m.FileCursor = 0
			if m.ShowingArchive {
				m.StatusMessage = "Showing archived files"
			} else {
				m.StatusMessage = "Showing active files"
			}
		}

	case "a":
		switch m.ActivePanel {
		case FilePanel:
			// Create new file (only in file panel, not in archive view)
			if !m.ShowingArchive {
				m.Mode = EditMode
				m.EditingIndex = -2 // Special value for new file
				m.InputText = ""
				m.StatusMessage = "Enter filename (without .json)"
			}
		case TodoPanel:
			// Add new todo (only in todo panel)
			m.Mode = EditMode
			m.EditingIndex = -1
			m.InputText = ""
			m.TodoCursor = 0
			m.StatusMessage = "Adding new todo (Enter to save, Esc to cancel)"
		}

	case "i":
		// Edit current todo (only in todo panel)
		if m.ActivePanel == TodoPanel && m.TodoCursor < len(m.TodoList.Todos) {
			m.Mode = EditMode
			m.EditingIndex = m.TodoCursor
			m.InputText = m.TodoList.Todos[m.TodoCursor].Title
			m.StatusMessage = "Editing todo (Enter to save, Esc to cancel)"
		}

	case "d":
		if m.ActivePanel == FilePanel {
			// Delete file (only in file panel, not in archive view)
			if !m.ShowingArchive && m.FileCursor < len(m.Files) {
				m.Mode = EditMode
				m.EditingIndex = -4 // Special value for delete file confirmation
				m.StatusMessage = "Delete this file? (y/n)"
			}
		} else if m.ActivePanel == TodoPanel && m.TodoCursor < len(m.TodoList.Todos) {
			// Delete current todo (only in todo panel)
			m.TodoList.Delete(m.TodoCursor)
			if m.TodoCursor >= len(m.TodoList.Todos) && m.TodoCursor > 0 {
				m.TodoCursor--
			}
			m.StatusMessage = "Deleted todo"
		}

	case " ":
		// Space key - toggle in todo panel, open file in file panel
		if m.ActivePanel == FilePanel {
			// Same as Enter in file panel
			if m.ShowingArchive && m.FileCursor < len(m.ArchivedFiles) {
				m.unarchiveFile(m.ArchivedFiles[m.FileCursor])
				m.ActivePanel = TodoPanel
				m.StatusMessage = fmt.Sprintf("Unarchived: %s", m.CurrentFile)
			} else if !m.ShowingArchive && m.FileCursor < len(m.Files) {
				m.CurrentFile = m.Files[m.FileCursor]
				m.TodoList = todo.NewTodoList(filepath.Join(m.TodoDir, m.CurrentFile))
				m.ActivePanel = TodoPanel
				m.TodoCursor = 0
				m.StatusMessage = fmt.Sprintf("Opened: %s", m.CurrentFile)
			}
		} else if m.ActivePanel == TodoPanel && m.TodoCursor < len(m.TodoList.Todos) {
			// Toggle completion in todo panel
			m.toggleTodoWithArchivePrompt()
		}

	case "x":
		// Toggle completion (only in todo panel)
		if m.ActivePanel == TodoPanel && m.TodoCursor < len(m.TodoList.Todos) {
			m.toggleTodoWithArchivePrompt()
		}

	case "A":
		// Manual archive (shift+a, only in file panel)
		if m.ActivePanel == FilePanel && !m.ShowingArchive {
			m.Mode = EditMode
			m.EditingIndex = -3
			m.StatusMessage = "Archive this file? (y/n)"
		}
	}

	return m, nil
}

// toggleTodoWithArchivePrompt toggles a todo and prompts for archiving if all are complete
func (m *Model) toggleTodoWithArchivePrompt() {
	wasCompleted := m.TodoList.Todos[m.TodoCursor].Completed
	m.TodoList.Toggle(m.TodoCursor)

	if !wasCompleted {
		if m.TodoCursor >= len(m.TodoList.Todos) {
			m.TodoCursor = len(m.TodoList.Todos) - 1
		}
	} else {
		if m.TodoCursor >= len(m.TodoList.Todos) {
			m.TodoCursor = len(m.TodoList.Todos) - 1
		}
	}

	// Check if all todos are completed
	if m.allTodosCompleted() {
		m.Mode = EditMode
		m.EditingIndex = -3
		m.StatusMessage = "All complete! Archive this list? (y/n)"
	} else {
		m.StatusMessage = "Toggled todo status"
	}
}

// handleEditMode handles keyboard input in edit mode
func (m Model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle delete file prompt (y/n)
	if m.EditingIndex == -4 {
		switch msg.String() {
		case "y", "Y":
			m.deleteCurrentFile()
			m.Mode = NormalMode
			m.ActivePanel = FilePanel
			m.StatusMessage = "File deleted!"
			return m, nil
		case "n", "N", "esc":
			m.Mode = NormalMode
			m.StatusMessage = "Cancelled"
			return m, nil
		}
		return m, nil
	}

	// Handle archive prompt (y/n)
	if m.EditingIndex == -3 {
		switch msg.String() {
		case "y", "Y":
			m.archiveCurrentFile()
			m.Mode = NormalMode
			m.ActivePanel = FilePanel // Go back to file panel
			m.StatusMessage = "File archived!"
			return m, nil
		case "n", "N", "esc":
			m.Mode = NormalMode
			m.StatusMessage = "Cancelled"
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "esc":
		m.Mode = NormalMode
		m.StatusMessage = "Cancelled"
		return m, nil

	case "enter":
		if m.InputText != "" {
			if m.EditingIndex == -2 {
				// Creating new file
				filename := m.InputText + ".json"
				newPath := filepath.Join(m.TodoDir, filename)
				m.TodoList = todo.NewTodoList(newPath)
				m.TodoList.Save() // Force save to create the file
				m.CurrentFile = filename
				m.Files = LoadTodoFiles(m.TodoDir) // Reload file list after save

				// Find index of new file
				for i, f := range m.Files {
					if f == filename {
						m.FileCursor = i
						break
					}
				}

				m.ActivePanel = TodoPanel
				m.TodoCursor = 0
				m.StatusMessage = fmt.Sprintf("Created: %s", filename)
			} else if m.EditingIndex == -1 {
				// Adding new todo at top
				m.TodoList.Insert(m.TodoCursor, m.InputText)
				m.TodoCursor = 0
			} else {
				// Editing existing todo
				m.TodoList.Update(m.EditingIndex, m.InputText)
			}
			m.Mode = NormalMode
			if m.EditingIndex >= 0 {
				m.StatusMessage = "Saved"
			}
		} else {
			m.StatusMessage = "Cannot be empty"
		}
		return m, nil

	case "backspace":
		if len(m.InputText) > 0 {
			m.InputText = m.InputText[:len(m.InputText)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.InputText += msg.String()
		}
	}

	return m, nil
}

// handleMouse handles mouse input
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action != tea.MouseActionPress && msg.Action != tea.MouseActionRelease {
		return m, nil
	}

	x, y := msg.X, msg.Y

	// Calculate panel boundaries
	leftWidth := m.Width / 4
	leftPanelEnd := leftWidth + 2

	// Click in left panel (files)
	if x >= 0 && x < leftPanelEnd {
		m.ActivePanel = FilePanel
		clickedLine := y - 3
		if clickedLine >= 0 && clickedLine < len(m.Files) {
			m.FileCursor = clickedLine
		}
		return m, nil
	}

	// Click in right panel (todos)
	if x >= leftPanelEnd && x < m.Width {
		m.ActivePanel = TodoPanel
		clickedLine := y - 3
		if clickedLine >= 0 && clickedLine < len(m.TodoList.Todos) {
			m.TodoCursor = clickedLine
			m.StatusMessage = fmt.Sprintf("Selected: %s", m.TodoList.Todos[clickedLine].Title)
		}
	}

	return m, nil
}
