package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func main() {
	pixImage := NewPixImage()
	a := app.New()
	w := a.NewWindow("QuickPix")

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
	cimg.FillMode = canvas.ImageFillOriginal
	cimg.ScaleMode = canvas.ImageScalePixels

	mode := "Move"
	dragging := false
	dragStart := fyne.NewPos(0, 0)
	originalPos := cimg.Position()

	mainCanvas := NewMainCanvas(cimg)
	mainCanvas.OnMouseDown = func(e *desktop.MouseEvent) {
		if mode == "Move" {
			dragging = true
			dragStart = e.Position
			originalPos = cimg.Position()
		}
		if mode == "Magnifier" {
			if e.Button == desktop.MouseButtonPrimary {
				mainCanvas.Ratio *= 2
				mainCanvas.Refresh()
			} else if e.Button == desktop.MouseButtonSecondary {
				mainCanvas.Ratio /= 2
				mainCanvas.Refresh()
			}
		}
		if mode == "Fill" {
			pos := e.Position
			x := int(float64(pos.X) / mainCanvas.Ratio)
			y := int(float64(pos.Y) / mainCanvas.Ratio)

			if x < 0 || y < 0 || x >= pixImage.Image.Bounds().Dx() || y >= pixImage.Image.Bounds().Dy() {
				return
			}

			pixImage.Fill(x, y, mainCanvas.Color)
			mainCanvas.Refresh()
		}
		if mode == "Pencil" {
			pos := e.Position
			x := int(float64(pos.X) / mainCanvas.Ratio)
			y := int(float64(pos.Y) / mainCanvas.Ratio)

			if x < 0 || y < 0 || x >= pixImage.Image.Bounds().Dx() || y >= pixImage.Image.Bounds().Dy() {
				return
			}

			pixImage.Image.Set(x, y, mainCanvas.Color)
			mainCanvas.Refresh()

			dragging = true
		}
	}
	mainCanvas.OnMouseMove = func(e *desktop.MouseEvent) {
		if mode == "Move" && dragging {
			cimg.Move(fyne.NewPos(float32(e.Position.X-dragStart.X)+originalPos.X, float32(e.Position.Y-dragStart.Y)+originalPos.Y))
		}
		if mode == "Pencil" && dragging {
			pos := e.Position
			x := int(float64(pos.X) / mainCanvas.Ratio)
			y := int(float64(pos.Y) / mainCanvas.Ratio)

			if x < 0 || y < 0 || x >= pixImage.Image.Bounds().Dx() || y >= pixImage.Image.Bounds().Dy() {
				return
			}

			pixImage.Image.Set(x, y, mainCanvas.Color)
			mainCanvas.Refresh()
		}
	}
	mainCanvas.OnMouseUp = func(e *desktop.MouseEvent) {
		if mode == "Move" {
			dragging = false
			cimg.Move(fyne.NewPos(float32(e.Position.X-dragStart.X)+originalPos.X, float32(e.Position.Y-dragStart.Y)+originalPos.Y))
		}
		if mode == "Pencil" {
			dragging = false
		}
	}

	content := container.NewHBox(
		container.NewVBox(
			widget.NewButton("Move", func() {
				mode = "Move"
			}),
			widget.NewButton("Magnifier", func() {
				mode = "Magnifier"
			}),
			widget.NewButton("Fill", func() {
				mode = "Fill"
			}),
			widget.NewButton("Pencil", func() {
				mode = "Pencil"
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
	w.SetContent(content)
	w.ShowAndRun()
}

type MainCanvas struct {
	widget.BaseWidget
	Image *canvas.Image
	Ratio float64
	Color color.Color

	OnMouseDown func(*desktop.MouseEvent)
	OnMouseUp   func(*desktop.MouseEvent)
	OnMouseMove func(*desktop.MouseEvent)
}

var _ fyne.Widget = (*MainCanvas)(nil)
var _ desktop.Cursorable = (*MainCanvas)(nil)
var _ desktop.Mouseable = (*MainCanvas)(nil)
var _ desktop.Hoverable = (*MainCanvas)(nil)

func NewMainCanvas(image *canvas.Image) *MainCanvas {
	item := &MainCanvas{
		Image: image,
		Ratio: 1.0,
		Color: color.Black,
	}
	item.ExtendBaseWidget(item)

	return item
}

func (m *MainCanvas) MinSize() fyne.Size {
	return fyne.NewSize(400, 400)
}

func (m *MainCanvas) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.Image)
}

func (m *MainCanvas) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (m *MainCanvas) MouseDown(e *desktop.MouseEvent) {
	m.OnMouseDown(e)
}

func (m *MainCanvas) MouseUp(e *desktop.MouseEvent) {
	m.OnMouseUp(e)
}

func (m *MainCanvas) MouseMoved(e *desktop.MouseEvent) {
	m.OnMouseMove(e)
}

func (m *MainCanvas) MouseIn(e *desktop.MouseEvent) {
}

func (m *MainCanvas) MouseOut() {
}

func (m *MainCanvas) Refresh() {
	m.Image.Resize(fyne.NewSize(float32(float64(m.Image.Image.Bounds().Dx())*m.Ratio), float32(float64(m.Image.Image.Bounds().Dy())*m.Ratio)))
	m.Image.Refresh()
	m.BaseWidget.Refresh()
}
