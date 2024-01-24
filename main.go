package main

import (
	"MusicPlayer/data_access"
	"MusicPlayer/gui"
	app_icon "MusicPlayer/icon"
	"MusicPlayer/playback"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func main() {
	data_access.EstablishConnection()
	go data_access.LoadHomeAlbums()

	fileName := "C:\\Users\\Asus\\Music\\MusicPlayer\\Gladiator Soundtrack\\17. Now We Are Free.flac"
	go playback.PlayAudio(fileName)

	a := app.New()
	w := a.NewWindow("MusicPlayer")

	icon, _ := app_icon.LoadResourceFromPath("icon/icon2.png")
	w.SetIcon(icon)

	gui.Menubar(&w, &a)
	header := gui.CreateHeader()
	gui.MainContent = container.New(layout.NewStackLayout())
	navPanel := gui.CreateNavPanel(gui.MainContent)
	playerControls := gui.CreatePlayerControls(&w)

	mainLayout := container.NewBorder(header, playerControls, navPanel, nil, gui.MainContent)

	w.SetContent(mainLayout)
	w.Resize(fyne.NewSize(800*1.3, 450*1.3))

	w.ShowAndRun()
}
