package baiducfel

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SVpc struct {
	multicloud.SVpc
	BaiduTags

	region *SRegion
	iwires []cloudprovider.ICloudWire

	// wires
	// secgroups

	Id        string `json:"vpcId"`
	Name      string
	Desc      string
	Cidr      string
	Ipv6Cidr  string
	IsDefault bool
	Tags      []BaiduTags
	Subnets   []SSubNet
}

var _ cloudprovider.ICloudVpc = (*SVpc)(nil)

type SSubNet struct {
	Name        string `json:"name"`
	SubnetId    string `json:"subnetId"`
	ZoneName    string `json:"zoneName"`
	Cidr        string `json:"cidr"`
	Ipv6Cidr    string `json:"ipv6Cidr"`
	VpcId       string `json:"vpcId"`
	SubnetType  string `json:"subnetType"`
	Description string `json:"description"`
	CreatedTime string `json:"createdTime"`
}

func (s *SVpc) Keys() string {
	return "vpc"
}

func (s *SVpc) KeysPlural() string {
	return "vpcs"
}

// Delete implements cloudprovider.ICloudVpc.
func (s *SVpc) Delete() error {
	return s.region.doDelete("vpc", "/v1/vpc/"+s.Id)
}

// GetCidrBlock implements cloudprovider.ICloudVpc.
func (s *SVpc) GetCidrBlock() string {
	return s.Cidr
}

// GetDescription implements cloudprovider.ICloudVpc.
// Subtle: this method shadows the method (SVpc).GetDescription of SVpc.SVpc.
func (s *SVpc) GetDescription() string {
	return s.Desc
}

// GetGlobalId implements cloudprovider.ICloudVpc.
func (s *SVpc) GetGlobalId() string {
	return s.Id
}

// GetIRouteTableById implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIRouteTableById(routeTableId string) (cloudprovider.ICloudRouteTable, error) {
	panic("unimplemented")
}

// GetIRouteTables implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIRouteTables() ([]cloudprovider.ICloudRouteTable, error) {
	panic("unimplemented")
}

// GetISecurityGroups implements cloudprovider.ICloudVpc.
func (s *SVpc) GetISecurityGroups() ([]cloudprovider.ICloudSecurityGroup, error) {
	panic("unimplemented")
}

// GetIWireById implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIWireById(wireId string) (cloudprovider.ICloudWire, error) {
	panic("unimplemented")
}

// GetIWires implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIWires() ([]cloudprovider.ICloudWire, error) {
	panic("unimplemented")
}

// GetId implements cloudprovider.ICloudVpc.
func (s *SVpc) GetId() string {
	return s.Id
}

// GetIsDefault implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIsDefault() bool {
	return s.IsDefault
}

// GetName implements cloudprovider.ICloudVpc.
func (s *SVpc) GetName() string {
	return s.Name
}

// GetRegion implements cloudprovider.ICloudVpc.
func (s *SVpc) GetRegion() cloudprovider.ICloudRegion {
	return s.region
}

// GetStatus implements cloudprovider.ICloudVpc.
func (s *SVpc) GetStatus() string {
	return "ready"
}

func (self *SRegion) CreateIVpc(opts *cloudprovider.VpcCreateOptions) (cloudprovider.ICloudVpc, error) {
	params := map[string]interface{}{
		"name":        opts.NAME,
		"description": opts.Desc,
		"cidr":        opts.CIDR,
		"enableIpv6":  false,
		//  "tags":[
		// 	{
		// 	 "tagKey":"tagKey",
		// 	  "tagValue":"tagValue"
		// 	}
		//  ]
	}
	res, err := self.doPost("vpc", "v1/vpc", params)
	if err != nil {
		return nil, err
	}
	var vpc = &SVpc{
		Id:   res.Interface().(string),
		Name: opts.NAME,
		Desc: opts.Desc,
		Cidr: opts.CIDR,
	}
	return vpc, nil
}

type vpcResp struct {
	NextMarker  string `json:"nextMarker"`
	Marker      string `json:"marker"`
	MaxKeys     int    `json:"maxKeys"`
	IsTruncated bool   `json:"isTruncated"`
	Vpcs        []SVpc `json:"vpcs"`
}

func (self *SRegion) GetIVpcs() ([]cloudprovider.ICloudVpc, error) {
	res, err := self.doList("vpc", "v1/vpc", nil)
	if err != nil {
		return nil, err
	}

	var vpcs []SVpc

	for {
		var r vpcResp
		if err := res.Unmarshal(&r); err != nil {
			continue
		}
		if r.IsTruncated {
			break
		}
		vpcs = append(vpcs, r.Vpcs...)
	}

	var ret []cloudprovider.ICloudVpc
	for i := range vpcs {
		vpcs[i].region = self
		ret = append(ret, &vpcs[i])
	}
	return ret, nil
}

func (self *SRegion) GetIVpcById(id string) (cloudprovider.ICloudVpc, error) {
	var vpc SVpc
	return &vpc, self.doGet("vpc", "/v1/vpc/"+id, nil, &vpc)
}
