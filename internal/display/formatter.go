package display

import (
	"fmt"
	"strings"
	"time"
	"taskmgr/internal/tasks"
)

type DisplayOptions struct {
	ShowColors    bool
	ShowIcons     bool
	TableFormat   bool
	ShowTags      bool
	ShowDueDate   bool
	ShowPriority  bool
	ColorScheme   ColorScheme
}

type TaskFormatter struct {
	options DisplayOptions
}

// NewTaskFormatter creates a new task formatter with the given options
func NewTaskFormatter(opts DisplayOptions) *TaskFormatter {
	// Set default color scheme if not provided
	if opts.ColorScheme == (ColorScheme{}) {
		opts.ColorScheme = DefaultColorScheme
	}
	return &TaskFormatter{options: opts}
}

// FormatTask formats a single task for display
func (tf *TaskFormatter) FormatTask(index int, task tasks.Task) string {
	if tf.options.TableFormat {
		return tf.formatTableRow(index, task)
	}
	return tf.formatListItem(index, task)
}

// formatListItem formats a task as a list item
func (tf *TaskFormatter) formatListItem(index int, task tasks.Task) string {
	var parts []string
	
	// Status icon
	statusIcon := tf.getStatusIcon(task)
	parts = append(parts, fmt.Sprintf("%s %d:", statusIcon, index))
	
	// Title with color
	title := tf.formatTitle(task)
	parts = append(parts, title)
	
	// Priority
	if tf.options.ShowPriority {
		priority := tf.formatPriority(task.Priority)
		parts = append(parts, priority)
	}
	
	// Tags
	if tf.options.ShowTags && len(task.Tags) > 0 {
		tags := tf.formatTags(task.Tags)
		parts = append(parts, tags)
	}
	
	// Due date
	if tf.options.ShowDueDate && task.DueDate != nil {
		dueDate := tf.formatDueDate(*task.DueDate, task.Done)
		parts = append(parts, dueDate)
	}
	
	return strings.Join(parts, " ")
}

// formatTableRow formats a task as a table row
func (tf *TaskFormatter) formatTableRow(index int, task tasks.Task) string {
	status := tf.getStatusIcon(task)
	priority := tf.formatPriority(task.Priority)
	title := tf.formatTitle(task)
	
	tags := ""
	if len(task.Tags) > 0 {
		tags = tf.formatTags(task.Tags)
	}
	
	dueDate := "N/A"
	if task.DueDate != nil {
		dueDate = tf.formatDueDate(*task.DueDate, task.Done)
	}
	
	// Format as table row with fixed widths
	return fmt.Sprintf("%-3d | %-6s | %-8s | %-25s | %-15s | %s", 
		index, status, priority, truncateString(title, 25), truncateString(tags, 15), dueDate)
}

// FormatTableHeader returns the table header
func (tf *TaskFormatter) FormatTableHeader() string {
	header := "ID  | Status | Priority | Title                     | Tags            | Due Date"
	if tf.options.ShowColors {
		return Colorize(Bold, header)
	}
	return header
}

// FormatTableSeparator returns the table separator line
func (tf *TaskFormatter) FormatTableSeparator() string {
	return "----+--------+----------+---------------------------+-----------------+----------"
}

// getStatusIcon returns the appropriate status icon for a task
func (tf *TaskFormatter) getStatusIcon(task tasks.Task) string {
	if !tf.options.ShowIcons {
		if task.Done {
			return "[x]"
		}
		return "[ ]"
	}
	
	if task.Done {
		if tf.options.ShowColors {
			return Colorize(tf.options.ColorScheme.Completed, "✅")
		}
		return "✅"
	}
	
	// Check if overdue
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		if tf.options.ShowColors {
			return Colorize(tf.options.ColorScheme.Overdue, "❌")
		}
		return "❌"
	}
	
	// Pending task
	if tf.options.ShowColors {
		return Colorize(tf.options.ColorScheme.Pending, "⭕")
	}
	return "⭕"
}

// formatTitle formats the task title with appropriate styling
func (tf *TaskFormatter) formatTitle(task tasks.Task) string {
	if !tf.options.ShowColors {
		return task.Title
	}
	
	if task.Done {
		// Completed tasks in green with strikethrough effect (using dim)
		return Colorize(tf.options.ColorScheme.Completed, task.Title)
	}
	
	// Check if overdue
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		return Colorize(tf.options.ColorScheme.Overdue, task.Title)
	}
	
	return task.Title
}

// formatPriority formats the priority with color and icon
func (tf *TaskFormatter) formatPriority(priority tasks.Priority) string {
	if !tf.options.ShowColors {
		return fmt.Sprintf("[%s]", strings.ToUpper(priority.String()[:3]))
	}
	
	var color Color
	var icon string
	
	switch priority {
	case tasks.Critical:
		color = tf.options.ColorScheme.Critical
		icon = "🔴"
	case tasks.High:
		color = tf.options.ColorScheme.High
		icon = "🟠"
	case tasks.Medium:
		color = tf.options.ColorScheme.Medium
		icon = "🟡"
	case tasks.Low:
		color = tf.options.ColorScheme.Low
		icon = "🟢"
	}
	
	text := fmt.Sprintf("[%s %s]", icon, strings.ToUpper(priority.String()[:3]))
	return Colorize(color, text)
}

// formatTags formats the task tags
func (tf *TaskFormatter) formatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	
	tagStr := strings.Join(tags, ", ")
	if !tf.options.ShowColors {
		return fmt.Sprintf("[%s]", tagStr)
	}
	
	return fmt.Sprintf("[%s]", Colorize(tf.options.ColorScheme.Tags, tagStr))
}

// formatDescription formats the task description as a tag-like element
func (tf *TaskFormatter) formatDescription(description string) string {
	if description == "" {
		return ""
	}
	
	if !tf.options.ShowColors {
		return fmt.Sprintf("[%s]", description)
	}
	
	return fmt.Sprintf("[%s]", Colorize(tf.options.ColorScheme.Tags, description))
}

// formatDueDate formats the due date with appropriate color coding
func (tf *TaskFormatter) formatDueDate(dueDate time.Time, isDone bool) string {
	now := time.Now()
	diff := dueDate.Sub(now)
	
	var text string
	var color Color = tf.options.ColorScheme.DueDate
	
	if diff < 0 && !isDone {
		// Overdue
		days := int(-diff.Hours() / 24)
		if days == 0 {
			text = "(OVERDUE: today)"
		} else {
			text = fmt.Sprintf("(OVERDUE: %d days)", days)
		}
		color = tf.options.ColorScheme.Overdue
	} else if diff < 24*time.Hour && diff >= 0 {
		// Due today
		text = "(Due: today)"
		color = tf.options.ColorScheme.Medium
	} else if diff < 48*time.Hour {
		// Due tomorrow
		text = "(Due: tomorrow)"
		color = tf.options.ColorScheme.Medium
	} else if diff >= 0 {
		// Due in future
		days := int(diff.Hours() / 24)
		text = fmt.Sprintf("(Due: %d days)", days)
		color = tf.options.ColorScheme.Low
	} else {
		// Completed overdue task
		text = fmt.Sprintf("(Was due: %s)", dueDate.Format("2006-01-02"))
		color = tf.options.ColorScheme.Completed
	}
	
	if tf.options.ShowColors {
		return Colorize(color, text)
	}
	return text
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}