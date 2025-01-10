package cucloudcfel

import (
	"encoding/json"
	"strconv"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
)

func (s *SZone) GetCapability() (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotImplemented
}

type storageInfo struct {
	MaxValue        string `json:"maxValue"`
	MeasurementUnit string `json:"measurementUnit"`
	MinValue        string `json:"minValue"`
	PrtyCode        string `json:"prtyCode"`
	PrtyName        string `json:"prtyName"`
	PrtyType        string `json:"prtyType"`
	StepLen         string `json:"stepLen"`
}

func (s *SZone) GetICfelDiskType(diskType string) ([]*cloudprovider.DiskInfo, error) {
	res, err := s.fetchDiskType()
	if err != nil {
		return nil, err
	}
	var ret = []*cloudprovider.DiskInfo{}
	for _,val := range res {
		var info []storageInfo
		min,max := 10,30720
		if err = json.Unmarshal([]byte(val.ProductDetail),&info);err == nil {
			for _,vv := range info {
				if vv.PrtyCode == "storageSize" {
					min,_ = strconv.Atoi(vv.MinValue)
					max,_ = strconv.Atoi(vv.MaxValue)
				}
			}
		}
		ret = append(ret, &cloudprovider.DiskInfo{
			Name:        val.ProductClassName,
			StorageType: val.StorageType,
			MinSizeGB:   min,
			MaxSizeGB:   max,
			StepLen:     0,
			IsSysDisk:   true,
			IsDataDisk:  true,
		})
	}
	return ret, nil
}
