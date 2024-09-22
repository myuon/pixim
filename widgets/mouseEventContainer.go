package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type MouseEventContainer struct {
	fyne.CanvasObject

	OnMouseDown func(*desktop.MouseEvent)
	OnMouseUp   func(*desktop.MouseEvent)
	OnMouseMove func(*desktop.MouseEvent)
	OnMouseOut  func()
}

var _ fyne.Widget = (*MouseEventContainer)(nil)
var _ desktop.Cursorable = (*MouseEventContainer)(nil)
var _ desktop.Mouseable = (*MouseEventContainer)(nil)
var _ desktop.Hoverable = (*MouseEventContainer)(nil)

func NewMouseEventContainer(chilren fyne.CanvasObject) *MouseEventContainer {
	item := &MouseEventContainer{
		CanvasObject: chilren,
		OnMouseDown:  func(e *desktop.MouseEvent) {},
		OnMouseUp:    func(e *desktop.MouseEvent) {},
		OnMouseMove:  func(e *desktop.MouseEvent) {},
		OnMouseOut:   func() {},
	}

	return item
}

func (m *MouseEventContainer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.CanvasObject)
}

func (m *MouseEventContainer) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (m *MouseEventContainer) MouseDown(e *desktop.MouseEvent) {
	m.OnMouseDown(e)
}

func (m *MouseEventContainer) MouseUp(e *desktop.MouseEvent) {
	m.OnMouseUp(e)
}

func (m *MouseEventContainer) MouseMoved(e *desktop.MouseEvent) {
	m.OnMouseMove(e)
}

func (m *MouseEventContainer) MouseIn(e *desktop.MouseEvent) {
}

func (m *MouseEventContainer) MouseOut() {
	m.OnMouseOut()
}
