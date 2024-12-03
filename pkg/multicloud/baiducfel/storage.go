package baiducfel

import (
	"errors"
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
)

type SStorage struct {
	multicloud.SStorageBase
	BaiduTags
	zone        *SZone
	storageType string

	ID string `json:"resourceId"` //创建磁盘返回的磁盘id
}

var storageTypes = []string{
	"enhanced_ssd_pl2",
	"enhanced_ssd_pl1",
	"premium_ssd",
	"hdd",
	"hp1",
}

var _ cloudprovider.ICloudStorage = (*SStorage)(nil)

type createDiskResp struct {
	VolumeIds []string `json:"volumeIds,omitempty"`
}

// CreateIDisk implements cloudprovider.ICloudStorage.
func (s *SStorage) CreateIDisk(conf *cloudprovider.DiskCreateConfig) (cloudprovider.ICloudDisk, error) {
	// https://cloud.baidu.com/doc/BCC/s/Ujwvyo1ta
	var params = map[string]interface{}{
		"storageType": s.storageType,
		"cdsSizeInGB": conf.SizeGb,
		// "snapshotId": snapshotId,
		// "purchaseCount" : purchaseCount,
		"name":        conf.Name,
		"description": conf.Desc,
	}
	res, err := s.zone.region.doPost(ServiceDisk, "/v2/volume", params)
	if err != nil {
		return nil, err
	}
	var resp createDiskResp
	if err = res.Unmarshal(&resp); err != nil {
		return nil, err
	}
	if len(resp.VolumeIds) == 0 {
		return nil, errors.New("not volume ids")
	}
	var disk = &SDisk{
		ID:           resp.VolumeIds[0],
		Name:         conf.Name,
		Desc:         conf.Desc,
		DiskSizeInGB: int64(conf.SizeGb),
	}
	return disk, nil
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
	return fmt.Sprintf("%s-%s-%s", "", s.zone.GetGlobalId(), s.storageType)
}

// GetIDiskById implements cloudprovider.ICloudStorage.
func (s *SStorage) GetIDiskById(idStr string) (cloudprovider.ICloudDisk, error) {
	var disk SDisk
	err := s.zone.region.doGet(ServiceDisk, "/v2/volume/"+idStr, nil, &disk)
	if err != nil {
		return nil,err
	}
	disk.storage = s
	return &disk, nil
}

type diskResp struct {
	NextMarker  string  `json:"nextMarker"`
	Marker      string  `json:"marker"`
	MaxKeys     int     `json:"maxKeys"`
	IsTruncated bool    `json:"isTruncated"`
	Disks       []SDisk `json:"volumes"`
}

// GetIDisks implements cloudprovider.ICloudStorage.
func (s *SStorage) GetIDisks() ([]cloudprovider.ICloudDisk, error) {
	var disks []SDisk
	var query = map[string]string{
		"maxKeys": "100",
	}

	var marker string
	for {
		var r diskResp
		if len(marker) > 0 {
			query["marker"] = marker
		}
		res, err := s.zone.region.doList(ServiceDisk, "/v2/volume", query)
		if err != nil {
			return nil, err
		}
		if err := res.Unmarshal(&r); err != nil {
			return nil, err
		}
		disks = append(disks, r.Disks...)
		if !r.IsTruncated {
			break
		}
		marker = r.NextMarker
	}

	var ret []cloudprovider.ICloudDisk
	for i := range disks {
		disks[i].storage = s
		ret = append(ret, &disks[i])
	}

	return ret, nil
}

// GetIStoragecache implements cloudprovider.ICloudStorage.
func (s *SStorage) GetIStoragecache() cloudprovider.ICloudStoragecache {
	panic("unimplemented")
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
	return s.storageType
}

// GetStatus implements cloudprovider.ICloudStorage.
func (s *SStorage) GetStatus() string {
	return "ready"
}

// GetStorageConf implements cloudprovider.ICloudStorage.
func (s *SStorage) GetStorageConf() jsonutils.JSONObject {
	panic("unimplemented")
}

// GetStorageType implements cloudprovider.ICloudStorage.
func (s *SStorage) GetStorageType() string {
	return s.storageType
}

// IsSysDiskStore implements cloudprovider.ICloudStorage.
func (s *SStorage) IsSysDiskStore() bool {
	return false
}
