package aliyun

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

func (self *SStoragecache) CfelGetICloudImages() ([]cloudprovider.ICloudImage, error) {
	images := make([]SImage, 0)
	for {
		parts, total, err := self.region.GetImages(ImageStatusType(""), ImageOwnerSystem, nil, "", len(images), 50)
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
		images[i].storageCache = self
		ret = append(ret, &images[i])
	}
	return ret, nil
}
