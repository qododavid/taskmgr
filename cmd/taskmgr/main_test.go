package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestEnvironment creates a clean test environment with a tasks file
func setupTestEnvironment(t *testing.T) string {
	dir, err := ioutil.TempDir("", "taskmgr_cli_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	// Change to the directory to simplify relative paths
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	
	// Return the original directory so we can change back in cleanup
	return oldDir
}

// cleanupTestEnvironment restores the original directory and cleans up temp files
func cleanupTestEnvironment(t *testing.T, originalDir string, tempDir string) {
	os.Chdir(originalDir)
	os.RemoveAll(tempDir)
}

// runCommand executes the taskmgr command with the given arguments
func runCommand(t *testing.T, args ...string) (string, error) {
	fullArgs := append([]string{"run", "."}, args...)
	cmd := exec.Command("go", fullArgs...)
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

func TestCliCommands(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	
	// Create temp directory and set it up
	dir, err := ioutil.TempDir("", "taskmgr_commands_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	
	// Copy main.go to temp directory for testing
	mainSrc, err := ioutil.ReadFile("main.go")
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}
	
	mainDest := filepath.Join(dir, "main.go")
	err = ioutil.WriteFile(mainDest, mainSrc, 0644)
	if err != nil {
		t.Fatalf("Failed to write main.go to temp dir: %v", err)
	}
	
	// Change to temp directory for tests
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	// Test the bulkadd command
	t.Run("TestBulkAdd", func(t *testing.T) {
		out, err := exec.Command("go", "run", "main.go", "bulkadd", "Task1,Task2,Task3").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run bulkadd command: %v (%s)", err, string(out))
		}
		if !strings.Contains(string(out), "Tasks added") {
			t.Errorf("Expected 'Tasks added.' output, got: %s", string(out))
		}
		
		// Verify using the list command
		out, err = exec.Command("go", "run", "main.go", "list").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run list command: %v (%s)", err, string(out))
		}
		
		output := string(out)
		for _, taskName := range []string{"Task1", "Task2", "Task3"} {
			if !strings.Contains(output, taskName) {
				t.Errorf("Expected task '%s' in list output, got: %s", taskName, output)
			}
		}
	})
	
	// Test the done command
	t.Run("TestDone", func(t *testing.T) {
		out, err := exec.Command("go", "run", "main.go", "done", "0").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run done command: %v (%s)", err, string(out))
		}
		if !strings.Contains(string(out), "Task marked as done") {
			t.Errorf("Expected 'Task marked as done.' output, got: %s", string(out))
		}
		
		// Verify using the list command
		out, err = exec.Command("go", "run", "main.go", "list").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run list command: %v (%s)", err, string(out))
		}
		
		if !strings.Contains(string(out), "[x] Task1") {
			t.Errorf("Expected '[x] Task1' in list output, got: %s", string(out))
		}
	})
	
	// Test the undodone command
	t.Run("TestUndoDone", func(t *testing.T) {
		out, err := exec.Command("go", "run", "main.go", "undodone", "0").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run undodone command: %v (%s)", err, string(out))
		}
		if !strings.Contains(string(out), "Task marked as not done") {
			t.Errorf("Expected 'Task marked as not done.' output, got: %s", string(out))
		}
		
		// Verify using the list command
		out, err = exec.Command("go", "run", "main.go", "list").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run list command: %v (%s)", err, string(out))
		}
		
		if !strings.Contains(string(out), "[ ] Task1") {
			t.Errorf("Expected '[ ] Task1' in list output, got: %s", string(out))
		}
	})
	
	// Test the find command
	t.Run("TestFind", func(t *testing.T) {
		out, err := exec.Command("go", "run", "main.go", "find", "Task2").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run find command: %v (%s)", err, string(out))
		}
		if !strings.Contains(string(out), "[ ] Task2") {
			t.Errorf("Expected '[ ] Task2' in find output, got: %s", string(out))
		}
	})
	
	// Test the countdone command
	t.Run("TestCountDone", func(t *testing.T) {
		// First mark a task as done
		_, err := exec.Command("go", "run", "main.go", "done", "1").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run done command: %v", err)
		}
		
		out, err := exec.Command("go", "run", "main.go", "countdone").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run countdone command: %v (%s)", err, string(out))
		}
		if !strings.Contains(string(out), "Completed tasks: 1") {
			t.Errorf("Expected 'Completed tasks: 1' in output, got: %s", string(out))
		}
	})
	
	// Test the markall command
	t.Run("TestMarkAll", func(t *testing.T) {
		out, err := exec.Command("go", "run", "main.go", "markall").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run markall command: %v (%s)", err, string(out))
		}
		if !strings.Contains(string(out), "All tasks marked as done") {
			t.Errorf("Expected 'All tasks marked as done.' output, got: %s", string(out))
		}
		
		// Verify all tasks are done
		out, err = exec.Command("go", "run", "main.go", "list").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run list command: %v (%s)", err, string(out))
		}
		
		output := string(out)
		if strings.Contains(output, "[ ]") { // Check if any task is not done
			t.Errorf("Expected all tasks to be marked done, got: %s", output)
		}
	})
	
	// Test the findbydesc command
	t.Run("TestFindByDesc", func(t *testing.T) {
		// First add a task with a description
		cmd := exec.Command("go", "run", "main.go", "add", "TaskWithDesc")
		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
		
		// Since there's no direct way to add a description via CLI, we need to
		// modify the tasks.json file directly to add a description
		
		// For now, we'll just test the command with a description that won't match
		out, err := exec.Command("go", "run", "main.go", "findbydesc", "no matching description").CombinedOutput()
		if err != nil {
			// This is expected since there are no tasks with this description
			if !strings.Contains(string(out), "No tasks found with that description") {
				t.Errorf("Expected 'No tasks found' message, got: %s", string(out))
			}
		}
	})
	
	// Test the remove command - do this last
	t.Run("TestRemove", func(t *testing.T) {
		// Get initial count
		out, err := exec.Command("go", "run", "main.go", "list").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run list command: %v (%s)", err, string(out))
		}
		initialTaskCount := strings.Count(string(out), "\n")
		
		// Remove a task
		out, err = exec.Command("go", "run", "main.go", "remove", "0").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run remove command: %v (%s)", err, string(out))
		}
		if !strings.Contains(string(out), "Task removed") {
			t.Errorf("Expected 'Task removed.' output, got: %s", string(out))
		}
		
		// Verify one less task
		out, err = exec.Command("go", "run", "main.go", "list").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run list command: %v (%s)", err, string(out))
		}
		newTaskCount := strings.Count(string(out), "\n")
		
		if newTaskCount >= initialTaskCount {
			t.Errorf("Expected task count to decrease after removal")
		}
	})
	
	// Test invalid commands and error handling
	t.Run("TestInvalidCommands", func(t *testing.T) {
		// Test done with invalid index
		out, err := exec.Command("go", "run", "main.go", "done", "999").CombinedOutput()
		if err == nil {
			t.Errorf("Expected error for invalid index, got success")
		}
		if !strings.Contains(string(out), "Error marking") {
			t.Errorf("Expected error message for invalid index, got: %s", string(out))
		}
		
		// Test remove with invalid index
		out, err = exec.Command("go", "run", "main.go", "remove", "999").CombinedOutput()
		if err == nil {
			t.Errorf("Expected error for invalid remove index, got success")
		}
		if !strings.Contains(string(out), "Error removing") {
			t.Errorf("Expected error message for invalid remove index, got: %s", string(out))
		}
		
		// Test find with non-existent title
		out, err = exec.Command("go", "run", "main.go", "find", "NonExistentTask").CombinedOutput()
		if err != nil {
			t.Fatalf("Find with non-existent task should not return error: %v", err)
		}
		if !strings.Contains(string(out), "No task found") {
			t.Errorf("Expected 'No task found' message, got: %s", string(out))
		}
	})
}