package cloudpods

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (self *SRegion) GetImageByID(id string) (cloudprovider.ICloudImage, error) {
	return self.GetImage(id)
}