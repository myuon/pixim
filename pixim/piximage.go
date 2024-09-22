package pixim

import (
	"image"
	"image/color"
)

type PixImage struct {
	Image *image.RGBA
}

func NewPixImage() *PixImage {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	size := 64

	blockSize := size / 8
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if (x/blockSize+y/blockSize)%2 == 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	return &PixImage{
		Image: img,
	}
}

func (i *PixImage) Fill(x, y int, color color.Color) {
	original := i.Image.At(x, y)

	i.Image.Set(x, y, color)

	stack := []struct{ x, y int }{{x, y}}
	for len(stack) > 0 {
		pos := stack[0]
		stack = stack[1:]

		for _, d := range []struct{ x, y int }{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
			next := struct{ x, y int }{pos.x + d.x, pos.y + d.y}
			if next.x < 0 || next.y < 0 || next.x >= i.Image.Bounds().Dx() || next.y >= i.Image.Bounds().Dy() {
				continue
			}

			if i.Image.At(next.x, next.y) == original {
				i.Image.Set(next.x, next.y, color)
				stack = append(stack, next)
			}
		}
	}
}

func (i *PixImage) DrawLine(x1, y1 int, x2, y2 int, color color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	steps := max(dx, dy)

	for s := 0; s < steps; s++ {
		x := x1 + s*dx/steps
		y := y1 + s*dy/steps

		i.Image.Set(x, y, color)
	}
}

func (i *PixImage) Set(x, y int, color color.Color) {
	i.Image.Set(x, y, color)
}
