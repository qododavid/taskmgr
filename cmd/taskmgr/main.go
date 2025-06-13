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
			fmt.Println("Usage: taskmgr add <title> [--priority=<low|medium|high|critical>] [--due=<date>]")
			fmt.Println("Examples:")
			fmt.Println("  taskmgr add \"Fix bug\" --priority=high --due=2024-01-15")
			fmt.Println("  taskmgr add \"Review PR\" --priority=medium --due=tomorrow")
			os.Exit(1)
		}
		
		t := tasks.Task{Title: opts.Title, Priority: tasks.Medium} // Default priority
		
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
		} else if opts.Overdue {
			tasksToShow = manager.ListOverdue()
		} else if opts.DueToday {
			tasksToShow = manager.ListDueToday()
		} else if opts.DueWithin > 0 {
			tasksToShow = manager.ListDueWithin(opts.DueWithin)
		} else {
			tasksToShow = manager.List()
		}
		
		// Display tasks with enhanced formatting
		for i, t := range tasksToShow {
			done := " "
			if t.Done {
				done = "x"
			}
			
			// Format priority with color
			priorityStr := fmt.Sprintf("%s[%s]%s", t.Priority.Color(), strings.ToUpper(t.Priority.String()[:3]), t.Priority.ColorReset())
			
			// Format due date
			dueDateStr := ""
			if t.DueDate != nil {
				now := time.Now()
				if t.DueDate.Before(now) && !t.Done {
					// Overdue - red
					dueDateStr = fmt.Sprintf(" \033[31m(Due: %s) ⚠️\033[0m", t.DueDate.Format("2006-01-02"))
				} else {
					dueDateStr = fmt.Sprintf(" (Due: %s)", t.DueDate.Format("2006-01-02"))
				}
			}
			
			fmt.Printf("%d: [%s] %s %s%s\n", i, done, priorityStr, t.Title, dueDateStr)
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
	case "error":
		// Trigger an error to test sentry.
		err := errors.New("test error. create gh issue?")
		sentry.CaptureException(err)
	default:
		fmt.Println("Usage: taskmgr [command] ...")
		fmt.Println("Available commands:")
		fmt.Println("  add <title> [--priority=<low|medium|high|critical>] [--due=<date>]")
		fmt.Println("                         - Add a new task with optional priority and due date")
		fmt.Println("  list [--priority=<priority>] [--overdue] [--due-today] [--due-within=<days>]")
		fmt.Println("                         - List tasks with optional filters")
		fmt.Println("  done <index>           - Mark a task as done")
		fmt.Println("  remove <index>         - Remove a task")
		fmt.Println("  undodone <index>       - Mark a completed task as not done")
		fmt.Println("  find <title>           - Find task by title")
		fmt.Println("  bulkadd <t1,t2,...>    - Add multiple tasks at once")
		fmt.Println("  countdone              - Count completed tasks")
		fmt.Println("  markall                - Mark all tasks as done")
		fmt.Println("  findbydesc <desc>      - Find tasks by description")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  taskmgr add \"Fix bug\" --priority=high --due=2024-01-15")
		fmt.Println("  taskmgr add \"Review PR\" --priority=medium --due=tomorrow")
		fmt.Println("  taskmgr list --priority=high")
		fmt.Println("  taskmgr list --overdue")
		fmt.Println("  taskmgr list --due-today")
		fmt.Println("  taskmgr list --due-within=7days")
		os.Exit(1)
	}
}
