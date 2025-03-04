package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var currentView = "Home"

var MainContent *fyne.Container

var homeButton, exploreButton, libraryButton, suggestButton, searchButton, playlistsButton *widget.Button

// UpdateContent simulates navigation by updating the main area
func UpdateContent(mainContent *fyne.Container, content fyne.CanvasObject) {
	mainContent.Objects = []fyne.CanvasObject{content}
	mainContent.Refresh()
}

func CreateNavPanel(mainContent *fyne.Container) *fyne.Container {
	UpdateContent(mainContent, CreateHomeView())

	homeButton = widget.NewButton("Home", nil)
	exploreButton = widget.NewButton("Explore", nil)
	libraryButton = widget.NewButton("Library", nil)
	suggestButton = widget.NewButton("Suggest", nil)
	searchButton = widget.NewButton("Search", nil)
	playlistsButton = widget.NewButton("Playlists", nil)

	homeButton.Importance = widget.HighImportance

	navPanel := container.NewVBox(
		homeButton,
		exploreButton,
		libraryButton,
		suggestButton,
		searchButton,
		playlistsButton,
		//widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
	)

	homeButton.OnTapped = func() {
		if currentView == "Home" {
			return
		}

		UpdateContent(mainContent, CreateHomeView())

		homeButton.Importance = widget.HighImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.MediumImportance
		suggestButton.Importance = widget.MediumImportance
		searchButton.Importance = widget.MediumImportance
		playlistsButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
		suggestButton.Refresh()
		searchButton.Refresh()
		playlistsButton.Refresh()
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

		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.HighImportance
		libraryButton.Importance = widget.MediumImportance
		suggestButton.Importance = widget.MediumImportance
		searchButton.Importance = widget.MediumImportance
		playlistsButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
		suggestButton.Refresh()
		searchButton.Refresh()
		playlistsButton.Refresh()

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

		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.HighImportance
		suggestButton.Importance = widget.MediumImportance
		searchButton.Importance = widget.MediumImportance
		playlistsButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
		suggestButton.Refresh()
		searchButton.Refresh()
		playlistsButton.Refresh()

	}

	suggestButton.OnTapped = func() {
		if currentView == "Suggest" {
			return
		}

		UpdateContent(mainContent, CreateSuggestView())

		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.MediumImportance
		suggestButton.Importance = widget.HighImportance
		searchButton.Importance = widget.MediumImportance
		playlistsButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
		suggestButton.Refresh()
		searchButton.Refresh()
		playlistsButton.Refresh()

	}

	searchButton.OnTapped = func() {
		if currentView == "Search" {
			return
		}

		UpdateContent(mainContent, CreateSearchView(mainContent))

		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.MediumImportance
		suggestButton.Importance = widget.MediumImportance
		searchButton.Importance = widget.HighImportance
		playlistsButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
		suggestButton.Refresh()
		searchButton.Refresh()
		playlistsButton.Refresh()
	}

	playlistsButton.OnTapped = func() {
		if currentView == "Search" {
			return
		}

		UpdateContent(mainContent, CreatePlaylistsViews())

		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.MediumImportance
		suggestButton.Importance = widget.MediumImportance
		searchButton.Importance = widget.MediumImportance
		playlistsButton.Importance = widget.HighImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
		suggestButton.Refresh()
		searchButton.Refresh()
		playlistsButton.Refresh()
	}

	return navPanel
}
