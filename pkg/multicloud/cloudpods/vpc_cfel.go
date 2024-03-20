package cloudpods

import "yunion.io/x/cloudmux/pkg/cloudprovider"

// CfelGetISecurityGroups 获取 vpc 下的安全组
func (self *SVpc) CfelGetISecurityGroups() ([]cloudprovider.ICloudSecurityGroup, error) {
	groups, err := self.region.CfelGetSecurityGroupsByVpcId(self.GetGlobalId())
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICloudSecurityGroup
	for i := range groups {
		groups[i].region = self.region
		ret = append(ret, &groups[i])
	}
	return ret, nil
}
