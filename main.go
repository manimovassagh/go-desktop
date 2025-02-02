package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("To-Do List")

	// Task list
	var tasks []string
	taskList := widget.NewList(
		func() int { return len(tasks) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(tasks[i])
		},
	)

	// Input field and add button
	taskEntry := widget.NewEntry()
	taskEntry.SetPlaceHolder("Enter new task...")

	addButton := widget.NewButton("Add Task", func() {
		if taskEntry.Text != "" {
			tasks = append(tasks, taskEntry.Text)
			taskEntry.SetText("")
			taskList.Refresh() // Update UI
		}
	})

	// Layout
	myWindow.SetContent(container.NewVBox(
		taskEntry,
		addButton,
		taskList,
	))

	myWindow.Resize(fyne.NewSize(400, 600)) // Set window size
	myWindow.ShowAndRun()
}