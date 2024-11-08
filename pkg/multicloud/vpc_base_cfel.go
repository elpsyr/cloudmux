package multicloud

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

func (self *SVpc) CfelCreateSubnet(opts *cloudprovider.SNetworkCreateOptions) (cloudprovider.ICloudNetwork, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelCreateSubnet")
}

func (self *SVpc) Update(opts *cloudprovider.VpcUpdateOptions) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "Update")
}
