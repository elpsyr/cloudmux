package cloudpods

import (
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

func (s *SZone) GetCapability() (jsonutils.JSONObject, error) {
	return modules.Zones.GetSpecific(s.region.cli.s,s.Id,"capability",nil)
}