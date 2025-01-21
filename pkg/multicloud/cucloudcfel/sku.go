package cucloudcfel

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/log"
)

type SServerSku struct {
	multicloud.SInstanceBase
	multicloud.SResourceBase
	CuCloudTags

	region *SRegion

	Name         string
	ProductId    string
	FlavorId     string
	CpuCount     int `json:"cpu"`
	MemorySize   int `json:"ram"`
	GpuSpec      string
	GpuCount     string
	GpuMemory    int
	OnDemand     bool // 是否支持按量付费
	Subscription bool // 是否支持包年包月
	SpecsName    string
	SpecsType    string
	VmType       string
	ZoneId       string
	Arch         string

	ProductIds []string
}

// Delete implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) Delete() error {
	panic("unimplemented")
}

// GetAttachedDiskCount implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetAttachedDiskCount() int {
	return 0
}

// GetAttachedDiskSizeGB implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetAttachedDiskSizeGB() int {
	return 0
}

// GetAttachedDiskType implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetAttachedDiskType() string {
	return ""
}

// GetCpuArch implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetCpuArch() string {
	return s.Arch
}

// GetCpuCoreCount implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetCpuCoreCount() int {
	return s.CpuCount
}

// GetDataDiskMaxCount implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetDataDiskMaxCount() int {
	return 0
}

// GetDataDiskTypes implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetDataDiskTypes() string {
	return ""
}

// GetGPUMemorySizeMB implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetGPUMemorySizeMB() int {
	return s.GpuMemory
}

// GetGlobalId implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetGlobalId() string {
	return s.FlavorId
}

// GetGpuAttachable implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetGpuAttachable() bool {
	return s.GpuCount != ""
}

// GetGpuCount implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetGpuCount() string {
	return s.GpuCount
}

// GetGpuMaxCount implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetGpuMaxCount() int {
	return 0
}

// GetGpuSpec implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetGpuSpec() string {
	return s.GpuSpec
}

// GetId implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetId() string {
	return s.FlavorId
}

// GetInstanceTypeCategory implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetInstanceTypeCategory() string {
	return ""
}

// GetInstanceTypeFamily implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetInstanceTypeFamily() string {
	return ""
}

// GetIsBareMetal implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetIsBareMetal() bool {
	return false
}

// GetMemorySizeMB implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetMemorySizeMB() int {
	return s.MemorySize
}

// GetName implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetName() string {
	return s.Name
}

// GetNicMaxCount implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetNicMaxCount() int {
	return 0
}

// GetNicType implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetNicType() string {
	return ""
}

// GetOsName implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetOsName() string {
	return ""
}

// GetPostpaidStatus implements cloudprovider.ICfelCloudSku.
// Subtle: this method shadows the method (SInstanceBase).GetPostpaidStatus of SServerSku.SInstanceBase.
func (s *SServerSku) GetPostpaidStatus() string {
	if s.OnDemand {
		return api.SkuStatusAvailable
	}
	return api.SkuStatusSoldout
}

// GetPrepaidStatus implements cloudprovider.ICfelCloudSku.
// Subtle: this method shadows the method (SInstanceBase).GetPrepaidStatus of SServerSku.SInstanceBase.
func (s *SServerSku) GetPrepaidStatus() string {
	if s.Subscription {
		return api.SkuStatusAvailable
	}
	return api.SkuStatusSoldout
}

// GetSpotpaidStatus implements cloudprovider.ICfelCloudSku.
// Subtle: this method shadows the method (SInstanceBase).GetSpotpaidStatus of SServerSku.SInstanceBase.
func (s *SServerSku) GetSpotpaidStatus() string {
	return api.SkuStatusSoldout
}

// GetStatus implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetStatus() string {
	return "ready"
}

// GetSysDiskMaxSizeGB implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetSysDiskMaxSizeGB() int {
	return 0
}

// GetSysDiskMinSizeGB implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetSysDiskMinSizeGB() int {
	return 0
}

// GetSysDiskResizable implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetSysDiskResizable() bool {
	return false
}

// GetSysDiskType implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetSysDiskType() string {
	return ""
}

// GetZoneID implements cloudprovider.ICfelCloudSku.
func (s *SServerSku) GetZoneID() string {
	return s.ZoneId
}

func (s *SServerSku) GetDescription() string {
	return s.ProductId
}

var _ cloudprovider.ICfelCloudSku = (*SServerSku)(nil)

type flavor struct {
	BillUnit         string  `json:"billUnit"`
	BmsSurplusNum    string  `json:"bmsSurplusNum"`
	BillType         string  `json:"billType"`
	CanSellAlone     string  `json:"canSellAlone"`
	ProductArchitect string  `json:"productArchitect"`
	ProductClassName string  `json:"productClassName"`
	ProductDesc      string  `json:"productDesc"`
	ProductDetail    string  `json:"productDetail"`
	ProductID        string  `json:"productId"`
	ProductMode      string  `json:"productMode"`
	ProductPrice     float64 `json:"productPrice"`
	ResType          string  `json:"resType"`
	ResourceID       string  `json:"resourceId"`
	SubBrandID       int64   `json:"subBrandId"`
	Version          string  `json:"version"`
}

type productDetail struct {
	MeasurementUnit string `json:"measurementUnit"`
	PrtyCode        string `json:"prtyCode"`
	PrtyName        string `json:"prtyName"`
	PrtyType        string `json:"prtyType"`
	PrtyValue       string `json:"prtyValue"`
}

type flavorResp struct {
	List          []*flavor `json:"list"`
	ZoneId        string
	ResourcesType string
}

func (r *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {

	if r.zones == nil {
		err := r.fetchZones()
		if err != nil {
			return nil, err
		}
	}

	var serviceType = []string{"1001", "1101", "8073"} // 1001 x86,1101 gpu,8073 arm
	var billType = []string{"0", "5"}                  // 0 包年包月，5 按量付费

	var ch = make(chan *flavorResp)
	defer close(ch)

	var wg sync.WaitGroup
	var tmp = make(map[string]*SServerSku)

	go func() {
		for val := range ch {
			for _, v := range val.List {
				if v.BmsSurplusNum == "0" {
					continue
				}
				if sku, ok := tmp[v.ProductMode]; ok {
					if v.BillType == "0" {
						sku.Subscription = true
					} else if v.BillType == "5" {
						sku.OnDemand = true
					}
					sku.ProductIds = append(sku.ProductIds, v.BillType+"-"+v.ProductID+"-"+val.ResourcesType)
				} else {

					sku = &SServerSku{
						region:     r,
						FlavorId:   v.ProductMode,
						SpecsName:  v.ProductMode,
						ZoneId:     val.ZoneId,
						ProductIds: []string{v.BillType + "-" + v.ProductID + "-" + val.ResourcesType},
					}
					// if v.ProductArchitect != "pGPU" {
					sku.Arch = v.ProductArchitect
					// }
					var detail []productDetail
					json.Unmarshal([]byte(v.ProductDetail), &detail)

					for _, d := range detail {
						if d.PrtyCode == "memorySize" {
							vv, _ := strconv.Atoi(d.PrtyValue)
							sku.MemorySize = vv * 1024
						} else if d.PrtyCode == "cpuNum" {
							vv, _ := strconv.Atoi(d.PrtyValue)
							sku.CpuCount = vv
						} else if d.PrtyCode == "vGpu" { // 4 * NVIDIA Tesla T4，4 * 16GB
							arr := strings.Split(d.PrtyValue, "，")
							if len(arr) >= 2 {
								arr1 := strings.Split(arr[0], "*")
								if len(arr1) >= 2 {
									sku.GpuCount = strings.TrimSpace(arr1[0])
									sku.GpuSpec = strings.TrimSpace(arr1[1])
								}
								arr2 := strings.Split(arr[0], "*")
								if len(arr2) >= 2 {
									vv, err := strconv.Atoi(strings.Trim(strings.TrimSpace(arr2[1]), "GB"))
									if err == nil {
										sku.GpuMemory = vv * 1024
									}
								}
							}
						}
					}
					if v.BillType == "0" {
						sku.Subscription = true
					} else if v.BillType == "5" {
						sku.OnDemand = true
					}
					tmp[v.ProductMode] = sku
				}
			}
		}
	}()

	for _, zone := range r.zones {
		for _, v := range serviceType {
			for _, vv := range billType {
				wg.Add(1)
				go func(billType, serviceType string, zone cloudprovider.ICloudZone) {
					defer wg.Done()
					params := map[string]interface{}{
						"currPage":      1,
						"pageSize":      10000,
						"specilZone":    zone.(*SZone).ZoneCode,
						"serviceTypeId": serviceType, // 1001 x86,1101 gpu,8073 arm
						"billType":      billType,    // 0 包年包月，5 按量付费
						"version":       "1",
						"cloudRegionId": r.CloudRegionId,
						"zoneId":        zone.(*SZone).ZoneID,
					}
					var ret = flavorResp{}
					res, err := r.client.postWithToken("product-cons/product/productquery/flavorsHasResChk", params)
					if err == nil {
						if err = res.Unmarshal(&ret, "result"); err == nil {
							if len(ret.List) > 0 {
								log.Infof("get flavors res:%d", len(ret.List))
								ret.ZoneId = zone.GetGlobalId()
								ret.ResourcesType = serviceType
								ch <- &ret
							}
						}
					} else {
						log.Warningf("get flavors error:%v zoneCode:%s", err, zone.(*SZone).ZoneCode)
					}
				}(vv, v, zone)
			}
		}
	}

	wg.Wait()
	time.Sleep(2 * time.Second)
	var result []cloudprovider.ICfelCloudSku
	for _, vv := range tmp {
		sort.Strings(vv.ProductIds)
		vv.Name = strings.Join(vv.ProductIds, "@")
		result = append(result, vv)
	}
	return result, nil
}

func (region *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, nil
}

func (region *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return region.GetPrice(zoneID, instanceType, true)
}

func (region *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	return region.GetPrice(zoneID, instanceType, false)
}

func (region *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	return api.SkuStatusSoldout, nil
}

// 目前获取 instanceType 获取到的数据皆为可购买字段，所以实例的可购买状态皆为  available
func (region *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	arr := strings.Split(instanceType, "@")
	if len(arr) == 2 {
		return api.SkuStatusAvailable, nil
	} else {
		arr = strings.Split(arr[0], "-")
		if len(arr) == 2 && arr[0] == "5" {
			return api.SkuStatusAvailable, nil
		} else {
			return api.SkuStatusSoldout, nil
		}
	}
}

func (region *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {

	arr := strings.Split(instanceType, "@")
	if len(arr) == 2 {
		return api.SkuStatusAvailable, nil
	} else {
		arr = strings.Split(arr[0], "-")
		if len(arr) == 2 && arr[0] == "0" {
			return api.SkuStatusAvailable, nil
		} else {
			return api.SkuStatusSoldout, nil
		}
	}
}

func (region *SRegion) GetPrice(zoneID, instanceType string, isOnDemand bool) (float64, error) {
	if region.zones == nil {
		if err := region.fetchZones(); err != nil {
			return 0, err
		}
	}
	var zone *SZone
	for _, z := range region.zones {
		if z.GetGlobalId() == zoneID {
			zone = z.(*SZone)
		}
	}
	if zone == nil {
		return 0, fmt.Errorf("zone not found")
	}
	var purchaseUnit = 3 // 3 包年包月 2 按量付费
	arr := strings.Split(instanceType, "@")
	// if len(arr) < 2 {
	// 	return 0, fmt.Errorf("instanceType format error;instanceType:%s", instanceType)
	// }
	if len(arr) == 1 {
		arr = append(arr, arr[0])
	}
	var productId = arr[0]
	if isOnDemand {
		purchaseUnit = 2
		productId = arr[1]
	}
	arr = strings.Split(productId, "-")
	if len(arr) < 3 {
		return 0, fmt.Errorf("instanceType format error;instanceType:%s", instanceType)
	}
	params := map[string]interface{}{
		"accountId":               region.client.accountId,
		"userId":                  region.client.userId,
		"isOnDemand":              isOnDemand,
		"orderChannel":            1,
		"orderType":               1,
		"payType":                 1,
		"autoRenew":               false,
		"autoRenewPrice":          0,
		"autoRenewPurchaseUnit":   3, // 3 包年包月 2 按量付费
		"autoRenewPurchaseNumber": 1,
		"autoRenewTimes":          1,
		"description":             "",
		"subOrders": []map[string]interface{}{
			{
				"packageCount":   1,
				"purchaseNumber": 1,
				"purchaseUnit":   purchaseUnit, // 3 包年包月 2 按量付费
				"subOrderItems": []map[string]interface{}{
					{
						"instanceType":          arr[2],
						"master":                true,
						"productId":             arr[1],
						"resourceType":          arr[2],
						"zone":                  zone.ZoneID,
						"resourceConfiguration": "{\"cpuNum\":10,\"memorySize\":1,\"dataType\":\"public\",\"diskSize\":50}",
						"isOnDemand":            false,
					},
				},
			},
		},
		"businessType": 0,
	}
	ret, err := region.client.postWithToken("bill-out-cons/bill/charge/order/getTotalProductPriceFromBill", params)
	if err != nil {
		return 0, err
	}
	var rr = struct {
		StandardPrice float64
		OriginalPrice float64
	}{}
	if err = ret.Unmarshal(&rr, "result"); err != nil {
		return 0, err
	}
	if rr.StandardPrice == 0 {
		return rr.OriginalPrice / 1000, nil
	}
	return rr.StandardPrice / 1000, nil
}
