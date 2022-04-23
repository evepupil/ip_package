package main

import (
	"IP_pkg_analyze/app/ip"
	"IP_pkg_analyze/app/util"
)

func main() {
	// Find all devices
	//b:= get_if_list()
	util.InitFont()
	//getPkg(b[0].NPFName)
	//go GetPkg(b[0].NPFName)
	//go saveAsPcap(ip.PkgInfos)
	//go ip.OpenPcap("test.pcap")
	ip.Run() //会阻塞
}
