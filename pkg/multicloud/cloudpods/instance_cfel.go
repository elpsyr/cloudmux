package cloudpods

import (
	"context"
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
)

var _ cloudprovider.ICfelCloudVM = (*SInstance)(nil)

func (self *SInstance) RebootVM(ctx context.Context) error {
	err := self.host.zone.region.RebootVM(self.Id)
	if err != nil {
		return err
	}
	return cloudprovider.WaitStatus(self, api.VM_RUNNING, 10*time.Second, 300*time.Second) // 5mintues
}

func (self *SRegion) RebootVM(instanceId string) error {
	instance, err := self.GetHostInstance(instanceId)
	if err != nil {
		log.Errorf("Fail to GetHostInstance : %s", err)
		return err
	}
	status := instance.GetStatus()

	if status != api.VM_RUNNING {
		log.Errorf("RebootVM: vm status is %s expect %s", status, api.VM_RUNNING)
		return cloudprovider.ErrInvalidStatus
	}
	err = instance.StopVM(context.Background(), &cloudprovider.ServerStopOptions{
		IsForce: true,
	})
	if err != nil {
		log.Errorf("Fail to RebootVM  , first step StopVM err : %s", err)
		return err
	}

	err = cloudprovider.WaitStatus(self, api.VM_READY, 10*time.Second, 300*time.Second) // 5mintues
	if err != nil {
		log.Errorf("Fail to RebootVM  , first step StopVM failed : %s", err)
		return err
	}
	err = instance.StartVM(context.Background())
	if err != nil {
		log.Errorf("Fail to RebootVM  , second step StartVM err : %s", err)
		return err
	}
	if err != nil {
		return err
	}
	return cloudprovider.WaitStatus(self, api.VM_RUNNING, 10*time.Second, 300*time.Second) // 5mintues
}

func (self *SRegion) GetHostInstance(instanceId string) (cloudprovider.ICloudVM, error) {
	instance, err := self.GetIVMById(instanceId)
	if err != nil {
		log.Errorf("Fail to GetInstance on RebootVM: %s", err)
		return instance, err
	}
	host, err := self.GetHost(instance.GetIHostId())
	if err != nil {
		log.Errorf("Fail to GetHost on RebootVM: %s", err)
		return instance, err
	}
	// add host
	instance, err = host.GetIVMById(instanceId)

	if err != nil {
		log.Errorf("Fail to get instance status on RebootVM: %s", err)
		return instance, err
	}

	return instance, err
}

func (self *SInstance) GetMonitorData(start, end string) ([]cloudprovider.ICfelMonitorData, error) {
	return nil, cloudprovider.ErrNotImplemented
}
