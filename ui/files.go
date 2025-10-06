package ui

import (
	"os"
	"path/filepath"

	"justdoit/todo"
)

// LoadTodoFiles loads all .json todo files from a directory
func LoadTodoFiles(dir string) []string {
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

// deleteCurrentFile deletes the currently active file
func (m *Model) deleteCurrentFile() {
	filePath := filepath.Join(m.TodoDir, m.CurrentFile)
	os.Remove(filePath)

	// Reload file lists
	m.Files = LoadTodoFiles(m.TodoDir)

	// Load next file or create default
	if len(m.Files) > 0 {
		if m.FileCursor >= len(m.Files) {
			m.FileCursor = len(m.Files) - 1
		}
		m.CurrentFile = m.Files[m.FileCursor]
		m.TodoList = todo.NewTodoList(filepath.Join(m.TodoDir, m.CurrentFile))
	} else {
		m.CurrentFile = "default.json"
		m.TodoList = todo.NewTodoList(filepath.Join(m.TodoDir, m.CurrentFile))
		m.TodoList.Save()
		m.Files = LoadTodoFiles(m.TodoDir)
		m.FileCursor = 0
	}
	m.TodoCursor = 0
}

// archiveCurrentFile moves the current file to the archive directory
func (m *Model) archiveCurrentFile() {
	srcPath := filepath.Join(m.TodoDir, m.CurrentFile)
	dstPath := filepath.Join(m.ArchiveDir, m.CurrentFile)

	os.Rename(srcPath, dstPath)

	// Reload file lists
	m.Files = LoadTodoFiles(m.TodoDir)
	m.ArchivedFiles = LoadTodoFiles(m.ArchiveDir)

	// Load next file or create default
	if len(m.Files) > 0 {
		m.FileCursor = 0
		m.CurrentFile = m.Files[0]
		m.TodoList = todo.NewTodoList(filepath.Join(m.TodoDir, m.CurrentFile))
	} else {
		m.CurrentFile = "default.json"
		m.TodoList = todo.NewTodoList(filepath.Join(m.TodoDir, m.CurrentFile))
		m.TodoList.Save()
		m.Files = LoadTodoFiles(m.TodoDir)
	}
	m.TodoCursor = 0
}

// unarchiveFile moves a file from the archive directory back to the main directory
func (m *Model) unarchiveFile(filename string) {
	srcPath := filepath.Join(m.ArchiveDir, filename)
	dstPath := filepath.Join(m.TodoDir, filename)

	os.Rename(srcPath, dstPath)

	// Reload file lists
	m.Files = LoadTodoFiles(m.TodoDir)
	m.ArchivedFiles = LoadTodoFiles(m.ArchiveDir)

	// Switch to the unarchived file
	m.CurrentFile = filename
	m.TodoList = todo.NewTodoList(dstPath)
	m.ShowingArchive = false

	// Find cursor position
	for i, f := range m.Files {
		if f == filename {
			m.FileCursor = i
			break
		}
	}
}

// previewFile loads a file for preview without switching the active panel
func (m *Model) previewFile() {
	var filename string
	var dir string

	if m.ShowingArchive {
		if m.FileCursor >= len(m.ArchivedFiles) {
			return
		}
		filename = m.ArchivedFiles[m.FileCursor]
		dir = m.ArchiveDir
	} else {
		if m.FileCursor >= len(m.Files) {
			return
		}
		filename = m.Files[m.FileCursor]
		dir = m.TodoDir
	}

	// Load the file for preview (without switching activePanel)
	previewPath := filepath.Join(dir, filename)
	m.TodoList = todo.NewTodoList(previewPath)
	m.TodoCursor = 0
}

// allTodosCompleted checks if all todos in the current list are completed
func (m *Model) allTodosCompleted() bool {
	if len(m.TodoList.Todos) == 0 {
		return false
	}
	for _, todo := range m.TodoList.Todos {
		if !todo.Completed {
			return false
		}
	}
	return true
}
