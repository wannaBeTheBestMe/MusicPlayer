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
	"golang.org/x/image/font/opentype"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"image/color"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var APIResp api.ArtistInfoQuery
var APIRespChan = make(chan bool)

var LastLibraryViewAlbum data_access.Album
var currentTrack data_access.Track

var trackList *widget.List

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

	trackList = getListOfTracks(tracks)

	return container.NewScroll(trackList)
}

func getListOfTracks(tracks []data_access.Track) *widget.List {
	listOfTracks := widget.NewList(
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

	listOfTracks.OnSelected = func(num widget.ListItemID) {
		track := tracks[num]
		currentTrack = track
		//fmt.Println(currentTrack)

		playback.PauseAudio(playback.PDevice)
		playback.CleanCP()
		go playback.PlayAudio(fmt.Sprintf("%s", track.Filepath))
		//playback.PauseAudio(playback.PDevice)
		playback.ResumeAudio(playback.PDevice)

		err := data_access.IncrementTrackFreq(track)
		if err != nil {
			log.Fatal(err)
		}
	}

	return listOfTracks
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
	currentView = "Explore"

	for {
		select {
		case <-APIRespChan:
			if APIResp.Search.Artists.Edges[0].Node.MBID == "" {
				return container.NewScroll(widget.NewLabel("No information was found about this artist."))
			}

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
	if alb.PictureData == nil {
		return &fyne.Container{}
	}
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
		suggestButton.Importance = widget.MediumImportance
		searchButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
		suggestButton.Refresh()
		searchButton.Refresh()

		go func() {
			var err error
			albumArtist := LastLibraryViewAlbum.AlbumArtist
			pos := strings.LastIndex(albumArtist, " &")
			if pos != -1 {
				firstAlbumArtist := albumArtist[:pos]
				albumArtist = firstAlbumArtist
			}
			APIResp, err = api.GetArtistInfo(albumArtist)
			if err != nil {
				APIResp = api.ArtistInfoQuery{}
				//log.Fatal(err)
			}
			APIRespChan <- true
		}()
	})
	button.Importance = widget.LowImportance
	button.SetIcon(nil)

	overlay := container.NewStack(button, card)

	return overlay
}

type FieldFreq []float64

func (ff FieldFreq) Len() int {
	return len(ff)
}

func (ff FieldFreq) Value(i int) float64 {
	if i >= 0 && i < len(ff) {
		return ff[i]
	}
	return 0
}

func PlotTrend(x []string, y []float64, name string) *canvas.Image {
	p := plot.New()

	p.BackgroundColor = color.RGBA{R: 23, G: 23, B: 24, A: 255}
	//p.Title.Text = "Bar chart"
	//p.Title.TextStyle.Color = color.White
	//p.X.Label.Text = "Genres"
	//p.X.Label.TextStyle.Color = color.White
	//p.Y.Label.Text = "Heights"
	//p.Y.Label.TextStyle.Color = color.White
	p.Title.TextStyle.Font.Size = vg.Points(30)
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.X.Tick.Label.Font.Size = vg.Points(12)
	p.X.Tick.Label.Color = color.White
	p.X.Tick.Label.Rotation = 45
	p.X.Tick.Label.XAlign = draw.XRight
	p.X.Tick.Label.YAlign = draw.YCenter
	p.Y.Tick.Label.Font.Size = vg.Points(12)
	p.Y.Tick.Color = color.White
	p.Y.Tick.Label.Color = color.White
	p.Y.LineStyle.Color = color.White

	barsA, err := plotter.NewBarChart(FieldFreq(y), vg.Points(20))
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = color.RGBA{R: 54, G: 119, B: 246, A: 255}

	p.Add(barsA)
	p.NominalX(x...)

	width := 10 * vg.Inch
	height := 6 * vg.Inch
	dpi := 300

	c := vgimg.NewWith(vgimg.UseWH(width, height), vgimg.UseDPI(dpi))
	dc := draw.New(c)
	p.Draw(dc)

	w, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	pngimg := vgimg.PngCanvas{Canvas: c}
	if _, err := pngimg.WriteTo(w); err != nil {
		panic(err)
	}

	chartImg := canvas.NewImageFromFile(name)
	chartImg.FillMode = canvas.ImageFillContain
	chartImg.SetMinSize(fyne.NewSize(512, 512))

	return chartImg
}

func CreateSuggestView() *container.Scroll {
	currentView = "Suggest"

	f, _ := os.ReadFile("NotoSans-VariableFont_wdth,wght.ttf")
	ttf, _ := opentype.Parse(f)
	noto := font.Font{Typeface: "Noto Sans"}
	font.DefaultCache.Add([]font.Face{
		{
			Font: noto,
			Face: ttf,
		},
	})
	if !font.DefaultCache.Has(noto) {
		log.Fatalf("no font %q!", noto.Typeface)
	}
	plot.DefaultFont = noto
	plotter.DefaultFont = noto

	genre, genreFreq, err := data_access.GetFreqByGenre()
	if err != nil {
		log.Fatal(err)
	}
	genreImg := PlotTrend(genre, genreFreq, "genre.png")

	artist, artistFreq, err := data_access.GetFreqByArtist()
	if err != nil {
		log.Fatal(err)
	}
	artistImg := PlotTrend(artist, artistFreq, "artist.png")

	album, albumFreq, err := data_access.GetFreqByAlbum()
	if err != nil {
		log.Fatal(err)
	}
	albumImg := PlotTrend(album, albumFreq, "album.png")

	genreTitle := canvas.NewText("Listens on Top 10 Genres", color.White)
	genreTitle.TextStyle.Bold = true
	genreTitle.TextSize = 24
	genreTitle.Alignment = fyne.TextAlignCenter

	artistTitle := canvas.NewText("Listens on Top 10 Artists", color.White)
	artistTitle.TextStyle.Bold = true
	artistTitle.TextSize = 24
	artistTitle.Alignment = fyne.TextAlignCenter

	albumTitle := canvas.NewText("Listens on Top 10 Albums", color.White)
	albumTitle.TextStyle.Bold = true
	albumTitle.TextSize = 24
	albumTitle.Alignment = fyne.TextAlignCenter

	genreTop, err := data_access.GetTopByGenreFreq()
	if err != nil {
		log.Fatal(err)
	}
	genreText := fmt.Sprintf("**%v** is the genre you listen to the most. Give **%v** by **%v** another listen!", genreTop.Genre, genreTop.Title, genreTop.Artist)
	suggestGenre := widget.NewRichTextFromMarkdown(genreText)

	artistTop, err := data_access.GetTopByArtistFreq()
	if err != nil {
		log.Fatal(err)
	}
	artistText := fmt.Sprintf("**%v** is the artist you've listened to the most. Give **%v** by **%v** another listen!", artistTop.Artist, artistTop.Title, artistTop.Artist)
	suggestArtist := widget.NewRichTextFromMarkdown(artistText)

	albumTop, err := data_access.GetTopByAlbumFreq()
	if err != nil {
		log.Fatal(err)
	}
	albumText := fmt.Sprintf("**%v** is the album you've listened to the most. Give **%v** by **%v** another listen!", albumTop.Album, albumTop.Title, albumTop.Artist)
	suggestAlbum := widget.NewRichTextFromMarkdown(albumText)

	cont := container.NewVBox(
		genreTitle, genreImg, container.NewCenter(suggestGenre),
		artistTitle, artistImg, container.NewCenter(suggestArtist),
		albumTitle, albumImg, container.NewCenter(suggestAlbum),
	)

	return container.NewScroll(cont)
}

func CreateSearchView(mainContent *fyne.Container) *fyne.Container {
	currentView = "Search"

	// Create the search entry
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Enter search term")

	var searchResults []data_access.Track

	// Create the list to display search results
	searchList := widget.NewList(
		func() int {
			return len(searchResults)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				newLabelWithFixedWidth("", 600), // Title
				newLabelWithFixedWidth("", 500), // Album
				newLabelWithFixedWidth("", 400), // Artist
				newLabelWithFixedWidth("", 100), // Year
				newLabelWithFixedWidth("", 200), // Genre
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			track := searchResults[id]
			hbox := item.(*fyne.Container)
			hbox.Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(track.Title)
			hbox.Objects[1].(*fyne.Container).Objects[0].(*widget.Label).SetText(track.Album)
			hbox.Objects[2].(*fyne.Container).Objects[0].(*widget.Label).SetText(track.Artist)
			hbox.Objects[3].(*fyne.Container).Objects[0].(*widget.Label).SetText(strconv.Itoa(track.Year))
			hbox.Objects[4].(*fyne.Container).Objects[0].(*widget.Label).SetText(track.Genre)
		},
	)

	// Set the OnSelected function to handle click events
	searchList.OnSelected = func(id widget.ListItemID) {
		track := searchResults[id]
		albFromTrack := data_access.Album{
			ID:          track.ID,
			Title:       track.Album,
			AlbumArtist: track.AlbumArtist,
			PictureData: track.PictureData,
		}
		UpdateContent(mainContent, CreateLibraryView(albFromTrack))
		LastLibraryViewAlbum = albFromTrack
		homeButton.Importance = widget.MediumImportance
		exploreButton.Importance = widget.MediumImportance
		libraryButton.Importance = widget.HighImportance
		suggestButton.Importance = widget.MediumImportance
		searchButton.Importance = widget.MediumImportance
		homeButton.Refresh()
		exploreButton.Refresh()
		libraryButton.Refresh()
		suggestButton.Refresh()
		searchButton.Refresh()
	}

	// Function to update the search results
	updateSearchResults := func(searchTerm string) {
		// Update searchResults with the query from the database based on searchTerm
		// Then refresh the list
		var err error
		searchResults, err = data_access.GetFulltextSearchResults(searchTerm)
		if err != nil {
			log.Fatal(err)
		}
		searchList.Refresh()
	}

	// Set the function to be executed when the user presses Enter
	searchEntry.OnSubmitted = func(s string) {
		updateSearchResults(s)
	}

	cont := container.NewBorder(searchEntry, nil, nil, nil, searchList)

	return cont
}
