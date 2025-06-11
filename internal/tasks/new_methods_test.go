package tasks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestClearCompleted(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_clear_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add some tasks
	manager.Add(Task{Title: "Task 1", Done: true})
	manager.Add(Task{Title: "Task 2", Done: false})
	manager.Add(Task{Title: "Task 3", Done: true})
	manager.Add(Task{Title: "Task 4", Done: false})

	// Clear completed tasks
	err = manager.ClearCompleted()
	if err != nil {
		t.Errorf("ClearCompleted returned error: %v", err)
	}

	// Check remaining tasks
	remaining := manager.List()
	if len(remaining) != 2 {
		t.Errorf("Expected 2 remaining tasks, got %d", len(remaining))
	}

	for _, task := range remaining {
		if task.Done {
			t.Errorf("Found completed task after ClearCompleted: %v", task)
		}
	}

	// Test with empty list
	manager.ClearCompleted()
	remaining = manager.List()
	if len(remaining) != 2 {
		t.Errorf("ClearCompleted on partial list should not affect undone tasks")
	}
}

func TestCountUndone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_count_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test empty list
	count := manager.CountUndone()
	if count != 0 {
		t.Errorf("Expected 0 undone tasks in empty list, got %d", count)
	}

	// Add mixed tasks
	manager.Add(Task{Title: "Task 1", Done: true})
	manager.Add(Task{Title: "Task 2", Done: false})
	manager.Add(Task{Title: "Task 3", Done: true})
	manager.Add(Task{Title: "Task 4", Done: false})
	manager.Add(Task{Title: "Task 5", Done: false})

	count = manager.CountUndone()
	if count != 3 {
		t.Errorf("Expected 3 undone tasks, got %d", count)
	}

	// Mark all as done
	manager.MarkAllDone()
	count = manager.CountUndone()
	if count != 0 {
		t.Errorf("Expected 0 undone tasks after MarkAllDone, got %d", count)
	}
}

func TestListUndone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_list_undone_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test empty list
	undone := manager.ListUndone()
	if len(undone) != 0 {
		t.Errorf("Expected empty list for undone tasks, got %d", len(undone))
	}

	// Add mixed tasks
	manager.Add(Task{Title: "Done Task 1", Done: true})
	manager.Add(Task{Title: "Undone Task 1", Done: false})
	manager.Add(Task{Title: "Done Task 2", Done: true})
	manager.Add(Task{Title: "Undone Task 2", Done: false})

	undone = manager.ListUndone()
	if len(undone) != 2 {
		t.Errorf("Expected 2 undone tasks, got %d", len(undone))
	}

	expectedTitles := []string{"Undone Task 1", "Undone Task 2"}
	for i, task := range undone {
		if task.Done {
			t.Errorf("Found done task in undone list: %v", task)
		}
		if task.Title != expectedTitles[i] {
			t.Errorf("Expected title %s, got %s", expectedTitles[i], task.Title)
		}
	}
}

func TestToggleDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_toggle_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add a task
	manager.Add(Task{Title: "Toggle Task", Done: false})

	// Toggle to done
	err = manager.ToggleDone("0")
	if err != nil {
		t.Errorf("ToggleDone returned error: %v", err)
	}

	tasks := manager.List()
	if !tasks[0].Done {
		t.Errorf("Expected task to be done after toggle")
	}

	// Toggle back to undone
	err = manager.ToggleDone("0")
	if err != nil {
		t.Errorf("ToggleDone returned error: %v", err)
	}

	tasks = manager.List()
	if tasks[0].Done {
		t.Errorf("Expected task to be undone after second toggle")
	}

	// Test invalid index
	err = manager.ToggleDone("10")
	if err == nil {
		t.Errorf("Expected error for invalid index")
	}

	// Test invalid index string
	err = manager.ToggleDone("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid index string")
	}
}

func TestIsEmpty(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_empty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test empty list
	if !manager.IsEmpty() {
		t.Errorf("Expected empty list to return true")
	}

	// Add a task
	manager.Add(Task{Title: "Test Task"})
	if manager.IsEmpty() {
		t.Errorf("Expected non-empty list to return false")
	}

	// Clear all tasks by removing the only task
	manager.Remove("0")
	if !manager.IsEmpty() {
		t.Errorf("Expected empty list after removal to return true")
	}
}

func TestGetTask(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_get_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add some tasks
	manager.Add(Task{Title: "First Task", Description: "First"})
	manager.Add(Task{Title: "Second Task", Description: "Second"})

	// Get valid task
	task, err := manager.GetTask("0")
	if err != nil {
		t.Errorf("GetTask returned error: %v", err)
	}
	if task == nil {
		t.Errorf("GetTask returned nil task")
	}
	if task.Title != "First Task" {
		t.Errorf("Expected 'First Task', got %s", task.Title)
	}

	// Get second task
	task, err = manager.GetTask("1")
	if err != nil {
		t.Errorf("GetTask returned error: %v", err)
	}
	if task.Title != "Second Task" {
		t.Errorf("Expected 'Second Task', got %s", task.Title)
	}

	// Test invalid index
	task, err = manager.GetTask("10")
	if err == nil {
		t.Errorf("Expected error for invalid index")
	}
	if task != nil {
		t.Errorf("Expected nil task for invalid index")
	}

	// Test invalid index string
	task, err = manager.GetTask("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid index string")
	}
	if task != nil {
		t.Errorf("Expected nil task for invalid index string")
	}
}

func TestUpdateTitle(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_update_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add a task
	manager.Add(Task{Title: "Original Title", Description: "Test"})

	// Update title
	err = manager.UpdateTitle("0", "New Title")
	if err != nil {
		t.Errorf("UpdateTitle returned error: %v", err)
	}

	// Verify update
	tasks := manager.List()
	if tasks[0].Title != "New Title" {
		t.Errorf("Expected 'New Title', got %s", tasks[0].Title)
	}
	if tasks[0].Description != "Test" {
		t.Errorf("Description should remain unchanged")
	}

	// Test persistence
	newStore := NewFileStore(testFile)
	newManager := NewTaskManager(newStore)
	newTasks := newManager.List()
	if newTasks[0].Title != "New Title" {
		t.Errorf("Title update not persisted")
	}

	// Test invalid index
	err = manager.UpdateTitle("10", "Should Fail")
	if err == nil {
		t.Errorf("Expected error for invalid index")
	}

	// Test invalid index string
	err = manager.UpdateTitle("invalid", "Should Fail")
	if err == nil {
		t.Errorf("Expected error for invalid index string")
	}

	// Test empty title
	err = manager.UpdateTitle("0", "")
	if err != nil {
		t.Errorf("UpdateTitle should allow empty title: %v", err)
	}
	tasks = manager.List()
	if tasks[0].Title != "" {
		t.Errorf("Expected empty title, got %s", tasks[0].Title)
	}
}