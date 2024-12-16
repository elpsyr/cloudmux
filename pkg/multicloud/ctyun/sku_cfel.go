package ctyun

import (
	"strconv"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

var _ cloudprovider.ICfelCloudSku = (*ServerSku)(nil)

// Delete implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) Delete() error {
	panic("unimplemented")
}

// GetAttachedDiskCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetAttachedDiskCount() int {
	panic("unimplemented")
}

// GetAttachedDiskSizeGB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetAttachedDiskSizeGB() int {
	panic("unimplemented")
}

// GetAttachedDiskType implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetAttachedDiskType() string {
	panic("unimplemented")
}

// GetCpuArch implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetCpuArch() string {
	return s.CPUInfo
}

// GetCpuCoreCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetCpuCoreCount() int {
	return s.FlavorCPU
}

// GetCreatedAt implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetCreatedAt() time.Time {
	panic("unimplemented")
}

// GetDataDiskMaxCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetDataDiskMaxCount() int {
	panic("unimplemented")
}

// GetDataDiskTypes implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetDataDiskTypes() string {
	panic("unimplemented")
}

// GetDescription implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetDescription() string {
	return s.FlavorName
}

// GetGPUMemorySizeMB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGPUMemorySizeMB() int {
	return s.VideoMemSize * 1024
}

// GetGlobalId implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGlobalId() string {
	return s.FlavorId
}

// GetGpuAttachable implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGpuAttachable() bool {
	panic("unimplemented")
}

// GetGpuCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGpuCount() string {
	return strconv.Itoa(s.GpuCount)
}

// GetGpuMaxCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGpuMaxCount() int {
	panic("unimplemented")
}

// GetGpuSpec implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGpuSpec() string {
	panic("unimplemented")
}

// GetId implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetId() string {
	return s.FlavorId
}

// GetInstanceTypeCategory implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetInstanceTypeCategory() string {
	panic("unimplemented")
}

// GetInstanceTypeFamily implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetInstanceTypeFamily() string {
	panic("unimplemented")
}

// GetIsBareMetal implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetIsBareMetal() bool {
	return false
}

// GetMemorySizeMB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetMemorySizeMB() int {
	return s.FlavorRAM*1024
}

// GetName implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetName() string {
	return s.FlavorName
}

// GetNicMaxCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetNicMaxCount() int {
	return 0
}

// GetNicType implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetNicType() string {
	panic("unimplemented")
}

// GetOsName implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetOsName() string {
	panic("unimplemented")
}

// GetPostpaidStatus implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetPostpaidStatus() string {
	panic("unimplemented")
}

// GetPrepaidStatus implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetPrepaidStatus() string {
	panic("unimplemented")
}

// GetSpotpaidStatus implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSpotpaidStatus() string {
	panic("unimplemented")
}

// GetStatus implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetStatus() string {
	return ""
}

// GetSysDiskMaxSizeGB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSysDiskMaxSizeGB() int {
	panic("unimplemented")
}

// GetSysDiskMinSizeGB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSysDiskMinSizeGB() int {
	panic("unimplemented")
}

// GetSysDiskResizable implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSysDiskResizable() bool {
	panic("unimplemented")
}

// GetSysDiskType implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSysDiskType() string {
	panic("unimplemented")
}

// GetZoneID implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetZoneID() string {
	panic("unimplemented")
}

// IsEmulated implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) IsEmulated() bool {
	panic("unimplemented")
}

// Refresh implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) Refresh() error {
	panic("unimplemented")
}

// SetTags implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) SetTags(tags map[string]string, replace bool) error {
	panic("unimplemented")
}

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	self.GetServerSkus("")
	return nil, nil
}
