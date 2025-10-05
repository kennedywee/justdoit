package main

import "github.com/charmbracelet/lipgloss"

// Enhanced Catppuccin Mocha color palette with LazyVim-inspired accents
var (
	// Base colors
	ColorBase       = lipgloss.Color("#181825") // deeper background
	ColorMantle     = lipgloss.Color("#1e1e2e") // surface background
	ColorCrust      = lipgloss.Color("#11111b") // darkest background
	ColorOverlay0   = lipgloss.Color("#6c7086") // muted text
	ColorOverlay1   = lipgloss.Color("#7f849c") // slightly less muted
	ColorText       = lipgloss.Color("#cdd6f4") // main text
	ColorSubtext0   = lipgloss.Color("#a6adc8") // dimmed text
	ColorSubtext1   = lipgloss.Color("#bac2de") // less dimmed text

	// Accent colors - more vibrant
	ColorLavender   = lipgloss.Color("#b4befe") // titles
	ColorBlue       = lipgloss.Color("#89b4fa") // active borders
	ColorSky        = lipgloss.Color("#89dceb") // info
	ColorSapphire   = lipgloss.Color("#74c7ec") // selection
	ColorTeal       = lipgloss.Color("#94e2d5") // selection alt
	ColorMauve      = lipgloss.Color("#cba6f7") // titles alt
	ColorPink       = lipgloss.Color("#f5c2e7") // special
	ColorMaroon     = lipgloss.Color("#eba0ac") // error alt
	ColorGreen      = lipgloss.Color("#a6e3a1") // success/hints
	ColorYellow     = lipgloss.Color("#f9e2af") // warning
	ColorRed        = lipgloss.Color("#f38ba8") // edit/danger
	ColorPeach      = lipgloss.Color("#fab387") // current file
	ColorFlamingo   = lipgloss.Color("#f2cdcd") // accent
	ColorRosewater  = lipgloss.Color("#f5e0dc") // subtle accent
)

// Custom border styles
var (
	// LazyVim-style double border
	LazyBorder = lipgloss.Border{
		Top:         "═",
		Bottom:      "═",
		Left:        "║",
		Right:       "║",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
	}

	// Thick border for active panels
	ThickBorder = lipgloss.Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}

	// Modern rounded border
	ModernBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}
)

// Styles returns all the lipgloss styles used in the application
type Styles struct {
	Selected      lipgloss.Style
	Border        lipgloss.Style
	ActiveBorder  lipgloss.Style
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	Completed     lipgloss.Style
	Hint          lipgloss.Style
	HintKey       lipgloss.Style
	Edit          lipgloss.Style
	Normal        lipgloss.Style
	Muted         lipgloss.Style
	Dimmed        lipgloss.Style
	CurrentFile   lipgloss.Style
	StatusBar     lipgloss.Style
	Shadow        lipgloss.Style
	Badge         lipgloss.Style
	Checkbox      lipgloss.Style
	CheckboxDone  lipgloss.Style
	Separator     lipgloss.Style
}

// NewStyles creates and returns all application styles
func NewStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().
			Foreground(ColorCrust).
			Background(ColorTeal).
			Bold(true).
			Padding(0, 1),

		Border: lipgloss.NewStyle().
			Border(ModernBorder).
			BorderForeground(ColorOverlay0),

		ActiveBorder: lipgloss.NewStyle().
			Border(ThickBorder).
			BorderForeground(ColorBlue).
			Bold(true),

		Title: lipgloss.NewStyle().
			Foreground(ColorLavender).
			Background(ColorMantle).
			Bold(true).
			Padding(0, 1).
			MarginBottom(0),

		Subtitle: lipgloss.NewStyle().
			Foreground(ColorSapphire).
			Italic(true),

		Completed: lipgloss.NewStyle().
			Foreground(ColorOverlay0).
			Strikethrough(true),

		Hint: lipgloss.NewStyle().
			Foreground(ColorSubtext1).
			Background(ColorMantle),

		HintKey: lipgloss.NewStyle().
			Foreground(ColorPeach).
			Background(ColorCrust).
			Bold(true).
			Padding(0, 1),

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
			Background(ColorCrust).
			Bold(true).
			Padding(0, 1),

		StatusBar: lipgloss.NewStyle().
			Foreground(ColorText).
			Background(ColorMantle).
			Padding(0, 1),

		Shadow: lipgloss.NewStyle().
			Foreground(ColorCrust),

		Badge: lipgloss.NewStyle().
			Foreground(ColorBase).
			Background(ColorMauve).
			Bold(true).
			Padding(0, 1),

		Checkbox: lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true),

		CheckboxDone: lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true),

		Separator: lipgloss.NewStyle().
			Foreground(ColorOverlay0),
	}
}
