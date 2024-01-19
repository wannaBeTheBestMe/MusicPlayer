package main

import (
	"MusicPlayer/data_access"
	"MusicPlayer/gui"
	"MusicPlayer/playback"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/machinebox/graphql"
)

var request = graphql.NewRequest(``)

func main() {
	data_access.EstablishConnection()
	go data_access.LoadHomeAlbums()

	fileName := "C:\\Users\\Asus\\Music\\MusicPlayer\\Gladiator Soundtrack\\17. Now We Are Free.flac"
	go playback.PlayAudio(fileName)

	//client := graphql.NewClient("http://localhost:3000/graphbrainz")
	//go func() {
	//	api.MetadataRequest(client, request)
	//}()

	//go func() {
	//	art, _ := caa.GetReleaseGroup("5e62fe2c-28d3-3508-9a5c-b6fe2fd4ca7a")
	//	fmt.Println(art)
	//}()

	a := app.New()
	w := a.NewWindow("Music Player")

	gui.Menubar(&w, &a)
	header := gui.CreateHeader()
	gui.MainContent = container.New(layout.NewStackLayout())
	navPanel := gui.CreateNavPanel(gui.MainContent)
	playerControls := gui.CreatePlayerControls()

	mainLayout := container.NewBorder(header, playerControls, navPanel, nil, gui.MainContent)

	w.SetContent(mainLayout)
	w.Resize(fyne.NewSize(960, 540))

	w.ShowAndRun()
}
