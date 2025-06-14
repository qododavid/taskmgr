package display

import (
	"fmt"
	"strings"
	"time"
	"taskmgr/internal/tasks"
)

// ProgressStats holds statistics about task completion
type ProgressStats struct {
	Total      int
	Completed  int
	Pending    int
	Overdue    int
	ByPriority map[tasks.Priority]int
}

// ProgressFormatter handles progress display formatting
type ProgressFormatter struct {
	options DisplayOptions
}

// NewProgressFormatter creates a new progress formatter
func NewProgressFormatter(opts DisplayOptions) *ProgressFormatter {
	if opts.ColorScheme == (ColorScheme{}) {
		opts.ColorScheme = DefaultColorScheme
	}
	return &ProgressFormatter{options: opts}
}

// CalculateStats calculates progress statistics from a list of tasks
func (pf *ProgressFormatter) CalculateStats(taskList []tasks.Task) ProgressStats {
	stats := ProgressStats{
		Total:      len(taskList),
		ByPriority: make(map[tasks.Priority]int),
	}
	
	for _, task := range taskList {
		// Count by priority
		stats.ByPriority[task.Priority]++
		
		// Count by status
		if task.Done {
			stats.Completed++
		} else {
			stats.Pending++
			
			// Check if overdue
			if task.DueDate != nil && task.DueDate.Before(time.Now()) {
				stats.Overdue++
			}
		}
	}
	
	return stats
}

// FormatProgress formats a progress bar
func (pf *ProgressFormatter) FormatProgress(stats ProgressStats) string {
	if stats.Total == 0 {
		return "No tasks found."
	}
	
	percentage := float64(stats.Completed) / float64(stats.Total) * 100
	barWidth := 20
	filled := int(percentage / 100 * float64(barWidth))
	
	var bar string
	if pf.options.ShowColors {
		filledBar := strings.Repeat("█", filled)
		emptyBar := strings.Repeat("░", barWidth-filled)
		bar = Colorize(pf.options.ColorScheme.Completed, filledBar) + 
			  Colorize(Gray, emptyBar)
	} else {
		bar = strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	}
	
	return fmt.Sprintf("Progress: [%s] %d/%d (%.1f%%)", 
		bar, stats.Completed, stats.Total, percentage)
}

// FormatDetailedStats formats detailed statistics
func (pf *ProgressFormatter) FormatDetailedStats(stats ProgressStats) string {
	var lines []string
	
	// Progress bar
	lines = append(lines, pf.FormatProgress(stats))
	lines = append(lines, "")
	
	// By Priority
	lines = append(lines, "By Priority:")
	priorities := []tasks.Priority{tasks.Critical, tasks.High, tasks.Medium, tasks.Low}
	for _, priority := range priorities {
		count := stats.ByPriority[priority]
		if count > 0 {
			icon := pf.getPriorityIcon(priority)
			priorityName := strings.Title(priority.String())
			line := fmt.Sprintf("  %s %s: %d tasks", icon, priorityName, count)
			lines = append(lines, line)
		}
	}
	
	lines = append(lines, "")
	
	// By Status
	lines = append(lines, "By Status:")
	if stats.Completed > 0 {
		icon := "✅"
		if pf.options.ShowColors {
			icon = Colorize(pf.options.ColorScheme.Completed, icon)
		}
		lines = append(lines, fmt.Sprintf("  %s Completed: %d tasks", icon, stats.Completed))
	}
	
	if stats.Pending > 0 {
		icon := "⭕"
		if pf.options.ShowColors {
			icon = Colorize(pf.options.ColorScheme.Pending, icon)
		}
		lines = append(lines, fmt.Sprintf("  %s Pending: %d tasks", icon, stats.Pending))
	}
	
	if stats.Overdue > 0 {
		icon := "❌"
		if pf.options.ShowColors {
			icon = Colorize(pf.options.ColorScheme.Overdue, icon)
		}
		lines = append(lines, fmt.Sprintf("  %s Overdue: %d tasks", icon, stats.Overdue))
	}
	
	return strings.Join(lines, "\n")
}

// getPriorityIcon returns the appropriate icon for a priority level
func (pf *ProgressFormatter) getPriorityIcon(priority tasks.Priority) string {
	var icon string
	var color Color
	
	switch priority {
	case tasks.Critical:
		icon = "🔴"
		color = pf.options.ColorScheme.Critical
	case tasks.High:
		icon = "🟠"
		color = pf.options.ColorScheme.High
	case tasks.Medium:
		icon = "🟡"
		color = pf.options.ColorScheme.Medium
	case tasks.Low:
		icon = "🟢"
		color = pf.options.ColorScheme.Low
	}
	
	if pf.options.ShowColors {
		return Colorize(color, icon)
	}
	return icon
}