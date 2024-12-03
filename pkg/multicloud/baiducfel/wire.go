package baiducfel

import (
	"fmt"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SWire struct {
	multicloud.SResourceBase
	BaiduTags
	vpc      *SVpc
	zone     *SZone
	networks []cloudprovider.ICloudNetwork
}

var _ cloudprovider.ICloudWire = (*SWire)(nil)

// CreateINetwork implements cloudprovider.ICloudWire.
func (s *SWire) CreateINetwork(opts *cloudprovider.SNetworkCreateOptions) (cloudprovider.ICloudNetwork, error) {
	params := map[string]interface{}{
		"zoneName":    opts.ZoneId,
		"subnetType":  "BCC",
		"name":        opts.Name,
		"description": opts.Desc,
		// "enableIpv6": true,
		"cidr":  opts.Cidr,
		"vpcId": s.vpc.Id,
		// "vpcSecondaryCidr":"172.17.0.0/16"
	}
	res, err := s.vpc.region.doPost(ServiceSubnet, "v1/subnet", params)
	if err != nil {
		return nil, err
	}
	id, _ := res.GetString("subnetId")
	var net = &SNetwork{
		Id:          id,
		Name:        opts.Name,
		Cidr:        opts.Cidr,
		Description: opts.Desc,
	}
	return net, nil
}

// GetBandwidth implements cloudprovider.ICloudWire.
func (s *SWire) GetBandwidth() int {
	return 1000000
}

// GetGlobalId implements cloudprovider.ICloudWire.
func (s *SWire) GetGlobalId() string {
	return fmt.Sprintf("%s-%s", s.vpc.GetGlobalId(), s.zone.GetGlobalId())
}

// GetINetworkById implements cloudprovider.ICloudWire.
func (s *SWire) GetINetworkById(id string) (cloudprovider.ICloudNetwork, error) {
	var net SNetwork
	err := s.vpc.region.doGet(ServiceSubnet, "/v1/subnet/"+id, nil, &net)
	if err != nil {
		return nil, err
	}
	net.wire = s
	return &net, nil
}

type subnetResp struct {
	NextMarker  string     `json:"nextMarker"`
	Marker      string     `json:"marker"`
	MaxKeys     int        `json:"maxKeys"`
	IsTruncated bool       `json:"isTruncated"`
	Networks    []SNetwork `json:"subnets"`
}



// GetIVpc implements cloudprovider.ICloudWire.
func (s *SWire) GetIVpc() cloudprovider.ICloudVpc {
	panic("unimplemented")
}

// GetIZone implements cloudprovider.ICloudWire.
func (s *SWire) GetIZone() cloudprovider.ICloudZone {
	return s.zone
}

// GetId implements cloudprovider.ICloudWire.
func (s *SWire) GetId() string {
	return fmt.Sprintf("%s-%s", s.vpc.GetGlobalId(), s.vpc.region.GetGlobalId())
}

// GetName implements cloudprovider.ICloudWire.
func (s *SWire) GetName() string {
	return s.GetId()
}

// GetStatus implements cloudprovider.ICloudWire.
func (s *SWire) GetStatus() string {
	return api.WIRE_STATUS_AVAILABLE
}
