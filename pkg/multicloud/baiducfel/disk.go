package baiducfel

import (
	"context"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SDisk struct {
	storage *SStorage
	// TODO instance

	multicloud.SDisk
	BaiduTags
	multicloud.SBillingBase

	Attachments    []Attachment `json:"attachments"`
	CreateTime     string       `json:"createTime"`
	Desc           string       `json:"desc"`
	DiskSizeInGB   int64        `json:"diskSizeInGB"`
	ExpireTime     time.Time    `json:"expireTime"`
	ID             string       `json:"volumeId"`
	Name           string       `json:"name"`
	PaymentTiming  string       `json:"paymentTiming"`
	Status         string       `json:"status"`
	StorageType    string       `json:"storageType"`
	Type           string       `json:"type"`
	ZoneName       string       `json:"zoneName"`
	IsSystemVolume bool         `json:"isSystemVolume"`
}

type Attachment struct {
	Device     string `json:"device"`
	InstanceID string `json:"instanceId"`
	Serial     string `json:"serial"`
	VolumeID   string `json:"volumeId"`
}

const ServiceDisk = "volumes"

var _ cloudprovider.ICloudDisk = (*SDisk)(nil)

// CreateISnapshot implements cloudprovider.ICloudDisk.
func (s *SDisk) CreateISnapshot(ctx context.Context, name string, desc string) (cloudprovider.ICloudSnapshot, error) {
	panic("unimplemented")
}

// Delete implements cloudprovider.ICloudDisk.
func (s *SDisk) Delete(ctx context.Context) error {
	var params = map[string]interface{}{
		"autoSnapshot":   "on",
		"manualSnapshot": "on",
		"recycle":        "on",
	}
	_, err := s.storage.zone.region.doPost(ServiceDisk, "/v2/volume/"+s.ID, params)
	return err
}

// GetAccessPath implements cloudprovider.ICloudDisk.
func (s *SDisk) GetAccessPath() string {
	return ""
}

// GetCacheMode implements cloudprovider.ICloudDisk.
func (s *SDisk) GetCacheMode() string {
	return ""
}

// GetDescription implements cloudprovider.ICloudDisk.
// Subtle: this method shadows the method (SDisk).GetDescription of SDisk.SDisk.
func (s *SDisk) GetDescription() string {
	return s.Desc
}

// GetDiskFormat implements cloudprovider.ICloudDisk.
func (s *SDisk) GetDiskFormat() string {
	return ""
}

// GetDiskSizeMB implements cloudprovider.ICloudDisk.
func (s *SDisk) GetDiskSizeMB() int {
	return int(s.DiskSizeInGB) * 1024
}

// GetDiskType implements cloudprovider.ICloudDisk.
func (s *SDisk) GetDiskType() string {
	return s.Type
}

// GetDriver implements cloudprovider.ICloudDisk.
func (s *SDisk) GetDriver() string {
	return ""
}

// GetFsFormat implements cloudprovider.ICloudDisk.
func (s *SDisk) GetFsFormat() string {
	return ""
}

// GetGlobalId implements cloudprovider.ICloudDisk.
func (s *SDisk) GetGlobalId() string {
	return s.ID
}

// GetISnapshots implements cloudprovider.ICloudDisk.
func (s *SDisk) GetISnapshots() ([]cloudprovider.ICloudSnapshot, error) {
	return nil, nil
}

// GetIStorage implements cloudprovider.ICloudDisk.
func (s *SDisk) GetIStorage() (cloudprovider.ICloudStorage, error) {
	return s.storage, nil
}

// GetId implements cloudprovider.ICloudDisk.
func (s *SDisk) GetId() string {
	return s.ID
}

// GetIsAutoDelete implements cloudprovider.ICloudDisk.
func (s *SDisk) GetIsAutoDelete() bool {
	return false
}

// GetIsNonPersistent implements cloudprovider.ICloudDisk.
func (s *SDisk) GetIsNonPersistent() bool {
	return false
}

// GetMountpoint implements cloudprovider.ICloudDisk.
func (s *SDisk) GetMountpoint() string {
	return ""
}

// GetName implements cloudprovider.ICloudDisk.
func (s *SDisk) GetName() string {
	return s.Name
}

// GetStatus implements cloudprovider.ICloudDisk.
func (s *SDisk) GetStatus() string {
	switch s.Status {
	case "Creating":
		return api.DISK_ALLOCATING
	case "Available", "InUse":
		return api.DISK_READY
	case "Scaling":
		return api.DISK_RESIZING
	case "Detaching":
		return api.DISK_DETACHING
	case "Attaching":
		return api.DISK_ATTACHING
	case "Deleting":
		return api.DISK_DEALLOC
	case "Error":
		return api.DISK_ALLOC_FAILED
	default:
		return api.DISK_READY
	}
}

// GetTemplateId implements cloudprovider.ICloudDisk.
func (s *SDisk) GetTemplateId() string {
	return ""
}

// Rebuild implements cloudprovider.ICloudDisk.
func (s *SDisk) Rebuild(ctx context.Context) error {
	panic("unimplemented")
}

// Reset implements cloudprovider.ICloudDisk.
func (s *SDisk) Reset(ctx context.Context, snapshotId string) (string, error) {
	panic("unimplemented")
}

// Resize implements cloudprovider.ICloudDisk.
func (s *SDisk) Resize(ctx context.Context, newSizeMB int64) error {
	panic("unimplemented")
}
