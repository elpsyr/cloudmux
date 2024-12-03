package baiducfel

import (
	"context"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SStoragecache struct {
	multicloud.SResourceBase
	BaiduTags
	region *SRegion
}

var _ cloudprovider.ICloudStoragecache = (*SStoragecache)(nil)

// GetDescription implements cloudprovider.ICloudStoragecache.
// Subtle: this method shadows the method (SResourceBase).GetDescription of SStoragecache.SResourceBase.
func (s *SStoragecache) GetDescription() string {
	panic("unimplemented")
}

// GetGlobalId implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) GetGlobalId() string {
	panic("unimplemented")
}

// GetICloudImages implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) GetICloudImages() ([]cloudprovider.ICloudImage, error) {
	panic("unimplemented")
}

// GetICustomizedCloudImages implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) GetICustomizedCloudImages() ([]cloudprovider.ICloudImage, error) {
	panic("unimplemented")
}

// GetIImageById implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) GetIImageById(extId string) (cloudprovider.ICloudImage, error) {
	panic("unimplemented")
}

// GetId implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) GetId() string {
	panic("unimplemented")
}

// GetName implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) GetName() string {
	panic("unimplemented")
}

// GetPath implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) GetPath() string {
	panic("unimplemented")
}

// GetStatus implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) GetStatus() string {
	panic("unimplemented")
}

// UploadImage implements cloudprovider.ICloudStoragecache.
func (s *SStoragecache) UploadImage(ctx context.Context, image *cloudprovider.SImageCreateOption, callback func(float32)) (string, error) {
	panic("unimplemented")
}
