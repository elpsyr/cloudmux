package aliyun

import "yunion.io/x/cloudmux/pkg/cloudprovider"

func (self *SVpc) Update(opts *cloudprovider.VpcUpdateOptions) error {
	return self.region.UpdateVpc(opts)
}
