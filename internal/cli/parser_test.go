package cli

import "testing"

func TestParseArgs(t *testing.T) {
	cmd, rest := ParseArgs([]string{"add", "MyTask"})
	if cmd != "add" {
		t.Errorf("Expected cmd 'add', got '%s'", cmd)
	}
	if len(rest) != 1 || rest[0] != "MyTask" {
		t.Errorf("Expected args ['MyTask'], got %v", rest)
	}
}

func TestParseArgsEmpty(t *testing.T) {
	cmd, rest := ParseArgs([]string{})
	if cmd != "" {
		t.Errorf("Expected empty cmd, got '%s'", cmd)
	}
	if rest != nil {
		t.Errorf("Expected nil rest, got %v", rest)
	}
}

func TestParseArgsSingleCommand(t *testing.T) {
	cmd, rest := ParseArgs([]string{"list"})
	if cmd != "list" {
		t.Errorf("Expected cmd 'list', got '%s'", cmd)
	}
	if len(rest) != 0 {
		t.Errorf("Expected empty rest, got %v", rest)
	}
}

func TestParseArgsMultipleArgs(t *testing.T) {
	cmd, rest := ParseArgs([]string{"bulkadd", "task1,task2,task3"})
	if cmd != "bulkadd" {
		t.Errorf("Expected cmd 'bulkadd', got '%s'", cmd)
	}
	if len(rest) != 1 || rest[0] != "task1,task2,task3" {
		t.Errorf("Expected args ['task1,task2,task3'], got %v", rest)
	}
}
