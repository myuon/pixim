package widgets

import "fyne.io/fyne/v2"

type StackingLayout struct {
}

var _ fyne.Layout = (*StackingLayout)(nil)

func (s *StackingLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (s *StackingLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)

	for _, o := range objects {
		size := o.Size()

		w = max(w, size.Width)
		h = max(h, size.Height)
	}

	return fyne.NewSize(w, h)
}
