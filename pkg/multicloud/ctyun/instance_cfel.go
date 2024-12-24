package ctyun

import (
	"context"
	"fmt"
	"strings"
)

func (self *SInstance) GetIHostId() string {
	if self.region != nil {
		return fmt.Sprintf("-%s/%s", self.region.GetGlobalId(), self.AzName)
	} else if self.host != nil {
		return self.host.GetGlobalId()
	}
	return ""
}

func (s *SInstance) RebootVM(ctx context.Context) error {
	params := map[string]interface{}{
		"instanceID": s.InstanceId,
		"regionID":   s.host.zone.region.RegionId,
	}
	_, err := s.host.zone.region.post(SERVICE_ECS, "/v4/ecs/reboot-instance", params)
	return err
}

func (s *SInstance) GetVncUrl() (string, error) {
	// https://www.ctyun.cn/document/10026730/10597652 单可用区和多可用区处理方式不同
	url, err := s.region.GetInstanceVnc(s.InstanceId)
	if s.region.IsMultiZones {
		arr := strings.Split(url, "?")
		if len(arr) < 2 {
			return "", fmt.Errorf("vnc url format error")
		}
		os := "Linux"
		if strings.Contains(strings.ToLower(s.Image.ImageName), "windows") {
			os = "Windows"
		}
		domain := arr[0][0:strings.LastIndex(arr[0], "/")]
		u := "https://console.ctyun.cn/compute/index/#/ecm/ecmtelnet"
		url = fmt.Sprintf("%s?%s&domainName=%s&os=%s", u, arr[1], domain, os)
	}
	return url, err
}
