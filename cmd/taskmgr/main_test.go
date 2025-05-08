package main

import (
	"os/exec"
	"strings"
	"testing"
	"os"
	"io/ioutil"
	"path/filepath"
)

// Helper function to clean up test environment
func setupTestEnvironment(t *testing.T) (string, func()) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "taskmgr-cli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Return a cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// Helper function to run CLI commands with a specific tasks file
func runCommand(t *testing.T, taskFile string, args ...string) (string, error) {
	cmdArgs := []string{"run", "./main.go"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("go", cmdArgs...)
	// Set the tasks file path in the environment
	cmd.Env = append(os.Environ(), "TASKMGR_FILE="+taskFile)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestMainCLI(t *testing.T) {
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

func TestNewCLICommands(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tasksFile := filepath.Join(tempDir, "tasks.json")

	// Test 'add' to setup some tasks first
	out, err := runCommand(t, tasksFile, "add", "Task1")
	if err != nil {
		t.Fatalf("Failed to add task: %v (%s)", err, out)
	}
	out, err = runCommand(t, tasksFile, "add", "Task2")
	if err != nil {
		t.Fatalf("Failed to add second task: %v (%s)", err, out)
	}

	// Test 'done' command to mark a task as done
	out, err = runCommand(t, tasksFile, "done", "0")
	if err != nil {
		t.Fatalf("Failed to mark task as done: %v (%s)", err, out)
	}
	if !strings.Contains(out, "Task marked as done") {
		t.Errorf("Expected 'Task marked as done' output, got: %s", out)
	}

	// Test 'undodone' command
	out, err = runCommand(t, tasksFile, "undodone", "0")
	if err != nil {
		t.Fatalf("Failed to undo task done: %v (%s)", err, out)
	}
	if !strings.Contains(out, "Task marked as not done") {
		t.Errorf("Expected 'Task marked as not done' output, got: %s", out)
	}

	// Test 'find' command
	out, err = runCommand(t, tasksFile, "find", "Task1")
	if err != nil {
		t.Fatalf("Failed to find task: %v (%s)", err, out)
	}
	if !strings.Contains(out, "Task1") {
		t.Errorf("Expected to find 'Task1', got: %s", out)
	}

	// Test 'bulkadd' command
	out, err = runCommand(t, tasksFile, "bulkadd", "Task3,Task4,Task5")
	if err != nil {
		t.Fatalf("Failed to bulk add tasks: %v (%s)", err, out)
	}
	if !strings.Contains(out, "Tasks added") {
		t.Errorf("Expected 'Tasks added' output, got: %s", out)
	}

	// Test 'list' to verify tasks were added
	out, err = runCommand(t, tasksFile, "list")
	if err != nil {
		t.Fatalf("Failed to list tasks: %v (%s)", err, out)
	}
	for _, task := range []string{"Task1", "Task2", "Task3", "Task4", "Task5"} {
		if !strings.Contains(out, task) {
			t.Errorf("Expected list to contain '%s', got: %s", task, out)
		}
	}

	// Test 'countdone' command
	out, err = runCommand(t, tasksFile, "countdone")
	if err != nil {
		t.Fatalf("Failed to count done tasks: %v (%s)", err, out)
	}
	if !strings.Contains(out, "Completed tasks: 0") {
		t.Errorf("Expected 'Completed tasks: 0' output, got: %s", out)
	}

	// Test 'markall' command
	out, err = runCommand(t, tasksFile, "markall")
	if err != nil {
		t.Fatalf("Failed to mark all tasks as done: %v (%s)", err, out)
	}
	if !strings.Contains(out, "All tasks marked as done") {
		t.Errorf("Expected 'All tasks marked as done' output, got: %s", out)
	}

	// Test 'countdone' again to verify all tasks are marked as done
	out, err = runCommand(t, tasksFile, "countdone")
	if err != nil {
		t.Fatalf("Failed to count done tasks: %v (%s)", err, out)
	}
	if !strings.Contains(out, "Completed tasks: 5") {
		t.Errorf("Expected 'Completed tasks: 5' output, got: %s", out)
	}

	// Test 'remove' command
	out, err = runCommand(t, tasksFile, "remove", "0")
	if err != nil {
		t.Fatalf("Failed to remove task: %v (%s)", err, out)
	}
	if !strings.Contains(out, "Task removed") {
		t.Errorf("Expected 'Task removed' output, got: %s", out)
	}

	// Create a separate test for testing FindByDescription since it needs special setup
}

func TestErrorHandling(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tasksFile := filepath.Join(tempDir, "tasks.json")

	// Test missing arguments for various commands
	testCases := []struct {
		name    string
		args    []string
		errorMsg string
	}{
		{"add missing arg", []string{"add"}, "Usage: taskmgr add"},
		{"done missing arg", []string{"done"}, "Usage: taskmgr done"},
		{"remove missing arg", []string{"remove"}, "Usage: taskmgr remove"},
		{"undodone missing arg", []string{"undodone"}, "Usage: taskmgr undodone"},
		{"find missing arg", []string{"find"}, "Usage: taskmgr find"},
		{"bulkadd missing arg", []string{"bulkadd"}, "Usage: taskmgr bulkadd"},
		{"findbydesc missing arg", []string{"findbydesc"}, "Usage: taskmgr findbydesc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, _ := runCommand(t, tasksFile, tc.args...)
			if !strings.Contains(out, tc.errorMsg) {
				t.Errorf("Expected error message containing '%s', got: %s", tc.errorMsg, out)
			}
		})
	}

	// Test invalid index
	// First add a task
	_, err := runCommand(t, tasksFile, "add", "TestTask")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	invalidIndexCases := []struct {
		name    string
		args    []string
		errorMsg string
	}{
		{"done invalid index", []string{"done", "999"}, "Error marking done"},
		{"remove invalid index", []string{"remove", "999"}, "Error removing task"},
		{"undodone invalid index", []string{"undodone", "999"}, "Error undoing done"},
	}

	for _, tc := range invalidIndexCases {
		t.Run(tc.name, func(t *testing.T) {
			out, _ := runCommand(t, tasksFile, tc.args...)
			if !strings.Contains(out, tc.errorMsg) {
				t.Errorf("Expected error message containing '%s', got: %s", tc.errorMsg, out)
			}
		})
	}
}

func TestInvalidCommand(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tasksFile := filepath.Join(tempDir, "tasks.json")

	// Test with an invalid command
	out, _ := runCommand(t, tasksFile, "invalidcommand")
	
	// Check that the usage information is shown
	if !strings.Contains(out, "Usage: taskmgr [command]") {
		t.Errorf("Expected usage information for invalid command, got: %s", out)
	}

	// Check that all available commands are listed
	expectedCommands := []string{
		"add", "list", "done", "remove", "undodone", 
		"find", "bulkadd", "countdone", "markall", "findbydesc",
	}

	for _, cmd := range expectedCommands {
		if !strings.Contains(out, cmd) {
			t.Errorf("Expected command '%s' to be listed in usage, got: %s", cmd, out)
		}
	}
}

func TestTaskNotFoundScenarios(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tasksFile := filepath.Join(tempDir, "tasks.json")

	// Add a single task for testing
	out, err := runCommand(t, tasksFile, "add", "ExistingTask")
	if err != nil {
		t.Fatalf("Failed to add task: %v (%s)", err, out)
	}

	// Test 'find' with a non-existent title
	out, err = runCommand(t, tasksFile, "find", "NonExistentTask")
	if err != nil {
		t.Fatalf("The find command should not return an error for non-existent tasks: %v (%s)", err, out)
	}
	if !strings.Contains(out, "No task found with that title") {
		t.Errorf("Expected 'No task found' message, got: %s", out)
	}

	// Test 'findbydesc' with a non-existent description
	out, err = runCommand(t, tasksFile, "findbydesc", "NonExistentDescription")
	if err != nil {
		t.Fatalf("The findbydesc command should not return an error for non-existent descriptions: %v (%s)", err, out)
	}
	if !strings.Contains(out, "No tasks found with that description") {
		t.Errorf("Expected 'No tasks found' message, got: %s", out)
	}
}

func TestFindByDescription(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tasksFile := filepath.Join(tempDir, "tasks.json")

	// Create a task file with a task that has a description
	task := `[
	{
		"Title": "TaskWithDesc",
		"Description": "TestDescription",
		"Done": false
	}
]`

	err := ioutil.WriteFile(tasksFile, []byte(task), 0644)
	if err != nil {
		t.Fatalf("Failed to create task file: %v", err)
	}

	// Test 'findbydesc' command
	out, err := runCommand(t, tasksFile, "findbydesc", "TestDescription")
	if err != nil {
		t.Fatalf("Failed to find task by description: %v (%s)", err, out)
	}
	if !strings.Contains(out, "TaskWithDesc") {
		t.Errorf("Expected to find task with description, got: %s", out)
	}
}
