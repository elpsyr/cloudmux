package qcloud

import (
	"fmt"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
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

// GetICfelCloudImage 获取腾讯云镜像
func (self *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	images := make([]SImage, 0)
	for {
		parts, total, err := self.GetImages("", "", nil, "", len(images), 50)
		if err != nil {
			return nil, errors.Wrapf(err, "GetImages")
		}
		images = append(images, parts...)
		if len(images) >= total {
			break
		}
	}
	ret := []cloudprovider.ICloudImage{}
	for i := range images {
		images[i].storageCache = self.getStoragecache()
		ret = append(ret, &images[i])
	}
	return ret, nil
}

// GetInstanceMatchImage  获取腾讯云实例规格可用镜像
func (self *SRegion) GetInstanceMatchImage(instancetype string) ([]cloudprovider.ICloudImage, error) {
	images := make([]SImage, 0)
	for {
		parts, total, err := self.GetImagesByInstanceType(instancetype, "", "", nil, "", len(images), 50)
		if err != nil {
			return nil, errors.Wrapf(err, "GetImagesByInstanceType")
		}
		images = append(images, parts...)
		if len(images) >= total {
			break
		}
	}
	ret := []cloudprovider.ICloudImage{}
	for i := range images {
		images[i].storageCache = self.getStoragecache()
		ret = append(ret, &images[i])
	}
	return ret, nil
}

func (self *SRegion) GetImagesByInstanceType(instancetype, status string, owner string, imageIds []string, name string, offset int, limit int) ([]SImage, int, error) {
	if limit > 50 || limit <= 0 {
		limit = 50
	}
	params := make(map[string]string)
	params["Limit"] = fmt.Sprintf("%d", limit)
	params["Offset"] = fmt.Sprintf("%d", offset)
	params["InstanceType"] = instancetype

	for index, imageId := range imageIds {
		params[fmt.Sprintf("ImageIds.%d", index)] = imageId
	}

	if len(imageIds) == 0 { // imageIds 不能和Filter同时查询
		filter := 0
		if len(status) > 0 {
			params[fmt.Sprintf("Filters.%d.Name", filter)] = "image-state"
			params[fmt.Sprintf("Filters.%d.Values.0", filter)] = status
			filter++
		}

		if len(owner) > 0 {
			params[fmt.Sprintf("Filters.%d.Name", filter)] = "image-type"
			params[fmt.Sprintf("Filters.%d.Values.0", filter)] = owner
			filter++
		}

		if len(name) > 0 {
			params[fmt.Sprintf("Filters.%d.Name", filter)] = "image-name"
			params[fmt.Sprintf("Filters.%d.Values.0", filter)] = name
			filter++
		}
	}

	images := make([]SImage, 0)
	body, err := self.cvmRequest("DescribeImages", params, true)
	if err != nil {
		return nil, 0, err
	}
	err = body.Unmarshal(&images, "ImageSet")
	if err != nil {
		return nil, 0, err
	}
	for i := 0; i < len(images); i++ {
		images[i].storageCache = self.getStoragecache()
	}
	total, _ := body.Float("TotalCount")
	return images, int(total), nil
}
