package ip

import (
	"fyne.io/fyne/v2"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
)

func OpenPcap(pcapFile string) {
	// Open file instead of device
	InitData()
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.Packets()
	No := 1
	for packet := range packetSource.Packets() {
		//sourceIpfilter = <-util.ipChan
		if No == 1 {
			//StartTime=
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
			Pkgs = append(Pkgs, p)
			PkgInfos = append(PkgInfos, packet)
			PkgStringList = append(PkgStringList, p.FormatePkgListInfo())
			if SortBySourceIp == INCRESE {
				ReLoadPkgList(PkgSortBySourceIp(Pkgs))
			}
		}
		No++
	}

}
func SaveAsPcap(path string, PkgInfos []gopacket.Packet) (string, error) {
	//length:=len(PkgInfos)
	// Open output pcap file and write header
	//FName := time.Now().Format("2006-01-02 15.04.05")
	//FName:="1"
	f, _ := os.Create(path)
	w := pcapgo.NewWriter(f)
	err := w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
	if err != nil {
		fyne.LogError("写入文件头失败", err)
	}
	defer f.Close()
	// Open the device for capturing
	// defer handle.Close()  //会阻塞
	// Start processing packets
	for _, packet := range PkgInfos {
		// Process packet here
		err = w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			fyne.LogError("将包写入文件失败", err)
		}
	}
	return path, nil

}

//func SaveAsPcap(PkgInfos []gopacket.Packet) {   //version:1.1
//	// Open output pcap file and write header
//	FName:=time.Now().Format("2006-01-02 15.04.05")
//	f, _ := os.Create(FName+".pcap")
//	//f, _ := os.Create(FName)
//	w := pcapgo.NewWriter(f)
//	err:=w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
//	if err!=nil{
//		fmt.Println("出错了1")
//	}
//	defer f.Close()
//	for  _,packet:=range PkgInfos {
//		// Process packet here
//		//fmt.Println(packet)
//		err=w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
//		packetCount++
//		// Only capture 100 and then stop
//	}
//	if err!=nil{
//		fmt.Println("出错了2")
//	}
//}
