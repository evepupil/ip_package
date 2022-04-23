package ip

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var OnSelectIface int

func LoadMenus(w fyne.Window, app fyne.App) {
	tools_key := []string{"文件(F)", "过滤(E)", "排序(V)", "切换(W)", "模式(M)", "发送(S)", "捕获(C)"}
	//"捕获(C)", "分析(A)", "统计(S)", "电话(Y)",
	//"无线(W)", "工具(T)", "帮助(H)",
	var tools = [][]string{{"保存", "打开"}, {"过滤器设置"}, {"源IP(递增)", "源IP(递减)", "目的IP(递增)", "目的IP(递减)", "源端口(递增)", "源端口(递减)", "目的端口(递增)",
		"目的端口(递减)", "长度(递增)", "长度(递减)"},
		{"选择网卡"}, {"严格模式", "混杂模式"}, {"发送数据包"}, {"暂停/继续"}}
	var toolsFunc = [][]func(){{
		func() {
			dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
				if closer != nil {
					uri := closer.URI().String()
					for i := 0; i < len(uri); i++ {
						if uri[i] == ':' {
							uri = uri[i-1 : len(uri)]
						}
					}
					path, err := SaveAsPcap(uri, AllPkgInfos)
					if err == nil {
						dialog.NewInformation("提示", "成功保存至 "+path, w).Show()
					}
				}
			}, w).Show()

		},
		func() {
			dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
				if err != nil {
					log.Fatal("打开文件错误")
				}
				if closer != nil {
					uri := closer.URI().String()
					for i := 0; i < len(uri); i++ {
						if uri[i] == ':' {
							uri = uri[i-1 : len(uri)]
						}
					}
					OpenPcap(uri)
					if IsGetPkg {
						Cancel()
					}
				}

			}, w).Show()
		},
	},
		{
			func() {
				filter_window := app.NewWindow("请输入过滤器")
				filter_window.SetIcon(Icon)
				source_ip_entry := widget.NewEntry()
				dest_ip_entry := widget.NewEntry()
				source_port_entry := widget.NewEntry()
				dest_port_entry := widget.NewEntry()
				form := &widget.Form{
					Items: []*widget.FormItem{ // we can specify items in the constructor
					},
					OnSubmit: func() { // optional, handle form submission
						flliterBySourceIp(source_ip_entry.Text, AllPkgs)
						SourceIp_filter = source_ip_entry.Text
						flliterByDestIp(dest_ip_entry.Text, Pkgs)
						DestIp_filter = dest_ip_entry.Text
						flliterBySourcePort(source_port_entry.Text, Pkgs)
						SourcePort_filter = source_port_entry.Text
						flliterByDestPort(dest_port_entry.Text, Pkgs)
						DestPort_filter = dest_port_entry.Text
					},
				}
				source_ip_entry.Text = SourceIp_filter
				dest_ip_entry.Text = DestIp_filter
				source_port_entry.Text = SourcePort_filter
				dest_port_entry.Text = DestPort_filter
				// we can also append items
				form.Append("源IP", source_ip_entry)
				form.Append("目的IP", dest_ip_entry)
				form.Append("源Port", source_port_entry)
				form.Append("目的Port", dest_port_entry)
				form.SubmitText = "确定"
				//form.Append("端口")
				filter_window.SetContent(form)
				filter_window.Resize(fyne.NewSize(300, 150))
				filter_window.Show()
				fmt.Println(len(Map_Pkg_Infos), len(Pkgs), len(PkgInfos))
				filter_window.Close()
			},
		},
		{
			func() {
				SortBySourceIp = INCRESE
				SortByDestIp = 0
				SortBySourcePort = 0
				SortByDestPort = 0
				SortByLength = 0
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortBySourceIp(Pkgs))
			}, func() {
				SortBySourceIp = DECRESE
				SortByDestIp = 0
				SortBySourcePort = 0
				SortByDestPort = 0
				SortByLength = 0
				//Sort_tpyes = 1
				fmt.Println("ww")
				ReLoadPkgList(PkgSortBySourceIpReverse(Pkgs))
			},
			func() {
				SortByDestIp = INCRESE
				SortBySourceIp = 0
				SortBySourcePort = 0
				SortByDestPort = 0
				SortByLength = 0
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortByDestIp(Pkgs))
			}, func() {
				SortByDestIp = DECRESE
				SortBySourceIp = 0
				SortBySourcePort = 0
				SortByDestPort = 0
				SortByLength = 0
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortByDestIpReverse(Pkgs))
			}, func() {
				SortBySourcePort = INCRESE
				SortBySourceIp = 0
				SortByDestIp = 0
				SortByDestPort = 0
				SortByLength = 0
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortBySourcePort(Pkgs))
			}, func() {
				SortBySourcePort = DECRESE
				SortByDestIp = 0
				SortBySourceIp = 0
				SortByDestPort = 0
				SortByLength = 0
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortBySourcePortReverse(Pkgs))
			}, func() {
				SortByDestIp = 0
				SortBySourceIp = 0
				SortBySourcePort = 0
				SortByDestPort = INCRESE
				SortByLength = 0
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortByDestPort(Pkgs))
			}, func() {
				SortByDestIp = 0
				SortBySourceIp = 0
				SortBySourcePort = 0
				SortByDestPort = DECRESE
				SortByLength = 0
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortByDestPortReverse(Pkgs))
			}, func() {
				SortByDestIp = 0
				SortBySourceIp = 0
				SortBySourcePort = 0
				SortByDestPort = 0
				SortByLength = INCRESE
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortByLength(Pkgs))
			}, func() {
				SortByDestIp = 0
				SortBySourceIp = 0
				SortBySourcePort = 0
				SortByDestPort = 0
				SortByLength = DECRESE
				//Sort_tpyes = 1
				ReLoadPkgList(PkgSortByLengthReverse(Pkgs))
			},
		},
		{
			func() {
				c := container.NewVBox()
				iface_window := app.NewWindow("选择网卡")
				iface_window.SetIcon(Icon)
				valid := []int{}
				info := []string{}
				for i, v := range IfaceInfoList {
					hint := "(不可用)"
					for _, vv := range VaildIfaces {
						if v.IPv4 == vv.IPv4 {
							hint = ""
							valid = append(valid, i)
						}
					}
					info = append(info, strconv.Itoa(i+1)+":"+v.NickName+":"+v.IPv4+hint)
				}
				radio := widget.NewRadioGroup(info, func(value string) {
					str_id := ""
					for i := 0; i < len(value); i++ {
						if value[i] == ':' {
							str_id = value[0:i]
							break
						}
					}
					id, _ := strconv.Atoi(str_id)
					id--
					for _, v := range valid {
						if id == v {
							if IsGetPkg {
								Cancel()
								IsGetPkg = false
							}
							Ctx, Cancel = context.WithCancel(context.Background())
							time.Sleep(1000 * time.Millisecond)
							go GetPkg(Ctx, IfaceInfoList[id].NPFName)
							IsGetPkg = true
							IfaceName.Set("当前选择的网卡为: " + value[2:])
							break
						}
					}
					iface_window.Close()
				})
				c.Add(radio)
				iface_window.SetContent(c)
				iface_window.Show()
			},
		},
		{
			func() {
				Promiscuous = false
				dialog.NewInformation("提示", "已开启严格模式，仅抓获与本机相关的IP包", w).Show()
			},
			func() {
				Promiscuous = true
				dialog.NewInformation("提示", "已开启混杂模式，抓获所有经过的IP包", w).Show()

			},
		},
		{
			func() {
				send_window := app.NewWindow("自定义数据包")
				send_window.SetIcon(Icon)
				source_ip_entry := widget.NewEntry()
				dest_ip_entry := widget.NewEntry()
				source_port_entry := widget.NewEntry()
				dest_port_entry := widget.NewEntry()
				source_mac_entry := widget.NewEntry()
				dest_mac_entry := widget.NewEntry()
				pay_load_entry := widget.NewEntry()
				form := &widget.Form{
					Items: []*widget.FormItem{ // we can specify items in the constructor
					},
					OnSubmit: func() { // optional, handle form submission
						mac_regexp := regexp.MustCompile("[A-F0-9]{12}")
						if !mac_regexp.MatchString(strings.ToUpper(source_mac_entry.Text)) {
							fmt.Println(strings.ToUpper(source_mac_entry.Text))
							dialog.NewInformation("提示", "源mac: "+source_mac_entry.Text+" 不符合mac地址格式", w).Show()
							return
						}
						if !mac_regexp.MatchString(strings.ToUpper(dest_mac_entry.Text)) {
							dialog.NewInformation("提示", "目的mac: "+dest_mac_entry.Text+" 不符合mac地址格式", w).Show()
							return
						}
						if source_ip_entry.Text != "" && dest_ip_entry.Text != "" {

							ip_regexip := regexp.MustCompile("((2((5[0-5])|([0-4]\\d)))|([0-1]?\\d{1,2}))(\\.((2((5[0-5])|([0-4]\\d)))|([0-1]?\\d{1,2}))){3}")
							if !ip_regexip.MatchString(source_ip_entry.Text) {
								dialog.NewInformation("提示", "源IP: +"+source_ip_entry.Text+" 不符合IP格式", w).Show()
								return
							}
							if !ip_regexip.MatchString(dest_ip_entry.Text) {
								dialog.NewInformation("提示", "目的IP: +"+dest_ip_entry.Text+" 不符合IP格式", w).Show()
								return
							}
							if source_port_entry.Text == "" || dest_port_entry.Text == "" {
								dialog.NewInformation("提示", "端口号不能为空！", w).Show()
								return
							}
						} else {
							source_ip_entry.SetText("")
							dest_ip_entry.SetText("")
						}
						pkg := SendPkgInfo{
							Source_ip:   source_ip_entry.Text,
							Dest_ip:     dest_ip_entry.Text,
							Source_port: source_port_entry.Text,
							Dest_port:   dest_port_entry.Text,
							Source_mac:  source_mac_entry.Text,
							Dest_mac:    dest_mac_entry.Text,
							Payload:     pay_load_entry.Text,
							Type:        "TCP",
						}
						status, err := SendPkg(pkg)
						if status != 0 {
							dialog.NewInformation("提示", "发送失败，原因:"+err.Error(), w).Show()
						} else {
							dialog.NewInformation("提示", "发送成功!", w).Show()
						}
					},
				}
				form.SubmitText = "发送"
				source_ip_entry.Text = SourceIp_filter
				dest_ip_entry.Text = DestIp_filter
				source_port_entry.Text = SourcePort_filter
				dest_port_entry.Text = DestPort_filter
				// we can also append items
				form.Append("源Mac", source_mac_entry)
				form.Append("目的Mac", dest_mac_entry)
				form.Append("源IP", source_ip_entry)
				form.Append("目的IP", dest_ip_entry)
				form.Append("源Port", source_port_entry)
				form.Append("目的Port", dest_port_entry)
				form.Append("负载", pay_load_entry)
				form.SubmitText = "确定"
				//form.Append("端口")

				source_ip_entry.SetText("192.168.15.5")
				dest_ip_entry.SetText("192.168.152.51")
				source_port_entry.SetText("830")
				dest_port_entry.SetText("934")
				source_mac_entry.SetText("0000ffffffff")
				dest_mac_entry.SetText("0000ffffffff")
				pay_load_entry.SetText("ffffffffffff")
				send_window.SetContent(form)
				send_window.Resize(fyne.NewSize(400, 500))
				send_window.Show()
				//send_window.Close()
			},
		},
		{
			func() {
				Interupt = !Interupt
			},
		},
	}
	Menus := []*fyne.Menu{}
	for i, key := range tools_key {
		tMenu := []*fyne.MenuItem{}
		for index, item := range tools[i] {
			newMenuItem := fyne.NewMenuItem(item, toolsFunc[i][index])
			newMenuItem.IsQuit = false
			tMenu = append(tMenu, newMenuItem)
		}
		Menus = append(Menus, fyne.NewMenu(key, tMenu...))
	}
	w.SetMainMenu(fyne.NewMainMenu(Menus...))

}
