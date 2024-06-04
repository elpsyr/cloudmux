package aliyun

import (
	"fmt"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

var _ cloudprovider.ICfelLoadbalancerListener = (*SLoadbalancerTCPListener)(nil)
var _ cloudprovider.ICfelLoadbalancerListener = (*SLoadbalancerUDPListener)(nil)
var _ cloudprovider.ICfelLoadbalancerListener = (*SLoadbalancerHTTPListener)(nil)
var _ cloudprovider.ICfelLoadbalancerListener = (*SLoadbalancerHTTPSListener)(nil)

// Update implements cloudprovider.ICfelLoadbalancerListener.
func (self *SLoadbalancerTCPListener) Update(opts *cloudprovider.SLoadbalancerListenerCreateOptions) error {
	params := map[string]string{
		"LoadBalancerId":    self.lb.LoadBalancerId,
		"ListenerPort":      fmt.Sprintf("%d", self.ListenerPort),
		"AclStatus":         "off",
		"Scheduler":         opts.Scheduler,
		"HealthCheckSwitch": "off",
	}
	if opts.AccessControlListStatus == api.LB_BOOL_ON {
		params["AclStatus"] = "on"
		params["AclType"] = opts.AccessControlListType
		params["AclId"] = opts.AccessControlListId
	}
	if opts.HealthCheck == api.LB_BOOL_ON {
		params["HealthCheckSwitch"] = "on"
		if opts.HealthCheckTimeout >= 1 && opts.HealthCheckTimeout <= 300 {
			params["HealthCheckConnectTimeout"] = fmt.Sprintf("%d", opts.HealthCheckTimeout)
		}

		params["HealthCheckDomain"] = opts.HealthCheckDomain
		params["HealthCheckType"] = opts.HealthCheckType
		params["HealthCheckHttpCode"] = opts.HealthCheckHttpCode
		params["HealthCheckURI"] = opts.HealthCheckURI

		if opts.HealthCheckRise >= 2 && opts.HealthCheckRise <= 10 {
			params["HealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckRise)
		}

		if opts.HealthCheckFail >= 2 && opts.HealthCheckFail <= 10 {
			params["UnhealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckFail)
		}
		if opts.HealthCheckInterval >= 1 && opts.HealthCheckInterval <= 50 {
			params["healthCheckInterval"] = fmt.Sprintf("%d", opts.HealthCheckInterval)
		}
	}
	_, err := self.lb.region.lbRequest("SetLoadBalancerTCPListenerAttribute", params)
	return err
}

// Update implements cloudprovider.ICfelLoadbalancerListener.
func (self *SLoadbalancerHTTPSListener) Update(opts *cloudprovider.SLoadbalancerListenerCreateOptions) error {
	params := map[string]string{
		"LoadBalancerId":    self.lb.LoadBalancerId,
		"ListenerPort":      fmt.Sprintf("%d", self.ListenerPort),
		"HealthCheckSwitch": "off",
		"AclStatus":         "off",
		"Scheduler":         opts.Scheduler,
		"ServerCertificateId": opts.CertificateId,
		"TLSCipherPolicy": opts.TLSCipherPolicy,
	}
	if opts.EnableHTTP2 {
		params["EnableHttp2"] = "on"
	} else {
		params["EnableHttp2"] = "off"
	}
	if opts.HealthCheck == api.LB_BOOL_ON {
		params["HealthCheckSwitch"] = "on"
		params["HealthCheckDomain"] = opts.HealthCheckDomain
		params["HealthCheckHttpCode"] = opts.HealthCheckHttpCode
		params["HealthCheckURI"] = opts.HealthCheckURI

		if opts.HealthCheckRise >= 2 && opts.HealthCheckRise <= 10 {
			params["HealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckRise)
		}

		if opts.HealthCheckFail >= 2 && opts.HealthCheckFail <= 10 {
			params["UnhealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckFail)
		}
		if opts.HealthCheckInterval >= 1 && opts.HealthCheckInterval <= 50 {
			params["healthCheckInterval"] = fmt.Sprintf("%d", opts.HealthCheckInterval)
		}
	}
	_, err := self.lb.region.lbRequest("SetLoadBalancerHTTPSListenerAttribute", params)
	return err
}

// Update implements cloudprovider.ICfelLoadbalancerListener.
func (self *SLoadbalancerHTTPListener) Update(opts *cloudprovider.SLoadbalancerListenerCreateOptions) error {
	params := map[string]string{
		"LoadBalancerId":    self.lb.LoadBalancerId,
		"ListenerPort":      fmt.Sprintf("%d", self.ListenerPort),
		"HealthCheckSwitch": "off",
		"Scheduler":         opts.Scheduler,
		"HealthCheck":opts.HealthCheck,
	}
	if opts.HealthCheck == api.LB_BOOL_ON {

		params["HealthCheckSwitch"] = "on"
		params["HealthCheckDomain"] = opts.HealthCheckDomain
		params["HealthCheckHttpCode"] = opts.HealthCheckHttpCode
		params["HealthCheckURI"] = opts.HealthCheckURI

		if opts.HealthCheckRise >= 2 && opts.HealthCheckRise <= 10 {
			params["HealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckRise)
		}

		if opts.HealthCheckFail >= 2 && opts.HealthCheckFail <= 10 {
			params["UnhealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckFail)
		}
		if opts.HealthCheckInterval >= 1 && opts.HealthCheckInterval <= 50 {
			params["healthCheckInterval"] = fmt.Sprintf("%d", opts.HealthCheckInterval)
		}
	}


	if opts.HealthCheck == "on" {
		
		params["HealthCheckURI"] = opts.HealthCheckURI
		//The HealthCheckTimeout parameter is required.
		if opts.HealthCheckTimeout < 1 || opts.HealthCheckTimeout > 300 {
			opts.HealthCheckTimeout = 5
		}
		params["HealthCheckTimeout"] = fmt.Sprintf("%d", opts.HealthCheckTimeout)
	}

	if opts.ClientRequestTimeout < 1 || opts.ClientRequestTimeout > 180 {
		opts.ClientRequestTimeout = 60
	}
	params["RequestTimeout"] = fmt.Sprintf("%d", opts.ClientRequestTimeout)

	if opts.ClientIdleTimeout < 1 || opts.ClientIdleTimeout > 60 {
		opts.ClientIdleTimeout = 15
	}
	params["IdleTimeout"] = fmt.Sprintf("%d", opts.ClientIdleTimeout)

	params["StickySession"] = opts.StickySession
	params["StickySessionType"] = opts.StickySessionType
	params["Cookie"] = opts.StickySessionCookie
	if opts.StickySessionCookieTimeout < 1 || opts.StickySessionCookieTimeout > 86400 {
		opts.StickySessionCookieTimeout = 500
	}
	params["CookieTimeout"] = fmt.Sprintf("%d", opts.StickySessionCookieTimeout)
	//params["ForwardPort"] = fmt.Sprintf("%d", listener.ForwardPort) //暂不支持
	params["Gzip"] = "off"
	if opts.Gzip {
		params["Gzip"] = "on"
	}
	params["XForwardedFor"] = "off"
	if opts.XForwardedFor {
		params["XForwardedFor"] = "on"
	}

	_, err := self.lb.region.lbRequest("SetLoadBalancerHTTPListenerAttribute", params)
	return err
}

// Update implements cloudprovider.ICfelLoadbalancerListener.
func (self *SLoadbalancerUDPListener) Update(opts *cloudprovider.SLoadbalancerListenerCreateOptions) error {
	params := map[string]string{
		"LoadBalancerId":    self.lb.LoadBalancerId,
		"ListenerPort":      fmt.Sprintf("%d", self.ListenerPort),
		"HealthCheckSwitch": "off",
		"Scheduler":         opts.Scheduler,
	}
	if opts.HealthCheck == api.LB_BOOL_ON {
		params["HealthCheckSwitch"] = "on"
		if opts.HealthCheckTimeout >= 1 && opts.HealthCheckTimeout <= 300 {
			params["HealthCheckConnectTimeout"] = fmt.Sprintf("%d", opts.HealthCheckTimeout)
		}
		params["healthCheckReq"] = opts.HealthCheckReq
		params["healthCheckExp"] = opts.HealthCheckExp
		params["HealthCheckDomain"] = opts.HealthCheckDomain
		params["HealthCheckHttpCode"] = opts.HealthCheckHttpCode
		params["HealthCheckURI"] = opts.HealthCheckURI

		if opts.HealthCheckRise >= 2 && opts.HealthCheckRise <= 10 {
			params["HealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckRise)
		}

		if opts.HealthCheckFail >= 2 && opts.HealthCheckFail <= 10 {
			params["UnhealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckFail)
		}
		if opts.HealthCheckInterval >= 1 && opts.HealthCheckInterval <= 50 {
			params["healthCheckInterval"] = fmt.Sprintf("%d", opts.HealthCheckInterval)
		}
	}
	_, err := self.lb.region.lbRequest("SetLoadBalancerUDPListenerAttribute", params)
	return err
}

// CfelCreateILoadBalancerListenerRule implements cloudprovider.ICfelLoadbalancerListener.
func (listener *SLoadbalancerUDPListener) CfelCreateILoadBalancerListenerRule(*cloudprovider.SCfelLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	panic("unimplemented")
}

// CfelCreateILoadBalancerListenerRule implements cloudprovider.ICfelLoadbalancerListener.
func (listener *SLoadbalancerHTTPSListener) CfelCreateILoadBalancerListenerRule(*cloudprovider.SCfelLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	panic("unimplemented")
}

// CfelCreateILoadBalancerListenerRule implements cloudprovider.ICfelLoadbalancerListener.
func (listener *SLoadbalancerHTTPListener) CfelCreateILoadBalancerListenerRule(*cloudprovider.SCfelLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	panic("unimplemented")
}

// CfelCreateILoadBalancerListenerRule implements cloudprovider.ICfelLoadbalancerListener.
func (listener *SLoadbalancerTCPListener) CfelCreateILoadBalancerListenerRule(*cloudprovider.SCfelLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	panic("unimplemented")
}
