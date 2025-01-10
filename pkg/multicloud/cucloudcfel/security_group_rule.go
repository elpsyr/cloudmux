package cucloudcfel

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/util/secrules"
)

type SecurityGroupRule struct {
	multicloud.SResourceBase
	CuCloudTags
	region *SRegion

	Action                string `json:"action"`
	CreateTime            string `json:"createTime"`
	Description           string `json:"description"`
	Direction             string `json:"direction"`
	Ethertype             string `json:"ethertype"`
	PortRangeMax          string `json:"portRangeMax"`
	PortRangeMin          string `json:"portRangeMin"`
	Priority              string `json:"priority"`
	Protocol              string `json:"protocol"`
	RemoteIPPrefix        string `json:"remoteIpPrefix"`
	RuleName              string `json:"ruleName"`
	SecurityGroupID       string `json:"securityGroupId"`
	SecurityGroupRuleID   string `json:"securityGroupRuleId"`
	SecurityGroupRuleUUID string `json:"securityGroupRuleUuid"`
	UpdateTime            string `json:"updateTime"`
}

// Delete implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) Delete() error {
	panic("unimplemented")
}

// GetAction implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetAction() secrules.TSecurityRuleAction {
	return secrules.TSecurityRuleAction(s.Action)
}

// GetCIDRs implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetCIDRs() []string {
	return []string{s.RemoteIPPrefix}
}

// GetDescription implements cloudprovider.ISecurityGroupRule.
// Subtle: this method shadows the method (SResourceBase).GetDescription of SecurityGroupRule.SResourceBase.
func (s *SecurityGroupRule) GetDescription() string {
	return s.Description
}

// GetDirection implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetDirection() secrules.TSecurityRuleDirection {
	return secrules.TSecurityRuleDirection(s.Direction)
}

// GetGlobalId implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetGlobalId() string {
	return s.SecurityGroupRuleID
}

// GetPorts implements cloudprovider.ISecurityGroupRule.
func (s *SecurityGroupRule) GetPorts() string {
	if s.PortRangeMin == s.PortRangeMax {
		return s.PortRangeMin
	}
	return s.PortRangeMin + "-" + s.PortRangeMax
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
	panic("unimplemented")
}

var _ cloudprovider.ISecurityGroupRule = (*SecurityGroupRule)(nil)
