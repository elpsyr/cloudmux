package qcloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (self *SLBBackendGroup) CfelAddBackendServer(opts *cloudprovider.SCfelLBListenerAddServer) (cloudprovider.ICloudLoadbalancerBackend, error) {
	requestId, err := self.appLBBackendServerCfel("RegisterTargets",opts.LocationId, opts.ServerId, opts.Weight, opts.Port)
	if err != nil {
		return nil, err
	}
	err = self.lb.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
	if err != nil {
		return nil, err
	}
	backends, err := self.GetBackends()
	if err != nil {
		return nil, err
	}
	for _, backend := range backends {
		if strings.HasSuffix(backend.GetId(), fmt.Sprintf("%s-%d", opts.ServerId, opts.Port)) {
			return &backend, nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SLBBackendGroup) CfelRemoveBackendServer(opts *cloudprovider.SCfelLBListenerRemoveServer) error {
	requestId, err := self.appLBBackendServerCfel("DeregisterTargets", opts.LocationId, opts.ServerId, opts.Weight, opts.Port)
	if err != nil {
		if strings.Contains(err.Error(), "not registered") {
			return nil
		}
		return err
	}
	return self.lb.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
}

func (self *SLBBackendGroup) appLBBackendServerCfel(action string, locationId string, serverId string, weight int, port int) (string, error) {
	params := map[string]string{
		"LoadBalancerId":       self.lb.LoadBalancerId,
		"ListenerId":           self.listener.ListenerId,
		"Targets.0.InstanceId": serverId,
		"Targets.0.Port":       strconv.Itoa(port),
		"Targets.0.Weight":     strconv.Itoa(weight),
	}
	if len(locationId) > 0 {
		params["LocationId"] = locationId
	}
	resp, err := self.lb.region.clbRequest(action, params)
	if err != nil {
		return "", err
	}

	return resp.GetString("RequestId")
}
