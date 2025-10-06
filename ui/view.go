package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI (Bubble Tea interface)
func (m Model) View() string {
	if m.Width == 0 {
		return "Loading..."
	}

	leftWidth := m.Width / 4
	rightWidth := m.Width - leftWidth - 4

	// Calculate panel height based on whether status bar is showing
	panelHeight := m.Height - 4
	if m.StatusMessage != "" && m.EditingIndex != -3 && m.EditingIndex != -4 {
		panelHeight = m.Height - 7 // Account for status bar extra lines
	}

	// Render panels
	leftPanel := m.renderFilePanelWithHeight(leftWidth, panelHeight)
	rightPanel := m.renderTodoPanelWithHeight(rightWidth, panelHeight)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Handle special confirmation dialogs
	if m.Mode == EditMode && m.EditingIndex == -4 {
		return m.renderDeleteConfirmation()
	}

	if m.Mode == EditMode && m.EditingIndex == -3 {
		return m.renderArchiveConfirmation()
	}

	// Render hints and status
	hints := m.renderHints()
	statusBar := m.renderStatusBar()

	return mainView + "\n\n" + hints + statusBar
}

// renderFilePanelWithHeight renders the left file panel with specified height
func (m Model) renderFilePanelWithHeight(width int, height int) string {
	content := ""

	if m.Mode == EditMode && m.EditingIndex == -2 {
		// Creating new file
		content = m.Styles.Edit.Render("  "+m.InputText+"█.json") + "\n"
		for _, file := range m.Files {
			content += m.Styles.Normal.Render("  󰈔 "+file) + "\n"
		}
	} else if m.ShowingArchive {
		// Show archived files
		content += m.Styles.Separator.Render("  ─── archived ───") + "\n\n"
		for i, file := range m.ArchivedFiles {
			if m.ActivePanel == FilePanel && i == m.FileCursor {
				cursor := lipgloss.NewStyle().Foreground(ColorTeal).Render("▊")
				content += m.Styles.Selected.Render(" " + cursor + " " + file + " ") + "\n"
			} else {
				content += m.Styles.Dimmed.Render("  󰃨 "+file) + "\n"
			}
		}
	} else {
		// Show active files
		for i, file := range m.Files {
			if m.ActivePanel == FilePanel && i == m.FileCursor && !m.ShowingArchive {
				cursor := lipgloss.NewStyle().Foreground(ColorTeal).Render("▊")
				content += m.Styles.Selected.Render(" " + cursor + " " + file + " ") + "\n"
			} else if file == m.CurrentFile {
				content += m.Styles.CurrentFile.Render("󰄲 "+file) + "\n"
			} else {
				content += m.Styles.Normal.Render("  󰈔 "+file) + "\n"
			}
		}

		// Show archive section
		if len(m.ArchivedFiles) > 0 {
			content += "\n"
			content += m.Styles.Separator.Render("  ─────────────") + "\n"
			content += m.Styles.Badge.Render(fmt.Sprintf(" %d archived ", len(m.ArchivedFiles))) + "\n"
		}
	}

	// Apply border
	borderStyle := m.Styles.Border
	if m.ActivePanel == FilePanel {
		borderStyle = m.Styles.ActiveBorder
	}

	// Title with icon
	titleIcon := "󰈔"
	if m.ShowingArchive {
		titleIcon = "󰃨"
	}
	title := m.Styles.Title.Render(fmt.Sprintf(" %s Files ", titleIcon))

	return borderStyle.
		Width(width).
		Height(height).
		Padding(1, 1).
		Render(title + "\n\n" + content)
}

// renderTodoPanelWithHeight renders the right todo panel with specified height
func (m Model) renderTodoPanelWithHeight(width int, height int) string {
	content := ""

	// Always show renderTodoList when adding new todo to show input preview
	if m.Mode == EditMode && m.EditingIndex == -1 {
		content = m.renderTodoList()
	} else if len(m.TodoList.Todos) == 0 {
		emptyIcon := "󰄱"
		emptyMsg := m.Styles.Dimmed.Italic(true).Render(fmt.Sprintf("  %s  No todos yet", emptyIcon))
		emptyHint := m.Styles.Muted.Render("  Press 'a' to add one")
		content = emptyMsg + "\n" + emptyHint
	} else {
		content = m.renderTodoList()
	}

	// Apply border
	borderStyle := m.Styles.Border
	if m.ActivePanel == TodoPanel {
		borderStyle = m.Styles.ActiveBorder
	}

	// Title with stats
	completed := 0
	for _, todo := range m.TodoList.Todos {
		if todo.Completed {
			completed++
		}
	}
	total := len(m.TodoList.Todos)

	titleIcon := " "
	stats := ""
	if total > 0 {
		stats = m.Styles.Badge.Render(fmt.Sprintf(" %d/%d ", completed, total))
	}

	title := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.Styles.Title.Render(fmt.Sprintf(" %s %s ", titleIcon, m.CurrentFile)),
		" ",
		stats,
	)

	return borderStyle.
		Width(width).
		Height(height).
		Padding(1, 1).
		Render(title + "\n\n" + content)
}

// renderTodoList renders the list of todos
func (m Model) renderTodoList() string {
	content := ""

	// Show new todo input inline at the top
	if m.Mode == EditMode && m.EditingIndex == -1 {
		newCheckbox := m.Styles.Checkbox.Render("")
		content += m.Styles.Edit.Render(fmt.Sprintf("  %s  %s█", newCheckbox, m.InputText)) + "\n"
	}

	for i, todo := range m.TodoList.Todos {
		var checkbox string
		var checkStyle lipgloss.Style

		if todo.Completed {
			checkbox = ""
			checkStyle = m.Styles.CheckboxDone
		} else {
			checkbox = ""
			checkStyle = m.Styles.Checkbox
		}

		checkboxStr := checkStyle.Render(checkbox)
		line := fmt.Sprintf("%s  %s", checkboxStr, todo.Title)

		// Apply style based on completion
		if todo.Completed {
			textStyle := m.Styles.Completed
			line = fmt.Sprintf("%s  %s", checkboxStr, textStyle.Render(todo.Title))
		} else {
			line = fmt.Sprintf("%s  %s", checkboxStr, m.Styles.Normal.Render(todo.Title))
		}

		// Handle editing mode
		if m.Mode == EditMode && m.EditingIndex == i {
			editIcon := m.Styles.Edit.Render("")
			line = m.Styles.Edit.Render(fmt.Sprintf(" %s  %s█", editIcon, m.InputText))
		} else if m.ActivePanel == TodoPanel && i == m.TodoCursor {
			cursor := lipgloss.NewStyle().Foreground(ColorTeal).Render("▊")
			line = m.Styles.Selected.Render(" " + cursor + " " + line + " ")
		} else {
			line = "  " + line
		}

		content += line + "\n"
	}

	return content
}

// renderDeleteConfirmation renders the delete confirmation dialog
func (m Model) renderDeleteConfirmation() string {
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
		Render(m.CurrentFile)

	question := m.Styles.Normal.Render("Permanently delete this file?")

	yesKey := m.Styles.HintKey.Render(" y ")
	yesText := m.Styles.Hint.Render(" Yes, delete")
	noKey := m.Styles.HintKey.Render(" n ")
	noText := m.Styles.Hint.Render(" No, cancel")

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
		m.Width,
		m.Height-4,
		lipgloss.Center,
		lipgloss.Center,
		confirmBox,
	)

	return panels
}

// renderArchiveConfirmation renders the archive confirmation dialog
func (m Model) renderArchiveConfirmation() string {
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
		Render(m.CurrentFile)

	question := m.Styles.Normal.Render("Archive this file?")

	yesKey := m.Styles.HintKey.Render(" y ")
	yesText := m.Styles.Hint.Render(" Yes, archive")
	noKey := m.Styles.HintKey.Render(" n ")
	noText := m.Styles.Hint.Render(" No, cancel")

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
		m.Width,
		m.Height-4,
		lipgloss.Center,
		lipgloss.Center,
		confirmBox,
	)

	return panels
}

// renderHints renders the hints bar at the bottom
func (m Model) renderHints() string {
	renderKey := func(key string) string {
		return m.Styles.HintKey.Render(fmt.Sprintf(" %s ", key))
	}
	renderDesc := func(desc string) string {
		return m.Styles.Hint.Render(fmt.Sprintf(" %s ", desc))
	}
	sep := m.Styles.Muted.Render(" │ ")

	var hints []string

	if m.Mode == EditMode {
		switch m.EditingIndex {
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
	} else if m.ActivePanel == FilePanel {
		if m.ShowingArchive {
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
func (m Model) renderStatusBar() string {
	if m.StatusMessage != "" && m.EditingIndex != -3 && m.EditingIndex != -4 {
		statusIcon := "󰙎 "
		statusStyle := lipgloss.NewStyle().
			Foreground(ColorGreen).
			Background(ColorMantle).
			Padding(0, 1).
			Bold(true)
		return "\n\n" + statusStyle.Render(statusIcon+m.StatusMessage)
	}
	return ""
}
