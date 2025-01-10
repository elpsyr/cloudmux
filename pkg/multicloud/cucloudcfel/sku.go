package cucloudcfel

import (
	"encoding/json"
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
	return s.GpuMemory > 0
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
	return s.FlavorId
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
	return ""
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
	List   []*flavor `json:"list"`
	ZoneId string
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
				if sku, ok := tmp[v.ProductMode]; ok {
					if v.BillType == "0" && v.BmsSurplusNum != "0" {
						sku.Subscription = true
					} else if v.BillType == "5" && v.BmsSurplusNum != "0" {
						sku.OnDemand = true
					}
				} else {

					sku = &SServerSku{
						region:    r,
						FlavorId:  v.ProductMode,
						SpecsName: v.ProductMode,
						ZoneId:    val.ZoneId,
					}
					if v.ProductArchitect != "pGPU" {
						sku.Arch = v.ProductArchitect
					}
					var detail []productDetail
					json.Unmarshal([]byte(v.ProductDetail), &detail)
					for _, d := range detail {
						if d.PrtyType == "memorySize" {
							vv, _ := strconv.Atoi(d.PrtyValue)
							sku.MemorySize = vv * 1024
						} else if d.PrtyType == "cpuNum" {
							vv, _ := strconv.Atoi(d.PrtyValue)
							sku.CpuCount = vv
						} else if d.PrtyType == "vGpu" { // 4 * NVIDIA Tesla T4，4 * 16GB
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
					if v.BillType == "0" && v.BmsSurplusNum != "0" {
						sku.Subscription = true
					} else if v.BillType == "5" && v.BmsSurplusNum != "0" {
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
								ret.ZoneId = zone.(*SZone).ZoneCode
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
		result = append(result, vv)
	}
	return result, nil
}
