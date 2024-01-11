package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"image/color"
	"math"
	"time"
)

// var MfilePath = "\"C:\\Users\\Asus\\Downloads\\Lame_Drivers_-_01_-_Frozen_Egg.mp3\""
var MfilePath = "C:\\Users\\Asus\\Music\\MusicPlayer\\Vangelis - Nocturne (2019) [24Bit Hi-Res]\\10 - To a Friend.flac"

var vol float64
var isSilent = false

func CreatePlayerControls(streamer *beep.StreamSeekCloser, format *beep.Format) *fyne.Container {
	skipPrevious := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), nil)

	playPauseButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil)
	playPauseButton.OnTapped = func() {
		if playPauseButton.Icon == theme.MediaPlayIcon() {
			playPauseButton.SetIcon(theme.MediaPauseIcon())
			go playAudio(streamer)
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
			go playAudio(streamer)
		}
	}
	volumeSliderCont := container.NewHBox(volumeSlider)

	mySpacer := canvas.NewRectangle(color.Transparent)
	mySpacer.SetMinSize(volumeSlider.MinSize())

	mediaButtonsCen := container.NewHBox(volumeSliderCont, layout.NewSpacer(), mediaButtons, layout.NewSpacer(), mySpacer)

	pos := 0.0
	playbackPos := binding.BindFloat(&pos)
	duration := func() time.Duration {
		return time.Second * time.Duration((*streamer).Len()) / time.Duration((*format).SampleRate)
	}()
	formattedDuration := formatTime(duration)
	playbackDurLabel := widget.NewLabel(formattedDuration)

	playbackSlider := widget.NewSliderWithData(0, duration.Seconds(), playbackPos)
	playbackSlider.Step = 0.001
	position := func() int {
		return (*streamer).Position()
	}()
	elapsed := time.Second * time.Duration(position) / time.Duration((*format).SampleRate)
	playbackPosText := binding.NewString()
	playbackPosText.Set(formatTime(elapsed))
	playbackPosLabel := widget.NewLabelWithData(playbackPosText)

	playbackSlider.OnChanged = func(t float64) {
		if t-elapsed.Seconds() > 0.1 {
			//speaker.Clear()
			//func() {
			//	if worked := StreamerMu.TryLock(); worked {
			//		speaker.Lock()
			//		defer speaker.Unlock()
			//		defer StreamerMu.Unlock()
			//
			//		fmt.Println("About to do it")
			//		err := (*streamer).Seek(0)
			//		if err != nil {
			//			log.Fatal(err)
			//		}
			//		fmt.Println("Done it")
			//	}
			//}()

			//	fmt.Println("User")
			//	speaker.Clear()
			//	speaker.Close()
			//	Streamer.Close()
			//	Streamer, Format, _ = flac.Decode(Mfile)
			//	err := Streamer.Seek(int(t))
			//	if err != nil {
			//		log.Fatal(err)
			//	}
			//	//go func() {
			//	//	Streamer.Seek(int(t))
			//	//	playAudio()
			//	//}()
			//
			//	//if err := speaker.Init(Format.SampleRate, Format.SampleRate.N(time.Second/10)); err != nil {
			//	//	log.Fatal(err)
			//	//}
			//	//go playAudio()
			//	//pt := time.Second * time.Duration(t) / time.Duration(Format.SampleRate) // int(pt.Seconds())
			//	//fmt.Println(seeker.Position())
		}
	}

	playbackCont := container.New(
		NewFixedSizeLayout(playbackPosLabel, playbackDurLabel),
		playbackPosLabel,
		playbackSlider,
		playbackDurLabel,
	)

	go func() {
		for range time.NewTicker(time.Millisecond).C {
			position = func() int {
				return (*streamer).Position()
			}()
			elapsed = time.Second * time.Duration(position) / time.Duration((*format).SampleRate)
			pos = elapsed.Seconds()
			playbackPos.Reload()
			playbackPosText.Set(formatTime(elapsed))
		}
	}()

	mediaControls := container.NewVBox(
		mediaButtonsCen,
		playbackCont,
	)
	mediaControls.Resize(fyne.NewSize(600, 100))
	return mediaControls
}

// playAudio plays audio (call with go to play in a separate goroutine)
func playAudio(streamer *beep.StreamSeekCloser) {
	donePlaying := make(chan bool)
	streamerWithVol := &effects.Volume{
		Streamer: *streamer,
		Base:     math.E,
		Volume:   vol,
		Silent:   isSilent,
	}

	func() {
		speaker.Play(beep.Seq(streamerWithVol, beep.Callback(func() {
			donePlaying <- true
		})))
	}()

	// Wait for the audio to finish playing
	<-donePlaying
}

func formatTime(duration time.Duration) string {
	totalSeconds := int(duration.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
