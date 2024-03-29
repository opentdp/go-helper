package psutil

import (
	"encoding/json"
)

// Go 内存信息

type GoMemoryStat struct {
	Alloc      uint64 `note:"已分配内存"`
	Sys        uint64 `note:"已申请内存"`
	HeapAlloc  uint64 `note:"堆已分配内存"`
	HeapInuse  uint64 `note:"堆已使用内存"`
	HeapSys    uint64 `note:"堆已申请内存"`
	StackInuse uint64 `note:"栈已使用内存"`
	StackSys   uint64 `note:"栈已申请内存"`
	TotalAlloc uint64 `note:"累计已分配内存"`
	LastGC     uint64 `note:"最后一次 GC 时间"`
	NumGC      uint32 `note:"GC 执行次数"`
}

// 系统概要信息

type SummaryStat struct {
	CreateAt     int64     `note:"创建时间"`
	HostId       string    `note:"主机 ID"`
	HostName     string    `note:"主机名"`
	Uptime       uint64    `note:"运行时间"`
	OS           string    `note:"操作系统"`
	Platform     string    `note:"平台"`
	KernelArch   string    `note:"内核架构"`
	CpuCore      int       `note:"CPU 核心数"`
	CpuCoreLogic int       `note:"CPU 逻辑核心数"`
	CpuPercent   []float64 `note:"CPU 使用率"`
	MemoryTotal  uint64    `note:"内存总量"`
	MemoryUsed   uint64    `note:"内存使用量"`
	PublicIpv4   string    `note:"公网 IPV4"`
	PublicIpv6   string    `note:"公网 IPV6"`
}

func (p *SummaryStat) From(s string) {
	json.Unmarshal([]byte(s), p)
}

func (p *SummaryStat) String() string {
	jsonbyte, _ := json.Marshal(p)
	return string(jsonbyte)
}

// 系统统计详情

type DetailStat struct {
	*SummaryStat
	CpuModel      []string        `note:"CPU 型号"`
	NetInterface  []NetInterface  `note:"网卡信息"`
	NetBytesRecv  uint64          `note:"网卡接收字节数"`
	NetBytesSent  uint64          `note:"网卡发送字节数"`
	DiskPartition []DiskPartition `note:"磁盘分区信息"`
	DiskTotal     uint64          `note:"磁盘总量"`
	DiskUsed      uint64          `note:"磁盘使用量"`
	SwapTotal     uint64          `note:"交换分区总量"`
	SwapUsed      uint64          `note:"交换分区使用量"`
}

// 硬盘分区信息

type DiskPartition struct {
	Device     string `note:"设备名"`
	Mountpoint string `note:"挂载点"`
	Fstype     string `note:"文件系统"`
	Total      uint64 `note:"总量"`
	Used       uint64 `note:"使用量"`
}

// 网卡信息

type NetInterface struct {
	Name      string   `note:"网卡名称"`
	BytesRecv uint64   `note:"接收字节数"`
	BytesSent uint64   `note:"发送字节数"`
	Dropin    uint64   `note:"丢弃的接收包"`
	Dropout   uint64   `note:"丢弃的发送包"`
	Ipv4List  []string `note:"IPV4 列表"`
	Ipv6List  []string `note:"IPV6 列表"`
}
