package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ProcedureWidget struct {
	widget.BaseWidget
	Result fyne.CanvasObject
	OnInit func() fyne.CanvasObject
}

var _ fyne.Widget = (*ProcedureWidget)(nil)

func NewProcedureWidget(
	onInit func() fyne.CanvasObject,
) *ProcedureWidget {
	item := &ProcedureWidget{
		OnInit: onInit,
		Result: onInit(),
	}
	item.ExtendBaseWidget(item)

	return item
}

func (c *ProcedureWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.Result)
}

func (c *ProcedureWidget) Refresh() {
	c.Result = c.OnInit()
	c.BaseWidget.Refresh()
}
