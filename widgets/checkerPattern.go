package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CheckerPattern struct {
	widget.BaseWidget
	ContainerSize fyne.Size
	CheckerSize   fyne.Size
}

func NewCheckerPattern(
	containerSize fyne.Size,
	checkerSize fyne.Size,
) *CheckerPattern {
	item := &CheckerPattern{
		ContainerSize: containerSize,
		CheckerSize:   checkerSize,
	}
	item.ExtendBaseWidget(item)

	return item
}

func (c *CheckerPattern) MinSize() fyne.Size {
	return c.ContainerSize
}

func (c *CheckerPattern) CreateRenderer() fyne.WidgetRenderer {
	objects := []fyne.CanvasObject{}

	for i := 0; i < int(c.ContainerSize.Width/c.CheckerSize.Width); i++ {
		for j := 0; j < int(c.ContainerSize.Height/c.CheckerSize.Height); j++ {
			if (i+j)%2 == 0 {
				rect := canvas.NewRectangle(color.RGBA{0, 0, 0, 0x20})
				rect.Resize(c.CheckerSize)
				rect.Move(fyne.NewPos(float32(i)*c.CheckerSize.Width, float32(j)*c.CheckerSize.Height))

				objects = append(objects, rect)
			} else {
				rect := canvas.NewRectangle(color.White)
				rect.Resize(c.CheckerSize)
				rect.Move(fyne.NewPos(float32(i)*c.CheckerSize.Width, float32(j)*c.CheckerSize.Height))

				objects = append(objects, rect)
			}
		}
	}

	return widget.NewSimpleRenderer(
		container.New(&StackingLayout{}, objects...),
	)
}
