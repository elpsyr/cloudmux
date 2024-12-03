package baiducfel

import "yunion.io/x/cloudmux/pkg/cloudprovider"

type SInstanceNic struct {
	cloudprovider.DummyICloudNic
	BaiduTags

	Az             string        `json:"az"`
	CreatedTime    string        `json:"createdTime"`
	Description    string        `json:"description"`
	DeviceID       string        `json:"deviceId"`
	EniID          string        `json:"eniId"`
	EniNum         int           `json:"eniNum"`
	EniUUID        string        `json:"eniUuid"`
	EriInfos       []interface{} `json:"eriInfos"`
	EriNum         int           `json:"eriNum"`
	Ips            []SEip        `json:"ips"`
	MacAddress     string        `json:"macAddress"`
	Name           string        `json:"name"`
	SecurityGroups []string      `json:"securityGroups"`
	Status         string        `json:"status"`
	SubnetID       string        `json:"subnetId"`
	SubnetType     string        `json:"subnetType"`
	Type           string        `json:"type"`
	VpcID          string        `json:"vpcId"`
}

var _ cloudprovider.ICloudNic = (*SInstanceNic)(nil)

// AssignAddress implements cloudprovider.ICloudNic.
func (s *SInstanceNic) AssignAddress(ipAddrs []string) error {
	panic("unimplemented")
}

// AssignNAddress implements cloudprovider.ICloudNic.
func (s *SInstanceNic) AssignNAddress(count int) ([]string, error) {
	panic("unimplemented")
}

// GetDriver implements cloudprovider.ICloudNic.
func (s *SInstanceNic) GetDriver() string {
	return ""
}

// GetINetworkId implements cloudprovider.ICloudNic.
func (s *SInstanceNic) GetINetworkId() string {
	return s.SubnetID
}

// GetIP implements cloudprovider.ICloudNic.
func (s *SInstanceNic) GetIP() string {
	if len(s.Ips) > 0 {
		return s.Ips[0].PrivateIp
	}
	return ""
}

// GetIP6 implements cloudprovider.ICloudNic.
func (s *SInstanceNic) GetIP6() string {
	return ""
}

// GetId implements cloudprovider.ICloudNic.
func (s *SInstanceNic) GetId() string {
	return s.DeviceID
}

// GetMAC implements cloudprovider.ICloudNic.
func (s *SInstanceNic) GetMAC() string {
	return s.MacAddress
}

// GetSubAddress implements cloudprovider.ICloudNic.
func (s *SInstanceNic) GetSubAddress() ([]string, error) {
	panic("unimplemented")
}

// InClassicNetwork implements cloudprovider.ICloudNic.
func (s *SInstanceNic) InClassicNetwork() bool {
	return false
}

// UnassignAddress implements cloudprovider.ICloudNic.
func (s *SInstanceNic) UnassignAddress(ipAddrs []string) error {
	panic("unimplemented")
}
