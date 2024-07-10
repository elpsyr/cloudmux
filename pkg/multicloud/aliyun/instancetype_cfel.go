package aliyun

import (
	"strconv"
	"strings"
	"unicode"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
)

// SInstanceTypeCFEL
// https://api.aliyun.com/api/Ecs/2014-05-26/DescribeInstanceTypes
type SInstanceTypeCFEL struct {
	multicloud.SInstanceBase
	AliyunTags
	ZoneID                      string  // zone
	CpuArchitecture             string  // CPU架构，可能值： X86。 ARM。
	CpuSpeedFrequency           float64 // CPU基频，单位GHz。
	CpuTurboFrequency           float64 // CPU睿频，单位GHz。
	DiskQuantity                int     // 支持挂载的云盘数量上限。
	EniIpv6AddressQuantity      int     // 单块弹性网卡的IPv6地址上限。
	MaximumQueueNumberPerEni    int     // 单块弹性网卡最大队列数。包括主网卡及辅助网卡支持的队列数。
	EniPrivateIpAddressQuantity int     // 单块弹性网卡的IPv4地址上限。
	EniTotalQuantity            int     // 支持挂载的所有网卡（包括主网卡、弹性网卡、中继网卡等）上限。
	EniTrunkSupported           bool    // 实例规格所挂载的网卡是否支持中继。
	EriQuantity                 int     // 弹性RDMA网卡（ERI）数量。
	InstanceCategory            string  // 实例规格分类。
	NetworkEncryptionSupport    bool    // 实例是否支持VPC网络流量加密
	NvmeSupport                 string  // 实例规格所挂载的云盘是否支持NVMe。可能值：required：支持。表示云盘以NVMe的方式挂载。 unsupported：不支持。表示云盘不以NVMe的方式挂载。
	PhysicalProcessorModel      string  // 处理器型号。示例值: Intel Xeon(Ice Lake) Platinum 8369B
	PrimaryEniQueueNumber       int     // 主网卡默认队列数。
	SecondaryEniQueueNumber     int     // 辅助弹性网卡默认队列数。
	TotalEniQueueQuantity       int     // 实例规格允许修改的弹性网卡队列数总配额。
	GPUMemorySize               float64 // 规格对应的单块GPU显存。单位：GiB。
	QueuePairNumber             int     // 单块弹性RDMA网卡（ERI）的QP（QueuePair）队列数上限。
	InitialCredit               int     // 突发性能实例t5、t6的初始vCPU积分值。
	NetworkCardQuantity         int     // 实例规格支持的物理网卡数量。
	PostpaidStatus              string
	PrepaidStatus               string
}

// GetICfelSkus 获取 aliyun ICfelCloudSku
func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	skus, err := self.GetRegionAvailableInstanceTypes()
	if err != nil {
		return nil, errors.Wrapf(err, "GetInstanceTypes")
	}
	var ret []cloudprovider.ICfelCloudSku
	for i := range skus {
		ret = append(ret, &skus[i])
	}
	return ret, nil
}

const (
	AliyunResourceAvailable = "Available"
	AliyunResourceSoldOut   = "SoldOut"
)

// SAvailableResource 某 Region 下的可用资源
// 用于反序列化 DescribeAvailableResource 接口的返回数据
// https://api.aliyun.com/api/Ecs/2014-05-26/DescribeAvailableResource
type SAvailableResource struct {
	AvailableZones struct {
		AvailableZone []struct {
			Status             string `json:"Status"`
			StatusCategory     string `json:"StatusCategory"`
			ZoneID             string `json:"ZoneId"`
			AvailableResources struct {
				AvailableResource []struct {
					Type               string `json:"Type"`
					SupportedResources struct {
						SupportedResource []struct {
							Status         string `json:"Status"`
							StatusCategory string `json:"StatusCategory"`
							Value          string `json:"Value"`
						} `json:"SupportedResource"`
					} `json:"SupportedResources"`
				} `json:"AvailableResource"`
			} `json:"AvailableResources"`
			RegionID string `json:"RegionId"`
		} `json:"AvailableZone"`
	} `json:"AvailableZones"`
}

// GetZoneID
// 获取 SInstanceType 所属 zone
func (self *SInstanceType) GetZoneID() string {
	return self.ZoneID
}

// GetRegionAvailableInstanceTypes 获取 region 下可以获取的 instanceType
// SInstanceType (server sku) 所属层级 ： region->zone->sku
// implement by cfel
func (self *SRegion) GetRegionAvailableInstanceTypes() ([]SInstanceType, error) {
	// 0. 获取 region 层级下所有 InstanceTypes
	params := make(map[string]string)
	params["RegionId"] = self.RegionId

	body, err := self.ecsRequest("DescribeInstanceTypes", params)
	if err != nil {
		log.Errorf("GetInstanceTypes fail %s", err)
		return nil, err
	}

	instanceTypes := make([]SInstanceType, 0)
	err = body.Unmarshal(&instanceTypes, "InstanceTypes", "InstanceType")
	if err != nil {
		log.Errorf("Unmarshal instance type details fail %s", err)
		return nil, err
	}
	// 0. instanceTypes列表 转存map[SInstanceType.InstanceTypeId] SInstanceType
	instanceTypeMap := make(map[string]SInstanceType, 0)
	for _, instanceType := range instanceTypes {
		instanceTypeMap[instanceType.InstanceTypeId] = instanceType
	}
	// 1. 获取 region->zone 层级下 可用资源 DescribeAvailableResource

	// 1.1 PostPaid
	postPaidAvailableResource, err := self.GetAvailableResource(DestinationResourceInstanceType, POSTPAID, "", "", false)
	if err != nil {
		log.Errorf("GetAvailableResource PostPaid fail %s", err)
		return nil, err
	}
	// 1.2 PrePaid
	prePaidAvailableResource, err := self.GetAvailableResource(DestinationResourceInstanceType, PREPAID, "", "", false)
	if err != nil {
		log.Errorf("GetAvailableResource PrePaid fail %s", err)
		return nil, err
	}
	prePaidAvailableResourceMap := make(map[string]map[string]string, 0)
	for _, zone := range prePaidAvailableResource.AvailableZones.AvailableZone {
		zoneAvailable := make(map[string]string, 0)
		for _, availableResource := range zone.AvailableResources.AvailableResource {
			for _, sku := range availableResource.SupportedResources.SupportedResource {
				zoneAvailable[sku.Value] = sku.Status
			}
		}
		prePaidAvailableResourceMap[zone.ZoneID] = zoneAvailable
	}
	// 2. 处理获得 关联zone的 []SInstanceType
	// 为 availableResource zone下的 InstanceTypes 创建 SInstanceType
	zonesInstanceType := make([]SInstanceType, 0)
	for _, zone := range postPaidAvailableResource.AvailableZones.AvailableZone {
		//zoneById, err := self.GetIZoneById(zone.ZoneID)
		if err != nil {
			log.Errorf("GetIZoneById fail %s", err)
			return nil, err
		}
		for _, resources := range zone.AvailableResources.AvailableResource {
			for _, resource := range resources.SupportedResources.SupportedResource {
				_instanceType, ok := instanceTypeMap[resource.Value]
				if ok {
					_instanceType.ZoneID = zone.ZoneID
					_instanceType.PostpaidStatus = resource.Status

					// set PrepaidStatus
					prePaidAvailableStatus, ok := prePaidAvailableResourceMap[zone.ZoneID][resource.Value]
					if ok {
						_instanceType.PrepaidStatus = prePaidAvailableStatus
					}
					zonesInstanceType = append(zonesInstanceType, _instanceType)
				}
			}

		}
	}
	return zonesInstanceType, nil
}

const (
	DestinationResourceSystemDisk   = "SystemDisk"
	DestinationResourceInstanceType = "InstanceType"
)

// GetAvailableResource 查询某一可用区（zone）的资源列表。
// AvailableResource 所属层级 ： region -> zone-> AvailableResource
// DestinationResource ：SystemDisk、InstanceType
func (self *SRegion) GetAvailableResource(DestinationResource, InstanceChargeType string, ZoneID, InstanceType string, spot bool) (*SAvailableResource, error) {
	// 0. 获取 region 层级下所有 InstanceTypes
	params := make(map[string]string)
	params["RegionId"] = self.RegionId
	params["DestinationResource"] = DestinationResource // 要查询的资源类型
	params["InstanceChargeType"] = InstanceChargeType   //资源的计费方式
	if InstanceType != "" {
		params["InstanceType"] = InstanceType
	}
	if ZoneID != "" {
		params["ZoneId"] = ZoneID
	}
	if spot {
		params["InstanceChargeType"] = POSTPAID
		params["SpotStrategy"] = SPOTPOSTPAID
	}

	body, err := self.ecsRequest("DescribeAvailableResource", params)
	if err != nil {
		log.Errorf("GetAvailableResource fail %s", err)
		return nil, err
	}

	availableResources := new(SAvailableResource)
	err = body.Unmarshal(&availableResources)
	if err != nil {
		log.Errorf("Unmarshal available resources details fail %s", err)
		return nil, err
	}
	return availableResources, nil
}

// GetInstanceTypeAvailableDiskType 获取实例规格可用的系统盘类型
func (self *SRegion) GetInstanceTypeAvailableDiskType(InstanceChargeType string, ZoneID, InstanceType string, spot bool) ([]string, error) {
	supportType := make([]string, 0)
	resource, err := self.GetAvailableResource(DestinationResourceSystemDisk, InstanceChargeType, ZoneID, InstanceType, spot)
	if err != nil {
		return supportType, err
	}
	availableZone := resource.AvailableZones.AvailableZone
	if availableZone != nil && len(availableZone) > 0 {
		availableResource := availableZone[0].AvailableResources.AvailableResource
		if availableResource != nil && len(availableResource) > 0 {
			supportedResource := availableResource[0].SupportedResources.SupportedResource
			if supportedResource != nil && len(supportedResource) > 0 {
				for _, sResource := range supportedResource {
					if sResource.Status == AliyunResourceAvailable {
						supportType = append(supportType, sResource.Value)
					}
				}
			}
		}
	}
	return supportType, nil

}

func (self *SInstanceType) GetStatus() string {
	return ""
}

func (self *SInstanceType) Delete() error {
	return nil
}

func (self *SInstanceType) GetName() string {
	return self.InstanceTypeId
}

func (self *SInstanceType) GetId() string {
	return self.InstanceTypeId
}

func (self *SInstanceType) GetGlobalId() string {
	return self.InstanceTypeId
}

func (self *SInstanceType) GetInstanceTypeFamily() string {
	return self.InstanceTypeFamily
}

func (self *SInstanceType) GetInstanceTypeCategory() string {
	return self.InstanceCategory
}

func (self *SInstanceType) GetPrepaidStatus() string {
	if self.PrepaidStatus == AliyunResourceAvailable {
		return api.SkuStatusAvailable
	} else if self.PrepaidStatus == AliyunResourceSoldOut {
		return api.SkuStatusSoldout
	}
	return api.SkuStatusSoldout
}

func (self *SInstanceType) GetPostpaidStatus() string {
	if self.PostpaidStatus == AliyunResourceAvailable {
		return api.SkuStatusAvailable
	} else if self.PostpaidStatus == AliyunResourceSoldOut {
		return api.SkuStatusSoldout
	}
	return api.SkuStatusSoldout
}

func (self *SInstanceType) GetCpuArch() string {
	return self.CpuArchitecture
}

func (self *SInstanceType) GetCpuCoreCount() int {
	return int(self.CpuCoreCount)
}

func (self *SInstanceType) GetMemorySizeMB() int {
	return int(self.MemorySize * 1024)
}

func (self *SInstanceType) GetOsName() string {
	return "Any"
}

func (self *SInstanceType) GetSysDiskResizable() bool {
	return true
}

func (self *SInstanceType) GetSysDiskType() string {
	return ""
}

func (self *SInstanceType) GetSysDiskMinSizeGB() int {
	return 0
}

func (self *SInstanceType) GetSysDiskMaxSizeGB() int {
	return 0
}

// 本地盘  https://help.aliyun.com/zh/ecs/user-guide/local-disks
// openapi  https://api.aliyun.com/api/Ecs/2014-05-26/DescribeInstanceTypes

// GetAttachedDiskType 本地盘类型
// local_hdd_pro：实例规格族 d1ne 和 d1 搭载的 SATA HDD 本地盘。
// local_ssd_pro：实例规格族 i2、i2g、i1、ga1 和 gn5 等搭载的 NVMe SSD 本地盘。
func (self *SInstanceType) GetAttachedDiskType() string {
	return self.LocalStorageCategory
}

// GetAttachedDiskSizeGB 实例挂载的本地盘的单盘容量
func (self *SInstanceType) GetAttachedDiskSizeGB() int {
	return int(self.LocalStorageCapacity)
}

// GetAttachedDiskCount 实例挂载的本地盘的数量。
func (self *SInstanceType) GetAttachedDiskCount() int {
	return self.LocalStorageAmount
}

func (self *SInstanceType) GetDataDiskTypes() string {
	return ""
}

func (self *SInstanceType) GetDataDiskMaxCount() int {
	// DiskQuantity为支持挂载的云盘数量上限，系统盘为 1 块云盘
	if self.DiskQuantity < 1 {
		return 0
	}
	return self.DiskQuantity - 1

}

func (self *SInstanceType) GetNicType() string {
	return "vpc"
}

func (self *SInstanceType) GetNicMaxCount() int {
	return 1
}

func (self *SInstanceType) GetGpuAttachable() bool {
	return self.GPUAmount != 0
}

func (self *SInstanceType) GetGpuSpec() string {
	// 解决部分机型  "GPUSpec": "0"
	if self.GPUAmount != 0 {
		return processGPUSpec(self.GPUSpec)
	} else {
		return ""
	}
}

// processGPUSpec
// 例如 ："NVIDIA A10*1"
func processGPUSpec(input string) string {
	// 判断字符串是否为空
	if input == "" {
		return input
	}

	// 判断字符串的第一个字符是否是字母
	firstChar := rune(input[0])
	if unicode.IsLetter(firstChar) {
		// 判断字符串中是否包含 "*"
		if index := strings.Index(input, "*"); index != -1 {
			// 判断"*"后是否有斜杠"/"
			if slashIndex := strings.Index(input[index:], "/"); slashIndex == -1 {
				// 如果没有斜杠"/"，则去除"*"后的内容
				return strings.TrimSpace(input[:index])
			}
		}
	}

	// 如果不符合条件，直接返回原始字符串
	return input
}

func (self *SInstanceType) GetGpuCount() string {
	return strconv.Itoa(self.GPUAmount)
}

func (self *SInstanceType) GetGpuMaxCount() int {
	return 0
}

func (self *SInstanceType) GetGPUMemorySizeMB() int {
	// GPUMemorySize 为单卡显存，需要与 GPUAmount 相乘得到总显存
	return int(self.GPUMemorySize * 1024 * float64(self.GPUAmount))
}

func (self *SInstanceType) GetIsBareMetal() bool {
	// aliyun 的裸金属规格定义在 InstanceCategory
	// https://api.aliyun.com/api/Ecs/2014-05-26/DescribeInstanceTypes
	return self.InstanceCategory == "ECS Bare Metal"
}
