package ctyun

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

func (r *SZone) GetCapability() (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SZone) GetICfelDiskType(dt string) ([]*cloudprovider.DiskInfo, error) {
	product, err := self.region.getProduct()
	if err != nil {
		return nil, errors.Wrapf(err, "getProduct")
	}
	ret := make([]*cloudprovider.DiskInfo, 0)
	for i := range product.Ebs.StorageType {
		min, max := 10, 32768
		if dt == "sys" {
			min, max = 40, 2048
		}
		ret = append(ret, &cloudprovider.DiskInfo{
			StorageType: product.Ebs.StorageType[i].Type, 
			Name: product.Ebs.StorageType[i].Name,
			MinSizeGB: min,
			MaxSizeGB: max,
		})
		// ret[product.Ebs.StorageType[i].Type] = product.Ebs.StorageType[i].Name
	}
	return ret, nil
}
