package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// captureOutput captures stdout and stderr for testing
func captureOutput(f func()) (string, string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	
	os.Stdout = wOut
	os.Stderr = wErr

	f()

	wOut.Close()
	wErr.Close()
	
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)
	
	return bufOut.String(), bufErr.String()
}

// runMainWithExec runs the main function via go run to test it properly
func runMainWithExec(args []string) (string, string, int) {
	cmdArgs := append([]string{"run", "./main.go"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = "/home/runner/work/taskmgr/taskmgr/cmd/taskmgr"
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}
	
	return stdout.String(), stderr.String(), exitCode
}

func TestMainAddCommand(t *testing.T) {
	// Clean up test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	stdout, stderr, exitCode := runMainWithExec([]string{"add", "Test Task"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Task added.") {
		t.Errorf("Expected 'Task added.' in output, got: %s", stdout)
	}
}

func TestMainAddCommandMissingArgs(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"add"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr add <title>") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
}

func TestMainListCommand(t *testing.T) {
	// Clean up and setup test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	// Add a task first
	runMainWithExec([]string{"add", "Test Task"})
	
	stdout, stderr, exitCode := runMainWithExec([]string{"list"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Test Task") {
		t.Errorf("Expected task in list output, got: %s", stdout)
	}
}

func TestMainDoneCommand(t *testing.T) {
	// Clean up and setup test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	// Add a task first
	runMainWithExec([]string{"add", "Test Task"})
	
	stdout, stderr, exitCode := runMainWithExec([]string{"done", "0"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Task marked as done.") {
		t.Errorf("Expected success message, got: %s", stdout)
	}
}

func TestMainDoneCommandMissingArgs(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"done"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr done <index>") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
}

func TestMainDoneCommandInvalidIndex(t *testing.T) {
	// Clean up test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	stdout, _, exitCode := runMainWithExec([]string{"done", "999"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Error marking done:") {
		t.Errorf("Expected error message, got: %s", stdout)
	}
}

func TestMainRemoveCommand(t *testing.T) {
	// Clean up and setup test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	// Add a task first
	runMainWithExec([]string{"add", "Test Task"})
	
	stdout, stderr, exitCode := runMainWithExec([]string{"remove", "0"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Task removed.") {
		t.Errorf("Expected success message, got: %s", stdout)
	}
}

func TestMainRemoveCommandMissingArgs(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"remove"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr remove <index>") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
}

func TestMainUndoDoneCommand(t *testing.T) {
	// Clean up and setup test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	// Add and mark a task as done first
	runMainWithExec([]string{"add", "Test Task"})
	runMainWithExec([]string{"done", "0"})
	
	stdout, stderr, exitCode := runMainWithExec([]string{"undodone", "0"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Task marked as not done.") {
		t.Errorf("Expected success message, got: %s", stdout)
	}
}

func TestMainUndoDoneCommandMissingArgs(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"undodone"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr undodone <index>") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
}

func TestMainFindCommand(t *testing.T) {
	// Clean up and setup test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	// Add a task first
	runMainWithExec([]string{"add", "Test Task"})
	
	stdout, stderr, exitCode := runMainWithExec([]string{"find", "Test Task"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Test Task") {
		t.Errorf("Expected task found, got: %s", stdout)
	}
}

func TestMainFindCommandNotFound(t *testing.T) {
	// Clean up test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	stdout, _, exitCode := runMainWithExec([]string{"find", "Nonexistent"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "No task found with that title.") {
		t.Errorf("Expected not found message, got: %s", stdout)
	}
}

func TestMainFindCommandMissingArgs(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"find"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr find <title>") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
}

func TestMainBulkAddCommand(t *testing.T) {
	// Clean up test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	stdout, stderr, exitCode := runMainWithExec([]string{"bulkadd", "Task 1,Task 2,Task 3"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Tasks added.") {
		t.Errorf("Expected success message, got: %s", stdout)
	}
}

func TestMainBulkAddCommandMissingArgs(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"bulkadd"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr bulkadd <title1,title2,...>") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
}

func TestMainCountDoneCommand(t *testing.T) {
	// Clean up and setup test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	// Add and mark some tasks as done
	runMainWithExec([]string{"add", "Task 1"})
	runMainWithExec([]string{"add", "Task 2"})
	runMainWithExec([]string{"done", "0"})
	
	stdout, stderr, exitCode := runMainWithExec([]string{"countdone"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Completed tasks: 1") {
		t.Errorf("Expected count message, got: %s", stdout)
	}
}

func TestMainMarkAllCommand(t *testing.T) {
	// Clean up and setup test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	// Add some tasks
	runMainWithExec([]string{"add", "Task 1"})
	runMainWithExec([]string{"add", "Task 2"})
	
	stdout, stderr, exitCode := runMainWithExec([]string{"markall"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "All tasks marked as done.") {
		t.Errorf("Expected success message, got: %s", stdout)
	}
}

func TestMainFindByDescCommand(t *testing.T) {
	// Clean up test file
	os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	defer os.Remove("/home/runner/work/taskmgr/taskmgr/cmd/taskmgr/tasks.json")
	
	stdout, _, exitCode := runMainWithExec([]string{"findbydesc", "test"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "No tasks found with that description.") {
		t.Errorf("Expected not found message, got: %s", stdout)
	}
}

func TestMainFindByDescCommandMissingArgs(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"findbydesc"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr findbydesc <description>") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
}

func TestMainErrorCommand(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"error"})
	
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	// Error command doesn't produce output, just triggers sentry
	_ = stdout
}

func TestMainDefaultCommand(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{"unknown"})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr [command] ...") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Available commands:") {
		t.Errorf("Expected available commands, got: %s", stdout)
	}
}

func TestMainEmptyArgs(t *testing.T) {
	stdout, _, exitCode := runMainWithExec([]string{})
	
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: taskmgr [command] ...") {
		t.Errorf("Expected usage message, got: %s", stdout)
	}
}