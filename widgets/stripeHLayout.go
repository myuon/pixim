package widgets

import "fyne.io/fyne/v2"

type StripeHLayout struct {
}

var _ fyne.Layout = (*StripeHLayout)(nil)

func (s *StripeHLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects)
	if n == 0 {
		return
	}

	h := size.Height / float32(n)

	for i, o := range objects {
		o.Move(fyne.NewPos(0, float32(i)*h))
		o.Resize(fyne.NewSize(size.Width, 0))
	}
}

func (s *StripeHLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}
