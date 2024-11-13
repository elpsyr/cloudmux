package ecloudcfel

import (
	"context"
	"slices"
	"sync"

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

type DiskStatus struct {
	ProductType string `json:"productType,omitempty"`
	Status      string `json:"status,omitempty"`
}

type SysDiskType struct {
	BootVolumeType     string `json:"bootVolumeType"`
	BootVolumeTypeName string `json:"bootVolumeTypeName"`
	SoldOut            string `json:"soldOut"`
	ZoneDesc           string `json:"zoneDesc"`
	ZoneName           string `json:"zoneName"`
}

func (r *SZone) GetCapability() (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (r *SZone) GetICfelDiskType(diskType string) (map[string]interface{}, error) {
	if diskType == "sys" {
		// https://ecloud.10086.cn/op-help-center/doc/article/75582
		query := map[string]string{
			"region": r.Region,
		}
		req := NewConsoleRequest(r.region.ID, "/api/openapi-ecs/acl/v3/server/system/disk/type", query, nil)
		var res []SysDiskType
		err := r.region.client.doList(context.Background(), req, &res)
		if err != nil {
			return nil, nil
		}
		var ret = make(map[string]interface{})
		for _, val := range res {
			if val.SoldOut == "0" {
				ret[val.BootVolumeType] = val.BootVolumeTypeName
			}
		}
		return ret, nil
	}

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
	var wg sync.WaitGroup

	var result = make(map[string]interface{})
	var lock = sync.Mutex{}

	for dt := range ret {
		wg.Add(1)
		go func(dt string) {
			var query = map[string]string{
				"productType": dt,
				"poolId":      r.PoolId,
			}
			req := NewConsoleRequest(r.region.ID, "/api/ebs/acl/v3/mop/common/getPoolInfo", query, nil)
			req.SetMethod("GET")
			var res []DiskStatus
			jsonRes, err := r.region.client.request(context.Background(), req)
			if err == nil {
				if err := jsonRes.Unmarshal(&res, "poolList"); err == nil {
					if res != nil {
						lock.Lock()
						for _, val := range res {
							result[val.ProductType] = val.Status
						}
						lock.Unlock()
						// ch <- res
					}
				}
			}
			defer wg.Done()
		}(dt)
	}
	wg.Wait()

	return result, nil
}
