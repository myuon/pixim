package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/myuon/pixim/pixim"
)

type ImageCanvas struct {
	PixImage *pixim.PixImage
	*canvas.Image
}

func NewImageCanvas(img *pixim.PixImage) *ImageCanvas {
	cimg := canvas.NewImageFromImage(img.Image)
	cimg.ScaleMode = canvas.ImageScalePixels
	cimg.Resize(fyne.NewSize(100, 100))

	item := &ImageCanvas{
		PixImage: img,
		Image:    cimg,
	}

	return item
}

func (m *ImageCanvas) ReplaceImage(pixImage *pixim.PixImage) {
	m.PixImage = pixImage
	m.Image.Image = pixImage.Image
	m.Refresh()
}

var _ fyne.Widget = (*ImageCanvas)(nil)

func (m *ImageCanvas) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.Image)
}
