package ip

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func Filter_Window(app fyne.App) func() { //保存
	return func() {
		w3 := A.NewWindow("Third")
		w3.SetContent(widget.NewLabel("Third"))
		w3.Resize(fyne.NewSize(100, 100))
		w3.Show()
	}
}
