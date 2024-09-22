package widgets

import (
	"image"

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
	cimg.Resize(fyne.NewSize(float32(img.Image.Bounds().Dx()), float32(img.Image.Bounds().Dy())))

	item := &ImageCanvas{
		PixImage: img,
		Image:    cimg,
	}

	return item
}

func (m *ImageCanvas) ReplaceImage(img *image.RGBA) {
	m.PixImage = &pixim.PixImage{Image: img}
	m.Image.Image = img
}

func (m *ImageCanvas) SetViewerRatio(ratio float64) {
	m.Image.Resize(fyne.NewSize(float32(m.PixImage.Image.Bounds().Dx())*float32(ratio), float32(m.PixImage.Image.Bounds().Dy())*float32(ratio)))
}

var _ fyne.Widget = (*ImageCanvas)(nil)

func (m *ImageCanvas) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.Image)
}
