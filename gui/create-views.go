package gui

import (
	"MusicPlayer/api"
	"MusicPlayer/data_access"
	"MusicPlayer/playback"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var APIResp api.ArtistInfoQuery
var APIRespChan = make(chan bool)

var LastLibraryViewAlbum data_access.Album

func CreateHomeView() *container.Scroll {
	currentView = "Home"
	loaded := data_access.HomeAlbumsLoaded
	if !loaded {
		for !loaded {
		}
	}

	cardGrid := container.New(layout.NewGridWrapLayout(fyne.NewSize(460, 460)))

	for _, alb := range data_access.HomeAlbumsArr {
		albumCard := createImgCard(alb)
		cardGrid.Add(albumCard)
	}
	//cardGridScroll.SetMinSize()

	return container.NewScroll(cardGrid)
}

func newLabelWithFixedWidth(text string, width float32) *fyne.Container {
	label := widget.NewLabel(text)
	//label.Wrapping = fyne.TextTruncate // Ensures the text does not exceed one line
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
			//fmt.Println(path)
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

// fetchImage fetches an image from a URL and returns a canvas.Image.
func fetchImage(url string, imgChan chan<- *canvas.Image) {
	go func() {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to fetch the image from URL '%s': %v", url, err)
			imgChan <- nil
			return
		}
		defer resp.Body.Close()

		imageData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read the image data from URL '%s': %v", url, err)
			imgChan <- nil
			return
		}

		imgResource := fyne.NewStaticResource("image", imageData)
		img := canvas.NewImageFromResource(imgResource)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(200, 200)) // Set a reasonable size for the image

		imgChan <- img
	}()
}

func CreateExploreView() *container.Scroll {
	for {
		select {
		case <-APIRespChan:
			mainAccordion := widget.NewAccordion()

			for _, edge := range APIResp.Search.Artists.Edges {
				artist := edge.Node

				artistInfo := container.NewVBox(
					//widget.NewLabel("MBID: "+artist.MBID),
					//widget.NewLabel("Name: "+artist.Name),
					//widget.NewLabel("LifeSpan: "+artist.LifeSpan.Begin+" - "+artist.LifeSpan.End),
					widget.NewLabel("Country: "+artist.Country),
					widget.NewLabel("Disambiguation: "+artist.Disambiguation),
					//widget.NewLabel("Biography: "+artist.TheAudioDB.Biography),
					widget.NewLabel("Type: "+artist.Type),
				)

				// Nested Accordion for Tags
				tagList := container.NewVBox()
				for _, tag := range artist.Tags.Nodes {
					tagList.Add(widget.NewLabel(fmt.Sprintf("%s (%d)", tag.Name, tag.Count)))
				}
				tagsAccordion := widget.NewAccordion(widget.NewAccordionItem("Tags", tagList))

				// Nested Accordion for Release Groups
				releaseGroupGrid := container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 300)))
				for _, rg := range artist.ReleaseGroups.Edges {
					titleLab := widget.NewLabel(rg.Node.Title)
					titleLab.Wrapping = fyne.TextWrapWord
					titleLab.Alignment = fyne.TextAlignCenter
					releaseGroup := container.NewVBox(titleLab)

					if rg.Node.CoverArtArchive.Artwork && rg.Node.CoverArtArchive.Front != "" {
						imgChan := make(chan *canvas.Image)
						fetchImage(rg.Node.CoverArtArchive.Front, imgChan)

						go func() {
							img := <-imgChan
							if img != nil {
								releaseGroup.Add(img)
								releaseGroup.Refresh()
							}
						}()
					}

					releaseGroupGrid.Add(releaseGroup)
				}
				releaseGroupsAccordion := widget.NewAccordion(widget.NewAccordionItem("Release Groups", releaseGroupGrid))

				// Add the nested accordions to the artist info
				artistInfo.Add(tagsAccordion)
				artistInfo.Add(releaseGroupsAccordion)

				// Add the artist info to the main accordion
				mainAccordion.Append(widget.NewAccordionItem(artist.Name, artistInfo))
			}

			return container.NewScroll(mainAccordion)
		default:
			return container.NewScroll(widget.NewLabel("Still fetching data..."))
		}
	}
}

func createImgCard(alb data_access.Album) *fyne.Container {
	reader := bytes.NewReader(alb.PictureData)
	img := canvas.NewImageFromReader(reader, "image.jpg")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(350, 350))

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

		go func() {
			fmt.Println("1")
			var err error
			albumArtist := LastLibraryViewAlbum.AlbumArtist
			fmt.Println(albumArtist)
			pos := strings.LastIndex(albumArtist, " &")
			if pos != -1 {
				firstAlbumArtist := albumArtist[:pos]
				albumArtist = firstAlbumArtist
				fmt.Println(firstAlbumArtist)
			}
			APIResp, err = api.GetArtistInfo(albumArtist)
			if err != nil {
				log.Fatal(err)
			}
			APIRespChan <- true
		}()
	})
	button.Importance = widget.LowImportance
	button.SetIcon(nil)

	overlay := container.NewStack(button, card)

	return overlay
}
