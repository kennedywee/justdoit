package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLargeFileLoad tests loading files of various sizes and measures time
func TestLargeFileLoad(t *testing.T) {
	sizes := []int{100, 1000, 10000, 50000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Load_%d_todos", size), func(t *testing.T) {
			filepath, cleanup := createTestJSONFile(t, size)
			defer cleanup()

			start := time.Now()
			tl := NewTodoList(filepath)
			elapsed := time.Since(start)

			if len(tl.Todos) != size {
				t.Errorf("Expected %d todos, got %d", size, len(tl.Todos))
			}

			t.Logf("Loaded %d todos in %v", size, elapsed)

			// Fail if loading takes more than reasonable time
			maxDuration := time.Second * 2
			if elapsed > maxDuration {
				t.Errorf("Loading %d todos took %v, expected < %v", size, elapsed, maxDuration)
			}
		})
	}
}

// TestLargeFileSave tests saving files of various sizes and measures time
func TestLargeFileSave(t *testing.T) {
	sizes := []int{100, 1000, 10000, 50000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Save_%d_todos", size), func(t *testing.T) {
			tmpDir := t.TempDir()
			filepath := filepath.Join(tmpDir, "test_save.json")

			tl := generateLargeTodoList(size)
			tl.filepath = filepath

			start := time.Now()
			err := tl.Save()
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("Save failed: %v", err)
			}

			// Verify file was created
			stat, err := os.Stat(filepath)
			if err != nil {
				t.Fatalf("Saved file not found: %v", err)
			}

			t.Logf("Saved %d todos (%d bytes) in %v", size, stat.Size(), elapsed)

			// Fail if saving takes more than reasonable time
			maxDuration := time.Second * 3
			if elapsed > maxDuration {
				t.Errorf("Saving %d todos took %v, expected < %v", size, elapsed, maxDuration)
			}
		})
	}
}

// TestLargeFileIntegrity tests that large files maintain data integrity
func TestLargeFileIntegrity(t *testing.T) {
	sizes := []int{1000, 10000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Integrity_%d_todos", size), func(t *testing.T) {
			tmpDir := t.TempDir()
			filepath := filepath.Join(tmpDir, "integrity_test.json")

			// Create and save original data
			original := generateLargeTodoList(size)
			original.filepath = filepath
			if err := original.Save(); err != nil {
				t.Fatalf("Failed to save: %v", err)
			}

			// Load and verify
			loaded := NewTodoList(filepath)
			if len(loaded.Todos) != len(original.Todos) {
				t.Errorf("Todo count mismatch: got %d, want %d", len(loaded.Todos), len(original.Todos))
			}

			// Verify a sample of todos
			samplesToCheck := []int{0, size / 4, size / 2, 3 * size / 4, size - 1}
			for _, idx := range samplesToCheck {
				if idx >= len(loaded.Todos) {
					continue
				}

				if loaded.Todos[idx].ID != original.Todos[idx].ID {
					t.Errorf("ID mismatch at index %d: got %d, want %d", idx, loaded.Todos[idx].ID, original.Todos[idx].ID)
				}
				if loaded.Todos[idx].Title != original.Todos[idx].Title {
					t.Errorf("Title mismatch at index %d", idx)
				}
				if loaded.Todos[idx].Completed != original.Todos[idx].Completed {
					t.Errorf("Completed mismatch at index %d", idx)
				}
			}
		})
	}
}

// TestMemoryUsage tests memory consumption with large files
func TestMemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory test in short mode")
	}

	sizes := []int{1000, 10000, 50000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Memory_%d_todos", size), func(t *testing.T) {
			filepath, cleanup := createTestJSONFile(t, size)
			defer cleanup()

			// Force GC before measurement
			// runtime.GC()

			tl := NewTodoList(filepath)

			if len(tl.Todos) != size {
				t.Errorf("Expected %d todos, got %d", size, len(tl.Todos))
			}

			// Get file size for comparison
			stat, _ := os.Stat(filepath)
			fileSize := stat.Size()

			t.Logf("Loaded %d todos - File size: %d bytes (~%.2f MB)",
				size, fileSize, float64(fileSize)/(1024*1024))

			// Approximate memory per todo (rough estimate)
			// Each Todo struct is ~56 bytes + string overhead
			approxMemoryPerTodo := 100 // bytes (conservative estimate)
			expectedMemory := size * approxMemoryPerTodo

			t.Logf("Estimated memory usage: ~%.2f MB",
				float64(expectedMemory)/(1024*1024))
		})
	}
}

// TestSortPerformance tests sorting performance with various list sizes
func TestSortPerformance(t *testing.T) {
	sizes := []int{100, 1000, 10000, 50000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Sort_%d_todos", size), func(t *testing.T) {
			tl := generateLargeTodoList(size)

			start := time.Now()
			tl.Sort()
			elapsed := time.Since(start)

			t.Logf("Sorted %d todos in %v", size, elapsed)

			// Verify sort correctness: incomplete first, then completed
			foundCompleted := false
			for i, todo := range tl.Todos {
				if todo.Completed {
					foundCompleted = true
				} else if foundCompleted {
					t.Errorf("Found incomplete todo at index %d after completed todos", i)
					break
				}
			}

			// Performance threshold
			maxDuration := time.Millisecond * 100
			if elapsed > maxDuration {
				t.Errorf("Sorting %d todos took %v, expected < %v", size, elapsed, maxDuration)
			}
		})
	}
}

// TestConcurrentAccess tests thread-safety concerns (if any)
func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	filepath := filepath.Join(tmpDir, "concurrent_test.json")

	tl := generateLargeTodoList(1000)
	tl.filepath = filepath

	// Note: Current implementation is NOT thread-safe
	// This test documents expected behavior
	t.Run("Sequential_operations", func(t *testing.T) {
		if err := tl.Save(); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		tl2 := NewTodoList(filepath)
		if len(tl2.Todos) != len(tl.Todos) {
			t.Errorf("Todo count mismatch after sequential load")
		}
	})
}

// TestFileCorruption tests recovery from corrupted files
func TestFileCorruptionHandling(t *testing.T) {
	tmpDir := t.TempDir()
	filepath := filepath.Join(tmpDir, "corrupted.json")

	// Write invalid JSON
	invalidJSON := []byte(`{"todos": [{"id": 1, "title": "broken"`)
	if err := os.WriteFile(filepath, invalidJSON, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	tl := &TodoList{filepath: filepath}
	err := tl.Load()

	if err == nil {
		t.Error("Expected error when loading corrupted file")
	}

	// Check if backup was created
	backupPath := filepath + ".corrupted"
	if _, err := os.Stat(backupPath); err != nil {
		t.Error("Expected corrupted file backup to be created")
	}

	t.Logf("Corruption handled correctly: %v", err)
}

// TestJSONFileSize tests actual file sizes produced
func TestJSONFileSize(t *testing.T) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("FileSize_%d_todos", size), func(t *testing.T) {
			tmpDir := t.TempDir()
			filepath := filepath.Join(tmpDir, "size_test.json")

			tl := generateLargeTodoList(size)
			tl.filepath = filepath

			if err := tl.Save(); err != nil {
				t.Fatalf("Save failed: %v", err)
			}

			stat, err := os.Stat(filepath)
			if err != nil {
				t.Fatalf("Stat failed: %v", err)
			}

			bytesPerTodo := float64(stat.Size()) / float64(size)
			t.Logf("%d todos = %d bytes (%.2f bytes/todo, %.2f MB total)",
				size, stat.Size(), bytesPerTodo, float64(stat.Size())/(1024*1024))

			// Check file contents are valid JSON
			data, _ := os.ReadFile(filepath)
			var check TodoList
			if err := json.Unmarshal(data, &check); err != nil {
				t.Errorf("Saved file contains invalid JSON: %v", err)
			}
		})
	}
}
