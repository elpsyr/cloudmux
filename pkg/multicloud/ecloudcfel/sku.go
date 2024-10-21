package ecloudcfel

import (
	"context"
	"strings"
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

	gpuInfo map[string]string
}

var _ cloudprovider.ICfelCloudSku = (*SServerSku)(nil)

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
			for _, val := range sku {
				id := val.FlavorId + val.ZoneId
				if _, ok := tmp[id]; !ok {
					res = append(res, &val)
					tmp[id] = struct{}{}
				}
			}
		}
	}()

	vmTypes := []string{"memImprove", "common", "gpu", "commonIntroductory", "commonNetImprove", "compute", "computeNetImprove", "memNetImprove", "localStorage", "xlargeMemory", "highFrequency", "vgpu", "fpga", "highIO", "exclusive", "normalComputeImprove", "normalNetEnhance", "storeEnhance", "computeEnhance", "npu", "universal"}
	for _, zone := range self.izones {
		for _, vmType := range vmTypes {
			wg.Add(1)
			go func(vmType string, region *SZone) {
				query := map[string]string{
					"vmType": vmType,
					"region": region.Region,
				}
				request := NewConsoleRequest(self.ID, "/api/openapi-ecs/acl/v3/server/serverSpecs", query, nil)
				// request := NewNovaRequest(NewApiRequest(self.ID, "/api/openapi-ecs/acl/v3/server/serverSpecs", query, nil))
				skus := make([]SServerSku, 0)
				err := self.client.doList(context.Background(), request, &skus)
				if len(skus) == 0 || err != nil {
					wg.Done()
					return
				}
				for i := range skus {
					skus[i].ZoneId = region.GetGlobalId()
					skus[i].gpuInfo = self.gpuSkuInfo
				}
				skuChan <- skus
				wg.Done()
			}(vmType, zone.(*SZone))
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
		arr := strings.Split(info, "*")
		return strings.Trim(arr[0], " ")
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
		arr := strings.Split(info, "*")
		return strings.Trim(arr[1], " ")
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
	return self.MemorySize
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
	return 0
}

func (self *SServerSku) GetIsBareMetal() bool {
	return false
}
