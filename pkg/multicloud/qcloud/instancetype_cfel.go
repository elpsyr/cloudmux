package qcloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
)

// Verify that *SInstanceType implements ICloudSku interface
var _ cloudprovider.ICloudSku = (*SInstanceType)(nil)

type SCfelInstanceType struct {
	// cfel information
	TypeName   string
	CPUType    float64
	GPUDesc    string
	GpuCount   float64 // GPU 数量
	Hypervisor string  // 用于判断是否为裸金属
}

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	skus, err := self.GetCfelInstanceTypes()
	if err != nil {
		return nil, errors.Wrapf(err, "GetCfelInstanceTypes")
	}
	ret := []cloudprovider.ICfelCloudSku{}
	for i := range skus {
		ret = append(ret, &skus[i])
	}
	return ret, nil
}

// GetCfelInstanceTypes 获取 instanceType 信息
// 原调用 DescribeInstanceTypeConfigs 获取基础信息
// 现调用 DescribeUserAvailableInstanceTypes 获取详细信息
func (self *SRegion) GetCfelInstanceTypes() ([]SInstanceType, error) {
	params := make(map[string]string)
	params["Region"] = self.Region
	params["Filters.0.Name"] = "zone"
	instanceTypes := make([]SInstanceType, 0)

	zones, err := self.GetIZones()
	if err != nil {
		log.Errorf("GetInstanceTypes GetIZones fail %s", err)
		return nil, err
	}
	for _, izone := range zones {
		// 判断是否是可用 zone
		if izone.GetStatus() != api.ZONE_ENABLE {
			// 该区域不可用
			continue
		}
		params["Filters.0.Values.0"] = izone.GetId()

		//body, err := self.cvmRequest("DescribeInstanceTypeConfigs", params, true)
		body, err := self.cvmRequest("DescribeUserAvailableInstanceTypes", params, true)
		if err != nil {
			log.Errorf("DescribeUserAvailableInstanceTypes fail %s", err)
			continue
			// 目前 可用 zone 数据不一定真实
			//return nil, err
		}

		//err = body.Unmarshal(&instanceTypes, "InstanceTypeConfigSet")

		allInfo := new(DescribeInstanceConfigInfosUnmarshal)
		err = json.Unmarshal([]byte(body.String()), &allInfo)
		//err = body.Unmarshal(&allInfo)
		if err != nil {
			log.Errorf("Unmarshal instance type details fail %s", err)
			return nil, err
		}

		isStore := make(map[string]bool)

		for _, info := range allInfo.InstanceTypeQuotaSet {

			instanceType := SInstanceType{
				Zone:              info.Zone,
				InstanceType:      info.InstanceType,
				InstanceFamily:    info.InstanceFamily,
				GPU:               info.Gpu,
				CPU:               info.CPU,
				Memory:            info.Memory,
				CbsSupport:        "TRUE",
				InstanceTypeState: info.Status,
				SCfelInstanceType: SCfelInstanceType{
					TypeName:   info.TypeName,
					GpuCount:   float64(info.GpuCount),
					GPUDesc:    info.Externals.GPUDesc,
					Hypervisor: info.Externals.Hypervisor,
				},
			}
			// 判断使用有 GPU
			if info.GpuCount != 0 {
				// 普通类型从 Externals.GPUDesc 取数据
				if info.Externals.GPUDesc != "" {
					instanceType.GPUDesc = info.Externals.GPUDesc
				} else {
					// 目前 baremetal 类型主机 GPU 存储在 Remark 字段
					instanceType.GPUDesc = info.Remark
				}

				// todo:
				// 裁剪数据：
				// 1 * NVIDIA V100
				// 8 颗 NVIDIA V100

				// 定义正则表达式，匹配 * 或者 颗 前后的空格
				re := regexp.MustCompile(`\s*[*颗]\s*`)

				// 使用正则表达式分割字符串
				parts := re.Split(instanceType.GPUDesc, -1)

				// 输出提取出的数据（目标数据在切片的最后一个元素）
				instanceType.GPUDesc = parts[len(parts)-1]

			}
			if ok := isStore[info.Zone+info.InstanceType]; !ok {
				instanceTypes = append(instanceTypes, instanceType)
				isStore[info.Zone+info.InstanceType] = true
			}

		}
	}

	return instanceTypes, nil
}

// GetInstanceTypesPrice 获取 instanceType 价格
// 获取 zone 下对应 instance-type 的 规格以及价格信息
func (self *SRegion) GetInstanceTypesPrice(zoneID, instanceType string) (*DescribeInstanceConfigInfosUnmarshal, error) {
	params := make(map[string]string)
	params["Region"] = self.Region
	params["Filters.0.Name"] = "zone"
	params["Filters.0.Values.0"] = zoneID
	params["Filters.1.Name"] = "instance-type"
	params["Filters.1.Values.0"] = instanceType

	//body, err := self.cvmRequest("DescribeInstanceTypeConfigs", params, true)
	body, err := self.cvmRequest("DescribeUserAvailableInstanceTypes", params, true)
	if err != nil {
		log.Errorf("DescribeUserAvailableInstanceTypes fail %s", err)
		return nil, err
	}

	//err = body.Unmarshal(&instanceTypes, "InstanceTypeConfigSet")

	allInfo := new(DescribeInstanceConfigInfosUnmarshal)
	//err = json.Unmarshal([]byte(body.String()), &allInfo)
	err = body.Unmarshal(&allInfo)
	if err != nil {
		log.Errorf("Unmarshal instance type details fail %s", err)
		return nil, err
	}

	for _, info := range allInfo.InstanceTypeQuotaSet {

		fmt.Println(info.InstanceChargeType, "-----", info.Price)

	}

	return allInfo, nil
}

func (self *SInstanceType) GetId() string {
	return self.InstanceType
}

func (self *SInstanceType) GetName() string {
	return self.TypeName
}

func (self *SInstanceType) GetGlobalId() string {
	return self.InstanceType
}

func (self *SInstanceType) GetCreatedAt() time.Time {

	return time.Now()
}

func (self *SInstanceType) GetDescription() string {
	return self.TypeName
}

func (self *SInstanceType) GetStatus() string {
	return ""
}

func (self *SInstanceType) Refresh() error {
	return nil
}

func (self *SInstanceType) IsEmulated() bool {
	return false
}

func (self *SInstanceType) GetSysTags() map[string]string {
	return nil
}

func (self *SInstanceType) GetTags() (map[string]string, error) {
	return nil, nil
}

func (self *SInstanceType) SetTags(tags map[string]string, replace bool) error {
	return nil
}

func (self *SInstanceType) GetZoneID() string {
	return self.Zone
}

func (self *SInstanceType) GetInstanceTypeFamily() string {
	return self.InstanceFamily
}

func (self *SInstanceType) GetInstanceTypeCategory() string {
	return ""
}

func (self *SInstanceType) GetPrepaidStatus() string {
	return ""
}

func (self *SInstanceType) GetPostpaidStatus() string {
	return ""
}

func (self *SInstanceType) GetCpuArch() string {
	return ""
}

func (self *SInstanceType) GetCpuCoreCount() int {
	return self.CPU
}

func (self *SInstanceType) GetMemorySizeMB() int {
	return int(self.Memory * 1024)
}

func (self *SInstanceType) GetOsName() string {
	return ""
}

func (self *SInstanceType) GetSysDiskResizable() bool {
	return false
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

func (self *SInstanceType) GetAttachedDiskType() string {
	return ""
}

func (self *SInstanceType) GetAttachedDiskSizeGB() int {
	return 0
}

func (self *SInstanceType) GetAttachedDiskCount() int {
	return 0
}

func (self *SInstanceType) GetDataDiskTypes() string {
	return ""
}

func (self *SInstanceType) GetDataDiskMaxCount() int {
	return 0
}

func (self *SInstanceType) GetNicType() string {
	return "vpc"
}

func (self *SInstanceType) GetNicMaxCount() int {
	return 0
}

func (self *SInstanceType) GetGpuAttachable() bool {
	return self.GpuCount != 0
}

func (self *SInstanceType) GetGpuSpec() string {
	return self.GPUDesc
}

func (self *SInstanceType) GetGpuCount() string {
	//return strconv.Itoa(self.GpuCount)
	return strconv.FormatFloat(self.GpuCount, 'f', -1, 64)
}

func (self *SInstanceType) GetGpuMaxCount() int {
	return 0
}

func (self *SInstanceType) Delete() error {
	return nil
}

func (self *SInstanceType) GetGPUMemorySizeMB() int {
	return 0
}

func (self *SInstanceType) GetIsBareMetal() bool {
	// Qcloud 的裸金属规格定义在 Externals:{ "Hypervisor": "baremetal" }
	return self.Hypervisor == "baremetal"
}

// DescribeInstanceConfigInfosResp 是
// DescribeInstanceConfigInfos 接口的返回格式
type DescribeInstanceConfigInfosResp struct {
	Data struct {
		Response struct {
			InstanceConfigInfos []struct {
				Type             string `json:"type"`
				TypeName         string `json:"typeName"`
				Order            int    `json:"order"`
				InstanceFamilies []struct {
					InstanceFamily string `json:"instanceFamily"`
					TypeName       string `json:"TypeName"`
					Order          int    `json:"order"`
					InstanceTypes  []struct {
						InstanceType string  `json:"InstanceType"`
						CPU          int     `json:"Cpu"`
						Memory       int     `json:"Memory"`
						Gpu          int     `json:"Gpu"`
						Fpga         int     `json:"Fpga"`
						StorageBlock int     `json:"StorageBlock"`
						NetworkCard  int     `json:"NetworkCard"`
						MaxBandwidth float64 `json:"MaxBandwidth"`
						Frequency    string  `json:"Frequency"`
						CPUModelName string  `json:"CpuModelName"`
						Pps          int     `json:"Pps"`
						Remark       string  `json:"Remark"`
						Externals    struct {
							UnsupportNetworks   []string `json:"UnsupportNetworks"`
							ForceAcrossNodeFlag bool     `json:"forceAcrossNodeFlag"`
							GpuAttr             struct {
								Type string
							} `json:"GpuAttr"`
							GPUDesc string
						} `json:"Externals,omitempty"`
					} `json:"instanceTypes"`
					Architecture string `json:"Architecture"`
				} `json:"instanceFamilies"`
			} `json:"InstanceConfigInfos"`
			RequestID string `json:"RequestId"`
		} `json:"Response"`
	} `json:"data"`
	Code int `json:"code"`
}

type DescribeInstanceConfigInfosUnmarshal struct {
	InstanceTypeQuotaSet []struct {
		Zone               string `json:"Zone"`
		InstanceType       string `json:"InstanceType"`
		InstanceChargeType string `json:"InstanceChargeType"`
		NetworkCard        int    `json:"NetworkCard"`
		Externals          struct {
			GpuAttr struct {
				Type string `json:"Type"`
			} `json:"GpuAttr"`
			GPUDesc                 string `json:"GPUDesc"`
			InquiryWithEntireServer string `json:"InquiryWithEntireServer"`
			ExtraSpecs              struct {
				ForceAcrossNodeFlag bool `json:"forceAcrossNodeFlag"`
			} `json:"extra_specs"`
			UnsupportNetworks []string `json:"UnsupportNetworks"`
			StorageBlockAttr  struct {
				Type    string `json:"Type"`
				MinSize int    `json:"MinSize"`
				MaxSize int    `json:"MaxSize"`
			} `json:"StorageBlockAttr"`
			RequireNetworkFeatures        []string `json:"RequireNetworkFeatures"`
			UnsupportNetworkFeature       []string `json:"UnsupportNetworkFeature"`
			UnsupportLoginSettingsFeature []string `json:"UnsupportLoginSettingsFeature"`
			RdmaNicCount                  int      `json:"RdmaNicCount"`
			Hypervisor                    string   `json:"Hypervisor"`
			HypervisorSpec                []string `json:"HypervisorSpec"`
			RequiredEnhancedService       struct {
				MonitorService struct {
					Enabled string `json:"Enabled"`
				} `json:"MonitorService"`
			} `json:"RequiredEnhancedService"`
		} `json:"Externals,omitempty"`
		CPU                int              `json:"Cpu"`
		Memory             int              `json:"Memory"`
		InstanceFamily     string           `json:"InstanceFamily"`
		Architecture       string           `json:"Architecture"`
		TypeName           string           `json:"TypeName"`
		LocalDiskTypeList  []SLocalDiskType `json:"LocalDiskTypeList"`
		Status             string           `json:"Status"`
		SoldOutReason      string           `json:"SoldOutReason"`
		StorageBlockAmount int              `json:"StorageBlockAmount"`
		InstanceBandwidth  float64          `json:"InstanceBandwidth"`
		InstancePps        int              `json:"InstancePps"`
		CPUType            string           `json:"CpuType"`
		Frequency          string           `json:"Frequency"`
		Gpu                int              `json:"Gpu"`
		GpuCount           int              `json:"GpuCount"`
		Fpga               int              `json:"Fpga"`
		Remark             string           `json:"Remark"`
		ExtraProperty      struct {
		} `json:"ExtraProperty"`
		Disable      string `json:"Disable"`
		DeviceClass  string `json:"DeviceClass"`
		StorageBlock int    `json:"StorageBlock"`
		Price        struct {
			OriginalPrice     float64 `json:"OriginalPrice"`
			Discount          float64 `json:"Discount"`
			DiscountPrice     float64 `json:"DiscountPrice"`     // PREPAID
			UnitPrice         float64 `json:"UnitPrice"`         //  SPOTPAID、POSTPAID_BY_HOUR
			UnitPriceDiscount float64 `json:"UnitPriceDiscount"` //  SPOTPAID、POSTPAID_BY_HOUR
		} `json:"Price"`
	} `json:"InstanceTypeQuotaSet"`
	RequestID string `json:"RequestId"`
}
