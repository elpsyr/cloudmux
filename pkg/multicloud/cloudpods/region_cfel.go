package cloudpods

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	disk "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/image"
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
	res, err := modules.Images.Create(self.cli.s, jsonutils.Marshal(params))
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
	res, err := disk.Disks.Create(self.cli.s, jsonutils.Marshal(params))
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
	networks,err := self.GetNetworks("")
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