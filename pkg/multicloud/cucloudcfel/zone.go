package cucloudcfel

import (
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SZone struct {
	multicloud.SResourceBase
	CuCloudTags
	region *SRegion
	// host   *SHost

	istorages []cloudprovider.ICloudStorage

	AzID            string `json:"azId"`
	AzName          string `json:"azName"`
	CloudRegionID   string `json:"cloudRegionId"`
	CloudRegionName string `json:"cloudRegionName"`
	PubStatus       string `json:"pubStatus"`
	RegionCode      string `json:"regionCode"`
	RegionID        string `json:"regionId"`
	RegionName      string `json:"regionName"`
	Status          string `json:"status"`
	ZoneCode        string `json:"zoneCode"`
	ZoneID          string `json:"zoneId"`
	ZoneName        string `json:"zoneName"`
}

var _ cloudprovider.ICloudZone = (*SZone)(nil)

// GetGlobalId implements cloudprovider.ICloudZone.
func (s *SZone) GetGlobalId() string {
	return fmt.Sprintf("%s/%s", s.region.GetGlobalId(), s.ZoneCode)
}

// GetI18n implements cloudprovider.ICloudZone.
func (s *SZone) GetI18n() cloudprovider.SModelI18nTable {
	table := cloudprovider.SModelI18nTable{}
	table["name"] = cloudprovider.NewSModelI18nEntry(s.GetName()).CN(s.GetName()).EN(fmt.Sprintf("%s %s", CLOUD_PROVIDER_CUCLOUD, s.ZoneName))
	return table
}

// GetIHostById implements cloudprovider.ICloudZone.
func (s *SZone) GetIHostById(id string) (cloudprovider.ICloudHost, error) {
	panic("unimplemented")
}

// GetIHosts implements cloudprovider.ICloudZone.
func (s *SZone) GetIHosts() ([]cloudprovider.ICloudHost, error) {
	panic("unimplemented")
}

// GetIRegion implements cloudprovider.ICloudZone.
func (s *SZone) GetIRegion() cloudprovider.ICloudRegion {
	return s.region
}

// GetIStorageById implements cloudprovider.ICloudZone.
func (s *SZone) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	return nil, nil
}

// GetIStorages implements cloudprovider.ICloudZone.
func (s *SZone) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	if s.istorages == nil {
		res, err := s.fetchDiskType()
		if err != nil {
			return nil, err
		}
		var result []cloudprovider.ICloudStorage

		for i := range res {
			res[i].zone = s
			result = append(result, &res[i])
		}
		s.istorages = result
	}
	return s.istorages, nil
}

func (s *SZone) fetchDiskType() ([]SStorage, error) {
	params := map[string]interface{}{
		"pageSize":      10000,
		"specilZone":    s.RegionCode,
		"billType":      "0",
		"serviceTypeId": "1002",
	}

	res, err := s.region.client.postWithToken("product-cons/product/productquery/flavors", params)
	if err != nil {
		return nil, err
	}
	var storage = struct {
		List []SStorage `json:"list"`
	}{}
	if err = res.Unmarshal(&storage, "result"); err != nil {
		return nil, err
	}

	return storage.List, nil
}

// GetId implements cloudprovider.ICloudZone.
func (s *SZone) GetId() string {
	return s.GetGlobalId()
}

// GetName implements cloudprovider.ICloudZone.
func (s *SZone) GetName() string {
	return s.ZoneName
}

// GetStatus implements cloudprovider.ICloudZone.
func (s *SZone) GetStatus() string {
	return "enable"
}
