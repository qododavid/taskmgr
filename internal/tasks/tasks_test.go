package tasks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTaskManager(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_tasks_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add a task
	if err := manager.Add(Task{Title: "First"}); err != nil {
		t.Fatal("Error adding task:", err)
	}

	list := manager.List()
	if len(list) != 1 || list[0].Title != "First" {
		t.Errorf("Expected one task with title 'First', got %v", list)
	}

	// Mark it done
	err = manager.MarkDone("0")
	if err != nil {
		t.Error("MarkDone returned an error:", err)
	}

	list = manager.List()
	if !list[0].Done {
		t.Error("Expected task to be marked done")
	}

	// Ensure persistence: create a new manager and check again
	newStore := NewFileStore(testFile)
	newManager := NewTaskManager(newStore)
	newList := newManager.List()
	if len(newList) != 1 || newList[0].Title != "First" || !newList[0].Done {
		t.Errorf("Expected persisted task 'First' to be done, got %v", newList)
	}
}

func TestClearCompleted(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_clear_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test with no tasks
	err = manager.ClearCompleted()
	if err != nil {
		t.Error("ClearCompleted should not error on empty list:", err)
	}

	// Add mixed tasks
	manager.Add(Task{Title: "Task1", Done: false})
	manager.Add(Task{Title: "Task2", Done: true})
	manager.Add(Task{Title: "Task3", Done: false})
	manager.Add(Task{Title: "Task4", Done: true})

	// Clear completed
	err = manager.ClearCompleted()
	if err != nil {
		t.Error("ClearCompleted failed:", err)
	}

	// Check remaining tasks
	list := manager.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 remaining tasks, got %d", len(list))
	}
	for _, task := range list {
		if task.Done {
			t.Error("Found completed task after ClearCompleted")
		}
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

	// Test with no tasks
	count := manager.CountUndone()
	if count != 0 {
		t.Errorf("Expected 0 undone tasks, got %d", count)
	}

	// Add mixed tasks
	manager.Add(Task{Title: "Task1", Done: false})
	manager.Add(Task{Title: "Task2", Done: true})
	manager.Add(Task{Title: "Task3", Done: false})

	count = manager.CountUndone()
	if count != 2 {
		t.Errorf("Expected 2 undone tasks, got %d", count)
	}

	// Mark all done
	manager.MarkDone("0")
	manager.MarkDone("2")

	count = manager.CountUndone()
	if count != 0 {
		t.Errorf("Expected 0 undone tasks after marking all done, got %d", count)
	}
}

func TestListUndone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_listundone_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test with no tasks
	undone := manager.ListUndone()
	if len(undone) != 0 {
		t.Errorf("Expected 0 undone tasks, got %d", len(undone))
	}

	// Add mixed tasks
	manager.Add(Task{Title: "Undone1", Done: false})
	manager.Add(Task{Title: "Done1", Done: true})
	manager.Add(Task{Title: "Undone2", Done: false})

	undone = manager.ListUndone()
	if len(undone) != 2 {
		t.Errorf("Expected 2 undone tasks, got %d", len(undone))
	}

	// Verify correct tasks are returned
	expectedTitles := map[string]bool{"Undone1": true, "Undone2": true}
	for _, task := range undone {
		if task.Done {
			t.Error("Found completed task in undone list")
		}
		if !expectedTitles[task.Title] {
			t.Errorf("Unexpected task in undone list: %s", task.Title)
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

	// Add tasks
	manager.Add(Task{Title: "Task1", Done: false})
	manager.Add(Task{Title: "Task2", Done: true})

	// Toggle undone task to done
	err = manager.ToggleDone("0")
	if err != nil {
		t.Error("ToggleDone failed:", err)
	}

	list := manager.List()
	if !list[0].Done {
		t.Error("Expected task 0 to be done after toggle")
	}

	// Toggle done task to undone
	err = manager.ToggleDone("1")
	if err != nil {
		t.Error("ToggleDone failed:", err)
	}

	list = manager.List()
	if list[1].Done {
		t.Error("Expected task 1 to be undone after toggle")
	}

	// Test invalid index
	err = manager.ToggleDone("10")
	if err == nil {
		t.Error("Expected error for invalid index")
	}

	// Test invalid index string
	err = manager.ToggleDone("invalid")
	if err == nil {
		t.Error("Expected error for invalid index string")
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

	// Test with empty list
	if !manager.IsEmpty() {
		t.Error("Expected empty task list")
	}

	// Add a task
	manager.Add(Task{Title: "Task1"})

	// Test with non-empty list
	if manager.IsEmpty() {
		t.Error("Expected non-empty task list")
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

	// Add tasks
	manager.Add(Task{Title: "Task1", Description: "First task"})
	manager.Add(Task{Title: "Task2", Description: "Second task"})

	// Test valid index
	task, err := manager.GetTask("0")
	if err != nil {
		t.Error("GetTask failed:", err)
	}
	if task == nil || task.Title != "Task1" {
		t.Errorf("Expected Task1, got %v", task)
	}

	// Test another valid index
	task, err = manager.GetTask("1")
	if err != nil {
		t.Error("GetTask failed:", err)
	}
	if task == nil || task.Title != "Task2" {
		t.Errorf("Expected Task2, got %v", task)
	}

	// Test invalid index
	task, err = manager.GetTask("10")
	if err == nil {
		t.Error("Expected error for invalid index")
	}
	if task != nil {
		t.Error("Expected nil task for invalid index")
	}

	// Test invalid index string
	task, err = manager.GetTask("invalid")
	if err == nil {
		t.Error("Expected error for invalid index string")
	}
	if task != nil {
		t.Error("Expected nil task for invalid index string")
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
	manager.Add(Task{Title: "Original Title", Description: "Test task"})

	// Update title
	err = manager.UpdateTitle("0", "New Title")
	if err != nil {
		t.Error("UpdateTitle failed:", err)
	}

	// Verify update
	list := manager.List()
	if list[0].Title != "New Title" {
		t.Errorf("Expected 'New Title', got '%s'", list[0].Title)
	}
	if list[0].Description != "Test task" {
		t.Error("Description should remain unchanged")
	}

	// Test invalid index
	err = manager.UpdateTitle("10", "Should Fail")
	if err == nil {
		t.Error("Expected error for invalid index")
	}

	// Test invalid index string
	err = manager.UpdateTitle("invalid", "Should Fail")
	if err == nil {
		t.Error("Expected error for invalid index string")
	}
}

// Test additional functions for better coverage
func TestFindByTitle(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_find_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks
	manager.Add(Task{Title: "Unique Task", Description: "Test"})
	manager.Add(Task{Title: "Another Task", Description: "Test"})

	// Find existing task
	task := manager.FindByTitle("Unique Task")
	if task == nil || task.Title != "Unique Task" {
		t.Errorf("Expected to find 'Unique Task', got %v", task)
	}

	// Find non-existing task
	task = manager.FindByTitle("Non-existent")
	if task != nil {
		t.Error("Expected nil for non-existent task")
	}
}

func TestBulkAdd(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_bulk_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Bulk add tasks
	tasks := []Task{
		{Title: "Task1", Description: "First"},
		{Title: "Task2", Description: "Second"},
		{Title: "Task3", Description: "Third"},
	}

	err = manager.BulkAdd(tasks)
	if err != nil {
		t.Error("BulkAdd failed:", err)
	}

	// Verify all tasks were added
	list := manager.List()
	if len(list) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(list))
	}

	// Test empty bulk add
	err = manager.BulkAdd([]Task{})
	if err != nil {
		t.Error("BulkAdd with empty slice should not error:", err)
	}
}

func TestCountDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_countdone_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test with no tasks
	count := manager.CountDone()
	if count != 0 {
		t.Errorf("Expected 0 done tasks, got %d", count)
	}

	// Add mixed tasks
	manager.Add(Task{Title: "Task1", Done: true})
	manager.Add(Task{Title: "Task2", Done: false})
	manager.Add(Task{Title: "Task3", Done: true})

	count = manager.CountDone()
	if count != 2 {
		t.Errorf("Expected 2 done tasks, got %d", count)
	}
}

func TestUndoDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_undo_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add done task
	manager.Add(Task{Title: "Done Task", Done: true})
	manager.Add(Task{Title: "Undone Task", Done: false})

	// Undo done task
	err = manager.UndoDone("0")
	if err != nil {
		t.Error("UndoDone failed:", err)
	}

	list := manager.List()
	if list[0].Done {
		t.Error("Expected task to be undone")
	}

	// Undo already undone task (should not error)
	err = manager.UndoDone("1")
	if err != nil {
		t.Error("UndoDone on already undone task should not error:", err)
	}

	// Test invalid index
	err = manager.UndoDone("10")
	if err == nil {
		t.Error("Expected error for invalid index")
	}
}

func TestFindByDescription(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_finddesc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks with different descriptions
	manager.Add(Task{Title: "Task1", Description: "unique description"})
	manager.Add(Task{Title: "Task2", Description: "another description"})
	manager.Add(Task{Title: "Task3", Description: "unique description"})

	// Find tasks by description
	results := manager.FindByDescription("unique description")
	if len(results) != 2 {
		t.Errorf("Expected 2 tasks with 'unique description', got %d", len(results))
	}

	// Find non-existing description
	results = manager.FindByDescription("non-existent")
	if len(results) != 0 {
		t.Errorf("Expected 0 tasks with non-existent description, got %d", len(results))
	}

	// Test with empty list
	manager2 := NewTaskManager(NewFileStore(filepath.Join(dir, "empty.json")))
	results = manager2.FindByDescription("any")
	if len(results) != 0 {
		t.Errorf("Expected 0 tasks from empty list, got %d", len(results))
	}
}

func TestMarkAllDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_markall_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test with empty list
	err = manager.MarkAllDone()
	if err != nil {
		t.Error("MarkAllDone should not error on empty list:", err)
	}

	// Add mixed tasks
	manager.Add(Task{Title: "Task1", Done: false})
	manager.Add(Task{Title: "Task2", Done: true})
	manager.Add(Task{Title: "Task3", Done: false})

	// Mark all done
	err = manager.MarkAllDone()
	if err != nil {
		t.Error("MarkAllDone failed:", err)
	}

	// Verify all tasks are done
	list := manager.List()
	for i, task := range list {
		if !task.Done {
			t.Errorf("Task %d should be done after MarkAllDone", i)
		}
	}

	// Test marking already done tasks (should not error)
	err = manager.MarkAllDone()
	if err != nil {
		t.Error("MarkAllDone on already done tasks should not error:", err)
	}
}

func TestRemove(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_remove_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks
	manager.Add(Task{Title: "Task1", Description: "First"})
	manager.Add(Task{Title: "Task2", Description: "Second"})
	manager.Add(Task{Title: "Task3", Description: "Third"})

	// Remove middle task
	err = manager.Remove("1")
	if err != nil {
		t.Error("Remove failed:", err)
	}

	// Verify task was removed
	list := manager.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 tasks after removal, got %d", len(list))
	}
	if list[0].Title != "Task1" || list[1].Title != "Task3" {
		t.Error("Wrong tasks remaining after removal")
	}

	// Test invalid index
	err = manager.Remove("10")
	if err == nil {
		t.Error("Expected error for invalid index")
	}

	// Test invalid index string
	err = manager.Remove("invalid")
	if err == nil {
		t.Error("Expected error for invalid index string")
	}
}
