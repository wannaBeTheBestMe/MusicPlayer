package gui

import (
	"MusicPlayer/data_access"
	"MusicPlayer/playback"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"math"
	"time"
)

var trackDuration time.Duration

var vol float64

var playPauseButton *widget.Button

func CreatePlayerControls(myWindow *fyne.Window) *fyne.Container {
	skipPrevious := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), nil)
	skipPrevOnTapped := func() {
		if LastLibraryViewAlbum.ID == 0 || currentTrack.ID == 0 {
			return
		}

		tracks, err := data_access.GetTracksInAlbum(LastLibraryViewAlbum)
		if err != nil {
			log.Fatal(err)
		}

		for trackNum, track := range tracks {
			if currentTrack.ID == track.ID {
				currentTrack = tracks[trackNum-1]
				trackList.UnselectAll()
				trackList.Select(trackNum - 1)
				break
			}
		}
	}
	skipPrevious.OnTapped = skipPrevOnTapped
	(*myWindow).Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyP {
			skipPrevOnTapped()
		}
	})

	playPauseButton = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil)
	playPauseOnTapped := func() {
		if playPauseButton.Icon == theme.MediaPlayIcon() {
			playPauseButton.SetIcon(theme.MediaPauseIcon())
			playback.ResumeAudio(playback.PDevice)
		} else {
			playPauseButton.SetIcon(theme.MediaPlayIcon())
			playback.PauseAudio(playback.PDevice)
		}
		playPauseButton.Refresh()
	}
	playPauseButton.OnTapped = playPauseOnTapped
	(*myWindow).Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeySpace {
			playPauseOnTapped()
		}
	})

	skipNext := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), nil)
	skipNextOnTapped := func() {
		if LastLibraryViewAlbum.ID == 0 || currentTrack.ID == 0 {
			return
		}

		tracks, err := data_access.GetTracksInAlbum(LastLibraryViewAlbum)
		if err != nil {
			log.Fatal(err)
		}

		for trackNum, track := range tracks {
			if currentTrack.ID == track.ID {
				currentTrack = tracks[trackNum+1]
				trackList.UnselectAll()
				trackList.Select(trackNum + 1)
				break
			}
		}
	}
	skipNext.OnTapped = skipNextOnTapped
	(*myWindow).Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyN {
			skipNextOnTapped()
		}
	})

	loopButton := widget.NewButton("Loop", nil)
	loopButton.OnTapped = func() {
		if loopButton.Importance == widget.MediumImportance {
			if playPauseButton.Icon == theme.MediaPauseIcon() {
				playback.PauseAudio(playback.PDevice)
				playback.SetLoopPlayback(true)
				playback.ResumeAudio(playback.PDevice)
			} else {
				playback.SetLoopPlayback(true)
			}

			loopButton.Importance = widget.HighImportance
			loopButton.Refresh()
		} else {
			if playPauseButton.Icon == theme.MediaPauseIcon() {
				playback.PauseAudio(playback.PDevice)
				playback.SetLoopPlayback(false)
				playback.ResumeAudio(playback.PDevice)
			} else {
				playback.SetLoopPlayback(false)
			}

			loopButton.Importance = widget.MediumImportance
			loopButton.Refresh()
		}
	}

	shuffleButton := widget.NewButton("Shuffle", nil)
	shuffleButton.OnTapped = func() {
		if shuffleButton.Importance == widget.MediumImportance {
			shuffleButton.Importance = widget.HighImportance
			shuffleButton.Refresh()
		} else {
			shuffleButton.Importance = widget.MediumImportance
			shuffleButton.Refresh()
		}
	}

	mediaButtons := container.NewHBox(loopButton, skipPrevious, playPauseButton, skipNext, shuffleButton)

	vol = 1
	volumeSlider := widget.NewSliderWithData(0, 1, binding.BindFloat(&vol))
	volumeSlider.Step = 0.01
	volumeSlider.OnChanged = func(v float64) {
		if v == 0 {
			playback.SetSilent(true)
		} else {
			playback.SetSilent(false)
			playback.SetVolume(float32(v))
		}
	}
	volumeSliderCont := container.NewHBox(volumeSlider)

	mySpacer := canvas.NewRectangle(color.Transparent)
	mySpacer.SetMinSize(volumeSlider.MinSize())

	mediaButtonsCen := container.NewHBox(volumeSliderCont, layout.NewSpacer(), mediaButtons, layout.NewSpacer(), mySpacer)

	trackDuration = time.Duration(playback.GetTotalLength(playback.PDecoder)) * time.Second
	playbackDurLabel := widget.NewLabel(FormatTime(trackDuration))

	elapsed := time.Duration(playback.GetCurrentPosition(playback.PDecoder)) * time.Second
	playbackPosLabel := widget.NewLabel(FormatTime(elapsed))

	percent := elapsed.Seconds() / trackDuration.Seconds()

	playbackSlider := widget.NewSlider(0, 1)
	playbackSlider.Step = 0.001
	playbackSlider.OnChanged = func(percentT float64) {
		t := percentT * trackDuration.Seconds()
		if math.Abs(t-elapsed.Seconds()) > 0.2 {
			if playPauseButton.Icon == theme.MediaPauseIcon() {
				playback.PauseAudio(playback.PDevice)
				playback.SeekToTime(playback.PDecoder, uint64(t/2)) // Using t/2 is necessary here, incredibly strange...
				playback.ResumeAudio(playback.PDevice)
			} else {
				playback.SeekToTime(playback.PDecoder, uint64(t/2)) // Using t/2 is necessary here, incredibly strange...
			}
		}
	}

	playbackCont := container.New(
		NewFixedSizeLayout(playbackPosLabel, playbackDurLabel),
		playbackPosLabel,
		playbackSlider,
		playbackDurLabel,
	)

	go func() {
		for range time.NewTicker(200 * time.Millisecond).C {
			elapsed = time.Duration(playback.GetCurrentPosition(playback.PDecoder)) * time.Second
			playbackPosLabel.SetText(FormatTime(elapsed))

			trackDuration = time.Duration(playback.GetTotalLength(playback.PDecoder)) * time.Second
			playbackDurLabel.SetText(FormatTime(trackDuration))

			percent = elapsed.Seconds() / trackDuration.Seconds()
			playbackSlider.SetValue(percent)
			playbackSlider.Refresh()
		}
	}()

	mediaControls := container.NewVBox(
		mediaButtonsCen,
		playbackCont,
	)
	mediaControls.Resize(fyne.NewSize(600, 100))
	return mediaControls
}

func FormatTime(duration time.Duration) string {
	totalSeconds := int(duration.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
