package ecloudcfel

import (
	"context"
	"encoding/json"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

func (self *SRegion) SetSkuExtInfo(info string) error {
	return json.Unmarshal([]byte(info), &self.gpuSkuInfo)
}

func (r *SRegion) GetInstanceMatchImage(instancetype string) ([]cloudprovider.ICloudImage, error) {
	if !strings.HasSuffix(instancetype, ".8") {// 移动云的bug
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
