package cloudpods

import (
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

func (self *SRegion) CfelUpdateNetworkTags(id string,tags map[string]string) error {
	_,err := modules.Networks.PerformAction(self.cli.s,id,"set-user-metadata",jsonutils.Marshal(tags))
	
	return err
}