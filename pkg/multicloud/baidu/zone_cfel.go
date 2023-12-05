package baidu

import (
	"fmt"
	"regexp"
	"strings"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SZone struct {
	multicloud.SResourceBase
	region   *SRegion
	ZoneName string
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

func (S SZone) GetGlobalId() string {
	return S.ZoneName
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
	return nil, nil
}

func (S SZone) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	return nil, nil
}

func (S SZone) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	return nil, nil
}
