package aws

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

type SCfelSku struct {
	multicloud.SInstanceBase
	// The instance type. For more information, see Instance types (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)
	// in the Amazon EC2 User Guide.
	// InstanceType string `xml:"instanceType"`

	// cfel 新增字段

	// 地域
	ZoneId string
	// 实例系列 ,根据实例名称规则获取 (https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/instance-types.html)
	InstanceFamily   string
	InstanceCategory string

	// Indicates whether Amazon CloudWatch action based recovery is supported.
	AutoRecoverySupported bool `xml:"autoRecoverySupported"`
	// Indicates whether the instance is a bare metal instance type.
	BareMetal bool `xml:"bareMetal"`
	// Indicates whether the instance type is a burstable performance T instance
	// type. For more information, see Burstable performance instances (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/burstable-performance-instances.html).
	BurstablePerformanceSupported bool `xml:"burstablePerformanceSupported"`
	// Indicates whether the instance type is current generation.
	CurrentGeneration bool `xml:"currentGeneration"`
	// Indicates whether Dedicated Hosts are supported on the instance type.
	DedicatedHostsSupported bool `xml:"dedicatedHostsSupported"`

	// Indicates whether the instance type is eligible for the free tier.
	FreeTierEligible bool `xml:"freeTierEligible" type:"boolean"`

	// Describes the GPU accelerator settings for the instance type.
	GpuInfo GpuInfo `xml:"gpuInfo"`

	// Indicates whether On-Demand hibernation is supported.
	HibernationSupported bool `xml:"hibernationSupported" type:"boolean"`

	// The hypervisor for the instance type.
	Hypervisor string `xml:"hypervisor"`

	// Describes the instance storage for the instance type.
	InstanceStorageInfo InstanceStorageInfo `xml:"instanceStorageInfo" type:"structure"`

	// Indicates whether instance storage is supported.
	InstanceStorageSupported bool `xml:"instanceStorageSupported" type:"boolean"`

	// Describes the memory for the instance type.
	MemoryInfo struct {
		SizeInMiB int64 `xml:"sizeInMiB"`
	} `xml:"memoryInfo" type:"structure"`

	// The supported boot modes. For more information, see Boot modes (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ami-boot.html)
	// in the Amazon EC2 User Guide.
	SupportedBootModes []string `xml:"supportedBootModes>item" locationNameList:"item" type:"list" enum:"BootModeType"`

	// The supported root device types.
	SupportedRootDeviceTypes []string `xml:"supportedRootDeviceTypes>item" locationNameList:"item" type:"list" enum:"RootDeviceType"`

	// Indicates whether the instance type is offered for spot or On-Demand.
	SupportedUsageClasses []string `xml:"supportedUsageClasses>item" locationNameList:"item" type:"list" enum:"UsageClassType"`

	// The supported virtualization types.
	SupportedVirtualizationTypes []string `xml:"supportedVirtualizationTypes>item" locationNameList:"item" type:"list" enum:"VirtualizationType"`

	// Describes the vCPU configurations for the instance type.
	VCpuInfo VCpuInfo `xml:"vCpuInfo"`

	//// Describes the network settings for the instance type.
	//NetworkInfo *NetworkInfo `locationName:"networkInfo" type:"structure"`
	//
	//// Indicates whether Nitro Enclaves is supported.
	//NitroEnclavesSupport *string `locationName:"nitroEnclavesSupport" type:"string" enum:"NitroEnclavesSupport"`
	//
	//// Describes the supported NitroTPM versions for the instance type.
	//NitroTpmInfo *NitroTpmInfo `locationName:"nitroTpmInfo" type:"structure"`
	//
	//// Indicates whether NitroTPM is supported.
	//NitroTpmSupport *string `locationName:"nitroTpmSupport" type:"string" enum:"NitroTpmSupport"`
	//
	//// Describes the placement group settings for the instance type.
	//PlacementGroupInfo *PlacementGroupInfo `locationName:"placementGroupInfo" type:"structure"`
	//
	// Describes the processor.
	ProcessorInfo *ProcessorInfo `xml:"processorInfo" locationName:"processorInfo" type:"structure"`

	//// Describes the Amazon EBS settings for the instance type.
	//EbsInfo *EbsInfo `locationName:"ebsInfo" type:"structure"`
	//// Describes the FPGA accelerator settings for the instance type.
	//FpgaInfo *FpgaInfo `locationName:"fpgaInfo" type:"structure"`
	//
	//// Describes the Inference accelerator settings for the instance type.
	//InferenceAcceleratorInfo *InferenceAcceleratorInfo `locationName:"inferenceAcceleratorInfo" type:"structure"`
}

func (s Sku) GetId() string {
	return s.InstanceType
}

func (s Sku) GetName() string {
	return s.InstanceType
}

func (s Sku) GetGlobalId() string {
	return s.InstanceType
}

func (s Sku) GetCreatedAt() time.Time {
	return time.Now()
}

func (s Sku) GetDescription() string {
	return ""
}

func (s Sku) GetStatus() string {
	return ""
}

func (s Sku) Refresh() error {
	return nil
}

func (s Sku) IsEmulated() bool {
	return false
}

func (s Sku) GetSysTags() map[string]string {
	return nil
}

func (s Sku) GetTags() (map[string]string, error) {
	return nil, nil
}

func (s Sku) SetTags(tags map[string]string, replace bool) error {
	return nil
}

func (s Sku) GetZoneID() string {
	return s.ZoneId
}

func (s Sku) GetInstanceTypeFamily() string {
	return s.InstanceFamily
}

func (s Sku) GetInstanceTypeCategory() string {
	return s.InstanceCategory
}

func (s Sku) GetPrepaidStatus() string {
	return ""
}

func (s Sku) GetPostpaidStatus() string {
	return ""
}

func (s Sku) GetCpuArch() string {
	if s.ProcessorInfo != nil {
		return strings.Join(s.ProcessorInfo.SupportedArchitectures, ",")
	}
	return ""
}

func (s Sku) GetCpuCoreCount() int {
	return int(s.VCpuInfo.DefaultVCpus)
}

func (s Sku) GetMemorySizeMB() int {
	return int(s.MemoryInfo.SizeInMiB)
}

func (s Sku) GetOsName() string {
	return ""
}

func (s Sku) GetSysDiskResizable() bool {
	return false
}

func (s Sku) GetSysDiskType() string {
	return ""
}

func (s Sku) GetSysDiskMinSizeGB() int {
	return int(s.InstanceStorageInfo.TotalSizeInGB * 1024)
}

func (s Sku) GetSysDiskMaxSizeGB() int {
	return int(s.InstanceStorageInfo.TotalSizeInGB * 1024)
}

func (s Sku) GetAttachedDiskType() string {
	return ""
}

func (s Sku) GetAttachedDiskSizeGB() int {
	return 0
}

func (s Sku) GetAttachedDiskCount() int {
	return 0
}

func (s Sku) GetDataDiskTypes() string {
	return ""
}

func (s Sku) GetDataDiskMaxCount() int {
	return 0
}

func (s Sku) GetNicType() string {
	return ""
}

func (s Sku) GetNicMaxCount() int {
	return 0
}

func (s Sku) GetGpuAttachable() bool {
	return s.GpuInfo.TotalGpuMemoryInMiB != 0
}

func (s Sku) GetGpuSpec() string {
	if s.GpuInfo.TotalGpuMemoryInMiB != 0 && len(s.GpuInfo.Gpus) > 0 {
		if s.GpuInfo.Gpus[0].Manufacturer != "" && s.GpuInfo.Gpus[0].Name != "" {
			return fmt.Sprintf("%s %s", s.GpuInfo.Gpus[0].Manufacturer, s.GpuInfo.Gpus[0].Name)
		}
		return ""
	}
	return ""
}

func (s Sku) GetGpuCount() string {
	if s.GpuInfo.TotalGpuMemoryInMiB != 0 && len(s.GpuInfo.Gpus) > 0 {
		return strconv.Itoa(s.GpuInfo.Gpus[0].Count)
	}
	return "0"
}

func (s Sku) GetGpuMaxCount() int {
	if s.GpuInfo.TotalGpuMemoryInMiB != 0 && len(s.GpuInfo.Gpus) > 0 {
		return s.GpuInfo.Gpus[0].Count
	}
	return 0
}

func (s Sku) Delete() error {
	return nil
}

func (s *Sku) GetGPUMemorySizeMB() int {
	return s.GpuInfo.TotalGpuMemoryInMiB
}

func (s *Sku) GetIsBareMetal() bool {
	return s.BareMetal
}

func (s *Sku) GetSpotpaidStatus() string {
	if slices.Contains(s.SupportedUsageClasses, "spot") {
		return api.SkuStatusAvailable
	}
	return api.CfelSkuStatusAbandon
}

type GpuInfo struct {
	Gpus []struct {
		Count        int    `xml:"count"`
		Manufacturer string `xml:"manufacturer"`
		MemoryInfo   struct {
			SizeInMiB int `xml:"sizeInMiB"`
		} `xml:"memoryInfo"`
		Name string `xml:"name"`
	} `xml:"gpus>item"`
	//Gpus                string `xml:"gpus"`
	TotalGpuMemoryInMiB int `xml:"totalGpuMemoryInMiB"`
}
type VCpuInfo struct {
	// The default number of cores for the instance type.
	DefaultCores int64 `locationName:"defaultCores" type:"integer" xml:"defaultCores"`
	// The default number of threads per core for the instance type.
	DefaultThreadsPerCore int64 `locationName:"defaultThreadsPerCore" type:"integer" xml:"defaultThreadsPerCore"`
	// The default number of vCPUs for the instance type.
	DefaultVCpus int64 `locationName:"defaultVCpus" type:"integer" xml:"defaultVCpus"`
	// The valid number of cores that can be configured for the instance type.
	ValidCores []int64 `locationName:"validCores" locationNameList:"item" type:"list" xml:"validCores>item"`
	// The valid number of threads per core that can be configured for the instance
	// type.
	ValidThreadsPerCore []int64 `locationName:"validThreadsPerCore" locationNameList:"item" type:"list" xml:"validThreadsPerCore>item"`
	// contains filtered or unexported fields
}

type ProcessorInfo struct {

	// The architectures supported by the instance type.
	SupportedArchitectures []string `xml:"supportedArchitectures>item" locationNameList:"item" type:"list" enum:"ArchitectureType"`

	// Indicates whether the instance type supports AMD SEV-SNP. If the request
	// returns amd-sev-snp, AMD SEV-SNP is supported. Otherwise, it is not supported.
	// For more information, see AMD SEV-SNP (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/sev-snp.html).
	SupportedFeatures []string `xml:"supportedFeatures>item" locationNameList:"item" type:"list" enum:"SupportedAdditionalProcessorFeature"`

	// The speed of the processor, in GHz.
	SustainedClockSpeedInGhz float64 `xml:"sustainedClockSpeedInGhz" type:"double"`
	// contains filtered or unexported fields
}

type InstanceStorageInfo struct {

	// Describes the disks that are available for the instance type.
	Disks []struct {
		// The number of disks with this configuration.
		Count int64 `xml:"count" type:"integer"`

		// The size of the disk in GB.
		SizeInGB int64 `xml:"sizeInGB" type:"long"`

		// The type of disk.
		Type string `xml:"type" type:"string" enum:"DiskType"`
		// contains filtered or unexported fields
	} `xml:"disks>item" locationNameList:"item" type:"list"`

	// Indicates whether data is encrypted at rest.
	EncryptionSupport string `xml:"encryptionSupport" type:"string" enum:"InstanceStorageEncryptionSupport"`

	// Indicates whether non-volatile memory express (NVMe) is supported.
	NvmeSupport string `xml:"nvmeSupport" type:"string" enum:"EphemeralNvmeSupport"`

	// The total size of the disks, in GB.
	TotalSizeInGB int64 `xml:"totalSizeInGB" type:"long"`
	// contains filtered or unexported fields
}

// DescribeInstanceTypesAll 获取某一个 region 下所有 InstanceTypes
func (self *SRegion) DescribeInstanceTypesAll() ([]Sku, error) {
	ret := []Sku{}
	var nextToken string
	for {
		parts, _nextToken, err := self.CfelDescribeInstanceTypes("", nextToken)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parts...)
		if len(_nextToken) == 0 || len(parts) == 0 {
			break
		}
		nextToken = _nextToken
	}

	return ret, nil
}

type InstanceTypeOffering struct {

	// The instance type. For more information, see Instance types (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)
	// in the Amazon EC2 User Guide.
	InstanceType string `xml:"instanceType" type:"string" enum:"InstanceType"`

	// The identifier for the location. This depends on the location type. For example,
	// if the location type is region, the location is the Region code (for example,
	// us-east-2.)
	Location string `xml:"location" type:"string"`

	// The location type.
	LocationType string `xml:"locationType" type:"string" enum:"LocationType"`
}

// DescribeInstanceTypeAvailable
// 查询指定 instanceType 在zone 下是否售卖
func (self *SRegion) DescribeInstanceTypeAvailable(instanceType, zone string) (bool, error) {
	self.mut.Lock()
	defer self.mut.Unlock()

	if self.instanceAvailable != nil {
		return *self.instanceAvailable, nil
	}
	params := map[string]string{
		"LocationType": "availability-zone",
	}

	params["Filter.1.Name"] = "location"
	params["Filter.1.Value.1"] = zone
	params["Filter.2.Name"] = "instance-type"
	params["Filter.2.Value.1"] = instanceType

	ret := struct {
		InstanceTypeSet []InstanceTypeOffering `xml:"instanceTypeOfferingSet>item"`
		NextToken       string                 `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeInstanceTypeOfferings", params, &ret)
	if err != nil {
		return false, err
	}
	var res = len(ret.InstanceTypeSet) > 0
	self.instanceAvailable = &res
	return res, nil
}

// DescribeReservedInstancesOfferings
// 查询指定 instanceType 在 zone 预付费 是否售卖
// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeReservedInstancesOfferings.html
func (self *SRegion) DescribeReservedInstancesOfferings(instanceType, zone string) (bool, error) {
	params := map[string]string{
		// "AvailabilityZone": zone,
		"OfferingType":       "All upfront",
		"OfferingClass":      "standard",
		"ProductDescription": "Linux/UNIX",
		"InstanceTenancy":    "default",
		"MaxInstanceCount":   "100",
	}

	// params["Filter.1.Name"] = "marketplace"
	// params["Filter.1.Value.1"] = "true"
	var idx = 1
	params[fmt.Sprintf("Filter.%d.Name", idx)] = "scope"
	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = "Region"
	idx++

	params[fmt.Sprintf("Filter.%d.Name", idx)] = "instance-type"
	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = instanceType
	idx++

	var nextToken string
	var res []InstanceTypeOffering
	for {
		if len(nextToken) > 0 {
			params["NextToken"] = nextToken
		}
		ret := struct {
			InstanceTypeSet []InstanceTypeOffering `xml:"reservedInstancesOfferingsSet>item"`
			NextToken       string                 `xml:"nextToken"`
		}{}
		_ = self.ec2Request("DescribeReservedInstancesOfferings", params, &ret)
		if len(ret.InstanceTypeSet) > 0 {
			res = ret.InstanceTypeSet
			break
		}
		nextToken = ret.NextToken
		if ret.NextToken == "" {
			break
		}
	}

	return len(res) > 0, nil
}

func (self *SRegion) DescribeInstanceTypeOfferings(nextToken string) ([]InstanceTypeOffering, string, error) {
	params := map[string]string{
		"LocationType": "availability-zone",
	}

	if len(nextToken) > 0 {
		params["NextToken"] = nextToken
	}
	//idx := 1
	//if len(arch) > 0 {
	//	params[fmt.Sprintf("Filter.%d.Name", idx)] = "processor-info.supported-architecture"
	//	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = arch
	//	idx++
	//}
	ret := struct {
		InstanceTypeSet []InstanceTypeOffering `xml:"instanceTypeOfferingSet>item"`
		NextToken       string                 `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeInstanceTypeOfferings", params, &ret)
	if err != nil {
		return nil, "", err
	}
	return ret.InstanceTypeSet, ret.NextToken, nil
}

// DescribeInstanceTypeOfferingsAll 获取某一个 region 下所有 InstanceTypeOfferings
func (self *SRegion) DescribeInstanceTypeOfferingsAll() ([]InstanceTypeOffering, error) {
	ret := []InstanceTypeOffering{}
	var nextToken string
	for {
		parts, _nextToken, err := self.DescribeInstanceTypeOfferings(nextToken)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parts...)
		if len(_nextToken) == 0 || len(parts) == 0 {
			break
		}
		nextToken = _nextToken
	}

	return ret, nil
}

// GetRegionZoneInstanceType 获取带zone 的 instanceType
func (self *SRegion) GetRegionZoneInstanceType() ([]Sku, error) {

	existMap := make(map[string]*Sku)
	instanceTypesAll, err := self.DescribeInstanceTypesAll()
	if err != nil {
		return nil, errors.Wrapf(err, "DescribeInstanceTypesAll")
	}
	for _, instanceType := range instanceTypesAll {
		_instanceType := instanceType
		existMap[instanceType.InstanceType] = &_instanceType
	}

	skus := make([]Sku, 0)
	offerings, err := self.DescribeInstanceTypeOfferingsAll()
	if err != nil {
		return nil, errors.Wrapf(err, "DescribeInstanceTypeOfferingsAll")
	}
	for _, offering := range offerings {
		sku, ok := existMap[offering.InstanceType]
		if ok {
			sku.ZoneId = offering.Location
			skus = append(skus, *sku)
		}
	}
	return skus, nil
}

func (self *SRegion) GetICfelSkus() ([]cloudprovider.ICfelCloudSku, error) {
	skus, err := self.GetRegionZoneInstanceType()
	if err != nil {
		return nil, errors.Wrapf(err, "GetInstanceTypes")
	}
	var ret []cloudprovider.ICfelCloudSku
	for i := range skus {
		ret = append(ret, &skus[i])
	}
	return ret, nil
}

func getParams(filters map[string]string) []ProductFilter {
	params := []ProductFilter{}

	for k, v := range filters {
		params = append(params, ProductFilter{
			Type:  "TERM_MATCH",
			Field: k,
			Value: v,
		})
	}

	return params
}

func (self *SRegion) GetICfelSkuPrice(opt *cloudprovider.CfelSkuPriceOptions) (map[string]string, error) {

	var volumeTotalPrice float64

	var wg sync.WaitGroup
	wg.Add(1)

	var diskMap = map[string]string{
		"gp2":      "General Purpose",
		"gp3":      "General Purpose",
		"standard": "Magnetic",
		"io1":      "Provisioned IOPS",
		"io2":      "Provisioned IOPS",
		"st1":      "Throughput Optimized HDD",
		"sc1":      "Cold HDD",
	}

	go func() {
		filters := map[string]string{
			"regionCode":    self.RegionId,
			"volumeType":    diskMap[opt.DiskType],
			"volumeApiName": opt.DiskType,
		}
		parts, _, err := self.GetProducts("AmazonEC2", getParams(filters), "")
		if err != nil {

		}
		if len(parts) <= 0 {
			return
		}
		var volumePrice float64
		for _, val := range parts[0].Terms.OnDemand {
			for _, vv := range val.PriceDimensions {
				if vv.Unit == "GB-Mo" {
					volumePrice = vv.PricePerUnit.Usd
				}
			}
		}
		if opt.FeeUnit == "month" {
			volumeTotalPrice = volumePrice * float64(opt.DiskSize) * float64(opt.Duration)
		} else if opt.FeeUnit == "year" {
			volumeTotalPrice = volumePrice * float64(opt.DiskSize) * float64(opt.Duration) * 12
		} else {
			volumeTotalPrice = volumePrice * float64(opt.DiskSize) * float64(opt.Duration) / 730
		}
		defer wg.Done()
	}()

	var serverTotalPrice float64

	if opt.ChargeType == cloudprovider.InstanceChargeTypeSpotPaid {
		price, err := self.DescribeSpotPriceHistory(opt.ZoneId, opt.InstanceType)
		if err == nil && len(price) > 0 {
			serverTotalPrice, _ = strconv.ParseFloat(strings.ReplaceAll(price[0].SpotPrice, ",", ""), 64)
		}
	} else {
		if opt.InstanceType != "" {
			img, err := self.GetImage(opt.ImageId)
			if err != nil {
				return nil, err
			}
			var ops string
			var dist = img.GetOsDist()
			if strings.Contains(dist, "Ubuntu") {
				ops = "Ubuntu Pro"
			} else if strings.Contains(dist, "Windows") {
				ops = "Windows"
			} else if strings.Contains(dist, "SUSE") {
				ops = "SUSE"
			} else if strings.Contains(dist, "RHEL") {
				ops = "RHEL"
			} else {
				ops = "Linux"
			}
			filters := map[string]string{
				"regionCode":      self.RegionId,
				"operatingSystem": ops,
				"licenseModel":    "No License required",
				"productFamily":   "Compute Instance",
				// "operation":       img.UsageOperation,
				"tenancy":        "Shared",
				"preInstalledSw": "NA",
				"capacitystatus": "Used",
				"instanceType":   opt.InstanceType,
			}

			opt.FeeUnit = strings.ToLower(opt.FeeUnit)
			//var volumePrice float64

			var ret = []SInstanceType{}
			var nextToken string
			for {
				parts, _nextToken, err := self.GetProducts("AmazonEC2", getParams(filters), nextToken)
				if err != nil {
					return nil, err
				}
				ret = append(ret, parts...)
				if len(_nextToken) == 0 || len(parts) == 0 {
					break
				}
				nextToken = _nextToken
			}

			// rr, _ := json.Marshal(ret)
			// fmt.Println(string(rr))
			// if len(rr) > 0 {
			// 	ioutil.WriteFile("./price.json", rr, os.ModeAppend)
			// }

			if len(ret) > 0 {
				// 小时价
				var hourPrice float64
				for _, val := range ret[0].Terms.OnDemand {
					for _, vv := range val.PriceDimensions {
						if vv.Unit == "Hrs" {
							hourPrice = vv.PricePerUnit.Usd
						}
					}
				}
				// 一年价和三年价
				// var oneYearPrice, threeYearPrice float64
				// for _, term := range ret[0].Terms.Reserved {
				// 	attributes := term.TermAttributes
				// 	if attributes.LeaseContractLength == "1yr" && attributes.OfferingClass == "standard" && attributes.PurchaseOption == "All Upfront" {
				// 		for _, dimension := range term.PriceDimensions {
				// 			if dimension.Unit == "Quantity" {
				// 				oneYearPrice = dimension.PricePerUnit.Usd
				// 			}
				// 		}
				// 	}
				// 	if attributes.LeaseContractLength == "3yr" && attributes.OfferingClass == "standard" && attributes.PurchaseOption == "All Upfront" {
				// 		for _, dimension := range term.PriceDimensions {
				// 			if dimension.Unit == "Quantity" {
				// 				threeYearPrice = dimension.PricePerUnit.Usd
				// 			}
				// 		}
				// 	}
				// }

				// 包年包月都按小时价算
				if opt.FeeUnit == "month" {
					serverTotalPrice = hourPrice * 730 * float64(opt.Quantity) * float64(opt.Duration)
				} else if opt.FeeUnit == "year" {
					serverTotalPrice = hourPrice * 730 * 12 * float64(opt.Quantity) * float64(opt.Duration)
				} else {
					serverTotalPrice = hourPrice * float64(opt.Quantity) * float64(opt.Duration)
				}
			}
		}
	}

	wg.Wait()

	var result = map[string]string{
		// "bootVolumePrice": fmt.Sprintf("%v", volumeTotalPrice),
		// "serverPrice":     fmt.Sprintf("%v", serverTotalPrice),
		"currency": "USD",
	}
	if opt.InstanceType == "" {
		result["dataVolumePrice"] = fmt.Sprintf("%v", volumeTotalPrice)
	} else {
		result["bootVolumePrice"] = fmt.Sprintf("%v", volumeTotalPrice)
		result["serverPrice"] = fmt.Sprintf("%v", serverTotalPrice)
	}
	return result, nil

}

func (self *SRegion) ListPriceLists() (*SInstanceType, error) {
	filters := map[string]string{
		"regionCode": self.RegionId,
		//"operatingSystem": "Linux",
		//"licenseModel":    "No License required",
		//"productFamily":   "Compute Instance",
		//"operation":       "RunInstances",
		//"preInstalledSw":  "NA",
		//"tenancy":         "Shared",
		//"capacitystatus":  "Used",
		//"instanceType":    instanceType,
	}

	params := []ProductFilter{}

	for k, v := range filters {
		params = append(params, ProductFilter{
			Type:  "TERM_MATCH",
			Field: k,
			Value: v,
		})
	}

	ret := []SInstanceType{}
	var nextToken string
	for {
		parts, _nextToken, err := self.GetProducts("AmazonEC2", params, nextToken)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parts...)
		if len(_nextToken) == 0 || len(parts) == 0 {
			break
		}
		nextToken = _nextToken
	}
	if len(ret) > 0 {
		return &ret[0], nil
	} else {
		return nil, nil
	}
}

var noResultErr = errors.Errorf("No Result")
// Price

func (self *SRegion) GetInstanceTypePrice(instanceType string) (*SInstanceType, error) {
	self.instanceMut.Lock()
	defer self.instanceMut.Unlock()

	if self.instanceType != nil {
		return self.instanceType, nil
	}
	filters := map[string]string{
		"regionCode":     self.RegionId,
		"operation":      "RunInstances",
		"capacitystatus": "Used",
		"tenancy":        "Shared",
		"instanceType":   instanceType,
	}

	params := []ProductFilter{}

	for k, v := range filters {
		params = append(params, ProductFilter{
			Type:  "TERM_MATCH",
			Field: k,
			Value: v,
		})
	}

	ret := []SInstanceType{}
	var nextToken string
	for {
		parts, _nextToken, err := self.GetProducts("AmazonEC2", params, nextToken)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parts...)
		if len(_nextToken) == 0 || len(parts) == 0 {
			break
		}
		nextToken = _nextToken
	}
	if len(ret) > 0 {
		self.instanceType = &ret[0]
		return &ret[0], nil
	} else {
		return nil, noResultErr
	}
}

type SpotPrice struct {

	// The Availability Zone.
	AvailabilityZone string `xml:"availabilityZone" type:"string"`

	// The instance type.
	InstanceType string `xml:"instanceType" type:"string" enum:"InstanceType"`

	// A general description of the AMI.
	ProductDescription string `xml:"productDescription" type:"string" enum:"RIProductDescription"`

	// The maximum price per unit hour that you are willing to pay for a Spot Instance.
	// We do not recommend using this parameter because it can lead to increased
	// interruptions. If you do not specify this parameter, you will pay the current
	// Spot price.
	//
	// If you specify a maximum price, your instances will be interrupted more frequently
	// than if you do not specify this parameter.
	SpotPrice string `xml:"spotPrice" type:"string"`

	// The date and time the request was created, in UTC format (for example, YYYY-MM-DDTHH:MM:SSZ).
	Timestamp time.Time `xml:"timestamp" type:"timestamp"`
	// contains filtered or unexported fields
}

func (self *SRegion) DescribeSpotPriceHistory(zone, instanceType string) ([]*SpotPrice, error) {

	currentTime := time.Now().UTC()
	startTime := currentTime.Add(-24 * 60 * 60 * time.Second)
	params := map[string]string{
		"AvailabilityZone": zone,
		"StartTime":        startTime.Format("2006-01-02T15:04:05Z"),
		"EndTime":          currentTime.Format("2006-01-02T15:04:05Z"),
	}

	params["ProductDescription.1"] = "Linux/UNIX"
	params["InstanceType.1"] = instanceType

	ret := struct {
		SpotPriceHistorySet []*SpotPrice `xml:"spotPriceHistorySet>item"`
		NextToken           string       `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeSpotPriceHistory", params, &ret)
	if err != nil {
		return nil, err
	}
	return ret.SpotPriceHistorySet, nil
}

func (self *SRegion) GetSpotPostPaidPrice(zoneID, instanceType string) (float64, error) {
	spotPriceHistory, err := self.DescribeSpotPriceHistory(zoneID, instanceType)
	if err != nil {
		return -1, err
	}
	if len(spotPriceHistory) > 0 {
		newSpotPrice, err := strconv.ParseFloat(spotPriceHistory[0].SpotPrice, 64)
		if err != nil {
			return 0, err
		}
		return newSpotPrice, nil
	}
	return -1, nil
}

func (self *SRegion) GetPostPaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetInstanceTypePrice(instanceType)
	if err != nil {
		if err == noResultErr {
			return -1,nil
		}
		return -1, errors.Wrapf(err, "GetInstanceTypePrice")
	}
	var value float64
	for _, term := range price.Terms.OnDemand {
		for _, dimension := range term.PriceDimensions {
			value = dimension.PricePerUnit.Usd
		}
	}
	return value, nil
}

// GetPrePaidPrice
// PurchaseOption: No Upfront  All Upfront  Partial Upfront
/**
在AWS（Amazon Web Services）中，包年包月的价格模型有三种支付类型，它们分别是：

No Upfront (无预付费)： 在这种支付类型下，无需提前支付任何费用，而是在每个计费周期结束时按月支付使用费用。这是一种相对灵活的选项，适合那些不愿意一次性支付较大费用的用户。

All Upfront (全部预付费)： 在这种支付类型下，需要在合同开始时一次性支付整个合同期的费用，这样您就能够获得更大的折扣。虽然一开始支付较多费用，但在合同期内不需要再支付额外费用。

Partial Upfront (部分预付费)： 这是一种折中的选择，您需要在合同开始时支付部分费用，而余下的费用则按月支付。这样可以降低一次性支付的压力，并且仍然能够获得一定的折扣。
*/
func (self *SRegion) GetPrePaidPrice(zoneID, instanceType string) (float64, error) {
	price, err := self.GetInstanceTypePrice(instanceType)
	if err != nil {
		if err == noResultErr {
			return -1,nil
		}
		return -1, errors.Wrapf(err, "GetInstanceTypePrice")
	}
	var value float64
	var hourPrice float64
	for _, val := range price.Terms.OnDemand {
		for _, vv := range val.PriceDimensions {
			if vv.Unit == "Hrs" {
				hourPrice = vv.PricePerUnit.Usd
			}
		}
	}
	value = 730 * hourPrice
	return value, nil
}

func (self *SRegion) GetSpotPostPaidStatus(zoneID, instanceType string) (string, error) {
	spotPriceHistory, err := self.DescribeSpotPriceHistory(zoneID, instanceType)

	if err != nil {
		return api.SkuStatusSoldout, errors.Wrapf(err, "DescribeSpotPriceHistory")
	}
	if len(spotPriceHistory) > 0 {
		return api.SkuStatusAvailable, nil
	}

	return api.SkuStatusSoldout, nil
}

func (self *SRegion) GetPostPaidStatus(zoneID, instanceType string) (string, error) {
	available, err := self.DescribeInstanceTypeAvailable(instanceType, zoneID)
	if err != nil {
		return api.SkuStatusSoldout, errors.Wrapf(err, "DescribeInstanceTypeAvailable")
	}
	if available {
		return api.SkuStatusAvailable, nil
	}
	return api.SkuStatusSoldout, nil
}

func (self *SRegion) GetPrePaidStatus(zoneID, instanceType string) (string, error) {
	return self.GetPostPaidStatus(zoneID, instanceType) // aws只有包年而且和普通的按量付费是两个东西，所以包年包月直接按小时算，和后付费一样
	available, err := self.DescribeReservedInstancesOfferings(instanceType, zoneID)
	if err != nil {
		return api.SkuStatusSoldout, errors.Wrapf(err, "DescribeReservedInstancesOfferings")
	}
	if available {
		return api.SkuStatusAvailable, nil
	}
	return api.SkuStatusSoldout, nil
}

func (self *SRegion) CfelDescribeInstanceTypes(arch string, nextToken string) ([]Sku, string, error) {
	params := map[string]string{}
	if len(nextToken) > 0 {
		params["NextToken"] = nextToken
	}
	idx := 1
	if len(arch) > 0 {
		params[fmt.Sprintf("Filter.%d.Name", idx)] = "processor-info.supported-architecture"
		params[fmt.Sprintf("Filter.%d.Value.1", idx)] = arch
		idx++
	}
	ret := struct {
		InstanceTypeSet []Sku  `xml:"instanceTypeSet>item"`
		NextToken       string `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeInstanceTypes", params, &ret)
	if err != nil {
		return nil, "", err
	}

	// 获取 InstanceFamilies
	skus := make([]Sku, 0)
	for _, instanceType := range ret.InstanceTypeSet {
		instanceType.InstanceCategory = getInstanceCategory(instanceType.InstanceType)
		instanceType.InstanceFamily = getInstanceFamily(instanceType.InstanceType)
		skus = append(skus, instanceType)
	}
	return skus, ret.NextToken, nil
}

// getInstanceCategory 根据 instanceType 命名规则获取 instanceFamilies
// 具体规则详见: (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)
func getInstanceCategory(instanceType string) string {
	re := regexp.MustCompile(`^([a-zA-Z]+)`)
	match := re.FindStringSubmatch(instanceType)

	if len(match) > 1 {
		switch match[1] {
		case "c":
			return "Compute optimized"
		case "d":
			return "Dense storage"
		case "f":
			return "FPGA"
		case "g":
			return "Graphics intensive"
		case "hpc":
			return "High performance computing"
		case "i":
			return "Storage optimized"
		case "im":
			return "Storage optimized with a one to four ratio of vCPU to memory"
		case "is":
			return "Storage optimized with a one to six ratio of vCPU to memory"
		case "inf":
			return "AWS Inferentia"
		case "m":
			return "General purpose"
		case "mac":
			return "macOS"
		case "p":
			return "GPU accelerated"
		case "r":
			return "Memory optimized"
		case "t":
			return "Burstable performance"
		case "trn":
			return "AWS Trainium"
		case "u":
			return "High memory"
		case "vt":
			return "Video transcoding"
		case "x":
			return "Memory intensive"
		default:
			return match[1]
		}
	}
	return ""
}

func getInstanceFamily(instanceType string) string {
	re := regexp.MustCompile(`^([^.]*)`)
	match := re.FindStringSubmatch(instanceType)

	if len(match) > 1 {
		return match[1]
	}
	return ""
}
