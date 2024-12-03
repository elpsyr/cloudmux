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
