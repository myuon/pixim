package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type CheckerPattern struct {
	ProcedureWidget
	ContainerSize fyne.Size
	CheckerSize   fyne.Size
}

var _ fyne.Widget = (*CheckerPattern)(nil)

func NewCheckerPattern(
	containerSize fyne.Size,
	checkerSize fyne.Size,
) *CheckerPattern {
	item := &CheckerPattern{
		ProcedureWidget: *NewProcedureWidget(func() fyne.CanvasObject {
			objects := []fyne.CanvasObject{}

			for i := 0; i < int(containerSize.Width/checkerSize.Width); i++ {
				for j := 0; j < int(containerSize.Height/checkerSize.Height); j++ {
					if (i+j)%2 == 0 {
						rect := canvas.NewRectangle(color.RGBA{0, 0, 0, 0x20})
						rect.Resize(checkerSize)
						rect.Move(fyne.NewPos(float32(i)*checkerSize.Width, float32(j)*checkerSize.Height))

						objects = append(objects, rect)
					} else {
						rect := canvas.NewRectangle(color.White)
						rect.Resize(checkerSize)
						rect.Move(fyne.NewPos(float32(i)*checkerSize.Width, float32(j)*checkerSize.Height))

						objects = append(objects, rect)
					}
				}
			}

			return container.New(&StackingLayout{}, objects...)
		}),
		ContainerSize: containerSize,
		CheckerSize:   checkerSize,
	}

	return item
}

func (c *CheckerPattern) MinSize() fyne.Size {
	return c.ContainerSize
}

func (c *CheckerPattern) Size() fyne.Size {
	return c.ContainerSize
}
