// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package huawei

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

// Verify that *SInstanceType implements ICloudSku、ICloudSkuUltra
var _ cloudprovider.ICloudSku = (*SInstanceType)(nil)
var _ cloudprovider.ICloudSkuUltra = (*SInstanceType)(nil)

// https://support.huaweicloud.com/api-ecs/zh-cn_topic_0020212656.html
type SInstanceType struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Vcpus        string       `json:"vcpus"`
	RamMB        int          `json:"ram"`            // 内存大小
	OSExtraSpecs OSExtraSpecs `json:"os_extra_specs"` // 扩展规格
	ZoneID       string       `json:"zoneId"`
}

type OSExtraSpecs struct {
	EcsPerformancetype             string `json:"ecs:performancetype"`                 // 云服务器规格
	HwNumaNodes                    string `json:"hw:numa_nodes"`                       // 主机的物理cpu数量。（该字段是否返回根据云服务器规格而定）
	ResourceType                   string `json:"resource_type"`                       // 资源类型。resource_type是为了区分云服务器的物理主机类型。
	HpetSupport                    string `json:"hpet_support"`                        // 云服务器高精度时钟是否开启，开启为true，否则为false。（该字段是否返回根据云服务器规格而定）
	InstanceVnicType               string `json:"instance_vnic:type"`                  // 网卡类型，值固定为“enhanced”，表示使用增强型网络的资源创建云服务器。
	InstanceVnicInstanceBandwidth  string `json:"instance_vnic:instance_bandwidth"`    // 最大带宽，单位Mbps，最大值为10000。
	InstanceVnicMaxCount           string `json:"instance_vnic:max_count"`             // 最大网卡个数，最大为4。
	QuotaLocalDisk                 string `json:"quota:local_disk"`                    // 磁盘增强型特有字段。
	QuotaNvmeSsd                   string `json:"quota:nvme_ssd"`                      // 超高I/O型特有字段。
	ExtraSpecIoPersistentGrant     string `json:"extra_spec:io:persistent_grant"`      // 密集存储D1型特有字段。 是否支持持久化，值为true。
	EcsGeneration                  string `json:"ecs:generation"`                      // 弹性云服务器类型的代数。
	EcsVirtualizationEnvTypes      string `json:"ecs:virtualization_env_types"`        // 虚拟化类型。
	PciPassthroughEnableGpu        string `json:"pci_passthrough:enable_gpu"`          // 显卡是否直通。
	PciPassthroughGpuSpecs         string `json:"pci_passthrough:gpu_specs"`           // G1型和G2型云服务器应用的技术，包括GPU虚拟化和GPU直通。
	PciPassthroughAlias            string `json:"pci_passthrough:alias"`               // PCI直通设备信息，格式为PCI设备名称:数量。多个设备信息以逗号隔开。 例如nvidia-a30:1，表示携带一张A30的GPU。
	CondOperationStatus            string `json:"cond:operation:status"`               // 此参数是Region级配置，某个AZ没有在cond:operation:az参数中配置时默认使用此参数的取值。不配置或无此参数时等同于“normal”。
	CondOperationAz                string `json:"cond:operation:az"`                   // 此参数是AZ级配置，某个AZ没有在此参数中配置时默认使用cond:operation:status参数的取值。
	QuotaMaxRate                   string `json:"quota:max_rate"`                      // 最大带宽
	QuotaMinRate                   string `json:"quota:min_rate"`                      // 基准带宽
	QuotaMaxPps                    string `json:"quota:max_pps"`                       // 内网最大收发包能力
	CondOperationChargeStop        string `json:"cond:operation:charge:stop"`          // 关机是否收费
	CondOperationCharge            string `json:"cond:operation:charge"`               // 计费类型 计费场景，不配置时都支持
	CondSpotOperationAz            string `json:"cond:spot:operation:az"`              // spot售卖信息请使用 查询规格销售策略 接口查询。
	CondOperationRoles             string `json:"cond:operation:roles"`                // 允许的角色
	CondSpotOperationStatus        string `json:"cond:spot:operation:status"`          // spot售卖信息请使用 查询规格销售策略 接口查询。
	CondNetwork                    string `json:"cond:network"`                        // 网络约束
	CondStorage                    string `json:"cond:storage"`                        // 存储约束
	CondComputeLiveResizable       string `json:"cond:compute:live_resizable"`         // 计算约束
	CondCompute                    string `json:"cond:compute"`                        // 计算约束
	EcsInstanceArchitecture        string `json:"ecs:instance_architecture"`           // 该规格对应的CPU架构，且仅鲲鹏实例架构规格返回该字段。
	InfoGpuName                    string `json:"info:gpu:name"`                       // GPU显卡数量和名称。
	InfoAscendName                 string `json:"info:ascend:name"`                    // 华为昇腾AI处理器数量和名称
	InfoCpuName                    string `json:"info:cpu:name"`                       // CPU名称
	QuotaGpu                       string `json:"quota:gpu"`                           // GPU显卡名称。
	QuotaVifMaxNum                 string `json:"quota:vif_max_num"`                   // 云服务器最多支持绑定的弹性网卡个数。
	QuotaSubNetworkInterfaceMaxNum string `json:"quota:sub_network_interface_max_num"` // 云服务器最多支持绑定的辅助弹性网卡个数。
}

func (S SInstanceType) GetIsBareMetal() bool {
	return false
}

func (S SInstanceType) GetGPUMemorySizeMB() int {

	//  N卡信息记录在 InfoGpuName
	//  HUAWEI卡记录在了 InfoAscendName
	if S.OSExtraSpecs.InfoAscendName != "" {
		S.OSExtraSpecs.InfoGpuName = handleSpecialCharacters(S.OSExtraSpecs.InfoAscendName)
	}

	memInfo := ""
	// 以斜杠为分隔符进行分割
	parts := strings.Split(S.OSExtraSpecs.InfoGpuName, "/")

	// 获取斜杠后面的字符，去除首尾空格
	if len(parts) > 1 {
		memInfo = strings.TrimSpace(parts[1])
	} else {
		return 0
	}

	// 定义更严格的正则表达式
	re := regexp.MustCompile(`(\d+)G`)

	// 在字符串中查找匹配项
	match := re.FindStringSubmatch(memInfo)

	// 输出匹配结果
	if len(match) > 1 {
		atoi, err := strconv.Atoi(match[1])
		if err != nil {
			return 0
		}
		gpuCount := S.GetGpuCount()
		gpuCountInt, err := strconv.Atoi(gpuCount)
		if err != nil {
			return 0
		}

		return atoi * gpuCountInt * 1024
	} else {
		return 0
	}
	return 0
}

func (S SInstanceType) GetId() string {
	return S.ID
}

func (S SInstanceType) GetName() string {
	return S.Name
}

func (S SInstanceType) GetGlobalId() string {
	return S.ID
}

func (S SInstanceType) GetCreatedAt() time.Time {
	return time.Now()
}

func (S SInstanceType) GetDescription() string {
	return ""
}

func (S SInstanceType) GetStatus() string {
	return ""
}

func (S SInstanceType) Refresh() error {
	return nil
}

func (S SInstanceType) IsEmulated() bool {
	return false
}

func (S SInstanceType) GetSysTags() map[string]string {
	return nil
}

func (S SInstanceType) GetTags() (map[string]string, error) {
	return nil, nil
}

func (S SInstanceType) SetTags(tags map[string]string, replace bool) error {
	return nil
}

func (S SInstanceType) GetZoneID() string {
	return S.ZoneID
}

func (S SInstanceType) GetInstanceTypeFamily() string {
	return S.OSExtraSpecs.EcsGeneration
}

func (S SInstanceType) GetInstanceTypeCategory() string {
	return S.OSExtraSpecs.EcsPerformancetype
}

func (S SInstanceType) GetPrepaidStatus() string {
	return ""
}

func (S SInstanceType) GetPostpaidStatus() string {
	return ""
}

func (S SInstanceType) GetCpuArch() string {
	return S.OSExtraSpecs.EcsInstanceArchitecture
}

func (S SInstanceType) GetCpuCoreCount() int {
	vCpus, err := strconv.Atoi(S.Vcpus)
	if err != nil {
		return 0
	}
	return vCpus
}

func (S SInstanceType) GetMemorySizeMB() int {
	return S.RamMB
}

func (S SInstanceType) GetOsName() string {
	return ""
}

func (S SInstanceType) GetSysDiskResizable() bool {
	return false
}

func (S SInstanceType) GetSysDiskType() string {
	return ""
}

func (S SInstanceType) GetSysDiskMinSizeGB() int {
	return 0
}

func (S SInstanceType) GetSysDiskMaxSizeGB() int {
	return 0
}

func (S SInstanceType) GetAttachedDiskType() string {
	return ""
}

func (S SInstanceType) GetAttachedDiskSizeGB() int {
	return 0
}

func (S SInstanceType) GetAttachedDiskCount() int {
	return 0
}

func (S SInstanceType) GetDataDiskTypes() string {
	return ""
}

func (S SInstanceType) GetDataDiskMaxCount() int {
	return 0
}

func (S SInstanceType) GetNicType() string {
	return S.OSExtraSpecs.InstanceVnicType
}

func (S SInstanceType) GetNicMaxCount() int {
	count, err := strconv.Atoi(S.OSExtraSpecs.InstanceVnicMaxCount)
	if err != nil {
		return 0
	}
	return count
}

func (S SInstanceType) GetGpuAttachable() bool {
	//  N卡信息记录在 InfoGpuName
	//  HUAWEI卡记录在了 InfoAscendName
	return S.OSExtraSpecs.InfoGpuName != "" || S.OSExtraSpecs.InfoAscendName != ""
}

func handleSpecialCharacters(info string) string {
	// 检查是否包含单引号
	if strings.Contains(info, "'") {
		// 去除单引号
		info = strings.ReplaceAll(info, "'", "")
	}
	// 匹配没有以/进行分割的情况
	count := strings.Count(info, "/")
	if count == 0 {
		if strings.Contains(info, "HUAWEI Ascend 310") {
			info = info + "/ 8G"
		} else {
			info = info + "/"
		}
	}
	return info
}

func (S SInstanceType) GetGpuSpec() string {

	//  N卡信息记录在 InfoGpuName
	//  HUAWEI卡记录在了 InfoAscendName
	if S.OSExtraSpecs.InfoAscendName != "" {
		S.OSExtraSpecs.InfoGpuName = handleSpecialCharacters(S.OSExtraSpecs.InfoAscendName)
	}

	// 这里需要处理特殊的显卡：
	// NVIDIA Professional Graphics Card
	if strings.Contains(S.OSExtraSpecs.InfoGpuName, "NVIDIA Professional Graphics Card") {
		// nvidia-rtx5000
		return strings.ReplaceAll(strings.ToUpper(S.OSExtraSpecs.QuotaGpu), "-", " ")
	}

	// 定义正则表达式
	// 匹配了第一个 * 后面的任意字符，直到遇到第一个 /。匹配的结果是正则表达式的第一个捕获组 .*?，即 NVIDIA P100。 \s* 用于匹配可能存在的空格。
	re := regexp.MustCompile(`\*\s*(.*?)\s*/`)

	// 在字符串中查找匹配项
	match := re.FindStringSubmatch(S.OSExtraSpecs.InfoGpuName)

	// 输出匹配的结果
	if len(match) > 1 {
		return match[1]
	} else {
		return S.OSExtraSpecs.InfoGpuName
	}
}

func (S SInstanceType) GetGpuCount() string {

	// 即使有 GPU 挂载 S.OSExtraSpecs.QuotaGpu 字段也不一定有值
	// 原方案从 pci_passthrough:alias 获取数量：
	// 		PCI直通设备信息，格式为PCI设备名称:数量。多个设备信息以逗号隔开。
	// 		例如nvidia-a30:1，表示携带一张A30的GPU。
	// 但是无法满足 S.OSExtraSpecs.QuotaGpu 为空时， PciPassthroughGpuSpecs 为 m60_2q:virt:1 的情况

	// 现改为从 InfoGpuName 获取相关参数，以下为情况枚举：
	// 		1 * NVIDIA P100 / 1 * 16G
	// 		1 * NVIDIA M60-2Q / 2G

	//  N卡信息记录在 InfoGpuName
	//  HUAWEI卡记录在了 InfoAscendName
	if S.OSExtraSpecs.InfoAscendName != "" {
		S.OSExtraSpecs.InfoGpuName = handleSpecialCharacters(S.OSExtraSpecs.InfoAscendName)
	}

	count := strings.Count(S.OSExtraSpecs.InfoGpuName, "*")
	switch count {
	case 1:
		// 定义正则表达式模式
		pattern := `(\d+)\s*\*`
		// 编译正则表达式
		regexpPattern := regexp.MustCompile(pattern)
		// 查找匹配项
		matches := regexpPattern.FindStringSubmatch(S.OSExtraSpecs.InfoGpuName)
		// 提取匹配的数字
		if len(matches) > 1 {
			firstNumber := matches[1]
			return firstNumber
		} else {
			return "0"
		}
	case 2:
		// 定义正则表达式模式
		pattern := `^([0-9]+)\s*\*`
		// 编译正则表达式
		regexpPattern := regexp.MustCompile(pattern)
		// 查找匹配项
		matches := regexpPattern.FindStringSubmatch(S.OSExtraSpecs.InfoGpuName)
		// 提取匹配的数字
		if len(matches) > 1 {
			firstNumber := matches[1]
			return firstNumber
		} else {
			return "0"
		}
	default:
		return "0"
	}
}

func (S SInstanceType) GetGpuMaxCount() int {

	// 即使有 GPU 挂载 S.OSExtraSpecs.QuotaGpu 字段也不一定有值
	// 原方案从 pci_passthrough:alias 获取数量：
	// 		PCI直通设备信息，格式为PCI设备名称:数量。多个设备信息以逗号隔开。
	// 		例如nvidia-a30:1，表示携带一张A30的GPU。
	// 但是无法满足 S.OSExtraSpecs.QuotaGpu 为空时， PciPassthroughGpuSpecs 为 m60_2q:virt:1 的情况

	// 现改为从 InfoGpuName 获取相关参数，以下为情况枚举：
	// 		1 * NVIDIA P100 / 1 * 16G
	// 		1 * NVIDIA M60-2Q / 2G

	//  N卡信息记录在 InfoGpuName
	//  HUAWEI卡记录在了 InfoAscendName
	if S.OSExtraSpecs.InfoAscendName != "" {
		S.OSExtraSpecs.InfoGpuName = S.OSExtraSpecs.InfoAscendName
	}

	count := strings.Count(S.OSExtraSpecs.InfoGpuName, "*")
	switch count {
	case 1:
		// 定义正则表达式模式
		pattern := `(\d+)\s*\*`
		// 编译正则表达式
		regexpPattern := regexp.MustCompile(pattern)
		// 查找匹配项
		matches := regexpPattern.FindStringSubmatch(S.OSExtraSpecs.InfoGpuName)
		// 提取匹配的数字
		if len(matches) > 1 {
			firstNumber := matches[1]
			atoi, err := strconv.Atoi(firstNumber)
			if err != nil {
				return 0
			}
			return atoi
		} else {
			return 0
		}
	case 2:
		// 定义正则表达式模式
		pattern := `^([0-9]+)\s*\*`
		// 编译正则表达式
		regexpPattern := regexp.MustCompile(pattern)
		// 查找匹配项
		matches := regexpPattern.FindStringSubmatch(S.OSExtraSpecs.InfoGpuName)
		// 提取匹配的数字
		if len(matches) > 1 {
			firstNumber := matches[1]
			atoi, err := strconv.Atoi(firstNumber)
			if err != nil {
				return 0
			}
			return atoi
		} else {
			return 0
		}
	default:
		return 0
	}
}

func (S SInstanceType) Delete() error {
	return nil
}

// https://support.huaweicloud.com/api-ecs/zh-cn_topic_0020212656.html
func (self *SRegion) fetchInstanceTypes(zoneId string) ([]SInstanceType, error) {
	querys := map[string]string{}
	if len(zoneId) > 0 {
		querys["availability_zone"] = zoneId
	}

	instanceTypes := make([]SInstanceType, 0)
	err := doListAll(self.ecsClient.Flavors.List, querys, &instanceTypes)
	return instanceTypes, err
}

func (self *SRegion) GetMatchInstanceTypes(cpu int, memMB int, zoneId string) ([]SInstanceType, error) {
	instanceTypes, err := self.fetchInstanceTypes(zoneId)
	if err != nil {
		return nil, err
	}

	ret := make([]SInstanceType, 0)
	for _, t := range instanceTypes {
		// cpu & mem & disk都匹配才行
		if t.Vcpus == strconv.Itoa(cpu) && t.RamMB == memMB {
			ret = append(ret, t)
		}
	}

	return ret, nil
}

// GetInstanceTypes 获取 zone 下 InstanceTypes(cloudservers/flavors)
func (self *SRegion) GetInstanceTypes(zoneId string) ([]SInstanceType, error) {
	instanceTypes, err := self.fetchInstanceTypes(zoneId)
	if err != nil {
		return nil, err
	}
	return instanceTypes, nil
}

// GetRegionInstanceTypes 获取 region 下 InstanceTypes(cloudservers/flavors)
func (self *SRegion) GetRegionInstanceTypes() ([]SInstanceType, error) {
	zones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	sInstanceTypes := make([]SInstanceType, 0)
	for _, zone := range zones {
		zoneInstanceTypes, err := self.fetchInstanceTypes(zone.GetId())
		if err != nil {
			continue
		}
		for _, instanceType := range zoneInstanceTypes {
			instanceType.ZoneID = zone.GetId()
			sInstanceTypes = append(sInstanceTypes, instanceType)
		}
	}

	return sInstanceTypes, nil
}

func (self *SRegion) GetISkus() ([]cloudprovider.ICloudSku, error) {
	skus, err := self.GetRegionInstanceTypes()
	if err != nil {
		return nil, errors.Wrapf(err, "GetInstanceTypes")
	}
	ret := []cloudprovider.ICloudSku{}
	for i := range skus {
		ret = append(ret, &skus[i])
	}
	return ret, nil
}

// GetInstanceTypeStatus
// 某个AZ没有在cond:operation:az参数中配置时默认使用此参数的取值。不配置或无此参数时等同于“normal”。取值范围：
//
// normal：正常商用
// abandon：下线（即不显示）
// sellout：售罄
// obt：公测
// obt_sellout：公测售罄
// promotion：推荐(等同 normal，也是商用 )
func (self *SRegion) GetInstanceTypeStatus(zoneId, instanceTypeName string) (string, error) {
	skus, err := self.fetchInstanceTypes(zoneId)
	if err != nil {
		return "", errors.Wrapf(err, "fetchInstanceTypes")
	}
	ret := make([]SInstanceType, 0)
	for _, sku := range skus {
		if sku.Name == instanceTypeName {
			ret = append(ret, sku)
			break
		}
	}
	if len(ret) > 0 {

		var retSkuStatus string
		// like: cn-east-3a(normal),cn-east-3c(sellout),cn-east-3d(normal)
		az := ret[0].OSExtraSpecs.CondOperationAz

		// like: abandon
		status := ret[0].OSExtraSpecs.CondOperationStatus
		// 定义正则表达式
		re := regexp.MustCompile(zoneId + `\(([^)]+)\)`)
		// 在字符串中查找匹配项
		match := re.FindStringSubmatch(az)
		if len(match) >= 2 {
			retSkuStatus = match[1]
		} else {
			retSkuStatus = status
		}
		switch retSkuStatus {
		case "normal", "obt", "promotion":
			return api.SkuStatusAvailable, nil
		case "abandon", "sellout", "obt_sellout":
			return api.SkuStatusSoldout, nil
		default:
			return api.SkuStatusSoldout, nil
		}
	}
	// 没有查询到对应 zone instanceType
	return api.SkuStatusSoldout, nil
}
