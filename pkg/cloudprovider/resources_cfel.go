package cloudprovider

import "context"

type ICfelCloudRegion interface {
	ICloudRegion
	GetICfelSkus() ([]ICfelCloudSku, error)
	GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) //  抢占付费 价格
	GetPostPaidPrice(zoneID, instanceType string) (float64, error)     //  按量付费 价格
	GetPrePaidPrice(zoneID, instanceType string) (float64, error)      //  包年包月 价格
	GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) //  抢占付费 售卖状态
	GetPostPaidStatus(zoneID, instanceType string) (string, error)     //  按量付费 售卖状态
	GetPrePaidStatus(zoneID, instanceType string) (string, error)      //  包年包月 售卖状态

	CreateBareMetal(desc *SManagedVMCreateConfig) (ICloudVM, error)
	CreateVM(desc *SManagedVMCreateConfig) (ICloudVM, error)
	CreateImageByUrl(params *CfelSImageCreateOption) (ICloudImage, error)
	GetImageByID(id string) (ICloudImage, error)
	ResetGuestPassword(params *CfelResetGuestPasswordOption) (ICloudVM, error)
	PingQga(guestId string, timeout int) (bool, error)

	CfelCreateDisk(params *CfelDiskCreateConfig) (ICloudDisk, error)
	CfelAttachDisk(instanceId,diskId string) error
	CfelDetachDisk(instanceId,diskId string) error
	CfelInstanceSettingChange(id string,params *CfelChangeSettingOption) error //虚拟机配置修改
}

type ICfelCloudSku interface {
	GetZoneID() string
	GetGPUMemorySizeMB() int // GPU 显存
	GetIsBareMetal() bool    // 获取是否为裸金属
	ICloudSku
}

// ICfelCloudVM vm接口
type ICfelCloudVM interface {
	RebootVM(ctx context.Context) error
	GetMonitorData(start, end string) ([]ICfelMonitorData, error) // 获取主机监控数据
	ICloudVM
}

type ICfelMonitorData interface {
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

// ICfelCloudVpc  vpc 额外功能接口
type ICfelCloudVpc interface {
	Update(opts *VpcUpdateOptions) error // 更新 vpc 名称以及描述
	ICloudVpc
}
