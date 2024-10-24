package ecloudcfel

import (
	"context"
	"sync"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SServerSku struct {
	multicloud.SInstanceBase
	multicloud.SResourceBase
	EcloudTags
	region *SRegion

	FlavorId   string
	CpuCount   int `json:"cpu"`
	MemorySize int `json:"ram"`
	// ZoneDesc   string
	// ZoneName   string
	SpecsName string
	SpecsType string
	VmType    string
	ZoneId    string

	gpuInfo map[string]*gpuInfo
}

var _ cloudprovider.ICfelCloudSku = (*SServerSku)(nil)

type vmSpecs struct {
	SpecsName string `json:"specsName"`
	SoldOut   string `json:"soldOut"`
}

// 规格类型 https://ecloud.10086.cn/op-help-center/doc/article/71795
var vmTypes = map[string]string{
	"1ab7650b904cef5ac2078a00fc2eca60": "commonNetImprove",
	"5691272bb306174e7b2ed84210a28fd3": "memNetImprove",
	"762db914c221ddb9fd8ce7f456e8bcf8": "gpu",
	"795ed2017b3f26a49e34b1cb2bd0ac3d": "commonIntroductory",
	"a715136ad3147ceb0bc0ba005b8e8897": "computeNetImprove",
	"aa40a113d03520947e613bae515c21ad": "highFrequency",
	"bbcb90abefd93fb38ce994be310e0854": "compute",
	"bf81ee49b985f547233726b3fbbad3a6": "memImprove",
	"d0a7122eacdd7f997285ff20f0301519": "xlargeMemory",
	"ebd890c0712a4b42bb147a55ef16057e": "storeEnhance",
	"fe0ca5a14938cfab007929bce91e06b9": "common",
}

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	err := self.fetchZones()
	if err != nil {
		return nil, err
	}
	var skuChan = make(chan []SServerSku)
	var wg sync.WaitGroup

	var res []cloudprovider.ICfelCloudSku
	var tmp = make(map[string]struct{})
	go func() {
		for sku := range skuChan {
			for i := range sku {
				id := sku[i].SpecsName + sku[i].ZoneId
				if _, ok := tmp[id]; !ok {
					res = append(res, &sku[i])
					tmp[id] = struct{}{}
				}
			}
		}
	}()

	// vmTypes := []string{"memImprove", "common", "gpu", "commonIntroductory", "commonNetImprove", "compute", "computeNetImprove", "memNetImprove", "localStorage", "xlargeMemory", "highFrequency", "vgpu", "fpga", "highIO", "exclusive", "normalComputeImprove", "normalNetEnhance", "storeEnhance", "computeEnhance", "npu", "universal"}
	for _, zone := range self.izones {
		for offerId, vmType := range vmTypes {
			wg.Add(1)
			go func(vmType, offerId string, region *SZone) {
				query := map[string]string{
					"offerId": offerId,
					"region":  region.Region,
				}
				req := NewConsoleRequest(self.ID, "/api/openapi-ecs/acl/v3/server/specsName", query, nil)
				var vmSpecs []vmSpecs
				err := self.client.doGet(context.Background(), req, &vmSpecs)
				if err != nil {
					wg.Done()
				}
				var tmp = make(map[string]struct{}) //已售罄规格
				for _, val := range vmSpecs {
					if val.SoldOut == "1" {
						tmp[val.SpecsName] = struct{}{}
					}
				}
				query = map[string]string{
					"vmType": vmType,
					"region": region.Region,
				}
				request := NewConsoleRequest(self.ID, "/api/openapi-ecs/acl/v3/server/serverSpecs", query, nil)
				// request := NewNovaRequest(NewApiRequest(self.ID, "/api/openapi-ecs/acl/v3/server/serverSpecs", query, nil))
				skus := make([]SServerSku, 0)
				err = self.client.doList(context.Background(), request, &skus)
				if len(skus) == 0 || err != nil {
					wg.Done()
					return
				}
				var ret = make([]SServerSku, 0)
				for i := range skus {
					if _, ok := tmp[skus[i].SpecsName]; ok {
						continue
					}
					skus[i].ZoneId = region.GetGlobalId()
					skus[i].gpuInfo = self.gpuSkuInfo
					ret = append(ret, skus[i])
				}
				skuChan <- ret
				wg.Done()
			}(vmType, offerId, zone.(*SZone))
		}
	}
	wg.Wait()
	close(skuChan)

	return res, nil
}

func (self *SServerSku) GetZoneID() string {
	return self.ZoneId
}

// GetAttachedDiskCount implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetAttachedDiskCount() int {
	return 0
}

// GetAttachedDiskSizeGB implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetAttachedDiskSizeGB() int {
	return 0
}

// GetAttachedDiskType implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetAttachedDiskType() string {
	return ""
}

// GetCpuArch implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetCpuArch() string {
	return "x86"
}

// GetCpuCoreCount implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetCpuCoreCount() int {
	return self.CpuCount
}

// GetDataDiskMaxCount implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetDataDiskMaxCount() int {
	return 0
}

// GetDataDiskTypes implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetDataDiskTypes() string {
	return ""
}

// GetGlobalId implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetGlobalId() string {
	return self.FlavorId
}

// GetGpuAttachable implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetGpuAttachable() bool {
	return false
}

// GetGpuCount implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetGpuCount() string {
	if info, ok := self.gpuInfo[self.SpecsName]; ok {
		return info.GpuCount
	}
	return ""
}

// GetGpuMaxCount implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetGpuMaxCount() int {
	return 0
}

// GetGpuSpec implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetGpuSpec() string {
	if info, ok := self.gpuInfo[self.SpecsName]; ok {
		return info.Spec
	}
	return ""
}

// GetId implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetId() string {
	return self.SpecsName
}

// GetInstanceTypeCategory implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetInstanceTypeCategory() string {
	return ""
}

// GetInstanceTypeFamily implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetInstanceTypeFamily() string {
	return ""
}

// GetMemorySizeMB implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetMemorySizeMB() int {
	return self.MemorySize * 1024
}

// GetName implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetName() string {
	return self.SpecsName
}

// GetNicMaxCount implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetNicMaxCount() int {
	return 0
}

// GetNicType implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetNicType() string {
	return ""
}

// GetOsName implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetOsName() string {
	return ""
}

// GetStatus implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetStatus() string {
	return "ready"
}

// GetSysDiskMaxSizeGB implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetSysDiskMaxSizeGB() int {
	return 0
}

// GetSysDiskMinSizeGB implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetSysDiskMinSizeGB() int {
	return 0
}

// GetSysDiskResizable implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetSysDiskResizable() bool {
	return false
}

// GetSysDiskType implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetSysDiskType() string {
	return ""
}

// Delete implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) Delete() error {
	panic("unimplemented")
}

func (self *SServerSku) GetGPUMemorySizeMB() int {
	if info, ok := self.gpuInfo[self.SpecsName]; ok {
		return info.GPUMemorySizeGB * 1024
	}
	return 0
}

func (self *SServerSku) GetIsBareMetal() bool {
	return false
}
