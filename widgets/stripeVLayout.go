package widgets

import (
	"fyne.io/fyne/v2"
)

type StripeVLayout struct {
}

var _ fyne.Layout = (*StripeVLayout)(nil)

func (s *StripeVLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects)
	if n == 0 {
		return
	}

	w := size.Width / float32(n)

	for i, o := range objects {
		o.Move(fyne.NewPos(float32(i)*w, 0))
		o.Resize(fyne.NewSize(0, size.Height))
	}
}

func (s *StripeVLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}
