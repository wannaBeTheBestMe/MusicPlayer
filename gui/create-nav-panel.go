package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateNavPanel(mainContent *fyne.Container) *fyne.Container {
	navPanel := container.NewVBox(
		widget.NewButton("Home", nil),
		widget.NewButton("Explore", nil),
		widget.NewButton("Library", nil),
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
	)
	// Function to update main content (simulating navigation)
	updateContent := func(content fyne.CanvasObject) {
		mainContent.Objects = []fyne.CanvasObject{content}
		mainContent.Refresh()
	}

	// Link navigation buttons to update main content
	navPanel.Objects[0].(*widget.Button).OnTapped = func() {
		updateContent(widget.NewLabel("Home View"))
	}
	navPanel.Objects[1].(*widget.Button).OnTapped = func() {
		updateContent(widget.NewLabel("Explore View"))
	}
	navPanel.Objects[2].(*widget.Button).OnTapped = func() {
		updateContent(widget.NewLabel("Library View"))
	}

	return navPanel
}
