package aliyun

import (
	"context"
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
)

func (self *SInstance) RebootVM(ctx context.Context) error {
	err := self.host.zone.region.RebootVM(self.InstanceId)
	if err != nil {
		return err
	}
	return cloudprovider.WaitStatus(self, api.VM_RUNNING, 10*time.Second, 300*time.Second) // 5mintues
}

func (self *SRegion) doRebootVM(instanceId string) error {
	return self.instanceOperation(instanceId, "RebootInstance", nil)
}

func (self *SRegion) RebootVM(instanceId string) error {
	status, err := self.GetInstanceStatus(instanceId)
	if err != nil {
		log.Errorf("Fail to get instance status on StartVM: %s", err)
		return err
	}
	if status != InstanceStatusRunning {
		log.Errorf("RebootVM: vm status is %s expect %s", status, InstanceStatusRunning)
		return cloudprovider.ErrInvalidStatus
	}
	return self.doRebootVM(instanceId)
	// if err != nil {
	//	return err
	// }
	// return self.waitInstanceStatus(instanceId, InstanceStatusRunning, time.Second*5, time.Second*180) // 3 minutes to timeout
}

func (self *SInstance) GetMonitorData(start, end string) ([]cloudprovider.ICfelMonitorData, error) {
	data, err := self.host.zone.region.DescribeInstanceMonitorData(self.InstanceId, start, end, "")
	if err != nil {
		return nil, errors.Wrap(err, "DescribeInstanceMonitorData")
	}
	// 将 []MonitorDataItem 转换成 []cloudprovider.MonitorData
	var providerData []cloudprovider.ICfelMonitorData
	for _, item := range data {
		providerData = append(providerData, cloudprovider.ICfelMonitorData(item))
	}
	return providerData, nil
}

func (self *SRegion) DescribeInstanceMonitorData(instanceId, startTime, endTime, period string) ([]MonitorDataItem, error) {
	params := map[string]string{
		"InstanceId": instanceId,
		"StartTime":  startTime,
		"EndTime":    endTime,
		"Period":     period,
	}
	if period == "" {
		params["Period"] = "600" // 默认值：60
	}
	resp, err := self.ecsRequest("DescribeInstanceMonitorData", params)
	if err != nil {
		return nil, errors.Wrapf(err, "DescribeInstanceMonitorData")
	}
	ret := DescribeInstanceMonitorDataResponse{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, errors.Wrapf(err, "resp.Unmarshal")
	}
	return ret.MonitorData.InstanceMonitorData, nil
}

// DescribeInstanceMonitorDataResponse
// DescribeInstanceMonitorData 接口返回数据结构
type DescribeInstanceMonitorDataResponse struct {
	RequestId   string `json:"RequestId"`
	MonitorData struct {
		InstanceMonitorData []MonitorDataItem `json:"InstanceMonitorData"`
	} `json:"MonitorData"`
}

type MonitorDataItem struct {
	IOPSRead          float64    `json:"IOPSRead,omitempty"`
	IntranetBandwidth float64    `json:"IntranetBandwidth"`
	IOPSWrite         float64    `json:"IOPSWrite,omitempty"`
	InstanceId        string `json:"InstanceId"`
	IntranetTX        float64    `json:"IntranetTX"`
	CPU               float64    `json:"CPU"`
	BPSRead           float64    `json:"BPSRead,omitempty"`
	IntranetRX        float64    `json:"IntranetRX"`
	TimeStamp         string `json:"TimeStamp"`
	InternetBandwidth float64    `json:"InternetBandwidth"`
	InternetTX        float64    `json:"InternetTX"`
	InternetRX        float64    `json:"InternetRX"`
	BPSWrite          float64    `json:"BPSWrite,omitempty"`
}

var _ cloudprovider.ICfelMonitorData = (*MonitorDataItem)(nil)

func (m MonitorDataItem) GetBPSRead() float64 {
	return m.IOPSRead
}

func (m MonitorDataItem) GetInternetTX() float64 {
	return m.InternetTX
}

func (m MonitorDataItem) GetCPU() float64 {
	return m.CPU
}

func (m MonitorDataItem) GetMem() float64 {
	return 0
}

func (m MonitorDataItem) GetDisk() float64 {
	return 0
}

func (m MonitorDataItem) GetIOPSWrite() float64 {
	return m.IOPSWrite
}

func (m MonitorDataItem) GetIntranetTX() float64 {
	return m.IntranetTX
}

func (m MonitorDataItem) GetInstanceId() string {
	return m.InstanceId
}

func (m MonitorDataItem) GetBPSWrite() float64 {
	return m.BPSWrite
}

func (m MonitorDataItem) GetIOPSRead() float64 {
	return m.IOPSRead
}

func (m MonitorDataItem) GetInternetBandwidth() float64 {
	return m.InternetBandwidth
}

func (m MonitorDataItem) GetInternetRX() float64 {
	return m.InternetRX
}

func (m MonitorDataItem) GetTimeStamp() string {
	return m.TimeStamp
}

func (m MonitorDataItem) GetIntranetRX() float64 {
	return m.IntranetRX
}

func (m MonitorDataItem) GetIntranetBandwidth() float64 {
	return m.IntranetBandwidth
}
