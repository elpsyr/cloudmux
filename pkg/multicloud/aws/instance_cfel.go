package aws

import (
	"context"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (self *SInstance) GetIHostId() string {
	return "-" + self.Placement.AvailabilityZone
}

func (self *SInstance) RebootVM(ctx context.Context) error {
	err := self.host.zone.region.RebootVM(self.InstanceId)
	if err != nil {
		return err
	}
	return cloudprovider.WaitStatus(self, api.VM_RUNNING, 10*time.Second, 300*time.Second) // 5mintues
}

func (self *SRegion) RebootVM(instanceId string) error {
	params := map[string]string{
		"InstanceId.1": instanceId,
	}
	ret := struct{}{}
	return self.ec2Request("RebootInstances", params, &ret)
}
