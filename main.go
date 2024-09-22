package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/myuon/pixim/pixim"
	"github.com/myuon/pixim/widgets"
)

type EditorMode string

const (
	Move      EditorMode = "Move"
	Magnifier EditorMode = "Magnifier"
	Fill      EditorMode = "Fill"
	Pencil    EditorMode = "Pencil"
)

func main() {
	pixImage := pixim.NewPixImage()
	a := app.New()
	w := a.NewWindow("Pixim")

	w.SetMainMenu(&fyne.MainMenu{
		Items: []*fyne.Menu{
			{
				Label: "File",
				Items: []*fyne.MenuItem{
					{
						Label: "New",
						Action: func() {
							w.SetContent(widget.NewLabel("New content"))
						},
					},
					{
						Label:  "Open",
						Action: func() { w.SetContent(widget.NewLabel("Open content")) },
					},
				},
			},
		},
	})

	cimg := canvas.NewImageFromImage(pixImage.Image)
	cimg.ScaleMode = canvas.ImageScalePixels

	mode := Move
	dragging := false
	dragStart := fyne.NewPos(0, 0)
	originalPos := cimg.Position()

	prevPos := fyne.NewPos(0, 0)

	mainCanvas := NewMainCanvas(cimg, fyne.NewSize(800, 800))
	mainCanvas.OnMouseDown = func(e *desktop.MouseEvent, x, y int, contains bool) {
		if mode == Move {
			dragging = true
			dragStart = e.Position
			originalPos = cimg.Position()
		}
		if mode == Magnifier {
			if e.Button == desktop.MouseButtonPrimary {
				*mainCanvas.Ratio *= 2
				mainCanvas.Refresh()
			} else if e.Button == desktop.MouseButtonSecondary {
				*mainCanvas.Ratio /= 2
				mainCanvas.Refresh()
			}
		}
		if mode == Fill {
			if !contains {
				return
			}

			pixImage.Fill(x, y, mainCanvas.Color)
			mainCanvas.Refresh()
		}
		if mode == Pencil {
			if !contains {
				return
			}

			pixImage.Image.Set(x, y, mainCanvas.Color)
			mainCanvas.Refresh()

			prevPos = fyne.NewPos(float32(x), float32(y))

			dragging = true
		}
	}
	mainCanvas.OnMouseMove = func(e *desktop.MouseEvent, x, y int, contains bool) {
		if mode == Move && dragging {
			mainCanvas.MoveImage(fyne.NewPos(float32(e.Position.X-dragStart.X)+originalPos.X, float32(e.Position.Y-dragStart.Y)+originalPos.Y))
		}
		if mode == Pencil && dragging {
			if !contains {
				return
			}

			pixImage.DrawLine(int(prevPos.X), int(prevPos.Y), x, y, mainCanvas.Color)
			mainCanvas.Refresh()

			prevPos = fyne.NewPos(float32(x), float32(y))
		}
	}
	mainCanvas.OnMouseUp = func(e *desktop.MouseEvent) {
		if mode == Move {
			dragging = false
			mainCanvas.MoveImage(fyne.NewPos(float32(e.Position.X-dragStart.X)+originalPos.X, float32(e.Position.Y-dragStart.Y)+originalPos.Y))
		}
		if mode == Pencil {
			dragging = false
		}
	}

	content := container.NewHBox(
		container.NewVBox(
			widget.NewButton("Move", func() {
				mode = Move
			}),
			widget.NewButton("Magnifier", func() {
				mode = Magnifier
			}),
			widget.NewButton("Fill", func() {
				mode = Fill
			}),
			widget.NewButton("Pencil", func() {
				mode = Pencil
			}),
			widget.NewButton("Color", func() {
				dialog.NewColorPicker(
					"Select a color",
					"foobar",
					func(c color.Color) {
						mainCanvas.Color = c
					},
					w,
				).Show()
			}),
		),
		mainCanvas,
	)
	mainCanvas.Refresh()

	w.SetContent(content)
	w.ShowAndRun()
}

type MainCanvas struct {
	widget.BaseWidget
	Image           *canvas.Image
	ScrollContainer *container.Scroll
	Grid            *canvas.Raster
	Widget          fyne.CanvasObject
	Ratio           *float64
	Color           color.Color
	ContainerSize   fyne.Size
	ImagePosition   fyne.Position

	OnMouseDown func(*desktop.MouseEvent, int, int, bool)
	OnMouseUp   func(*desktop.MouseEvent)
	OnMouseMove func(*desktop.MouseEvent, int, int, bool)
}

var _ fyne.Widget = (*MainCanvas)(nil)
var _ desktop.Cursorable = (*MainCanvas)(nil)
var _ desktop.Mouseable = (*MainCanvas)(nil)
var _ desktop.Hoverable = (*MainCanvas)(nil)

func NewMainCanvas(img *canvas.Image, containerSize fyne.Size) *MainCanvas {
	ratio := 1.0

	background := widgets.NewCheckerPattern(
		fyne.NewSize(containerSize.Width, containerSize.Height),
		fyne.NewSize(40, 40),
	)
	grid := canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, int(containerSize.Width), int(containerSize.Height)))
		if ratio < 5 {
			return img
		}

		for i := 0; i < int(containerSize.Width); i++ {
			for j := 0; j < int(containerSize.Height); j++ {
				if i%int(ratio) == 0 || j%int(ratio) == 0 {
					img.Set(i, j, color.RGBA{0xd0, 0xd0, 0xd0, 0xff})
				}
			}
		}

		return img
	})
	grid.Resize(containerSize)
	grid.ScaleMode = canvas.ImageScalePixels

	imgContainer := container.New(&widgets.StackingLayout{}, img)
	imgContainer.Resize(containerSize)

	scrollContainer := container.NewScroll(imgContainer)
	scrollContainer.Resize(containerSize)

	widget := container.New(
		&widgets.StackingLayout{},
		background,
		// grid,
		scrollContainer,
	)

	item := &MainCanvas{
		Image:           img,
		ScrollContainer: scrollContainer,
		Grid:            grid,
		Widget:          widget,
		Ratio:           &ratio,
		Color:           color.Black,
		ContainerSize:   containerSize,
	}
	item.ExtendBaseWidget(item)

	return item
}

func (m *MainCanvas) MinSize() fyne.Size {
	return m.ContainerSize
}

func (m *MainCanvas) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.Widget)
}

func (m *MainCanvas) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (m *MainCanvas) MouseDown(e *desktop.MouseEvent) {
	pos := e.Position
	x := int(float64(pos.X) / *m.Ratio)
	y := int(float64(pos.Y) / *m.Ratio)
	contains := !(x < 0 || y < 0 || x >= m.Image.Image.Bounds().Dx() || y >= m.Image.Image.Bounds().Dy())

	m.OnMouseDown(e, x, y, contains)
}

func (m *MainCanvas) MouseUp(e *desktop.MouseEvent) {
	m.OnMouseUp(e)
}

func (m *MainCanvas) MouseMoved(e *desktop.MouseEvent) {
	pos := e.Position
	x := int(float64(pos.X) / *m.Ratio)
	y := int(float64(pos.Y) / *m.Ratio)
	contains := !(x < 0 || y < 0 || x >= m.Image.Image.Bounds().Dx() || y >= m.Image.Image.Bounds().Dy())

	m.OnMouseMove(e, x, y, contains)
}

func (m *MainCanvas) MouseIn(e *desktop.MouseEvent) {
}

func (m *MainCanvas) MouseOut() {
}

func (m *MainCanvas) Refresh() {
	m.Image.SetMinSize(fyne.NewSize(float32(float64(m.Image.Image.Bounds().Dx())**m.Ratio), float32(float64(m.Image.Image.Bounds().Dy())**m.Ratio)))
	m.Image.Resize(fyne.NewSize(float32(float64(m.Image.Image.Bounds().Dx())**m.Ratio), float32(float64(m.Image.Image.Bounds().Dy())**m.Ratio)))
	m.Widget.Refresh()
}

func (m *MainCanvas) MoveImage(pos fyne.Position) {
	m.ImagePosition = pos
}
