package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// generateLargeTodoList creates a TodoList with n todos for testing
func generateLargeTodoList(n int) *TodoList {
	todos := make([]Todo, n)
	for i := 0; i < n; i++ {
		todos[i] = Todo{
			ID:        i + 1,
			Title:     fmt.Sprintf("Todo item number %d with some descriptive text", i+1),
			Completed: i%3 == 0, // ~33% completed
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		}
	}
	return &TodoList{
		Todos:  todos,
		NextID: n + 1,
	}
}

// createTestJSONFile creates a temporary JSON file with n todos
func createTestJSONFile(t *testing.T, n int) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	filepath := filepath.Join(tmpDir, "test_todos.json")

	tl := generateLargeTodoList(n)
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cleanup := func() {
		os.Remove(filepath)
	}

	return filepath, cleanup
}

// BenchmarkLoad_Small tests loading a small JSON file (100 todos)
func BenchmarkLoad_Small(b *testing.B) {
	benchmarkLoad(b, 100)
}

// BenchmarkLoad_Medium tests loading a medium JSON file (1,000 todos)
func BenchmarkLoad_Medium(b *testing.B) {
	benchmarkLoad(b, 1000)
}

// BenchmarkLoad_Large tests loading a large JSON file (10,000 todos)
func BenchmarkLoad_Large(b *testing.B) {
	benchmarkLoad(b, 10000)
}

// BenchmarkLoad_VeryLarge tests loading a very large JSON file (100,000 todos)
func BenchmarkLoad_VeryLarge(b *testing.B) {
	benchmarkLoad(b, 100000)
}

func benchmarkLoad(b *testing.B, numTodos int) {
	tmpDir := b.TempDir()
	filepath := filepath.Join(tmpDir, "benchmark_todos.json")

	// Create test file once
	tl := generateLargeTodoList(numTodos)
	data, _ := json.MarshalIndent(tl, "", "  ")
	os.WriteFile(filepath, data, 0644)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tl := &TodoList{filepath: filepath}
		if err := tl.Load(); err != nil {
			b.Fatalf("Load failed: %v", err)
		}
	}
}

// BenchmarkSort_Small tests sorting 100 todos
func BenchmarkSort_Small(b *testing.B) {
	benchmarkSort(b, 100)
}

// BenchmarkSort_Medium tests sorting 1,000 todos
func BenchmarkSort_Medium(b *testing.B) {
	benchmarkSort(b, 1000)
}

// BenchmarkSort_Large tests sorting 10,000 todos
func BenchmarkSort_Large(b *testing.B) {
	benchmarkSort(b, 10000)
}

// BenchmarkSort_VeryLarge tests sorting 100,000 todos
func BenchmarkSort_VeryLarge(b *testing.B) {
	benchmarkSort(b, 100000)
}

func benchmarkSort(b *testing.B, numTodos int) {
	tl := generateLargeTodoList(numTodos)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tl.Sort()
	}
}

// BenchmarkSave_Small tests saving 100 todos
func BenchmarkSave_Small(b *testing.B) {
	benchmarkSave(b, 100)
}

// BenchmarkSave_Medium tests saving 1,000 todos
func BenchmarkSave_Medium(b *testing.B) {
	benchmarkSave(b, 1000)
}

// BenchmarkSave_Large tests saving 10,000 todos
func BenchmarkSave_Large(b *testing.B) {
	benchmarkSave(b, 10000)
}

// BenchmarkSave_VeryLarge tests saving 100,000 todos
func BenchmarkSave_VeryLarge(b *testing.B) {
	benchmarkSave(b, 100000)
}

func benchmarkSave(b *testing.B, numTodos int) {
	tmpDir := b.TempDir()
	filepath := filepath.Join(tmpDir, "benchmark_save.json")

	tl := generateLargeTodoList(numTodos)
	tl.filepath = filepath

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := tl.Save(); err != nil {
			b.Fatalf("Save failed: %v", err)
		}
	}
}

// BenchmarkAddTodo tests adding todos to existing lists of various sizes
func BenchmarkAddTodo_Small(b *testing.B) {
	benchmarkAddTodo(b, 100)
}

func BenchmarkAddTodo_Medium(b *testing.B) {
	benchmarkAddTodo(b, 1000)
}

func BenchmarkAddTodo_Large(b *testing.B) {
	benchmarkAddTodo(b, 10000)
}

func benchmarkAddTodo(b *testing.B, numTodos int) {
	tl := generateLargeTodoList(numTodos)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tl.Add("New benchmark todo item")
	}
}

// BenchmarkToggleTodo tests toggling completion status
func BenchmarkToggleTodo_Small(b *testing.B) {
	benchmarkToggleTodo(b, 100)
}

func BenchmarkToggleTodo_Medium(b *testing.B) {
	benchmarkToggleTodo(b, 1000)
}

func BenchmarkToggleTodo_Large(b *testing.B) {
	benchmarkToggleTodo(b, 10000)
}

func benchmarkToggleTodo(b *testing.B, numTodos int) {
	tl := generateLargeTodoList(numTodos)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Toggle middle item
		tl.Toggle(numTodos / 2)
	}
}

// BenchmarkJSONUnmarshal tests raw JSON unmarshaling performance
func BenchmarkJSONUnmarshal_Small(b *testing.B) {
	benchmarkJSONUnmarshal(b, 100)
}

func BenchmarkJSONUnmarshal_Medium(b *testing.B) {
	benchmarkJSONUnmarshal(b, 1000)
}

func BenchmarkJSONUnmarshal_Large(b *testing.B) {
	benchmarkJSONUnmarshal(b, 10000)
}

func BenchmarkJSONUnmarshal_VeryLarge(b *testing.B) {
	benchmarkJSONUnmarshal(b, 100000)
}

func benchmarkJSONUnmarshal(b *testing.B, numTodos int) {
	tl := generateLargeTodoList(numTodos)
	data, _ := json.MarshalIndent(tl, "", "  ")

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		var result TodoList
		if err := json.Unmarshal(data, &result); err != nil {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}

// BenchmarkJSONMarshal tests raw JSON marshaling performance
func BenchmarkJSONMarshal_Small(b *testing.B) {
	benchmarkJSONMarshal(b, 100)
}

func BenchmarkJSONMarshal_Medium(b *testing.B) {
	benchmarkJSONMarshal(b, 1000)
}

func BenchmarkJSONMarshal_Large(b *testing.B) {
	benchmarkJSONMarshal(b, 10000)
}

func BenchmarkJSONMarshal_VeryLarge(b *testing.B) {
	benchmarkJSONMarshal(b, 100000)
}

func benchmarkJSONMarshal(b *testing.B, numTodos int) {
	tl := generateLargeTodoList(numTodos)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		data, err := json.MarshalIndent(tl, "", "  ")
		if err != nil {
			b.Fatalf("Marshal failed: %v", err)
		}
		b.SetBytes(int64(len(data)))
	}
}
