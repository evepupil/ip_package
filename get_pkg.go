package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"strconv"
)

//设备名：pcap.FindAllDevs()返回的设备的Name
//snaplen：捕获一个数据包的多少个字节，一般来说对任何情况65535是一个好的实践，如果不关注全部内容，只关注数据包头，可以设置成1024
//promisc：设置网卡是否工作在混杂模式，即是否接收目的地址不为本机的包
//timeout：设置抓到包返回的超时。如果设置成30s，那么每30s才会刷新一次数据包；设置成负数，会立刻刷新数据包，即不做等待
//要记得释放掉handle

var (
	device       string = "eth0"
	snapshot_len int32  = 1024
)

func GetPkg(device string) {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process packet here
		//fmt.Println(packet.String())
		packet.Data()
		packet.String()
	}

}

//func packetString(p *packet) string {
//	var b bytes.Buffer
//	fmt.Fprintf(&b, "PACKET: %d bytes", len(p.Data()))
//	if p.metadata.Truncated {
//		b.WriteString(", truncated")
//	}
//	if p.metadata.Length > 0 {
//		fmt.Fprintf(&b, ", wire length %d cap length %d", p.metadata.Length, p.metadata.CaptureLength)
//	}
//	if !p.metadata.Timestamp.IsZero() {
//		fmt.Fprintf(&b, " @ %v", p.metadata.Timestamp)
//	}
//	b.WriteByte('\n')
//	for i, l := range p.layers {
//		fmt.Fprintf(&b, "- Layer %d (%02d bytes) = %s\n", i+1, len(l.LayerContents()), LayerString(l))
//	}
//	return b.String()
//}

func setFlliter(device string, p string, port int) {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// Set filter
	var filter string = p + " and port " + strconv.Itoa(port)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing " + p + " port " + strconv.Itoa(port) + " packets.")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Do something with a packet here.
		fmt.Println(packet)
	}
}
