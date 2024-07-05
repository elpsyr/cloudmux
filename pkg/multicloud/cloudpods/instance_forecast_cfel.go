package cloudpods

import (
	"strings"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/apis"
	"yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/modulebase"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
)

type SchedulerManager struct {
	modulebase.ResourceManager
}

var (
	Schedulers SchedulerManager
)

func init() {
	Schedulers = SchedulerManager{modules.NewResourceManager(apis.SERVICE_TYPE_SCHEDULER, "scheduler", "",
		[]string{},
		[]string{}, "v1")}

	modules.RegisterCompute(&Schedulers)
}

func (this *SchedulerManager) DoForecast(s *mcclient.ClientSession, params jsonutils.JSONObject) (jsonutils.JSONObject, error) {
	return this.PerformClassAction(s, "forecast", params)
}

func (self *SRegion) SchedulerForecast(hypervisor string, opts *cloudprovider.CfelSManagedVMCreateConfig) (jsonutils.JSONObject, error) {
	return self.preCreateInstance("", hypervisor, opts)
}

func (self *SRegion) preCreateInstance(hostId, hypervisor string, opts *cloudprovider.CfelSManagedVMCreateConfig) (jsonutils.JSONObject, error) {
	input := compute.ServerCreateInput{
		ServerConfigs: &compute.ServerConfigs{},
	}
	// 数量
	if opts.Count >= 0 {
		input.Count = opts.Count
	} else {
		input.Count = 1 // default 1
	}

	var isolatedDevice []*compute.IsolatedDeviceConfig
	if opts.IsolatedDevice != nil {
		for _, v := range opts.IsolatedDevice {
			isolatedDevice = append(isolatedDevice, &compute.IsolatedDeviceConfig{
				DevType: v.DevType,
				Model:   v.Model,
				Vendor:  v.Vendor,
			})
		}

	}
	input.GenerateName = opts.Name
	// input.Name = opts.Name
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

	if opts.EipBw > 0 {
		input.EipBw = opts.EipBw
		input.EipAutoDellocate = opts.EipAutoDellocate
	}
	// system disk
	sysDiskSize := 0
	if opts.SysDisk.SizeGB > 0 {
		sysDiskSize = opts.SysDisk.SizeGB * 1024
	} else {
		sysDiskSize = -1
	}
	if hypervisor == api.HYPERVISOR_KVM {
		arr := strings.Split(opts.SysDisk.StorageType, "_")
		input.Disks = append(input.Disks, &compute.DiskConfig{
			Index:    0,
			ImageId:  opts.ExternalImageId,
			DiskType: api.DISK_TYPE_SYS,
			SizeMb:   sysDiskSize,
			Backend:  arr[0],
			Storage:  opts.SysDisk.StorageExternalId,
			Medium:   arr[1],
		})
	} else if hypervisor == api.HYPERVISOR_BAREMETAL {
		input.Disks = append(input.Disks, &compute.DiskConfig{
			Index:    0,
			ImageId:  opts.ExternalImageId,
			DiskType: api.DISK_TYPE_SYS,
			SizeMb:   sysDiskSize,
			Backend:  opts.SysDisk.StorageType,
			Storage:  opts.SysDisk.StorageExternalId,
		})
	}

	if hypervisor == api.HYPERVISOR_KVM {
		for idx, disk := range opts.DataDisks {
			arr := strings.Split(disk.StorageType, "_")
			input.Disks = append(input.Disks, &compute.DiskConfig{
				Index:    idx + 1,
				DiskType: api.DISK_TYPE_DATA,
				SizeMb:   disk.SizeGB * 1024,
				Backend:  arr[0],
				Storage:  disk.StorageExternalId,
				Medium:   arr[1],
			})
		}
	}
	if opts.ExternalNetworkId != "" {
		input.Networks = append(input.Networks, &compute.NetworkConfig{
			Index:   0,
			Network: opts.ExternalNetworkId,
			Address: opts.IpAddr,
		})
	} else if len(opts.Networks) != 0 {
		for i := 0; i < len(opts.Networks); i++ {
			input.Networks = append(input.Networks, &compute.NetworkConfig{
				Index:          i,
				Network:        opts.Networks[i].NetworkId,
				Address:        opts.Networks[i].Address,
				RequireTeaming: opts.Networks[i].RequireTeaming,
			})
		}
	}
	// raid
	ebmDiskConfig := make([]*compute.BaremetalDiskConfig, 0)
	for _, config := range opts.BaremetalDiskConfigs {
		ebmDiskConfig = append(ebmDiskConfig, &compute.BaremetalDiskConfig{
			Type:         config.Type,
			Conf:         config.Conf,
			Count:        config.Count,
			Range:        config.Range,
			Splits:       config.Splits,
			Size:         config.Size,
			Adapter:      config.Adapter,
			Driver:       config.Driver,
			Cachedbadbbu: config.Cachedbadbbu,
			Strip:        config.Strip,
			RA:           config.RA,
			WT:           config.WT,
			Direct:       config.Direct,
		})
	}
	input.BaremetalDiskConfigs = ebmDiskConfig

	return Schedulers.DoForecast(self.cli.s, jsonutils.Marshal(input))
}
