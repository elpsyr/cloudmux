package cloudpods

import "yunion.io/x/cloudmux/pkg/cloudprovider"

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
