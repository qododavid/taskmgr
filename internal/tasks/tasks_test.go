package tasks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTaskManager(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_tasks_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add a task
	if err := manager.Add(Task{Title: "First"}); err != nil {
		t.Fatal("Error adding task:", err)
	}

	list := manager.List()
	if len(list) != 1 || list[0].Title != "First" {
		t.Errorf("Expected one task with title 'First', got %v", list)
	}

	// Mark it done
	err = manager.MarkDone("0")
	if err != nil {
		t.Error("MarkDone returned an error:", err)
	}

	list = manager.List()
	if !list[0].Done {
		t.Error("Expected task to be marked done")
	}

	// Ensure persistence: create a new manager and check again
	newStore := NewFileStore(testFile)
	newManager := NewTaskManager(newStore)
	newList := newManager.List()
	if len(newList) != 1 || newList[0].Title != "First" || !newList[0].Done {
		t.Errorf("Expected persisted task 'First' to be done, got %v", newList)
	}
}

func TestPriorityMethods(t *testing.T) {
	tests := []struct {
		priority Priority
		expectedString string
		expectedColor string
	}{
		{Low, "low", "\033[32m"},
		{Medium, "medium", "\033[33m"},
		{High, "high", "\033[31m"},
		{Critical, "critical", "\033[35m"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedString, func(t *testing.T) {
			if tt.priority.String() != tt.expectedString {
				t.Errorf("Expected String() '%s', got '%s'", tt.expectedString, tt.priority.String())
			}
			if tt.priority.Color() != tt.expectedColor {
				t.Errorf("Expected Color() '%s', got '%s'", tt.expectedColor, tt.priority.Color())
			}
			if tt.priority.ColorReset() != "\033[0m" {
				t.Errorf("Expected ColorReset() '\033[0m', got '%s'", tt.priority.ColorReset())
			}
		})
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		input    string
		expected Priority
		hasError bool
	}{
		{"low", Low, false},
		{"l", Low, false},
		{"LOW", Low, false},
		{"medium", Medium, false},
		{"med", Medium, false},
		{"m", Medium, false},
		{"MEDIUM", Medium, false},
		{"high", High, false},
		{"h", High, false},
		{"HIGH", High, false},
		{"critical", Critical, false},
		{"crit", Critical, false},
		{"c", Critical, false},
		{"CRITICAL", Critical, false},
		{"invalid", Medium, true}, // returns Medium as default but with error
		{"", Medium, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParsePriority(tt.input)
			if tt.hasError && err == nil {
				t.Errorf("Expected error for input '%s', got nil", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("Expected no error for input '%s', got %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("Expected priority %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseDueDate(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	nextWeek := now.AddDate(0, 0, 7)

	tests := []struct {
		name     string
		input    string
		expected *time.Time
		hasError bool
	}{
		{"empty string", "", nil, false},
		{"today", "today", &now, false},
		{"tomorrow", "tomorrow", &tomorrow, false},
		{"next week", "next week", &nextWeek, false},
		{"YYYY-MM-DD format", "2024-01-15", parseTime("2024-01-15"), false},
		{"MM/DD/YYYY format", "01/15/2024", parseTime("2024-01-15"), false},
		{"invalid format", "invalid-date", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDueDate(tt.input)
			if tt.hasError && err == nil {
				t.Errorf("Expected error for input '%s', got nil", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("Expected no error for input '%s', got %v", tt.input, err)
			}
			if tt.expected == nil && result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
			if tt.expected != nil && result == nil {
				t.Errorf("Expected non-nil result, got nil")
			}
			if tt.expected != nil && result != nil {
				// For relative dates like "today", "tomorrow", allow some tolerance
				if tt.input == "today" || tt.input == "tomorrow" || tt.input == "next week" {
					diff := result.Sub(*tt.expected)
					if diff < -time.Minute || diff > time.Minute {
						t.Errorf("Expected time close to %v, got %v (diff: %v)", *tt.expected, *result, diff)
					}
				} else {
					// For absolute dates, expect exact match
					if !result.Equal(*tt.expected) {
						t.Errorf("Expected time %v, got %v", *tt.expected, *result)
					}
				}
			}
		})
	}
}

// Helper function to parse time for tests
func parseTime(dateStr string) *time.Time {
	parsed, _ := time.Parse("2006-01-02", dateStr)
	return &parsed
}

func TestTaskManagerWithPriority(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_priority_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks with different priorities
	tasks := []Task{
		{Title: "Low priority task", Priority: Low},
		{Title: "Medium priority task", Priority: Medium},
		{Title: "High priority task", Priority: High},
		{Title: "Critical priority task", Priority: Critical},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Error adding task: %v", err)
		}
	}

	// Test ListByPriority
	highPriorityTasks := manager.ListByPriority(High)
	if len(highPriorityTasks) != 1 || highPriorityTasks[0].Title != "High priority task" {
		t.Errorf("Expected 1 high priority task, got %v", highPriorityTasks)
	}

	criticalTasks := manager.ListByPriority(Critical)
	if len(criticalTasks) != 1 || criticalTasks[0].Title != "Critical priority task" {
		t.Errorf("Expected 1 critical priority task, got %v", criticalTasks)
	}

	// Test with priority that has no tasks
	emptyResult := manager.ListByPriority(Low)
	if len(emptyResult) != 1 || emptyResult[0].Title != "Low priority task" {
		t.Errorf("Expected 1 low priority task, got %v", emptyResult)
	}
}

func TestTaskManagerWithDueDates(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_due_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	tomorrow := now.AddDate(0, 0, 1)
	nextWeek := now.AddDate(0, 0, 7)

	// Add tasks with different due dates
	tasks := []Task{
		{Title: "Overdue task", DueDate: &yesterday},
		{Title: "Due today task", DueDate: &today},
		{Title: "Due tomorrow task", DueDate: &tomorrow},
		{Title: "Due next week task", DueDate: &nextWeek},
		{Title: "No due date task"},
		{Title: "Overdue but done task", DueDate: &yesterday, Done: true},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Error adding task: %v", err)
		}
	}

	// Test ListOverdue
	overdueTasks := manager.ListOverdue()
	if len(overdueTasks) != 1 || overdueTasks[0].Title != "Overdue task" {
		t.Errorf("Expected 1 overdue task, got %v", overdueTasks)
	}

	// Test ListDueToday
	dueTodayTasks := manager.ListDueToday()
	if len(dueTodayTasks) != 1 || dueTodayTasks[0].Title != "Due today task" {
		t.Errorf("Expected 1 due today task, got %v", dueTodayTasks)
	}

	// Test ListDueWithin
	// ListDueWithin includes tasks from now until cutoff, so "due today" task at 12:00 and "due tomorrow" should both be included
	dueWithin2Days := manager.ListDueWithin(2)
	if len(dueWithin2Days) != 2 {
		t.Errorf("Expected 2 tasks due within 2 days (today + tomorrow), got %d: %v", len(dueWithin2Days), dueWithin2Days)
	}

	// Within 10 days should include: due today, due tomorrow, due next week (3 tasks)
	dueWithin10Days := manager.ListDueWithin(10)
	if len(dueWithin10Days) != 3 {
	t.Errorf("Expected 3 tasks due within 10 days, got %d: %v", len(dueWithin10Days), dueWithin10Days)
	}
	}

func TestTaskTagMethods(t *testing.T) {
	task := Task{Title: "Test Task"}

	// Test HasTag on empty task
	if task.HasTag("work") {
		t.Error("Empty task should not have any tags")
	}

	// Test AddTag
	task.AddTag("work")
	if !task.HasTag("work") {
		t.Error("Task should have 'work' tag after adding it")
	}
	if len(task.Tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(task.Tags))
	}

	// Test adding duplicate tag
	task.AddTag("work")
	if len(task.Tags) != 1 {
		t.Error("Adding duplicate tag should not increase tag count")
	}

	// Test adding multiple tags
	task.AddTag("urgent")
	task.AddTag("bug")
	if len(task.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(task.Tags))
	}
	if !task.HasTag("urgent") || !task.HasTag("bug") {
		t.Error("Task should have all added tags")
	}

	// Test case insensitive HasTag
	if !task.HasTag("WORK") || !task.HasTag("Work") {
		t.Error("HasTag should be case insensitive")
	}

	// Test RemoveTag
	task.RemoveTag("work")
	if task.HasTag("work") {
		t.Error("Task should not have 'work' tag after removing it")
	}
	if len(task.Tags) != 2 {
		t.Errorf("Expected 2 tags after removal, got %d", len(task.Tags))
	}

	// Test removing non-existent tag
	task.RemoveTag("nonexistent")
	if len(task.Tags) != 2 {
		t.Error("Removing non-existent tag should not change tag count")
	}

	// Test case insensitive RemoveTag
	task.RemoveTag("URGENT")
	if task.HasTag("urgent") {
		t.Error("RemoveTag should be case insensitive")
	}
	if len(task.Tags) != 1 {
		t.Errorf("Expected 1 tag after case insensitive removal, got %d", len(task.Tags))
	}

	// Test tag normalization (lowercase and trimmed)
	task.AddTag(" PERSONAL ")
	if !task.HasTag("personal") {
		t.Error("Tag should be normalized to lowercase and trimmed")
	}
	expectedTag := "personal"
	found := false
	for _, tag := range task.Tags {
		if tag == expectedTag {
			found = true
			break
		}
	}
	if !found {
		t.Error("Tag should be stored in normalized form")
	}
}

func TestTaskManagerTagMethods(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := ioutil.TempDir("", "taskmgr_tag_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "tasks.json")
	store := NewFileStore(testFile)
	manager := NewTaskManager(store)

	// Add tasks with different tags
	tasks := []Task{
		{Title: "Work task 1", Tags: []string{"work", "urgent"}},
		{Title: "Work task 2", Tags: []string{"work", "meeting"}},
		{Title: "Personal task", Tags: []string{"personal", "shopping"}},
		{Title: "Bug fix", Tags: []string{"work", "bug", "urgent"}},
		{Title: "No tags task"},
	}

	for _, task := range tasks {
		if err := manager.Add(task); err != nil {
			t.Fatalf("Error adding task: %v", err)
		}
	}

	// Test ListByTag
	workTasks := manager.ListByTag("work")
	if len(workTasks) != 3 {
		t.Errorf("Expected 3 work tasks, got %d", len(workTasks))
	}

	urgentTasks := manager.ListByTag("urgent")
	if len(urgentTasks) != 2 {
		t.Errorf("Expected 2 urgent tasks, got %d", len(urgentTasks))
	}

	personalTasks := manager.ListByTag("personal")
	if len(personalTasks) != 1 {
		t.Errorf("Expected 1 personal task, got %d", len(personalTasks))
	}

	// Test case insensitive ListByTag
	workTasksUpper := manager.ListByTag("WORK")
	if len(workTasksUpper) != 3 {
		t.Errorf("Expected 3 work tasks (case insensitive), got %d", len(workTasksUpper))
	}

	// Test ListByTag with non-existent tag
	nonExistentTasks := manager.ListByTag("nonexistent")
	if len(nonExistentTasks) != 0 {
		t.Errorf("Expected 0 tasks for non-existent tag, got %d", len(nonExistentTasks))
	}

	// Test GetAllTags
	allTags := manager.GetAllTags()
	expectedTags := []string{"bug", "meeting", "personal", "shopping", "urgent", "work"} // sorted alphabetically
	if len(allTags) != len(expectedTags) {
		t.Errorf("Expected %d unique tags, got %d", len(expectedTags), len(allTags))
	}
	for i, expectedTag := range expectedTags {
		if i >= len(allTags) || allTags[i] != expectedTag {
			t.Errorf("Expected tag[%d] '%s', got '%s'", i, expectedTag, allTags[i])
		}
	}

	// Test AddTagToTask
	err = manager.AddTagToTask("4", "important") // Add tag to "No tags task"
	if err != nil {
		t.Errorf("Error adding tag to task: %v", err)
	}

	// Verify tag was added
	updatedTasks := manager.List()
	if !updatedTasks[4].HasTag("important") {
		t.Error("Task should have 'important' tag after adding it")
	}

	// Test AddTagToTask with invalid index
	err = manager.AddTagToTask("99", "invalid")
	if err == nil {
		t.Error("Expected error when adding tag to invalid index")
	}

	err = manager.AddTagToTask("invalid", "tag")
	if err == nil {
		t.Error("Expected error when using invalid index format")
	}

	// Test RemoveTagFromTask
	err = manager.RemoveTagFromTask("0", "urgent") // Remove 'urgent' from first work task
	if err != nil {
		t.Errorf("Error removing tag from task: %v", err)
	}

	// Verify tag was removed
	updatedTasks = manager.List()
	if updatedTasks[0].HasTag("urgent") {
		t.Error("Task should not have 'urgent' tag after removing it")
	}
	if !updatedTasks[0].HasTag("work") {
		t.Error("Task should still have 'work' tag after removing 'urgent'")
	}

	// Test RemoveTagFromTask with non-existent tag
	err = manager.RemoveTagFromTask("0", "nonexistent")
	if err != nil {
		t.Errorf("Removing non-existent tag should not return error: %v", err)
	}

	// Test RemoveTagFromTask with invalid index
	err = manager.RemoveTagFromTask("99", "work")
	if err == nil {
		t.Error("Expected error when removing tag from invalid index")
	}

	err = manager.RemoveTagFromTask("invalid", "work")
	if err == nil {
		t.Error("Expected error when using invalid index format")
	}

	// Test persistence of tags
	newStore := NewFileStore(testFile)
	newManager := NewTaskManager(newStore)
	persistentTasks := newManager.List()

	// Verify the tag changes persisted
	if persistentTasks[0].HasTag("urgent") {
		t.Error("Removed tag should not persist")
	}
	if !persistentTasks[0].HasTag("work") {
		t.Error("Remaining tags should persist")
	}
	if !persistentTasks[4].HasTag("important") {
		t.Error("Added tag should persist")
	}
}
