package ip

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"strconv"
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
	if device == "eth0" {
		return 1, errors.New("请先选择网卡")
	}
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
		return 1, err
	}
	defer handle.Close()
	rawBytes := []byte(pkg.Payload)
	if pkg.Source_mac == "" || pkg.Dest_mac == "" {
		return 1, errors.New("mac地址不能为空")
	}
	source_mac_bytes := String2HexBytes(pkg.Source_mac)
	dest_mac_bytes := String2HexBytes(pkg.Dest_mac)
	ethernetLayer := &layers.Ethernet{
		SrcMAC:       append(net.HardwareAddr{}, source_mac_bytes...),
		DstMAC:       append(net.HardwareAddr{}, dest_mac_bytes...),
		EthernetType: layers.EthernetTypeIPv4,
	}
	ipLayer := &layers.IPv4{}
	tcpLayer := &layers.TCP{}
	if pkg.Source_ip != "" {
		src_port, _ := strconv.Atoi(pkg.Source_port)
		dest_port, _ := strconv.Atoi(pkg.Dest_port)
		ipLayer = &layers.IPv4{
			IHL:      5,
			SrcIP:    append(net.IP{}, IpStr2Bytes(pkg.Source_ip)...),
			DstIP:    append(net.IP{}, IpStr2Bytes(pkg.Dest_ip)...),
			Version:  4,
			TTL:      255,
			Protocol: layers.IPProtocolTCP,
			Length:   20,
		}
		tcpLayer = &layers.TCP{
			SrcPort:  layers.TCPPort(src_port),
			DstPort:  layers.TCPPort(dest_port),
			Checksum: 0,
		}
		tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	}
	// And create the packet with the layers
	options.FixLengths = true
	buffer = gopacket.NewSerializeBuffer()
	var pack_err error = nil
	if ipLayer.SrcIP == nil {
		pack_err = gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			gopacket.Payload(rawBytes),
		)
	} else {
		pack_err = gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer,
			tcpLayer,
			gopacket.Payload(rawBytes),
		)
	}
	fmt.Println(buffer.Layers())
	if pack_err != nil {
		return 1, pack_err
	}
	outgoingPacket := buffer.Bytes()
	err := handle.WritePacketData(outgoingPacket)
	if err == nil {
		return 0, nil
	}
	return 1, err
}
func String2HexBytes(str string) []byte {
	mac_bytes := []byte{}
	for i := 0; i < len(str); i++ {
		s := str[i : i+2]
		i++
		hex_data, _ := hex.DecodeString(s)
		mac_bytes = append(mac_bytes, hex_data[0])
	}
	return mac_bytes
}
func IpStr2Bytes(str string) []byte {
	res := []byte{}
	dot_index := 0
	for i := 0; i < len(str); i++ {
		if str[i] == '.' {
			v, _ := strconv.Atoi(str[dot_index:i])
			res = append(res, byte(v))
			dot_index = i + 1
		}
	}
	v, _ := strconv.Atoi(str[dot_index:len(str)])
	res = append(res, byte(v))
	return res
}
