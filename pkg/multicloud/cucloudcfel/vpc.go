package cucloudcfel

import (
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SVpc struct {
	multicloud.SVpc
	CuCloudTags

	region *SRegion
	iwires []cloudprovider.ICloudWire

	BandWidth        string `json:"bandWidth"`
	Cidr             string `json:"cidr"`
	CloudRegionID    string `json:"cloudRegionId"`
	CloudRegionName  string `json:"cloudRegionName"`
	CreateTime       string `json:"createTime"`
	Description      string `json:"description"`
	DnatID           string `json:"dnatId"`
	Flag             string `json:"flag"`
	InstanceID       string `json:"instanceId"`
	IsDefaultNetwork string `json:"isDefaultNetwork"`
	ResourceGroupID  string `json:"resourceGroupId"`
	RouterID         string `json:"routerId"`
	RouterName       string `json:"routerName"`
	RouterNum        int64  `json:"routerNum"`
	RouterTableNum   int64  `json:"routerTableNum"`
	Status           string `json:"status"`
	SubNetNum        int64  `json:"subNetNum"`
	SubNetworkNum    int64  `json:"subNetworkNum"`
	VpcID            string `json:"vpcId"`
	VpcName          string `json:"vpcName"`
	VpcType          string `json:"vpcType"`
}

// Delete implements cloudprovider.ICloudVpc.
func (s *SVpc) Delete() error {
	params := map[string]interface{}{
		"cloudRegionCode": s.region.CloudRegionCode,
	}
	return s.region.client.delete(fmt.Sprintf("/instance/v1/product/vpcs/%s", s.VpcID), params)
}

// GetCidrBlock implements cloudprovider.ICloudVpc.
func (s *SVpc) GetCidrBlock() string {
	return s.Cidr
}

// GetDescription implements cloudprovider.ICloudVpc.
// Subtle: this method shadows the method (SVpc).GetDescription of SVpc.SVpc.
func (s *SVpc) GetDescription() string {
	return s.Description
}

// GetGlobalId implements cloudprovider.ICloudVpc.
func (s *SVpc) GetGlobalId() string {
	return s.VpcID
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
	sgs,err :=  s.region.GetSecurityGroups()
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICloudSecurityGroup
	for i := range sgs {
		sgs[i].region = s.region
		ret = append(ret, &sgs[i])
	}
	return ret, nil
}

// GetIWireById implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIWireById(wireId string) (cloudprovider.ICloudWire, error) {
	wires, err := s.GetIWires()
	if err != nil {
		return nil, err
	}
	for i := range wires {
		if wires[i].GetGlobalId() == wireId {
			return wires[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

// GetIWires implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIWires() ([]cloudprovider.ICloudWire, error) {
	if s.iwires != nil {
		return s.iwires, nil
	}
	zones, err := s.region.GetIZones()
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICloudWire
	for i := range zones {
		wire := &SWire{zone: zones[i].(*SZone), vpc: s}
		ret = append(ret, wire)
	}
	s.iwires = ret
	return ret, nil
}

// GetId implements cloudprovider.ICloudVpc.
func (s *SVpc) GetId() string {
	return s.VpcID
}

// GetIsDefault implements cloudprovider.ICloudVpc.
func (s *SVpc) GetIsDefault() bool {
	return s.IsDefaultNetwork == "true"
}

// GetName implements cloudprovider.ICloudVpc.
func (s *SVpc) GetName() string {
	return s.VpcName
}

// GetRegion implements cloudprovider.ICloudVpc.
func (s *SVpc) GetRegion() cloudprovider.ICloudRegion {
	return s.region
}

// GetStatus implements cloudprovider.ICloudVpc.
func (s *SVpc) GetStatus() string {
	return s.Status
}

var _ cloudprovider.ICloudVpc = (*SVpc)(nil)

func (self *SRegion) CreateIVpc(opts *cloudprovider.VpcCreateOptions) (cloudprovider.ICloudVpc, error) {
	params := map[string]interface{}{
		"cloudRegionCode": self.CloudRegionCode,
		"networkName":     opts.NAME,
		"description":     opts.Desc,
		"networkCidr":     opts.CIDR,
		// "networkType":     "vlan",
	}
	res, err := self.client.post("/instance/v1/product/vpcs", params)
	if err != nil {
		return nil, err
	}
	id, err := res.GetString("resourceId")
	if err != nil {
		return nil, err
	}
	var vpc = &SVpc{
		VpcID: id,
		VpcName: opts.NAME,
		Cidr:   opts.CIDR,
		region: self,
	}
	return vpc, nil
}

func (self *SRegion) GetIVpcById(id string) (cloudprovider.ICloudVpc, error) {
	params := map[string]interface{}{
		"cloudRegionCode": self.CloudRegionCode,
		"networkId":       id,
	}
	res, err := self.client.get("/instance/v1/product/vpcs", params)
	if err != nil {
		return nil, err
	}
	var r []SVpc
	err = res.Unmarshal(&r, "list")
	if err != nil {
		return nil, err
	}
	if len(r) == 0 || r[0].VpcID != id {
		return nil, fmt.Errorf("vpc %s not found", id)
	}
	r[0].region = self
	return &r[0], nil
}

func (self *SRegion) GetIVpcs() ([]cloudprovider.ICloudVpc, error) {
	params := map[string]interface{}{
		"cloudRegionCode": self.CloudRegionCode,
		"pageNum":         "1",
		"pageSize":        "1000",
	}
	res, err := self.client.get("/instance/v1/product/vpcs", params)
	if err != nil {
		return nil, err
	}
	var r []SVpc
	err = res.Unmarshal(&r, "list")
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICloudVpc
	for i := range r {
		r[i].region = self
		ret = append(ret, &r[i])
	}
	return ret, nil
}
