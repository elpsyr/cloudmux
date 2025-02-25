package qcloud

import (
	"strconv"
	"strings"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (self *SLBListener) CfelCreateILoadBalancerListenerRule(*cloudprovider.SCfelLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	return nil, nil
}
func (self *SLBListener) CfelUpdateILoadBalancerListenerRule(rule *cloudprovider.SCfelUpdateLoadbalancerListenerRule) error {
	hc := getListenerRuleHealthCheck(&rule.SLoadbalancerListenerRule)

	switch strings.ToLower(rule.Scheduler) {
	case api.LB_SCHEDULER_WRR:
		rule.Scheduler = "WRR"
	case api.LB_SCHEDULER_WLC:
		rule.Scheduler = "LEAST_CONN"
	case api.LB_SCHEDULER_SCH:
		rule.Scheduler = "IP_HASH"
	}

	params := map[string]string{
		"LoadBalancerId": self.lb.LoadBalancerId,
		"ListenerId":     self.ListenerId,
		// "Domain":         rule.Domain,
		"Url":        rule.Path,
		"LocationId": rule.LocationId,
	}

	params["Scheduler"] = rule.Scheduler
	params["SessionExpireTime"] = strconv.Itoa(rule.StickySessionCookieTimeout)

	// health check
	params["HealthCheck.HealthSwitch"] = strconv.Itoa(hc.HealthSwitch)
	params["HealthCheck.TimeOut"] = strconv.Itoa(hc.TimeOut)
	params["HealthCheck.IntervalTime"] = strconv.Itoa(hc.IntervalTime)
	params["HealthCheck.HealthNum"] = strconv.Itoa(hc.HealthNum)
	params["HealthCheck.UnHealthNum"] = strconv.Itoa(hc.UnHealthNum)

	if hc.HTTPCode > 0 {
		params["HealthCheck.HttpCode"] = strconv.Itoa(hc.HTTPCode)
		params["HealthCheck.HttpCheckPath"] = hc.HTTPCheckPath
		params["HealthCheck.HttpCheckDomain"] = hc.HTTPCheckDomain
		params["HealthCheck.HttpCheckMethod"] = hc.HTTPCheckMethod
	}

	resp, err := self.lb.region.clbRequest("ModifyRule", params)
	if err != nil {
		return err
	}

	requestId, err := resp.GetString("RequestId")
	if err != nil {
		return err
	}

	return self.lb.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)

}

func (self *SLBListener) Update(*cloudprovider.SLoadbalancerListenerCreateOptions) error {
	return nil
}
