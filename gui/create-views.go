package gui

import (
	"MusicPlayer/data_access"
	"MusicPlayer/playback"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"time"
)

var LastLibraryViewAlbum data_access.Album

func CreateHomeView() *container.Scroll {
	currentView = "Home"
	loaded := data_access.HomeAlbumsLoaded
	if !loaded {
		for !loaded {
		}
	}

	cardGrid := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 400)))

	for _, alb := range data_access.HomeAlbumsArr {
		albumCard := createImgCard(alb)
		cardGrid.Add(albumCard)
	}
	//cardGridScroll.SetMinSize()

	return container.NewScroll(cardGrid)
}

func newLabelWithFixedWidth(text string, width float32) *fyne.Container {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextTruncate // Ensures the text does not exceed one line
	return container.NewGridWrap(fyne.NewSize(width, label.MinSize().Height), label)
}

func CreateLibraryView(alb data_access.Album) *container.Scroll {
	currentView = "Library"
	tracks, err := data_access.GetTracksInAlbum(alb)
	if err != nil {
		log.Fatal(err)
	}

	trackList := widget.NewList(
		func() int {
			return len(tracks)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				newLabelWithFixedWidth("", 500), // Title
				newLabelWithFixedWidth("", 400), // Artist
				newLabelWithFixedWidth("", 100), // Year
				newLabelWithFixedWidth("", 200), // Genre
				newLabelWithFixedWidth("", 100), // Duration
			)
		},
		func(num widget.ListItemID, item fyne.CanvasObject) {
			track := tracks[num]
			hbox := item.(*fyne.Container)
			hbox.Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(strconv.Itoa(num+1) + ". " + track.Title)
			hbox.Objects[1].(*fyne.Container).Objects[0].(*widget.Label).SetText(track.Artist)
			hbox.Objects[2].(*fyne.Container).Objects[0].(*widget.Label).SetText(strconv.Itoa(track.Year))
			hbox.Objects[3].(*fyne.Container).Objects[0].(*widget.Label).SetText(track.Genre)

			//fmt.Println(track)
			path := fmt.Sprintf("%s", track.Filepath)
			fmt.Println(path)
			duration := time.Duration(playback.GetTotalLengthFromPath(path)) * time.Second
			formattedDuration := FormatTime(duration)
			hbox.Objects[4].(*fyne.Container).Objects[0].(*widget.Label).SetText(formattedDuration)
		},
	)

	trackList.OnSelected = func(num widget.ListItemID) {
		track := tracks[num]
		playback.PauseAudio(playback.PDevice)
		playback.CleanCP()
		go playback.PlayAudio(fmt.Sprintf("%s", track.Filepath))
		playback.PauseAudio(playback.PDevice)
		playback.ResumeAudio(playback.PDevice)
	}

	return container.NewScroll(trackList)
}

func createImgCard(alb data_access.Album) *fyne.Container {
	reader := bytes.NewReader(alb.PictureData)
	img := canvas.NewImageFromReader(reader, "image.jpg")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(250, 250))

	titleLab := widget.NewLabel(alb.Title)
	titleLab.TextStyle = fyne.TextStyle{Bold: true}
	titleLab.Wrapping = fyne.TextWrapWord
	titleLab.Alignment = fyne.TextAlignCenter

	card := container.NewVBox(titleLab, widget.NewSeparator(), img)

	button := widget.NewButton("", func() {
		UpdateContent(MainContent, CreateLibraryView(alb))
		LastLibraryViewAlbum = alb

		currentView = "Library"
		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.HighImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
	})
	button.Importance = widget.LowImportance
	button.SetIcon(nil)

	overlay := container.NewStack(button, card)

	return overlay
}
