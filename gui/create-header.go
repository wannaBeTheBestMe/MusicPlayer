package gui

import (
	"MusicPlayer/data_access"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
)

func CreateHeader() *fyne.Container {
	title := canvas.NewText("MusicPlayer", colornames.BlueA100)
	title.TextSize = 48
	titleCont := container.NewVBox(title)

	musicSize := canvas.NewText(data_access.GetMusicDirSize(), colornames.White)
	cont := container.NewHBox(titleCont, layout.NewSpacer(), musicSize)
	header := container.NewStack()
	header.Add(canvas.NewRectangle(color.NRGBA{
		R: 59,
		G: 59,
		B: 62,
		A: 20,
	}))
	header.Add(cont)
	return header
}
