package baiducfel

import (
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SEip struct {
	multicloud.SEipBase
	BaiduTags `json:"baidu_tags,omitempty"`

	Eip             string    `json:"eip,omitempty"`
	ID              string    `json:"eipId,omitempty"`
	BandwidthInMbps int64     `json:"bandwidthInMbps,omitempty"`
	CreateTime      time.Time `json:"createTime,omitempty"`
	EipInstanceType string    `json:"eipInstanceType,omitempty"`
	ExpireTime      time.Time `json:"expireTime,omitempty"`
	InstanceID      string    `json:"instanceId,omitempty"`
	InstanceType    string    `json:"instanceType,omitempty"`
	Name            string    `json:"name,omitempty"`
	PaymentTiming   string    `json:"paymentTiming,omitempty"`
	Region          string    `json:"region,omitempty"`
	RouteType       string    `json:"routeType,omitempty"`
	ShareGroupID    string    `json:"shareGroupId,omitempty"`
	Status          string    `json:"status,omitempty"`
	BillingMethod   string    `json:"billingMethod,omitempty"`
	PrivateIp       string    `json:"privateIp,omitempty"`
}

// Associate implements cloudprovider.ICloudEIP.
func (s *SEip) Associate(conf *cloudprovider.AssociateConfig) error {
	panic("unimplemented")
}

// ChangeBandwidth implements cloudprovider.ICloudEIP.
func (s *SEip) ChangeBandwidth(bw int) error {
	panic("unimplemented")
}

// Delete implements cloudprovider.ICloudEIP.
func (s *SEip) Delete() error {
	panic("unimplemented")
}

// Dissociate implements cloudprovider.ICloudEIP.
func (s *SEip) Dissociate() error {
	panic("unimplemented")
}

// GetAssociationExternalId implements cloudprovider.ICloudEIP.
func (s *SEip) GetAssociationExternalId() string {
	return s.InstanceID
}

// GetAssociationType implements cloudprovider.ICloudEIP.
func (s *SEip) GetAssociationType() string {
	return ""
}

// GetBandwidth implements cloudprovider.ICloudEIP.
func (s *SEip) GetBandwidth() int {
	return int(s.BandwidthInMbps)
}

// GetCreatedAt implements cloudprovider.ICloudEIP.
// Subtle: this method shadows the method (SEipBase).GetCreatedAt of SEip.SEipBase.
func (s *SEip) GetCreatedAt() time.Time {
	return s.CreateTime
}

// GetDescription implements cloudprovider.ICloudEIP.
// Subtle: this method shadows the method (SEipBase).GetDescription of SEip.SEipBase.
func (s *SEip) GetDescription() string {
	return ""
}

// GetExpiredAt implements cloudprovider.ICloudEIP.
// Subtle: this method shadows the method (SEipBase).GetExpiredAt of SEip.SEipBase.
func (s *SEip) GetExpiredAt() time.Time {
	return s.ExpireTime
}

// GetGlobalId implements cloudprovider.ICloudEIP.
func (s *SEip) GetGlobalId() string {
	return s.ID
}

// GetINetworkId implements cloudprovider.ICloudEIP.
// Subtle: this method shadows the method (SEipBase).GetINetworkId of SEip.SEipBase.
func (s *SEip) GetINetworkId() string {
	panic("unimplemented")
}

// GetId implements cloudprovider.ICloudEIP.
func (s *SEip) GetId() string {
	return s.ID
}

// GetInternetChargeType implements cloudprovider.ICloudEIP.
func (s *SEip) GetInternetChargeType() string {
	return ""
}

// GetIpAddr implements cloudprovider.ICloudEIP.
func (s *SEip) GetIpAddr() string {
	return s.Eip
}

// GetMode implements cloudprovider.ICloudEIP.
func (s *SEip) GetMode() string {
	return s.BillingMethod
}

// GetName implements cloudprovider.ICloudEIP.
func (s *SEip) GetName() string {
	return s.Name
}

// GetProjectId implements cloudprovider.ICloudEIP.
func (s *SEip) GetProjectId() string {
	return ""
}

// GetStatus implements cloudprovider.ICloudEIP.
func (s *SEip) GetStatus() string {
	return s.Status
}

var _ cloudprovider.ICloudEIP = (*SEip)(nil)
