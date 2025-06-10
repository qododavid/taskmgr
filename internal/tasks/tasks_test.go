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

	// Add multiple tasks
	tasks := []Task{
		{Title: "Task 1", Done: true},
		{Title: "Task 2", Done: false},
		{Title: "Task 3", Done: true},
		{Title: "Task 4", Done: false},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	// Clear completed tasks
	if err := manager.ClearCompleted(); err != nil {
		t.Fatalf("ClearCompleted failed: %v", err)
	}

	// Check remaining tasks
	remaining := manager.List()
	if len(remaining) != 2 {
		t.Errorf("Expected 2 remaining tasks, got %d", len(remaining))
	}

	for _, task := range remaining {
		if task.Done {
			t.Errorf("Found completed task that should have been cleared: %v", task)
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

	// Test empty list
	if count := manager.CountUndone(); count != 0 {
		t.Errorf("Expected 0 undone tasks in empty list, got %d", count)
	}

	// Add tasks with mixed completion status
	tasks := []Task{
		{Title: "Task 1", Done: true},
		{Title: "Task 2", Done: false},
		{Title: "Task 3", Done: true},
		{Title: "Task 4", Done: false},
		{Title: "Task 5", Done: false},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	count := manager.CountUndone()
	if count != 3 {
		t.Errorf("Expected 3 undone tasks, got %d", count)
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

	// Add tasks with mixed completion status
	tasks := []Task{
		{Title: "Done Task 1", Done: true},
		{Title: "Undone Task 1", Done: false},
		{Title: "Done Task 2", Done: true},
		{Title: "Undone Task 2", Done: false},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	undone = manager.ListUndone()
	if len(undone) != 2 {
		t.Errorf("Expected 2 undone tasks, got %d", len(undone))
	}

	for _, task := range undone {
		if task.Done {
			t.Errorf("Found completed task in undone list: %v", task)
		}
		if task.Title != "Undone Task 1" && task.Title != "Undone Task 2" {
			t.Errorf("Unexpected task in undone list: %v", task)
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
	task := Task{Title: "Toggle Test", Done: false}
	if err := manager.Add(task); err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Toggle from false to true
	if err := manager.ToggleDone("0"); err != nil {
		t.Fatalf("ToggleDone failed: %v", err)
	}

	tasks := manager.List()
	if !tasks[0].Done {
		t.Error("Expected task to be done after toggle")
	}

	// Toggle from true to false
	if err := manager.ToggleDone("0"); err != nil {
		t.Fatalf("ToggleDone failed: %v", err)
	}

	tasks = manager.List()
	if tasks[0].Done {
		t.Error("Expected task to be undone after second toggle")
	}

	// Test invalid index
	if err := manager.ToggleDone("invalid"); err == nil {
		t.Error("Expected error for invalid index")
	}

	if err := manager.ToggleDone("99"); err == nil {
		t.Error("Expected error for out of range index")
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
		t.Error("Expected empty task list to return true")
	}

	// Add a task
	if err := manager.Add(Task{Title: "Test"}); err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Test non-empty list
	if manager.IsEmpty() {
		t.Error("Expected non-empty task list to return false")
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
	tasks := []Task{
		{Title: "First Task", Description: "First description"},
		{Title: "Second Task", Description: "Second description"},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	// Test valid index
	task, err := manager.GetTask("0")
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if task.Title != "First Task" {
		t.Errorf("Expected 'First Task', got '%s'", task.Title)
	}

	task, err = manager.GetTask("1")
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if task.Title != "Second Task" {
		t.Errorf("Expected 'Second Task', got '%s'", task.Title)
	}

	// Test invalid index
	_, err = manager.GetTask("invalid")
	if err == nil {
		t.Error("Expected error for invalid index")
	}

	_, err = manager.GetTask("99")
	if err == nil {
		t.Error("Expected error for out of range index")
	}

	_, err = manager.GetTask("-1")
	if err == nil {
		t.Error("Expected error for negative index")
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
	task := Task{Title: "Original Title", Description: "Test description"}
	if err := manager.Add(task); err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Update title
	newTitle := "Updated Title"
	if err := manager.UpdateTitle("0", newTitle); err != nil {
		t.Fatalf("UpdateTitle failed: %v", err)
	}

	// Verify update
	tasks := manager.List()
	if tasks[0].Title != newTitle {
		t.Errorf("Expected title '%s', got '%s'", newTitle, tasks[0].Title)
	}
	if tasks[0].Description != "Test description" {
		t.Error("Description should remain unchanged")
	}

	// Test invalid index
	if err := manager.UpdateTitle("invalid", "New Title"); err == nil {
		t.Error("Expected error for invalid index")
	}

	if err := manager.UpdateTitle("99", "New Title"); err == nil {
		t.Error("Expected error for out of range index")
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

	// Test bulk add
	tasks := []Task{
		{Title: "Bulk Task 1"},
		{Title: "Bulk Task 2"},
		{Title: "Bulk Task 3"},
	}

	if err := manager.BulkAdd(tasks); err != nil {
		t.Fatalf("BulkAdd failed: %v", err)
	}

	list := manager.List()
	if len(list) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(list))
	}

	for i, task := range tasks {
		if list[i].Title != task.Title {
			t.Errorf("Expected task %d to have title '%s', got '%s'", i, task.Title, list[i].Title)
		}
	}

	// Test empty bulk add
	if err := manager.BulkAdd([]Task{}); err != nil {
		t.Fatalf("BulkAdd with empty slice failed: %v", err)
	}
}

func TestCountDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_count_done_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test empty list
	if count := manager.CountDone(); count != 0 {
		t.Errorf("Expected 0 done tasks in empty list, got %d", count)
	}

	// Add tasks with mixed completion status
	tasks := []Task{
		{Title: "Task 1", Done: true},
		{Title: "Task 2", Done: false},
		{Title: "Task 3", Done: true},
		{Title: "Task 4", Done: false},
		{Title: "Task 5", Done: true},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	count := manager.CountDone()
	if count != 3 {
		t.Errorf("Expected 3 done tasks, got %d", count)
	}
}

func TestFindByTitle(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_find_title_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks
	tasks := []Task{
		{Title: "Unique Task", Description: "First"},
		{Title: "Common Task", Description: "Second"},
		{Title: "Another Task", Description: "Third"},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	// Test finding existing task
	found := manager.FindByTitle("Unique Task")
	if found == nil {
		t.Error("Expected to find task with title 'Unique Task'")
	} else if found.Description != "First" {
		t.Errorf("Expected description 'First', got '%s'", found.Description)
	}

	// Test finding non-existing task
	notFound := manager.FindByTitle("Non-existent Task")
	if notFound != nil {
		t.Error("Expected nil for non-existent task")
	}
}

func TestFindByDescription(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_find_desc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks
	tasks := []Task{
		{Title: "Task 1", Description: "unique description"},
		{Title: "Task 2", Description: "common description"},
		{Title: "Task 3", Description: "common description"},
		{Title: "Task 4", Description: "another description"},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	// Test finding tasks with common description
	found := manager.FindByDescription("common description")
	if len(found) != 2 {
		t.Errorf("Expected 2 tasks with 'common description', got %d", len(found))
	}

	// Test finding task with unique description
	found = manager.FindByDescription("unique description")
	if len(found) != 1 {
		t.Errorf("Expected 1 task with 'unique description', got %d", len(found))
	}

	// Test finding non-existing description
	found = manager.FindByDescription("non-existent description")
	if len(found) != 0 {
		t.Errorf("Expected 0 tasks with non-existent description, got %d", len(found))
	}
}

func TestMarkAllDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_mark_all_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks with mixed completion status
	tasks := []Task{
		{Title: "Task 1", Done: false},
		{Title: "Task 2", Done: true},
		{Title: "Task 3", Done: false},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	// Mark all done
	if err := manager.MarkAllDone(); err != nil {
		t.Fatalf("MarkAllDone failed: %v", err)
	}

	// Verify all tasks are done
	list := manager.List()
	for i, task := range list {
		if !task.Done {
			t.Errorf("Expected task %d to be done", i)
		}
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

	// Add a done task
	task := Task{Title: "Done Task", Done: true}
	if err := manager.Add(task); err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Undo done
	if err := manager.UndoDone("0"); err != nil {
		t.Fatalf("UndoDone failed: %v", err)
	}

	tasks := manager.List()
	if tasks[0].Done {
		t.Error("Expected task to be undone")
	}

	// Test undoing already undone task (should be no-op)
	if err := manager.UndoDone("0"); err != nil {
		t.Fatalf("UndoDone on already undone task failed: %v", err)
	}

	tasks = manager.List()
	if tasks[0].Done {
		t.Error("Expected task to remain undone")
	}

	// Test invalid index
	if err := manager.UndoDone("invalid"); err == nil {
		t.Error("Expected error for invalid index")
	}

	if err := manager.UndoDone("99"); err == nil {
		t.Error("Expected error for out of range index")
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
	tasks := []Task{
		{Title: "Task 1"},
		{Title: "Task 2"},
		{Title: "Task 3"},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	}

	// Remove middle task
	if err := manager.Remove("1"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	list := manager.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 tasks after removal, got %d", len(list))
	}

	if list[0].Title != "Task 1" || list[1].Title != "Task 3" {
		t.Errorf("Unexpected tasks after removal: %v", list)
	}

	// Test invalid index
	if err := manager.Remove("invalid"); err == nil {
		t.Error("Expected error for invalid index")
	}

	if err := manager.Remove("99"); err == nil {
		t.Error("Expected error for out of range index")
	}
}

func TestMarkDoneErrors(t *testing.T) {
	dir, err := ioutil.TempDir("", "taskmgr_mark_done_errors_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Test invalid index
	if err := manager.MarkDone("invalid"); err == nil {
		t.Error("Expected error for invalid index")
	}

	// Test out of range index
	if err := manager.MarkDone("99"); err == nil {
		t.Error("Expected error for out of range index")
	}

	// Add a task and test negative index
	if err := manager.Add(Task{Title: "Test"}); err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	if err := manager.MarkDone("-1"); err == nil {
		t.Error("Expected error for negative index")
	}
}