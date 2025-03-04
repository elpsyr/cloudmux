package qcloud

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (s *SLBBackendGroup) CfelBackendServerPortAndWeight(opts *cloudprovider.SCfelUpdateServerPortAndWeight) (cloudprovider.ICloudLoadbalancerBackend, error) {
	var (
		portErr   error
		wigthErr  error
		requestId string
	)
	var server *SLBBackend
	backends, err := s.GetBackends()
	if err != nil {
		return nil, err
	}
	for _, val := range backends {
		if val.GetId() == opts.ServerId {
			server = &val
			break
		}
	}
	if server == nil {
		return nil, cloudprovider.ErrNotFound
	}

	if opts.Port > 0 && opts.Port != server.Port {
		requestId, portErr = s.UpdateBackendServerWeight("ModifyTargetPort", server.InstanceId, opts.LocationId, opts.Weight, opts.Port, server.Port)
		if portErr != nil {
			return nil, portErr
		}
		portErr = s.lb.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
		if portErr != nil {
			return nil, portErr
		}
		server.Port = opts.Port
	}
	if opts.Weight != server.Weight {
		requestId, wigthErr = s.UpdateBackendServerWeight("ModifyTargetWeight", server.InstanceId, opts.LocationId, opts.Weight, 0, server.Port)
		if wigthErr != nil {
			return nil, wigthErr
		}
		wigthErr = s.lb.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
		if wigthErr != nil {
			return nil, wigthErr
		}
		server.Weight = opts.Weight
	}

	return server, nil
}

func (self *SLBBackendGroup) CfelAddBackendServer(opts *cloudprovider.SCfelLBListenerAddServer) (cloudprovider.ICloudLoadbalancerBackend, error) {
	backends, err := self.GetBackends()
	if err != nil {
		return nil, err
	}
	for _, backend := range backends {
		if backend.InstanceId == opts.ServerId && backend.Port == opts.Port {
			return &backend, nil
		}
	}

	tag := md5Str(uuid.NewString())
	requestId, err := self.appLBBackendServerCfel("RegisterTargets", opts.LocationId, opts.ServerId, tag, opts.Weight, opts.Port)
	if err != nil {
		return nil, err
	}
	err = self.lb.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
	if err != nil {
		return nil, err
	}
	backends, err = self.GetBackends()
	if err != nil {
		return nil, err
	}
	for _, backend := range backends {
		if strings.HasSuffix(backend.GetGlobalId(), fmt.Sprintf("%s/%s", opts.ServerId, tag)) {
			return &backend, nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SLBBackendGroup) CfelRemoveBackendServer(opts *cloudprovider.SCfelLBListenerRemoveServer) error {
	// 传进来的serverId格式 fmt.Sprintf("%s/%s/%d", self.group.GetId(), self.InstanceId, self.Port)
	// 要分割出真的serverId
	arr := strings.Split(opts.ServerId, "/")
	if len(arr) < 2 {
		return fmt.Errorf("serverId format error")
	}
	backends, err := self.GetBackends()
	if err != nil {
		return err
	}
	var found bool
	for _, backend := range backends {
		if backend.InstanceId == arr[1] && backend.Port == opts.Port {
			found = true
		}
	}

	if found {
		requestId, err := self.appLBBackendServerCfel("DeregisterTargets", opts.LocationId, arr[1], "", opts.Weight, opts.Port)
		if err != nil {
			if strings.Contains(err.Error(), "not registered") {
				return nil
			}
			return err
		}
		return self.lb.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
	}
	return nil
}

func md5Str(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:8])
}

func (self *SLBBackendGroup) appLBBackendServerCfel(action string, locationId string, serverId string, tag string, weight int, port int) (string, error) {
	params := map[string]string{
		"LoadBalancerId":       self.lb.LoadBalancerId,
		"ListenerId":           self.listener.ListenerId,
		"Targets.0.InstanceId": serverId,
		"Targets.0.Port":       strconv.Itoa(port),
		"Targets.0.Weight":     strconv.Itoa(weight),
	}
	if action == "RegisterTargets" && len(tag) > 0 {
		params["Targets.0.Tag"] = tag
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
