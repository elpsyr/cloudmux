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
	"fmt"
	"strconv"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/util/secrules"
)

type SecurityGroupRule struct {
	multicloud.SResourceBase
	EcloudTags
	region *SRegion

	Id             string `json:"id"`
	Protocol       string `json:"protocol"`
	Direction      string `json:"direction"`
	Description    string `json:"description"`
	CreatedTime    string `json:"createdTime"`
	SecgroupId     string `json:"secgroupId"`
	Status         string `json:"status"`
	AimSgid        string `json:"aimSgid"`
	MaxPortRange   int    `json:"maxPortRange"`
	MinPortRange   int    `json:"minPortRange"`
	EtherType      string `json:"etherType"`
	RemoteIpPrefix string `json:"remoteIpPrefix"`
	DefaultRule    bool   `json:"defaultRule"`
}

func (self *SecurityGroupRule) GetGlobalId() string {
	return self.Id
}

func (self *SecurityGroupRule) GetAction() secrules.TSecurityRuleAction {
	return secrules.TSecurityRuleAction("allow")
}

func (self *SecurityGroupRule) GetDescription() string {
	return self.Description
}

func (self *SecurityGroupRule) GetDirection() secrules.TSecurityRuleDirection {
	direction := "out"
	if self.Direction == "ingress" {
		direction = "in"
	}
	return secrules.TSecurityRuleDirection(direction)
}

func (self *SecurityGroupRule) GetCIDRs() []string {
	return []string{self.RemoteIpPrefix}
}

func (self *SecurityGroupRule) GetProtocol() string {
	return self.Protocol
}

func (self *SecurityGroupRule) GetPorts() string {
	if self.MaxPortRange == self.MinPortRange {
		return strconv.Itoa(self.MinPortRange)
	}
	return fmt.Sprintf("%d-%d", self.MinPortRange, self.MaxPortRange)
}

func (self *SecurityGroupRule) GetPriority() int {
	return 0
}

func (self *SecurityGroupRule) Delete() error {
	req := NewConsoleRequest(self.region.ID, "/api/openapi-vpc/customer/v3/SecurityGroupRule/"+self.Id, nil, nil)
	return self.region.client.doDelete(req)
}

func (self *SecurityGroupRule) Update(opts *cloudprovider.SecurityGroupRuleUpdateOptions) error {

	var minPort, maxPort string
	if strings.Contains(opts.Ports, "-") {
		arr := strings.Split(opts.Ports, "-")
		minPort, maxPort = arr[0], arr[1]
	} else {
		minPort, maxPort = opts.Ports, opts.Ports
	}
	params := map[string]interface{}{
		"direction":       self.Direction,
		"description":     opts.Desc,
		"etherType":       "IPv4",
		"maxPortRange":    maxPort,
		"minPortRange":    minPort,
		"protocol":        strings.ToUpper(opts.Protocol),
		"remoteType":      "cidr",
		"remoteIpPrefix":  opts.CIDR,
		"securityGroupId": self.SecgroupId,
	}
	req := NewConsoleRequest(self.region.ID, "/api/openapi-vpc/customer/v3/SecurityGroupRule/update/"+self.Id, nil, jsonutils.Marshal(params))
	_, err := self.region.client.doPut(req)
	return err
}
