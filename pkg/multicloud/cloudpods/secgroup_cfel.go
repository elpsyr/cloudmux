package cloudpods

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	api "yunion.io/x/onecloud/pkg/apis/compute"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

func (self *SRegion) CfelGetSecurityGroups() ([]SSecurityGroup, error) {
	params := map[string]interface{}{
		"cloud_env": "",
	}
	ret := []SSecurityGroup{}
	return ret, self.cli.list(&modules.SecGroups, params, &ret)
}

func (self *SSecurityGroup) CreateRule(opts *cloudprovider.SecurityGroupRuleCreateOptions) (cloudprovider.ISecurityGroupRule, error) {
	return self.region.CfelCreateSecRule(self.Id, opts)
}

func (self *SRegion) CfelCreateSecRule(secId string, opts *cloudprovider.SecurityGroupRuleCreateOptions) (cloudprovider.ISecurityGroupRule, error) {
	input := api.SSecgroupRuleCreateInput{}
	input.SecgroupId = secId
	input.Priority = &opts.Priority
	input.Action = string(opts.Action)
	input.Protocol = string(opts.Protocol)
	input.Direction = string(opts.Direction)
	input.Description = opts.Desc
	input.CIDR = opts.CIDR

	input.Ports = opts.Ports
	var ret = &SecurityGroupRule{region: self}
	return ret, self.create(&modules.SecGroupRules, input, &ret)
}
