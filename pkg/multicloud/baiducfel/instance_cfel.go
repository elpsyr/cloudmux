package baiducfel

import (
	"context"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

var _ cloudprovider.ICfelCloudVM = (*SInstance)(nil)

func (s *SInstance) RebootVM(ctx context.Context) error {
	var params = map[string]interface{}{
		"forceStop": true,
	}
	_, err := s.region.doPut(ServiceInstance, "/v2/instance/"+s.ID+"?reboot", params)
	return err
}

func (s *SInstance) GetVncUrl() (string, error) {
	res, err := s.region.doGetWithoutVal(ServiceInstance, "/v2/instance/"+s.ID+"/vnc", nil)
	if err != nil {
		return "", nil
	}
	url, err := res.GetString("vncUrl")
	if err != nil {
		return "", nil
	}
	return url, nil
}
