// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package baiducfel

import (
	"fmt"
	"sync"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
)

var regions = map[string]string{
	"bj":  "北京",
	"gz":  "广州",
	"su":  "苏州",
	"hkg": "香港",
	"fwh": "武汉",
	"bd":  "保定",
	"fsh": "上海",
}

type SRegion struct {
	multicloud.SRegion
	multicloud.SNoObjectStorageRegion
	multicloud.SNoLbRegion
	client *SBaiduClient

	Region     string
	RegionName string
	SCfelRegion

	instanceType []Sku
	mut          sync.Mutex
}

func (self *SRegion) GetId() string {
	return self.Region
}

func (self *SRegion) GetGlobalId() string {
	return fmt.Sprintf("%s/%s", api.CLOUD_PROVIDER_BAIDU, self.Region)
}

func (self *SRegion) GetProvider() string {
	return api.CLOUD_PROVIDER_BAIDU
}

func (self *SRegion) GetCloudEnv() string {
	return api.CLOUD_PROVIDER_BAIDU
}

func (self *SRegion) GetGeographicInfo() cloudprovider.SGeographicInfo {
	geo, ok := map[string]cloudprovider.SGeographicInfo{
		"bj":  api.RegionBeijing,
		"gz":  api.RegionGuangzhou,
		"su":  api.RegionSuzhou,
		"hkg": api.RegionHongkong,
		"fwh": api.RegionHangzhou,
		"bd":  api.RegionBaoDing,
		"sin": api.RegionSingapore,
		"fsh": api.RegionShanghai,
	}[self.Region]
	if ok {
		return geo
	}
	return cloudprovider.SGeographicInfo{}
}

func (self *SRegion) GetName() string {
	return self.RegionName
}

func (self *SRegion) GetI18n() cloudprovider.SModelI18nTable {
	table := cloudprovider.SModelI18nTable{}
	table["name"] = cloudprovider.NewSModelI18nEntry(self.GetName()).CN(self.GetName()).EN(self.Region)
	return table
}

func (self *SRegion) GetStatus() string {
	return api.CLOUD_REGION_STATUS_INSERVER
}

func (self *SRegion) GetClient() *SBaiduClient {
	return self.client
}

func (self *SRegion) CreateEIP(opts *cloudprovider.SEip) (cloudprovider.ICloudEIP, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetCapabilities() []string {
	return region.client.GetCapabilities()
}

func (self *SRegion) GetIEipById(eipId string) (cloudprovider.ICloudEIP, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetIEips() ([]cloudprovider.ICloudEIP, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetIHostById(id string) (cloudprovider.ICloudHost, error) {
	zone, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := range zone {
		host, err := zone[i].GetIHostById(id)
		if err == nil {
			return host, nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) GetIVMById(id string) (cloudprovider.ICloudVM, error) {
	var ret SInstance
	err := self.doGet(ServiceInstance, "/v2/instance/"+id, nil, &ret)
	if err != nil {
		return nil, err
	}
	ret.region = self
	return &ret, nil
}

func (r *SRegion) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	istores, err := r.GetIStorages()
	if err != nil {
		return nil, err
	}
	for i := range istores {
		if istores[i].GetGlobalId() == id {
			return istores[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (r *SRegion) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	iStores := make([]cloudprovider.ICloudStorage, 0)

	izones, err := r.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		iZoneStores, err := izones[i].GetIStorages()
		if err != nil {
			return nil, err
		}
		iStores = append(iStores, iZoneStores...)
	}
	return iStores, nil
}

func (self *SRegion) GetIZones() ([]cloudprovider.ICloudZone, error) {
	return self.GetICfelZones()
}

func (self *SRegion) GetIZoneById(id string) (cloudprovider.ICloudZone, error) {
	if self.iZones == nil {
		self._fetchZones()
	}
	for i := range self.iZones {
		if self.iZones[i].GetGlobalId() == id {
			return self.iZones[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) doGet(service, resource string, query map[string]string, ret interface{}) error {
	res, err := self.client.request("GET", service, self.Region, resource, query, nil)
	if err != nil {
		return err
	}
	return res.Unmarshal(ret, service[:len(service)-1])
}

func (self *SRegion) doGetWithoutVal(service, resource string, query map[string]string) (jsonutils.JSONObject, error) {
	res, err := self.client.request("GET", service, self.Region, resource, query, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (self *SRegion) doList(service, resource string, query map[string]string) (jsonutils.JSONObject, error) {
	return self.client.request("GET", service, self.Region, resource, query, nil)
}

func (self *SRegion) doList1(service, resource string, query map[string]string, val interface{}) ([]interface{}, error) {
	if query == nil {
		query = make(map[string]string, 0)
	}
	var result []interface{}
	for {
		res, err := self.client.request("GET", service, self.Region, resource, query, nil)
		if err != nil {
			return nil, err
		}
		var ret = make(map[string]interface{})
		err = res.Unmarshal(&ret)
		if err != nil {
			return nil, err
		}
		res.Unmarshal(&val, service)

		result = append(result, val)
		if isTruncate, ok := ret["isTruncated"]; ok && isTruncate.(bool) {
			continue
		} else {
			break
		}
	}
	return result, nil
}

func (self *SRegion) doPost(service, resource string, param map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.client.request("POST", service, self.Region, resource, nil, param)
}

func (self *SRegion) doDelete(service, resource string) error {
	_, err := self.client.request("DELETE", service, self.Region, resource, nil, nil)
	return err
}

func (self *SRegion) doPut(service, resource string, param map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.client.request("PUT", service, self.Region, resource, nil, param)
}
