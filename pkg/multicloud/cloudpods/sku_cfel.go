package cloudpods

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	api "yunion.io/x/onecloud/pkg/apis/compute"
)

// Verify that *SRegion implements ICfelCloudRegion
var _ cloudprovider.ICfelCloudSku = (*SServerSku)(nil)

type SInstanceTypeCFEL struct {
	ZoneID string // zone
}

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	skus, err := self.GetServerSkus()
	if err != nil {
		return nil, err
	}
	var ret []cloudprovider.ICfelCloudSku
	for i := range skus {
		skus[i].region = self
		ret = append(ret, &skus[i])
	}
	return ret, nil
}

func (self *SServerSku) GetZoneID() string {
	return self.ZoneID
}

func (self *SServerSku) GetGPUMemorySizeMB() int {
	return 0
}

func (self *SServerSku) GetIsBareMetal() bool {
	return false
}

func (self *SRegion) CreateBareMetal(opts *cloudprovider.SManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	hypervisor := api.HYPERVISOR_BAREMETAL
	ins, err := self.CreateInstance("", hypervisor, opts)
	if err != nil {
		return nil, err
	}

	hostInstance, err := self.GetHostInstance(ins.GetId())
	return hostInstance, nil
}
