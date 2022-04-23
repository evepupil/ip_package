package ip

import (
	"bytes"
	"fmt"
	fyne2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"strconv"
	//"time"
)

type LayersInfo struct {
	PkgInfos                                                                  []string
	LinkLayersInfos, NetWorkLayersInfos, TransportLayersInfos, AppLayersInfos []string
}

var (
	LayersWidget fyne2.Widget
	LayersData   map[string][]string
)

func LoadLayers() *fyne2.Container {
	LayersData = map[string][]string{}
	LayersWidget = widget.NewTreeWithStrings(LayersData)
	scrollC := container.NewScroll(LayersWidget)
	speparator := widget.NewSeparator()
	PkgInfoContainer := container.NewBorder(nil, speparator, nil, nil, scrollC)
	PkgInfoContainer.Resize(fyne2.NewSize(1600, 240))
	return PkgInfoContainer
}
func NewLayersData(FrameNo int, packet gopacket.Packet) map[string][]string {
	for k := range LayersData { //每次都将map初始化为空
		delete(LayersData, k)
	}
	PkgInfoBbranch, L0Infos := getPkgInfoData(FrameNo, packet)
	layersAll := []string{}
	layersAll = append(layersAll, PkgInfoBbranch)
	LayersData[PkgInfoBbranch] = L0Infos
	if packet.LinkLayer() != nil {
		LinkLayerBranch, L1Infos := getLinkLayerData(packet)
		LayersData[LinkLayerBranch] = L1Infos
		layersAll = append(layersAll, LinkLayerBranch)
	}
	if packet.NetworkLayer() != nil {
		NetWorkLayerBranch, L2Infos := getNetWorkLayerData(packet)
		LayersData[NetWorkLayerBranch] = L2Infos
		layersAll = append(layersAll, NetWorkLayerBranch)
	}
	if packet.TransportLayer() != nil {
		TransportLayerBranch, L3Infos := getTransportLayerData(packet)
		LayersData[TransportLayerBranch] = L3Infos
		layersAll = append(layersAll, TransportLayerBranch)
	}
	if packet.ApplicationLayer() != nil {
		AppLayerBranch, L4Infos := getAppLayerData(packet)
		LayersData[AppLayerBranch] = L4Infos
		layersAll = append(layersAll, AppLayerBranch)
	}
	LayersData[""] = layersAll
	return LayersData
}
func getPkgInfoData(FrameNo int, packet gopacket.Packet) (branch string, nodes []string) {
	var PkgInfoBuffer, InterfaceBuffer bytes.Buffer
	metadata := packet.Metadata()
	fmt.Fprintf(&PkgInfoBuffer, "Frame %d: %d bytes on wire (%d bits),%d bytes captured(%d bits) "+
		"on interface %s , id:%d", FrameNo, metadata.Length, metadata.Length*8, metadata.CaptureLength, metadata.CaptureLength*8,
		DeviceName, metadata.InterfaceIndex)
	branch = PkgInfoBuffer.String()
	fmt.Fprintf(&InterfaceBuffer, "Interface id: %d (%s)",
		metadata.InterfaceIndex, DeviceName) //设备信息
	time := metadata.Timestamp.Format("2006-01-02T15:04:05") //时间
	No := "Frame Number: " + strconv.Itoa(FrameNo)           //No
	FrameLength := "Frame Length: " + strconv.Itoa(metadata.Length) + "(" + strconv.Itoa(metadata.Length*8) + "bits)"
	CaptureLength := "Capture Length: " + strconv.Itoa(metadata.CaptureLength) + "(" + strconv.Itoa(metadata.CaptureLength*8) + "bits)"
	nodes = append(nodes, InterfaceBuffer.String(), "Arrival Time: "+time, No, FrameLength, CaptureLength)
	return

}
func getLinkLayerData(packet gopacket.Packet) (branch string, nodes []string) {
	linkLayerMetaData := packet.LinkLayer()
	var linkLayerInfoBuffer bytes.Buffer
	fmt.Fprintf(&linkLayerInfoBuffer, "%s , Src: %s , Dst: %s", linkLayerMetaData.LayerType().String(),
		linkLayerMetaData.LinkFlow().Src().String(), linkLayerMetaData.LinkFlow().Dst().String())
	branch = linkLayerInfoBuffer.String()
	Dst := "Destination: " + linkLayerMetaData.LinkFlow().Dst().String()
	Src := "Source: " + linkLayerMetaData.LinkFlow().Src().String()
	Type := "Type: IPV6(0x86dd)" //IPV6
	if linkLayerMetaData.LayerContents()[12] == 8 {
		Type = "Type: IPV4(0x0800)" //IPV4
	}
	nodes = append(nodes, Dst, Src, Type)
	return
}
func getNetWorkLayerData(packet gopacket.Packet) (branch string, nodes []string) {
	networkLayerMetaData := packet.NetworkLayer()
	src, dst := networkLayerMetaData.NetworkFlow().Src().String(), networkLayerMetaData.NetworkFlow().Dst().String()
	var networkLayerInfoBuffer bytes.Buffer
	fmt.Fprintf(&networkLayerInfoBuffer, "Internet Protocol Version %d, Src: %s, Dst: %s",
		networkLayerMetaData.LayerContents()[0]/16, src, dst)
	branch = networkLayerInfoBuffer.String()
	Version := hex2(4, networkLayerMetaData.LayerContents()[0]/16) + " .... = Version : " +
		strconv.Itoa(int(networkLayerMetaData.LayerContents()[0]/16))
	headlengthMetadata := int(networkLayerMetaData.LayerContents()[0] % 16)
	HeaderLength := " ...." + hex2(4, networkLayerMetaData.LayerContents()[0]%16) +
		" = Header Length " + strconv.Itoa(headlengthMetadata*4) + " bytes (" + strconv.Itoa(headlengthMetadata) + ")"
	TotalLength := "Total Length: " + strconv.Itoa(int(networkLayerMetaData.LayerContents()[2])*256+int(networkLayerMetaData.LayerContents()[3]))
	Identification := "Identification: 0x" + byte2HexString(networkLayerMetaData.LayerContents()[4]) + byte2HexString(networkLayerMetaData.LayerContents()[5]) +
		" (" + strconv.Itoa(int(networkLayerMetaData.LayerContents()[4])*256+int(networkLayerMetaData.LayerContents()[5])) + ")"
	FlagsMetaData := networkLayerMetaData.LayerContents()[6]
	Flags := "Flags: 0x" + byte2HexString(FlagsMetaData) + ", Don't fragment"
	if FlagsMetaData>>6&1 != 1 {
		Flags = "Flags: 0x" + byte2HexString(FlagsMetaData) + ", Set fragment"
	}
	FragmentOffset := "..." + strconv.Itoa(int(FlagsMetaData)/(1<<4)%2) + hex2(4, (FlagsMetaData)%(1<<4)) +
		hex2(8, networkLayerMetaData.LayerContents()[7]) + " = Fragment Offset: " + strconv.Itoa(int(networkLayerMetaData.LayerContents()[7])+int(FlagsMetaData)%(1<<4)*256+int(FlagsMetaData/(1<<4)%2)*4096)
	TimeToLive := "Time to Live: " + strconv.Itoa(int(networkLayerMetaData.LayerContents()[8]))
	Protocol := "Protocol: UDP (17)"
	if networkLayerMetaData.LayerContents()[9] == 6 {
		Protocol = "Protocol: TCP (6)"
	} else if networkLayerMetaData.LayerContents()[9] == 1 {
		Protocol = "Protocol: ICMP (1)"
	}
	HeaderChecksum := "0x" + byte2HexString(networkLayerMetaData.LayerContents()[10]) + byte2HexString(networkLayerMetaData.LayerContents()[11]) +
		" [validation disabled]"
	SourceAddress := "Source Address: " + src
	DestinationAddress := "Destination Address: " + dst
	nodes = append(nodes, Version, HeaderLength, TotalLength, Identification, Flags, FragmentOffset,
		TimeToLive, Protocol, HeaderChecksum, SourceAddress, DestinationAddress)
	return
}
func getTransportLayerData(packet gopacket.Packet) (string, []string) {
	if packet.NetworkLayer() != nil {
		networkLayerMetaData := packet.NetworkLayer()
		if networkLayerMetaData.LayerContents()[9] == 6 { //TCP
			return getTCPData(packet)
		} else if networkLayerMetaData.LayerContents()[9] == 1 { //ICMP
			return getICMPData(packet)
		} else if networkLayerMetaData.LayerContents()[9] == 17 { //UDP
			return getUDPData(packet)
		}
	}
	return "Cannot anlyse transport layer type", []string{}
}
func getAppLayerData(packet gopacket.Packet) (branch string, nodes []string) {
	//AppLayerMetaData:=packet.ApplicationLayer().LayerContents()
	var TransportProtcal gopacket.LayerType
	//layers.LayerTypeUDP
	if packet.TransportLayer() != nil {
		TransportProtcal = packet.TransportLayer().LayerType()
	}
	app := packet.ApplicationLayer().LayerContents()
	appBytes := packet.ApplicationLayer().LayerContents()
	branch = fmt.Sprintf("Data (%d bytes)", len(app))
	Data := fmt.Sprintf("Data:%s", PkgBytes2String(app))
	Length := fmt.Sprintf("[Length: %d bytes]", len(app))
	nodes = append(nodes, Data, Length)
	switch TransportProtcal {
	case layers.LayerTypeUDP:
		branch = fmt.Sprintf("Data (%d bytes)", len(app))
		Data := fmt.Sprintf("Data:%s", PkgBytes2String(app))
		Length := fmt.Sprintf("[Length: %d bytes]", len(app))
		nodes = append(nodes, Data, Length)
	case layers.LayerTypeTCP:
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		tcp, _ := tcpLayer.(*layers.TCP)
		if tcp.SrcPort.String() == "80(http)" || tcp.DstPort.String() == "80(http)" {
			branch = "Hypertext Transfer Protocol"
			infos := []string{}
			info := []byte{}
			for i := 0; i < len(app); i++ {
				if i == len(app)-2 {
					break
				}
				if app[i] == 13 && app[i+1] == 10 {
					t := ""
					for _, b := range info {
						t += Byte2AscllString(b)
					}
					infos = append(infos, t)
					info = []byte{}
				}
				if app[i] != 13 && app[i] != 10 {
					info = append(info, app[i])
				}
			}
			nodes = infos
		} else {
			branch = fmt.Sprintf("Data (%d bytes)", len(app))
			Data := fmt.Sprintf("Data:%s", PkgBytes2String(app))
			Length := fmt.Sprintf("[Length: %d bytes]", len(app))
			nodes = append(nodes, Data, Length)
		}

	}
	if appBytes[0] == 2 {
		branch = "OICQ - IM software, popular in China"
		flag := "Flag: Oicq packet (0x02)"
		Version := "Version: 0x" + byte2HexString(appBytes[1]) + byte2HexString(appBytes[2])
		Command := "cannot anlyse"
		status := int(appBytes[3])*256 + int(appBytes[4])
		if status == 129 {
			Command = "Command: Get status of friend (129)"
		} else if status == 23 {
			Command = "Command: Receive message (23)"
		} else if status == 88 {
			Command = "Command: Download group friend (88)"
		}
		Sequence := "Sequence: " + strconv.Itoa(int(appBytes[5])*256+int(appBytes[6]))
		Data := "Data(OICQ Number,if sender is client): " + strconv.Itoa(int(appBytes[7])*256+
			int(appBytes[8])+int(appBytes[9])*256+int(appBytes[10]))
		nodes = append(nodes, flag, Version, Command, Sequence, Data)
	}
	return
}
func getTCPData(packet gopacket.Packet) (branch string, nodes []string) {
	//TCPFlow:=packet.TransportLayer().TransportFlow()
	TCPMetaData := packet.TransportLayer().LayerContents()
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	tcp, _ := tcpLayer.(*layers.TCP)
	branch = fmt.Sprintf("Transmission Control Protocol, Src Port: %s, Dst Port: %s, Seq: %d, Ack: %d, Len: %d",
		tcp.SrcPort.String(), tcp.DstPort.String(), tcp.Seq, tcp.Ack, len(packet.TransportLayer().LayerPayload()))
	SourcePort := fmt.Sprintf("Source Port: %s", tcp.SrcPort.String())
	DestinationPort := fmt.Sprintf("Dest Port: %s", tcp.DstPort.String())
	TCPSegmentLen := fmt.Sprintf("[TCP Segment Len: %d]", len(tcp.Payload))
	SequenceNumber := fmt.Sprintf("Sequence Number: %d    (relative sequence number)")
	SequenceNumberRaw := fmt.Sprintf("Sequence Number (raw): %d", tcp.Seq)
	NextSequenceNumber := fmt.Sprintf("[Next Sequence Number: %d    (relative sequence number)]")
	AcknowledgmentNumber := fmt.Sprintf("Acknowledgment Number: 1513    (relative ack number)")
	AcknowledgmentNumberRaw := fmt.Sprintf("Acknowledgment number (raw): %d", tcp.Ack)
	HeaderLength := fmt.Sprintf("%s .... = Header Length: %d bytes (%d)", hex2(4, TCPMetaData[12]>>4), TCPMetaData[12]>>4*byte(4), TCPMetaData[12]>>4)
	Flags := fmt.Sprintf("Flags: 0x011 (FIN, ACK)")
	Window := fmt.Sprintf("Window: %d", tcp.Window)
	Checksum := fmt.Sprintf("Checksum: 0x%s%s [unverified]", byte2HexString(byte(tcp.Checksum>>4)), byte2HexString(byte(tcp.Checksum)))
	UrgentPointer := fmt.Sprintf("Urgent Pointer: %d", tcp.Urgent)
	nodes = append(nodes, SourcePort, DestinationPort, TCPSegmentLen, SequenceNumber,
		SequenceNumberRaw, NextSequenceNumber, AcknowledgmentNumber, AcknowledgmentNumberRaw, HeaderLength, Flags, Window, Checksum,
		UrgentPointer)
	return
}
func getUDPData(packet gopacket.Packet) (branch string, nodes []string) {
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	udp := udpLayer.(*layers.UDP)
	branch = fmt.Sprintf("User Datagram Protocol, Src Port: %s, Dst Port: %s",
		udp.SrcPort.String(), udp.DstPort.String())
	SourcePort := fmt.Sprintf("Source Port: %s", udp.SrcPort.String())
	DestinationPort := fmt.Sprintf("Destination Port: %s", udp.DstPort.String())
	Length := fmt.Sprintf("Length: %d", udp.Length)
	Checksum := fmt.Sprintf("Checksum: 0x%s%s [unverified]", byte2HexString(byte(udp.Checksum>>4)), byte2HexString(byte(udp.Checksum)))
	UDPpayload := fmt.Sprintf("UDP payload (%d bytes)", len(udp.Payload))
	nodes = append(nodes, SourcePort, DestinationPort, Length, Checksum, UDPpayload)
	return
}
func getICMPData(packet gopacket.Packet) (branch string, nodes []string) {
	icmp := packet.NetworkLayer().LayerContents()
	branch = "Internet Control Message Protocol"
	Type := fmt.Sprintf("Type: %d (Echo (ping) request)", icmp[0])
	Code := fmt.Sprintf("Code: %d", icmp[1])
	Checksum := fmt.Sprintf("0x%s%s [correct]", byte2HexString(icmp[2]), byte2HexString(icmp[3]))
	IdentifierBE := fmt.Sprintf("Identifier (BE): %d (0x%s%s)", int(icmp[4])*256+int(icmp[5]), byte2HexString(icmp[4]), byte2HexString(icmp[5]))
	IdentifierLE := fmt.Sprintf("Identifier (LE): %d (0x%s%s)", int(icmp[5])*256+int(icmp[4]), byte2HexString(icmp[5]), byte2HexString(icmp[4]))
	SequenceNumberBE := fmt.Sprintf("Sequence Number (BE): %d (0x%s%s)", int(icmp[6])*256+int(icmp[7]), byte2HexString(icmp[6]), byte2HexString(icmp[7]))
	SequenceNumberLE := fmt.Sprintf("Sequence Number (BE): %d (0x%s%s)", int(icmp[7])*256+int(icmp[6]), byte2HexString(icmp[7]), byte2HexString(icmp[6]))
	//Data
	nodes = append(nodes, Type, Code, Checksum, IdentifierBE, IdentifierLE, SequenceNumberBE, SequenceNumberLE)
	return
}
func getWebsocketData(packet gopacket.Packet) (branch string, nodes []string) {
	branch = "WebSocket"
	return
}

func hex2(length int, hex byte) (res string) {
	hex1 := int(hex)
	c := 0
	for length > 0 {
		res = strconv.Itoa(hex1%2) + res
		hex1 /= 2
		c++
		if c == 4 {
			res = " " + res
			c = 0
		}
		length--
	}
	return
}
func byte2HexString(b byte) string {
	return byte2Hex(b)
}
func byte2d() {

}
func getNewWorkFlagsData(packet gopacket.Packet) (branch string, nodes []string) {
	return
}
