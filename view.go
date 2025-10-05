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

	return mainView + "\n\n" + hints + statusBar
}

// renderFilePanel renders the left file panel
func (m model) renderFilePanel(width int) string {
	content := ""

	if m.mode == EditMode && m.editingIndex == -2 {
		// Creating new file
		content = m.styles.Edit.Render("  "+m.inputText+"█.json") + "\n"
		for _, file := range m.files {
			content += m.styles.Normal.Render("  󰈔 "+file) + "\n"
		}
	} else if m.showingArchive {
		// Show archived files
		content += m.styles.Separator.Render("  ─── archived ───") + "\n\n"
		for i, file := range m.archivedFiles {
			if m.activePanel == FilePanel && i == m.fileCursor {
				content += m.styles.Selected.Render(" "+file) + "\n"
			} else {
				content += m.styles.Dimmed.Render("  󰃨 "+file) + "\n"
			}
		}
	} else {
		// Show active files
		for i, file := range m.files {
			if m.activePanel == FilePanel && i == m.fileCursor && !m.showingArchive {
				content += m.styles.Selected.Render(" "+file) + "\n"
			} else if file == m.currentFile {
				content += m.styles.CurrentFile.Render("󰄲 "+file) + "\n"
			} else {
				content += m.styles.Normal.Render("  󰈔 "+file) + "\n"
			}
		}

		// Show archive section
		if len(m.archivedFiles) > 0 {
			content += "\n"
			content += m.styles.Separator.Render("  ─────────────") + "\n"
			content += m.styles.Badge.Render(fmt.Sprintf(" %d archived ", len(m.archivedFiles))) + "\n"
		}
	}

	// Apply border
	borderStyle := m.styles.Border
	if m.activePanel == FilePanel {
		borderStyle = m.styles.ActiveBorder
	}

	// Title with icon
	titleIcon := "󰈔"
	if m.showingArchive {
		titleIcon = "󰃨"
	}
	title := m.styles.Title.Render(fmt.Sprintf(" %s Files ", titleIcon))

	return borderStyle.
		Width(width).
		Height(m.height - 4).
		Padding(1, 1).
		Render(title + "\n\n" + content)
}

// renderTodoPanel renders the right todo panel
func (m model) renderTodoPanel(width int) string {
	content := ""

	if len(m.todoList.Todos) == 0 {
		emptyIcon := "󰄱"
		emptyMsg := m.styles.Dimmed.Italic(true).Render(fmt.Sprintf("  %s  No todos yet", emptyIcon))
		emptyHint := m.styles.Muted.Render("  Press 'a' to add one")
		content = emptyMsg + "\n" + emptyHint
	} else {
		content = m.renderTodoList()
	}

	// Apply border
	borderStyle := m.styles.Border
	if m.activePanel == TodoPanel {
		borderStyle = m.styles.ActiveBorder
	}

	// Title with stats
	completed := 0
	for _, todo := range m.todoList.Todos {
		if todo.Completed {
			completed++
		}
	}
	total := len(m.todoList.Todos)

	titleIcon := " "
	stats := ""
	if total > 0 {
		stats = m.styles.Badge.Render(fmt.Sprintf(" %d/%d ", completed, total))
	}

	title := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Title.Render(fmt.Sprintf(" %s %s ", titleIcon, m.currentFile)),
		" ",
		stats,
	)

	return borderStyle.
		Width(width).
		Height(m.height - 4).
		Padding(1, 1).
		Render(title + "\n\n" + content)
}

// renderTodoList renders the list of todos
func (m model) renderTodoList() string {
	content := ""

	// Show new todo input inline at the top
	if m.mode == EditMode && m.editingIndex == -1 {
		newCheckbox := m.styles.Checkbox.Render("")
		content += m.styles.Edit.Render(fmt.Sprintf("  %s  %s█", newCheckbox, m.inputText)) + "\n"
	}

	for i, todo := range m.todoList.Todos {
		var checkbox string
		var checkStyle lipgloss.Style

		if todo.Completed {
			checkbox = ""
			checkStyle = m.styles.CheckboxDone
		} else {
			checkbox = ""
			checkStyle = m.styles.Checkbox
		}

		checkboxStr := checkStyle.Render(checkbox)
		line := fmt.Sprintf("%s  %s", checkboxStr, todo.Title)

		// Apply style based on completion
		if todo.Completed {
			textStyle := m.styles.Completed
			line = fmt.Sprintf("%s  %s", checkboxStr, textStyle.Render(todo.Title))
		} else {
			line = fmt.Sprintf("%s  %s", checkboxStr, m.styles.Normal.Render(todo.Title))
		}

		// Handle editing mode
		if m.mode == EditMode && m.editingIndex == i {
			editIcon := m.styles.Edit.Render("")
			line = m.styles.Edit.Render(fmt.Sprintf(" %s  %s█", editIcon, m.inputText))
		} else if m.activePanel == TodoPanel && i == m.todoCursor {
			line = m.styles.Selected.Render(" " + line)
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
		Border(ThickBorder).
		BorderForeground(ColorRed).
		Padding(2, 4).
		Align(lipgloss.Center)

	titleIcon := lipgloss.NewStyle().
		Foreground(ColorRed).
		Bold(true).
		Render(" ")

	title := lipgloss.NewStyle().
		Foreground(ColorRed).
		Bold(true).
		Render("Delete Confirmation")

	titleBar := lipgloss.JoinHorizontal(lipgloss.Left, titleIcon, " ", title)

	filename := lipgloss.NewStyle().
		Foreground(ColorLavender).
		Background(ColorCrust).
		Bold(true).
		Padding(0, 1).
		Render(m.currentFile)

	question := m.styles.Normal.Render("Permanently delete this file?")

	yesKey := m.styles.HintKey.Render(" y ")
	yesText := m.styles.Hint.Render(" Yes, delete")
	noKey := m.styles.HintKey.Render(" n ")
	noText := m.styles.Hint.Render(" No, cancel")

	options := lipgloss.JoinHorizontal(
		lipgloss.Left,
		yesKey, yesText, "    ", noKey, noText,
	)

	confirmContent := lipgloss.JoinVertical(
		lipgloss.Center,
		titleBar,
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

	return panels
}

// renderArchiveConfirmation renders the archive confirmation dialog
func (m model) renderArchiveConfirmation() string {
	confirmStyle := lipgloss.NewStyle().
		Border(ThickBorder).
		BorderForeground(ColorSapphire).
		Padding(2, 4).
		Align(lipgloss.Center)

	titleIcon := lipgloss.NewStyle().
		Foreground(ColorSapphire).
		Bold(true).
		Render("󰃨 ")

	title := lipgloss.NewStyle().
		Foreground(ColorSapphire).
		Bold(true).
		Render("Archive Confirmation")

	titleBar := lipgloss.JoinHorizontal(lipgloss.Left, titleIcon, title)

	filename := lipgloss.NewStyle().
		Foreground(ColorLavender).
		Background(ColorCrust).
		Bold(true).
		Padding(0, 1).
		Render(m.currentFile)

	question := m.styles.Normal.Render("Archive this file?")

	yesKey := m.styles.HintKey.Render(" y ")
	yesText := m.styles.Hint.Render(" Yes, archive")
	noKey := m.styles.HintKey.Render(" n ")
	noText := m.styles.Hint.Render(" No, cancel")

	options := lipgloss.JoinHorizontal(
		lipgloss.Left,
		yesKey, yesText, "    ", noKey, noText,
	)

	confirmContent := lipgloss.JoinVertical(
		lipgloss.Center,
		titleBar,
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

	return panels
}

// renderHints renders the hints bar at the bottom
func (m model) renderHints() string {
	renderKey := func(key string) string {
		return m.styles.HintKey.Render(fmt.Sprintf(" %s ", key))
	}
	renderDesc := func(desc string) string {
		return m.styles.Hint.Render(fmt.Sprintf(" %s ", desc))
	}
	sep := m.styles.Muted.Render(" │ ")

	var hints []string

	if m.mode == EditMode {
		switch m.editingIndex {
		case -2:
			hints = []string{
				renderKey("Enter") + renderDesc("create"),
				renderKey("Esc") + renderDesc("cancel"),
			}
		case -3, -4:
			hints = []string{
				renderKey("y") + renderDesc("yes"),
				renderKey("n") + renderDesc("no"),
			}
		default:
			hints = []string{
				renderKey("Enter") + renderDesc("save"),
				renderKey("Esc") + renderDesc("cancel"),
			}
		}
	} else if m.activePanel == FilePanel {
		if m.showingArchive {
			hints = []string{
				renderKey("j/k") + renderDesc("navigate"),
				renderKey("Enter") + renderDesc("unarchive"),
				renderKey("z") + renderDesc("show active"),
				renderKey("h/l") + renderDesc("switch"),
				renderKey("q") + renderDesc("quit"),
			}
		} else {
			hints = []string{
				renderKey("j/k") + renderDesc("navigate"),
				renderKey("a") + renderDesc("new"),
				renderKey("d") + renderDesc("delete"),
				renderKey("Enter") + renderDesc("open"),
				renderKey("A") + renderDesc("archive"),
				renderKey("z") + renderDesc("archived"),
				renderKey("h/l") + renderDesc("switch"),
				renderKey("q") + renderDesc("quit"),
			}
		}
	} else {
		hints = []string{
			renderKey("j/k") + renderDesc("navigate"),
			renderKey("a") + renderDesc("add"),
			renderKey("i") + renderDesc("edit"),
			renderKey("d") + renderDesc("delete"),
			renderKey("x/Space") + renderDesc("toggle"),
			renderKey("h/l") + renderDesc("switch"),
			renderKey("q") + renderDesc("quit"),
		}
	}

	// Join hints with separator
	result := " "
	for i, hint := range hints {
		result += hint
		if i < len(hints)-1 {
			result += sep
		}
	}
	return result
}

// renderStatusBar renders the status message
func (m model) renderStatusBar() string {
	if m.statusMessage != "" && m.editingIndex != -3 && m.editingIndex != -4 {
		statusIcon := "󰙎 "
		statusStyle := lipgloss.NewStyle().
			Foreground(ColorGreen).
			Background(ColorMantle).
			Padding(0, 1).
			Bold(true)
		return "\n\n" + statusStyle.Render(statusIcon+m.statusMessage)
	}
	return ""
}
