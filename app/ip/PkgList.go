package ip

import (
	"fmt"
	fyne2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	//"time"
)

var PkgStringList []string
var PkgList *widget.List

func LoadPkgList() *fyne2.Container {
	pkgText := []string{
		"No.     ", "Time                          ", "Source              ", "Dst                      ",
		"Protocol   ", "Length     ", "Info                                 ",
	}
	pkgTextContainer := container.NewHBox()
	for _, v := range pkgText {
		pkgTextContainer.Add(widget.NewButton(v, func() {

		}))
	}
	//PkgStringList := binding.NewStringList()
	//PkgList := widget.NewListWithData(PkgStringList,
	//	func() fyne2.CanvasObject {
	//		return widget.NewLabel("")
	//	},
	//	func(item binding.DataItem, object fyne2.CanvasObject) {
	//		i := item.(binding.String)
	//		l := object.(*widget.Label)
	//		s, _ := i.Get()
	//		l.SetText(s)
	//	})
	PkgList = widget.NewList(
		func() int {
			return len(PkgStringList)
		},
		func() fyne2.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne2.CanvasObject) {
			o.(*widget.Label).SetText(PkgStringList[i])
		})

	PkgList.OnSelected = func(id widget.ListItemID) {
		packet := Map_Pkg_Infos[Pkgs[id]]
		fmt.Println(len(Map_Pkg_Infos), len(Pkgs), len(AllPkgInfos))
		NewLayersData(id+1, packet)
		LayersWidget.Refresh()
		NewPkgInfoData(packet)
		//pkg.Set(PkgBytes2String(ip.PkgInfos[id].Data()))

	}
	s := widget.NewSeparator()
	PkgListContainer := container.NewBorder(pkgTextContainer, nil, nil, s, PkgList)
	PkgListContainer.Resize(fyne2.NewSize(1600, 280))
	return PkgListContainer

}
func ReLoadPkgList(Pkgs []PkgRow) {
	pkgStringList := []string{}
	//
	for _, v := range Pkgs {
		pkgStringList = append(pkgStringList, v.FormatePkgListInfo())
	}
	PkgStringList = pkgStringList
}
func AddList(strings []string, string2 string) {
	strings = append(strings, string2)
}
