package qcloud

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
)

func (r *SZone) GetCapability() (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotImplemented
}


func (r *SZone) GetICfelDiskType(diskType string) ([]*cloudprovider.DiskInfo, error) {
	return nil, nil
}
