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

func TestParseAddCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected AddOptions
	}{
		{
			name: "title only",
			args: []string{"Fix bug"},
			expected: AddOptions{Title: "Fix bug"},
		},
		{
			name: "title with priority flag",
			args: []string{"Fix bug", "--priority=high"},
			expected: AddOptions{Title: "Fix bug", Priority: "high"},
		},
		{
			name: "title with due date flag",
			args: []string{"Fix bug", "--due=2024-01-15"},
			expected: AddOptions{Title: "Fix bug", Due: "2024-01-15"},
		},
		{
			name: "title with both flags",
			args: []string{"Fix bug", "--priority=high", "--due=tomorrow"},
			expected: AddOptions{Title: "Fix bug", Priority: "high", Due: "tomorrow"},
		},
		{
			name: "priority with space separator",
			args: []string{"Fix bug", "--priority", "medium"},
			expected: AddOptions{Title: "Fix bug", Priority: "medium"},
		},
		{
			name: "due with space separator",
			args: []string{"Fix bug", "--due", "today"},
			expected: AddOptions{Title: "Fix bug", Due: "today"},
		},
		{
			name: "empty args",
			args: []string{},
			expected: AddOptions{},
		},
		{
			name: "title with tags flag",
			args: []string{"Fix bug", "--tags=work,urgent"},
			expected: AddOptions{Title: "Fix bug", Tags: []string{"work", "urgent"}},
		},
		{
			name: "title with tags space separator",
			args: []string{"Fix bug", "--tags", "personal,shopping"},
			expected: AddOptions{Title: "Fix bug", Tags: []string{"personal", "shopping"}},
		},
		{
			name: "title with all flags including tags",
			args: []string{"Fix bug", "--priority=high", "--due=tomorrow", "--tags=work,urgent,bug"},
			expected: AddOptions{Title: "Fix bug", Priority: "high", Due: "tomorrow", Tags: []string{"work", "urgent", "bug"}},
		},
		{
			name: "tags with spaces and mixed case",
			args: []string{"Fix bug", "--tags=Work, URGENT , Bug"},
			expected: AddOptions{Title: "Fix bug", Tags: []string{"work", "urgent", "bug"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAddCommand(tt.args)
			if result.Title != tt.expected.Title {
				t.Errorf("Expected title '%s', got '%s'", tt.expected.Title, result.Title)
			}
			if result.Priority != tt.expected.Priority {
				t.Errorf("Expected priority '%s', got '%s'", tt.expected.Priority, result.Priority)
			}
			if result.Due != tt.expected.Due {
				t.Errorf("Expected due '%s', got '%s'", tt.expected.Due, result.Due)
			}
			// Check tags
			if len(result.Tags) != len(tt.expected.Tags) {
				t.Errorf("Expected %d tags, got %d", len(tt.expected.Tags), len(result.Tags))
			}
			for i, tag := range tt.expected.Tags {
				if i >= len(result.Tags) || result.Tags[i] != tag {
					t.Errorf("Expected tag[%d] '%s', got '%s'", i, tag, result.Tags[i])
				}
			}
		})
	}
}

func TestParseListCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected ListOptions
	}{
		{
			name: "no args",
			args: []string{},
			expected: ListOptions{},
		},
		{
			name: "priority filter",
			args: []string{"--priority=high"},
			expected: ListOptions{Priority: "high"},
		},
		{
			name: "priority with space",
			args: []string{"--priority", "low"},
			expected: ListOptions{Priority: "low"},
		},
		{
			name: "overdue flag",
			args: []string{"--overdue"},
			expected: ListOptions{Overdue: true},
		},
		{
			name: "due today flag",
			args: []string{"--due-today"},
			expected: ListOptions{DueToday: true},
		},
		{
			name: "due within days",
			args: []string{"--due-within=7days"},
			expected: ListOptions{DueWithin: 7},
		},
		{
			name: "due within number only",
			args: []string{"--due-within=5"},
			expected: ListOptions{DueWithin: 5},
		},
		{
			name: "multiple flags",
			args: []string{"--priority=medium", "--overdue"},
			expected: ListOptions{Priority: "medium", Overdue: true},
		},
		{
			name: "tag filter",
			args: []string{"--tag=work"},
			expected: ListOptions{Tag: "work"},
		},
		{
			name: "tag with space separator",
			args: []string{"--tag", "personal"},
			expected: ListOptions{Tag: "personal"},
		},
		{
			name: "tag with priority filter",
			args: []string{"--tag=urgent", "--priority=high"},
			expected: ListOptions{Tag: "urgent", Priority: "high"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseListCommand(tt.args)
			if result.Priority != tt.expected.Priority {
				t.Errorf("Expected priority '%s', got '%s'", tt.expected.Priority, result.Priority)
			}
			if result.Overdue != tt.expected.Overdue {
				t.Errorf("Expected overdue %v, got %v", tt.expected.Overdue, result.Overdue)
			}
			if result.DueToday != tt.expected.DueToday {
				t.Errorf("Expected due today %v, got %v", tt.expected.DueToday, result.DueToday)
			}
			if result.DueWithin != tt.expected.DueWithin {
				t.Errorf("Expected due within %d, got %d", tt.expected.DueWithin, result.DueWithin)
			}
			if result.Tag != tt.expected.Tag {
				t.Errorf("Expected tag '%s', got '%s'", tt.expected.Tag, result.Tag)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"7", 7},
		{"abc", 0}, // invalid input should return 0
		{"", 0},    // empty string should return 0
		{"12a", 0}, // mixed input should return 0
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseInt(tt.input)
			if result != tt.expected {
				t.Errorf("parseInt(%s) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}
