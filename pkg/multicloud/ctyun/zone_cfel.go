package ctyun

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

func (r *SZone) GetCapability() (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SZone) GetICfelDiskType(dt string) (map[string]interface{}, error) {
	product, err := self.region.getProduct()
	if err != nil {
		return nil, errors.Wrapf(err, "getProduct")
	}
	ret := make(map[string]interface{})
	for i := range product.Ebs.StorageType {
		ret[product.Ebs.StorageType[i].Type] = product.Ebs.StorageType[i].Name
	}
	return ret, nil
}
