# IP_pkg_analyze

基于 Go 语言和 fyne GUI 框架实现的跨平台网络抓包工具。

## 主要功能
- 支持多平台（Windows/Linux）图形化抓包
- 网络设备选择与抓包
- 支持混杂/严格模式切换
- 实时显示数据包列表，支持过滤与排序
- 数据包详细内容和分层信息展示
- 支持保存/打开 pcap 文件
- 支持自定义数据包发送
- 实时流量速率监控

## 依赖环境
- Go 1.17 及以上
- fyne.io/fyne/v2
- github.com/google/gopacket
- github.com/flopp/go-findfont

## 安装与运行
1. 安装依赖：
   ```bash
   go mod tidy
   ```
2. 运行程序：
   ```bash
   go run main.go
   ```

## 目录结构简介
- `main.go`         程序入口
- `app/ip/`         核心抓包、界面与功能实现
- `app/util/`       字体等工具
- `func/`           设备选择等辅助功能

## 注意事项
- Windows 下需保证有合适的字体（如 simhei.ttf），程序会自动查找。
- 抓包需管理员权限。

---
如有问题欢迎反馈。
