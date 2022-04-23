package util

import (
	"fmt"
	"github.com/flopp/go-findfont"
	"os"
	"strings"
)

func InitFont() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//楷体:simkai.ttf
		//黑体:simhei.ttf
		if strings.Contains(path, "simhei.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
	fmt.Println("=====字体初始化成功========")
}
