package huawei

import (
	"net/url"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

type SCfelLoadbalancerSku struct {
	Name string
	Id   string
	Type string
}

// GetID implements cloudprovider.ICfelLoadbalancerSku.
func (s *SCfelLoadbalancerSku) GetID() string {
	return s.Id
}

// GetName implements cloudprovider.ICfelLoadbalancerSku.
func (s *SCfelLoadbalancerSku) GetName() string {
	return s.Name
}

// GetType implements cloudprovider.ICfelLoadbalancerSku.
func (s *SCfelLoadbalancerSku) GetType() string {
	return s.Type
}

var _ cloudprovider.ICfelLoadbalancerSku = (*SCfelLoadbalancerSku)(nil)

func (region *SRegion) GetLoadbalancerSkus() ([]cloudprovider.ICfelLoadbalancerSku, error) {
	ret := jsonutils.NewArray()
	query := url.Values{}
	var res []SCfelLoadbalancerSku
	for {
		resp, err := region.list(SERVICE_ELB, "elb/flavors", query)
		if err != nil {
			return nil, err
		}
		arr, err := resp.GetArray("flavors")
		if err != nil {
			return nil, err
		}
		ret.Add(arr...)
		marker, _ := resp.GetString("page_info", "next_marker")
		if len(marker) == 0 {
			break
		}
		query.Set("marker", marker)
	}
	err := ret.Unmarshal(&res)
	if err != nil {
		return nil, err
	}
	var result []cloudprovider.ICfelLoadbalancerSku
	for i := range res {
		result = append(result, &res[i])
	}
	return result, nil
}

// GetICfelCloudImage 获取华为云镜像
func (self *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	images := make([]SImage, 0)
	images, err := self.GetImages("", "", "", "")
	if err != nil {
		return nil, errors.Wrapf(err, "GetImages")
	}

	ret := []cloudprovider.ICloudImage{}
	for i := range images {
		images[i].storageCache = self.getStoragecache()
		ret = append(ret, &images[i])
	}
	return ret, nil
}
