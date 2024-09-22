package main

import (
	"image"
	"image/color"
	"image/png"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
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
	a := app.New()
	w := a.NewWindow("Pixim")

	pixImage := pixim.NewPixImage()
	cimg := canvas.NewImageFromImage(pixImage.Image)
	cimg.ScaleMode = canvas.ImageScalePixels

	ratio := 1.0
	containerSize := fyne.NewSize(800, 800)

	background := widgets.NewCheckerPattern(
		fyne.NewSize(containerSize.Width, containerSize.Height),
		fyne.NewSize(40, 40),
	)

	var gridCache image.Image
	cachedRatio := 0.0
	grid := canvas.NewRaster(func(w, h int) image.Image {
		if gridCache != nil && cachedRatio == ratio {
			return gridCache
		}

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

		gridCache = img
		cachedRatio = ratio

		return img
	})
	grid.Resize(containerSize)
	grid.ScaleMode = canvas.ImageScalePixels

	imageCanvas := widgets.NewImageCanvas(pixim.NewPixImage())

	imgContainer := container.New(&widgets.StackingLayout{}, imageCanvas, grid)
	imgContainer.Resize(containerSize)

	scrollContainer := container.NewScroll(imgContainer)
	scrollContainer.Resize(containerSize)

	children := container.New(
		&widgets.StackingLayout{},
		background,
		scrollContainer,
	)

	mode := Move
	dragging := false
	// dragStart := fyne.NewPos(0, 0)
	// originalPos := fyne.NewPos(0, 0)

	prevPos := fyne.NewPos(0, 0)

	var currentColor color.Color = color.Transparent

	mainCanvas := NewMainCanvas(children)
	mainCanvas.OnMouseDown = func(e *desktop.MouseEvent) {
		pos := e.Position
		x := int(float64(pos.X) / ratio)
		y := int(float64(pos.Y) / ratio)
		contains := !(x < 0 || y < 0 || x >= imageCanvas.Image.Image.Bounds().Dx() || y >= imageCanvas.Image.Image.Bounds().Dy())

		if mode == Move {
			dragging = true
			// dragStart = e.Position
			// originalPos = cimg.Position()
		}
		if mode == Magnifier {
			if e.Button == desktop.MouseButtonPrimary {
				ratio *= 2
				mainCanvas.Refresh()
			} else if e.Button == desktop.MouseButtonSecondary {
				ratio /= 2
				mainCanvas.Refresh()
			}
		}
		if mode == Fill {
			if !contains {
				return
			}

			pixImage.Fill(x, y, currentColor)
			mainCanvas.Refresh()
		}
		if mode == Pencil {
			if !contains {
				return
			}

			pixImage.Image.Set(x, y, currentColor)
			mainCanvas.Refresh()

			prevPos = fyne.NewPos(float32(x), float32(y))

			dragging = true
		}
	}
	mainCanvas.OnMouseMove = func(e *desktop.MouseEvent) {
		pos := e.Position
		x := int(float64(pos.X) / ratio)
		y := int(float64(pos.Y) / ratio)
		contains := !(x < 0 || y < 0 || x >= imageCanvas.Image.Image.Bounds().Dx() || y >= imageCanvas.Image.Image.Bounds().Dy())

		if mode == Move && dragging {
			// mainCanvas.MoveImage(fyne.NewPos(float32(e.Position.X-dragStart.X)+originalPos.X, float32(e.Position.Y-dragStart.Y)+originalPos.Y))
		}
		if mode == Pencil && dragging {
			if !contains {
				return
			}

			pixImage.DrawLine(int(prevPos.X), int(prevPos.Y), x, y, currentColor)
			mainCanvas.Refresh()

			prevPos = fyne.NewPos(float32(x), float32(y))
		}
	}
	mainCanvas.OnMouseUp = func(e *desktop.MouseEvent) {
		if mode == Move {
			dragging = false
			// mainCanvas.MoveImage(fyne.NewPos(float32(e.Position.X-dragStart.X)+originalPos.X, float32(e.Position.Y-dragStart.Y)+originalPos.Y))
		}
		if mode == Pencil {
			dragging = false
		}
	}

	w.SetMainMenu(&fyne.MainMenu{
		Items: []*fyne.Menu{
			{
				Label: "File",
				Items: []*fyne.MenuItem{
					{
						Label: "New",
						Action: func() {
							width := widget.NewEntry()
							width.Validator = validation.NewRegexp(`\d+`, "Width must be a number")

							height := widget.NewEntry()
							height.Validator = validation.NewRegexp(`\d+`, "Height must be a number")

							dialog.ShowForm("Create new image", "Create", "Cancel", []*widget.FormItem{
								widget.NewFormItem("Width", width),
								widget.NewFormItem("Height", height),
							}, func(b bool) {
								if !b {
									return
								}

								w, _ := strconv.Atoi(width.Text)
								h, _ := strconv.Atoi(height.Text)

								img := image.NewRGBA(image.Rect(0, 0, w, h))
								for i := 0; i < w; i++ {
									for j := 0; j < h; j++ {
										img.Set(i, j, color.White)
									}
								}

								imageCanvas.ReplaceImage(img)
								mainCanvas.Refresh()
							}, w)
						},
					},
					{
						Label: "Open",
						Action: func() {
							dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
								if err != nil {
									dialog.ShowError(err, w)
									return
								}

								img, _, err := image.Decode(f)
								if err != nil {
									dialog.ShowError(err, w)
									return
								}

								pixImage.Image = img.(*image.RGBA)
								mainCanvas.Refresh()
							}, w).Show()
						},
					},
					{
						Label: "Save",
						Action: func() {
							dialog.NewFileSave(func(f fyne.URIWriteCloser, err error) {
								if err != nil {
									dialog.ShowError(err, w)
									return
								}

								if err := png.Encode(f, pixImage.Image); err != nil {
									dialog.ShowError(err, w)
									return
								}
							}, w).Show()
						},
					},
				},
			},
		},
	})

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
						currentColor = c
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
	*fyne.Container

	OnMouseDown func(*desktop.MouseEvent)
	OnMouseUp   func(*desktop.MouseEvent)
	OnMouseMove func(*desktop.MouseEvent)
}

var _ fyne.Widget = (*MainCanvas)(nil)
var _ desktop.Cursorable = (*MainCanvas)(nil)
var _ desktop.Mouseable = (*MainCanvas)(nil)
var _ desktop.Hoverable = (*MainCanvas)(nil)

func NewMainCanvas(chilren *fyne.Container) *MainCanvas {
	item := &MainCanvas{
		Container: chilren,
	}

	return item
}

func (m *MainCanvas) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.Container)
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
	m.Container.Refresh()
}

// func (m *MainCanvas) Refresh() {
// 	m.Image.SetMinSize(fyne.NewSize(float32(float64(m.Image.Image.Bounds().Dx())**m.Ratio), float32(float64(m.Image.Image.Bounds().Dy())**m.Ratio)))
// 	m.Image.Resize(fyne.NewSize(float32(float64(m.Image.Image.Bounds().Dx())**m.Ratio), float32(float64(m.Image.Image.Bounds().Dy())**m.Ratio)))
// 	m.Image.Refresh()
// 	m.Widget.Refresh()
// }

// func (m *MainCanvas) MoveImage(pos fyne.Position) {
// 	m.ImagePosition = pos
// }
