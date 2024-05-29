package cloudpods

import (
	"context"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type SCloudLoadbalancerListener struct {
	multicloud.SResourceBase
	CloudpodsTags

	// SCloudloadbalancerHealthCheck

	// SCloudLoadbalancerRedirect

	loadbalancer *SLoadbalancer

	ACLID                 string    `json:"acl_id"`
	ACLStatus             string    `json:"acl_status"`
	ACLType               string    `json:"acl_type"`
	BackendConnectTimeout int       `json:"backend_connect_timeout"`
	BackendGroupID        string    `json:"backend_group_id"`
	BackendIdleTimeout    int       `json:"backend_idle_timeout"`
	BackendServerPort     int       `json:"backend_server_port"`
	CertificateID         string    `json:"certificate_id"`
	ClientIdleTimeout     int       `json:"client_idle_timeout"`
	ClientRequestTimeout  int       `json:"client_request_timeout"`
	CreatedAt             time.Time `json:"created_at"`
	Deleted               bool      `json:"deleted"`
	Description           string    `json:"description"`
	EgressMbps            int       `json:"egress_mbps"`
	EnableHTTP2           bool      `json:"enable_http2"`
	Gzip                  bool      `json:"gzip"`

	HTTPRequestRate       int    `json:"http_request_rate"`
	HTTPRequestRatePerSrc int    `json:"http_request_rate_per_src"`
	ID                    string `json:"id"`
	// ImportedAt                 time.Time `json:"imported_at"`
	// IsEmulated                 bool      `json:"is_emulated"`
	ListenerPort   int    `json:"listener_port"`
	ListenerType   string `json:"listener_type"`
	LoadbalancerID string `json:"loadbalancer_id"`
	Name           string `json:"name"`
	Progress       int    `json:"progress"`

	HealthCheck         string `json:"health_check"`
	HealthCheckDomain   string `json:"health_check_domain"`
	HealthCheckFall     int    `json:"health_check_fall"`
	HealthCheckHTTPCode string `json:"health_check_http_code"`
	HealthCheckInterval int    `json:"health_check_interval"`
	HealthCheckRise     int    `json:"health_check_rise"`
	HealthCheckTimeout  int    `json:"health_check_timeout"`
	HealthCheckType     string `json:"health_check_type"`
	HealthCheckURI      string `json:"health_check_uri"`

	Scheduler                  string    `json:"scheduler"`
	SendProxy                  string    `json:"send_proxy"`
	Source                     string    `json:"source"`
	Status                     string    `json:"status"`
	StickySession              string    `json:"sticky_session"`
	StickySessionType          string    `json:"sticky_session_type"`
	StickySessionCookieTimeout int       `json:"sticky_session_cookie_timeout"`
	TLSCipherPolicy            string    `json:"tls_cipher_policy"`
	UpdateVersion              int       `json:"update_version"`
	UpdatedAt                  time.Time `json:"updated_at"`
	XforwardedFor              bool      `json:"xforwarded_for"`
}

type SCloudLoadbalancerListenerRule struct {
	multicloud.SResourceBase
	CloudpodsTags

	listener *SCloudLoadbalancerListener

	BackendGroupID        string    `json:"backend_group_id"`
	CreatedAt             time.Time `json:"created_at"`
	Deleted               bool      `json:"deleted"`
	Domain                string    `json:"domain"`
	HealthCheckFall       int       `json:"health_check_fall"`
	HealthCheckInterval   int       `json:"health_check_interval"`
	HealthCheckRise       int       `json:"health_check_rise"`
	HealthCheckTimeout    int       `json:"health_check_timeout"`
	HTTPRequestRate       int       `json:"http_request_rate"`
	HTTPRequestRatePerSrc int       `json:"http_request_rate_per_src"`
	ID                    string    `json:"id"`
	ListenerID            string    `json:"listener_id"`
	Default               bool      `json:"is_default"`
	Name                  string    `json:"name"`
	Path                  string    `json:"path"`
	Progress              int       `json:"progress"`
	Redirect              string    `json:"redirect"`
	RedirectCode          int       `json:"redirect_code"`
	RedirectHost          string    `json:"redirect_host"`
	RedirectPath          string    `json:"redirect_path"`
	RedirectScheme        string    `json:"redirect_scheme"`

	Status string `json:"status"`

	UpdatedAt time.Time `json:"updated_at"`
}

// Delete implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) Delete(ctx context.Context) error {
	return s.listener.loadbalancer.region.cli.delete(&modules.LoadbalancerListenerRules, s.ID)
}

// GetBackendGroupId implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetBackendGroupId() string {
	return s.BackendGroupID
}

// GetCondition implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetCondition() string {
	return ""
}

// GetCreatedAt implements cloudprovider.ICloudLoadbalancerListenerRule.
// Subtle: this method shadows the method (SResourceBase).GetCreatedAt of SCloudLoadbalancerListenerRule.SResourceBase.
func (s *SCloudLoadbalancerListenerRule) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetDomain implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetDomain() string {
	return s.Domain
}

// GetGlobalId implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetGlobalId() string {
	return s.ID
}

// GetId implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetId() string {
	return s.ID
}

// GetName implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetName() string {
	return s.Name
}

// GetPath implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetPath() string {
	return s.Path
}

// GetRedirect implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetRedirect() string {
	return s.Redirect
}

// GetRedirectCode implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetRedirectCode() int64 {
	return int64(s.RedirectCode)
}

// GetRedirectHost implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetRedirectHost() string {
	return s.RedirectHost
}

// GetRedirectPath implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetRedirectPath() string {
	return s.RedirectPath
}

// GetRedirectScheme implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetRedirectScheme() string {
	return s.RedirectScheme
}

// GetStatus implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) GetStatus() string {
	return s.Status
}

// IsDefault implements cloudprovider.ICloudLoadbalancerListenerRule.
func (s *SCloudLoadbalancerListenerRule) IsDefault() bool {
	return s.Default
}

var _ cloudprovider.ICloudLoadbalancerListenerRule = (*SCloudLoadbalancerListenerRule)(nil)

// CfelCreateILoadBalancerListenerRule implements cloudprovider.ICfelLoadbalancerListener.
func (s *SCloudLoadbalancerListener) CfelCreateILoadBalancerListenerRule(rule *cloudprovider.SCfelLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	params := map[string]interface{}{
		// "listener_type": "http",
		"name":                      rule.Name,
		"domain":                    rule.Domain,
		"path":                      rule.Path,
		"backend_group":             rule.BackendGroupId,
		"http_request_rate":         rule.HttpRequestRate,
		"http_request_rate_per_src": rule.HttpRequestRatePerSrc,
		"listener":                  rule.ListenerId,

		"redirect":        rule.Redirect, // off raw
		"redirect_code":   rule.RedirectCode,
		"redirect_scheme": rule.RedirectScheme,
		"redirect_host":   rule.RedirectHost,
		"redirect_path":   rule.RedirectPath,
	}
	var res SCloudLoadbalancerListenerRule

	err := s.loadbalancer.region.cli.create(&modules.LoadbalancerListenerRules, params, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetCreatedAt implements cloudprovider.ICloudLoadbalancerListener.
// Subtle: this method shadows the method (SResourceBase).GetCreatedAt of SCloudLoadbalancerListener.SResourceBase.
func (s *SCloudLoadbalancerListener) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetDescription implements cloudprovider.ICloudLoadbalancerListener.
// Subtle: this method shadows the method (SResourceBase).GetDescription of SCloudLoadbalancerListener.SResourceBase.
func (s *SCloudLoadbalancerListener) GetDescription() string {
	return s.Description
}

// GetHealthCheck implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheck() string {
	return s.HealthCheck
}

// GetHealthCheckCode implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckCode() string {
	return s.HealthCheckHTTPCode
}

// GetHealthCheckDomain implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckDomain() string {
	return s.HealthCheckDomain
}

// GetHealthCheckExp implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckExp() string {
	return s.GetHealthCheckExp()
}

// GetHealthCheckFail implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckFail() int {
	return s.HealthCheckFall
}

// GetHealthCheckInterval implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckInterval() int {
	return s.HealthCheckInterval
}

// GetHealthCheckReq implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckReq() string {
	return ""
}

// GetHealthCheckRise implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckRise() int {
	return 0
}

// GetHealthCheckTimeout implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckTimeout() int {
	return s.HealthCheckTimeout
}

// GetHealthCheckType implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckType() string {
	return s.HealthCheckType
}

// GetHealthCheckURI implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetHealthCheckURI() string {
	return s.HealthCheckURI
}

// GetRedirect implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetRedirect() string {
	panic("unimplemented")
}

// GetRedirectCode implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetRedirectCode() int64 {
	panic("unimplemented")
}

// GetRedirectHost implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetRedirectHost() string {
	panic("unimplemented")
}

// GetRedirectPath implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetRedirectPath() string {
	panic("unimplemented")
}

// GetRedirectScheme implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetRedirectScheme() string {
	panic("unimplemented")
}

// ChangeCertificate implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) ChangeCertificate(ctx context.Context, opts *cloudprovider.ListenerCertificateOptions) error {
	panic("unimplemented")
}

// ChangeScheduler implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) ChangeScheduler(ctx context.Context, opts *cloudprovider.ChangeListenerSchedulerOptions) error {
	panic("unimplemented")
}

// CreateILoadBalancerListenerRule implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) CreateILoadBalancerListenerRule(rule *cloudprovider.SLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	panic("unimplemented")
}

// Delete implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) Delete(ctx context.Context) error {
	return s.loadbalancer.region.cli.delete(&modules.LoadbalancerListeners, s.ID)
}

// GetAclId implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetAclId() string {
	return s.ACLID
}

// GetAclStatus implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetAclStatus() string {
	return s.ACLStatus
}

// GetAclType implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetAclType() string {
	return s.ACLType
}

// GetBackendConnectTimeout implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetBackendConnectTimeout() int {
	return s.BackendConnectTimeout
}

// GetBackendGroupId implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetBackendGroupId() string {
	return s.BackendGroupID
}

// GetBackendServerPort implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetBackendServerPort() int {
	return s.BackendServerPort
}

// GetCertificateId implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetCertificateId() string {
	return s.CertificateID
}

// GetClientIdleTimeout implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetClientIdleTimeout() int {
	return s.ClientIdleTimeout
}

// GetEgressMbps implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetEgressMbps() int {
	return s.EgressMbps
}

// GetGlobalId implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetGlobalId() string {
	return s.ID
}

// GetILoadBalancerListenerRuleById implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetILoadBalancerListenerRuleById(ruleId string) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	params := map[string]string{
		"listener": s.ID,
	}
	var res SCloudLoadbalancerListenerRule
	err := s.loadbalancer.region.cli.get(&modules.LoadbalancerListenerRules, ruleId, params, &res)
	if err != nil {
		return nil, err
	}
	res.listener = s
	return &res, nil
}

// GetILoadbalancerListenerRules implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetILoadbalancerListenerRules() ([]cloudprovider.ICloudLoadbalancerListenerRule, error) {
	params := map[string]interface{}{
		"listener": s.ID,
	}
	var res []SCloudLoadbalancerListenerRule
	err := s.loadbalancer.region.cli.list(&modules.LoadbalancerListenerRules, params, &res)
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICloudLoadbalancerListenerRule
	for i := range res {
		res[i].listener = s
		ret = append(ret, &res[i])
	}
	return ret, nil
}

// GetId implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetId() string {
	return s.ID
}

// GetListenerPort implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetListenerPort() int {
	return s.ListenerPort
}

// GetListenerType implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetListenerType() string {
	return s.ListenerType
}

// GetName implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetName() string {
	return s.Name
}

// GetScheduler implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetScheduler() string {
	return s.Scheduler
}

// GetStatus implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetStatus() string {
	return s.Status
}

// GetStickySession implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetStickySession() string {
	return s.StickySession
}

// GetStickySessionCookie implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetStickySessionCookie() string {
	return s.StickySession
}

// GetStickySessionCookieTimeout implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetStickySessionCookieTimeout() int {
	return s.StickySessionCookieTimeout
}

// GetStickySessionType implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetStickySessionType() string {
	return s.StickySessionType
}

// GetTLSCipherPolicy implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GetTLSCipherPolicy() string {
	return s.TLSCipherPolicy
}

// GzipEnabled implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) GzipEnabled() bool {
	return s.Gzip
}

// HTTP2Enabled implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) HTTP2Enabled() bool {
	return s.EnableHTTP2
}

// SetAcl implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) SetAcl(ctx context.Context, opts *cloudprovider.ListenerAclOptions) error {
	params := map[string]interface{}{
		"acl_status": opts.AclStatus,
		"acl_type":   opts.AclType,
		"acl":        opts.AclId,
	}

	_, err := modules.LoadbalancerListeners.Update(s.loadbalancer.region.cli.s, s.ID, jsonutils.Marshal(params))
	// _,err := s.loadbalancer.region.perform(&modules.LoadbalancerListeners,s.ID,"",params)
	return err

}

// SetHealthCheck implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) SetHealthCheck(ctx context.Context, opts *cloudprovider.ListenerHealthCheckOptions) error {
	panic("unimplemented")
}

// Start implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) Start() error {
	params := map[string]interface{}{
		"status": "enabled",
	}
	_, err := s.loadbalancer.region.cli.perform(&modules.LoadbalancerListeners, s.ID, "status", params)
	return err
}

// Stop implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) Stop() error {
	params := map[string]interface{}{
		"status": "disabled",
	}
	_, err := s.loadbalancer.region.cli.perform(&modules.LoadbalancerListeners, s.ID, "status", params)
	return err
}

// XForwardedForEnabled implements cloudprovider.ICloudLoadbalancerListener.
func (s *SCloudLoadbalancerListener) XForwardedForEnabled() bool {
	return s.XforwardedFor
}

var _ cloudprovider.ICloudLoadbalancerListener = (*SCloudLoadbalancerListener)(nil)

var _ cloudprovider.ICfelLoadbalancerListener = (*SCloudLoadbalancerListener)(nil)
