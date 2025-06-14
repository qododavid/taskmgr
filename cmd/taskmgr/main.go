package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	"taskmgr/internal/cli"
	"taskmgr/internal/display"
	"taskmgr/internal/tasks"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://c41e794f3c4d70cd5616e8586b60545f@o4509316990959616.ingest.us.sentry.io/4509317167579136",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	cmd, args := cli.ParseArgs(os.Args[1:])
	store := tasks.NewFileStore("tasks.json") // Updated to file-based store
	manager := tasks.NewTaskManager(store)

	switch cmd {
	case "add":
		opts := cli.ParseAddCommand(args)
		if opts.Title == "" {
			fmt.Println("Usage: taskmgr add <title> [--priority=<low|medium|high|critical>] [--due=<date>] [--tags=<tag1,tag2,...>]")
			fmt.Println("Examples:")
			fmt.Println("  taskmgr add \"Fix bug\" --priority=high --due=2024-01-15 --tags=work,urgent")
			fmt.Println("  taskmgr add \"Review PR\" --priority=medium --due=tomorrow --tags=work,code-review")
			fmt.Println("  taskmgr add \"Buy groceries\" --tags=personal,shopping")
			os.Exit(1)
		}
		
		t := tasks.Task{Title: opts.Title, Priority: tasks.Medium, Tags: opts.Tags} // Default priority
		
		// Parse priority if provided
		if opts.Priority != "" {
			if priority, err := tasks.ParsePriority(opts.Priority); err != nil {
				fmt.Println("Error parsing priority:", err)
				os.Exit(1)
			} else {
				t.Priority = priority
			}
		}
		
		// Parse due date if provided
		if opts.Due != "" {
			if dueDate, err := tasks.ParseDueDate(opts.Due); err != nil {
				fmt.Println("Error parsing due date:", err)
				os.Exit(1)
			} else {
				t.DueDate = dueDate
			}
		}
		
		err := manager.Add(t)
		if err != nil {
			fmt.Println("Error adding task:", err)
			sentry.CaptureException(err)
			os.Exit(1)
		}
		fmt.Println("Task added.")
	case "list":
		opts := cli.ParseListCommand(args)
		var tasksToShow []tasks.Task
		
		// Apply filters
		if opts.Priority != "" {
			if priority, err := tasks.ParsePriority(opts.Priority); err != nil {
				fmt.Println("Error parsing priority:", err)
				os.Exit(1)
			} else {
				tasksToShow = manager.ListByPriority(priority)
			}
		} else if opts.Tag != "" {
			tasksToShow = manager.ListByTag(opts.Tag)
		} else if opts.Overdue {
			tasksToShow = manager.ListOverdue()
		} else if opts.DueToday {
			tasksToShow = manager.ListDueToday()
		} else if opts.DueWithin > 0 {
			tasksToShow = manager.ListDueWithin(opts.DueWithin)
		} else {
			tasksToShow = manager.List()
		}
		
		// Create display options
		displayOpts := display.DisplayOptions{
			ShowColors:   display.IsColorSupported(),
			ShowIcons:    true,
			TableFormat:  false,
			ShowTags:     true,
			ShowDueDate:  true,
			ShowPriority: true,
			ColorScheme:  display.DefaultColorScheme,
		}
		
		// Check for format flags
		for _, arg := range args {
			if arg == "--table" {
				displayOpts.TableFormat = true
			}
			if arg == "--no-color" {
				displayOpts.ShowColors = false
			}
			if arg == "--no-icons" {
				displayOpts.ShowIcons = false
			}
			if arg == "--minimal" {
				displayOpts.ShowTags = false
				displayOpts.ShowDueDate = false
				displayOpts.ShowPriority = false
				displayOpts.ShowIcons = false
			}
		}
		
		formatter := display.NewTaskFormatter(displayOpts)
		
		// Display table header if in table format
		if displayOpts.TableFormat {
			fmt.Println(formatter.FormatTableHeader())
			fmt.Println(formatter.FormatTableSeparator())
		}
		
		// Display tasks with enhanced formatting
		for i, t := range tasksToShow {
			fmt.Println(formatter.FormatTask(i, t))
		}
	case "done":
		if len(args) < 1 {
			fmt.Println("Usage: taskmgr done <index>")
			os.Exit(1)
		}
		err := manager.MarkDone(args[0])
		if err != nil {
			fmt.Println("Error marking done:", err)
			sentry.CaptureException(err)
			os.Exit(1)
		}
		fmt.Println("Task marked as done.")
	case "remove":
		if len(args) < 1 {
			fmt.Println("Usage: taskmgr remove <index>")
			os.Exit(1)
		}
		err := manager.Remove(args[0])
		if err != nil {
			fmt.Println("Error removing task:", err)
			os.Exit(1)
		}
		fmt.Println("Task removed.")
	case "undodone":
		if len(args) < 1 {
			fmt.Println("Usage: taskmgr undodone <index>")
			os.Exit(1)
		}
		err := manager.UndoDone(args[0])
		if err != nil {
			fmt.Println("Error undoing done:", err)
			os.Exit(1)
		}
		fmt.Println("Task marked as not done.")
	case "find":
		if len(args) < 1 {
			fmt.Println("Usage: taskmgr find <title>")
			os.Exit(1)
		}
		task := manager.FindByTitle(args[0])
		if task == nil {
			fmt.Println("No task found with that title.")
			os.Exit(0)
		}
		done := " "
		if task.Done {
			done = "x"
		}
		fmt.Printf("[%s] %s\n", done, task.Title)
	case "bulkadd":
		if len(args) < 1 {
			fmt.Println("Usage: taskmgr bulkadd <title1,title2,...>")
			os.Exit(1)
		}
		titles := strings.Split(args[0], ",")
		tasksToAdd := make([]tasks.Task, len(titles))
		for i, title := range titles {
			tasksToAdd[i] = tasks.Task{Title: strings.TrimSpace(title)}
		}
		err := manager.BulkAdd(tasksToAdd)
		if err != nil {
			fmt.Println("Error bulk adding tasks:", err)
			os.Exit(1)
		}
		fmt.Println("Tasks added.")
	case "countdone":
		count := manager.CountDone()
		fmt.Printf("Completed tasks: %d\n", count)
	case "markall":
		err := manager.MarkAllDone()
		if err != nil {
			fmt.Println("Error marking all tasks done:", err)
			os.Exit(1)
		}
		fmt.Println("All tasks marked as done.")
	case "findbydesc":
		if len(args) < 1 {
			fmt.Println("Usage: taskmgr findbydesc <description>")
			os.Exit(1)
		}
		results := manager.FindByDescription(args[0])
		if len(results) == 0 {
			fmt.Println("No tasks found with that description.")
			os.Exit(0)
		}
		for i, t := range results {
			done := " "
			if t.Done {
				done = "x"
			}
			fmt.Printf("%d: [%s] %s\n", i, done, t.Title)
		}
	case "stats":
		tasks := manager.List()
		displayOpts := display.DisplayOptions{
			ShowColors:   display.IsColorSupported(),
			ShowIcons:    true,
			ColorScheme:  display.DefaultColorScheme,
		}
		
		// Check for no-color flag
		for _, arg := range args {
			if arg == "--no-color" {
				displayOpts.ShowColors = false
			}
		}
		
		progressFormatter := display.NewProgressFormatter(displayOpts)
		stats := progressFormatter.CalculateStats(tasks)
		fmt.Println(progressFormatter.FormatDetailedStats(stats))
	case "tags":
		allTags := manager.GetAllTags()
		if len(allTags) == 0 {
			fmt.Println("No tags found.")
			return
		}
		fmt.Println("Available tags:")
		for _, tag := range allTags {
			fmt.Printf("  - %s\n", tag)
		}
	case "tag":
		if len(args) < 2 {
			fmt.Println("Usage: taskmgr tag <index> <tag>")
			os.Exit(1)
		}
		err := manager.AddTagToTask(args[0], args[1])
		if err != nil {
			fmt.Println("Error adding tag:", err)
			os.Exit(1)
		}
		fmt.Printf("Tag '%s' added to task.\n", args[1])
	case "untag":
		if len(args) < 2 {
			fmt.Println("Usage: taskmgr untag <index> <tag>")
			os.Exit(1)
		}
		err := manager.RemoveTagFromTask(args[0], args[1])
		if err != nil {
			fmt.Println("Error removing tag:", err)
			os.Exit(1)
		}
		fmt.Printf("Tag '%s' removed from task.\n", args[1])
	case "error":
		// Trigger an error to test sentry.
		err := errors.New("test error. create gh issue?")
		sentry.CaptureException(err)
	default:
		fmt.Println("Usage: taskmgr [command] ...")
		fmt.Println("Available commands:")
		fmt.Println("  add <title> [--priority=<low|medium|high|critical>] [--due=<date>] [--tags=<tag1,tag2,...>]")
		fmt.Println("                         - Add a new task with optional priority, due date, and tags")
		fmt.Println("  list [filters] [options] - List tasks with optional filters and formatting")
		fmt.Println("    Filters:")
		fmt.Println("      --priority=<priority>  - Filter by priority level")
		fmt.Println("      --tag=<tag>            - Filter by tag")
		fmt.Println("      --overdue              - Show only overdue tasks")
		fmt.Println("      --due-today            - Show tasks due today")
		fmt.Println("      --due-within=<days>    - Show tasks due within N days")
		fmt.Println("    Display Options:")
		fmt.Println("      --table                - Display in table format")
		fmt.Println("      --no-color             - Disable colored output")
		fmt.Println("      --no-icons             - Disable emoji icons")
		fmt.Println("      --minimal              - Minimal output (no colors, icons, or extra info)")
		fmt.Println("  stats [--no-color]      - Show progress statistics and task breakdown")
		fmt.Println("  tags                     - List all available tags")
		fmt.Println("  tag <index> <tag>        - Add a tag to an existing task")
		fmt.Println("  untag <index> <tag>      - Remove a tag from a task")
		fmt.Println("  done <index>             - Mark a task as done")
		fmt.Println("  remove <index>           - Remove a task")
		fmt.Println("  undodone <index>         - Mark a completed task as not done")
		fmt.Println("  find <title>             - Find task by title")
		fmt.Println("  bulkadd <t1,t2,...>      - Add multiple tasks at once")
		fmt.Println("  countdone                - Count completed tasks")
		fmt.Println("  markall                  - Mark all tasks as done")
		fmt.Println("  findbydesc <desc>        - Find tasks by description")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  taskmgr add \"Fix bug\" --priority=high --due=2024-01-15 --tags=work,urgent")
		fmt.Println("  taskmgr add \"Review PR\" --priority=medium --due=tomorrow --tags=work,code-review")
		fmt.Println("  taskmgr add \"Buy groceries\" --tags=personal,shopping")
		fmt.Println("  taskmgr list --priority=high --table")
		fmt.Println("  taskmgr list --tag=work")
		fmt.Println("  taskmgr list --overdue --no-color")
		fmt.Println("  taskmgr list --due-today")
		fmt.Println("  taskmgr list --due-within=7days --minimal")
		fmt.Println("  taskmgr tags")
		fmt.Println("  taskmgr tag 0 urgent")
		fmt.Println("  taskmgr untag 0 urgent")
		fmt.Println("  taskmgr stats")
		os.Exit(1)
	}
}
