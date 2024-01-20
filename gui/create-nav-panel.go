package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var currentView = "Home"

var MainContent *fyne.Container

var homeButton, exploreButton, libraryButton *widget.Button

// UpdateContent simulates navigation by updating the main area
func UpdateContent(mainContent *fyne.Container, content fyne.CanvasObject) {
	if currentView == "Library" {
		fmt.Println("2")
	}
	mainContent.Objects = []fyne.CanvasObject{content}
	mainContent.Refresh()
}

func CreateNavPanel(mainContent *fyne.Container) *fyne.Container {
	UpdateContent(mainContent, CreateHomeView())

	homeButton = widget.NewButton("Home", nil)
	exploreButton = widget.NewButton("Explore", nil)
	libraryButton = widget.NewButton("Library", nil)

	homeButton.Importance = widget.HighImportance

	navPanel := container.NewVBox(
		homeButton,
		exploreButton,
		libraryButton,
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
	)

	homeButton.OnTapped = func() {
		if currentView == "Home" {
			return
		}

		UpdateContent(mainContent, CreateHomeView())

		currentView = "Home"
		homeButton.Importance = widget.HighImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
	}

	exploreButton.OnTapped = func() {
		if currentView == "Explore" {
			return
		}

		if LastLibraryViewAlbum.ID == 0 {
			UpdateContent(mainContent, widget.NewLabel("Select an album from the home view to see more about the artist."))
		} else {
			UpdateContent(mainContent, CreateExploreView())
		}

		currentView = "Explore"
		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.HighImportance
		libraryButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
	}

	libraryButton.OnTapped = func() {
		if currentView == "Library" {
			return
		}

		if LastLibraryViewAlbum.ID == 0 {
			UpdateContent(mainContent, widget.NewLabel("Select an album from the home view to see tracks here."))
		} else {
			UpdateContent(mainContent, CreateLibraryView(LastLibraryViewAlbum))
		}

		currentView = "Library"
		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.HighImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
	}

	return navPanel
}
