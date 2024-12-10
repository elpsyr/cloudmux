package baiducfel

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type SCfelRegion struct {
	iZones []cloudprovider.ICloudZone
}

var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

type imgResp struct {
	NextMarker  string   `json:"nextMarker"`
	Marker      string   `json:"marker"`
	MaxKeys     int      `json:"maxKeys"`
	IsTruncated bool     `json:"isTruncated"`
	Images      []SImage `json:"images"`
}

func (r *SRegion) GetICfelCloudImageById(id string) (cloudprovider.ICloudImage, error) {
	return nil, nil
}

func (self *SRegion) GetInstanceMatchImage(instanceType string) ([]cloudprovider.ICloudImage, error) {
	var query = map[string]string{
		"spec":    instanceType,
		"maxKeys": "100",
	}
	var imgs []SImage

	var marker string
	for {
		if len(marker) > 0 {
			query["marker"] = marker
		}
		res, err := self.doList(ServiceImage, "/v2/image/getAvailableImagesBySpec", query)
		if err != nil {
			return nil, err
		}
		var r imgResp
		err = res.Unmarshal(&r)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, r.Images...)
		if !r.IsTruncated {
			break
		}
		marker = r.NextMarker
	}

	var ret []cloudprovider.ICloudImage
	for i := range imgs {
		ret = append(ret, &imgs[i])
	}
	return ret, nil
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
	body, err := region.client.list("bcc", region.Region, "/v2/zone", nil, nil)
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


