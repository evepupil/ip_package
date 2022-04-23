package ip

import "github.com/google/gopacket"

//1 源ip 2 目的ip
func Refind(slice []string, dest_string string, types int) int {
	//i := len(slice) - 1
	//for ; i >= 0; i-- {
	//	if types == 1 && dest_string == slice[i].Source {
	//		return i
	//	}
	//	if types == 2 && dest_string == slice[i].Dest {
	//		return i
	//	}
	//}
	return -1
}
func InsertStr(slice *[]string, index int, value string) {
	s := *slice
	s = append(s, value)
	copy(s[index+1:], s[index:])
	s[index] = value
}
func InsertPkgRow(slice *[]PkgRow, index int, value PkgRow) {
	s := *slice
	s = append(s, value)
	copy(s[index+1:], s[index:])
	s[index] = value
}
func InsertPkgInfo(slice *[]gopacket.Packet, index int, value gopacket.Packet) {
	s := *slice
	s = append(s, value)
	copy(s[index+1:], s[index:])
	s[index] = value
}

// 1 源ip 2 目的ip
// 二分查找插入
func GetIndexInsert(slice []string, string2 string, types int) int {
	l, r := 0, len(slice)
	for l < r {
		mid := (l + r) / 2
		if get_info_from_str(slice[mid], types) > get_info_from_str(string2, types) {
			r = mid - 1
		} else if get_info_from_str(slice[mid], types) < get_info_from_str(string2, types) {
			l = mid
		} else {
			for get_info_from_str(slice[mid], types) == get_info_from_str(string2, types) {
				mid++
			}
			return mid
		}
	}
	return r

}

// 1 源ip 2 目的ip
func get_info_from_str(str string, types int) string {
	k := 0
	index := 0
	res := ""
	for i := 0; i < len(str); i++ {
		if str[i] == ' ' {
			k++
			for str[i] == ' ' {
				i++
			}
			i--
		}
		if types == 1 {
			if k == 2 {
				index = i
				for str[i] != ' ' {
					i++
				}
				res = str[index:i]
			}
		} else if types == 2 {
			if k == 3 {
				index = i
				for str[i] != ' ' {
					i++
				}
				res = str[index:i]
			}
		}
	}
	return res
}
