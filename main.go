package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
	"github.com/go-sql-driver/mysql"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/speaker"
	"github.com/machinebox/graphql"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"time"
)

var db *sql.DB

func EstablishConnection() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "MusicPlayer",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

var request = graphql.NewRequest(``)

var mFilePath = "C:\\Users\\Asus\\Music\\MusicPlayer\\Vangelis - Nocturne (2019) [24Bit Hi-Res]\\10 - To a Friend.flac"
var mFile *os.File
var streamer beep.StreamSeekCloser
var format beep.Format
var meta tag.Metadata

func main() {
	EstablishConnection()

	//client := graphql.NewClient("http://localhost:3000/graphbrainz")
	//go func() {
	//	api.MetadataRequest(client, request)
	//}()

	//go func() {
	//	art, _ := caa.GetReleaseGroup("5e62fe2c-28d3-3508-9a5c-b6fe2fd4ca7a")
	//	fmt.Println(art)
	//}()

	// Load the FLAC file
	mFile, _ = os.Open(mFilePath)

	meta, _ = tag.ReadFrom(mFile)
	if _, err := mFile.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	}

	streamer, format, _ = flac.Decode(mFile)
	defer streamer.Close()

	// Initialize the speaker
	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		log.Fatal(err)
	}

	a := app.New()
	w := a.NewWindow("Music Player")

	// Menu bar
	menubar(&w, &a)

	// Main Content Area
	contentLabel := widget.NewLabel(meta.Album())
	mainContent := container.New(layout.NewStackLayout(), contentLabel)

	// Navigation Panel
	navPanel := createNavPanel(mainContent)

	// Player Controls
	playerControls := createPlayerControls()

	// Overall Layout
	mainLayout := container.NewBorder(nil, playerControls, navPanel, nil, mainContent)

	w.SetContent(mainLayout)
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}

var streamerWithVol *effects.Volume
var vol float64
var isSilent = false

// playAudio plays audio (call with go to play in a separate goroutine)
func playAudio() {
	donePlaying := make(chan bool)
	streamerWithVol = &effects.Volume{
		Streamer: streamer,
		Base:     math.E,
		Volume:   vol,
		Silent:   isSilent,
	}

	speaker.Play(beep.Seq(streamerWithVol, beep.Callback(func() {
		donePlaying <- true
	})))

	// Wait for the audio to finish playing
	<-donePlaying
}

func menubar(myWindow *fyne.Window, myApp *fyne.App) {
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

	// Create the main menu with the menus
	mainMenu := fyne.NewMainMenu(fileMenu, editMenu)

	// Set the main menu to the window
	(*myWindow).SetMainMenu(mainMenu)
}

func createNavPanel(mainContent *fyne.Container) *fyne.Container {
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

type fixedSizeLayout struct {
	leftLabel  fyne.CanvasObject
	rightLabel fyne.CanvasObject
}

func (f *fixedSizeLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	leftLabelSize := f.leftLabel.MinSize()
	rightLabelSize := f.rightLabel.MinSize()

	leftLabelOffset := fyne.NewPos(0, (size.Height-leftLabelSize.Height)/2) // Center left label vertically
	f.leftLabel.Resize(leftLabelSize)
	f.leftLabel.Move(leftLabelOffset)

	rightLabelOffset := fyne.NewPos(size.Width-rightLabelSize.Width, (size.Height-rightLabelSize.Height)/2) // Center right label vertically
	f.rightLabel.Resize(rightLabelSize)
	f.rightLabel.Move(rightLabelOffset)

	sliderOffset := fyne.NewPos(leftLabelSize.Width, 0)
	sliderSize := fyne.NewSize(size.Width-leftLabelSize.Width-rightLabelSize.Width, size.Height)
	objects[1].Resize(sliderSize)
	objects[1].Move(sliderOffset)
}

func (f *fixedSizeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	leftLabelSize := f.leftLabel.MinSize()
	rightLabelSize := f.rightLabel.MinSize()
	sliderMin := objects[1].MinSize()

	return fyne.NewSize(
		leftLabelSize.Width+sliderMin.Width+rightLabelSize.Width,
		fyne.Max(
			leftLabelSize.Height,
			fyne.Max(
				sliderMin.Height,
				rightLabelSize.Height,
			),
		),
	)
}

func newFixedSizeLayout(leftLabel, rightLabel fyne.CanvasObject) fyne.Layout {
	return &fixedSizeLayout{leftLabel: leftLabel, rightLabel: rightLabel}
}

func createPlayerControls() *fyne.Container {
	//// Stop button
	//stopButton := widget.NewButton("Stop", func() {
	//	speaker.Close()
	//})

	skipPrevious := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), nil)

	playPauseButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil)
	playPauseButton.OnTapped = func() {
		if playPauseButton.Icon == theme.MediaPlayIcon() {
			playPauseButton.SetIcon(theme.MediaPauseIcon())
			go playAudio()
		} else {
			playPauseButton.SetIcon(theme.MediaPlayIcon())
			speaker.Clear()
		}
		playPauseButton.Refresh()
	}

	skipNext := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), nil)

	mediaButtons := container.NewHBox(
		widget.NewButton("Loop", func() {}),
		skipPrevious,
		playPauseButton,
		skipNext,
		widget.NewButton("Shuffle", func() {}),
	)

	vol = 0.0
	volumeSlider := widget.NewSliderWithData(-2.5, 0, binding.BindFloat(&vol))
	volumeSlider.Step = (0 - (-2.5)) / 100
	volumeSlider.OnChanged = func(v float64) {
		if v == -2.5 {
			isSilent = true
		} else {
			isSilent = false
			vol = v
		}
		if playPauseButton.Icon == theme.MediaPauseIcon() {
			speaker.Clear()
			go playAudio()
		}
	}
	volumeSliderCont := container.NewHBox(volumeSlider)

	mySpacer := canvas.NewRectangle(color.Transparent)
	mySpacer.SetMinSize(volumeSlider.MinSize())

	mediaButtonsCen := container.NewHBox(volumeSliderCont, layout.NewSpacer(), mediaButtons, layout.NewSpacer(), mySpacer)

	duration := time.Second * time.Duration(streamer.Len()) / time.Duration(format.SampleRate)
	pos := 0.0
	playbackPos := binding.BindFloat(&pos)
	playbackSlider := widget.NewSliderWithData(0, duration.Seconds(), playbackPos)
	playbackSlider.Step = 0.01
	playbackSlider.Value = 0
	playbackPosLabel := widget.NewLabelWithData(binding.FloatToString(playbackPos))
	formattedDuration := formatTime(duration)
	playbackDurLabel := widget.NewLabel(formattedDuration)
	playbackCont := container.New(
		newFixedSizeLayout(playbackPosLabel, playbackDurLabel),
		playbackPosLabel,
		playbackSlider,
		playbackDurLabel,
	)

	mediaControls := container.NewVBox(
		mediaButtonsCen,
		playbackCont,
	)
	mediaControls.Resize(fyne.NewSize(600, 100))
	return mediaControls
}

func formatTime(duration time.Duration) string {
	totalSeconds := int(duration.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
