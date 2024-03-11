package cloudpods

import (
	"context"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type SCloudLoadbalancerBackend struct {
	multicloud.SVirtualResourceBase
	CloudpodsTags

	backendGroup *SCloudLoadbalancerBackendGroup

	Address        string    `json:"address"`
	BackendGroupID string    `json:"backend_group_id"`
	BackendID      string    `json:"backend_id"`
	BackendType    string    `json:"backend_type"`
	CreatedAt      time.Time `json:"created_at"`
	Deleted        bool      `json:"deleted"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Port           int       `json:"port"`
	SendProxy      string    `json:"send_proxy"`
	Ssl            string    `json:"ssl"`
	Status         string    `json:"status"`
	UpdatedAt      time.Time `json:"updated_at"`
	Weight         int       `json:"weight"`
}

// GetCreatedAt implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetBackendId implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetBackendId() string {
	return s.BackendGroupID
}

// GetBackendRole implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetBackendRole() string {
	return ""
}

// GetBackendType implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetBackendType() string {
	return s.BackendType
}

// GetGlobalId implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetGlobalId() string {
	return s.ID
}

// GetId implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetId() string {
	return s.ID
}

// GetIpAddress implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetIpAddress() string {
	return s.Address
}

// GetName implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetName() string {
	return s.Name
}

// GetPort implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetPort() int {
	return s.Port
}

// GetStatus implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetStatus() string {
	return s.Status
}

// GetWeight implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) GetWeight() int {
	return s.Weight
}

// SyncConf implements cloudprovider.ICloudLoadbalancerBackend.
func (s *SCloudLoadbalancerBackend) SyncConf(ctx context.Context, port int, weight int) error {
	params := map[string]interface{}{
		"weight": weight,
		"port":   port,
	}
	return s.backendGroup.loadbalancer.region.cli.update(&modules.LoadbalancerBackends, s.ID, params)
}

var _ cloudprovider.ICloudLoadbalancerBackend = (*SCloudLoadbalancerBackend)(nil)

type SCloudLoadbalancerBackendGroup struct {
	multicloud.SResourceBase
	CloudpodsTags

	loadbalancer *SLoadbalancer

	CreatedAt      time.Time `json:"created_at,omitempty"`
	Deleted        bool      `json:"deleted,omitempty"`
	ID             string    `json:"id,omitempty"`
	Default        bool      `json:"is_default,omitempty"`
	Loadbalancer   string    `json:"loadbalancer,omitempty"`
	LoadbalancerID string    `json:"loadbalancer_id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Region         string    `json:"region,omitempty"`
	RegionID       string    `json:"region_id,omitempty"`
	Status         string    `json:"status,omitempty"`
	Type           string    `json:"type,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	Vpc            string    `json:"vpc,omitempty"`
	VpcID          string    `json:"vpc_id,omitempty"`
	Zone           string    `json:"zone,omitempty"`
	ZoneID         string    `json:"zone_id,omitempty"`
}

// CfelAddBackendServer implements cloudprovider.ICfelLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) CfelAddBackendServer(serverType string, serverId string, ssl string, weight int, port int) (cloudprovider.ICloudLoadbalancerBackend, error) {
	params := map[string]interface{}{
		"backend_type":  serverType,
		"guest_backend": serverId,
		"port":          port,
		"weight":        weight,
		"ssl":           ssl,
		"backend_group": s.ID,
		"backend":       serverId,
	}
	var ret SCloudLoadbalancerBackend
	err := s.loadbalancer.region.create(&modules.LoadbalancerBackends, params, &ret)
	return &ret, err
}

// AddBackendServer implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) AddBackendServer(serverId string, weight int, port int) (cloudprovider.ICloudLoadbalancerBackend, error) {
	panic("unimplemented")
}

// Delete implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) Delete(ctx context.Context) error {
	return s.loadbalancer.region.cli.delete(&modules.LoadbalancerBackendGroups, s.ID)
}

// GetGlobalId implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) GetGlobalId() string {
	return s.ID
}

// GetILoadbalancerBackendById implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) GetILoadbalancerBackendById(backendId string) (cloudprovider.ICloudLoadbalancerBackend, error) {
	res, err := modules.LoadbalancerBackends.GetById(s.loadbalancer.region.cli.s, backendId, nil)
	if err != nil {
		return nil, err
	}
	var ret SCloudLoadbalancerBackend
	ret.backendGroup = s
	return &ret, res.Unmarshal(&ret)
}

// GetILoadbalancerBackends implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) GetILoadbalancerBackends() ([]cloudprovider.ICloudLoadbalancerBackend, error) {
	params := map[string]interface{}{
		"backend_group": s.ID,
	}
	var ret []SCloudLoadbalancerBackend
	err := s.loadbalancer.region.cli.list(&modules.LoadbalancerBackends, params, &ret)
	if err != nil {
		return nil, nil
	}
	var res []cloudprovider.ICloudLoadbalancerBackend
	for i := range ret {
		ret[i].backendGroup = s
		res = append(res, &ret[i])
	}
	return res, nil
}

// GetCreatedAt implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetId implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) GetId() string {
	return s.ID
}

// GetName implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) GetName() string {
	return s.Name
}

// GetStatus implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) GetStatus() string {
	return s.Status
}

// GetType implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) GetType() string {
	return s.Type
}

// IsDefault implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) IsDefault() bool {
	return s.Default
}

// RemoveBackendServer implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) RemoveBackendServer(serverId string, weight int, port int) error {
	return s.loadbalancer.region.cli.delete(&modules.LoadbalancerBackends, serverId)
}

// Sync implements cloudprovider.ICloudLoadbalancerBackendGroup.
func (s *SCloudLoadbalancerBackendGroup) Sync(ctx context.Context, group *cloudprovider.SLoadbalancerBackendGroup) error {
	panic("unimplemented")
}

var _ cloudprovider.ICloudLoadbalancerBackendGroup = (*SCloudLoadbalancerBackendGroup)(nil)
var _ cloudprovider.ICfelLoadbalancerBackendGroup = (*SCloudLoadbalancerBackendGroup)(nil)
