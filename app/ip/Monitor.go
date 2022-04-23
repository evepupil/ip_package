package ip

import (
	"fmt"
	fyne2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var FlowsStr binding.String
var IfaceName binding.String

func LoadMonitor() *fyne2.Container {
	FlowsStr = binding.NewString()
	IfaceName = binding.NewString()
	IfaceName.Set("当前未选择网卡")
	FlowsStr.Set(fmt.Sprintf("\rDown:%.2fkb/s \t Up:%.2fkb/s                                                                                                                             ", float32(0)/1024/1, float32(0)/1024/1))
	monitorWidget := widget.NewLabelWithData(FlowsStr)
	ifaceWidget := widget.NewLabelWithData(IfaceName)
	monitorContainer := container.NewHBox(monitorWidget, ifaceWidget)
	//monitorContainer := container.NewBorder(nil, nil, nil, ifaceWidget, monitorWidget)
	monitorContainer.Resize(fyne2.NewSize(560, 70))
	return monitorContainer
}
