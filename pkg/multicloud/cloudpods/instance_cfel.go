package cloudpods

import (
	"context"
	"sync"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	compute "yunion.io/x/onecloud/pkg/apis/compute"
	monitor_input "yunion.io/x/onecloud/pkg/apis/monitor"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
	monitor "yunion.io/x/onecloud/pkg/mcclient/modules/monitor"
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

func (self *SInstance) GetCfelHypervisor() string {
	return self.Hypervisor
}

type Monitor struct {
	Series []struct {
		Columns []string    `json:"columns"`
		Name    string      `json:"name"`
		Points  [][]float64 `json:"points"`
		RawName string      `json:"raw_name"`
	} `json:"series"`
	SeriesTotal int                `json:"series_total"`
	Type        string             `json:"-"`
	Item        []*MonitorDataItem `json:"-"`
}

func (self *SRegion) GetMonitorData(vmId, start, end, interval string) ([]cloudprovider.ICfelMonitorData, []string, error) {
	tmp := map[string][]string{
		"vm_diskio": {"write_bps", "read_bps"},
		"vm_cpu":    {"usage_active"},
		"vm_mem":    {"used_percent"},
		"vm_disk":   {"used_percent"},
		"vm_netio":  {"bps_send", "bps_recv"},
	}
	var wg sync.WaitGroup
	var m sync.Mutex
	var faild []string

	var mute sync.Mutex
	var item []*MonitorDataItem

	// var mm sync.Mutex
	var ms []*Monitor
	// var ch = make(chan Monitor)
	for measure, val := range tmp {
		for _, v := range val {

			wg.Add(1)
			go func(p, vv string) {
				params := monitor_input.MetricQueryInput{
					From:     start,
					To:       end,
					Unit:     false,
					Interval: interval,
					MetricQuery: []*monitor_input.AlertQuery{{
						Model: monitor_input.MetricQuery{
							Measurement: p, // vm_diskio[write_bps,read_bps] vm_cpu[usage_active] vm_mem[used_percent] vm_disk[used_percent] vm_netio[bps_send,bps_recv]
							Tags: []monitor_input.MetricQueryTag{{
								Key:      "vm_id",
								Operator: "=",
								Value:    vmId,
							}},
							GroupBy: []monitor_input.MetricQueryPart{{
								Type:   "tag",
								Params: []string{"vm_id"},
							}},
							Selects: []monitor_input.MetricQuerySelect{
								[]monitor_input.MetricQueryPart{
									{
										Type:   "field",
										Params: []string{vv},
									},
									{
										Type:   "mean",
										Params: []string{},
									},
									{
										Type:   "alias",
										Params: []string{"result"},
									},
								},
							},
						},
					}},
					SkipCheckSeries: true,
				}
				res, err := monitor.UnifiedMonitorManager.PerformQuery(self.cli.s, &params)
				if err != nil {
					m.Lock()
					faild = append(faild, p+"_"+vv)
					m.Unlock()
					// fmt.Println(vv ," error: ",err)
				} else {
					// fmt.Println("************"+p + "_"+ vv+"************",res.String())
					var obj Monitor
					_ = res.Unmarshal(&obj)
					mute.Lock()

					if item == nil && obj.SeriesTotal > 0 {
						item = make([]*MonitorDataItem, len(obj.Series[0].Points))
						for i := range item {
							item[i] = &MonitorDataItem{}
						}
					}
					obj.Type = p + "_" + vv
					ms = append(ms, &obj)
					mute.Unlock()
				}

				wg.Done()
			}(measure, v)
		}
	}
	wg.Wait()

	for _, v := range ms {

		if v.SeriesTotal == 0 {
			continue
		}

		for i := range item {
			item[i].TimeStamp = time.Unix(int64((v.Series[0].Points[i][1]/1000)-8*3600), 0).Format("2006-01-02 15:04:05")
			item[i].InstanceId = v.Series[0].RawName

			val := v.Series[0].Points[i][0]
			if v.Type == "vm_diskio_write_bps" {
				item[i].BPSWrite = val
			} else if v.Type == "vm_diskio_read_bps" {
				item[i].BPSRead = val
			} else if v.Type == "vm_cpu_usage_active" {
				item[i].CPU = val
			} else if v.Type == "vm_mem_used_percent" {
				item[i].Mem = val
			} else if v.Type == "vm_disk_used_percent" {
				item[i].Disk = val
			} else if v.Type == "vm_netio_bps_send" {
				item[i].InternetTX = val
				item[i].IntranetTX = val
			} else if v.Type == "vm_netio_bps_recv" {
				item[i].InternetRX = val
				item[i].IntranetRX = val
			}
		}
	}

	var res []cloudprovider.ICfelMonitorData
	for i := range item {
		res = append(res, item[i])
	}
	return res, faild, nil

}

func (self *SRegion) GetMonitorDataJSON(opts *cloudprovider.MonitorDataJSONOption) (jsonutils.JSONObject, error) {
	params := monitor_input.MetricQueryInput{
		From:     opts.Start,
		To:       opts.End,
		Unit:     false,
		Interval: opts.Interval,
		MetricQuery: []*monitor_input.AlertQuery{{
			Model: monitor_input.MetricQuery{
				Measurement: opts.Measure, // vm_diskio[write_bps,read_bps] vm_cpu[usage_active] vm_mem[used_percent] vm_disk[used_percent] vm_netio[bps_send,bps_recv]
				Tags: []monitor_input.MetricQueryTag{{
					Key:      "vm_id",
					Operator: "=",
					Value:    opts.GuestID,
				}},
				GroupBy: []monitor_input.MetricQueryPart{{
					Type:   "tag",
					Params: []string{"vm_id"},
				}},
				Selects: []monitor_input.MetricQuerySelect{
					[]monitor_input.MetricQueryPart{
						{
							Type:   "field",
							Params: []string{opts.Field},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type:   "alias",
							Params: []string{"result"},
						},
					},
				},
			},
		}},
		SkipCheckSeries: true,
	}
	return monitor.UnifiedMonitorManager.PerformQuery(self.cli.s, &params)
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

func (self *SRegion) CfelInstanceSettingChange(id string, opts *cloudprovider.CfelChangeSettingOption) error {

	_, err := modules.Servers.Update(self.cli.s, id, jsonutils.Marshal(opts))
	return err
}

func (self *SRegion) cfelCreateInstance(hostId, hypervisor string, opts *cloudprovider.CfelSManagedVMCreateConfig) (*SInstance, error) {
	input := compute.ServerCreateInput{
		ServerConfigs: &compute.ServerConfigs{},
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

	ins := &SInstance{}
	return ins, self.create(&modules.Servers, input, ins)
}

func (self *SInstance) GetIsolatedDevice() ([]*cloudprovider.IsolatedDeviceInfo, error) {

	var res []*cloudprovider.IsolatedDeviceInfo
	for _, v := range self.IsolatedDevices {
		res = append(res, &cloudprovider.IsolatedDeviceInfo{
			DevType:        v.DevType,
			Model:          v.Model,
			VendorDeviceId: v.VendorDeviceId,
		})
	}
	return res, nil
}
