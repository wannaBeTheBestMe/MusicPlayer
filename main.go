package main

import (
	"MusicPlayer/data_access"
	"MusicPlayer/gui"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/speaker"
	"github.com/machinebox/graphql"
	"log"
	"os"
	"time"
)

var request = graphql.NewRequest(``)

func main() {
	data_access.EstablishConnection()

	data_access.BatchAddTracks()

	data_access.GetAlbumDirectories()

	//client := graphql.NewClient("http://localhost:3000/graphbrainz")
	//go func() {
	//	api.MetadataRequest(client, request)
	//}()

	//go func() {
	//	art, _ := caa.GetReleaseGroup("5e62fe2c-28d3-3508-9a5c-b6fe2fd4ca7a")
	//	fmt.Println(art)
	//}()

	// Load the FLAC file
	gui.Mfile, _ = os.Open(gui.MfilePath)
	gui.Meta, _ = data_access.ReadTags(gui.Mfile, &gui.MfilePath)
	fmt.Println(gui.Meta.Genre())

	func() {
		gui.StreamerMu.Lock()
		defer gui.StreamerMu.Unlock()
		gui.Streamer, gui.Format, _ = flac.Decode(gui.Mfile)
	}()
	defer gui.Streamer.Close()

	// Initialize the speaker
	if err := speaker.Init(gui.Format.SampleRate, gui.Format.SampleRate.N(time.Second/10)); err != nil {
		log.Fatal(err)
	}

	a := app.New()
	w := a.NewWindow("Music Player")

	gui.Menubar(&w, &a)
	header := gui.CreateHeader()
	contentLabel := widget.NewLabel(gui.Meta.Album())
	mainContent := container.New(layout.NewStackLayout(), contentLabel)
	navPanel := gui.CreateNavPanel(mainContent)
	playerControls := gui.CreatePlayerControls()

	mainLayout := container.NewBorder(header, playerControls, navPanel, nil, mainContent)

	w.SetContent(mainLayout)
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}
