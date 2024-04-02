package cloudpods

import (
	"context"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	compute "yunion.io/x/onecloud/pkg/apis/compute"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
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

	err = cloudprovider.WaitStatus(instance, api.VM_READY, 10*time.Second, 300*time.Second) // 5mintues
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
	return cloudprovider.WaitStatus(instance, api.VM_RUNNING, 10*time.Second, 300*time.Second) // 5mintues
}

func (self *SRegion) GetHostInstance(instanceId string) (cloudprovider.ICloudVM, error) {
	instance, err := self.GetIVMById(instanceId)
	if err != nil {
		log.Errorf("GetIVMById: %s", err)
		return instance, err
	}
	host, err := self.GetIHostById(instance.GetIHostId())
	if err != nil {
		log.Errorf("GetHost err: %s", err)
		return instance, err
	}
	// add host
	instance, err = host.GetIVMById(instanceId)

	if err != nil {
		log.Errorf("GetIVMById err: %s", err)
		return instance, err
	}

	return instance, err
}

func (self *SInstance) GetMonitorData(start, end string) ([]cloudprovider.ICfelMonitorData, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SRegion) CreateBareMetal(opts *cloudprovider.CfelSManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	hypervisor := api.HYPERVISOR_BAREMETAL
	ins, err := self.cfelCreateInstance("", hypervisor, opts)
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (self *SRegion) CreateVM(opts *cloudprovider.CfelSManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	hypervisor := api.HYPERVISOR_KVM
	ins, err := self.cfelCreateInstance("", hypervisor, opts)
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (self *SRegion) ResetGuestPassword(params *cloudprovider.CfelResetGuestPasswordOption) (cloudprovider.ICloudVM, error) {
	instance, err := self.GetIVMById(params.GuestID)
	if err != nil {
		log.Errorf("Fail to GetIVMById on ResetGuestPassword: %s", err)
		return instance, err
	}
	param := map[string]interface{}{
		"reset_password": params.ResetPassword,
		"auto_start":     params.AutoStart,
		"password":       params.Password,
		"username":       params.UserName,
	}
	var vm SInstance
	res, err := self.perform(&modules.Servers, params.GuestID, "set-password", jsonutils.Marshal(param))
	if err != nil {
		log.Errorf("Fail  ResetGuestPassword: %s", err)
		return instance, err
	}
	return &vm, res.Unmarshal(&vm)
}

func (self *SRegion) PingQga(guestId string, timeout int) (bool, error) {
	param := map[string]interface{}{"timeout": timeout}
	res, err := self.perform(&modules.Servers, guestId, "qga-ping", jsonutils.Marshal(param))
	if err != nil {
		log.Errorf("Fail PingQga: %s", err)
		return false, err
	}
	return res.IsZero(), nil
}

func (self *SRegion) CfelAttachDisk(instanceId, diskId string) error {
	// input := api.ServerAttachDiskInput{}
	// input.DiskId = diskId
	input := map[string]interface{}{
		"disk_id": diskId,
	}
	_, err := self.perform(&modules.Servers, instanceId, "attachdisk", input)
	return err
}

func (self *SRegion) CfelDetachDisk(instanceId, diskId string) error {
	input := map[string]interface{}{
		"disk_id":   diskId,
		"keep_disk": true,
	}
	_, err := self.perform(&modules.Servers, instanceId, "detachdisk", input)
	return err
}

func (self *SRegion) CfelInstanceSettingChange(id string,opts *cloudprovider.CfelChangeSettingOption) error {
	
	_,err := modules.Servers.Update(self.cli.s,id,jsonutils.Marshal(opts))
	return err
}



func (self *SRegion) cfelCreateInstance(hostId, hypervisor string, opts *cloudprovider.CfelSManagedVMCreateConfig) (*SInstance, error) {
	input := compute.ServerCreateInput{
		ServerConfigs: &compute.ServerConfigs{},
	}
	var isolatedDevice []*compute.IsolatedDeviceConfig
	if opts.IsolatedDevice != nil {
		for _,v := range opts.IsolatedDevice {
			isolatedDevice = append(isolatedDevice, &compute.IsolatedDeviceConfig{
				DevType:      v.DevType,
				Model:        v.Model,
				Vendor:       v.Vendor,
			})
		}
		
	}
	input.Name = opts.Name
	input.Hostname = opts.Hostname
	input.Description = opts.Description
	input.InstanceType = opts.InstanceType
	input.VcpuCount = opts.Cpu
	input.VmemSize = opts.MemoryMB
	input.Password = opts.Password
	input.PublicIpBw = opts.PublicIpBw
	input.PublicIpChargeType = string(opts.PublicIpChargeType)
	input.ProjectId = opts.ProjectId
	input.Metadata = opts.Tags
	input.UserData = opts.UserData
	input.PreferHost = hostId
	input.Hypervisor = hypervisor
	input.AutoStart = true
	input.DisableDelete = new(bool)
	input.IsolatedDevices = isolatedDevice
	if len(input.UserData) > 0 {
		input.EnableCloudInit = true
	}
	input.Secgroups = opts.ExternalSecgroupIds
	if opts.BillingCycle != nil {
		input.Duration = opts.BillingCycle.String()
	}
	input.Disks = append(input.Disks, &compute.DiskConfig{
		Index:    0,
		ImageId:  opts.ExternalImageId,
		DiskType: api.DISK_TYPE_SYS,
		SizeMb:   opts.SysDisk.SizeGB * 1024,
		Backend:  opts.SysDisk.StorageType,
		Storage:  opts.SysDisk.StorageExternalId,
	})
	for idx, disk := range opts.DataDisks {
		input.Disks = append(input.Disks, &compute.DiskConfig{
			Index:    idx + 1,
			DiskType: api.DISK_TYPE_DATA,
			SizeMb:   disk.SizeGB * 1024,
			Backend:  disk.StorageType,
			Storage:  disk.StorageExternalId,
		})
	}
	input.Networks = append(input.Networks, &compute.NetworkConfig{
		Index:   0,
		Network: opts.ExternalNetworkId,
		Address: opts.IpAddr,
	})
	ins := &SInstance{}
	return ins, self.create(&modules.Servers, input, ins)
}
