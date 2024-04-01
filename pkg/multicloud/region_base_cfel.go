package multicloud

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

func (self *SRegion) CreateBareMetal(opts *cloudprovider.SManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CreateBareMetal")

}

func (self *SRegion) CreateVM(opts *cloudprovider.SManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CreateVM")

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

func (self *SRegion) CfelCreateDisk(params *cloudprovider.CfelDiskCreateConfig)(cloudprovider.ICloudDisk,error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelCreateDisk")
}

func (self *SRegion) CfelAttachDisk(instanceId,diskId string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelCreateDisk")
}

func (self *SRegion) CfelDetachDisk(instanceId,diskId string) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelCreateDisk")
}

func (self *SRegion) CfelInstanceSettingChange(id string,params *cloudprovider.CfelChangeSettingOption) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "CfelInstanceSettingChange")
}

func (self *SRegion) CfelGetINetworks() ([]cloudprovider.ICloudNetwork, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetINetworks")
}