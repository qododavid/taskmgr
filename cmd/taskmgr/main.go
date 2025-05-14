package main

import (
	"errors"
	"fmt"
	"log"
	"os"
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
	case "error":
		// Trigger an error to test sentry.
		err := errors.New("test error. create gh issue?")
		sentry.CaptureException(err)
	default:
		fmt.Println("Usage: taskmgr [add|list|done] ...")
		os.Exit(1)
	}
}
