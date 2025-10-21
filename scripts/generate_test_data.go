// Package main provides a utility to generate large test JSON files for performance testing
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoList struct {
	Todos  []Todo `json:"todos"`
	NextID int    `json:"next_id"`
}

var (
	count      = flag.Int("count", 1000, "Number of todos to generate")
	output     = flag.String("output", "", "Output file path (default: ~/.tui_todos/test_<count>.json)")
	completion = flag.Float64("completion", 33.0, "Percentage of todos marked as completed (0-100)")
)

func main() {
	flag.Parse()

	if *count <= 0 {
		log.Fatal("Count must be positive")
	}

	if *completion < 0 || *completion > 100 {
		log.Fatal("Completion percentage must be between 0 and 100")
	}

	// Determine output path
	outputPath := *output
	if outputPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v", err)
		}
		todoDir := filepath.Join(homeDir, ".tui_todos")
		if err := os.MkdirAll(todoDir, 0755); err != nil {
			log.Fatalf("Failed to create todo directory: %v", err)
		}
		outputPath = filepath.Join(todoDir, fmt.Sprintf("test_%d.json", *count))
	}

	fmt.Printf("Generating %d todos...\n", *count)
	start := time.Now()

	// Generate todos
	todos := make([]Todo, *count)
	completedThreshold := int(float64(*count) * (*completion / 100.0))

	for i := 0; i < *count; i++ {
		// Vary todo titles for realism
		titleVariants := []string{
			"Buy groceries for the week",
			"Complete project documentation",
			"Review pull request #%d",
			"Fix bug in authentication module",
			"Update dependencies to latest versions",
			"Write unit tests for %s module",
			"Refactor legacy code in %s package",
			"Schedule team meeting for Q%d planning",
			"Optimize database query performance",
			"Deploy hotfix to production",
			"Research new framework alternatives",
			"Create user onboarding flow",
			"Implement dark mode toggle",
			"Add error logging to API endpoints",
			"Update README with installation instructions",
		}

		variant := titleVariants[i%len(titleVariants)]
		title := fmt.Sprintf(variant+" [Item #%d]", i+1)

		todos[i] = Todo{
			ID:        i + 1,
			Title:     title,
			Completed: i < completedThreshold,
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Minute),
		}
	}

	todoList := TodoList{
		Todos:  todos,
		NextID: *count + 1,
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(todoList, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	elapsed := time.Since(start)

	// Get file size
	stat, _ := os.Stat(outputPath)
	fileSizeMB := float64(stat.Size()) / (1024 * 1024)

	fmt.Printf("\n✓ Successfully generated test file\n")
	fmt.Printf("  Path: %s\n", outputPath)
	fmt.Printf("  Todos: %d (%d completed, %d incomplete)\n",
		*count, completedThreshold, *count-completedThreshold)
	fmt.Printf("  Size: %.2f MB (%d bytes)\n", fileSizeMB, stat.Size())
	fmt.Printf("  Time: %v\n", elapsed)
}
