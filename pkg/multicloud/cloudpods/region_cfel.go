package cloudpods

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
	image "yunion.io/x/onecloud/pkg/mcclient/modules/image"
)

// Verify that *SRegion implements ICfelCloudRegion
// 私有云 无需实现价格接口
var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

func (self *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	return "", cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	return "", cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	return "", cloudprovider.ErrNotImplemented
}

// GetIVMs
// region 返回的 vm 无 host *SHost 操作对象
func (self *SRegion) GetIVMs() ([]cloudprovider.ICloudVM, error) {
	instances, err := self.GetInstances("")
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICloudVM
	for i := range instances {
		ret = append(ret, &instances[i])
	}
	return ret, nil
}

func (self *SRegion) CreateImageByUrl(opts *cloudprovider.CfelSImageCreateOption) (cloudprovider.ICloudImage, error) {
	params := map[string]interface{}{
		"generate_name": opts.ImageName,
		"protected":     opts.IsProtected,
		"copy_from":     opts.CopyFrom,
		"properties": map[string]string{
			"os_type":         opts.OsType,
			"os_distribution": opts.OsDistribution,
			"os_arch":         opts.OsArch,
			"os_version":      opts.OsVersion,
		},
	}
	res, err := image.Images.Create(self.cli.s, jsonutils.Marshal(params))
	if err != nil {
		return nil, err
	}
	var image SImage
	if err := res.Unmarshal(&image); err != nil {
		return nil, err
	}
	return &image, nil
}

func (self *SRegion) CfelCreateDisk(opts *cloudprovider.CfelDiskCreateConfig) (cloudprovider.ICloudDisk, error) {
	params := map[string]interface{}{
		"name":        opts.Name,
		"size":        opts.SizeGb,
		"backend":     opts.Backend,
		"medium":      opts.Medium,
		"description": opts.Desc,
	}
	res, err := modules.Disks.Create(self.cli.s, jsonutils.Marshal(params))
	if err != nil {
		return nil, err
	}
	var dd SDisk
	if err := res.Unmarshal(&dd); err != nil {
		return nil, err
	}
	return &dd, nil
}

func (self *SRegion) CfelGetINetworks() ([]cloudprovider.ICloudNetwork, error) {
	networks, err := self.GetNetworks("")
	if err != nil {
		return nil, err
	}
	ret := []cloudprovider.ICloudNetwork{}
	for i := range networks {
		// networks[i].wire = self
		ret = append(ret, &networks[i])
	}
	return ret, nil
}

func (self *SRegion) GetIHostsByCondition(opts *cloudprovider.FilterOption) ([]cloudprovider.ICloudHost, error) {
	
	params := map[string]interface{}{
		// "scope":                 "system",
		"show_fail_reason": "true",
		"host_type":        opts.HostType,
		"limit":            opts.Limit,
		"enabled":               1,
		"host_status": "online",
		"os_arch":     opts.OsArch,
		//"field":       ,
		// "server_id_for_network": "f13faa78-5a46-4236-80ee-f427defd947e",
		// "project_domain":        "default",
		//"filter":                "id.notin(7d09d25e-87ef-44db-8bf5-bf42b8554388,7096846e-4341-4267-874e-d047838e2c99)",
		"details": false,
	}
	if len(opts.FilterIds) > 0 {
		params["filter"] = "id.notin(" + opts.FilterIds + ")"
	}
	if len(opts.Field) > 0 {
		params["field"] = opts.Field
	}
	var ret []SHost

	err := self.cli.list(&modules.Hosts, params, &ret)
	var res []cloudprovider.ICloudHost
	for i := range ret {
		res = append(res, &ret[i])
	}
	return res, err
}

func (self *SRegion) MigrateForecast(opts *cloudprovider.MigrateForecastOption) ([]cloudprovider.ICfelFilter, error) {
	params := map[string]interface{}{
		"live_migrate":opts.LiveMigrate,
		"skip_cpu_check":false,
		"skip_kernel_check":opts.SkipKernelCheck,
		// "is_rescue_mode":true,
	}
	res,err := self.cli.perform(&modules.Servers,opts.GuestId,"migrate-forecast",params)
	if err != nil {
		return nil,err
	}
	rr,err := res.Get("filtered_candidates")
	if err != nil {
		return nil,err
	}
	var filter []SCfelFilter
	err = rr.Unmarshal(&filter)
	if err != nil {
		return nil,err
	}
	var ret []cloudprovider.ICfelFilter
	for i := range filter {
		ret = append(ret, &filter[i])
	}
	return ret, nil
}
