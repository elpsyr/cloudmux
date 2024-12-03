package baiducfel

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/util/secrules"
)

type SecurityGroupRule struct {
	multicloud.SResourceBase
	BaiduTags
	region *SRegion

	Id            string `json:"id,omitempty"`
	Direction     string `json:"direction,omitempty"`
	PortRange     string `json:"portRange,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	Remark        string `json:"remark,omitempty"`
	SourceGroupID string `json:"sourceGroupId,omitempty"`
	SourceIP      string `json:"sourceIp,omitempty"`
}

var _ cloudprovider.ISecurityGroupRule = (*SecurityGroupRule)(nil)

// Delete implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) Delete() error {
	return s.region.doDelete(ServiceSecurityGroups, "/v2/securityGroup/rule/"+s.Id)
}

// GetAction implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetAction() secrules.TSecurityRuleAction {
	return "allow"
}

// GetCIDRs implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetCIDRs() []string {
	return []string{s.SourceIP}
}

// GetDescription implements cloudprovider.ISecurityGroupRule.
// Subtle: this method shadows the method (SResourceBase).GetDescription of SecurityGroupRule.SResourceBase.
func (s *SecurityGroupRule) GetDescription() string {
	return s.Remark
}

// GetDirection implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetDirection() secrules.TSecurityRuleDirection {
	return secrules.TSecurityRuleDirection(s.Direction)
}

// GetGlobalId implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetGlobalId() string {
	return s.Id
}

// GetPorts implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetPorts() string {
	return s.PortRange
}

// GetPriority implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetPriority() int {
	return 0
}

// GetProtocol implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetProtocol() string {
	return s.Protocol
}

// Update implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) Update(opts *cloudprovider.SecurityGroupRuleUpdateOptions) error {
	var params = map[string]interface{}{
		"securityGroupRuleId": s.Id,
		"protocol":            opts.Protocol,
		"portRange":           opts.Ports,
		"remark":              opts.Desc,
	}
	if s.Direction == "ingress" && len(opts.CIDR) > 0 {
		params["sourceIp"] = opts.CIDR
	}
	if s.Direction == "egress" && len(opts.CIDR) > 0 {
		params["destIp"] = opts.CIDR
	}

	_, err := s.region.doPut(ServiceSecurityGroups, "/v2/securityGroup/rule/update", params)
	return err
}
