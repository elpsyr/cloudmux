package ecloudcfel

import (
	"context"

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
}

var _ cloudprovider.ICfelCloudSku = (*SServerSku)(nil)

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	query := map[string]string{"vmType": "common"}
	request := NewNovaRequest(NewApiRequest(self.ID, "/api/openapi-ecs/acl/v3/server/serverSpecs", query, nil))
	skus := make([]SServerSku, 0, 5)
	err := self.client.doList(context.Background(), request, &skus)
	if err != nil {
		return nil, err
	}
	var res []cloudprovider.ICfelCloudSku
	for _, val := range skus {
		res = append(res, &val)
	}
	return res, nil
}

func (self *SServerSku) GetZoneID() string {
	return ""
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
	return ""
}

// GetGpuMaxCount implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetGpuMaxCount() int {
	return 0
}

// GetGpuSpec implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetGpuSpec() string {
	return ""
}

// GetId implements cloudprovider.ICfelCloudSku.
func (self *SServerSku) GetId() string {
	return self.FlavorId
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
