package multicloud

import (
	"context"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (self *SInstanceBase) CfelRebuildRoot(ctx context.Context, opts *cloudprovider.CfelSManagedVMRebuildRootConfig) (string, error) {
	return "", cloudprovider.ErrNotImplemented
}

func (self *SInstanceBase) RebootVM(ctx context.Context) error {
	return cloudprovider.ErrNotImplemented
}
func (self *SInstanceBase) GetMonitorData(start, end string) ([]cloudprovider.ICfelMonitorData, error) { // 获取主机监控数据
	return nil, cloudprovider.ErrNotImplemented
}
func (self *SInstanceBase) GetIsolatedDevice() ([]*cloudprovider.IsolatedDeviceInfo, error) {
	return nil, cloudprovider.ErrNotImplemented
}
func (self *SInstanceBase) GetCfelHypervisor() string {
	return ""
}
func (self *SInstanceBase) GetSSHInfo() (*cloudprovider.ServerSSHInfo, error) {
	return nil, cloudprovider.ErrNotImplemented
}
