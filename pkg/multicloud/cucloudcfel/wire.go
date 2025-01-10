package cucloudcfel

import (
	"fmt"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SWire struct {
	multicloud.SResourceBase
	CuCloudTags
	vpc      *SVpc
	zone     *SZone
	networks []cloudprovider.ICloudNetwork
}

// CreateINetwork implements cloudprovider.ICloudWire.
func (s *SWire) CreateINetwork(opts *cloudprovider.SNetworkCreateOptions) (cloudprovider.ICloudNetwork, error) {
	params := map[string]interface{}{
		"cloudRegionCode":     s.vpc.region.GetId(),
		"zoneCode":            s.zone.ZoneCode,
		"networkId":           s.vpc.GetGlobalId(),
		"subNetworkName":      strings.ReplaceAll(opts.Name, "-", "_"),
		"subNetworkIpVersion": "4",
		"subNetworkCidr":      opts.Cidr,
		"isEnableIpv6":        false,
	}
	res, err := s.vpc.region.client.post("/instance/v1/product/subnets", params)
	if err != nil {
		return nil, err
	}
	id, err := res.GetString("resourceId")
	if err != nil {
		return nil, err
	}
	return &SNetwork{
		wire:           s,
		SubNetworkCidr: opts.Cidr,
		// Description: opts.Desc,
		SubNetworkName: opts.Name,
		SubNetworkID:   id,
		VpcID:          s.vpc.GetGlobalId(),
		ZoneCode:       s.zone.ZoneCode,
	}, nil
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
func (s *SWire) GetINetworkById(netid string) (cloudprovider.ICloudNetwork, error) {
	params := map[string]interface{}{
		"cloudRegionCode": s.vpc.region.GetId(),
		"subNetworkId":    netid,
	}
	res, err := s.vpc.region.client.get("/instance/v1/product/subnets", params)
	if err != nil {
		return nil, err
	}
	var r []SNetwork
	err = res.Unmarshal(&r, "list")
	if err != nil {
		return nil, err
	}
	if len(r) == 0 || r[0].SubNetworkID != netid {
		return nil, fmt.Errorf("network %s not found", netid)
	}
	r[0].wire = s
	return &r[0], nil
}

// GetINetworks implements cloudprovider.ICloudWire.
func (s *SWire) GetINetworks() ([]cloudprovider.ICloudNetwork, error) {
	params := map[string]interface{}{
		"cloudRegionCode": s.vpc.region.GetId(),
		"pageNum":         "1",
		"pageSize":        "1000",
	}
	res, err := s.vpc.region.client.list("/instance/v1/product/subnets", params)
	if err != nil {
		return nil, err
	}
	var r []SNetwork
	err = res.Unmarshal(&r, "list")
	if err != nil {
		return nil, err
	}
	var result []cloudprovider.ICloudNetwork
	for i := range r {
		r[i].wire = s
		result = append(result, &r[i])
	}
	return result, nil
}

// GetIVpc implements cloudprovider.ICloudWire.
func (s *SWire) GetIVpc() cloudprovider.ICloudVpc {
	return s.vpc
}

// GetIZone implements cloudprovider.ICloudWire.
func (s *SWire) GetIZone() cloudprovider.ICloudZone {
	return s.zone
}

// GetId implements cloudprovider.ICloudWire.
func (s *SWire) GetId() string {
	return s.GetGlobalId()
}

// GetName implements cloudprovider.ICloudWire.
func (s *SWire) GetName() string {
	return ""
}

// GetStatus implements cloudprovider.ICloudWire.
func (s *SWire) GetStatus() string {
	return ""
}

var _ cloudprovider.ICloudWire = (*SWire)(nil)
