package cucloudcfel

import (
	"fmt"
	"strings"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/log"
	"yunion.io/x/pkg/utils"
)

type SSecurityGroup struct {
	multicloud.SSecurityGroup
	CuCloudTags
	region *SRegion

	CloudRegionCode            string              `json:"cloudRegionCode"`
	Description                string              `json:"description"`
	RegionCode                 string              `json:"regionCode"`
	RegionVersion              string              `json:"regionVersion"`
	RemoteSecurityGroupRuleNum int64               `json:"remoteSecurityGroupRuleNum"`
	SecurityGroupID            string              `json:"securityGroupId"`
	SecurityGroupName          string              `json:"securityGroupName"`
	SecurityGroupRules         []SecurityGroupRule `json:"securityGroupRules"`
	SecurityGroupUUID          string              `json:"securityGroupUuid"`
}

// CreateRule implements cloudprovider.ICloudSecurityGroup.
// Subtle: this method shadows the method (SSecurityGroup).CreateRule of SSecurityGroup.SSecurityGroup.
func (s *SSecurityGroup) CreateRule(opts *cloudprovider.SecurityGroupRuleCreateOptions) (cloudprovider.ISecurityGroupRule, error) {
	direction := "egress"
	if opts.Direction == "in" {
		direction = "ingress"
	}
	var minPort, maxPort string
	if strings.Contains(opts.Ports, "-") {
		arr := strings.Split(opts.Ports, "-")
		minPort, maxPort = arr[0], arr[1]
	} else {
		minPort, maxPort = opts.Ports, opts.Ports
	}
	param := map[string]interface{}{
		"securityGroupRules": []map[string]interface{}{
			{
				"direction":       direction,
				"protocol":        opts.Protocol,
				"ethertype":       "ipv4",
				"remoteIpPrefix":  opts.CIDR,
				"securityGroupId": s.SecurityGroupID,
				"remoteType":      "0",
				"portRangeMax":    maxPort,
				"portRangeMin":    minPort,
			},
		},
		"cloudRegionCode": s.region.CloudRegionCode,
	}
	res, err := s.region.client.post("/instance/v1/product/securitygrouprule", param)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	return nil, nil
}

// Delete implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) Delete() error {
	params := map[string]interface{}{
		"cloudRegionCode": s.region.CloudRegionCode,
	}
	return s.region.client.delete("/instance/v1/product/securitygroup/"+s.SecurityGroupID, params)
}

// GetDescription implements cloudprovider.ICloudSecurityGroup.
// Subtle: this method shadows the method (SSecurityGroup).GetDescription of SSecurityGroup.SSecurityGroup.
func (s *SSecurityGroup) GetDescription() string {
	return s.Description
}

// GetGlobalId implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetGlobalId() string {
	return s.SecurityGroupID
}

// GetId implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetId() string {
	return s.SecurityGroupID
}

// GetName implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetName() string {
	return s.SecurityGroupName
}

// GetRules implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetRules() ([]cloudprovider.ISecurityGroupRule, error) {
	var rules []cloudprovider.ISecurityGroupRule
	for i := range s.SecurityGroupRules {
		rules = append(rules, &s.SecurityGroupRules[i])
	}
	return rules, nil
}

// GetStatus implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetStatus() string {
	return "ready"
}

// GetVpcId implements cloudprovider.ICloudSecurityGroup.
func (s *SSecurityGroup) GetVpcId() string {
	return ""
}

var _ cloudprovider.ICloudSecurityGroup = (*SSecurityGroup)(nil)

func (region *SRegion) CreateISecurityGroup(conf *cloudprovider.SecurityGroupCreateInput) (cloudprovider.ICloudSecurityGroup, error) {
	uuid := utils.GenRequestId(20)
	params := map[string]interface{}{
		"cloudRegionCode": region.CloudRegionCode,
		// "productNo": "",
		"securityGroupName":      strings.ReplaceAll(conf.Name, "-", "_"),
		"securityGroupDesc":      conf.Desc + ";" + uuid,
		"securityGroupRuleModel": "default",
	}
	_, err := region.client.post("/instance/v1/product/securitygroup", params)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 3; i++ {
		time.Sleep(3 * time.Second)
		sgs, err := region.GetSecurityGroups()
		if err != nil {
			log.Warningf("get ses err:%v", err)
			// time.Sleep(3 * time.Second)
			continue
		}
		for j := range sgs {
			if strings.HasSuffix(sgs[j].Description, uuid) {
				return &sgs[j], nil
			}
		}
	}

	return nil, cloudprovider.ErrNotSupported
}

func (region *SRegion) GetSecurityGroups() ([]SSecurityGroup, error) {
	params := map[string]interface{}{
		"cloudRegionCode": region.CloudRegionCode,
		"pageNum":         "1",
		"pageSize":        "1000",
	}
	res, err := region.client.list("/instance/v1/product/securitygroup", params)
	if err != nil {
		return nil, err
	}
	var r []SSecurityGroup
	err = res.Unmarshal(&r, "list")
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (region *SRegion) GetISecurityGroupById(secgroupId string) (cloudprovider.ICloudSecurityGroup, error) {
	params := map[string]interface{}{
		"cloudRegionCode": region.CloudRegionCode,
		"securityGroupId": secgroupId,
	}
	res, err := region.client.list("/instance/v1/product/securitygroup", params)
	if err != nil {
		return nil, err
	}
	var r []SSecurityGroup
	err = res.Unmarshal(&r, "list")
	if err != nil {
		return nil, err
	}
	if len(r) == 0 || r[0].SecurityGroupID != secgroupId {
		return nil, cloudprovider.ErrNotFound
	}
	r[0].region = region
	return &r[0], nil
}
