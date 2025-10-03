# Project Plan: Go TUI Todo App (Lazyvim-inspired)

## Architecture Overview
- **TUI Library**: Bubbletea (recommended - modern, composable, Elm-inspired)
- **Layout**:
  - Left panel: Category navigation (Ongoing/Completed)
  - Right panel: Todo list display
  - Bottom: Hint bar with keybindings
- **Storage**: JSON file persistence
- **Aesthetic**: Lazyvim-inspired colors, borders, vim keybindings

## Key Features
1. **Dual-panel layout** with vim-like navigation
2. **Keybindings**: j/k (navigate), a (add), d (delete), x (toggle complete), tab (switch panels)
3. **Context hints** showing available actions
4. **Color coding** for different todo states
5. **Persistent storage** for todos

## Tech Stack
- **bubbletea**: TUI framework
- **lipgloss**: Styling/colors
- **bubbles**: UI components (list, textarea)
- Standard library for JSON persistence

## Implementation Tasks
1. Research and select Go TUI library (bubbletea/tview/termui)
2. Design the layout architecture (left panel + right panel + hint bar)
3. Implement core data models (Todo struct, persistence layer)
4. Build left panel navigation (Ongoing/Completed sections)
5. Build right panel todo list view with vim-like keybindings
6. Implement CRUD operations (Create, Read, Update, Delete todos)
7. Add bottom hint bar with context-aware navigation shortcuts
8. Implement todo state management (mark complete/incomplete)
9. Add vim-inspired keybindings (j/k navigation, dd delete, etc.)
10. Implement data persistence (JSON/SQLite file storage)
11. Add color scheme similar to lazyvim aesthetic
12. Polish UI with borders, titles, and status indicators
