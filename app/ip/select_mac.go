package ip

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"strings"
)

type IfaceInfo struct {
	NPFName     string
	Description string
	NickName    string
	IPv4        string
}

var IPs []string
var IfaceInfoList []IfaceInfo
var VaildIfaces []IfaceInfo

func Get_if_list() []IfaceInfo {
	var ifaceInfoList []IfaceInfo

	// 得到所有的(网络)设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	interface_list, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range interface_list {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			log.Fatal(err)
		}
		address, err := byName.Addrs()
		ifaceInfoList = append(ifaceInfoList, IfaceInfo{NickName: byName.Name, IPv4: address[1].String()})
		IfaceInfoList = append(IfaceInfoList, IfaceInfo{NickName: byName.Name, IPv4: address[1].String()})
	}

	//打印设备信息
	//fmt.Println("Devices found:")
	//for _, device := range devices {
	//	fmt.Println("\nName: ", device.Name)
	//	fmt.Println("Description: ", device.Description)
	//	fmt.Println("Devices addresses: ", device.Description)
	//	for _, address := range device.Addresses {
	//		fmt.Println("- IP address: ", address.IP)
	//		fmt.Println("- Subnet mask: ", address.Netmask)
	//	}
	//}
	var vaildIfaces []IfaceInfo
	for _, device := range devices {
		for _, address := range device.Addresses {
			//fmt.Println("add",address)
			for i, ifaceinfo := range ifaceInfoList {
				if strings.Contains(ifaceinfo.IPv4, address.IP.String()) {
					vaildIfaces = append(vaildIfaces, IfaceInfo{NPFName: device.Name, Description: device.Description, NickName: ifaceinfo.NickName, IPv4: ifaceinfo.IPv4})
					VaildIfaces = append(VaildIfaces, IfaceInfo{NPFName: device.Name, Description: device.Description, NickName: ifaceinfo.NickName, IPv4: ifaceinfo.IPv4})
					IfaceInfoList[i].NPFName = device.Name
					IfaceInfoList[i].Description = device.Description
					ip := ifaceinfo.IPv4
					for i := 0; i < len(ip); i++ {
						if ip[i] == '/' {
							ip = ip[0:i]
							break
						}
					}
					IPs = append(IPs, ip)
					IP = ip
					break
				}
			}
		}
	}
	return vaildIfaces
}

func selectMac() {
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print device information
	fmt.Println("Devices found:", len(devices))
	s := 0
	for _, d := range devices {
		s += 1
		fmt.Println("\nName: ", d.Name)
		fmt.Println("Description: ", d.Description)
		fmt.Println("Devices addresses: ", d.Flags)

		//for _, address := range d.Addresses {
		//	fmt.Println("- IP address: ", address.IP)
		//	fmt.Println("- Subnet mask: ", address.Netmask)
		//}
		//setFlliter(d.Name,"tcp",3306)
		fmt.Println("---------------------------------------- ", s)
	}
}
