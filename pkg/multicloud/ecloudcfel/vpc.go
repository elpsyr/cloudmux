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

package ecloudcfel

import (
	"context"
	"fmt"
	"slices"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SVpc struct {
	multicloud.SVpc
	EcloudTags

	region *SRegion
	iwires []cloudprovider.ICloudWire

	// wires
	// secgroups

	Id       string
	Name     string
	Region   string
	EcStatus string
	RouterId string
	Scale    string
	UserId   string
	UserName string
	Cidr     string

	FirstNetworkId string `json:"firstNetworkId"`
}

func (v *SVpc) GetId() string {
	return v.Id
}

func (v *SVpc) GetName() string {
	return v.Name
}

func (v *SVpc) GetGlobalId() string {
	return v.GetId()
}

func (v *SVpc) GetStatus() string {
	switch v.EcStatus {
	case "ready":
		return "ready"
	case "ACTIVE":
		return api.VPC_STATUS_AVAILABLE
	case "DOWN", "BUILD", "ERROR":
		return api.VPC_STATUS_UNAVAILABLE
	case "PENDING_DELETE":
		return api.VPC_STATUS_DELETING
	case "PENDING_CREATE", "PENDING_UPDATE":
		return api.VPC_STATUS_PENDING
	default:
		return api.VPC_STATUS_UNKNOWN
	}
}

func (v *SVpc) Refresh() error {
	n, err := v.region.getVpcById(v.Id)
	if err != nil {
		return err
	}
	return jsonutils.Update(v, n)
	// TODO? v.fetchWires()
}

func (v *SVpc) IsEmulated() bool {
	return false
}

func (v *SVpc) GetRegion() cloudprovider.ICloudRegion {
	return v.region
}

func (v *SVpc) GetIsDefault() bool {
	return false
}

func (v *SVpc) GetCidrBlock() string {
	return v.Cidr
}

func (v *SVpc) GetIWires() ([]cloudprovider.ICloudWire, error) {
	if v.iwires == nil {
		err := v.fetchWires()
		if err != nil {
			return nil, err
		}
	}
	return v.iwires, nil
}

func (v *SVpc) fetchWires() error {
	networks, err := v.region.GetNetworks(v.RouterId, "")
	if err != nil {
		return err
	}
	izones, err := v.region.GetIZones()
	if err != nil {
		return errors.Wrap(err, "unable to GetZones")
	}
	findZone := func(zoneRegion string) *SZone {
		for i := range izones {
			zone := izones[i].(*SZone)
			if zone.Region == zoneRegion {
				return zone
			}
		}
		return nil
	}
	zoneRegion2Wire := map[string]*SWire{}
	for i := range networks {
		zoneRegion := networks[i].Region
		zone := findZone(zoneRegion)
		var (
			wire *SWire
			ok   bool
		)
		if wire, ok = zoneRegion2Wire[zoneRegion]; !ok {
			wire = &SWire{
				vpc:  v,
				zone: zone,
			}
			zoneRegion2Wire[zoneRegion] = wire
		}
		wire.inetworks = append(wire.inetworks, &networks[i])
	}
	iwires := make([]cloudprovider.ICloudWire, 0, len(zoneRegion2Wire))
	for _, wire := range zoneRegion2Wire {
		iwires = append(iwires, wire)
	}
	v.iwires = iwires
	return nil
}

func (v *SVpc) GetISecurityGroups() ([]cloudprovider.ICloudSecurityGroup, error) {
	return nil, nil
}

func (v *SVpc) GetIRouteTables() ([]cloudprovider.ICloudRouteTable, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (v *SVpc) GetIRouteTableById(routeTableId string) (cloudprovider.ICloudRouteTable, error) {
	return nil, nil
}

func (v *SVpc) Delete() error {
	param := map[string]interface{}{
		"productType": "router",
		"resourceId":  v.RouterId,
	}
	req := NewConsoleRequest(v.region.ID, "/api/openapi-vpc/customer/v3/order/delete", nil, jsonutils.Marshal(param))
	_, err := v.region.client.doPost(req)
	return err
}

// 用来返回创建vpc时默认创建的网络id
func (v *SVpc) GetAuthorityOwnerId() string {
	return v.FirstNetworkId
}

func (v *SVpc) GetIWireById(wireId string) (cloudprovider.ICloudWire, error) {
	iwires, err := v.GetIWires()
	if err != nil {
		return nil, err
	}
	for i := range iwires {
		if iwires[i].GetGlobalId() == wireId {
			return iwires[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) CreateIVpc(opts *cloudprovider.VpcCreateOptions) (cloudprovider.ICloudVpc, error) {
	var zone *SZone
	for _, z := range self.izones {
		if opts.Zone == z.GetGlobalId() {
			zone = z.(*SZone)
		}
	}
	if zone == nil {
		return nil, fmt.Errorf("zone: %s not exist", opts.Zone)
	}
	params := map[string]interface{}{
		"name":            opts.NAME,
		"specs":           "high", //enum(normal,high,middle,mega)
		"networkName":     opts.NAME + "-" + "first-network",
		"cidr":            opts.CIDR,
		"region":          zone.Region,
		"networkTypeEnum": "VM",
	}
	req := NewConsoleRequest(self.ID, "/api/openapi-vpc/customer/v3/order/create/vpc", nil, jsonutils.Marshal(params))
	res, err := self.client.doPost(req)
	if err != nil {
		return nil, err
	}
	var ret createResp
	if err := res.Unmarshal(&ret); err != nil {
		return nil, err
	}

	var ids []string
	for i := 1; i < 4; i++ {
		time.Sleep(3 * time.Second)
		ids, err = self.getOrderInfo(ret.OrderId)
		if err != nil {
			fmt.Printf("getOrderInfo [%d] err: %v", i, err)
		}
		if len(ids) > 0 {
			break
		}
	}

	var vpc SVpc
	if len(ids) > 0 {
		req = NewConsoleRequest(self.ID, "/api/openapi-vpc/customer/v3/vpc/router/"+ids[0], nil, nil)
		err = self.client.doGet(context.Background(), req, &vpc)
		if err != nil {
			return nil, err
		}
		vpc.EcStatus = "ready"
		var isReady bool
		for i := 1; i < 4; i++ {
			net, _ := self.GetNetworkById(ids[0], zone.Region, vpc.FirstNetworkId)
			if net != nil && net.EcStatus == "ACTIVE" {
				isReady = true
				break
			}
			time.Sleep(3 * time.Second)
		}
		if !isReady {
			vpc.EcStatus = "creat_fail"
		}
	} else {
		vpc.Id = ret.OrderId
		vpc.EcStatus = "creat_fail"
	}

	vpc.region = self

	vpc.Cidr = opts.CIDR
	// vpc := &SVpc{region: self, Id: order[0].InstanceId, Name: opts.NAME, EcStatus: "ready", Cidr: opts.CIDR}
	return &vpc, nil
}

func (self *SVpc) CfelCreateSubnet(opts *cloudprovider.SNetworkCreateOptions) (cloudprovider.ICloudNetwork, error) {
	self.region.fetchZones()

	index := slices.IndexFunc(self.region.izones, func(iz cloudprovider.ICloudZone) bool {
		return iz.GetGlobalId() == opts.ZoneId
	})
	if index == -1 {
		return nil, fmt.Errorf("zone not found")
	}
	var zone = self.region.izones[index].(*SZone)
	params := map[string]interface{}{
		"availabilityZoneHints": zone.Region,
		"networkName":           opts.Name,
		"networkTypeEnum":       "VM",
		"routerId":              self.RouterId,
		"subnets": []map[string]string{
			{
				"cidr":      opts.Cidr,
				"ipVersion": "4",
			},
		},
	}
	req := NewConsoleRequest(self.region.ID, "/api/openapi-vpc/customer/v3/network", nil, jsonutils.Marshal(params))
	res, err := self.region.client.doPost(req)
	if err != nil {
		return nil, err
	}
	return &SNetwork{
		Id: res.Interface().(string),
		Subnets: []SSubnet{{
			Id:         "",
			Name:       "",
			NetworkId:  "",
			Region:     "",
			GatewayIp:  "",
			EnableDHCP: false,
			Cidr:       opts.Cidr,
			IpVersion:  "",
		}}}, nil
	return nil, nil
}
