package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Todo represents a single todo item
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoList holds all todos and manages persistence
type TodoList struct {
	Todos    []Todo `json:"todos"`
	NextID   int    `json:"next_id"`
	filepath string
}

// NewTodoList creates a new TodoList
func NewTodoList(filepath string) *TodoList {
	tl := &TodoList{
		Todos:    []Todo{},
		NextID:   1,
		filepath: filepath,
	}
	tl.Load()
	return tl
}

// Add adds a new todo at the top
func (tl *TodoList) Add(title string) {
	todo := Todo{
		ID:        tl.NextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	// Insert at beginning
	tl.Todos = append([]Todo{todo}, tl.Todos...)
	tl.NextID++
	tl.Sort() // Keep completed at bottom
}

// Insert inserts a new todo at the top (always)
func (tl *TodoList) Insert(index int, title string) {
	todo := Todo{
		ID:        tl.NextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	tl.NextID++

	// Always insert at top
	tl.Todos = append([]Todo{todo}, tl.Todos...)
	tl.Sort() // Keep completed at bottom
}

// Delete removes a todo by index
func (tl *TodoList) Delete(index int) {
	if index >= 0 && index < len(tl.Todos) {
		tl.Todos = append(tl.Todos[:index], tl.Todos[index+1:]...)
		tl.Save()
	}
}

// Toggle toggles the completion status of a todo
func (tl *TodoList) Toggle(index int) {
	if index >= 0 && index < len(tl.Todos) {
		tl.Todos[index].Completed = !tl.Todos[index].Completed
		tl.Sort() // Auto-sort after toggling
	}
}

// Update updates a todo's title at a specific index
func (tl *TodoList) Update(index int, title string) {
	if index >= 0 && index < len(tl.Todos) {
		tl.Todos[index].Title = title
		tl.Save()
	}
}

// Sort sorts todos so completed ones are at the bottom
func (tl *TodoList) Sort() {
	// Stable sort: incomplete todos first, completed todos last
	// Preserves order within each group
	var incomplete []Todo
	var completed []Todo

	for _, todo := range tl.Todos {
		if todo.Completed {
			completed = append(completed, todo)
		} else {
			incomplete = append(incomplete, todo)
		}
	}

	tl.Todos = append(incomplete, completed...)
	tl.Save()
}

// Save persists the todo list to disk using atomic writes
func (tl *TodoList) Save() error {
	// Marshal data to JSON
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}

	// Create a temporary file in the same directory
	dir := filepath.Dir(tl.filepath)
	tmpFile, err := os.CreateTemp(dir, ".tui_todo_*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	// Close the temp file
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomically rename temp file to actual file
	// If this fails, the original file is unchanged
	if err := os.Rename(tmpPath, tl.filepath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// Load loads the todo list from disk with error recovery
func (tl *TodoList) Load() error {
	data, err := os.ReadFile(tl.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's ok
		}
		return fmt.Errorf("failed to read todo file: %w", err)
	}

	// Try to parse the JSON
	if err := json.Unmarshal(data, tl); err != nil {
		// If parsing fails, backup the corrupted file
		backupPath := tl.filepath + ".corrupted"
		if backupErr := os.WriteFile(backupPath, data, 0644); backupErr == nil {
			return fmt.Errorf("corrupted todo file backed up to %s: %w", backupPath, err)
		}
		return fmt.Errorf("corrupted todo file (backup failed): %w", err)
	}

	return nil
}
