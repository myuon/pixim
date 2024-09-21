package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	size := 64
	ratio := 1.0

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

	a := app.New()
	w := a.NewWindow("Hello")

	cimg := canvas.NewImageFromImage(img)
	cimg.FillMode = canvas.ImageFillOriginal
	cimg.ScaleMode = canvas.ImageScalePixels

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
							ratio *= 2
							cimg.Resize(fyne.NewSize(float32(float64(size)*ratio), float32(float64(size)*ratio)))
						},
					},
					{
						Label: "Zoom out",
						Action: func() {
							ratio /= 2
							cimg.Resize(fyne.NewSize(float32(float64(size)*ratio), float32(float64(size)*ratio)))
						},
					},
				},
			},
		},
	})

	w.Resize(fyne.NewSize(400, 400))

	content := container.New(layout.NewCenterLayout(), cimg)
	w.SetContent(content)

	w.ShowAndRun()
}
