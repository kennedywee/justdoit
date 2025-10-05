package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI (Bubble Tea interface)
func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - 4

	// Render panels
	leftPanel := m.renderFilePanel(leftWidth)
	rightPanel := m.renderTodoPanel(rightWidth)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Handle special confirmation dialogs
	if m.mode == EditMode && m.editingIndex == -4 {
		return m.renderDeleteConfirmation()
	}

	if m.mode == EditMode && m.editingIndex == -3 {
		return m.renderArchiveConfirmation()
	}

	// Render hints and status
	hints := m.renderHints()
	statusBar := m.renderStatusBar()

	return mainView + "\n" + hints + statusBar
}

// renderFilePanel renders the left file panel
func (m model) renderFilePanel(width int) string {
	content := ""

	if m.mode == EditMode && m.editingIndex == -2 {
		// Creating new file
		content = m.styles.Edit.Render("▶ "+m.inputText+"█.json") + "\n"
		for _, file := range m.files {
			content += "  " + file + "\n"
		}
	} else if m.showingArchive {
		// Show archived files
		content += m.styles.Muted.Render("╌╌╌ archived ╌╌╌") + "\n"
		for i, file := range m.archivedFiles {
			if m.activePanel == FilePanel && i == m.fileCursor {
				content += m.styles.Selected.Render("▶ "+file) + "\n"
			} else {
				content += m.styles.Dimmed.Render("  "+file) + "\n"
			}
		}
	} else {
		// Show active files
		for i, file := range m.files {
			if m.activePanel == FilePanel && i == m.fileCursor && !m.showingArchive {
				content += m.styles.Selected.Render("▶ "+file) + "\n"
			} else if file == m.currentFile {
				content += m.styles.CurrentFile.Render("● "+file) + "\n"
			} else {
				content += m.styles.Normal.Render("  "+file) + "\n"
			}
		}

		// Show archive section
		if len(m.archivedFiles) > 0 {
			content += "\n"
			content += m.styles.Muted.Render("╌╌╌ archived ╌╌╌") + "\n"
			content += m.styles.Muted.Render(fmt.Sprintf("  %d files", len(m.archivedFiles))) + "\n"
		}
	}

	// Apply border
	borderStyle := m.styles.Border
	if m.activePanel == FilePanel {
		borderStyle = m.styles.ActiveBorder
	}

	return borderStyle.
		Width(width).
		Height(m.height - 4).
		Padding(0, 1).
		Render(m.styles.Title.Render("📁 Files") + "\n\n" + content)
}

// renderTodoPanel renders the right todo panel
func (m model) renderTodoPanel(width int) string {
	content := ""

	if len(m.todoList.Todos) == 0 {
		content = m.styles.Muted.Italic(true).Render("No todos yet. Press 'a' to add one.")
	} else {
		content = m.renderTodoList()
	}

	// Apply border
	borderStyle := m.styles.Border
	if m.activePanel == TodoPanel {
		borderStyle = m.styles.ActiveBorder
	}

	return borderStyle.
		Width(width).
		Height(m.height - 4).
		Padding(0, 1).
		Render(m.styles.Title.Render(fmt.Sprintf("✓ %s", m.currentFile)) + "\n\n" + content)
}

// renderTodoList renders the list of todos
func (m model) renderTodoList() string {
	content := ""

	// Show new todo input inline at the top
	if m.mode == EditMode && m.editingIndex == -1 {
		content += m.styles.Edit.Render("▶ ○ "+m.inputText+"█") + "\n"
	}

	for i, todo := range m.todoList.Todos {
		checkbox := "○"
		if todo.Completed {
			checkbox = "●"
		}

		line := fmt.Sprintf("%s %s", checkbox, todo.Title)

		// Apply style based on completion
		if todo.Completed {
			line = m.styles.Completed.Render(line)
		} else {
			line = m.styles.Normal.Render(line)
		}

		// Handle editing mode
		if m.mode == EditMode && m.editingIndex == i {
			line = m.styles.Edit.Render("✎ " + m.inputText + "█")
		} else if m.activePanel == TodoPanel && i == m.todoCursor {
			line = m.styles.Selected.Render("▶ " + line)
		} else {
			line = "  " + line
		}

		content += line + "\n"
	}

	return content
}

// renderDeleteConfirmation renders the delete confirmation dialog
func (m model) renderDeleteConfirmation() string {
	confirmStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorRed).
		Padding(2, 4).
		Align(lipgloss.Center)

	title := lipgloss.NewStyle().
		Foreground(ColorRed).
		Bold(true).
		Render("⚠ Delete Confirmation")

	filename := lipgloss.NewStyle().
		Foreground(ColorLavender).
		Bold(true).
		Render(fmt.Sprintf("'%s'", m.currentFile))

	question := m.styles.Normal.Render("Permanently delete this file?")
	options := m.styles.Hint.Render("[y] Yes, delete    [n] No, cancel")

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

	return panels + "\n" + m.styles.Hint.Render("y: delete | n: cancel")
}

// renderArchiveConfirmation renders the archive confirmation dialog
func (m model) renderArchiveConfirmation() string {
	confirmStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBlue).
		Padding(2, 4).
		Align(lipgloss.Center)

	title := lipgloss.NewStyle().
		Foreground(ColorBlue).
		Bold(true).
		Render("Archive Confirmation")

	filename := lipgloss.NewStyle().
		Foreground(ColorLavender).
		Bold(true).
		Render(fmt.Sprintf("'%s'", m.currentFile))

	question := m.styles.Normal.Render("Archive this file?")
	options := m.styles.Hint.Render("[y] Yes, archive    [n] No, cancel")

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

	return panels + "\n" + m.styles.Hint.Render("y: archive | n: cancel")
}

// renderHints renders the hints bar at the bottom
func (m model) renderHints() string {
	if m.mode == EditMode {
		switch m.editingIndex {
		case -2:
			return m.styles.Hint.Render("Enter: create file | Esc: cancel")
		case -3:
			return m.styles.Hint.Render("y: archive | n: cancel")
		case -4:
			return m.styles.Hint.Render("y: delete | n: cancel")
		default:
			return m.styles.Hint.Render("Enter: save | Esc: cancel")
		}
	} else if m.activePanel == FilePanel {
		if m.showingArchive {
			return m.styles.Hint.Render("j/k: navigate | Enter: unarchive | z: show active | h/l: switch panel | q: quit")
		}
		return m.styles.Hint.Render("j/k: navigate | a: new | d: delete | Enter: open | A: archive | z: show archived | h/l: switch | q: quit")
	}
	return m.styles.Hint.Render("j/k: navigate | a: add | i: edit | d: delete | x/space: toggle | h/l: switch panel | q: quit")
}

// renderStatusBar renders the status message
func (m model) renderStatusBar() string {
	if m.statusMessage != "" && m.editingIndex != -3 && m.editingIndex != -4 {
		return "\n" + m.styles.Hint.Render(m.statusMessage)
	}
	return ""
}
