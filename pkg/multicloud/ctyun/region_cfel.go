package ctyun

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

func (self *SRegion) GetInstanceMatchImage(instanceType string) ([]cloudprovider.ICloudImage, error) {
	pageNo := 1
	params := map[string]interface{}{
		"pageNo":     pageNo,
		"pageSize":   50,
		"visibility": 1,
		"flavorName": instanceType,
	}

	ret := []SImage{}
	for {
		resp, err := self.list(SERVICE_IMAGE, "/v4/image/list", params)
		if err != nil {
			return nil, err
		}
		part := struct {
			ReturnObj struct {
				Images []SImage
			}
			TotalCount int
		}{}
		err = resp.Unmarshal(&part)
		if err != nil {
			return nil, errors.Wrapf(err, "Unmarshal")
		}
		ret = append(ret, part.ReturnObj.Images...)
		if len(part.ReturnObj.Images) == 0 || len(ret) >= part.TotalCount {
			break
		}
		pageNo++
		params["pageNo"] = pageNo
	}
	var res []cloudprovider.ICloudImage
	for i := range ret {
		res = append(res, &ret[i])
	}
	return res, nil
}
