package cucloudcfel

import (
	"context"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SImage struct {
	multicloud.SImageBase
	CuCloudTags

	storageCache *SStoragecache

	ImageId         string
	ServerId        string
	ImageAlias      string
	Name            string
	Url             string
	SrourceImageId  string
	Status          string
	SizeMb          int `json:"size"`
	IsPublic        int
	Note            string
	OsType          string
	MinDiskGB       int `json:"minDisk"`
	ImageType       string
	PublicImageType string
	BackupType      string
	BackupWay       string
	SnapshotId      string
	OsName          string
}

// Delete implements cloudprovider.ICloudImage.
func (s *SImage) Delete(ctx context.Context) error {
	panic("unimplemented")
}

// Export implements cloudprovider.ICloudImage.
// Subtle: this method shadows the method (SImageBase).Export of SImage.SImageBase.
func (s *SImage) Export(opts *cloudprovider.SImageExportOptions) ([]cloudprovider.SImageExportInfo, error) {
	panic("unimplemented")
}

// GetBios implements cloudprovider.ICloudImage.
func (s *SImage) GetBios() cloudprovider.TBiosType {
	panic("unimplemented")
}

// GetFullOsName implements cloudprovider.ICloudImage.
func (s *SImage) GetFullOsName() string {
	panic("unimplemented")
}

// GetGlobalId implements cloudprovider.ICloudImage.
func (s *SImage) GetGlobalId() string {
	panic("unimplemented")
}

// GetIStoragecache implements cloudprovider.ICloudImage.
func (s *SImage) GetIStoragecache() cloudprovider.ICloudStoragecache {
	panic("unimplemented")
}

// GetId implements cloudprovider.ICloudImage.
func (s *SImage) GetId() string {
	panic("unimplemented")
}

// GetImageFormat implements cloudprovider.ICloudImage.
func (s *SImage) GetImageFormat() string {
	panic("unimplemented")
}

// GetImageStatus implements cloudprovider.ICloudImage.
func (s *SImage) GetImageStatus() string {
	panic("unimplemented")
}

// GetImageType implements cloudprovider.ICloudImage.
func (s *SImage) GetImageType() cloudprovider.TImageType {
	panic("unimplemented")
}

// GetMinOsDiskSizeGb implements cloudprovider.ICloudImage.
func (s *SImage) GetMinOsDiskSizeGb() int {
	panic("unimplemented")
}

// GetMinRamSizeMb implements cloudprovider.ICloudImage.
func (s *SImage) GetMinRamSizeMb() int {
	panic("unimplemented")
}

// GetName implements cloudprovider.ICloudImage.
func (s *SImage) GetName() string {
	panic("unimplemented")
}

// GetOsArch implements cloudprovider.ICloudImage.
func (s *SImage) GetOsArch() string {
	panic("unimplemented")
}

// GetOsDist implements cloudprovider.ICloudImage.
func (s *SImage) GetOsDist() string {
	panic("unimplemented")
}

// GetOsLang implements cloudprovider.ICloudImage.
func (s *SImage) GetOsLang() string {
	panic("unimplemented")
}

// GetOsType implements cloudprovider.ICloudImage.
func (s *SImage) GetOsType() cloudprovider.TOsType {
	panic("unimplemented")
}

// GetOsVersion implements cloudprovider.ICloudImage.
func (s *SImage) GetOsVersion() string {
	panic("unimplemented")
}

// GetSizeByte implements cloudprovider.ICloudImage.
func (s *SImage) GetSizeByte() int64 {
	panic("unimplemented")
}

// GetStatus implements cloudprovider.ICloudImage.
func (s *SImage) GetStatus() string {
	panic("unimplemented")
}

var _ cloudprovider.ICloudImage = (*SImage)(nil)
