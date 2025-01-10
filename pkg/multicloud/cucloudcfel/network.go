package cucloudcfel

import (
	"fmt"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/util/rbacscope"
)

type SNetwork struct {
	multicloud.SNetworkBase
	CuCloudTags

	wire *SWire

	NetworkID           string `json:"networkId"`
	NetworkType         string `json:"networkType"`
	NetworkUUID         string `json:"networkUuid"`
	StatusEn            string `json:"statusEn"`
	SubNetworkCidr      string `json:"subNetworkCidr"`
	SubNetworkGatewayIP string `json:"subNetworkGatewayIp"`
	SubNetworkID        string `json:"subNetworkId"`
	SubNetworkIPRange   string `json:"subNetworkIpRange"`
	SubNetworkIPVersion string `json:"subNetworkIpVersion"`
	SubNetworkName      string `json:"subNetworkName"`
	SubNetworkUUID      string `json:"subNetworkUuid"`
	VlanPlanID          string `json:"vlanPlanId"`
	VpcID               string `json:"vpcId"`
	ZoneCode            string `json:"zoneCode"`
	ZoneID              string `json:"zoneId"`
}

// Delete implements cloudprovider.ICloudNetwork.
func (s *SNetwork) Delete() error {
	params := map[string]interface{}{
		"cloudRegionCode": s.wire.vpc.region.GetId(),
	}
	return s.wire.vpc.region.client.delete(fmt.Sprintf("/instance/v1/product/subnets/%s", s.VpcID), params)
}

// GetAllocTimeoutSeconds implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetAllocTimeoutSeconds() int {
	return 0
}

// GetDescription implements cloudprovider.ICloudNetwork.
// Subtle: this method shadows the method (SNetworkBase).GetDescription of SNetwork.SNetworkBase.
func (s *SNetwork) GetDescription() string {
	return ""
}

// GetGateway implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetGateway() string {
	return s.SubNetworkGatewayIP
}

// GetGlobalId implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetGlobalId() string {
	return s.SubNetworkID
}

// GetIWire implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetIWire() cloudprovider.ICloudWire {
	return s.wire
}

// GetId implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetId() string {
	return s.SubNetworkID
}

// GetIpEnd implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetIpEnd() string {
	arr := strings.Split(s.SubNetworkIPRange, "~")
	if len(arr) == 2 {
		return arr[1]
	}
	return ""
}

// GetIpMask implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetIpMask() int8 {
	return 0
}

// GetIpStart implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetIpStart() string {
	arr := strings.Split(s.SubNetworkIPRange, "~")
	if len(arr) == 2 {
		return arr[0]
	}
	return ""
}

// GetName implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetName() string {
	return s.SubNetworkName
}

// GetProjectId implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetProjectId() string {
	return ""
}

// GetPublicScope implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetPublicScope() rbacscope.TRbacScope {
	return ""
}

// GetServerType implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetServerType() string {
	return ""
}

// GetStatus implements cloudprovider.ICloudNetwork.
func (s *SNetwork) GetStatus() string {
	return "ready"
}

var _ cloudprovider.ICloudNetwork = (*SNetwork)(nil)
