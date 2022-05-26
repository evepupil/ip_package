package ip

import (
	"context"
	fyne2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var DeviceName string
var A fyne2.App
var IP string
var Icon fyne2.Resource
var Ctx context.Context
var Cancel context.CancelFunc

func Run() {
	Ctx, Cancel = context.WithCancel(context.Background())
	A := app.New()
	w := A.NewWindow("简单的IP抓包工具")
	Icon, err = fyne2.LoadResourceFromPath("C:\\Users\\Evepupil\\Pictures\\QQ图片20220418155447.jpg")
	if err != nil {
		fyne2.LogError("icon加载失败", err)
	}
	w.SetIcon(Icon)
	w.SetMaster()
	LoadMenus(w, A)
	ifaces := Get_if_list()
	Layers := LoadLayers()
	PkgInfo := LoadPkgInfo()
	Monitor := LoadMonitor()
	PkgList := LoadPkgList()
	PkgListContainer := container.NewWithoutLayout(PkgList)
	Layers.Move(fyne2.NewPos(0, 280))
	PkgInfo.Move(fyne2.NewPos(0, 520))
	Monitor.Move(fyne2.NewPos(0, 760))
	AllContainer := container.NewWithoutLayout(PkgListContainer, Layers, PkgInfo, Monitor)
	w.SetContent(AllContainer)
	w.Resize(fyne2.NewSize(1600, 830))
	DeviceName = ifaces[0].NPFName
	w.Show()
	A.Run()
}
