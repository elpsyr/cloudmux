package ecloudcfel

import (
	"context"
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

var _ cloudprovider.ICfelCloudVM = (*SInstance)(nil)

func (in *SInstance) RebootVM(ctx context.Context) error {
	req := NewConsoleRequest(in.host.zone.region.ID, fmt.Sprintf("/api/openapi-ecs/acl/v3/server/%s/reboot", in.Id), nil, nil)
	_, err := in.host.zone.region.client.doPut(req)
	return err
}

func (s *SInstance) GetVncUrl() (string, error) {
	// https://ecloud.10086.cn/op-help-center/doc/article/58481
	req := NewConsoleRequest(s.region.ID, fmt.Sprintf("/api/openapi-ecs/acl/v3/server/%s/vnc", s.Id), nil, nil)

	res, err := s.region.client.doGetWithoutVal(context.Background(), req)
	if err != nil {
		return "", nil
	}

	return res.Interface().(string), err
}
