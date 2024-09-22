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

type Editor struct {
	Image        *pixim.PixImage
	Ratio        float64
	CurrentColor color.Color
	View         fyne.CanvasObject

	OnChangeImage        func(*pixim.PixImage)
	OnChangeRatio        func(float64)
	OnChangeCurrentColor func(color.Color)
}

func (e *Editor) SetRatio(ratio float64) {
	e.Ratio = ratio
	e.OnChangeRatio(ratio)
}

func (e *Editor) SetImage(img *pixim.PixImage) {
	e.Image = img
	e.OnChangeImage(img)
}

func (e *Editor) Fill(x, y int) {
	e.Image.Fill(x, y, e.CurrentColor)
	e.OnChangeImage(e.Image)
}

func (e *Editor) Paint(x, y int) {
	e.Image.Set(x, y, e.CurrentColor)
	e.OnChangeImage(e.Image)
}

func (e *Editor) DrawLine(x1, y1, x2, y2 int) {
	e.Image.DrawLine(x1, y1, x2, y2, e.CurrentColor)
	e.OnChangeImage(e.Image)
}

func (e *Editor) SetCurrentColor(c color.Color) {
	e.CurrentColor = c
	e.OnChangeCurrentColor(c)
}

func main() {
	a := app.New()
	w := a.NewWindow("Pixim")

	editor := Editor{
		Image:        pixim.NewPixImage(),
		Ratio:        1.0,
		CurrentColor: color.Black,
		View:         nil,
	}

	containerSize := fyne.NewSize(800, 800)

	background := widgets.NewCachedRaster(
		func() any {
			return true
		},
		func(w, h int) image.Image {
			img := image.NewRGBA(image.Rect(0, 0, w, h))

			blockSize := 10
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					if (x/blockSize+y/blockSize)%2 == 1 {
						img.Set(x, y, color.RGBA{0, 0, 0, 0x20})
					}
				}
			}

			return img
		},
	)

	gridColor := color.RGBA{0xd0, 0xd0, 0xd0, 0xff}

	linesV := []fyne.CanvasObject{}
	for i := 0; i < editor.Image.Image.Bounds().Dx(); i++ {
		if i%int(editor.Ratio) == 0 {
			line := canvas.NewLine(gridColor)
			line.StrokeWidth = 1
			line.Resize(fyne.NewSize(0, containerSize.Height))

			linesV = append(linesV, line)
		}
	}

	linesH := []fyne.CanvasObject{}
	for i := 0; i < editor.Image.Image.Bounds().Dy(); i++ {
		if i%int(editor.Ratio) == 0 {
			line := canvas.NewLine(gridColor)
			line.StrokeWidth = 1
			line.Resize(fyne.NewSize(containerSize.Width, 0))

			linesH = append(linesH, line)
		}
	}

	gridLines := container.New(
		&widgets.StackingLayout{},
		container.New(&widgets.StripeVLayout{}, linesV...),
		container.New(&widgets.StripeHLayout{}, linesH...),
	)
	gridLines.Hide()

	imageCanvas := widgets.NewImageCanvas(editor.Image)

	imgContainer := container.New(&widgets.StackingLayout{}, imageCanvas, gridLines)

	scrollContainer := container.NewScroll(imgContainer)

	children := container.New(&widgets.StackingLayout{}, background, scrollContainer)
	children.Resize(containerSize)

	mode := Move
	dragging := false
	// dragStart := fyne.NewPos(0, 0)
	// originalPos := fyne.NewPos(0, 0)

	prevPos := fyne.NewPos(0, 0)

	editor.OnChangeRatio = func(ratio float64) {
		size := fyne.NewSize(float32(float64(imageCanvas.Image.Image.Bounds().Dx())*editor.Ratio), float32(float64(imageCanvas.Image.Image.Bounds().Dy())*editor.Ratio))
		imageCanvas.Resize(size)
		imageCanvas.Refresh()

		if editor.Ratio < 5 {
			gridLines.Hide()
		} else {
			gridLines.Resize(size)
			gridLines.Show()
		}
	}
	editor.OnChangeImage = func(img *pixim.PixImage) {
		imageCanvas.ReplaceImage(img.Image)
		imageCanvas.Refresh()
	}

	mouseEventContainer := widgets.NewMouseEventContainer(children)
	mouseEventContainer.OnMouseDown = func(e *desktop.MouseEvent) {
		pos := e.Position
		x := int(float64(pos.X) / editor.Ratio)
		y := int(float64(pos.Y) / editor.Ratio)
		contains := !(x < 0 || y < 0 || x >= imageCanvas.Image.Image.Bounds().Dx() || y >= imageCanvas.Image.Image.Bounds().Dy())

		if mode == Move {
			dragging = true
			// dragStart = e.Position
			// originalPos = cimg.Position()
		}
		if mode == Magnifier {
			if e.Button == desktop.MouseButtonPrimary {
				editor.SetRatio(editor.Ratio * 2)
				mouseEventContainer.Refresh()
			} else if e.Button == desktop.MouseButtonSecondary {
				editor.SetRatio(editor.Ratio / 2)
				mouseEventContainer.Refresh()
			}
		}
		if mode == Fill {
			if !contains {
				return
			}

			editor.Fill(x, y)
		}
		if mode == Pencil {
			if !contains {
				return
			}

			editor.Paint(x, y)
			prevPos = fyne.NewPos(float32(x), float32(y))

			dragging = true
		}
	}
	mouseEventContainer.OnMouseMove = func(e *desktop.MouseEvent) {
		pos := e.Position
		x := int(float64(pos.X) / editor.Ratio)
		y := int(float64(pos.Y) / editor.Ratio)
		contains := !(x < 0 || y < 0 || x >= imageCanvas.Image.Image.Bounds().Dx() || y >= imageCanvas.Image.Image.Bounds().Dy())

		if mode == Move && dragging {
			// mainCanvas.MoveImage(fyne.NewPos(float32(e.Position.X-dragStart.X)+originalPos.X, float32(e.Position.Y-dragStart.Y)+originalPos.Y))
		}
		if mode == Pencil && dragging {
			if !contains {
				return
			}

			editor.DrawLine(int(prevPos.X), int(prevPos.Y), x, y)
			prevPos = fyne.NewPos(float32(x), float32(y))
		}
	}
	mouseEventContainer.OnMouseUp = func(e *desktop.MouseEvent) {
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

								editor.SetImage(&pixim.PixImage{Image: img})
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

								editor.SetImage(&pixim.PixImage{Image: img.(*image.RGBA)})
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

								if err := png.Encode(f, imageCanvas.PixImage.Image); err != nil {
									dialog.ShowError(err, w)
									return
								}

								if err := f.Close(); err != nil {
									dialog.ShowError(err, w)
									return
								}

								dialog.ShowInformation("Saved", "Image saved successfully", w)
							}, w).Show()
						},
					},
				},
			},
		},
	})

	colorRect := canvas.NewRectangle(editor.CurrentColor)
	colorRect.SetMinSize(fyne.NewSize(40, 40))

	editor.OnChangeCurrentColor = func(c color.Color) {
		colorRect.FillColor = c
	}

	content := container.NewHBox(
		mouseEventContainer,
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
						editor.SetCurrentColor(c)
					},
					w,
				).Show()
			}),
			colorRect,
		),
	)

	w.SetContent(content)
	w.ShowAndRun()
}
