package huawei

import (
	"fmt"
	"testing"
)

func TestGetGpuCount(t *testing.T) {
	instanceType := SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoGpuName: "1 * NVIDIA M60-2Q / 2G"}}
	fmt.Println(instanceType.GetGpuCount())
	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoGpuName: "1 * NVIDIA P100 / 1 * 16G"}}
	fmt.Println(instanceType.GetGpuCount())
	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoAscendName: "4 * HUAWEI Ascend 310/4 * 8G"}}
	fmt.Println(instanceType.GetGpuCount())

	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoAscendName: "2 * HUAWEI Ascend 310"}}
	fmt.Println(instanceType.GetGpuCount())

}

func TestGetGpuSpec(t *testing.T) {
	instanceType := SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoGpuName: "1 * NVIDIA M60-2Q / 2G"}}
	fmt.Println(instanceType.GetGpuSpec())
	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoGpuName: "1 * NVIDIA P100 / 1 * 16G"}}
	fmt.Println(instanceType.GetGpuSpec())
	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoAscendName: "4 * HUAWEI Ascend 310/4 * 8G"}}
	fmt.Println(instanceType.GetGpuSpec())

	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoAscendName: "2 * HUAWEI Ascend 310"}}
	fmt.Println(instanceType.GetGpuSpec())
}

func TestGetGPUMemorySizeMB(t *testing.T) {
	instanceType := SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoGpuName: "1 * NVIDIA M60-2Q / 2G"}}
	fmt.Println(instanceType.GetGPUMemorySizeMB())
	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoGpuName: "1 * NVIDIA P100 / 1 * 16G"}}
	fmt.Println(instanceType.GetGPUMemorySizeMB())
	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoAscendName: "4 * HUAWEI Ascend 310/4 * 8G"}}
	fmt.Println(instanceType.GetGPUMemorySizeMB())
	instanceType = SInstanceType{OSExtraSpecs: OSExtraSpecs{InfoAscendName: "2 * HUAWEI Ascend 310"}}
	fmt.Println(instanceType.GetGPUMemorySizeMB())
}
