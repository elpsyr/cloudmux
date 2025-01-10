package baiducfel

import (
	"fmt"
	"regexp"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
)

type SZone struct {
	multicloud.SResourceBase
	region   *SRegion
	ZoneName string

	host cloudprovider.ICloudHost

	istorages []cloudprovider.ICloudStorage
}

func (S SZone) GetId() string {
	return S.ZoneName
}

func (S SZone) GetName() string {

	pattern := S.region.Region + `-(.+)`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(S.ZoneName)

	if len(match) > 1 {
		result := match[1]
		return fmt.Sprintf("可用区%s", strings.ToUpper(result))
	} else {
		return S.ZoneName
	}
	return S.ZoneName
}

func (s SZone) GetGlobalId() string {
	return fmt.Sprintf("%s/%s", s.region.GetGlobalId(), s.ZoneName)
}

func (S SZone) GetStatus() string {
	return ""
}

func (S SZone) GetSysTags() map[string]string {
	return nil
}

func (S SZone) GetTags() (map[string]string, error) {
	return nil, nil
}

func (S SZone) SetTags(tags map[string]string, replace bool) error {
	return nil
}

func (S SZone) GetI18n() cloudprovider.SModelI18nTable {
	return nil
}

func (S SZone) GetIRegion() cloudprovider.ICloudRegion {
	return S.region
}

func (S SZone) GetIHosts() ([]cloudprovider.ICloudHost, error) {
	return nil, nil
}

func (S SZone) GetIHostById(id string) (cloudprovider.ICloudHost, error) {
	if S.host == nil {
		S.host = &SHost{zone: &S}
	}
	return S.host, nil
}

func (z SZone) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	if z.istorages == nil {
		err := z.fetchStorages()
		if err != nil {
			return nil, err
		}
	}
	return z.istorages, nil
}

func (S SZone) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	return nil, nil
}

func (r *SZone) GetCapability() (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotImplemented
}

type diskType struct {
	DiskInfos []diskInfo `json:"diskInfos"`
	ZoneName  string     `json:"zoneName"`
}

type diskInfo struct {
	MaxDiskSize int    `json:"maxDiskSize"`
	MinDiskSize int    `json:"minDiskSize"`
	StorageType string `json:"storageType"`
}

func (r *SZone) GetICfelDiskType(dt string) ([]*cloudprovider.DiskInfo, error) {
	var query = map[string]string{
		"zoneName": r.ZoneName,
	}
	var ret []diskType
	res, err := r.region.doGetWithoutVal(ServiceDisk, "/v2/volume/disk", query)
	if err != nil {
		return nil, err
	}
	if err = res.Unmarshal(&ret, "diskZoneResources"); err != nil {
		return nil, err
	}
	var result = []*cloudprovider.DiskInfo{}
	if len(ret) == 0 {
		return result, nil
	}

	for _, val := range ret[0].DiskInfos {
		if dt == "sys" && val.StorageType == "hdd" {
			continue
		}
		min, max := val.MinDiskSize, val.MaxDiskSize
		if dt == "sys" {
			min, max = 40, 500
		}
		if val.StorageType == "enhanced_ssd_pl2" {
			min = 461
		}
		if val.StorageType == "enhanced_ssd_pl1" {
			min = 20
		}
		if val.StorageType == "ssd" {
			val.StorageType = "hp1"
		}
		// ssd 50G 起售 https://cloud.baidu.com/doc/BCC/s/Ujwvyo1ta
		if dt == "data" && val.StorageType == "hp1" {
			min = 50
		}
		result = append(result, &cloudprovider.DiskInfo{
			Name:        "",
			StorageType: val.StorageType,
			MinSizeGB:   min,
			MaxSizeGB:   max,
			StepLen:     0,
		})
	}
	return result, nil
}

func (z *SZone) fetchStorages() error {
	istorages := make([]cloudprovider.ICloudStorage, len(storageTypes))
	for i := range istorages {
		istorages[i] = &SStorage{
			zone:        z,
			storageType: storageTypes[i],
		}
	}
	z.istorages = istorages
	return nil
}
