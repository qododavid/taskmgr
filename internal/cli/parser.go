package cli

import (
	"strings"
)

type AddOptions struct {
	Title    string
	Priority string
	Due      string
}

type ListOptions struct {
	Priority   string
	Overdue    bool
	DueToday   bool
	DueWithin  int
}

func ParseArgs(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}
	return args[0], args[1:]
}

// ParseAddCommand parses arguments for the add command
func ParseAddCommand(args []string) AddOptions {
	opts := AddOptions{}
	
	for i, arg := range args {
		if strings.HasPrefix(arg, "--priority=") {
			opts.Priority = strings.TrimPrefix(arg, "--priority=")
		} else if strings.HasPrefix(arg, "--due=") {
			opts.Due = strings.TrimPrefix(arg, "--due=")
		} else if arg == "--priority" && i+1 < len(args) {
			opts.Priority = args[i+1]
		} else if arg == "--due" && i+1 < len(args) {
			opts.Due = args[i+1]
		} else if !strings.HasPrefix(arg, "--") && opts.Title == "" {
			// First non-flag argument is the title
			opts.Title = arg
		}
	}
	
	return opts
}

// ParseListCommand parses arguments for the list command
func ParseListCommand(args []string) ListOptions {
	opts := ListOptions{}
	
	for i, arg := range args {
		if strings.HasPrefix(arg, "--priority=") {
			opts.Priority = strings.TrimPrefix(arg, "--priority=")
		} else if arg == "--priority" && i+1 < len(args) {
			opts.Priority = args[i+1]
		} else if arg == "--overdue" {
			opts.Overdue = true
		} else if arg == "--due-today" {
			opts.DueToday = true
		} else if strings.HasPrefix(arg, "--due-within=") {
			// Parse number from --due-within=7days or --due-within=7
			value := strings.TrimPrefix(arg, "--due-within=")
			value = strings.TrimSuffix(value, "days")
			value = strings.TrimSuffix(value, "day")
			if days := parseInt(value); days > 0 {
				opts.DueWithin = days
			}
		}
	}
	
	return opts
}

// Helper function to parse integer
func parseInt(s string) int {
	var result int
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		} else {
			return 0
		}
	}
	return result
}
