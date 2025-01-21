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

package cucloudcfel

import (
	"fmt"
	"sync"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SRegion struct {
	multicloud.SRegion
	multicloud.SNoObjectStorageRegion
	multicloud.SNoLbRegion
	client *SChinaUnionClient

	zones []cloudprovider.ICloudZone

	CloudRegionId   string
	CloudRegionName string
	CloudRegionCode string
	Status          string

	mut sync.Mutex
	postPaidSku *SServerSku

	mut1 sync.Mutex
	prePaidSku *SServerSku
}

func (self *SRegion) GetId() string {
	return self.CloudRegionCode
}

func (self *SRegion) GetGlobalId() string {
	return fmt.Sprintf("%s/%s", CLOUD_PROVIDER_CUCLOUD, self.CloudRegionCode)
}

func (self *SRegion) GetProvider() string {
	return api.CLOUD_PROVIDER_CUCLOUD
}

func (self *SRegion) GetCloudEnv() string {
	return api.CLOUD_PROVIDER_CUCLOUD
}

func (self *SRegion) GetGeographicInfo() cloudprovider.SGeographicInfo {
	geo, ok := map[string]cloudprovider.SGeographicInfo{
		"cn-chongqing-1": api.RegionChongqing,
		"cn-wuhan-2":     api.RegionWuhan,
		"cn-guangzhou-1": api.RegionGuangzhou,
		"cn-beijing-1":   api.RegionBeijing,
	}[self.CloudRegionCode]
	if ok {
		return geo
	}
	return cloudprovider.SGeographicInfo{}
}

func (self *SRegion) GetName() string {
	return self.CloudRegionName
}

func (self *SRegion) GetI18n() cloudprovider.SModelI18nTable {
	table := cloudprovider.SModelI18nTable{}
	table["name"] = cloudprovider.NewSModelI18nEntry(self.GetName()).CN(self.GetName()).EN(self.CloudRegionName)
	return table
}

func (self *SRegion) GetStatus() string {
	return api.CLOUD_REGION_STATUS_INSERVER
}

func (self *SRegion) GetClient() *SChinaUnionClient {
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

func (self *SRegion) GetIZones() ([]cloudprovider.ICloudZone, error) {
	if self.zones == nil {
		err := self.fetchZones()
		if err != nil {
			return nil, err
		}
	}
	return self.zones, nil
}

func (self *SRegion) GetIZoneById(id string) (cloudprovider.ICloudZone, error) {
	if self.zones == nil {
		err := self.fetchZones()
		if err != nil {
			return nil, err
		}
	}
	for i := range self.zones {
		if self.zones[i].GetGlobalId() == id {
			return self.zones[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) fetchZones() error {
	params := map[string]interface{}{
		"cloudRegionCode": self.CloudRegionCode,
	}
	res, err := self.client.list("/instance/v1/product/zones", params)
	if err != nil {
		return err
	}

	var rr = []SZone{}
	// var zones []SZone
	if err = res.Unmarshal(&rr, "list"); err != nil {
		return err
	}
	var ret []cloudprovider.ICloudZone
	for i := range rr {
		if rr[i].RegionCode != self.CloudRegionCode {
			continue
		}
		rr[i].region = self
		ret = append(ret, &rr[i])
	}
	self.zones = ret
	return nil
}
