package ip

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"strconv"
	"time"
)

var (
	PkgInfos      []gopacket.Packet
	StartTime     time.Time
	Pkgs          []PkgRow
	AllPkgInfos   []gopacket.Packet
	AllPkgs       []PkgRow
	Map_Pkg_Infos map[PkgRow]gopacket.Packet
)

type PkgRow struct {
	No       int
	Time     time.Time
	Source   string
	Dest     string
	Protocol string
	Length   int
	Info     string
}

var (
	downStreamDataSize = 0 // 单位时间内下行的总字节数
	upStreamDataSize   = 0 // 单位时间内上行的总字节数
	deviceName         = flag.String("i", "eth0", "network interface device name")
)
var (
	SourceIp_filter   = ""
	DestIp_filter     = ""
	SourcePort_filter = ""
	DestPort_filter   = ""
)

//设备名：pcap.FindAllDevs()返回的设备的Name
//snaplen：捕获一个数据包的多少个字节，一般来说对任何情况65535是一个好的实践，如果不关注全部内容，只关注数据包头，可以设置成1024
//promisc：设置网卡是否工作在混杂模式，即是否接收目的地址不为本机的包
//timeout：设置抓到包返回的超时。如果设置成30s，那么每30s才会刷新一次数据包；设置成负数，会立刻刷新数据包，即不做等待
//要记得释放掉handle

var (
	device       string = "eth0"
	snapshot_len int32  = 10240
	snapshotLen  int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = -1 * time.Second
	handle       *pcap.Handle
	packetCount  int = 0
)
var (
	SortBySourceIp   = 0
	SortByDestIp     = 0
	SortBySourcePort = 0
	SortByDestPort   = 0
	SortByLength     = 0
)

const (
	//0默认 1递增 2递减
	INCRESE = 1
	DECRESE = 2
)

var IsGetPkg bool = false

//混杂模式
var Promiscuous bool = false
var (
	Sort_tpyes = 0
)

func GetPkg(ctx context.Context, device_str string) {
	InitData()
	//dstIpfilter:=""
	No := 1
	// Open device
	device = device_str
	handle, err = pcap.OpenLive(device_str, snapshot_len, Promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	//查找所有设备
	devices, err := pcap.FindAllDevs()
	var device pcap.Interface
	for _, d := range devices {
		if d.Name == device_str {
			device = d
		}
	}
	// 根据网卡的ipv4地址获取网卡的mac地址，用于后面判断数据包的方向
	macAddr, err := findMacAddrByIp(findDeviceIpv4(device))
	if err != nil {
		panic(err)
	}
	//启动计时器
	ctx1, cancel1 := context.WithCancel(context.Background())
	go monitor(ctx1)
	defer cancel1()
	for packet := range packetSource.Packets() {
		//sourceIpfilter = <-util.ipChan
		if No == 10 {
			time.Sleep(5 * time.Second)
		}
		//packet:=<-packetSource.Packets()
		p := anlysePacket(packet)
		p.No = No
		AllPkgs = append(AllPkgs, p)
		AllPkgInfos = append(AllPkgInfos, packet)
		Map_Pkg_Infos[p] = packet
		if (SourceIp_filter == "" || p.Source == SourceIp_filter) &&
			(DestIp_filter == "" || p.Dest == DestIp_filter) &&
			(SourcePort_filter == "" || get_source_port(p.Info) == SourcePort_filter) &&
			(DestPort_filter == "" || get_dest_port(p.Info) == DestPort_filter) {
			if Sort_tpyes != 0 {
				index := GetIndexInsert(PkgStringList, p.FormatePkgListInfo(), Sort_tpyes)
				InsertPkgRow(&Pkgs, index, p)
				InsertPkgInfo(&PkgInfos, index, packet)
				InsertStr(&PkgStringList, index, p.FormatePkgListInfo())
			} else {
				Pkgs = append(Pkgs, p)
				PkgInfos = append(PkgInfos, packet)
				PkgStringList = append(PkgStringList, p.FormatePkgListInfo())
			}
			if SortBySourceIp == INCRESE {
				ReLoadPkgList(PkgSortBySourceIp(Pkgs))
			} else if SortBySourceIp == DECRESE {
				ReLoadPkgList(PkgSortBySourceIpReverse(Pkgs))
			} else if SortByDestIp == INCRESE {
				ReLoadPkgList(PkgSortByDestIp(Pkgs))
			} else if SortByDestIp == DECRESE {
				ReLoadPkgList(PkgSortByDestIpReverse(Pkgs))
			} else if SortByLength == INCRESE {
				ReLoadPkgList(PkgSortByLength(Pkgs))
			} else if SortByLength == DECRESE {
				ReLoadPkgList(PkgSortByLengthReverse(Pkgs))
			} else if SortBySourcePort == INCRESE {
				ReLoadPkgList(PkgSortBySourcePort(Pkgs))
			} else if SortBySourcePort == DECRESE {
				ReLoadPkgList(PkgSortBySourcePortReverse(Pkgs))
			} else if SortByDestPort == INCRESE {
				ReLoadPkgList(PkgSortByDestPort(Pkgs))
			} else if SortByDestPort == DECRESE {
				ReLoadPkgList(PkgSortByDestPortReverse(Pkgs))
			}
			//统计流量
			ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
			if ethernetLayer != nil {
				ethernet := ethernetLayer.(*layers.Ethernet)
				// 如果封包的目的MAC是本机则表示是下行的数据包，否则为上行
				if ethernet.DstMAC.String() == macAddr {
					downStreamDataSize += len(packet.Data()) // 统计下行封包总大小
				} else {
					upStreamDataSize += len(packet.Data()) // 统计上行封包总大小
				}
			}
			//fmt.Println(packet.Layers()[len(packet.Layers())-1].LayerContents())
		}
		No++
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}

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
func anlysePacket(p gopacket.Packet) PkgRow {
	var nilEndpoint gopacket.Endpoint = gopacket.Endpoint{}
	var nilFlow gopacket.Flow = gopacket.Flow{}
	pkgrow := PkgRow{Time: p.Metadata().Timestamp,
		Length: p.Metadata().Length,
	}
	if p.NetworkLayer() != nil { //网际层
		if p.NetworkLayer().NetworkFlow().Src() != nilEndpoint {
			pkgrow.Source = p.NetworkLayer().NetworkFlow().Src().String()
		}
		if p.NetworkLayer().NetworkFlow().Dst() != nilEndpoint {
			pkgrow.Dest = p.NetworkLayer().NetworkFlow().Dst().String()
		}
		if p.NetworkLayer().NetworkFlow() != nilFlow {
			pkgrow.Protocol = p.NetworkLayer().NetworkFlow().EndpointType().String()
		}
	}
	if p.TransportLayer() != nil { //传输层
		if p.TransportLayer().TransportFlow() != nilFlow {
			//s:=gopacket.LayerString(p.TransportLayer())
			//ipLayer := p.Layer(layers.LayerTypeIPv4)
			//if ipLayer != nil {
			//	fmt.Println("IPv4 layer detected.")
			//	ip, _ := ipLayer.(*layers.IPv4)
			//
			//	// IP layer variables:
			//	// Version (Either 4 or 6)
			//	// IHL (IP Header Length in 32-bit words)
			//	// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
			//	// Checksum, SrcIP, DstIP
			//	fmt.Println(ip.Checksum)
			//	fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			//	fmt.Println("Protocol: ", ip.Protocol)
			//}
			//tcpLayer := p.Layer(layers.LayerTypeTCP)
			//if tcpLayer != nil {
			//	fmt.Println("TCP layer detected.")
			//	tcp, _ := tcpLayer.(*layers.TCP)
			//	// TCP layer variables:
			//	// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
			//	// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
			//	fmt.Println(tcp.Options)
			//	fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
			//	fmt.Println("Sequence number: ", tcp.Seq)
			//	fmt.Println()
			//}
			pkgrow.Protocol = p.TransportLayer().TransportFlow().EndpointType().String()
			pkgrow.Info = p.TransportLayer().TransportFlow().String()
		}
		if p.ApplicationLayer() != nil {
			switch p.ApplicationLayer().LayerType() {
			case layers.LayerTypeTLS:
				pkgrow.Protocol = "TLS"
			case layers.LayerTypeDNS:
				pkgrow.Protocol = "DNS"
			}

		}
	}
	return pkgrow

}
func (p PkgRow) FormatePkgListInfo() string {
	res := ""
	res += strconv.Itoa(p.No)
	t := p.Time.Format("\"2006-01-02T15:04:05\"")
	res += blankAdd(15-len(res)) + t[1:len(t)-1]
	res += blankAdd(45-len(res)) + p.Source
	res += blankAdd(70-len(res)) + p.Dest
	res += blankAdd(98-len(res)) + p.Protocol
	res += blankAdd(115-len(res)) + strconv.Itoa(p.Length)
	res += blankAdd(129-len(res)) + p.Info
	return res
}
func blankAdd(n int) string {
	res := ""
	for n > 0 {
		n--
		res += " "
	}
	return res

}

// 获取网卡的IPv4地址
func findDeviceIpv4(device pcap.Interface) string {
	for _, addr := range device.Addresses {
		if ipv4 := addr.IP.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	panic("device has no IPv4")
}

// 根据网卡的IPv4地址获取MAC地址
// 有此方法是因为gopacket内部未封装获取MAC地址的方法，所以这里通过找到IPv4地址相同的网卡来寻找MAC地址
func findMacAddrByIp(ip string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(interfaces)
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			if a, ok := addr.(*net.IPNet); ok {
				if ip == a.IP.String() {
					return i.HardwareAddr.String(), nil
				}
			}
		}
	}
	return "", errors.New(fmt.Sprintf("no device has given ip: %s", ip))
}

// 每一秒计算一次该秒内的数据包大小平均值，并将下载、上传总量置零
func monitor(ctx context.Context) {
	for {
		FlowsStr.Set(fmt.Sprintf("\rDown:%.2fkb/s \t Up:%.2fkb/s                                                                                                                             ", float32(downStreamDataSize)/1024/1, float32(upStreamDataSize)/1024/1))
		downStreamDataSize = 0
		upStreamDataSize = 0
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
func InitData() {
	fmt.Println("初始化数据成功")
	Map_Pkg_Infos = make(map[PkgRow]gopacket.Packet)
	AllPkgs = []PkgRow{}
	AllPkgInfos = []gopacket.Packet{}
	Pkgs = []PkgRow{}
	PkgInfos = []gopacket.Packet{}
	PkgStringList = []string{}
}
