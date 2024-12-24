package ctyun

import (
	"sync"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (self *SCtyunClient) CfelGetRegions() ([]SRegion, error) {
	resp, err := self.list(SERVICE_ECS, "/v4/region/list-regions", nil)
	if err != nil {
		return nil, err
	}
	ret := struct {
		ReturnObj struct {
			RegionList []SRegion
		}
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	var ch = make(chan struct{}, 10)
	defer close(ch)
	for i := range ret.ReturnObj.RegionList {
		wg.Add(1)
		ch <- struct{}{}
		go func(region *SRegion) {
			defer func() {
				<-ch
				wg.Done()
			}()
			region.client = self
			product, err := region.getProduct()
			if err != nil {
				region.Stat = api.CLOUD_REGION_STATUS_OUTOFSERVICE
				return
			}
			if len(product.Other.Region) == 0 {
				region.Stat = api.CLOUD_REGION_STATUS_OUTOFSERVICE
				return
			}
			region.Stat = api.CLOUD_REGION_STATUS_INSERVER
		}(&ret.ReturnObj.RegionList[i])
	}
	wg.Wait()

	return ret.ReturnObj.RegionList, nil
}

func (self *SCtyunClient) GetIRegions() ([]cloudprovider.ICloudRegion, error) {
	ret := []cloudprovider.ICloudRegion{}
	regions, err := self.CfelGetRegions()
	if err != nil {
		return nil, err
	}
	for i := range regions {
		if regions[i].Stat != api.CLOUD_REGION_STATUS_INSERVER {
			continue
		}
		ret = append(ret, &self.regions[i])
	}
	return ret, nil
}
