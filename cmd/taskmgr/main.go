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
		if len(args) < 1 {
			fmt.Println("Usage: taskmgr add <title>")
			os.Exit(1)
		}
		t := tasks.Task{Title: args[0]}
		err := manager.Add(t)
		if err != nil {
			fmt.Println("Error adding task:", err)
			sentry.CaptureException(err)
			os.Exit(1)
		}
		fmt.Println("Task added.")
	case "list":
		all := manager.List()
		for i, t := range all {
			done := " "
			if t.Done {
				done = "x"
			}
			fmt.Printf("%d: [%s] %s\n", i, done, t.Title)
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
		fmt.Println("  add <title>            - Add a new task")
		fmt.Println("  list                   - List all tasks")
		fmt.Println("  done <index>           - Mark a task as done")
		fmt.Println("  remove <index>         - Remove a task")
		fmt.Println("  undodone <index>       - Mark a completed task as not done")
		fmt.Println("  find <title>           - Find task by title")
		fmt.Println("  bulkadd <t1,t2,...>    - Add multiple tasks at once")
		fmt.Println("  countdone              - Count completed tasks")
		fmt.Println("  markall                - Mark all tasks as done")
		fmt.Println("  findbydesc <desc>      - Find tasks by description")
		os.Exit(1)
	}
}
