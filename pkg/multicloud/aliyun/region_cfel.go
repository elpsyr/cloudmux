package aliyun

import (
	"encoding/json"
	"fmt"
	alierr "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"strings"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
)

// Verify that *SRegion implements ICfelCloudRegion
var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

const (
	PREPAID      = "PrePaid"      // 包年包月
	POSTPAID     = "PostPaid"     // 按量付费
	SPOTPOSTPAID = "SpotPostPaid" // 抢占

	SpotStrategySpotAsPriceGo = "SpotAsPriceGo"
)

// SInstancePrice DescribePrice 接口返回
type SInstancePrice struct {
	RequestID string    `json:"RequestId"`
	PriceInfo PriceInfo `json:"PriceInfo"`
}
type DetailInfo struct {
	OriginalPrice float64 `json:"OriginalPrice"`
	DiscountPrice float64 `json:"DiscountPrice"`
	Resource      string  `json:"Resource"`
	TradePrice    float64 `json:"TradePrice"`
}
type DetailInfos struct {
	DetailInfo []DetailInfo `json:"DetailInfo"`
}
type Price struct {
	OriginalPrice             float64     `json:"OriginalPrice"`
	ReservedInstanceHourPrice float64     `json:"ReservedInstanceHourPrice"`
	DiscountPrice             float64     `json:"DiscountPrice"`
	Currency                  string      `json:"Currency"`
	DetailInfos               DetailInfos `json:"DetailInfos"`
	TradePrice                float64     `json:"TradePrice"`
}
type PriceInfoRules struct {
	Description string `json:"Description"`
	RuleID      int    `json:"RuleId"`
}
type PriceInfo struct {
	Price Price          `json:"Price"`
	Rules PriceInfoRules `json:"Rules"`
}

// GetDescribePrice 查询云服务器ECS资源的最新价格。
// 只获取实例规格价格，不包含系统盘
func (self *SRegion) GetDescribePrice(zoneID, InstanceType, paidType string) (*DetailInfo, error) {

	var spot bool
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
		params["SpotStrategy"] = SpotStrategySpotAsPriceGo // 系统自动出价，最高按量付费价格。
		params["ZoneId"] = zoneID                          // 抢占式实例不同可用区价格可能不同，查询抢占式实例价格时，建议传入ZoneId查询指定可用区的抢占式实例价格。
		spot = true
	default:
	}

	body, err := self.ecsRequest("DescribePrice", params)
	if err != nil {

		// import "github.com/pkg/errors"
		// errUnwrap := gerrors.Unwrap(err)
		// log.Errorf("Unwrap err %s", errUnwrap)
		if e, ok := errors.Cause(err).(*alierr.ServerError); ok {
			switch e.ErrorCode() {
			case "InvalidSystemDiskCategory.ValueNotSupported", "InvalidInstanceType.NotSupportDiskCategory":
				Category, err := self.GetInstanceTypeAvailableDiskType(paidType, zoneID, InstanceType, spot)
				if err != nil {
					fmt.Println("aliyun GetInstanceTypeAvailableDiskType fail params: ", err.Error())
				}
				if len(Category) > 0 {
					params["SystemDisk.Category"] = Category[0]
				} else { //售罄
					price := new(DetailInfo)
					price.TradePrice = -1
					return price, nil
				}
				body, err = self.ecsRequest("DescribePrice", params)
				if err != nil {
					log.Errorf("DescribePrice fail %s", err)
					return nil, err
				}
			case "PriceNotFound", "InvalidInstanceType.ValueNotSupported":
				// 部分 region 下查询 instanceType 出现 PriceNotFound ，返回空数据
				// 部分 instanceType 可按量付费无法包年包月
				price := new(DetailInfo)
				price.TradePrice = -1
				return price, nil
			default:
				// 打印出 请求参数
				jsonParams, _ := json.Marshal(params)
				log.Errorf("aliyun DescribePrice fail %s", err)
				log.Errorf("aliyun DescribePrice fail ,params: %s ", string(jsonParams))
				return nil, err
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

	// 获取并返回实例规格价格
	for _, detailInfo := range instancePrice.PriceInfo.Price.DetailInfos.DetailInfo {
		if detailInfo.Resource == "instanceType" {
			return &detailInfo, nil
		}
	}
	jsonParams, _ := json.Marshal(params)
	log.Errorf("aliyun DescribePrice no detailInfo ,fail params: %s", string(jsonParams))
	return nil, errors.Errorf("instanceType %s price detailInfo notfound ", InstanceType)
}

func (self *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetDescribePrice(zoneID, instanceType, SPOTPOSTPAID)
	if err != nil {
		return -1, err
	}
	return price.TradePrice, nil

}

func (self *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetDescribePrice(zoneID, instanceType, POSTPAID)
	if err != nil {
		return -1, err
	}
	return price.TradePrice, nil
}

func (self *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetDescribePrice(zoneID, instanceType, PREPAID)
	if err != nil {
		return -1, err
	}
	return price.TradePrice, nil
}

func (self *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	resource, err := self.GetAvailableResource(DestinationResourceInstanceType, POSTPAID, zoneID, instanceType, true)
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
	resource, err := self.GetAvailableResource(DestinationResourceInstanceType, POSTPAID, zoneID, instanceType, false)
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
	resource, err := self.GetAvailableResource(DestinationResourceInstanceType, PREPAID, zoneID, instanceType, false)
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

// GetICfelCloudImage 获取阿里云镜像
func (self *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	images := make([]SImage, 0)
	for {
		parts, total, err := self.GetImages(ImageStatusType(""), ImageOwnerSystem, nil, "", len(images), 50)
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

// GetInstanceMatchImage  获取阿里云实例规格可用镜像
func (self *SRegion) GetInstanceMatchImage(instancetype string) ([]cloudprovider.ICloudImage, error) {
	images := make([]SImage, 0)
	for {
		parts, total, err := self.GetImagesByInstanceType(instancetype, ImageStatusType(""), ImageOwnerSystem, nil, "", len(images), 50)
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

func (self *SRegion) GetImagesByInstanceType(instanceType string, status ImageStatusType, owner ImageOwnerType, imageId []string, name string, offset int, limit int) ([]SImage, int, error) {
	if limit > 50 || limit <= 0 {
		limit = 50
	}
	params := make(map[string]string)
	params["RegionId"] = self.RegionId
	params["PageSize"] = fmt.Sprintf("%d", limit)
	params["PageNumber"] = fmt.Sprintf("%d", (offset/limit)+1)

	if len(instanceType) > 0 {
		params["instanceType"] = instanceType
	}

	if len(status) > 0 {
		params["Status"] = string(status)
	} else {
		params["Status"] = "Creating,Available,UnAvailable,CreateFailed"
	}
	if imageId != nil && len(imageId) > 0 {
		params["ImageId"] = strings.Join(imageId, ",")
	}
	if len(owner) > 0 {
		params["ImageOwnerAlias"] = string(owner)
	}

	if len(name) > 0 {
		params["ImageName"] = name
	}

	return self.getImages(params)
}
