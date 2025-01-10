package cucloudcfel

import "yunion.io/x/cloudmux/pkg/cloudprovider"

var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

func (r *SRegion) GetICfelCloudImageById(id string) (cloudprovider.ICloudImage, error) {
	return nil, nil
}

func (self *SRegion) GetInstanceMatchImage(instanceType string) ([]cloudprovider.ICloudImage, error) {
	var query = map[string]interface{}{
		"pageSize": "1000",
		"cloudRegionCode":self.CloudRegionCode,
		"architecture":"",
		"dataType":"public",
	}
	var imgs []SImage

	var marker string
	for {
		if len(marker) > 0 {
			query["marker"] = marker
		}
		res, err := self.client.get("/instance/v1/product/images", query)
		if err != nil {
			return nil, err
		}
		var r interface{}
		err = res.Unmarshal(&r)
		if err != nil {
			return nil, err
		}

	}

	var ret []cloudprovider.ICloudImage
	for i := range imgs {
		ret = append(ret, &imgs[i])
	}
	return ret, nil
}
