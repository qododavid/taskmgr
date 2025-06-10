package tasks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileStore(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)

	// Initially empty
	list := store.List()
	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d tasks", len(list))
	}

	// Add a task
	if err := store.Add(Task{Title: "Example"}); err != nil {
		t.Fatalf("Add returned an error: %v", err)
	}

	list = store.List()
	if len(list) != 1 || list[0].Title != "Example" {
		t.Errorf("Expected one task 'Example', got %v", list)
	}

	// Update the task
	err = store.Update(0, Task{Title: "Updated", Done: true})
	if err != nil {
		t.Errorf("Update returned an error: %v", err)
	}

	list = store.List()
	if len(list) != 1 || list[0].Title != "Updated" || !list[0].Done {
		t.Errorf("Expected updated task 'Updated' with done=true, got %v", list)
	}

	// Test invalid update
	err = store.Update(1, Task{Title: "Invalid"})
	if err == nil {
		t.Error("Expected error updating invalid index, got nil")
	}

	// Ensure data persists by creating a new store and reading again
	newStore := NewFileStore(testFile)
	list = newStore.List()
	if len(list) != 1 || list[0].Title != "Updated" || !list[0].Done {
		t.Errorf("Expected persisted 'Updated' task, got %v", list)
	}
}

func TestFileStoreRemove(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_remove_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)

	// Add multiple tasks
	store.Add(Task{Title: "Task1"})
	store.Add(Task{Title: "Task2"})
	store.Add(Task{Title: "Task3"})

	list := store.List()
	if len(list) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(list))
	}

	// Remove middle task
	err = store.Remove(1)
	if err != nil {
		t.Errorf("Remove returned an error: %v", err)
	}

	list = store.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 tasks after removal, got %d", len(list))
	}
	if list[0].Title != "Task1" || list[1].Title != "Task3" {
		t.Errorf("Expected tasks 'Task1' and 'Task3', got %v", list)
	}

	// Test invalid remove
	err = store.Remove(5)
	if err == nil {
		t.Error("Expected error removing invalid index, got nil")
	}

	err = store.Remove(-1)
	if err == nil {
		t.Error("Expected error removing negative index, got nil")
	}
}

func TestFileStoreEmptyFile(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_empty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up

	testFile := filepath.Join(dir, "empty.json")
	// Create empty file
	ioutil.WriteFile(testFile, []byte(""), 0644)

	store := NewFileStore(testFile)
	list := store.List()
	if len(list) != 0 {
		t.Errorf("Expected empty list from empty file, got %d tasks", len(list))
	}
}

func TestFileStoreCorruptedFile(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_corrupt_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up

	testFile := filepath.Join(dir, "corrupt.json")
	// Create corrupted JSON file
	ioutil.WriteFile(testFile, []byte("invalid json"), 0644)

	store := NewFileStore(testFile)
	// Should handle corrupted file gracefully
	err = store.Add(Task{Title: "Test"})
	if err == nil {
		t.Error("Expected error when adding to corrupted file, got nil")
	}
}
