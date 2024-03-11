package cloudpods

import (
	"context"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type SLoadbalancer struct {
	multicloud.SLoadbalancerBase
	CloudpodsTags

	region *SRegion

	Address       string    `json:"address"`
	AddressType   string    `json:"address_type"`
	Cluster       string    `json:"cluster"`
	ClusterID     string    `json:"cluster_id"`
	CreatedAt     time.Time `json:"created_at"`
	Deleted       bool      `json:"deleted"`
	DeletedAt     time.Time `json:"deleted_at"`
	Description   string    `json:"description"`
	DisableDelete bool      `json:"disable_delete"`
	ID            string    `json:"id"`
	IsDefaultVpc  bool      `json:"is_default_vpc"`
	Name          string    `json:"name"`
	Network       string    `json:"network"`
	NetworkID     string    `json:"network_id"`
	NetworkType   string    `json:"network_type"`
	Region        string    `json:"region"`
	RegionID      string    `json:"region_id"`
	Status        string    `json:"status"`
	Tenant        string    `json:"tenant"`
	Vpc           string    `json:"vpc"`
	VpcID         string    `json:"vpc_id"`
	Zone          string    `json:"zone"`
	ZoneID        string    `json:"zone_id"`
}

// CfelCreateILoadBalancerBackendGroup implements cloudprovider.ICfelLoadbalancer.
func (s *SLoadbalancer) CfelCreateILoadBalancerBackendGroup(bg *cloudprovider.SCfelLoadbalancerBackendGroup) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	params := map[string]interface{}{
		"name":         bg.Name,
		"loadbalancer": bg.LoadBalancerId,
	}
	res, err := modules.LoadbalancerBackendGroups.Create(s.region.cli.s, jsonutils.Marshal(params))
	if err != nil {
		return nil, err
	}
	var ret SCloudLoadbalancerBackendGroup

	return &ret, res.Unmarshal(&ret)
}

// CreateILoadBalancerBackendGroup implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) CreateILoadBalancerBackendGroup(group *cloudprovider.SLoadbalancerBackendGroup) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	params := map[string]interface{}{
		"name":         group.Name,
		"loadbalancer": s.ID,
	}
	var ret SCloudLoadbalancerBackendGroup
	err := s.region.create(&modules.LoadbalancerBackendGroups, params, &ret)
	return &ret, err
}

// CreateILoadBalancerListener implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) CreateILoadBalancerListener(ctx context.Context, listener *cloudprovider.SLoadbalancerListenerCreateOptions) (cloudprovider.ICloudLoadbalancerListener, error) {
	params := map[string]interface{}{
		"acl_id":     listener.AccessControlListId,
		"acl_status": listener.AccessControlListStatus,
		"acl_type":   listener.AccessControlListType,
		// "backend_connect_timeout":       5,
		"backend_group":          listener.BackendGroupId,
		"backend_idle_timeout":   90,
		"certificate_id":         listener.CertificateId,
		"client_idle_timeout":    90,
		"client_request_timeout": 10,
		"description":            listener.Description,
		"enable_http2":           listener.EnableHTTP2,
		"gzip":                   listener.Gzip,
		"health_check":           listener.HealthCheck,
		"health_check_fall":      listener.HealthCheckFail,
		"health_check_http_code": listener.HealthCheckHttpCode,
		"health_check_interval":  listener.HealthCheckInterval,
		"health_check_path":      listener.HealthCheckURI,
		"health_check_rise":      listener.HealthCheckRise,
		"health_check_timeout":   listener.HealthCheckTimeout,
		"health_check_type":      listener.HealthCheckType,
		"health_check_domain":    listener.HealthCheckDomain,
		"health_check_uri":       listener.HealthCheckURI,
		"health_check_req":       listener.HealthCheckReq,
		"health_check_exp":       listener.HealthCheckExp,

		"http_request_rate":     0,
		"http_request_rate_src": 0,
		"listener_port":         listener.ListenerPort,
		"listener_type":         listener.ListenerType,
		"loadbalancer":          s.Name,
		"loadbalancer_id":       s.ID,
		"name":                  listener.Name,
		"redirect":              "off",

		"scheduler": listener.Scheduler,
		// "send_proxy":                    "",

		"sticky_session": listener.StickySession,

		"tls_cipher_policy": listener.TLSCipherPolicy,
		"x_forwarded_for":   listener.XForwardedFor,
	}
	if (listener.ListenerType == "http" || listener.ListenerType == "https") && listener.StickySession == "on" {
		params["sticky_session_cookie"] = listener.StickySessionCookie
		params["sticky_session_cookie_timeout"] = listener.StickySessionCookieTimeout
		params["sticky_session_type"] = listener.StickySessionType
	}
	var ret SCloudLoadbalancerListener
	err := s.region.create(&modules.LoadbalancerListeners, params, &ret)
	return &ret, err
}

// Delete implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) Delete(ctx context.Context) error {
	return s.region.cli.delete(&modules.Loadbalancers, s.ID)
}

// GetAddress implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetAddress() string {
	return s.Address
}

// GetAddressType implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetAddressType() string {
	return s.AddressType
}

// GetChargeType implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetChargeType() string {
	return ""
}

// GetEgressMbps implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetEgressMbps() int {
	return 0
}

// GetGlobalId implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetGlobalId implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetGlobalId() string {
	return s.ID
}

// GetILoadBalancerBackendGroupById implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetILoadBalancerBackendGroupById(groupId string) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	params := map[string]interface{}{}
	var ret SCloudLoadbalancerBackendGroup
	res, err := modules.LoadbalancerBackendGroups.GetById(s.region.cli.s, groupId, jsonutils.Marshal(params))
	if err != nil {
		return nil, err
	}
	ret.loadbalancer = s
	return &ret, res.Unmarshal(&ret)
}

// GetILoadBalancerBackendGroups implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetILoadBalancerBackendGroups() ([]cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	params := map[string]interface{}{
		"loadbalancer": s.ID,
	}
	var ret []SCloudLoadbalancerBackendGroup
	err := s.region.list(&modules.LoadbalancerBackendGroups, params, &ret)
	if err != nil {
		return nil, err
	}
	var res []cloudprovider.ICloudLoadbalancerBackendGroup
	for i := range ret {
		res = append(res, &ret[i])
	}
	return res, nil
}

// GetILoadBalancerListenerById implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetILoadBalancerListenerById(listenerId string) (cloudprovider.ICloudLoadbalancerListener, error) {
	params := map[string]interface{}{}
	res, err := modules.LoadbalancerListeners.GetById(s.region.cli.s, listenerId, jsonutils.Marshal(params))
	if err != nil {
		return nil, err
	}
	var ret SCloudLoadbalancerListener
	ret.loadbalancer = s
	return &ret, res.Unmarshal(&ret)
}

// GetILoadBalancerListeners implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetILoadBalancerListeners() ([]cloudprovider.ICloudLoadbalancerListener, error) {
	params := map[string]interface{}{
		"loadbalancer": s.ID,
	}
	var ret []SCloudLoadbalancerListener
	err := s.region.list(&modules.LoadbalancerListeners, params, &ret)
	if err != nil {
		return nil, err
	}
	var res []cloudprovider.ICloudLoadbalancerListener
	for i := range ret {
		ret[i].loadbalancer = s
		res = append(res, &ret[i])
	}
	return res, nil
}

// GetId implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetId() string {
	return s.ID
}

// GetLoadbalancerSpec implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetLoadbalancerSpec() string {
	return ""
}

// GetName implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetName() string {
	return s.Name
}

// GetNetworkIds implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetNetworkIds() []string {
	return []string{s.NetworkID}
}

// GetNetworkType implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetNetworkType() string {
	return s.NetworkType
}

// GetProjectId implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetProjectId() string {
	return ""
}

// GetStatus implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetStatus() string {
	return s.Status
}

// GetVpcId implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetVpcId() string {
	return s.VpcID
}

// GetZone1Id implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetZone1Id() string {
	return ""
}

// GetZoneId implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) GetZoneId() string {
	return s.ZoneID
}

// Start implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) Start() error {
	params := map[string]interface{}{
		"id":     s.ID,
		"status": "enabled",
	}
	_, err := s.region.cli.perform(&modules.Loadbalancers, s.ID, "status", params)
	return err
}

// Stop implements cloudprovider.ICloudLoadbalancer.
func (s *SLoadbalancer) Stop() error {
	params := map[string]interface{}{
		"id":     s.ID,
		"status": "disabled",
	}
	_, err := s.region.cli.perform(&modules.Loadbalancers, s.ID, "status", params)
	return err
}

var _ cloudprovider.ICloudLoadbalancer = (*SLoadbalancer)(nil)
var _ cloudprovider.ICfelLoadbalancer = (*SLoadbalancer)(nil)

func (self *SRegion) CreateILoadBalancer(loadbalancer *cloudprovider.SLoadbalancerCreateOptions) (cloudprovider.ICloudLoadbalancer, error) {
	params := map[string]interface{}{
		// "cluster_id":  "",
		"__meta__":       loadbalancer.Tags,
		"disable_delete": false,
		"name":           loadbalancer.Name,
		"vpc":            loadbalancer.VpcId,
		"network":        loadbalancer.NetworkIds[0],
		"description":    loadbalancer.Desc,
		"zone_id":        loadbalancer.ZoneId,
		"eip_id":         loadbalancer.EipId,
		// "project": "3be335319d384121800c1c7da7fae686",
	}
	ret, err := modules.Loadbalancers.Create(self.cli.s, jsonutils.Marshal(params))
	if err != nil {
		return nil, err
	}
	var lb SLoadbalancer

	return &lb, ret.Unmarshal(&lb)
}

func (self *SRegion) GetILoadBalancers() ([]cloudprovider.ICloudLoadbalancer, error) {
	params := map[string]interface{}{}
	var ret []SLoadbalancer
	err := self.list(&modules.Loadbalancers, params, &ret)

	if err != nil {
		return nil, err
	}
	var res []cloudprovider.ICloudLoadbalancer
	for i := range ret {
		ret[i].region = self
		res = append(res, &ret[i])
	}
	return res, nil
}

func (self *SRegion) GetILoadBalancerById(loadbalancerId string) (cloudprovider.ICloudLoadbalancer, error) {
	params := map[string]interface{}{}
	var ret SLoadbalancer
	ret.region = self
	res, err := modules.Loadbalancers.GetById(self.cli.s, loadbalancerId, jsonutils.Marshal(params))

	if err != nil {
		return nil, err
	}
	return &ret, res.Unmarshal(&ret)
}
