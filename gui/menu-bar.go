package gui

import (
	"MusicPlayer/data_access"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func Menubar(myWindow *fyne.Window, myApp *fyne.App) {
	// Define actions for menu items
	newFile := func() { dialog.ShowInformation("New", "Create a new file", *myWindow) }
	openFile := func() { dialog.ShowInformation("Open", "Open a file", *myWindow) }
	saveFile := func() { dialog.ShowInformation("Save", "Save the file", *myWindow) }
	quitApp := func() { (*myApp).Quit() }

	// Create the menu items
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("New", newFile),
		fyne.NewMenuItem("Open", openFile),
		fyne.NewMenuItem("Save", saveFile),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", quitApp),
	)

	// Add more menus as needed
	editMenu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Cut", nil),
		fyne.NewMenuItem("Copy", nil),
		fyne.NewMenuItem("Paste", nil),
	)

	importMenu := fyne.NewMenu("Import",
		fyne.NewMenuItem("Import tracks to database",
			func() { data_access.BatchAddTracks("C:\\Users\\Asus\\Music\\MusicPlayer") }),
	)

	// Create the main menu with the menus
	mainMenu := fyne.NewMainMenu(fileMenu, editMenu, importMenu)

	// Set the main menu to the window
	(*myWindow).SetMainMenu(mainMenu)
}
