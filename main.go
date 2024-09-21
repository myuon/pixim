package main

import (
	"image"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ImageView struct {
	Image   *image.RGBA
	Ratio   float64
	Refresh func()
}

func NewImageView() *ImageView {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	size := 64

	// 市松模様を描画
	blockSize := size / 8
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// x, y の座標に応じて色を決める（黒と白の市松模様）
			if (x/blockSize+y/blockSize)%2 == 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	return &ImageView{
		Image:   img,
		Ratio:   1.0,
		Refresh: func() {},
	}
}

func main() {
	imageView := NewImageView()
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
			{
				Label: "View",
				Items: []*fyne.MenuItem{
					{
						Label: "Zoom in",
						Action: func() {
							imageView.Ratio *= 2
							imageView.Refresh()
						},
					},
					{
						Label: "Zoom out",
						Action: func() {
							imageView.Ratio /= 2
							imageView.Refresh()
						},
					},
				},
			},
		},
	})

	w.Resize(fyne.NewSize(400, 400))

	cimg := canvas.NewImageFromImage(imageView.Image)
	cimg.FillMode = canvas.ImageFillOriginal
	cimg.ScaleMode = canvas.ImageScalePixels
	imageView.Refresh = func() {
		cimg.Resize(fyne.NewSize(float32(float64(imageView.Image.Bounds().Dx())*imageView.Ratio), float32(float64(imageView.Image.Bounds().Dy())*imageView.Ratio)))
	}
	imageView.Refresh()

	mode := "Move"
	dragStart := fyne.NewPos(0, 0)

	mainCanvas := NewMainCanvas(cimg)
	mainCanvas.OnMouseDown = func(e *desktop.MouseEvent) {
		if mode == "Move" {
			dragStart = e.Position
			log.Printf("Move: %v", dragStart)
		}
		if mode == "Magnifier" {
			if e.Button == desktop.MouseButtonPrimary {
				imageView.Ratio *= 2
				imageView.Refresh()
			} else if e.Button == desktop.MouseButtonSecondary {
				imageView.Ratio /= 2
				imageView.Refresh()
			}
		}
	}
	mainCanvas.OnMouseUp = func(e *desktop.MouseEvent) {
		if mode == "Move" {
			log.Printf("Move: %v %v", e.Position, dragStart)
			cimg.Move(fyne.NewPos(float32(e.Position.X-dragStart.X), float32(e.Position.Y-dragStart.Y)))
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
		),
		container.New(layout.NewCenterLayout(), mainCanvas),
	)
	w.SetContent(content)

	w.ShowAndRun()
}

type MainCanvas struct {
	widget.BaseWidget
	Image       *canvas.Image
	OnMouseDown func(*desktop.MouseEvent)
	OnMouseUp   func(*desktop.MouseEvent)
}

var _ fyne.Widget = (*MainCanvas)(nil)
var _ desktop.Mouseable = (*MainCanvas)(nil)

func NewMainCanvas(image *canvas.Image) *MainCanvas {
	item := &MainCanvas{
		Image: image,
	}
	item.ExtendBaseWidget(item)

	return item
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
