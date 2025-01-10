package cucloudcfel

import (
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
)

type SStorage struct {
	multicloud.SStorageBase
	CuCloudTags
	zone        *SZone

	StorageType string `json:"productClass"`

	BillType      string `json:"billType"`
	BillUnit      string `json:"billUnit"`
	BmsSurplusNum string `json:"bmsSurplusNum"`
	BrandID       int64  `json:"brandId"`
	CanSellAlone  string `json:"canSellAlone"`
	ProductCata   int64  `json:"productCata"`
	// ProductClass     string      `json:"productClass"`
	ProductClassName string `json:"productClassName"`
	ProductDesc      string `json:"productDesc"`
	ProductDetail    string `json:"productDetail"`
	ProductID        string `json:"productId"`
	ProductMode      string `json:"productMode"`
	ProductName      string `json:"productName"`
	ProductPrice     int64  `json:"productPrice"`
	ResType          string `json:"resType"`
}

// CreateIDisk implements cloudprovider.ICloudStorage.
func (s *SStorage) CreateIDisk(conf *cloudprovider.DiskCreateConfig) (cloudprovider.ICloudDisk, error) {
	panic("unimplemented")
}

// GetCapacityMB implements cloudprovider.ICloudStorage.
func (s *SStorage) GetCapacityMB() int64 {
	return 0
}

// GetCapacityUsedMB implements cloudprovider.ICloudStorage.
func (s *SStorage) GetCapacityUsedMB() int64 {
	return 0
}

// GetEnabled implements cloudprovider.ICloudStorage.
func (s *SStorage) GetEnabled() bool {
	return true
}

// GetGlobalId implements cloudprovider.ICloudStorage.
func (s *SStorage) GetGlobalId() string {
	return fmt.Sprintf("%s-%s-%s", "", s.zone.GetGlobalId(), s.StorageType)
}

// GetIDiskById implements cloudprovider.ICloudStorage.
func (s *SStorage) GetIDiskById(idStr string) (cloudprovider.ICloudDisk, error) {
	panic("unimplemented")
}

// GetIDisks implements cloudprovider.ICloudStorage.
func (s *SStorage) GetIDisks() ([]cloudprovider.ICloudDisk, error) {
	panic("unimplemented")
}

// GetIStoragecache implements cloudprovider.ICloudStorage.
func (s *SStorage) GetIStoragecache() cloudprovider.ICloudStoragecache {
	return nil
}

// GetIZone implements cloudprovider.ICloudStorage.
func (s *SStorage) GetIZone() cloudprovider.ICloudZone {
	return s.zone
}

// GetId implements cloudprovider.ICloudStorage.
func (s *SStorage) GetId() string {
	return s.GetGlobalId()
}

// GetMediumType implements cloudprovider.ICloudStorage.
func (s *SStorage) GetMediumType() string {
	return ""
}

// GetMountPoint implements cloudprovider.ICloudStorage.
func (s *SStorage) GetMountPoint() string {
	return ""
}

// GetName implements cloudprovider.ICloudStorage.
func (s *SStorage) GetName() string {
	return s.StorageType
}

// GetStatus implements cloudprovider.ICloudStorage.
func (s *SStorage) GetStatus() string {
	return "ready"
}

// GetStorageConf implements cloudprovider.ICloudStorage.
func (s *SStorage) GetStorageConf() jsonutils.JSONObject {
	return nil
}

// GetStorageType implements cloudprovider.ICloudStorage.
func (s *SStorage) GetStorageType() string {
	return s.StorageType
}

// IsSysDiskStore implements cloudprovider.ICloudStorage.
func (s *SStorage) IsSysDiskStore() bool {
	return true
}

var _ cloudprovider.ICloudStorage = (*SStorage)(nil)
