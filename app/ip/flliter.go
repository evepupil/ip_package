package ip

var portChan chan int
var ipChan chan string

func flliterBySourceIp(ip string, PkgRows []PkgRow) {
	var res []PkgRow
	for _, v := range PkgRows {
		if ip == "" || v.Source == ip {
			res = append(res, v)
		}
	}
	Pkgs = res
	ReLoadPkgList(res)
}
func flliterByDestIp(ip string, PkgRows []PkgRow) {
	var res []PkgRow
	for _, v := range PkgRows {
		if ip == "" || v.Dest == ip {
			res = append(res, v)
		}
	}
	Pkgs = res
	ReLoadPkgList(res)
}
func flliterBySourcePort(port string, PkgInfos []PkgRow) {
	var res []PkgRow
	for _, v := range PkgInfos {
		if port == "" || get_source_port(v.Info) == port {
			res = append(res, v)
		}
	}
	Pkgs = res
	ReLoadPkgList(res)
}
func flliterByDestPort(port string, PkgInfos []PkgRow) {
	var res []PkgRow
	for _, v := range PkgInfos {
		if port == "" || get_dest_port(v.Info) == port {
			res = append(res, v)
		}
	}
	Pkgs = res
	ReLoadPkgList(res)
}
func get_source_port(str string) string {
	res := ""
	for i := 0; i < len(str); i++ {
		if str[i] == '-' {
			res = str[0:i]
		}
	}
	return res
}
func get_dest_port(str string) string {
	res := ""
	for i := 0; i < len(str); i++ {
		if str[i] == '>' {
			res = str[i+1 : len(str)]
		}
	}
	return res
}
