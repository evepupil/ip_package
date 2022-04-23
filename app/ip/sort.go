package ip

import (
	"sort"
	"strconv"
)

type rowsBySourceIp []PkgRow
type rowsByDstIp []PkgRow
type rowsByLength []PkgRow
type rowsBySourcePort []PkgRow
type rowsByDestPort []PkgRow

func PkgSortBySourceIp(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsBySourceIp(PkgInfos))
	return PkgInfos
}
func PkgSortBySourceIpReverse(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsBySourceIp(PkgInfos))
	sort.Sort(sort.Reverse(rowsBySourceIp(PkgInfos)))
	return PkgInfos
}
func PkgSortByDestIp(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsByDstIp(PkgInfos))
	return PkgInfos
}
func PkgSortByDestIpReverse(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsByDstIp(PkgInfos))
	sort.Sort(sort.Reverse(rowsByDstIp(PkgInfos)))
	return PkgInfos
}

func PkgSortByLength(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsByLength(PkgInfos))
	return PkgInfos
}
func PkgSortByLengthReverse(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsByLength(PkgInfos))
	sort.Sort(sort.Reverse(rowsByLength(PkgInfos)))
	return PkgInfos
}

func PkgSortBySourcePort(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsBySourcePort(PkgInfos))
	return PkgInfos
}
func PkgSortBySourcePortReverse(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsBySourcePort(PkgInfos))
	sort.Sort(sort.Reverse(rowsBySourcePort(PkgInfos)))
	return PkgInfos
}
func PkgSortByDestPort(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsByDestPort(PkgInfos))
	return PkgInfos
}
func PkgSortByDestPortReverse(PkgInfos []PkgRow) []PkgRow {
	sort.Sort(rowsByDestPort(PkgInfos))
	sort.Sort(sort.Reverse(rowsByDestPort(PkgInfos)))
	return PkgInfos
}
func (receiver rowsBySourceIp) Len() int {
	return len(receiver)
}
func (receiver rowsBySourceIp) Swap(i, j int) {
	receiver[i], receiver[j] = receiver[j], receiver[i]
}
func (receiver rowsBySourceIp) Less(i, j int) bool {
	return receiver[i].Source < receiver[j].Source
}

func (receiver rowsByDstIp) Len() int {
	return len(receiver)
}
func (receiver rowsByDstIp) Swap(i, j int) {
	receiver[i], receiver[j] = receiver[j], receiver[i]
}
func (receiver rowsByDstIp) Less(i, j int) bool {
	return receiver[i].Dest < receiver[j].Dest
}

func (receiver rowsByLength) Len() int {
	return len(receiver)
}
func (receiver rowsByLength) Swap(i, j int) {
	receiver[i], receiver[j] = receiver[j], receiver[i]
}
func (receiver rowsByLength) Less(i, j int) bool {
	return receiver[i].Length < receiver[j].Length
}

func (receiver rowsBySourcePort) Len() int {
	return len(receiver)
}
func (receiver rowsBySourcePort) Swap(i, j int) {
	receiver[i], receiver[j] = receiver[j], receiver[i]
}
func (receiver rowsBySourcePort) Less(i, j int) bool {
	a, _ := strconv.Atoi(get_source_port(receiver[i].Info))
	b, _ := strconv.Atoi(get_source_port(receiver[j].Info))
	return a < b
}
func (receiver rowsByDestPort) Len() int {
	return len(receiver)
}
func (receiver rowsByDestPort) Swap(i, j int) {
	receiver[i], receiver[j] = receiver[j], receiver[i]
}
func (receiver rowsByDestPort) Less(i, j int) bool {
	a, _ := strconv.Atoi(get_dest_port(receiver[i].Info))
	b, _ := strconv.Atoi(get_dest_port(receiver[j].Info))
	return a < b
}
