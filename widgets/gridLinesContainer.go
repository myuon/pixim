package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func NewGridLinesContainer(size fyne.Size, gap int, color color.Color) *fyne.Container {
	linesV := []fyne.CanvasObject{}
	for i := 0; i < int(size.Width); i++ {
		if i%int(gap) == 0 {
			line := canvas.NewLine(color)
			line.StrokeWidth = 1
			line.Resize(fyne.NewSize(0, size.Height))

			linesV = append(linesV, line)
		}
	}

	linesH := []fyne.CanvasObject{}
	for i := 0; i < int(size.Height); i++ {
		if i%int(gap) == 0 {
			line := canvas.NewLine(color)
			line.StrokeWidth = 1
			line.Resize(fyne.NewSize(size.Width, 0))

			linesH = append(linesH, line)
		}
	}

	gridLines := container.New(
		&StackingLayout{},
		container.New(&StripeVLayout{}, linesV...),
		container.New(&StripeHLayout{}, linesH...),
	)

	return gridLines
}
