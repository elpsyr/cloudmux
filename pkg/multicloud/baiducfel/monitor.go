package baiducfel

import (
	"fmt"
	"strconv"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

func (self *SBaiduClient) GetMetrics(opts *cloudprovider.MetricListOptions) ([]cloudprovider.MetricValues, error) {
	switch opts.ResourceType {
	case cloudprovider.METRIC_RESOURCE_TYPE_SERVER:
		return self.GetEcsMetrics(opts)
	default:
		return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "%s", opts.ResourceType)
	}
}

func (self *SBaiduClient) GetEcsMetrics(opts *cloudprovider.MetricListOptions) ([]cloudprovider.MetricValues, error) {
	// metricTags, tagKey := map[string]string{}, ""

	ret := []cloudprovider.MetricValues{}
	for metricType, metricNames := range map[cloudprovider.TMetricType]map[string]string{
		cloudprovider.VM_METRIC_TYPE_CPU_USAGE: {
			"CpuIdlePercent": "", // cpu 空闲率
		},
		cloudprovider.VM_METRIC_TYPE_MEM_USAGE: {
			"MemUsedPercent": "",
		},
		cloudprovider.VM_METRIC_TYPE_NET_BPS_TX: {
			"WebOutBitsPerSecond": cloudprovider.METRIC_TAG_NET_TYPE + ":" + cloudprovider.METRIC_TAG_NET_TYPE_INTRANET,
			// "WanOuttraffic": cloudprovider.METRIC_TAG_NET_TYPE + ":" + cloudprovider.METRIC_TAG_NET_TYPE_INTERNET,
		},
		cloudprovider.VM_METRIC_TYPE_NET_BPS_RX: {
			"WebInBitsPerSecond": cloudprovider.METRIC_TAG_NET_TYPE + ":" + cloudprovider.METRIC_TAG_NET_TYPE_INTRANET,
			// "WanIntraffic": cloudprovider.METRIC_TAG_NET_TYPE + ":" + cloudprovider.METRIC_TAG_NET_TYPE_INTERNET,
		},
	} {
		for metricName := range metricNames {
			result, err := self.ListMetrics(opts.RegionExtId, opts.ResourceId, metricName, opts.Interval, opts.StartTime, opts.EndTime)
			if err != nil {
				// log.Errorf("ListMetric(%s) error: %v", metric, err)
				continue
			}

			var values = []cloudprovider.MetricValue{}
			for i := range result {
				// dataTag := result[i].GetTags()
				var val = result[i].GetValue()
				if metricType == "CpuIdlePercent" {
					val = 1 - val
				}
				t, _ := time.Parse("2006-01-02T15:04:05Z", result[i].Timestamp)
				values = append(values, cloudprovider.MetricValue{
					Timestamp: t,
					Value:     val,
					// Tags:      map[string]string{},
				})
			}
			var metricValue = cloudprovider.MetricValues{
				Id:         opts.ResourceId,
				Unit:       "",
				MetricType: metricType,
				Values:     values,
			}
			ret = append(ret, metricValue)
		}
	}
	return ret, nil
}

func (self *SBaiduClient) ListMetrics(region, instanceId string, metricName string, interval int, start, end time.Time) ([]MetricData, error) {
	var result = []MetricData{}
	var query = map[string]string{
		"dimensions":     fmt.Sprintf("InstanceId:%s", instanceId),
		"periodInSecond": strconv.Itoa(interval),
		"startTime":      start.Format("2006-01-02T15:04:05Z"),
		"endTime":        end.Format("2006-01-02T15:04:05Z"),
		"statistics[]":   "average",
	}
	resp, err := self.list("bcm", region, fmt.Sprintf("/json-api/v1/metricdata/%s/BCE_BCC/%s", self.ownerId, metricName), query, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeMetricList")
	}
	var ret = respMetrics{}

	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, errors.Wrapf(err, "resp.Unmarshal")
	}
	for _, val := range ret.DataPoints {
		result = append(result, MetricData{
			InstanceId: instanceId,
			Timestamp:  val.Timestamp,
			Average:    val.Average,
		})
	}
	return result, nil
}

type MetricData struct {
	InstanceId string
	Timestamp  string
	Average    float64
}

func (d MetricData) GetValue() float64 {
	return d.Average
}

func (d MetricData) GetTags() map[string]string {
	ret := map[string]string{}
	return ret
}

type respMetrics struct {
	Code       string       `json:"code"`
	DataPoints []dataPoints `json:"dataPoints"`
	Message    string       `json:"message"`
	RequestID  string       `json:"requestId"`
}

type dataPoints struct {
	Average   float64 `json:"average"`
	Timestamp string  `json:"timestamp"`
}
