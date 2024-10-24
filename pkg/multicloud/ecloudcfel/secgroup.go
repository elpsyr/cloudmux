package ecloudcfel

// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"

	api "yunion.io/x/onecloud/pkg/apis/compute"
)

type SSecurityGroup struct {
	multicloud.SSecurityGroup
	EcloudTags
	region *SRegion

	Id          string `json:"id"`
	Name        string `json:"name"`
	CreatedTime string `json:"createdTime"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Defaulted   bool   `json:"defaulted"`
	Stateful    bool   `json:"stateful"`
	Scale       string `json:"scale"`  // 安全组规格：high为高性能安全组
	PortId      string `json:"portId"` //网卡id
	Region      string `json:"region"`
	VpoolId     string `json:"vpoolId"`
	Vaz         string `json:"vaz"`
	CountEcs    int    `json:"countEcs"`
	Status      string
}

func (self *SSecurityGroup) GetName() string {
	return self.Name
}

func (self *SSecurityGroup) GetDescription() string {
	return self.Description
}

func (self *SSecurityGroup) GetId() string {
	return self.Id
}

func (self *SSecurityGroup) GetGlobalId() string {
	return self.Id
}

func (self *SSecurityGroup) GetStatus() string {
	return "ready"
}

func (self *SSecurityGroup) GetProjectId() string {
	return ""
}

func (self *SSecurityGroup) GetRules() ([]cloudprovider.ISecurityGroupRule, error) {
	// /api/openapi-vpc/customer/v3/SecurityGroupRule
	query := map[string]string{
		"securityGroupId": self.Id,
		"page":            "1",
		"pageSize":        "100",
	}
	request := NewConsoleRequest(self.region.ID, "/api/openapi-vpc/customer/v3/SecurityGroupRule", query, nil)
	ret := []cloudprovider.ISecurityGroupRule{}
	rules := []SecurityGroupRule{}
	err := self.region.client.doList(context.Background(), request, &rules)
	if err != nil {
		return nil, err
	}
	for i := range rules {
		rules[i].region = self.region
		ret = append(ret, &rules[i])
	}
	return ret, nil
}

func (self *SSecurityGroup) GetVpcId() string {
	return api.NORMAL_VPC_ID
}

func (self *SSecurityGroup) GetReferences() ([]cloudprovider.SecurityGroupReference, error) {

	servers := []SInstance{}

	ret := []cloudprovider.SecurityGroupReference{}
	for i := range servers {
		ret = append(ret, cloudprovider.SecurityGroupReference{
			Id:   servers[i].Id,
			Name: servers[i].Name,
		})
	}
	return ret, nil
}

func (self *SRegion) DeleteSecRule(id string) error {
	req := NewConsoleRequest(self.ID, "/api/openapi-vpc/customer/v3/SecurityGroupRule/"+id, nil, nil)
	return self.client.doDelete(req)
}

func (self *SRegion) CreateSecRule(secId string, opts *cloudprovider.SecurityGroupRuleCreateOptions) error {
	directon := "egress"
	if opts.Direction == "in" {
		directon = "ingress"
	}
	var minPort, maxPort string
	if strings.Contains(opts.Ports, "-") {
		arr := strings.Split(opts.Ports, "-")
		minPort, maxPort = arr[0], arr[1]
	} else {
		minPort, maxPort = opts.Ports, opts.Ports
	}
	// https://ecloud.10086.cn/op-help-center/doc/article/73827
	params := map[string]interface{}{
		"description":    opts.Desc,
		"direction":      directon,
		"etherType":      "IPv4",
		"maxPortRange":   maxPort,
		"minPortRange":   minPort,
		"protocol":       strings.ToUpper(opts.Protocol),
		"remoteIpPrefix": opts.CIDR,
		// "remoteSecurityGroupId":"83da3346-8bb8-4e69-b614-3cd9bcb762ed",
		"remoteType":      "cidr",
		"securityGroupId": secId,
	}
	request := NewConsoleRequest(self.ID, "/api/openapi-vpc/customer/v3/SecurityGroupRule", nil, jsonutils.Marshal(params))
	_, err := self.client.doPost(request)
	return err
}

func (self *SSecurityGroup) CreateRule(opts *cloudprovider.SecurityGroupRuleCreateOptions) (cloudprovider.ISecurityGroupRule, error) {
	directon := "egress"
	if opts.Direction == "in" {
		directon = "ingress"
	}
	var minPort, maxPort string
	if strings.Contains(opts.Ports, "-") {
		arr := strings.Split(opts.Ports, "-")
		minPort, maxPort = arr[0], arr[1]
	} else {
		minPort, maxPort = opts.Ports, opts.Ports
	}
	// https://ecloud.10086.cn/op-help-center/doc/article/73827
	params := map[string]interface{}{
		"description":  opts.Desc,
		"direction":    directon,
		"etherType":    "IPv4",
		"maxPortRange": maxPort,
		"minPortRange": minPort,
		"protocol":     strings.ToUpper(opts.Protocol),
		// "remoteSecurityGroupId":"83da3346-8bb8-4e69-b614-3cd9bcb762ed",
		"remoteType":      "cidr",
		"remoteIpPrefix":  opts.CIDR,
		"securityGroupId": self.Id,
	}
	request := NewConsoleRequest(self.region.ID, "/api/openapi-vpc/customer/v3/SecurityGroupRule", nil, jsonutils.Marshal(params))
	ret, err := self.region.client.doPost(request)
	if err != nil {
		return nil, err
	}
	// ret.String()
	// id, _ := ret.GetString("body")
	rule := SecurityGroupRule{
		region:         self.region,
		Id:             ret.Interface().(string),
		Direction:      string(opts.Direction),
		RemoteIpPrefix: opts.CIDR,
	}
	return &rule, err
}

func (self *SSecurityGroup) Delete() error {
	req := NewConsoleRequest(self.region.ID, "/api/openapi-vpc/customer/v3/SecurityGroup/"+self.Id, nil, nil)
	return self.region.client.doDelete(req)
}

func (self *SRegion) GetSecurityGroups() ([]SSecurityGroup, error) {
	// /api/openapi-vpc/customer/v3/SecurityGroup
	request := NewConsoleRequest(self.ID, "/api/openapi-vpc/customer/v3/SecurityGroup", nil, nil)
	sec := []SSecurityGroup{}
	self.client.doList(context.Background(), request, &sec)
	ret := []SSecurityGroup{}
	return ret, nil
}

func (self *SRegion) GetSecurityGroup(id string) (*SSecurityGroup, error) {
	secgroup := &SSecurityGroup{region: self}
	req := NewConsoleRequest(self.ID, "/api/openapi-vpc/customer/v3/SecurityGroup/"+id, nil, nil)
	return secgroup, self.client.doGet(context.Background(), req, &secgroup)
}

func (self *SRegion) CreateISecurityGroup(opts *cloudprovider.SecurityGroupCreateInput) (cloudprovider.ICloudSecurityGroup, error) {

	// https://ecloud.10086.cn/op-help-center/doc/article/73819
	params := map[string]interface{}{
		"name":        opts.Name,
		"description": opts.Desc,
		"type":        "VM", //enum(VM,IRONIC,NAS,EW)
	}
	request := NewConsoleRequest(self.ID, "/api/openapi-vpc/customer/v3/SecurityGroup", nil, jsonutils.Marshal(params))

	ret, err := self.client.doPost(request)
	if err != nil {
		return nil, err
	}
	secgroup := &SSecurityGroup{region: self, Id: ret.Interface().(string), Status: "ready"}
	return secgroup, nil
}

func (self *SRegion) GetISecurityGroups() ([]cloudprovider.ICloudSecurityGroup, error) {
	// /api/openapi-vpc/customer/v3/SecurityGroup
	query := map[string]string{
		"types": "VM",
		"size":  "100",
		"page":  "1",
	}
	request := NewConsoleRequest(self.ID, "/api/openapi-vpc/customer/v3/SecurityGroup", query, nil)
	sec := []SSecurityGroup{}
	err := self.client.doList(context.Background(), request, &sec)
	if err != nil {
		return nil, err
	}
	ret := []cloudprovider.ICloudSecurityGroup{}
	for i := range sec {
		sec[i].region = self
		ret = append(ret, &sec[i])
	}
	return ret, nil
}

func (self *SRegion) GetISecurityGroupById(secgroupId string) (cloudprovider.ICloudSecurityGroup, error) {
	secgroup, err := self.GetSecurityGroup(secgroupId)
	if err != nil {
		return nil, err
	}
	return secgroup, nil
}
