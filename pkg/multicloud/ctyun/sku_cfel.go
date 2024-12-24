package ctyun

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

var _ cloudprovider.ICfelCloudSku = (*ServerSku)(nil)

// https://www.ctyun.cn/document/10029787/10348867#section-0b0dc719ac08c0ae
// Delete implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) Delete() error {
	panic("unimplemented")
}

// GetAttachedDiskCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetAttachedDiskCount() int {
	return 0
}

// GetAttachedDiskSizeGB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetAttachedDiskSizeGB() int {
	return 0
}

// GetAttachedDiskType implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetAttachedDiskType() string {
	return ""
}

// GetCpuArch implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetCpuArch() string {
	return s.CPUInfo
}

// GetCpuCoreCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetCpuCoreCount() int {
	return s.FlavorCPU
}

// GetCreatedAt implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetCreatedAt() time.Time {
	return time.Now()
}

// GetDataDiskMaxCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetDataDiskMaxCount() int {
	return 0
}

// GetDataDiskTypes implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetDataDiskTypes() string {
	return ""
}

// GetDescription implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetDescription() string {
	return s.FlavorName
}

// GetGPUMemorySizeMB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGPUMemorySizeMB() int {
	return s.VideoMemSize * 1024
}

// GetGlobalId implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGlobalId() string {
	return s.FlavorId
}

// GetGpuAttachable implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGpuAttachable() bool {
	return s.GpuCount > 0
}

// GetGpuCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGpuCount() string {
	return strconv.Itoa(s.GpuCount)
}

// GetGpuMaxCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGpuMaxCount() int {
	return s.GpuCount
}

// GetGpuSpec implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetGpuSpec() string {
	arr := strings.Split(s.FlavorName, ".")
	if len(arr) > 0 {
		if spec, ok := s.skuExtInfo[strings.ToUpper(arr[0])]; ok {
			return spec
		}
	}
	return ""
}

// GetId implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetId() string {
	return s.FlavorName
}

// GetInstanceTypeCategory implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetInstanceTypeCategory() string {
	return s.FlavorType
}

// GetInstanceTypeFamily implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetInstanceTypeFamily() string {
	return s.FlavorSeries
}

// GetIsBareMetal implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetIsBareMetal() bool {
	return false
}

// GetMemorySizeMB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetMemorySizeMB() int {
	return s.FlavorRAM * 1024
}

// GetName implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetName() string {
	return s.FlavorName
}

// GetNicMaxCount implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetNicMaxCount() int {
	return s.NicCount
}

// GetNicType implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetNicType() string {
	return ""
}

// GetOsName implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetOsName() string {
	return ""
}

// GetPostpaidStatus implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetPostpaidStatus() string {
	if s.Available {
		return api.SkuStatusAvailable
	}
	return api.SkuStatusSoldout
}

// GetPrepaidStatus implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetPrepaidStatus() string {
	if s.Available {
		return api.SkuStatusAvailable
	}
	return api.SkuStatusSoldout
}

// GetSpotpaidStatus implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSpotpaidStatus() string {
	return api.SkuStatusSoldout
}

// GetStatus implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetStatus() string {
	return "ready"
}

// GetSysDiskMaxSizeGB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSysDiskMaxSizeGB() int {
	return 0
}

// GetSysDiskMinSizeGB implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSysDiskMinSizeGB() int {
	return 0
}

// GetSysDiskResizable implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSysDiskResizable() bool {
	return false
}

// GetSysDiskType implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetSysDiskType() string {
	return ""
}

// GetZoneID implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) GetZoneID() string {
	return s.ZoneId
}

// IsEmulated implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) IsEmulated() bool {
	return false
}

// Refresh implements cloudprovider.ICfelCloudSku.
func (s *ServerSku) Refresh() error {
	panic("unimplemented")
}

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	// https://www.ctyun.cn/document/10026730/10597667
	params := map[string]interface{}{
		"regionID": self.RegionId,
	}

	resp, err := self.post(SERVICE_ECS, "/v4/ecs/flavor/list", params)
	if err != nil {
		return nil, err
	}
	ret := struct {
		ReturnObj struct {
			FlavorList []ServerSku
		}
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	var result []cloudprovider.ICfelCloudSku
	for i := range ret.ReturnObj.FlavorList {
		for _, zoneName := range ret.ReturnObj.FlavorList[i].AzList {
			sku := ret.ReturnObj.FlavorList[i]
			sku.skuExtInfo = self.gpuSkuInfo
			sku.ZoneId = fmt.Sprintf("%s/%s", self.GetGlobalId(), zoneName)
			result = append(result, &sku)
		}
	}
	return result, nil
}

type priceResp struct {
	ReturnObj struct {
		DiscountPrice  float64 `json:"discountPrice"`
		FinalPrice     float64 `json:"finalPrice"`
		IsSucceed      bool    `json:"isSucceed"`
		SubOrderPrices []struct {
			FinalPrice      float64 `json:"finalPrice"`
			OrderItemPrices []struct {
				FinalPrice   float64 `json:"finalPrice"`
				ItemID       string  `json:"itemId"`
				ResourceType string  `json:"resourceType"`
				TotalPrice   int64   `json:"totalPrice"`
			} `json:"orderItemPrices"`
			ServiceTag string `json:"serviceTag"`
			TotalPrice int64  `json:"totalPrice"`
		} `json:"subOrderPrices"`
		TotalPrice int64 `json:"totalPrice"`
	} `json:"returnObj"`
}

func (self *SRegion) GetICfelSkuPrice(opt *cloudprovider.CfelSkuPriceOptions) (map[string]string, error) {
	params := map[string]interface{}{
		"count":    opt.Quantity,
		"onDemand": false,
		"regionID": self.RegionId,
	}
	if len(opt.InstanceType) > 0 {
		params["flavorName"] = opt.InstanceType
		params["resourceType"] = "VM"
		params["imageUUID"] = opt.ImageId
		params["sysDiskSize"] = opt.DiskSize
		params["sysDiskType"] = opt.DiskType
	} else {
		params["resourceType"] = "EBS"
		params["diskMode"] = "VBD"
		params["diskType"] = opt.DiskType
		params["diskSize"] = opt.DiskSize
	}
	if opt.ChargeType == cloudprovider.InstanceChargeTypePostPaid { // 按量付费
		params["onDemand"] = true
	} else if opt.ChargeType == cloudprovider.InstanceChargeTypePrePaid { // 包年包月
		params["cycleType"] = strings.ToUpper(opt.FeeUnit)
		params["cycleCount"] = opt.Duration
	}
	// https://www.ctyun.cn/document/10026730/10597642
	resp, err := self.post(SERVICE_ECS, "/v4/order/new-query-price", params)
	if err != nil {
		return nil, err
	}
	var ret = priceResp{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]string)
	if opt.InstanceType == "" {
		result["dataVolumePrice"] = fmt.Sprintf("%v", ret.ReturnObj.FinalPrice)
	} else {
		// result["bootVolumePrice"] = fmt.Sprintf("%v", volumeTotalPrice)
		result["serverPrice"] = fmt.Sprintf("%v", ret.ReturnObj.FinalPrice)
	}
	return result, nil
}

func (region *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return 0, nil
}

func (region *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	return region.getSkuPrice(zoneID, instanceType, true)
}

func (region *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	return region.getSkuPrice(zoneID, instanceType, false)
}

func (region *SRegion) getSkuPrice(zoneID, instanceType string, onDemand bool) (float64, error) {
	params := map[string]interface{}{
		"count":    1,
		"onDemand": onDemand,
		"regionID": region.RegionId,
	}
	imageId, err := region.getSkuImage(zoneID, instanceType)
	if err != nil {
		return 0, err
	}
	product, err := region.getProduct()
	if err != nil {
		return 0, err
	}
	if len(product.Ebs.StorageType) == 0 {
		return 0, fmt.Errorf("storage type is empty")
	}
	var diskType = product.Ebs.StorageType[0].Type
	for i := range product.Ebs.StorageType {
		if product.Ebs.StorageType[i].Type == "SSD-genric" {
			diskType = "SSD-genric"
			break
		}
	}
	params["flavorName"] = instanceType
	params["resourceType"] = "VM"
	params["imageUUID"] = imageId
	params["sysDiskSize"] = 40
	params["sysDiskType"] = diskType
	if !onDemand { // 包年包月
		params["cycleType"] = "MONTH"
		params["cycleCount"] = 1
	}
	// https://www.ctyun.cn/document/10026730/10597642
	resp, err := region.post(SERVICE_ECS, "/v4/order/new-query-price", params)
	if err != nil {
		return 0, err
	}
	var ret = priceResp{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return 0, err
	}
	return ret.ReturnObj.FinalPrice, nil
}

func (region *SRegion) getSkuImage(zoneId, instanceType string) (string, error) {
	region.mutImage.Lock()
	defer region.mutImage.Unlock()

	if region.skuImageId != nil {
		return *region.skuImageId, nil
	}
	images, err := region.getImages(instanceType, zoneId, "CTyunOS")
	if err != nil {
		return "", err
	}
	if len(images) == 0 {
		images, err = region.getImages(instanceType, zoneId, "")
		if err != nil {
			return "", err
		}
		if len(images) == 0 {
			return "", fmt.Errorf("no match image")
		}
	}
	region.skuImageId = &images[0].ImageId
	return images[0].ImageId, nil
}

func (region *SRegion) getImages(instanceType, zoneId, imageName string) ([]SImage, error) {
	params := map[string]interface{}{
		"pageNo":     1,
		"pageSize":   1,
		"visibility": 1,
		"flavorName": instanceType,
		// "queryContent": "CTyunOS",
		"azName":   zoneId[strings.LastIndex(zoneId, "/")+1:],
		"regionID": region.RegionId,
	}
	if len(imageName) > 0 {
		params["queryContent"] = imageName
	}
	resp, err := region.list(SERVICE_IMAGE, "/v4/image/list", params)
	if err != nil {
		return nil, err
	}
	part := struct {
		ReturnObj struct {
			Images []SImage
		}
		TotalCount int
	}{}
	err = resp.Unmarshal(&part)
	if err != nil {
		return nil, err
	}
	return part.ReturnObj.Images, nil
}

func (region *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {

	return api.SkuStatusSoldout, nil
}

// 目前获取 instanceType 获取到的数据皆为可购买字段，所以实例的可购买状态皆为  available

func (region *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {

	return region.GetPrePaidStatus(zoneID, instanceType)
}

func (region *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	sku, err := region.getSkuInfo(zoneID, instanceType)
	if err != nil {
		return "", err
	}
	if sku.Available {
		return api.SkuStatusAvailable, nil
	}
	return api.SkuStatusSoldout, nil
}

func (region *SRegion) getSkuInfo(zoneId, instanceType string) (*ServerSku, error) {
	region.mut.Lock()
	defer region.mut.Unlock()

	if region.skuInfo != nil {
		return region.skuInfo, nil
	}
	params := map[string]interface{}{
		"regionID":   region.RegionId,
		"flavorName": instanceType,
		"azName":     zoneId[strings.LastIndex(zoneId, "/")+1:],
	}

	resp, err := region.post(SERVICE_ECS, "/v4/ecs/flavor/list", params)
	if err != nil {
		return nil, err
	}
	ret := struct {
		ReturnObj struct {
			FlavorList []ServerSku
		}
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	if len(ret.ReturnObj.FlavorList) == 0 {
		return nil, fmt.Errorf("no sku")
	}
	region.skuInfo = &ret.ReturnObj.FlavorList[0]
	return region.skuInfo, nil
}
