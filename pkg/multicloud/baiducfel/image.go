package baiducfel

import (
	"context"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/util/rbacscope"
)

type SImage struct {
	multicloud.SImageBase
	BaiduTags
	// SZoneRegionBase

	storageCache *SStoragecache
	// imgInfo      *imagetools.ImageInfo

	ImageID      string `json:"imageId"`
	ImageName    string `json:"imageName"`
	MinSizeInGiB int64  `json:"minSizeInGiB"`
	OsArch       string `json:"osArch"`
	OsLang       string `json:"osLang"`
	OsName       string `json:"osName"`
	OsType       string `json:"osType"`
	OsVersion    string `json:"osVersion"`
}

var _ cloudprovider.ICloudImage = (*SImage)(nil)

const ServiceImage = "images"

// Delete implements cloudprovider.ICloudImage.
func (s *SImage) Delete(ctx context.Context) error {
	panic("unimplemented")
}

// Export implements cloudprovider.ICloudImage.
func (s *SImage) Export(opts *cloudprovider.SImageExportOptions) ([]cloudprovider.SImageExportInfo, error) {
	panic("unimplemented")
}

// GetBios implements cloudprovider.ICloudImage.
func (s *SImage) GetBios() cloudprovider.TBiosType {
	return cloudprovider.BIOS
}

// GetCreatedAt implements cloudprovider.ICloudImage.
func (s *SImage) GetCreatedAt() time.Time {
	return time.Now()
}

// GetDescription implements cloudprovider.ICloudImage.
func (s *SImage) GetDescription() string {
	return s.ImageName
}

// GetFullOsName implements cloudprovider.ICloudImage.
func (s *SImage) GetFullOsName() string {
	return s.OsName + " " + s.OsVersion
}

// GetGlobalId implements cloudprovider.ICloudImage.
func (s *SImage) GetGlobalId() string {
	return s.ImageID
}

// GetIStoragecache implements cloudprovider.ICloudImage.
func (s *SImage) GetIStoragecache() cloudprovider.ICloudStoragecache {
	return nil
}

// GetId implements cloudprovider.ICloudImage.
func (s *SImage) GetId() string {
	return s.ImageID
}

// GetImageFormat implements cloudprovider.ICloudImage.
func (s *SImage) GetImageFormat() string {
	return ""
}

// GetImageStatus implements cloudprovider.ICloudImage.
func (s *SImage) GetImageStatus() string {
	return "ready"
}

// GetImageType implements cloudprovider.ICloudImage.
func (s *SImage) GetImageType() cloudprovider.TImageType {
	return cloudprovider.ImageTypeShared
}

// GetMinOsDiskSizeGb implements cloudprovider.ICloudImage.
func (s *SImage) GetMinOsDiskSizeGb() int {
	return int(s.MinSizeInGiB)
}

// GetMinRamSizeMb implements cloudprovider.ICloudImage.
func (s *SImage) GetMinRamSizeMb() int {
	return 0
}

// GetName implements cloudprovider.ICloudImage.
func (s *SImage) GetName() string {
	return s.ImageName
}

// GetOsArch implements cloudprovider.ICloudImage.
func (s *SImage) GetOsArch() string {
	return s.OsArch
}

// GetOsDist implements cloudprovider.ICloudImage.
func (s *SImage) GetOsDist() string {
	return s.OsName
}

// GetOsLang implements cloudprovider.ICloudImage.
func (s *SImage) GetOsLang() string {
	return s.OsLang
}

// GetOsType implements cloudprovider.ICloudImage.
func (s *SImage) GetOsType() cloudprovider.TOsType {
	return cloudprovider.TOsType(s.OsType)
}

// GetOsVersion implements cloudprovider.ICloudImage.
func (s *SImage) GetOsVersion() string {
	return s.OsVersion
}

// GetProjectId implements cloudprovider.ICloudImage.
func (s *SImage) GetProjectId() string {
	return ""
}

// GetPublicScope implements cloudprovider.ICloudImage.
func (s *SImage) GetPublicScope() rbacscope.TRbacScope {
	panic("unimplemented")
}

// GetSizeByte implements cloudprovider.ICloudImage.
func (s *SImage) GetSizeByte() int64 {
	return -1
}

// GetStatus implements cloudprovider.ICloudImage.
func (s *SImage) GetStatus() string {
	return "ready"
}

// GetSubImages implements cloudprovider.ICloudImage.
func (s *SImage) GetSubImages() []cloudprovider.SSubImage {
	panic("unimplemented")
}

// IsEmulated implements cloudprovider.ICloudImage.
func (s *SImage) IsEmulated() bool {
	panic("unimplemented")
}

// Refresh implements cloudprovider.ICloudImage.
func (s *SImage) Refresh() error {
	panic("unimplemented")
}
