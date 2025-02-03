package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/therecipe/qt/widgets"
	  "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to SQLite
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task TEXT NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Load tasks from the database
	var tasks []string
	loadTasks := func() {
		tasks = nil
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

	// Initialize the Qt application
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Window Setup
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Todo List")
	window.SetMinimumSize2(400, 600)

	// Layout
	layout := widgets.NewQVBoxLayout()

	// Task Entry
	taskEntry := widgets.NewQLineEdit(nil)
	layout.AddWidget(taskEntry, 0, 0)

	// Task List
	taskList := widgets.NewQListWidget(nil)
	layout.AddWidget(taskList, 0, 0)

	// Load tasks
	loadTasks()
	for _, task := range tasks {
		taskList.AddItem(task)
	}

	// Add Button
	addButton := widgets.NewQPushButton2("Add Task", nil)
	addButton.ConnectClicked(func(bool) {
		if taskEntry.Text() != "" {
			// Insert the new task into the database
			_, err := db.Exec("INSERT INTO tasks (task) VALUES (?)", taskEntry.Text())
			if err != nil {
				log.Fatal(err)
			}

			// Reload tasks and update UI
			loadTasks()
			taskList.Clear()
			for _, task := range tasks {
				taskList.AddItem(task)
			}
			taskEntry.Clear()
		}
	})
	layout.AddWidget(addButton, 0, 0)

	// Set Layout and Show Window
	container := widgets.NewQWidget(nil, 0)
	container.SetLayout(layout)
	window.SetCentralWidget(container)

	window.Show()
	app.Exec()
}
