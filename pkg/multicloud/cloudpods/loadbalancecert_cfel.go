package cloudpods

import (
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type SCloudLoadbalancerCertificate struct {
	multicloud.SVirtualResourceBase
	CloudpodsTags

	region *SRegion

	Certificate string
	CreatedAt   time.Time
	Deleted     bool
	Freezed     bool
	ID          string
	Name        string
	PrivateKey  string
	Region      string
	RegionID    string
	Status      string
	UpdatedAt   time.Time
	Fingerprint string
	IsPublic    bool
	SubjectAlternativeNames string
	CommonName string
	NotAfter    time.Time
	NotBefore   time.Time
}

// GetSubjectAlternativeNames implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetSubjectAlternativeNames() string {
	return s.SubjectAlternativeNames
}

// Delete implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) Delete() error {
	return s.region.cli.delete(&modules.LoadbalancerCertificates,s.ID)
}

// GetCreatedAt implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetCommonName implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetCommonName() string {
	return s.CommonName
}

// GetExpireTime implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetExpireTime() time.Time {
	return s.NotAfter
}

// GetFingerprint implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetFingerprint() string {
	return s.Fingerprint
}

// GetGlobalId implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetGlobalId() string {
	return s.ID
}

// GetId implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetId() string {
	return s.ID
}

// GetName implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetName() string {
	return s.Name
}

// GetPrivateKey implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetPrivateKey() string {
	return s.PrivateKey
}

// GetPublickKey implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetPublickKey() string {
	return ""
}

// GetStatus implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) GetStatus() string {
	return s.Status
}

// Sync implements cloudprovider.ICloudLoadbalancerCertificate.
func (s *SCloudLoadbalancerCertificate) Sync(name string, privateKey string, publickKey string) error {
	panic("unimplemented")
}

var _ cloudprovider.ICloudLoadbalancerCertificate = (*SCloudLoadbalancerCertificate)(nil)

func (self *SRegion) CreateILoadBalancerCertificate(cert *cloudprovider.SLoadbalancerCertificate) (cloudprovider.ICloudLoadbalancerCertificate, error) {
	params := map[string]interface{}{
		"disable_delete": false,
		"name":           cert.Name,
		"certificate":    cert.Certificate,
		"private_key":    cert.PrivateKey,
	}
	var res SCloudLoadbalancerCertificate
	err := self.create(&modules.LoadbalancerCertificates, params, &res)

	return &res, err
}

func (self *SRegion) GetILoadBalancerCertificates() ([]cloudprovider.ICloudLoadbalancerCertificate, error) {
	params := map[string]interface{}{}
	var ret []SCloudLoadbalancerCertificate
	err := self.list(&modules.LoadbalancerAcls, params, &ret)

	if err != nil {
		return nil, err
	}
	var res []cloudprovider.ICloudLoadbalancerCertificate
	for i := range ret {
		ret[i].region = self
		res = append(res, &ret[i])
	}
	return res, nil
}

func (self *SRegion) GetILoadBalancerCertificateById(certId string) (cloudprovider.ICloudLoadbalancerCertificate, error) {
	params := map[string]interface{}{}
	var ret SCloudLoadbalancerCertificate
	res, err := modules.LoadbalancerCertificates.GetById(self.cli.s, certId, jsonutils.Marshal(params))

	if err != nil {
		return nil, err
	}
	ret.region = self
	return &ret, res.Unmarshal(&ret)
}
