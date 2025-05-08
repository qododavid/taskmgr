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
	
	// Test UndoDone
	err = manager.UndoDone("0")
	if err != nil {
		t.Errorf("UndoDone returned an error: %v", err)
	}
	
	list = manager.List()
	if list[0].Done {
		t.Error("Expected task to be marked as not done after UndoDone")
	}
	
	// Test FindByTitle - positive case
	task := manager.FindByTitle("First")
	if task == nil {
		t.Error("FindByTitle returned nil for existing task")
	} else if task.Title != "First" {
		t.Errorf("FindByTitle returned wrong task, got %v", task)
	}
	
	// Test FindByTitle - negative case
	task = manager.FindByTitle("Nonexistent")
	if task != nil {
		t.Errorf("FindByTitle returned task for nonexistent title: %v", task)
	}
	
	// Test BulkAdd
	tasksToAdd := []Task{
		{Title: "Task1"},
		{Title: "Task2"},
		{Title: "Task3"},
	}
	
	err = manager.BulkAdd(tasksToAdd)
	if err != nil {
		t.Errorf("BulkAdd returned an error: %v", err)
	}
	
	list = manager.List()
	if len(list) != 4 { // 1 original + 3 added
		t.Errorf("Expected 4 tasks after BulkAdd, got %d", len(list))
	}
	
	// Test CountDone
	// First, mark some tasks as done
	err = manager.MarkDone("1")
	if err != nil {
		t.Errorf("MarkDone returned an error: %v", err)
	}
	
	err = manager.MarkDone("2")
	if err != nil {
		t.Errorf("MarkDone returned an error: %v", err)
	}
	
	count := manager.CountDone()
	if count != 2 {
		t.Errorf("Expected CountDone to return 2, got %d", count)
	}
	
	// Test MarkAllDone
	err = manager.MarkAllDone()
	if err != nil {
		t.Errorf("MarkAllDone returned an error: %v", err)
	}
	
	list = manager.List()
	for i, task := range list {
		if !task.Done {
			t.Errorf("Task at index %d not marked done after MarkAllDone", i)
		}
	}
	
	// Test FindByDescription
	// Add a task with a description
	taskWithDesc := Task{
		Title:       "WithDesc",
		Description: "Test Description",
	}
	
	err = manager.Add(taskWithDesc)
	if err != nil {
		t.Errorf("Error adding task with description: %v", err)
	}
	
	// Positive test
	results := manager.FindByDescription("Test Description")
	if len(results) != 1 {
		t.Errorf("Expected 1 result from FindByDescription, got %d", len(results))
	} else if results[0].Title != "WithDesc" {
		t.Errorf("FindByDescription returned wrong task, got %v", results[0])
	}
	
	// Negative test
	results = manager.FindByDescription("Nonexistent Description")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for nonexistent description, got %d", len(results))
	}
	
	// Test Remove
	initialCount := len(manager.List())
	err = manager.Remove("0")
	if err != nil {
		t.Errorf("Remove returned an error: %v", err)
	}
	
	afterRemoveCount := len(manager.List())
	if afterRemoveCount != initialCount-1 {
		t.Errorf("Expected task count to be %d after Remove, got %d", initialCount-1, afterRemoveCount)
	}
	
	// Test error cases
	// Invalid index for MarkDone
	err = manager.MarkDone("999")
	if err == nil {
		t.Error("Expected error for invalid index in MarkDone, got nil")
	}
	
	// Invalid index for UndoDone
	err = manager.UndoDone("999")
	if err == nil {
		t.Error("Expected error for invalid index in UndoDone, got nil")
	}
	
	// Invalid index for Remove
	err = manager.Remove("999")
	if err == nil {
		t.Error("Expected error for invalid index in Remove, got nil")
	}
}