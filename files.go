package main

import (
	"os"
	"path/filepath"
)

// loadTodoFiles loads all .json todo files from a directory
func loadTodoFiles(dir string) []string {
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
func (m *model) deleteCurrentFile() {
	filePath := filepath.Join(m.todoDir, m.currentFile)
	os.Remove(filePath)

	// Reload file lists
	m.files = loadTodoFiles(m.todoDir)

	// Load next file or create default
	if len(m.files) > 0 {
		if m.fileCursor >= len(m.files) {
			m.fileCursor = len(m.files) - 1
		}
		m.currentFile = m.files[m.fileCursor]
		m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
	} else {
		m.currentFile = "default.json"
		m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
		m.todoList.Save()
		m.files = loadTodoFiles(m.todoDir)
		m.fileCursor = 0
	}
	m.todoCursor = 0
}

// archiveCurrentFile moves the current file to the archive directory
func (m *model) archiveCurrentFile() {
	srcPath := filepath.Join(m.todoDir, m.currentFile)
	dstPath := filepath.Join(m.archiveDir, m.currentFile)

	os.Rename(srcPath, dstPath)

	// Reload file lists
	m.files = loadTodoFiles(m.todoDir)
	m.archivedFiles = loadTodoFiles(m.archiveDir)

	// Load next file or create default
	if len(m.files) > 0 {
		m.fileCursor = 0
		m.currentFile = m.files[0]
		m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
	} else {
		m.currentFile = "default.json"
		m.todoList = NewTodoList(filepath.Join(m.todoDir, m.currentFile))
		m.todoList.Save()
		m.files = loadTodoFiles(m.todoDir)
	}
	m.todoCursor = 0
}

// unarchiveFile moves a file from the archive directory back to the main directory
func (m *model) unarchiveFile(filename string) {
	srcPath := filepath.Join(m.archiveDir, filename)
	dstPath := filepath.Join(m.todoDir, filename)

	os.Rename(srcPath, dstPath)

	// Reload file lists
	m.files = loadTodoFiles(m.todoDir)
	m.archivedFiles = loadTodoFiles(m.archiveDir)

	// Switch to the unarchived file
	m.currentFile = filename
	m.todoList = NewTodoList(dstPath)
	m.showingArchive = false

	// Find cursor position
	for i, f := range m.files {
		if f == filename {
			m.fileCursor = i
			break
		}
	}
}

// previewFile loads a file for preview without switching the active panel
func (m *model) previewFile() {
	var filename string
	var dir string

	if m.showingArchive {
		if m.fileCursor >= len(m.archivedFiles) {
			return
		}
		filename = m.archivedFiles[m.fileCursor]
		dir = m.archiveDir
	} else {
		if m.fileCursor >= len(m.files) {
			return
		}
		filename = m.files[m.fileCursor]
		dir = m.todoDir
	}

	// Load the file for preview (without switching activePanel)
	previewPath := filepath.Join(dir, filename)
	m.todoList = NewTodoList(previewPath)
	m.todoCursor = 0
}

// allTodosCompleted checks if all todos in the current list are completed
func (m *model) allTodosCompleted() bool {
	if len(m.todoList.Todos) == 0 {
		return false
	}
	for _, todo := range m.todoList.Todos {
		if !todo.Completed {
			return false
		}
	}
	return true
}
