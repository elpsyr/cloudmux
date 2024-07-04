package multicloud

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

func (self *SRegion) CreateBareMetal(opts *cloudprovider.CfelSManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CreateBareMetal")

}

func (self *SRegion) CreateVM(opts *cloudprovider.CfelSManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CreateVM")
}

func (self *SRegion) DeleteVM(instanceId string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "DeleteVM")

}

func (self *SRegion) CreateImageByUrl(params *cloudprovider.CfelSImageCreateOption) (cloudprovider.ICloudImage, error) {

	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CreateImageByUrl")
}

func (self *SRegion) GetImageByID(id string) (cloudprovider.ICloudImage, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetImage")
}

func (self *SRegion) ResetGuestPassword(params *cloudprovider.CfelResetGuestPasswordOption) (cloudprovider.ICloudVM, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "ResetGuestPassword")
}

func (self *SRegion) PingQga(guestId string, timeout int) (bool, error) {
	return false, errors.Wrapf(cloudprovider.ErrNotImplemented, "ResetGuestPassword")
}

func (self *SRegion) CfelCreateDisk(params *cloudprovider.CfelDiskCreateConfig) (cloudprovider.ICloudDisk, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelCreateDisk")
}

func (self *SRegion) CfelAttachDisk(instanceId, diskId string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelCreateDisk")
}

func (self *SRegion) CfelDetachDisk(instanceId, diskId string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelCreateDisk")
}

func (self *SRegion) CfelInstanceSettingChange(id string, params *cloudprovider.CfelChangeSettingOption) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelInstanceSettingChange")
}

func (self *SRegion) CfelGetINetworks(*cloudprovider.GetNetworkOptions) ([]cloudprovider.ICloudNetwork, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetINetworks")
}

func (self *SRegion) GetIHostsByCondition(*cloudprovider.FilterOption) ([]cloudprovider.ICloudHost, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetIHostsByCondition")
}

func (self *SRegion) MigrateForecast(*cloudprovider.MigrateForecastOption) ([]cloudprovider.ICfelFilter, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "MigrateForecast")
}

// func (self *SRegion) GetMonitorData(vmId, start, end, interval string) ([]cloudprovider.ICfelMonitorData, []string, error) {
// 	return nil, []string{}, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetMonitorData")
// }

func (self *SRegion) GetMonitorDataJSON(*cloudprovider.MonitorDataJSONOption) (jsonutils.JSONObject, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetMonitorDataJSON")
}

func (self *SRegion) GetGeneralUsage() (cloudprovider.ICfelGeneralUsage, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetGeneralUsage")
}

func (self *SRegion) ICfelDeleteImage(id string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "ICfelDeleteImage")
}
func (self *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetICfelCloudImage")
}

func (self *SRegion) SetImageUserTag(*cloudprovider.CfelSetImageUserTag) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "SetImageUserTag")
}

func (self *SRegion) GetUsableIEip() ([]cloudprovider.ICloudEIP, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetUsableIEip")
}

func (self *SRegion) GetLoadbalancerSkus() ([]cloudprovider.ICfelLoadbalancerSku, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetLoadbalancerSkus")
}

func (self *SRegion) CfelCreateILoadBalancerCertificate(cert *cloudprovider.SCfelLoadbalancerCertificate) (cloudprovider.ICloudLoadbalancerCertificate, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelCreateILoadBalancerCertificate")
}