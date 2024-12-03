package baiducfel

import (
	"errors"
	"fmt"
	"slices"

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

var duration = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}

type specPrice struct {
	SpecID     string  `json:"specId"`
	SpecPrices []price `json:"specPrices"`
}

type price struct {
	Spec       string  `json:"spec"`
	SpecPrice  float32 `json:"specPrice"`
	Status     string  `json:"status"`
	TradePrice float32 `json:"tradePrice"`
}

type volumePrice struct {
	CdsSizeInGB int     `json:"cdsSizeInGB"`
	Price       float32 `json:"price"`
	StorageType string  `json:"storageType"`
	Unit        string  `json:"unit"`
}

func (s *SRegion) GetICfelSkuPrice(opt *cloudprovider.CfelSkuPriceOptions) (map[string]string, error) {
	var serverTotalPrice, volumeTotalPrice float32
	var paymentTiming = "Postpaid" //包括Postpaid(后付费)，Prepaid(预付费)两种
	var length = opt.Duration
	if opt.FeeUnit == "month" {
		paymentTiming = "Prepaid"
	} else if opt.FeeUnit == "year" {
		paymentTiming = "Prepaid"
		length = opt.Duration * 12
	}
	if !slices.Contains(duration, length) {
		return nil, errors.New("only support Month and duration in [1,2,3,4,5,6,7,8,9,12,24,36]")
	}

	if len(opt.InstanceType) > 0 {
		var params = map[string]interface{}{
			"purchaseCount":  opt.Quantity,
			"purchaseLength": length,
			"spec":           opt.InstanceType,
			"paymentTiming":  paymentTiming,
			"zoneName":       opt.ZoneId,
		}
		res, err := s.doPost(ServiceInstance, "/v2/instance/price", params)
		if err != nil {
			return nil, err
		}
		var r []specPrice
		if err = res.Unmarshal(&r, "price"); err != nil {
			return nil, err
		}
		if len(r) == 0 {
			return nil, errors.New("price is empty")
		}
		for _, val := range r[0].SpecPrices {
			if val.Spec == opt.InstanceType && val.Status == "available" {
				serverTotalPrice = val.TradePrice
				break
			}
		}
	} else { // 云盘询价
		var params = map[string]interface{}{
			"purchaseCount":  opt.Quantity,
			"purchaseLength": length,
			"storageType":    opt.SysDiskType,
			"cdsSizeInGB":    opt.SysDiskSize,
			"paymentTiming":  paymentTiming,
			"zoneName":       opt.ZoneId,
		}
		res, err := s.doPost(ServiceInstance, "/v2/volume/getPrice", params)
		if err != nil {
			return nil, err
		}
		var r []volumePrice
		if err = res.Unmarshal(&r, "price"); err != nil {
			return nil, err
		}
		if len(r) == 0 {
			return nil, errors.New("price is empty")
		}
		volumeTotalPrice = r[0].Price
	}
	var result = make(map[string]string)
	if opt.InstanceType == "" {
		result["dataVolumePrice"] = fmt.Sprintf("%v", volumeTotalPrice)
	} else {
		// result["bootVolumePrice"] = fmt.Sprintf("%v", volumeTotalPrice)
		result["serverPrice"] = fmt.Sprintf("%v", serverTotalPrice)
	}

	return result, nil
}
