package cloudpods

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
	image "yunion.io/x/onecloud/pkg/mcclient/modules/image"
	"yunion.io/x/pkg/errors"
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

// GetBareMetalIHosts 获取可用区下的裸金属 host
func (region *SRegion) GetBareMetalIHosts(zoneId string) ([]cloudprovider.ICloudHost, error) {
	hosts, err := region.getBareMetalHosts(zoneId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetBareMetalHosts")
	}
	zone, err := region.GetIZoneById(zoneId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetIZoneById")
	}
	var ret []cloudprovider.ICloudHost
	for i := range hosts {
		hosts[i].zone = zone.(*SZone)
		ret = append(ret, &hosts[i])
	}
	return ret, nil
}

func (region *SRegion) getBareMetalHosts(zoneId string) ([]SHost, error) {
	params := map[string]interface{}{
		"baremetal": true,
	}
	if len(zoneId) > 0 {
		params["zone_id"] = zoneId
	}
	ret := []SHost{}
	err := region.list(&modules.Hosts, params, &ret)
	if err != nil {
		return nil, errors.Wrap(err, "list")
	}
	return ret, nil
}

func (self *SRegion) CreateImageByUrl(opts *cloudprovider.CfelSImageCreateOption) (cloudprovider.ICloudImage, error) {
	params := map[string]interface{}{
		"name":          opts.ImageName,
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

func (self *SRegion) CfelGetINetworks(opts *cloudprovider.GetNetworkOptions) ([]cloudprovider.ICloudNetwork, error) {

	params := map[string]interface{}{}
	if opts != nil {
		if len(opts.WireId) > 0 {
			params["wire_id"] = opts.WireId
		}
		if len(opts.ServerType) > 0 {
			params["server_type"] = opts.ServerType
		}
		if len(opts.ZoneId) > 0 {
			params["zone_id"] = opts.ZoneId
		}
		if opts.WithUserMeta {
			params["with_user_meta"] = opts.WithUserMeta
		}
		if len(opts.VpcId) > 0 {
			params["vpc_id"] = opts.VpcId
		}
		if len(opts.Ids) > 0 {
			params["id"] = opts.Ids
		}
	}
	networks := []SNetwork{}
	err := self.list(&modules.Networks, params, &networks)
	if err != nil {
		return nil, err
	}

	ret := []cloudprovider.ICloudNetwork{}
	for i := range networks {
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
		"enabled":          1,
		"host_status":      "online",
		"os_arch":          opts.OsArch,
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
		"live_migrate":      opts.LiveMigrate,
		"skip_cpu_check":    false,
		"skip_kernel_check": opts.SkipKernelCheck,
		// "is_rescue_mode":true,
	}
	res, err := self.cli.perform(&modules.Servers, opts.GuestId, "migrate-forecast", params)
	if err != nil {
		return nil, err
	}
	rr, err := res.Get("filtered_candidates")
	if err != nil {
		return nil, err
	}
	var filter []SCfelFilter
	err = rr.Unmarshal(&filter)
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICfelFilter
	for i := range filter {
		ret = append(ret, &filter[i])
	}
	return ret, nil
}

func (self *SRegion) GetGeneralUsage() (cloudprovider.ICfelGeneralUsage, error) {
	var usage GeneralUsage
	res, err := modules.Usages.GetGeneralUsage(self.cli.s, nil)
	if err != nil {
		return nil, err
	}
	_ = res.Unmarshal(&usage)
	return &usage, nil
}

func (self *SRegion) ICfelDeleteImage(id string) error {
	_, err := image.Images.Delete(self.cli.s, id, nil)
	return err
}

func (self *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	var params = map[string]interface{}{
		"is_public":      true,
		"is_guest_image": false,
		"with_user_meta": withUserMeta,
	}
	var rr []SImage
	err := self.list(&image.Images, params, &rr)

	if err != nil {
		return nil, nil
	}
	var ret []cloudprovider.ICloudImage
	for i := range rr {
		ret = append(ret, &rr[i])
	}
	return ret, nil

}

// GetICfelCloudImageById 获取 cloudpods 镜像
func (self *SRegion) GetICfelCloudImageById(id string) (cloudprovider.ICloudImage, error) {
	return self.GetImage(id)
}

func (self *SRegion) SetImageUserTag(opts *cloudprovider.CfelSetImageUserTag) error {
	_, err := image.Images.PerformAction(self.cli.s, opts.ImageId, "set-user-metadata", jsonutils.Marshal(opts.Tags))
	return err
}

func (self *SRegion) GetUsableIEip() ([]cloudprovider.ICloudEIP, error) {
	eips := []SEip{}
	params := map[string]interface{}{
		"usable": true,
	}

	err := self.list(&modules.Elasticips, params, &eips)
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICloudEIP
	for i := range eips {
		ret = append(ret, &eips[i])
	}
	return ret, nil
}

func (self *SRegion) GetSshKeypair(project string, isAdmin bool) (string, error) {
	query := jsonutils.NewDict()
	if isAdmin {
		query.Add(jsonutils.JSONTrue, "admin")
	}
	var keys jsonutils.JSONObject
	if len(project) == 0 {
		listResult, err := modules.Sshkeypairs.List(self.cli.s, query)
		if err != nil {
			return "", err
		}
		keys = listResult.Data[0]
	} else {
		result, err := modules.Sshkeypairs.GetById(self.cli.s, project, query)
		if err != nil {
			return "", err
		}
		keys = result
	}
	privKey, _ := keys.GetString("private_key")
	return privKey, nil
}
