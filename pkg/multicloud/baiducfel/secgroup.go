package baiducfel

import (
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SSecurityGroup struct {
	multicloud.SSecurityGroup
	BaiduTags
	region *SRegion

	BindInstanceNum int64               `json:"bindInstanceNum"`
	CreatedTime     string              `json:"createdTime"`
	Desc            string              `json:"desc"`
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Rules           []SecurityGroupRule `json:"rules"`
	SgVersion       int64               `json:"sgVersion"`
	// Tags            []interface{} `json:"tags"`
	VpcID string `json:"vpcId"`
}

var _ cloudprovider.ICloudSecurityGroup = (*SSecurityGroup)(nil)

const ServiceSecurityGroups = "securityGroups"

// CreateRule implements cloudprovider.ICloudSecurityGroup.
// Subtle: this method shadows the method (SSecurityGroup).CreateRule of SSecurityGroup.SSecurityGroup.
func (s *SSecurityGroup) CreateRule(opts *cloudprovider.SecurityGroupRuleCreateOptions) (cloudprovider.ISecurityGroupRule, error) {
	var params = map[string]interface{}{
		"rule": map[string]string{
			"remark":        opts.Desc,
			"protocol":      strings.ToLower(opts.Protocol),
			"portRange":     opts.Ports,
			"direction":     string(opts.Direction),
			"sourceIp":      opts.CIDR,
			// "sourceGroupId": s.ID,
		},
		"securityGroupId":s.ID,
	}
	_, err := s.region.doPut(ServiceSecurityGroups, "/v2/securityGroup/"+s.ID+"?authorizeRule", params)
	if err != nil {
		return nil, err
	}
	return &SecurityGroupRule{}, nil
}

// Delete implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) Delete() error {
	return s.region.doDelete(ServiceSecurityGroups,"/v2/securityGroup/" + s.ID)
}

// GetGlobalId implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetGlobalId() string {
	return s.ID
}

// GetId implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetId() string {
	return s.ID
}

// GetName implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetName() string {
	return s.Name
}

// GetRules implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetRules() ([]cloudprovider.ISecurityGroupRule, error) {
	var ret []cloudprovider.ISecurityGroupRule
	for i := range s.Rules {
		ret = append(ret, &s.Rules[i])
	}
	return ret, nil
}

// GetStatus implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetStatus() string {
	return "ready"
}

// GetVpcId implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetVpcId() string {
	return s.VpcID
}

func (region *SRegion) CreateISecurityGroup(conf *cloudprovider.SecurityGroupCreateInput) (cloudprovider.ICloudSecurityGroup, error) {
	var param = map[string]interface{}{
		"name":  conf.Name,
		"desc":  conf.Name,
		"vpcId": conf.VpcId,
	}
	res, err := region.doPost(ServiceSecurityGroups, "/v2/securityGroup", param)
	if err != nil {
		return nil, err
	}
	id, _ := res.GetString("securityGroupId")
	return &SSecurityGroup{ID: id, Name: conf.Name}, nil
}

func (region *SRegion) GetISecurityGroupById(secgroupId string) (cloudprovider.ICloudSecurityGroup, error) {
	var segc SSecurityGroup
	res, err := region.doGetWithoutVal(ServiceSecurityGroups, "/v2/securityGroup/"+secgroupId, nil)
	if err != nil {
		return nil, err
	}
	if err = res.Unmarshal(&segc); err != nil {
		return nil, err
	}
	segc.region = region
	return &segc, nil
}
