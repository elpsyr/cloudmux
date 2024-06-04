package cloudprovider

import (
	"context"

	"yunion.io/x/jsonutils"
)

type ICfelCloudRegion interface {
	ICloudRegion
	GetICfelSkus() ([]ICfelCloudSku, error)
	GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) //  抢占付费 价格
	GetPostPaidPrice(zoneID, instanceType string) (float64, error)     //  按量付费 价格
	GetPrePaidPrice(zoneID, instanceType string) (float64, error)      //  包年包月 价格
	GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) //  抢占付费 售卖状态
	GetPostPaidStatus(zoneID, instanceType string) (string, error)     //  按量付费 售卖状态
	GetPrePaidStatus(zoneID, instanceType string) (string, error)      //  包年包月 售卖状态

	CreateBareMetal(desc *CfelSManagedVMCreateConfig) (ICloudVM, error)
	CreateVM(desc *CfelSManagedVMCreateConfig) (ICloudVM, error)
	DeleteVM(instanceId string) error
	CreateImageByUrl(params *CfelSImageCreateOption) (ICloudImage, error)
	GetImageByID(id string) (ICloudImage, error)
	ResetGuestPassword(params *CfelResetGuestPasswordOption) (ICloudVM, error)
	PingQga(guestId string, timeout int) (bool, error)

	CfelCreateDisk(params *CfelDiskCreateConfig) (ICloudDisk, error)
	CfelAttachDisk(instanceId, diskId string) error
	CfelDetachDisk(instanceId, diskId string) error
	CfelInstanceSettingChange(id string, params *CfelChangeSettingOption) error //虚拟机配置修改
	CfelGetINetworks() ([]ICloudNetwork, error)
	GetIHostsByCondition(*FilterOption) ([]ICloudHost, error)
	MigrateForecast(*MigrateForecastOption) ([]ICfelFilter, error)
	// GetMonitorData(vmId, start, end, interval string) ([]ICfelMonitorData,[]string ,error) // 获取主机监控数据
	GetMonitorDataJSON(*MonitorDataJSONOption) (jsonutils.JSONObject, error) // 获取主机监控数据
	GetGeneralUsage() (ICfelGeneralUsage, error)

	ICfelDeleteImage(id string) error
	GetICfelCloudImage(withUserMeta bool) ([]ICloudImage, error)
	SetImageUserTag(*CfelSetImageUserTag) error

	GetUsableIEip() ([]ICloudEIP, error)

	GetLoadbalancerSkus() ([]ICfelLoadbalancerSku, error)
}

type ICfelCloudSku interface {
	GetZoneID() string
	GetGPUMemorySizeMB() int // GPU 显存
	GetIsBareMetal() bool    // 获取是否为裸金属
	ICloudSku
}

type ICfelZone interface {
	ICloudZone
	GetCapability() (jsonutils.JSONObject, error)
}

// ICfelCloudVM vm接口
type ICfelCloudVM interface {
	RebootVM(ctx context.Context) error
	GetMonitorData(start, end string) ([]ICfelMonitorData, error) // 获取主机监控数据
	GetIsolatedDevice() ([]*IsolatedDeviceInfo, error)
	GetCfelHypervisor() string
	ICloudVM
}

type ICfelMonitorData interface {
	GetBPSRead() float64           // 实例云盘（包括系统盘和数据盘）的读带宽，单位：Byte/s。
	GetInternetTX() float64        // 在查询监控信息时（TimeStamp），实例在指定的间隔时间（Period）内发送的公网数据流量。单位：kbits。
	GetCPU() float64               //实例 vCPU 的使用比例，单位：百分比（%）。
	GetMem() float64               //实例 mem 的使用比例，单位：百分比（%）。
	GetDisk() float64              //实例 disk 的使用比例，单位：百分比（%）。
	GetIOPSWrite() float64         // 实例云盘（包括系统盘和数据盘）的 I/O 写操作，单位：次/s。
	GetIntranetTX() float64        // 在查询监控信息时（TimeStamp），实例在指定的间隔时间（Period）内发送的内网数据流量。单位：kbits。
	GetInstanceId() string         // 实例 ID。
	GetBPSWrite() float64          // 实例云盘（包括系统盘和数据盘）的写带宽，单位：Byte/s。
	GetIOPSRead() float64          // 实例云盘（包括系统盘和数据盘）的 I/O 读操作，单位：次/s。
	GetInternetBandwidth() float64 // 实例的公网带宽，单位时间内的网络流量，单位：kbits/s。
	GetInternetRX() float64        // 在查询监控信息时（TimeStamp），实例在指定的间隔时间（Period）内接收的公网数据流量。单位：kbits。
	GetTimeStamp() string          // 查询监控信息的时间戳。 2014-10-30T05:00:00Z
	GetIntranetRX() float64        // 在查询监控信息时（TimeStamp），实例在指定的间隔时间（Period）内接收的内网数据流量。单位：kbits。
	GetIntranetBandwidth() float64 // 实例的内网带宽，单位时间内的网络流量，单位：kbits/s。
}

// ICfelCloudVpc  vpc 额外功能接口
type ICfelCloudVpc interface {
	Update(opts *VpcUpdateOptions) error // 更新 vpc 名称以及描述
	ICloudVpc
}

type ICfelFilter interface {
	GetFilterName() string
	GetID() string
	GetName() string
	GetReason() []string
}

type ICfelGeneralUsage interface {
	GetAllServers() int
	GetAllServersCpu() int
	GetAllServersMem() int
	GetAllServersDisk() int
	GetHosts() int
	GetHostsCpuTotal() int
	GetBaremetals() int
	GetHostsMem() int
	GetStorages() int
	GetIsolatedDevices() int
	GetAllServersIsolatedDevices() int
	GetRunningServersIsolatedDevices() int
}

type ICfelLoadbalancerBackendGroup interface {
	ICloudLoadbalancerBackendGroup
	CfelAddBackendServer(serverType, serverId, ssl string, weight int, port int) (ICloudLoadbalancerBackend, error)
}

type ICfelLoadbalancerListener interface {
	ICloudLoadbalancerListener
	CfelCreateILoadBalancerListenerRule(*SCfelLoadbalancerListenerRule) (ICloudLoadbalancerListenerRule, error)
	Update(*SLoadbalancerListenerCreateOptions) error
}

type ICfelLoadbalancer interface {
	ICloudLoadbalancer
	CfelCreateILoadBalancerBackendGroup(*SCfelLoadbalancerBackendGroup) (ICloudLoadbalancerBackendGroup, error)
}

type ICfelLoadbalancerSku interface {
	GetName() string
	GetType() string
	GetID() string
}
