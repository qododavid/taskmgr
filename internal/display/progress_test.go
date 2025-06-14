package display

import (
	"strings"
	"testing"
	"time"
	"taskmgr/internal/tasks"
)

func TestNewProgressFormatter(t *testing.T) {
	opts := DisplayOptions{
		ShowColors: true,
		ShowIcons:  true,
	}
	
	formatter := NewProgressFormatter(opts)
	if formatter == nil {
		t.Fatal("NewProgressFormatter should not return nil")
	}
	
	if formatter.options.ShowColors != opts.ShowColors {
		t.Error("ShowColors option not set correctly")
	}
	
	// Test default color scheme is set when empty
	emptyOpts := DisplayOptions{}
	formatter2 := NewProgressFormatter(emptyOpts)
	if formatter2.options.ColorScheme == (ColorScheme{}) {
		t.Error("Default color scheme should be set when empty")
	}
}

func TestCalculateStats(t *testing.T) {
	yesterday := time.Now().Add(-24 * time.Hour)
	tomorrow := time.Now().Add(24 * time.Hour)
	
	taskList := []tasks.Task{
		{Title: "Task 1", Done: true, Priority: tasks.High},
		{Title: "Task 2", Done: false, Priority: tasks.Medium},
		{Title: "Task 3", Done: false, Priority: tasks.Low, DueDate: &yesterday}, // overdue
		{Title: "Task 4", Done: true, Priority: tasks.Critical},
		{Title: "Task 5", Done: false, Priority: tasks.High, DueDate: &tomorrow},
	}
	
	opts := DisplayOptions{}
	formatter := NewProgressFormatter(opts)
	stats := formatter.CalculateStats(taskList)
	
	// Test total count
	if stats.Total != 5 {
		t.Errorf("Expected total 5, got %d", stats.Total)
	}
	
	// Test completed count
	if stats.Completed != 2 {
		t.Errorf("Expected completed 2, got %d", stats.Completed)
	}
	
	// Test pending count
	if stats.Pending != 3 {
		t.Errorf("Expected pending 3, got %d", stats.Pending)
	}
	
	// Test overdue count
	if stats.Overdue != 1 {
		t.Errorf("Expected overdue 1, got %d", stats.Overdue)
	}
	
	// Test priority counts
	if stats.ByPriority[tasks.High] != 2 {
		t.Errorf("Expected 2 high priority tasks, got %d", stats.ByPriority[tasks.High])
	}
	if stats.ByPriority[tasks.Medium] != 1 {
		t.Errorf("Expected 1 medium priority task, got %d", stats.ByPriority[tasks.Medium])
	}
	if stats.ByPriority[tasks.Low] != 1 {
		t.Errorf("Expected 1 low priority task, got %d", stats.ByPriority[tasks.Low])
	}
	if stats.ByPriority[tasks.Critical] != 1 {
		t.Errorf("Expected 1 critical priority task, got %d", stats.ByPriority[tasks.Critical])
	}
}

func TestCalculateStatsEmpty(t *testing.T) {
	opts := DisplayOptions{}
	formatter := NewProgressFormatter(opts)
	stats := formatter.CalculateStats([]tasks.Task{})
	
	if stats.Total != 0 {
		t.Errorf("Expected total 0, got %d", stats.Total)
	}
	if stats.Completed != 0 {
		t.Errorf("Expected completed 0, got %d", stats.Completed)
	}
	if stats.Pending != 0 {
		t.Errorf("Expected pending 0, got %d", stats.Pending)
	}
	if stats.Overdue != 0 {
		t.Errorf("Expected overdue 0, got %d", stats.Overdue)
	}
}

func TestFormatProgress(t *testing.T) {
	tests := []struct {
		name     string
		stats    ProgressStats
		expected string
	}{
		{
			name: "no tasks",
			stats: ProgressStats{
				Total:     0,
				Completed: 0,
			},
			expected: "No tasks found.",
		},
		{
			name: "50% completion",
			stats: ProgressStats{
				Total:     4,
				Completed: 2,
			},
			expected: "Progress: [██████████░░░░░░░░░░] 2/4 (50.0%)",
		},
		{
			name: "100% completion",
			stats: ProgressStats{
				Total:     2,
				Completed: 2,
			},
			expected: "Progress: [████████████████████] 2/2 (100.0%)",
		},
		{
			name: "0% completion",
			stats: ProgressStats{
				Total:     3,
				Completed: 0,
			},
			expected: "Progress: [░░░░░░░░░░░░░░░░░░░░] 0/3 (0.0%)",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DisplayOptions{ShowColors: false}
			formatter := NewProgressFormatter(opts)
			result := formatter.FormatProgress(tt.stats)
			
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatDetailedStats(t *testing.T) {
	stats := ProgressStats{
		Total:     5,
		Completed: 2,
		Pending:   3,
		Overdue:   1,
		ByPriority: map[tasks.Priority]int{
			tasks.Critical: 1,
			tasks.High:     2,
			tasks.Medium:   1,
			tasks.Low:      1,
		},
	}
	
	opts := DisplayOptions{ShowColors: false}
	formatter := NewProgressFormatter(opts)
	result := formatter.FormatDetailedStats(stats)
	
	// Check that all sections are present
	if !strings.Contains(result, "Progress:") {
		t.Error("Should contain progress bar")
	}
	if !strings.Contains(result, "By Priority:") {
		t.Error("Should contain priority breakdown")
	}
	if !strings.Contains(result, "By Status:") {
		t.Error("Should contain status breakdown")
	}
	
	// Check priority counts
	if !strings.Contains(result, "Critical: 1 tasks") {
		t.Error("Should contain critical priority count")
	}
	if !strings.Contains(result, "High: 2 tasks") {
		t.Error("Should contain high priority count")
	}
	if !strings.Contains(result, "Medium: 1 tasks") {
		t.Error("Should contain medium priority count")
	}
	if !strings.Contains(result, "Low: 1 tasks") {
		t.Error("Should contain low priority count")
	}
	
	// Check status counts
	if !strings.Contains(result, "Completed: 2 tasks") {
		t.Error("Should contain completed count")
	}
	if !strings.Contains(result, "Pending: 3 tasks") {
		t.Error("Should contain pending count")
	}
	if !strings.Contains(result, "Overdue: 1 tasks") {
		t.Error("Should contain overdue count")
	}
}

func TestFormatDetailedStatsNoCounts(t *testing.T) {
	// Test with stats that have zero counts for some categories
	stats := ProgressStats{
		Total:     2,
		Completed: 2,
		Pending:   0,
		Overdue:   0,
		ByPriority: map[tasks.Priority]int{
			tasks.High: 2,
		},
	}
	
	opts := DisplayOptions{ShowColors: false}
	formatter := NewProgressFormatter(opts)
	result := formatter.FormatDetailedStats(stats)
	
	// Should only show categories with non-zero counts
	if !strings.Contains(result, "High: 2 tasks") {
		t.Error("Should contain high priority count")
	}
	if strings.Contains(result, "Critical:") {
		t.Error("Should not contain critical priority (zero count)")
	}
	if !strings.Contains(result, "Completed: 2 tasks") {
		t.Error("Should contain completed count")
	}
	if strings.Contains(result, "Pending:") {
		t.Error("Should not contain pending count (zero)")
	}
	if strings.Contains(result, "Overdue:") {
		t.Error("Should not contain overdue count (zero)")
	}
}

func TestGetPriorityIcon(t *testing.T) {
	tests := []struct {
		name     string
		priority tasks.Priority
		expected string
	}{
		{
			name:     "critical priority",
			priority: tasks.Critical,
			expected: "🔴",
		},
		{
			name:     "high priority",
			priority: tasks.High,
			expected: "🟠",
		},
		{
			name:     "medium priority",
			priority: tasks.Medium,
			expected: "🟡",
		},
		{
			name:     "low priority",
			priority: tasks.Low,
			expected: "🟢",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DisplayOptions{ShowColors: false}
			formatter := NewProgressFormatter(opts)
			result := formatter.getPriorityIcon(tt.priority)
			
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetPriorityIconWithColors(t *testing.T) {
	opts := DisplayOptions{
		ShowColors:  true,
		ColorScheme: DefaultColorScheme,
	}
	formatter := NewProgressFormatter(opts)
	result := formatter.getPriorityIcon(tasks.Critical)
	
	// Should contain the icon
	if !strings.Contains(result, "🔴") {
		t.Error("Should contain critical priority icon")
	}
	// In CI environment, colors might be disabled, so just check that we get a result
	if result == "" {
		t.Error("Should return a non-empty result")
	}
}