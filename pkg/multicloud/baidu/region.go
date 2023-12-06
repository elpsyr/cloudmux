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

package baidu

import (
	"fmt"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

// https://cloud.baidu.com/doc/BOS/s/akrqd2wcx
var regions = map[string]string{
	"bj":  "北京",
	"gz":  "广州",
	"su":  "苏州",
	"hkg": "香港",
	"fwh": "武汉",
	"bd":  "保定",
	"sin": "新加坡",
	"fsh": "上海",
}

type SRegion struct {
	multicloud.SRegion
	multicloud.SNoObjectStorageRegion
	multicloud.SNoLbRegion
	client *SBaiduClient

	Region     string
	RegionName string
	iZones     []cloudprovider.ICloudZone
}

func (region *SRegion) GetId() string {
	return region.Region
}

func (region *SRegion) GetGlobalId() string {
	return fmt.Sprintf("%s/%s", api.CLOUD_PROVIDER_BAIDU, region.Region)
}

func (region *SRegion) GetProvider() string {
	return api.CLOUD_PROVIDER_BAIDU
}

func (region *SRegion) GetCloudEnv() string {
	return api.CLOUD_PROVIDER_BAIDU
}

func (region *SRegion) GetGeographicInfo() cloudprovider.SGeographicInfo {
	geo, ok := map[string]cloudprovider.SGeographicInfo{
		"bj":  api.RegionBeijing,
		"gz":  api.RegionGuangzhou,
		"su":  api.RegionSuzhou,
		"hkg": api.RegionHongkong,
		"fwh": api.RegionHangzhou,
		"bd":  api.RegionBaoDing,
		"sin": api.RegionSingapore,
		"fsh": api.RegionShanghai,
	}[region.Region]
	if ok {
		return geo
	}
	return cloudprovider.SGeographicInfo{}
}

func (region *SRegion) GetName() string {
	return region.RegionName
}

func (region *SRegion) GetI18n() cloudprovider.SModelI18nTable {
	table := cloudprovider.SModelI18nTable{}
	table["name"] = cloudprovider.NewSModelI18nEntry(region.GetName()).CN(region.GetName()).EN(region.Region)
	return table
}

func (region *SRegion) GetStatus() string {
	return api.CLOUD_REGION_STATUS_INSERVER
}

func (region *SRegion) GetClient() *SBaiduClient {
	return region.client
}

func (region *SRegion) CreateEIP(opts *cloudprovider.SEip) (cloudprovider.ICloudEIP, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) CreateISecurityGroup(conf *cloudprovider.SecurityGroupCreateInput) (cloudprovider.ICloudSecurityGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetISecurityGroupById(secgroupId string) (cloudprovider.ICloudSecurityGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) CreateIVpc(opts *cloudprovider.VpcCreateOptions) (cloudprovider.ICloudVpc, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetIVpcs() ([]cloudprovider.ICloudVpc, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetIVpcById(id string) (cloudprovider.ICloudVpc, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetCapabilities() []string {
	return region.client.GetCapabilities()
}

func (region *SRegion) GetIEipById(eipId string) (cloudprovider.ICloudEIP, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetIEips() ([]cloudprovider.ICloudEIP, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetIZones() ([]cloudprovider.ICloudZone, error) {
	if region.iZones == nil {
		var err error
		err = region.fetchInfrastructure()
		if err != nil {
			return nil, err
		}
	}
	return region.iZones, nil
}

func (region *SRegion) fetchInfrastructure() error {
	err := region._fetchZones()
	if err != nil {
		return err
	}
	return nil
}

// https://cloud.baidu.com/doc/BCC/s/ijwvyo9im
func (region *SRegion) _fetchZones() error {
	body, err := region.client.list("bcc", region.Region, "/v2/zone", nil)
	if err != nil {
		return err
	}

	zones := make([]SZone, 0)
	err = body.Unmarshal(&zones, "zones")
	if err != nil {
		return err
	}

	region.iZones = make([]cloudprovider.ICloudZone, len(zones))

	for i := 0; i < len(zones); i += 1 {
		zones[i].region = region
		region.iZones[i] = &zones[i]
	}
	return nil
}

func (region *SRegion) GetIZoneById(id string) (cloudprovider.ICloudZone, error) {
	return nil, cloudprovider.ErrNotImplemented
}
