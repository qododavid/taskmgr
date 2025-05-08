package main

import (
	"os/exec"
	"strings"
	"testing"
)

// Existing test
func TestMainCLI_AddCommand(t *testing.T) {
	// Test 'add' command
	cmd := exec.Command("go", "run", "./main.go", "add", "TestTask")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run add command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Task added") {
		t.Error("Expected 'Task added.' output")
	}
}

// New tests for added commands

func TestMainCLI_RemoveCommand(t *testing.T) {
	// First add a task to remove
	addCmd := exec.Command("go", "run", "./main.go", "add", "TaskToRemove")
	_, err := addCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to setup test by adding task: %v", err)
	}

	// Get the index of the task (assuming it's the only one or the last one)
	listCmd := exec.Command("go", "run", "./main.go", "list")
	listOut, err := listCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	// Now remove the task
	removeCmd := exec.Command("go", "run", "./main.go", "remove", "0")
	removeOut, err := removeCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run remove command: %v (%s)", err, string(removeOut))
	}
	if !strings.Contains(string(removeOut), "Task removed") {
		t.Error("Expected 'Task removed.' output")
	}
}

func TestMainCLI_UndoDoneCommand(t *testing.T) {
	// Add a task
	addCmd := exec.Command("go", "run", "./main.go", "add", "UndoDoneTask")
	_, err := addCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to setup test by adding task: %v", err)
	}

	// Mark it as done
	doneCmd := exec.Command("go", "run", "./main.go", "done", "0")
	_, err = doneCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to mark task as done: %v", err)
	}

	// Now undo the done status
	undoCmd := exec.Command("go", "run", "./main.go", "undodone", "0")
	undoOut, err := undoCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run undodone command: %v (%s)", err, string(undoOut))
	}
	if !strings.Contains(string(undoOut), "Task marked as not done") {
		t.Error("Expected 'Task marked as not done.' output")
	}
}

func TestMainCLI_FindCommand(t *testing.T) {
	// Add a task with a unique name
	uniqueTitle := "UniqueTaskToFind"
	addCmd := exec.Command("go", "run", "./main.go", "add", uniqueTitle)
	_, err := addCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to setup test by adding task: %v", err)
	}

	// Find the task by title
	findCmd := exec.Command("go", "run", "./main.go", "find", uniqueTitle)
	findOut, err := findCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run find command: %v (%s)", err, string(findOut))
	}
	if !strings.Contains(string(findOut), uniqueTitle) {
		t.Errorf("Expected to find '%s' in output: %s", uniqueTitle, string(findOut))
	}

	// Test not found case
	notFoundCmd := exec.Command("go", "run", "./main.go", "find", "NonExistentTask")
	notFoundOut, err := notFoundCmd.CombinedOutput()
	if err != nil {
		// This is expected to fail with exit code 1
		if !strings.Contains(string(notFoundOut), "No task found") {
			t.Errorf("Expected 'No task found' in output: %s", string(notFoundOut))
		}
	}
}

func TestMainCLI_BulkAddCommand(t *testing.T) {
	// Test bulk add with multiple tasks
	bulkCmd := exec.Command("go", "run", "./main.go", "bulkadd", "Task1,Task2,Task3")
	bulkOut, err := bulkCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run bulkadd command: %v (%s)", err, string(bulkOut))
	}
	if !strings.Contains(string(bulkOut), "Tasks added") {
		t.Error("Expected 'Tasks added.' output")
	}

	// Verify tasks were added by listing
	listCmd := exec.Command("go", "run", "./main.go", "list")
	listOut, err := listCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v (%s)", err, string(listOut))
	}
	
	// Check that all bulk added tasks exist in the output
	output := string(listOut)
	if !strings.Contains(output, "Task1") || !strings.Contains(output, "Task2") || !strings.Contains(output, "Task3") {
		t.Errorf("Not all bulk-added tasks were found in list output: %s", output)
	}
}

func TestMainCLI_CountDoneCommand(t *testing.T) {
	// Add a task
	addCmd := exec.Command("go", "run", "./main.go", "add", "CountTask")
	_, err := addCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to setup test by adding task: %v", err)
	}

	// Mark it as done
	doneCmd := exec.Command("go", "run", "./main.go", "done", "0")
	_, err = doneCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to mark task as done: %v", err)
	}

	// Now count done tasks
	countCmd := exec.Command("go", "run", "./main.go", "countdone")
	countOut, err := countCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run countdone command: %v (%s)", err, string(countOut))
	}
	if !strings.Contains(string(countOut), "Completed tasks:") {
		t.Error("Expected 'Completed tasks:' output")
	}
}

func TestMainCLI_MarkAllCommand(t *testing.T) {
	// Add a few tasks
	addCmd1 := exec.Command("go", "run", "./main.go", "add", "AllTask1")
	_, err := addCmd1.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to setup test by adding task 1: %v", err)
	}
	
	addCmd2 := exec.Command("go", "run", "./main.go", "add", "AllTask2")
	_, err = addCmd2.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to setup test by adding task 2: %v", err)
	}

	// Mark all as done
	markAllCmd := exec.Command("go", "run", "./main.go", "markall")
	markAllOut, err := markAllCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run markall command: %v (%s)", err, string(markAllOut))
	}
	if !strings.Contains(string(markAllOut), "All tasks marked as done") {
		t.Error("Expected 'All tasks marked as done.' output")
	}

	// Verify all tasks are marked done
	listCmd := exec.Command("go", "run", "./main.go", "list")
	listOut, err := listCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}
	
	// Count occurrences of "[x]" in the output which indicates done tasks
	if !strings.Contains(string(listOut), "[x]") {
		t.Error("Expected to see marked tasks [x] in the output")
	}
}

func TestMainCLI_FindByDescCommand(t *testing.T) {
	// This is a bit tricky to test since we don't have a way to add tasks with descriptions
	// We would need to modify the store directly, but for now we'll test the negative case
	
	findDescCmd := exec.Command("go", "run", "./main.go", "findbydesc", "NonExistentDescription")
	findDescOut, err := findDescCmd.CombinedOutput()
	if err != nil {
		// This should exit with code 0 even when nothing is found
		t.Fatalf("Find by description should not error out: %v (%s)", err, string(findDescOut))
	}
	if !strings.Contains(string(findDescOut), "No tasks found") {
		t.Error("Expected 'No tasks found' output")
	}
}

func TestMainCLI_UsageOutput(t *testing.T) {
	// Test the default case with no recognized command
	usageCmd := exec.Command("go", "run", "./main.go", "unknowncommand")
	usageOut, _ := usageCmd.CombinedOutput()
	// We expect this to fail with exit code 1, so we don't check error
	
	output := string(usageOut)
	// Check that all command descriptions are present
	expectedPhrases := []string{
		"Usage:", "Available commands:",
		"add", "list", "done", "remove", "undodone", 
		"find", "bulkadd", "countdone", "markall", "findbydesc",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected '%s' in usage output", phrase)
		}
	}
}