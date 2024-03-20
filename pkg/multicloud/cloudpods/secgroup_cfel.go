package cloudpods

import modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"

func (self *SRegion) CfelGetSecurityGroups() ([]SSecurityGroup, error) {
	params := map[string]interface{}{
		"cloud_env": "",
	}
	ret := []SSecurityGroup{}
	return ret, self.cli.list(&modules.SecGroups, params, &ret)
}
