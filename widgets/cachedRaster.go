package widgets

import (
	"image"

	"fyne.io/fyne/v2/canvas"
)

func NewCachedRaster(
	calculateKey func() any,
	render func(w, h int) image.Image,
) *canvas.Raster {
	var cache image.Image
	var prevKey any

	return canvas.NewRaster(func(w, h int) image.Image {
		key := calculateKey()
		if cache != nil && prevKey == key {
			return cache
		}

		cache = render(w, h)
		prevKey = key

		return cache
	})
}
