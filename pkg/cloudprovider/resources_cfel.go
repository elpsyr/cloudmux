package cloudprovider

import "context"

// resource interface defined by cfel

// ICloudVMUltra vm操作接口进阶操作
// define by cfel
type ICloudVMUltra interface {
	RebootVM(ctx context.Context) error
	GetMonitorData(start, end string) ([]MonitorData, error) // 获取主机监控数据
}

type MonitorData interface {
	GetBPSRead() int           // 实例云盘（包括系统盘和数据盘）的读带宽，单位：Byte/s。
	GetInternetTX() int        // 在查询监控信息时（TimeStamp），实例在指定的间隔时间（Period）内发送的公网数据流量。单位：kbits。
	GetCPU() int               //实例 vCPU 的使用比例，单位：百分比（%）。
	GetIOPSWrite() int         // 实例云盘（包括系统盘和数据盘）的 I/O 写操作，单位：次/s。
	GetIntranetTX() int        // 在查询监控信息时（TimeStamp），实例在指定的间隔时间（Period）内发送的内网数据流量。单位：kbits。
	GetInstanceId() string     // 实例 ID。
	GetBPSWrite() int          // 实例云盘（包括系统盘和数据盘）的写带宽，单位：Byte/s。
	GetIOPSRead() int          // 实例云盘（包括系统盘和数据盘）的 I/O 读操作，单位：次/s。
	GetInternetBandwidth() int // 实例的公网带宽，单位时间内的网络流量，单位：kbits/s。
	GetInternetRX() int        // 在查询监控信息时（TimeStamp），实例在指定的间隔时间（Period）内接收的公网数据流量。单位：kbits。
	GetTimeStamp() string      // 查询监控信息的时间戳。 2014-10-30T05:00:00Z
	GetIntranetRX() int        // 在查询监控信息时（TimeStamp），实例在指定的间隔时间（Period）内接收的内网数据流量。单位：kbits。
	GetIntranetBandwidth() int // 实例的内网带宽，单位时间内的网络流量，单位：kbits/s。
}
