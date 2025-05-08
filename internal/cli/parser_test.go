package cli

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	// Test with normal command and arguments
	cmd, rest := ParseArgs([]string{"add", "MyTask"})
	if cmd != "add" {
		t.Errorf("Expected cmd 'add', got '%s'", cmd)
	}
	if len(rest) != 1 || rest[0] != "MyTask" {
		t.Errorf("Expected args ['MyTask'], got %v", rest)
	}
	
	// Test with empty arguments
	cmd, rest = ParseArgs([]string{})
	if cmd != "" {
		t.Errorf("Expected empty cmd, got '%s'", cmd)
	}
	if rest != nil {
		t.Errorf("Expected nil rest args, got %v", rest)
	}
	
	// Test with just command, no arguments
	cmd, rest = ParseArgs([]string{"list"})
	if cmd != "list" {
		t.Errorf("Expected cmd 'list', got '%s'", cmd)
	}
	if len(rest) != 0 {
		t.Errorf("Expected empty args, got %v", rest)
	}
	
	// Test with multiple arguments
	cmd, rest = ParseArgs([]string{"add", "Task1", "Priority:High"})
	if cmd != "add" {
		t.Errorf("Expected cmd 'add', got '%s'", cmd)
	}
	expectedArgs := []string{"Task1", "Priority:High"}
	if !reflect.DeepEqual(rest, expectedArgs) {
		t.Errorf("Expected args %v, got %v", expectedArgs, rest)
	}
}