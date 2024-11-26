package aws

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/log"
)

type DatapointsCfel struct {
	Datapoints DatapointSlice `xml:"Datapoints>member"`
	Label      string         `xml:"Label"`
}

type DatapointSlice []Datapoint

func (a DatapointSlice) Len() int           { return len(a) }
func (a DatapointSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DatapointSlice) Less(i, j int) bool { return a[i].Timestamp.Before(a[j].Timestamp) }

func (self *SAwsClient) cfelGetMetrics(regionId, ns string, period int, metricType cloudprovider.TMetricType, metrics map[string]string, dimensionName, dimensionValue string, start, end time.Time) ([]cloudprovider.MetricValues, error) {
	result := []cloudprovider.MetricValues{}
	for metricName, tagValue := range metrics {
		params := map[string]string{
			"EndTime":                   end.Format(time.RFC3339),
			"MetricName":                metricName,
			"Dimensions.member.1.Name":  dimensionName,
			"Dimensions.member.1.Value": dimensionValue,
			"Namespace":                 ns,
			"Period":                    strconv.Itoa(period),
			"StartTime":                 start.Format(time.RFC3339),
			"Statistics.member.1":       "Average",
		}
		ret := DatapointsCfel{}
		err := self.monitorRequest(regionId, "GetMetricStatistics", params, &ret)
		if err != nil {
			log.Errorf("GetMetricStatistics error: %v", err)
			continue
		}
		sort.Sort(ret.Datapoints)
		metric := cloudprovider.MetricValues{}
		metric.Id = dimensionValue
		metric.MetricType = metricType
		tags := map[string]string{}
		idx := strings.Index(tagValue, ":")
		if idx > 0 {
			tags[tagValue[:idx]] = tagValue[idx+1:]
		}
		for _, data := range ret.Datapoints {
			metricValue := cloudprovider.MetricValue{}
			metricValue.Tags = tags
			metricValue.Timestamp = data.Timestamp
			metricValue.Value = data.GetValue()
			metric.Values = append(metric.Values, metricValue)
		}
		result = append(result, metric)
	}
	return result, nil
}
