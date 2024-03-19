package cloudpods

import "yunion.io/x/cloudmux/pkg/cloudprovider"

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
