package gui

import "fyne.io/fyne/v2"

type FixedSizeLayout struct {
	leftLabel  fyne.CanvasObject
	rightLabel fyne.CanvasObject
}

func (f *FixedSizeLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
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

func (f *FixedSizeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
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

func NewFixedSizeLayout(leftLabel, rightLabel fyne.CanvasObject) fyne.Layout {
	return &FixedSizeLayout{leftLabel: leftLabel, rightLabel: rightLabel}
}
