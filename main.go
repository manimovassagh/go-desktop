package main

import (
	"database/sql"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("To-Do List")

	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Use slice to store tasks
	var tasks []string

	// Function to load tasks from database into tasks slice
	loadTasks := func() {
		tasks = nil // Clear the previous list
		rows, err := db.Query("SELECT task FROM tasks")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var task string
			if err := rows.Scan(&task); err != nil {
				log.Fatal(err)
			}
			tasks = append(tasks, task)
		}
	}

	// Load initial tasks from the database
	loadTasks()

	// Task list with tasks from the tasks slice
	taskList := widget.NewList(
		func() int { return len(tasks) }, // Returns number of tasks
		func() fyne.CanvasObject { return widget.NewLabel("") }, // Creates new label for each task
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(tasks[i]) // Sets task text to the label
		},
	)

	// Input field and add button
	taskEntry := widget.NewEntry()
	taskEntry.SetPlaceHolder("Enter new task...")

	// Function to add task and persist it in the database
	addTask := func() {
		if taskEntry.Text != "" {
			// Insert the new task into the SQLite database
			_, err := db.Exec("INSERT INTO tasks (task) VALUES (?)", taskEntry.Text)
			if err != nil {
				log.Fatal(err)
			}

			// Reload the tasks from the database and update the list
			loadTasks()
			taskList.Refresh() // Refresh the UI list to show the updated tasks
			taskEntry.SetText("") // Clear the input field
		}
	}

	// Add button with better styling
	addButton := widget.NewButtonWithIcon("Add Task", theme.ConfirmIcon(), func() {
		addTask()
	})
	addButton.Importance = widget.HighImportance // Makes it more visually prominent

	// Allow "Enter" to add task
	taskEntry.OnSubmitted = func(string) {
		addTask()
	}

	// Layout with spacing
	myWindow.SetContent(container.NewVBox(
		taskEntry,
		addButton,
		taskList,
	))

	myWindow.Resize(fyne.NewSize(400, 600)) // Set window size
	myWindow.ShowAndRun()
}