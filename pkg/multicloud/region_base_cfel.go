package multicloud

import (
	"context"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

func (self *SRegion) SetSkuExtInfo(string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "SetSkuExtInfo")
}
func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetICfelSkus")
}
func (self *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetSpotPostPaidPrice")
}
func (self *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetPostPaidPrice")
}
func (self *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetPrePaidPrice")
}
func (self *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	return "", errors.Wrapf(cloudprovider.ErrNotImplemented, "GetSpotPostPaidStatus")
}
func (self *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	return "", errors.Wrapf(cloudprovider.ErrNotImplemented, "GetPostPaidStatus")
}
func (self *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	return "", errors.Wrapf(cloudprovider.ErrNotImplemented, "GetPrePaidStatus")
}

func (self *SRegion) CreateBareMetal(opts *cloudprovider.CfelSManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CreateBareMetal")

}

func (self *SRegion) CreateVM(opts *cloudprovider.CfelSManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CreateVM")
}

func (self *SRegion) SchedulerForecast(hypervisor string, opts *cloudprovider.CfelSManagedVMCreateConfig) (jsonutils.JSONObject, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CreateBareMetal")
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

func (self *SRegion) CfelGetDiskById(opts *cloudprovider.FilterOption) (cloudprovider.ICloudDisk, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelGetDiskById")
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

func (self *SRegion) GetICfelCloudImageById(id string) (cloudprovider.ICloudImage, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetICfelCloudImageById")
}

func (self *SRegion) GetInstanceMatchImage(instanceType string) ([]cloudprovider.ICloudImage, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetInstanceMatchImage")
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

func (self *SRegion) GetSshKeypair(project string, isAdmin bool) (string, error) {
	return "", errors.Wrapf(cloudprovider.ErrNotImplemented, "GetSshKeypair")
}

func (self *SRegion) CfelUpdateNetworkTags(id string, tags map[string]string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelUpdateNetworkTags")
}

func (self *SRegion) ICfelSetImageCanDelete(id string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "ICfelSetImageCanDelete")
}

func (self *SRegion) GetICfelSkuPrice(opt *cloudprovider.CfelSkuPriceOptions) (map[string]string, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetICfelSkuPrice")
}

func (self *SRegion) ExecHttp(ctx context.Context, opts *cloudprovider.CfelExecHttpInput) (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotImplemented
}
