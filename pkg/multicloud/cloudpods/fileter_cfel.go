package cloudpods

import "yunion.io/x/cloudmux/pkg/cloudprovider"

type SCfelFilter struct {
	FilterName string   `json:"filter_name,omitempty"`
	Id         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Reason     []string `json:"reason,omitempty"`
}

func (s *SCfelFilter) GetFilterName() string {
	return s.FilterName
}

func (s *SCfelFilter) GetID() string {
	return s.Id
}

func (s *SCfelFilter) GetName() string {
	return s.Name
}

func (s *SCfelFilter) GetReason() []string {
	return s.Reason
}

var _ cloudprovider.ICfelFilter = (*SCfelFilter)(nil)

type MonitorDataItem struct {
	IOPSRead          float64 `json:"IOPSRead,omitempty"`
	IntranetBandwidth float64 `json:"IntranetBandwidth"`
	IOPSWrite         float64 `json:"IOPSWrite,omitempty"`
	InstanceId        string  `json:"InstanceId"`
	IntranetTX        float64 `json:"IntranetTX"`
	CPU               float64 `json:"CPU"`
	Mem               float64 `json:"men"`
	Disk              float64 `json:"disk"`
	BPSRead           float64 `json:"BPSRead,omitempty"`
	IntranetRX        float64 `json:"IntranetRX"`
	TimeStamp         string  `json:"TimeStamp"`
	InternetBandwidth float64 `json:"InternetBandwidth"`
	InternetTX        float64 `json:"InternetTX"`
	InternetRX        float64 `json:"InternetRX"`
	BPSWrite          float64 `json:"BPSWrite,omitempty"`
}

func (m *MonitorDataItem) GetBPSRead() float64 {
	return m.BPSRead
}

func (m *MonitorDataItem) GetInternetTX() float64 {
	return m.InternetTX
}

func (m *MonitorDataItem) GetCPU() float64 {
	return m.CPU
}

func (m *MonitorDataItem) GetMem() float64 {
	return m.Mem
}

func (m *MonitorDataItem) GetDisk() float64 {
	return m.Disk
}

func (m *MonitorDataItem) GetIOPSWrite() float64 {
	return m.IOPSWrite
}

func (m *MonitorDataItem) GetIntranetTX() float64 {
	return m.IntranetTX
}

func (m *MonitorDataItem) GetInstanceId() string {
	return m.InstanceId
}

func (m *MonitorDataItem) GetBPSWrite() float64 {
	return m.BPSWrite
}

func (m *MonitorDataItem) GetIOPSRead() float64 {
	return m.IOPSRead
}

func (m *MonitorDataItem) GetInternetBandwidth() float64 {
	return m.InternetBandwidth
}

func (m *MonitorDataItem) GetInternetRX() float64 {
	return m.InternetRX
}

func (m *MonitorDataItem) GetTimeStamp() string {
	return m.TimeStamp
}

func (m *MonitorDataItem) GetIntranetRX() float64 {
	return m.IntranetRX
}

func (m *MonitorDataItem) GetIntranetBandwidth() float64 {
	return m.IntranetBandwidth
}

var _ cloudprovider.ICfelMonitorData = (*MonitorDataItem)(nil)
