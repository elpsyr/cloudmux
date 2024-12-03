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
	// Subnets   []SSubNet
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

const ServiceVpcs = "vpcs"

func (s *SVpc) Keys() string {
	return "vpc"
}

func (s *SVpc) KeysPlural() string {
	return "vpcs"
}

// Delete implements cloudprovider.ICloudVpc.
func (s *SVpc) Delete() error {
	return s.region.doDelete(ServiceVpcs, "/v1/vpc/"+s.Id)
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

type secgResp struct {
	NextMarker  string           `json:"nextMarker"`
	Marker      string           `json:"marker"`
	MaxKeys     int              `json:"maxKeys"`
	IsTruncated bool             `json:"isTruncated"`
	Secgs       []SSecurityGroup `json:"securityGroups"`
}

// GetISecurityGroups implements cloudprovider.ICloudVpc.
func (s *SVpc) GetISecurityGroups() ([]cloudprovider.ICloudSecurityGroup, error) {
	var query = map[string]string{
		"vpcId": s.Id,
	}

	var secg []SSecurityGroup
	var marker string
	for {
		var r secgResp
		if len(marker) > 0 {
			query["marker"] = marker
		}
		res, err := s.region.doList(ServiceSecurityGroups, "v2/securityGroup", query)
		if err != nil {
			// if err.Error() == `{"statusCode":404}` { // 没有安全组
			// 	return nil, nil
			// }
			return nil, err
		}
		if err := res.Unmarshal(&r); err != nil {
			return nil, err
		}
		secg = append(secg, r.Secgs...)
		if !r.IsTruncated {
			break
		}
		marker = r.Marker
	}

	var ret []cloudprovider.ICloudSecurityGroup
	for i := range secg {
		ret = append(ret, &secg[i])
	}
	return ret, nil
}

// GetIWireById implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIWireById(wireId string) (cloudprovider.ICloudWire, error) {
	if s.iwires == nil {
		_, err := s.GetIWires()
		if err != nil {
			return nil, err
		}
	}
	for _, wire := range s.iwires {
		if wire.GetGlobalId() == wireId {
			return wire, nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

// GetIWires implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIWires() ([]cloudprovider.ICloudWire, error) {
	zones, err := s.region.GetICfelZones()
	if err != nil {
		return nil, err
	}
	ret := []cloudprovider.ICloudWire{}
	for i := range zones {
		wire := &SWire{zone: zones[i].(*SZone), vpc: s}
		ret = append(ret, wire)
	}
	s.iwires = ret
	return ret, nil
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
	}
	res, err := self.doPost(ServiceVpcs, "v1/vpc", params)
	if err != nil {
		return nil, err
	}
	id, _ := res.GetString("vpcId")
	var vpc = &SVpc{
		Id:   id,
		Name: opts.NAME,
		Desc: opts.Desc,
		Cidr: opts.CIDR,
		region: self,
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

	var vpcs []SVpc
	var query = map[string]string{
		"maxKeys": "100",
	}
	var marker string
	for {
		var r vpcResp
		if len(marker) > 0 {
			query["marker"] = marker
		}
		res, err := self.doList(ServiceVpcs, "v1/vpc", query)
		if err != nil {
			return nil,err
		}
		if err := res.Unmarshal(&r); err != nil {
			return nil,err
		}
		vpcs = append(vpcs, r.Vpcs...)
		if !r.IsTruncated {
			break
		}
		marker = r.NextMarker
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
	err := self.doGet(ServiceVpcs, "/v1/vpc/"+id, nil, &vpc)
	if err != nil {
		return nil, err
	}
	vpc.region = self
	return &vpc, nil
}
