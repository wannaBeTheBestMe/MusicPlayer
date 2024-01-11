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

	//client := graphql.NewClient("http://localhost:3000/graphbrainz")
	//go func() {
	//	api.MetadataRequest(client, request)
	//}()

	//go func() {
	//	art, _ := caa.GetReleaseGroup("5e62fe2c-28d3-3508-9a5c-b6fe2fd4ca7a")
	//	fmt.Println(art)
	//}()

	// Load the FLAC file
	Mfile, _ := os.Open(gui.MfilePath)
	Meta, _ := data_access.ReadTags(Mfile, &gui.MfilePath)
	fmt.Println(Meta.Genre())

	Streamer, Format, _ := flac.Decode(Mfile)
	defer Streamer.Close()

	// Initialize the speaker
	if err := speaker.Init(Format.SampleRate, Format.SampleRate.N(time.Second/10)); err != nil {
		log.Fatal(err)
	}

	a := app.New()
	w := a.NewWindow("Music Player")

	gui.Menubar(&w, &a)
	header := gui.CreateHeader()
	contentLabel := widget.NewLabel(Meta.Album())
	mainContent := container.New(layout.NewStackLayout(), contentLabel)
	navPanel := gui.CreateNavPanel(mainContent)
	playerControls := gui.CreatePlayerControls(&Streamer, &Format)

	mainLayout := container.NewBorder(header, playerControls, navPanel, nil, mainContent)

	w.SetContent(mainLayout)
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}
