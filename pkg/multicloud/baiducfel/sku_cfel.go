package baiducfel

import (
	"fmt"
	"slices"

	"yunion.io/x/pkg/errors"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	skus, err := self.fetchFlavorSpec()
	if err != nil {
		return nil, errors.Wrapf(err, "fetchFlavorSpec")
	}
	var ret []cloudprovider.ICfelCloudSku
	for i := range skus {
		ret = append(ret, &skus[i])
	}
	return ret, nil
}

var duration = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}

type specPrice struct {
	SpecID     string  `json:"specId"`
	SpecPrices []price `json:"specPrices"`
}

type price struct {
	Spec       string  `json:"spec"`
	SpecPrice  float64 `json:"specPrice"`
	Status     string  `json:"status"`
	TradePrice float64 `json:"tradePrice"`
}

type volumePrice struct {
	CdsSizeInGB int     `json:"cdsSizeInGB"`
	Price       float64 `json:"price"`
	StorageType string  `json:"storageType"`
	Unit        string  `json:"unit"`
}

func (s *SRegion) GetICfelSkuPrice(opt *cloudprovider.CfelSkuPriceOptions) (map[string]string, error) {
	var serverTotalPrice, volumeTotalPrice float64
	// var paymentTiming = opt.ChargeType //包括Postpaid(后付费)，Prepaid(预付费)两种
	var length = opt.Duration
	if opt.ChargeType == cloudprovider.InstanceChargeTypePrePaid && opt.FeeUnit == "year" {
		length = opt.Duration * 12
	}

	if !slices.Contains(duration, length) {
		return nil, errors.Errorf("only support Month and duration in [1,2,3,4,5,6,7,8,9,12,24,36]")
	}

	if len(opt.InstanceType) > 0 {
		if opt.ChargeType == cloudprovider.InstanceChargeTypeSpotPaid {
			var err error
			serverTotalPrice, err = s.GetSpotPrice(opt.ZoneId, opt.InstanceType, opt.DiskType, opt.DiskSize)
			if err != nil {
				return nil, err
			}
		} else {
			var params = map[string]interface{}{
				"purchaseCount":  opt.Quantity,
				"purchaseLength": length,
				"spec":           opt.InstanceType,
				"paymentTiming":  opt.ChargeType,
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
				return nil, errors.Errorf("price is empty")
			}
			for _, val := range r[0].SpecPrices {
				if val.Spec == opt.InstanceType && val.Status == "available" {
					serverTotalPrice = val.TradePrice
					break
				}
			}
			if opt.ChargeType == cloudprovider.InstanceChargeTypePostPaid {
				serverTotalPrice = 60 * serverTotalPrice // 返回为每分钟价格，小时价*60
			}
		}
	} else { // 云盘询价
		var params = map[string]interface{}{
			"purchaseCount":  opt.Quantity,
			"purchaseLength": length,
			"storageType":    opt.DiskType,
			"cdsSizeInGB":    opt.DiskSize,
			"paymentTiming":  opt.ChargeType,
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
			return nil, errors.Errorf("price is empty")
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
