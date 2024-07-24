package aliyun

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

func (self *SAliyunClient) fetchRegionsCfel() error {
	if len(self.iregions) > 0 {
		return nil
	}
	body, err := self.ecsRequest("DescribeRegions", map[string]string{"AcceptLanguage": "zh-CN", "InstanceChargeType": "PostPaid"})
	if err != nil {
		return errors.Wrapf(err, "DescribeRegions")
	}

	regions := make([]SRegion, 0)
	err = body.Unmarshal(&regions, "Regions", "Region")
	if err != nil {
		return errors.Wrapf(err, "resp.Unmarshal")
	}
	self.iregions = make([]cloudprovider.ICloudRegion, len(regions))
	for i := 0; i < len(regions); i += 1 {
		regions[i].client = self
		self.iregions[i] = &regions[i]
	}
	return nil
}
