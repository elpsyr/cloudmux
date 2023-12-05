package qcloud

import (
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

// Verify that *SRegion implements ICfelCloudRegion
var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

const (
	PrePaid        = "PREPAID"          // 预付费
	PostPaidByHour = "POSTPAID_BY_HOUR" // 按量付费
	SpotPaid       = "SPOTPAID"         // 抢占付费

	StatusSoldOut = "SOLD_OUT"
	StatusSell    = "SELL"
)

// GetSpotPostPaidPrice 获取 zone 下 对应 instanceType 抢占付费价格
func (region *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := region.GetInstanceTypesPrice(zoneID, instanceType)
	if err != nil {
		return 0, err
	}
	for _, instance := range price.InstanceTypeQuotaSet {
		if instance.InstanceChargeType == SpotPaid {
			return instance.Price.UnitPriceDiscount, err
		}
		continue
	}
	// 如果没有抢占价格，则以按量付费作为价格展示
	for _, instance := range price.InstanceTypeQuotaSet {
		if instance.InstanceChargeType == PostPaidByHour {
			return instance.Price.UnitPriceDiscount, err
		}
		continue
	}
	return 0, err
}

func (region *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := region.GetInstanceTypesPrice(zoneID, instanceType)
	if err != nil {
		return 0, err
	}
	for _, instance := range price.InstanceTypeQuotaSet {
		if instance.InstanceChargeType == PostPaidByHour {
			return instance.Price.UnitPriceDiscount, err
		}
		continue
	}
	return 0, err
}

func (region *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := region.GetInstanceTypesPrice(zoneID, instanceType)
	if err != nil {
		return 0, err
	}
	for _, instance := range price.InstanceTypeQuotaSet {
		if instance.InstanceChargeType == PrePaid {
			return instance.Price.DiscountPrice, err
		}
		continue
	}
	return 0, err
}

// GetSpotPostPaidStatus 抢占付费是否可购买
// "SPOTPAID"
func (region *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	price, err := region.GetInstanceTypesPrice(zoneID, instanceType)
	if err != nil {
		return api.SkuStatusSoldout, err
	}

	for _, instance := range price.InstanceTypeQuotaSet {
		if instance.InstanceChargeType == SpotPaid {
			switch instance.Status {
			case StatusSoldOut:
				return api.SkuStatusSoldout, err
			case StatusSell:
				return api.SkuStatusAvailable, err
			default:
				return api.SkuStatusSoldout, err
			}
		}
		continue
	}
	return api.SkuStatusSoldout, nil
}

// GetPostPaidStatus 按量付费是否可购买
// "POSTPAID_BY_HOUR"
func (region *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	price, err := region.GetInstanceTypesPrice(zoneID, instanceType)
	if err != nil {
		return api.SkuStatusSoldout, err
	}

	for _, instance := range price.InstanceTypeQuotaSet {
		if instance.InstanceChargeType == PostPaidByHour {
			switch instance.Status {
			case StatusSoldOut:
				return api.SkuStatusSoldout, err
			case StatusSell:
				return api.SkuStatusAvailable, err
			default:
				return api.SkuStatusSoldout, err
			}
		}
		continue
	}
	return api.SkuStatusSoldout, nil
}

// GetPrePaidStatus 包年包月是否可购买
// "PREPAID" 单位为 MONTH
func (region *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	price, err := region.GetInstanceTypesPrice(zoneID, instanceType)
	if err != nil {
		return api.SkuStatusSoldout, err
	}

	for _, instance := range price.InstanceTypeQuotaSet {
		if instance.InstanceChargeType == PrePaid {
			switch instance.Status {
			case StatusSoldOut:
				return api.SkuStatusSoldout, err
			case StatusSell:
				return api.SkuStatusAvailable, err
			default:
				return api.SkuStatusSoldout, err
			}
		}
		continue
	}
	return api.SkuStatusSoldout, nil
}
