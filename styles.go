package main

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha color palette
var (
	// Base colors
	ColorBase     = lipgloss.Color("#1e1e2e") // background
	ColorOverlay0 = lipgloss.Color("#6c7086") // muted text
	ColorText     = lipgloss.Color("#cdd6f4") // main text
	ColorSubtext0 = lipgloss.Color("#a6adc8") // dimmed text

	// Accent colors
	ColorLavender = lipgloss.Color("#b4befe") // titles
	ColorBlue     = lipgloss.Color("#89b4fa") // active borders
	ColorSapphire = lipgloss.Color("#74c7ec") // selection
	ColorMauve    = lipgloss.Color("#cba6f7") // titles alt
	ColorGreen    = lipgloss.Color("#a6e3a1") // success/hints
	ColorRed      = lipgloss.Color("#f38ba8") // edit/danger
	ColorPeach    = lipgloss.Color("#fab387") // current file
)

// Styles returns all the lipgloss styles used in the application
type Styles struct {
	Selected      lipgloss.Style
	Border        lipgloss.Style
	ActiveBorder  lipgloss.Style
	Title         lipgloss.Style
	Completed     lipgloss.Style
	Hint          lipgloss.Style
	Edit          lipgloss.Style
	Normal        lipgloss.Style
	Muted         lipgloss.Style
	Dimmed        lipgloss.Style
	CurrentFile   lipgloss.Style
}

// NewStyles creates and returns all application styles
func NewStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().
			Foreground(ColorBase).
			Background(ColorSapphire).
			Bold(true),

		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorOverlay0),

		ActiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBlue).
			Bold(true),

		Title: lipgloss.NewStyle().
			Foreground(ColorMauve).
			Bold(true),

		Completed: lipgloss.NewStyle().
			Foreground(ColorOverlay0).
			Strikethrough(true),

		Hint: lipgloss.NewStyle().
			Foreground(ColorGreen),

		Edit: lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true),

		Normal: lipgloss.NewStyle().
			Foreground(ColorText),

		Muted: lipgloss.NewStyle().
			Foreground(ColorOverlay0),

		Dimmed: lipgloss.NewStyle().
			Foreground(ColorSubtext0),

		CurrentFile: lipgloss.NewStyle().
			Foreground(ColorPeach).
			Bold(true),
	}
}
