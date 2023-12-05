package aliyun

import (
	alierr "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
)

// Verify that *SRegion implements ICfelCloudRegion
var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

const (
	PREPAID      = "Prepaid"      // 包年包月
	POSTPAID     = "PostPaid"     // 按量付费
	SPOTPOSTPAID = "SpotPostPaid" // 抢占
)

// SInstancePrice DescribePrice 接口返回
type SInstancePrice struct {
	RequestID string `json:"RequestId"`
	PriceInfo struct {
		Price struct {
			OriginalPrice             int     `json:"OriginalPrice"`
			ReservedInstanceHourPrice int     `json:"ReservedInstanceHourPrice"`
			DiscountPrice             float64 `json:"DiscountPrice"`
			Currency                  string  `json:"Currency"`
			TradePrice                float64 `json:"TradePrice"`
		} `json:"Price"`
		Rules struct {
			Rule []struct {
				Description string `json:"Description"`
				RuleID      int    `json:"RuleId"`
			} `json:"Rule"`
		} `json:"Rules"`
	} `json:"PriceInfo"`
}

// / GetDescribePrice 查询云服务器ECS资源的最新价格。
// implement by cfel
func (self *SRegion) GetDescribePrice(zoneID, InstanceType, paidType string) (*SInstancePrice, error) {

	params := make(map[string]string)
	params["RegionId"] = self.RegionId
	params["ResourceType"] = "instance" // 目标资源的类型
	params["InstanceType"] = InstanceType
	switch paidType {
	case POSTPAID:
	case PREPAID:
		params["Period"] = "1"
		params["PriceUnit"] = "Month"
	case SPOTPOSTPAID:
		params["SpotStrategy"] = "SpotAsPriceGo" // 系统自动出价，最高按量付费价格。
		params["ZoneId"] = zoneID                // 抢占式实例不同可用区价格可能不同，查询抢占式实例价格时，建议传入ZoneId查询指定可用区的抢占式实例价格。
	default:
	}

	body, err := self.ecsRequest("DescribePrice", params)
	if err != nil {
		// import "github.com/pkg/errors"
		// errUnwrap := gerrors.Unwrap(err)
		// log.Errorf("Unwrap err %s", errUnwrap)
		if e, ok := errors.Cause(err).(*alierr.ServerError); ok {
			switch e.ErrorCode() {
			case "InvalidSystemDiskCategory.ValueNotSupported":
				params["SystemDisk.Category"] = "cloud_essd"
				body, err = self.ecsRequest("DescribePrice", params)
				if err != nil {
					log.Errorf("DescribePrice fail %s", err)
					return nil, err
				}
			case "PriceNotFound", "InvalidInstanceType.ValueNotSupported":
				// 部分 region 下查询 instanceType 出现 PriceNotFound ，返回空数据
				// 部分 instanceType 可按量付费无法包年包月
				price := new(SInstancePrice)
				price.PriceInfo.Price.TradePrice = -1
				return price, nil
			}

		} else {
			log.Errorf("DescribePrice fail %s", err)
			return nil, err
		}

	}

	instancePrice := new(SInstancePrice)
	err = body.Unmarshal(&instancePrice)
	if err != nil {
		log.Errorf("Unmarshal available resources details fail %s", err)
		return nil, err
	}
	if instancePrice == nil {
		return nil, errors.Errorf("body.Unmarshal err: return nil")
	}

	return instancePrice, nil
}

func (self *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetDescribePrice(zoneID, instanceType, SPOTPOSTPAID)
	if err != nil {
		return 0, err
	}
	return price.PriceInfo.Price.TradePrice, nil

}

func (self *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetDescribePrice(zoneID, instanceType, POSTPAID)
	if err != nil {
		return 0, err
	}
	return price.PriceInfo.Price.TradePrice, nil
}

func (self *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetDescribePrice(zoneID, instanceType, PREPAID)
	if err != nil {
		return 0, err
	}
	return price.PriceInfo.Price.TradePrice, nil
}

func (self *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	resource, err := self.GetAvailableResource("PostPaid", zoneID, instanceType, true)
	if err != nil {
		return "", err
	}
	availableZone := resource.AvailableZones.AvailableZone
	if availableZone != nil && len(availableZone) > 0 {
		availableResource := availableZone[0].AvailableResources.AvailableResource
		if availableResource != nil && len(availableResource) > 0 {
			supportedResource := availableResource[0].SupportedResources.SupportedResource
			if supportedResource != nil && len(supportedResource) > 0 {
				status := supportedResource[0].Status
				if status == AliyunResourceAvailable {
					return api.SkuStatusAvailable, nil
				} else if status == AliyunResourceSoldOut {
					return api.SkuStatusSoldout, nil
				}
				return api.SkuStatusSoldout, nil
			}
		}
	}
	return api.SkuStatusSoldout, nil

}

func (self *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	resource, err := self.GetAvailableResource("PostPaid", zoneID, instanceType, false)
	if err != nil {
		return "", err
	}
	availableZone := resource.AvailableZones.AvailableZone
	if availableZone != nil && len(availableZone) > 0 {
		availableResource := availableZone[0].AvailableResources.AvailableResource
		if availableResource != nil && len(availableResource) > 0 {
			supportedResource := availableResource[0].SupportedResources.SupportedResource
			if supportedResource != nil && len(supportedResource) > 0 {
				status := supportedResource[0].Status
				if status == AliyunResourceAvailable {
					return api.SkuStatusAvailable, nil
				} else if status == AliyunResourceSoldOut {
					return api.SkuStatusSoldout, nil
				}
				return api.SkuStatusSoldout, nil
			}
		}
	}
	return api.SkuStatusSoldout, nil
}

func (self *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	resource, err := self.GetAvailableResource("PrePaid", zoneID, instanceType, false)
	if err != nil {
		return "", err
	}
	availableZone := resource.AvailableZones.AvailableZone
	if availableZone != nil && len(availableZone) > 0 {
		availableResource := availableZone[0].AvailableResources.AvailableResource
		if availableResource != nil && len(availableResource) > 0 {
			supportedResource := availableResource[0].SupportedResources.SupportedResource
			if supportedResource != nil && len(supportedResource) > 0 {
				status := supportedResource[0].Status
				if status == AliyunResourceAvailable {
					return api.SkuStatusAvailable, nil
				} else if status == AliyunResourceSoldOut {
					return api.SkuStatusSoldout, nil
				}
				return api.SkuStatusSoldout, nil
			}
		}
	}
	return api.SkuStatusSoldout, nil
}

// UpdateVpc 更新VPC名称以及描述
func (self *SRegion) UpdateVpc(opts *cloudprovider.VpcUpdateOptions) error {
	params := make(map[string]string)
	params["VpcId"] = opts.ID
	if opts.NAME != "" {
		params["VpcName"] = opts.NAME
	}
	if opts.Desc != "" {
		params["Description"] = opts.Desc
	}
	_, err := self.ecsRequest("ModifyVpcAttribute", params)
	return err
}
