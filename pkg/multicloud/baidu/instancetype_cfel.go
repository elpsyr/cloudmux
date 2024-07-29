package baidu

import (
	"strconv"
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

// FlavorSpecList 实例套餐规格列表对象
type FlavorSpecList struct {
	ZoneResources []ZoneResources `json:"zoneResources"`
}
type ZoneResources struct {
	ZoneName     string       `json:"zoneName"`
	BccResources BccResources `json:"bccResources"`
}
type BccResources struct {
	FlavorGroups []FlavorGroups `json:"flavorGroups"`
}
type FlavorGroups struct {
	GroupID string     `json:"groupId"`
	Flavors []SFlavors `json:"flavors"`
}

// SFlavors  实例规格 （instanceType）
type SFlavors struct {
	multicloud.SInstanceBase
	CpuCount            int    `json:"cpuCount"`           // cpu数量
	MemoryCapacityInGB  int    `json:"memoryCapacityInGB"` // 内存容量（单位：GB）
	EphemeralDiskInGb   int    `json:"ephemeralDiskInGb"`  // 本地数据盘容量（单位：GB）
	EphemeralDiskCount  int    `json:"ephemeralDiskCount"` // 本地数据盘数量
	EphemeralDiskType   string `json:"ephemeralDiskType"`  // 本地数据盘类型
	GpuCardType         string `json:"gpuCardType"`        // gpu卡类型
	GpuCardCount        int    `json:"gpuCardCount"`       // gpu卡数量
	FpgaCardType        string `json:"fpgaCardType"`       // fpga卡类型
	FpgaCardCount       int    `json:"fpgaCardCount"`      // fpga卡数量
	ProductType         string `json:"productType"`        // 支持计费类型（PrePaid：预付费；PostPaid：后付费；both：预付费/后付费）
	Spec                string `json:"spec"`               // 实例套餐规格
	SpecID              string `json:"specId"`             // 实例套餐规格类型
	CpuModel            string `json:"cpuModel"`           // 处理器型号
	CpuGHz              string `json:"cpuGHz"`             // 处理器主频
	NetworkBandwidth    string `json:"networkBandwidth"`   // 内网带宽(Gbps)
	NetworkPackage      string `json:"networkPackage"`     // 网络收发包
	NetEthQueueCount    int    `json:"netEthQueueCount"`
	NetEthMaxQueueCount int    `json:"netEthMaxQueueCount"`
	// 赋值
	GroupID  string `json:"groupId"`
	ZoneName string `json:"zoneName"`
}

// Verify that *SInstanceType implements ICloudSku、ICloudSkuUltra
var _ cloudprovider.ICfelCloudSku = (*SFlavors)(nil)

func (S SFlavors) GetId() string {
	return S.Spec
}

func (S SFlavors) GetName() string {
	return S.Spec
}

func (S SFlavors) GetGlobalId() string {
	return S.Spec
}

func (S SFlavors) GetCreatedAt() time.Time {
	return time.Now()
}

func (S SFlavors) GetDescription() string {
	return ""
}

func (S SFlavors) GetStatus() string {
	return ""
}

func (S SFlavors) Refresh() error {
	return nil
}

func (S SFlavors) IsEmulated() bool {
	return false
}

func (S SFlavors) GetSysTags() map[string]string {
	return nil
}

func (S SFlavors) GetTags() (map[string]string, error) {
	return nil, nil
}

func (S SFlavors) SetTags(tags map[string]string, replace bool) error {
	return nil
}

func (S SFlavors) GetZoneID() string {
	return S.ZoneName
}

func (S SFlavors) GetInstanceTypeFamily() string {
	return S.GroupID
}

func (S SFlavors) GetInstanceTypeCategory() string {
	return S.GroupID
}

func (S SFlavors) GetPrepaidStatus() string {
	return ""
}

func (S SFlavors) GetPostpaidStatus() string {
	return ""
}

func (S SFlavors) GetCpuArch() string {
	return ""
}

func (S SFlavors) GetCpuCoreCount() int {
	return S.CpuCount
}

func (S SFlavors) GetMemorySizeMB() int {
	return S.MemoryCapacityInGB * 1024
}

func (S SFlavors) GetOsName() string {
	return ""
}

func (S SFlavors) GetSysDiskResizable() bool {
	return false
}

func (S SFlavors) GetSysDiskType() string {
	return S.EphemeralDiskType
}

func (S SFlavors) GetSysDiskMinSizeGB() int {
	return S.EphemeralDiskInGb
}

func (S SFlavors) GetSysDiskMaxSizeGB() int {
	return S.EphemeralDiskInGb
}

func (S SFlavors) GetAttachedDiskType() string {
	return S.EphemeralDiskType
}

func (S SFlavors) GetAttachedDiskSizeGB() int {
	return S.EphemeralDiskInGb
}

func (S SFlavors) GetAttachedDiskCount() int {
	return S.EphemeralDiskCount
}

func (S SFlavors) GetDataDiskTypes() string {
	return S.EphemeralDiskType
}

func (S SFlavors) GetDataDiskMaxCount() int {
	return 0
}

func (S SFlavors) GetNicType() string {
	return ""
}

func (S SFlavors) GetNicMaxCount() int {
	return S.NetEthMaxQueueCount
}

func (S SFlavors) GetGpuAttachable() bool {
	return S.GpuCardCount != 0
}

func (S SFlavors) GetGpuSpec() string {

	cardType, ok := gpuCardTypeMap[S.GpuCardType]
	if ok {
		return cardType
	}
	return S.GpuCardType
}

func (S SFlavors) GetGpuCount() string {
	return strconv.Itoa(S.GpuCardCount)
}

func (S SFlavors) GetGpuMaxCount() int {
	return S.GpuCardCount
}

func (S SFlavors) Delete() error {
	return nil
}

func (S SFlavors) GetIsBareMetal() bool {
	return false

}

func (S SFlavors) GetGPUMemorySizeMB() int {
	if S.GpuCardCount != 0 {
		memorySize, ok := gpuCardMemorySizeMap[S.GpuCardType]
		if ok {
			return memorySize * S.GpuCardCount * 1024
		}
	}
	return 0
}

var (
	// GPU 名称映射
	gpuCardTypeMap = map[string]string{
		"nTeslaP4":       "NVIDIA P4",
		"nTeslaA10Intel": "NVIDIA A10",
		"nTeslaA100-40G": "NVIDIA A100",
		"nTeslaA100-80G": "NVIDIA A100 80G", // A100 default 40 GB
		"nTeslaV100-16":  "NVIDIA V100 16G",
		"nTeslaV100-32":  "NVIDIA V100 32G", // V100 default 16 GB
		"nTeslaT4":       "NVIDIA T4",
	}
	// GPU 显存映射，单位 GB
	gpuCardMemorySizeMap = map[string]int{
		"nTeslaP4":       8,
		"nTeslaA10Intel": 24,
		"nTeslaA100-40G": 40,
		"nTeslaA100-80G": 80, // A100 default 40 GB
		"nTeslaV100-16":  16,
		"nTeslaV100-32":  32,
		"nTeslaT4":       16,
	}
)

// fetchFlavorSpec 查询实例套餐规格列表
// https://cloud.baidu.com/doc/BCC/s/Xk3pb75k1
func (region *SRegion) fetchFlavorSpec() ([]SFlavors, error) {
	// zoneName	String	Query参数	可用区名称  非必填
	body, err := region.client.list("bcc", region.Region, "/v2/instance/flavorSpec", nil)
	if err != nil {
		return nil, err
	}

	flavorSpecList := new(FlavorSpecList)
	err = body.Unmarshal(&flavorSpecList)
	if err != nil {
		return nil, err
	}

	flavors := make([]SFlavors, 0)
	flavorMap := map[string]bool{}
	for _, resource := range flavorSpecList.ZoneResources {
		for _, group := range resource.BccResources.FlavorGroups {
			for _, flavor := range group.Flavors {

				// 检查是否已经存入 list
				if flavorMap[flavor.Spec] {
					continue
				}
				flavorMap[flavor.Spec] = true

				flavor.ZoneName = resource.ZoneName
				flavor.GroupID = group.GroupID
				flavors = append(flavors, flavor)
			}
		}
	}

	return flavors, nil
}

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	skus, err := self.fetchFlavorSpec()
	if err != nil {
		return nil, errors.Wrapf(err, "fetchFlavorSpec")
	}
	var ret []cloudprovider.ICfelCloudSku
	for i := range skus {
		ret = append(ret, &skus[i])
	}
	return ret, nil
}

// fetchInstanceTypePrice 查询实例套餐价格
// https://cloud.baidu.com/doc/BCC/s/ijwvyo9im
// 请求参数
// version	String	是	URL参数	API版本号
// specId	String	是	RequestBody参数	实例规格族
// spec	String	是	RequestBody参数	实例套餐规格
// paymentTiming	String	是	RequestBody参数	付费方式，包括Postpaid(后付费)，Prepaid(预付费)两种。
// zoneName	String	是	RequestBody参数	可用区名称
// purchaseCount	int	否	RequestBody参数	查询在指定实例套餐规格下，任意数量实例的总价格，必须为大于0的整数，可选参数，缺省为1
// purchaseLength	int	是	RequestBody参数	时长，[1,2,3,4,5,6,7,8,9,12,24,36]，单位：月
func (region *SRegion) fetchInstanceTypePrice(instanceType, zoneName, paymentTiming string) (float64, error) {

	//split := strings.Split(instanceType, ".")

	params := map[string]interface{}{
		//"specId": instanceType,
		"spec":           instanceType,
		"paymentTiming":  paymentTiming,
		"zoneName":       zoneName,
		"purchaseCount":  1,
		"purchaseLength": 1,
	}

	body, err := region.client.post("bcc", region.Region, "/v2/instance/price", params)
	if err != nil {
		return -1, err
	}

	priceResponse := new(PriceResponse)
	err = body.Unmarshal(&priceResponse)
	if err != nil {
		return -1, err
	}

	if priceResponse != nil && len(priceResponse.Price) != 0 {
		return priceResponse.Price[0].SpecPrices[0].SpecPrice, nil
	}
	return -1, nil
}

type PriceResponse struct {
	Price []Price `json:"price"`
}
type SpecPrices struct {
	Spec      string  `json:"spec"`
	SpecPrice float64 `json:"specPrice"`
	Status    string  `json:"status"`
}
type Price struct {
	SpecID     string       `json:"specId"`
	SpecPrices []SpecPrices `json:"specPrices"`
}

// Verify that *SRegion implements ICfelCloudRegion
var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

func (region *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return -1, nil
}

func (region *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	pricePerMin, err := region.fetchInstanceTypePrice(instanceType, zoneID, "Postpaid")
	if err != nil {
		return -1, err
	}
	// baidu Postpaid 返回每分钟价格，需要返回小时计价
	return pricePerMin * 60, nil
}

func (region *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	return region.fetchInstanceTypePrice(instanceType, zoneID, "Prepaid")
}

func (region *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	return api.SkuStatusSoldout, nil
}

// 目前获取 instanceType 获取到的数据皆为可购买字段，所以实例的可购买状态皆为  available

func (region *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	return api.SkuStatusAvailable, nil
}

func (region *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	return api.SkuStatusAvailable, nil
}
