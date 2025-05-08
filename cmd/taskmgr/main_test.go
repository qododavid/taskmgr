package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// cleanupTasks removes the tasks.json file before/after tests
func cleanupTasks() {
	os.Remove("tasks.json")
}

func TestMainCLI(t *testing.T) {
	// Clean up before tests
	cleanupTasks()
	// Clean up after tests
	defer cleanupTasks()

	// Test 'add' command
	cmd := exec.Command("go", "run", "./main.go", "add", "TestTask")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run add command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Task added") {
		t.Error("Expected 'Task added.' output")
	}

	// Test 'list' command
	cmd = exec.Command("go", "run", "./main.go", "list")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run list command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "TestTask") {
		t.Error("Expected task to be listed")
	}

	// Test 'find' command
	cmd = exec.Command("go", "run", "./main.go", "find", "TestTask")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run find command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "TestTask") {
		t.Error("Expected to find task")
	}

	// Test 'done' command
	cmd = exec.Command("go", "run", "./main.go", "done", "0")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run done command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "marked as done") {
		t.Error("Expected 'Task marked as done.' output")
	}

	// Test 'undodone' command
	cmd = exec.Command("go", "run", "./main.go", "undodone", "0")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run undodone command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "marked as not done") {
		t.Error("Expected 'Task marked as not done.' output")
	}

	// Test 'bulkadd' command
	cmd = exec.Command("go", "run", "./main.go", "bulkadd", "Task1,Task2,Task3")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run bulkadd command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Tasks added") {
		t.Error("Expected 'Tasks added.' output")
	}

	// Test 'countdone' command
	// First mark a task as done
	cmd = exec.Command("go", "run", "./main.go", "done", "0")
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run done command: %v", err)
	}

	cmd = exec.Command("go", "run", "./main.go", "countdone")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run countdone command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Completed tasks:") {
		t.Error("Expected 'Completed tasks:' output")
	}

	// Test 'markall' command
	cmd = exec.Command("go", "run", "./main.go", "markall")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run markall command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "All tasks marked as done") {
		t.Error("Expected 'All tasks marked as done.' output")
	}

	// Add a task with description for findbydesc test
	cmd = exec.Command("go", "run", "./main.go", "add", "TaskWithDesc")
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}
	
	// Test 'remove' command
	cmd = exec.Command("go", "run", "./main.go", "remove", "1")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run remove command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Task removed") {
		t.Error("Expected 'Task removed.' output")
	}

	// Test invalid command (help message)
	cmd = exec.Command("go", "run", "./main.go", "invalid")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected error for invalid command")
	}
	if !strings.Contains(string(out), "Available commands") {
		t.Error("Expected help message with available commands")
	}
}

// Test commands with missing arguments
func TestMainCLIErrors(t *testing.T) {
	// Clean up before tests
	cleanupTasks()
	// Clean up after tests
	defer cleanupTasks()

	// Test 'add' with missing arguments
	cmd := exec.Command("go", "run", "./main.go", "add")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected error for add without title")
	}
	if !strings.Contains(string(out), "Usage: taskmgr add") {
		t.Error("Expected usage message for add")
	}

	// Test 'done' with missing arguments
	cmd = exec.Command("go", "run", "./main.go", "done")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected error for done without index")
	}
	if !strings.Contains(string(out), "Usage: taskmgr done") {
		t.Error("Expected usage message for done")
	}

	// Test 'remove' with missing arguments
	cmd = exec.Command("go", "run", "./main.go", "remove")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected error for remove without index")
	}
	if !strings.Contains(string(out), "Usage: taskmgr remove") {
		t.Error("Expected usage message for remove")
	}

	// Test 'undodone' with missing arguments
	cmd = exec.Command("go", "run", "./main.go", "undodone")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected error for undodone without index")
	}
	if !strings.Contains(string(out), "Usage: taskmgr undodone") {
		t.Error("Expected usage message for undodone")
	}

	// Test 'find' with missing arguments
	cmd = exec.Command("go", "run", "./main.go", "find")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected error for find without title")
	}
	if !strings.Contains(string(out), "Usage: taskmgr find") {
		t.Error("Expected usage message for find")
	}

	// Test 'bulkadd' with missing arguments
	cmd = exec.Command("go", "run", "./main.go", "bulkadd")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected error for bulkadd without titles")
	}
	if !strings.Contains(string(out), "Usage: taskmgr bulkadd") {
		t.Error("Expected usage message for bulkadd")
	}

	// Test 'findbydesc' with missing arguments
	cmd = exec.Command("go", "run", "./main.go", "findbydesc")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected error for findbydesc without description")
	}
	if !strings.Contains(string(out), "Usage: taskmgr findbydesc") {
		t.Error("Expected usage message for findbydesc")
	}
}