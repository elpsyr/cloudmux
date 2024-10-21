package ecloudcfel

import (
	"context"
	"slices"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
)

type DiskType struct {
	CinderType        string   `json:"cinderType"`
	AttachServerTypes []string `json:"attachServerTypes"`
	OpType            string   `json:"opType"`
	Region            string   `json:"region"`
	OnlineStatus      string   `json:"onlineStatus"`
	OnlineMode        string   `json:"onlineMode"`
}

func (r *SZone) GetCapability() (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (r *SZone) GetICfelDiskType() (map[string]interface{}, error) {
	req := NewConsoleRequest(r.region.ID, "/api/v2/volume/customer/volumeType/list", nil, nil)
	var res []DiskType
	err := r.region.client.doList(context.Background(), req, &res)
	if err != nil {
		return nil, nil
	}
	var ret = make(map[string]interface{})
	for _, val := range res {
		if val.Region == r.Region && slices.Contains(val.AttachServerTypes, "VM") {
			if v, ok := ret[val.OpType]; ok {
				vv, _ := v.([]string)
				vv = append(vv, val.CinderType)
				ret[val.OpType] = vv
			} else {
				ret[val.OpType] = []string{val.CinderType}
			}
		}
	}
	return ret, nil
}
