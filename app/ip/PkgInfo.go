package ip

import (
	"fyne.io/fyne/v2"
	fyne2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gopacket"
	//"github.com/google/gopacket"
)

var (
	PkgInfoWidget fyne2.Widget
	PkgCharWidget fyne2.Widget
	PkgMetaData   [][]string
	PkgCharData   [][]string
)

func LoadPkgInfo() *fyne2.Container {
	speparator := widget.NewSeparator()
	PkgInfoWidget = widget.NewTable(func() (int, int) {
		if len(PkgMetaData) != 0 {
			return len(PkgMetaData), len(PkgMetaData[0])
		}
		return 0, 0
	},
		func() fyne.CanvasObject {
			return widget.NewLabel("00")
		},
		func(tableCellID widget.TableCellID, object fyne.CanvasObject) {
			l := object.(*widget.Label)
			l.SetText(PkgMetaData[tableCellID.Row][tableCellID.Col])
		})
	PkgCharWidget = widget.NewTable(func() (int, int) {
		if len(PkgCharData) != 0 {
			return len(PkgCharData), len(PkgCharData[0])
		}
		return 0, 0
	},
		func() fyne.CanvasObject {
			return widget.NewLabel("00")
		},
		func(tableCellID widget.TableCellID, object fyne.CanvasObject) {
			l := object.(*widget.Label)
			l.SetText(PkgCharData[tableCellID.Row][tableCellID.Col])
		})
	PkgInfoWidget.(*widget.Table).OnSelected = func(id widget.TableCellID) {
		PkgCharWidget.(*widget.Table).Select(id)
	}
	PkgCharWidget.(*widget.Table).OnSelected = func(id widget.TableCellID) {
		PkgInfoWidget.(*widget.Table).Select(id)
	}
	InfoContainer := container.NewBorder(nil, nil, nil, nil, PkgInfoWidget)
	CharContainer := container.NewBorder(nil, nil, nil, nil, PkgCharWidget)
	InfoContainer.Resize(fyne2.NewSize(560, 240))
	CharContainer.Resize(fyne2.NewSize(560, 240))
	InfoContainer.Move(fyne2.NewPos(0, 0))
	CharContainer.Move(fyne2.NewPos(600, 0))
	InfoAndCharContainer := container.NewWithoutLayout(InfoContainer, CharContainer)
	PkgInfoContainer := container.NewBorder(nil, speparator, nil, nil, InfoAndCharContainer)
	PkgInfoContainer.Resize(fyne2.NewSize(1600, 240))
	return PkgInfoContainer
}
func NewPkgInfoData(packet gopacket.Packet) {
	PkgMetaData = PkgBytes2StringSlice(packet.Data())
	PkgCharData = PkgBytes2AsciiSlice(packet.Data())
	PkgInfoWidget.Refresh()
	PkgCharWidget.Refresh()
}
func PkgBytes2String(PkgBytes []byte) string {
	res := ""
	for _, v := range PkgBytes {
		res += byte2Hex(v)
	}
	return res
}
func PkgBytes2AsciiSlice(PkgBytes []byte) [][]string {
	res := [][]string{}
	for i := 0; i < len(PkgBytes); {
		temp := []string{}
		c := 0
		for i < len(PkgBytes) && c != 16 {
			temp = append(temp, Byte2AscllString(PkgBytes[i]))
			if i%16 == 7 {
				temp = append(temp, " ")
			}
			i++
			c++
		}
		if c == 16 {
			res = append(res, temp)
			c = 0
		} else if i == len(PkgBytes) { //如果最后不足一行则补齐空string
			for len(temp) < 17 {
				temp = append(temp, " ")
			}
			res = append(res, temp)
		}
	}
	return res
}

func Byte2AscllString(b byte) string {
	bint := int(b)
	if bint >= 33 && bint <= 126 {
		return string(b)
	}
	return "."

}
func PkgBytes2StringSlice(PkgBytes []byte) [][]string {
	res := [][]string{}
	for i := 0; i < len(PkgBytes); {
		temp := []string{}
		c := 0
		for i < len(PkgBytes) && c != 16 {
			temp = append(temp, byte2Hex(PkgBytes[i]))
			if i%16 == 7 {
				temp = append(temp, " ")
			}
			i++
			c++
		}
		if c == 16 {
			res = append(res, temp)
			c = 0
		} else if i == len(PkgBytes) { //如果最后不足一行则补齐空string
			for len(temp) < 17 {
				temp = append(temp, " ")
			}
			res = append(res, temp)
		}
	}
	return res
}
func byte2Hex(b byte) string {
	care := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	bb := int(b)
	res := ""
	if b == 0 {
		return "00"
	}
	if b < 16 {
		return "0" + care[bb%16]
	}
	for bb > 0 {
		res = care[bb%16] + res
		bb /= 16
	}
	return res
}
