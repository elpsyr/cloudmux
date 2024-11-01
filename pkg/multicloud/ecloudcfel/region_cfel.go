package ecloudcfel

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
)

var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

func (self *SRegion) SetSkuExtInfo(info string) error {
	return json.Unmarshal([]byte(info), &self.gpuSkuInfo)
}

func (r *SRegion) GetInstanceMatchImage(instancetype string) ([]cloudprovider.ICloudImage, error) {
	if strings.Count(instancetype, ".") == 1 { // 移动云的bug
		instancetype = instancetype + ".8"
	}
	query := map[string]string{
		"specsName": instancetype,
		"pageSize":  "100",
		"page":      "1",
	}
	request := NewNovaRequest(NewApiRequest(r.ID, "/api/openapi-ims/user/v5/image/public", query, nil))
	images := make([]SImage, 0)
	err := r.client.doList(context.Background(), request, &images)
	if err != nil {
		return nil, err
	}
	var img []cloudprovider.ICloudImage
	for i := range images {
		img = append(img, &images[i])
	}
	return img, nil
}

func (r *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	request := NewNovaRequest(NewApiRequest(r.ID, "/api/openapi-ims/user/v5/image/public", nil, nil))
	images := make([]SImage, 0, 5)
	err := r.client.doList(context.Background(), request, &images)
	if err != nil {
		return nil, err
	}
	var img []cloudprovider.ICloudImage
	for _, val := range images {
		img = append(img, &val)
	}
	return img, nil
}

func (r *SRegion) GetICfelCloudImageById(id string) (cloudprovider.ICloudImage, error) {
	return nil, nil
}

type price struct {
	ServerPrice     string `json:"serverPrice"`
	BootVolumePrice string `json:"bootVolumePrice"`
}

func (self *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, nil
}

func (self *SRegion) getPrice(instanceType, feeUnit string) (float64, error) {
	param := map[string]interface{}{
		"productType": "vm",
		"specsName":   instanceType,
		"feeUnit":     feeUnit,
	}
	req := NewConsoleRequest(self.ID, "/api/openapi-ecs/acl/v3/server/query/price", nil, jsonutils.Marshal(param))
	res, err := self.client.doPost(req)
	if err != nil {
		return 0, err
	}
	var ret price
	if err := res.Unmarshal(&ret); err != nil {
		return 0, err
	}
	price, err := strconv.ParseFloat(strings.ReplaceAll(ret.ServerPrice, ",", ""), 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (self *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return self.getPrice(instanceType, "hour")
}

func (self *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	return self.getPrice(instanceType, "month")
}

func (self *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	return "soldout", nil
}
func (self *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	return "available", nil
}
func (self *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	return "available", nil
}

func (self *SRegion) GetICfelSkuPrice(opt *cloudprovider.CfelSkuPriceOptions) (map[string]string, error) {
	param := map[string]interface{}{
		"productType": "vm",
		"specsName":   opt.InstanceType,
		"feeUnit":     opt.FeeUnit,
		"duration":    opt.Duration,
		"quantity":    opt.Quantity,
		"bootVolume": map[string]interface{}{
			"size":       opt.SysDiskSize,
			"volumeType": opt.SysDiskType,
		},
	}
	req := NewConsoleRequest(self.ID, "/api/openapi-ecs/acl/v3/server/query/price", nil, jsonutils.Marshal(param))
	res, err := self.client.doPost(req)
	if err != nil {
		return nil, err
	}
	var ret = make(map[string]string)
	if err := res.Unmarshal(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}
