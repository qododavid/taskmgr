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

func TestTaskManagerBulkAdd(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_bulk_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	tasks := []Task{
		{Title: "Task1"},
		{Title: "Task2"},
		{Title: "Task3"},
	}

	err = manager.BulkAdd(tasks)
	if err != nil {
		t.Errorf("BulkAdd returned error: %v", err)
	}

	list := manager.List()
	if len(list) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(list))
	}

	for i, expected := range []string{"Task1", "Task2", "Task3"} {
		if list[i].Title != expected {
			t.Errorf("Expected task %d to be '%s', got '%s'", i, expected, list[i].Title)
		}
	}
}

func TestTaskManagerCountDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_count_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks
	manager.Add(Task{Title: "Task1"})
	manager.Add(Task{Title: "Task2"})
	manager.Add(Task{Title: "Task3"})

	// Initially no tasks done
	count := manager.CountDone()
	if count != 0 {
		t.Errorf("Expected 0 done tasks, got %d", count)
	}

	// Mark some tasks done
	manager.MarkDone("0")
	manager.MarkDone("2")

	count = manager.CountDone()
	if count != 2 {
		t.Errorf("Expected 2 done tasks, got %d", count)
	}
}

func TestTaskManagerFindByDescription(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_finddesc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks with descriptions
	manager.Add(Task{Title: "Task1", Description: "urgent"})
	manager.Add(Task{Title: "Task2", Description: "normal"})
	manager.Add(Task{Title: "Task3", Description: "urgent"})

	// Find by description
	results := manager.FindByDescription("urgent")
	if len(results) != 2 {
		t.Errorf("Expected 2 urgent tasks, got %d", len(results))
	}

	results = manager.FindByDescription("normal")
	if len(results) != 1 {
		t.Errorf("Expected 1 normal task, got %d", len(results))
	}

	results = manager.FindByDescription("nonexistent")
	if len(results) != 0 {
		t.Errorf("Expected 0 nonexistent tasks, got %d", len(results))
	}
}

func TestTaskManagerMarkAllDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_markall_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks
	manager.Add(Task{Title: "Task1"})
	manager.Add(Task{Title: "Task2"})
	manager.Add(Task{Title: "Task3"})

	// Mark all done
	err = manager.MarkAllDone()
	if err != nil {
		t.Errorf("MarkAllDone returned error: %v", err)
	}

	list := manager.List()
	for i, task := range list {
		if !task.Done {
			t.Errorf("Expected task %d to be done, but it wasn't", i)
		}
	}

	count := manager.CountDone()
	if count != 3 {
		t.Errorf("Expected 3 done tasks, got %d", count)
	}
}

func TestTaskManagerUndoDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_undo_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add and mark task done
	manager.Add(Task{Title: "Task1"})
	manager.MarkDone("0")

	list := manager.List()
	if !list[0].Done {
		t.Error("Expected task to be done")
	}

	// Undo done
	err = manager.UndoDone("0")
	if err != nil {
		t.Errorf("UndoDone returned error: %v", err)
	}

	list = manager.List()
	if list[0].Done {
		t.Error("Expected task to be undone")
	}

	// Test undoing already undone task
	err = manager.UndoDone("0")
	if err != nil {
		t.Errorf("UndoDone on already undone task returned error: %v", err)
	}

	// Test invalid index
	err = manager.UndoDone("5")
	if err == nil {
		t.Error("Expected error for invalid index")
	}

	err = manager.UndoDone("invalid")
	if err == nil {
		t.Error("Expected error for invalid index format")
	}
}

func TestTaskManagerRemove(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_remove_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks
	manager.Add(Task{Title: "Task1"})
	manager.Add(Task{Title: "Task2"})
	manager.Add(Task{Title: "Task3"})

	// Remove middle task
	err = manager.Remove("1")
	if err != nil {
		t.Errorf("Remove returned error: %v", err)
	}

	list := manager.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 tasks after removal, got %d", len(list))
	}
	if list[0].Title != "Task1" || list[1].Title != "Task3" {
		t.Errorf("Expected tasks 'Task1' and 'Task3', got %v", list)
	}

	// Test invalid index
	err = manager.Remove("5")
	if err == nil {
		t.Error("Expected error for invalid index")
	}

	err = manager.Remove("invalid")
	if err == nil {
		t.Error("Expected error for invalid index format")
	}
}

func TestTaskManagerFindByTitle(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_findtitle_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks
	manager.Add(Task{Title: "Unique Task"})
	manager.Add(Task{Title: "Another Task"})

	// Find existing task
	task := manager.FindByTitle("Unique Task")
	if task == nil {
		t.Error("Expected to find task, got nil")
	} else if task.Title != "Unique Task" {
		t.Errorf("Expected 'Unique Task', got '%s'", task.Title)
	}

	// Find non-existing task
	task = manager.FindByTitle("Nonexistent")
	if task != nil {
		t.Errorf("Expected nil for nonexistent task, got %v", task)
	}
}

func TestTaskManagerMarkDoneInvalidIndex(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_markdone_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add one task
	manager.Add(Task{Title: "Task1"})

	// Test invalid index
	err = manager.MarkDone("5")
	if err == nil {
		t.Error("Expected error for invalid index")
	}

	err = manager.MarkDone("invalid")
	if err == nil {
		t.Error("Expected error for invalid index format")
	}
}
