package cucloudcfel

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

type CuCloudTags struct {
}

func (self *CuCloudTags) GetTags() (map[string]string, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetTags")
}

func (self *CuCloudTags) GetSysTags() map[string]string {
	return nil
}

func (self *CuCloudTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}
