package cloudpods

import (
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type SLoadbalancerAcl struct {
	multicloud.SVirtualResourceBase
	CloudpodsTags

	region *SRegion

	AclEntries []AclEntry
	CreatedAt  time.Time
	Deleted    bool
	Freezed    bool
	ID         string
	Name       string
	Region     string
	RegionID   string
	Status     string
	UpdatedAt  time.Time
}

type AclEntry struct {
	Comment string
	Cidr    string
}

// Delete implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) Delete() error {
	return s.region.cli.delete(&modules.LoadbalancerAcls,s.ID)
}

// GetCreatedAt implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetAclEntries implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) GetAclEntries() []cloudprovider.SLoadbalancerAccessControlListEntry {
	var res []cloudprovider.SLoadbalancerAccessControlListEntry
	for _, val := range s.AclEntries {
		res = append(res, cloudprovider.SLoadbalancerAccessControlListEntry{
			CIDR:    val.Cidr,
			Comment: val.Comment,
		})
	}
	return res
}

// GetAclListenerID implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) GetAclListenerID() string {
	return ""
}

// GetGlobalId implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) GetGlobalId() string {
	return s.ID
}

// GetId implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) GetId() string {
	return s.ID
}

// GetName implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) GetName() string {
	return s.Name
}

// GetStatus implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) GetStatus() string {
	return s.Status
}

// Sync implements cloudprovider.ICloudLoadbalancerAcl.
func (s *SLoadbalancerAcl) Sync(acl *cloudprovider.SLoadbalancerAccessControlList) error {
	panic("unimplemented")
}

var _ cloudprovider.ICloudLoadbalancerAcl = (*SLoadbalancerAcl)(nil)

func (self *SRegion) CreateILoadBalancerAcl(acl *cloudprovider.SLoadbalancerAccessControlList) (cloudprovider.ICloudLoadbalancerAcl, error) {
	entry := make([]map[string]string, 0, len(acl.Entrys))
	for _, val := range acl.Entrys {
		entry = append(entry, map[string]string{"cidr": val.CIDR, "comment": val.Comment})
	}
	params := map[string]interface{}{
		"disable_delete": false,
		"name":           acl.Name,
		"acl_entries":    entry,
	}
	var res SLoadbalancerAcl
	err := self.create(&modules.LoadbalancerAcls, params, &res)

	return &res, err
}

func (self *SRegion) GetILoadBalancerAcls() ([]cloudprovider.ICloudLoadbalancerAcl, error) {
	params := map[string]interface{}{}
	var ret []SLoadbalancerAcl
	err := self.list(&modules.LoadbalancerAcls, params, &ret)

	if err != nil {
		return nil, err
	}
	var res []cloudprovider.ICloudLoadbalancerAcl
	for i := range ret {
		ret[i].region = self
		res = append(res, &ret[i])
	}
	return res, nil
}

func (self *SRegion) GetILoadBalancerAclById(aclId string) (cloudprovider.ICloudLoadbalancerAcl, error) {
	params := map[string]interface{}{}
	var ret SLoadbalancerAcl
	res, err := modules.LoadbalancerAcls.GetById(self.cli.s, aclId, jsonutils.Marshal(params))

	if err != nil {
		return nil, err
	}
	ret.region = self
	return &ret, res.Unmarshal(&ret)
}
