package ip

import (
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
)

var (
	buffer  gopacket.SerializeBuffer
	options gopacket.SerializeOptions
)

type SendPkgInfo struct {
	Source_ip   string
	Dest_ip     string
	Source_port string
	Dest_port   string
	Source_mac  string
	Dest_mac    string
	Payload     string
	Type        string
}

func SendPkg(pkg SendPkgInfo) (int, error) {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
		return 1, err
	}
	defer handle.Close()
	// Send raw bytes over wire
	rawBytes := []byte(pkg.Payload)
	// Create a properly formed packet, just with
	// empty details. Should fill out MAC addresses,
	// IP addresses, etc.
	// This time lets fill out some information
	if pkg.Source_mac == "" || pkg.Dest_mac == "" {
		return 1, errors.New("mac地址不能为空")
	}
	ethernetLayer := &layers.Ethernet{
		SrcMAC: net.HardwareAddr{0xFF, 0xAA, 0xFA, 0xAA, 0xFF, 0xAA},
		DstMAC: net.HardwareAddr{0xBD, 0xBD, 0xBD, 0xBD, 0xBD, 0xBD},
	}
	if pkg.Source_ip != "" {

	}
	ipLayer := &layers.IPv4{
		SrcIP: net.IP{127, 0, 0, 1},
		DstIP: net.IP{8, 8, 8, 8},
	}

	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(4321),
		DstPort: layers.TCPPort(80),
	}
	// And create the packet with the layers
	buffer = gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, options,
		ethernetLayer,
		ipLayer,
		tcpLayer,
		gopacket.Payload(rawBytes),
	)
	return 0, nil
}
