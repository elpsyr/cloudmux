package ecloudcfel

import (
	"context"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

func (r *SRegion) GetInstanceMatchImage(instancetype string) ([]cloudprovider.ICloudImage, error) {
	query := map[string]string{"specsName": instancetype}
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

func (r *SRegion) GetICfelDiskType() (string, error) {
	req := NewConsoleRequest(r.ID, "/api/v2/volume/customer/volumeType/list", nil, nil)
	var res map[string]interface{}
	_ = r.client.doGet(context.Background(), req, &res)
	return "", nil
}
