package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestMainCLIList(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Add a task first
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "ListTestTask")
	cmd.Run()

	// Test 'list' command
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run list command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "ListTestTask") {
		t.Error("Expected task to be listed")
	}
}

func TestMainCLIDone(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Add a task first
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "DoneTestTask")
	cmd.Run()

	// Test 'done' command
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "done", "0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run done command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Task marked as done") {
		t.Error("Expected 'Task marked as done.' output")
	}
}

func TestMainCLIRemove(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Add a task first
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "RemoveTestTask")
	cmd.Run()

	// Test 'remove' command
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "remove", "0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run remove command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Task removed") {
		t.Error("Expected 'Task removed.' output")
	}
}

func TestMainCLIUndoDone(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Add and mark task done first
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "UndoTestTask")
	cmd.Run()
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "done", "0")
	cmd.Run()

	// Test 'undodone' command
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "undodone", "0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run undodone command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Task marked as not done") {
		t.Error("Expected 'Task marked as not done.' output")
	}
}

func TestMainCLIFind(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Add a task first
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "FindTestTask")
	cmd.Run()

	// Test 'find' command
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "find", "FindTestTask")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run find command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "FindTestTask") {
		t.Error("Expected to find the task")
	}
}

func TestMainCLIFindNotFound(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Test 'find' command for non-existent task
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "find", "NonExistentTask")
	out, err := cmd.CombinedOutput()
	// Note: exit code 0 is expected for "not found" case
	if !strings.Contains(string(out), "No task found") {
		t.Error("Expected 'No task found' message")
	}
}

func TestMainCLIBulkAdd(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Test 'bulkadd' command
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "bulkadd", "Task1,Task2,Task3")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run bulkadd command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Tasks added") {
		t.Error("Expected 'Tasks added.' output")
	}
}

func TestMainCLICountDone(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Add and mark some tasks done
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "CountTask1")
	cmd.Run()
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "CountTask2")
	cmd.Run()
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "done", "0")
	cmd.Run()

	// Test 'countdone' command
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "countdone")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run countdone command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "Completed tasks: 1") {
		t.Error("Expected 'Completed tasks: 1' output")
	}
}

func TestMainCLIMarkAll(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Add some tasks
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "MarkAllTask1")
	cmd.Run()
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add", "MarkAllTask2")
	cmd.Run()

	// Test 'markall' command
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "markall")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run markall command: %v (%s)", err, string(out))
	}
	if !strings.Contains(string(out), "All tasks marked as done") {
		t.Error("Expected 'All tasks marked as done.' output")
	}
}

func TestMainCLIFindByDesc(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Test 'findbydesc' command (will find no tasks since we don't set descriptions in CLI)
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "findbydesc", "urgent")
	out, err := cmd.CombinedOutput()
	// Note: exit code 0 is expected for "not found" case
	if !strings.Contains(string(out), "No tasks found") {
		t.Error("Expected 'No tasks found' message")
	}
}

func TestMainCLIInvalidCommands(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Test invalid command
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "invalid")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for invalid command")
	}
	if !strings.Contains(string(out), "Usage: taskmgr") {
		t.Error("Expected usage message")
	}
}

func TestMainCLIMissingArgs(t *testing.T) {
	// Create temp dir for test
	dir, err := ioutil.TempDir("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to temp dir
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Test add without title
	cmd := exec.Command("go", "run", filepath.Join(origDir, "main.go"), "add")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for add without title")
	}
	if !strings.Contains(string(out), "Usage: taskmgr add") {
		t.Error("Expected add usage message")
	}

	// Test done without index
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "done")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for done without index")
	}
	if !strings.Contains(string(out), "Usage: taskmgr done") {
		t.Error("Expected done usage message")
	}

	// Test remove without index
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "remove")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for remove without index")
	}
	if !strings.Contains(string(out), "Usage: taskmgr remove") {
		t.Error("Expected remove usage message")
	}

	// Test undodone without index
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "undodone")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for undodone without index")
	}
	if !strings.Contains(string(out), "Usage: taskmgr undodone") {
		t.Error("Expected undodone usage message")
	}

	// Test find without title
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "find")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for find without title")
	}
	if !strings.Contains(string(out), "Usage: taskmgr find") {
		t.Error("Expected find usage message")
	}

	// Test bulkadd without titles
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "bulkadd")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for bulkadd without titles")
	}
	if !strings.Contains(string(out), "Usage: taskmgr bulkadd") {
		t.Error("Expected bulkadd usage message")
	}

	// Test findbydesc without description
	cmd = exec.Command("go", "run", filepath.Join(origDir, "main.go"), "findbydesc")
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for findbydesc without description")
	}
	if !strings.Contains(string(out), "Usage: taskmgr findbydesc") {
		t.Error("Expected findbydesc usage message")
	}
}
