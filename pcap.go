package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"time"
)

var (
	deviceName  string = "eth0"
	snapshotLen int32  = 1024
	promiscuous bool   = false
	err         error
	timeout     time.Duration = -1 * time.Second
	handle      *pcap.Handle
	packetCount int = 0
)

func openPcap(pcapFile string) {
	// Open file instead of device
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
	}

}
func SaveAsPcap(PkgInfos []gopacket.Packet) {
	// Open output pcap file and write header
	f, _ := os.Create("test.pcap")
	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
	defer f.Close()

	// Open the device for capturing
	//handle, err = pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		fmt.Printf("Error opening device %s: %v", deviceName, err)
		os.Exit(1)
	}
	defer handle.Close()

	// Start processing packets
	//packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for _, packet := range PkgInfos {
		// Process packet here
		fmt.Println(packet)
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packetCount++

		// Only capture 100 and then stop
		if packetCount > 100 {
			break
		}
	}
}

//
//func SaveAsPcap(PkgInfos []gopacket.Packet) {
//	length:=len(PkgInfos)
//	// Open output pcap file and write header
//	FName:=time.Now().Format("2006-01-02T15:04:05")
//	f, _ := os.Create(FName+".pcap")
//	w := pcapgo.NewWriter(f)
//	err:=w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
//	if err!=nil{
//		fmt.Println("出错了1")
//	}
//	defer f.Close()
//	// Open the device for capturing
//	defer handle.Close()
//	// Start processing packets
//	for _,packet:=range PkgInfos{
//		// Process packet here
//		err =w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
//		if length<0{
//			break
//		}
//	}
//	if err!=nil{
//		fmt.Println("出错了2")
//	}
//	fmt.Println("保存成功为"+FName+".pcap")
//	return
//}
