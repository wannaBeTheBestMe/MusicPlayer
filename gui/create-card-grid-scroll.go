package gui

import (
	"MusicPlayer/data_access"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func CreateCardGridScroll() *container.Scroll {
	loaded := data_access.HomeAlbumsLoaded
	if !loaded {
		for !loaded {
		}
	}

	cardGrid := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 400)))

	for _, alb := range data_access.HomeAlbumsArr {
		albumCard := createImgCard(alb.Title, "Artist", alb.PictureData)
		cardGrid.Add(albumCard)
	}
	//cardGridScroll.SetMinSize()

	return container.NewScroll(cardGrid)
}

func createImgCard(title string, artist string, imgData []byte) *fyne.Container {
	reader := bytes.NewReader(imgData)
	img := canvas.NewImageFromReader(reader, "image.jpg")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(250, 250))
	//img.Resize(fyne.NewSize(10, 10))
	//img.SetMinSize(fyne.NewSize(100, 100))

	//card := widget.NewCard(title, artist, img)
	//return card

	titleLab := widget.NewLabel(title)
	titleLab.TextStyle = fyne.TextStyle{Bold: true}
	titleLab.Wrapping = fyne.TextWrapWord
	titleLab.Alignment = fyne.TextAlignCenter

	card := container.NewVBox(titleLab, widget.NewSeparator(), img)
	return card
}
