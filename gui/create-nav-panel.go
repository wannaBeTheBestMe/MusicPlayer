package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateNavPanel(mainContent *fyne.Container) *fyne.Container {
	currentView := ""

	// Function to update main content (simulating navigation)
	updateContent := func(content fyne.CanvasObject) {
		mainContent.Objects = []fyne.CanvasObject{content}
		mainContent.Refresh()
	}

	updateContent(CreateCardGridScroll())

	navPanel := container.NewVBox(
		widget.NewButton("Home", func() {
			if currentView != "Home" {
				currentView = "Home"
				updateContent(CreateCardGridScroll())
			}
		}),
		widget.NewButton("Explore", func() {
			if currentView != "Explore" {
				currentView = "Explore"
				updateContent(widget.NewLabel("Explore View"))
			}
		}),
		widget.NewButton("Library", func() {
			if currentView != "Library" {
				currentView = "Library"
				updateContent(widget.NewLabel("Library View"))
			}
		}),
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
	)

	return navPanel
}
