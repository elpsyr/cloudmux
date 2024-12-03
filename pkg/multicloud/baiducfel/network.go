package baiducfel

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/util/netutils"
	"yunion.io/x/pkg/util/rbacscope"
)

type SNetwork struct {
	multicloud.SNetworkBase
	BaiduTags

	wire *SWire

	AvailableIP           int64  `json:"availableIp"`
	AvailableUnreservedIP int64  `json:"availableUnreservedIp"`
	Cidr                  string `json:"cidr"`
	CreatedTime           string `json:"createdTime"`
	Description           string `json:"description"`
	Ipv6Cidr              string `json:"ipv6Cidr"`
	Name                  string `json:"name"`
	Id                    string `json:"subnetId"`
	SubnetType            string `json:"subnetType"`
	VpcID                 string `json:"vpcId"`
	ZoneName              string `json:"zoneName"`
}

var _ cloudprovider.ICloudNetwork = (*SNetwork)(nil)

const ServiceSubnet = "subnets"

// Delete implements cloudprovider.ICloudNetwork.
func (s *SNetwork) Delete() error {
	return s.wire.vpc.region.doDelete(ServiceSubnet, "/v1/subnet/"+s.Id)
}

// GetAllocTimeoutSeconds implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetAllocTimeoutSeconds() int {
	panic("unimplemented")
}

// GetDescription implements cloudprovider.ICloudNetwork.
// Subtle: this method shadows the method (SNetworkBase).GetDescription of SNetwork.SNetworkBase.
func (s *SNetwork) GetDescription() string {
	return s.Description
}

// GetGateway implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetGateway() string {
	pref, _ := netutils.NewIPV4Prefix(s.Cidr)
	startIp := pref.Address.NetAddr(pref.MaskLen) // 0
	startIp = startIp.StepUp()
	return startIp.String()
}

// GetGlobalId implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetGlobalId() string {
	return s.Id
}

// GetIWire implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetIWire() cloudprovider.ICloudWire {
	return s.wire
}

// GetId implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetId() string {
	return s.Id
}

// GetIpEnd implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetIpEnd() string {
	pref, _ := netutils.NewIPV4Prefix(s.Cidr)
	endIp := pref.Address.BroadcastAddr(pref.MaskLen) // 255
	endIp = endIp.StepDown()                          // 254
	return endIp.String()
}

// GetIpMask implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetIpMask() int8 {
	pref, _ := netutils.NewIPV4Prefix(s.Cidr)
	return pref.MaskLen
}

// GetIpStart implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetIpStart() string {
	pref, _ := netutils.NewIPV4Prefix(s.Cidr)
	startIp := pref.Address.NetAddr(pref.MaskLen) // 0
	startIp = startIp.StepUp()                    // 1
	startIp = startIp.StepUp()                    // 2
	return startIp.String()
}

// GetName implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetName() string {
	return s.Name
}

// GetProjectId implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetProjectId() string {
	return ""
}

// GetPublicScope implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetPublicScope() rbacscope.TRbacScope {
	panic("unimplemented")
}

// GetServerType implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetServerType() string {
	return "server"
}

// GetStatus implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetStatus() string {
	return "available"
}

// GetINetworks implements cloudprovider.ICloudWire.
func (s *SWire) GetINetworks() ([]cloudprovider.ICloudNetwork, error) {
	var subnets []SNetwork
	var query = map[string]string{
		"vpcId":      s.vpc.Id,
		"zoneName":   s.zone.ZoneName,
		"subnetType": "BCC",
		"maxKeys":    "100",
	}
	// res, err := s.vpc.region.doList1(ServiceSubnet, "/v1/subnet", query,&subnets)
	// if err != nil {
	// 	return nil, err
	// }

	// var ok bool
	// subnets,ok = res.([]SNetwork)

	var marker string
	for {
		var r subnetResp
		if len(marker) > 0 {
			query["marker"] = marker
		}
		res, err := s.vpc.region.doList(ServiceSubnet, "/v1/subnet", query)
		if err != nil {
			return nil,err
		}
		if err := res.Unmarshal(&r); err != nil {
			return nil,err
		}
		subnets = append(subnets, r.Networks...)
		if !r.IsTruncated {
			break
		}
		marker = r.NextMarker
	}

	var ret []cloudprovider.ICloudNetwork
	for i := range subnets {
		subnets[i].wire = s
		subnets[i].wire.vpc = s.vpc
		ret = append(ret, &subnets[i])
	}

	return ret, nil
}
