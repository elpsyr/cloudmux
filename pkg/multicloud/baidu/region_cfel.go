package baidu

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type SCfelRegion struct {
	iZones []cloudprovider.ICloudZone
}

func (region *SRegion) GetICfelZones() ([]cloudprovider.ICloudZone, error) {
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
